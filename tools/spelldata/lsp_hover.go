package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"
)

// Each pattern captures the id in group 1.
var spellIDPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\bMustFind\(\s*(\d+)\s*\)`),
	regexp.MustCompile(`\bFind\(\s*(\d+)\s*\)`),
	regexp.MustCompile(`\.ByID\(\s*(\d+)\s*\)`),
	regexp.MustCompile(`\bSpellID:\s*(\d+)`),
	regexp.MustCompile(`"spellId":\s*(\d+)`),
	regexp.MustCompile(`\bfromSpellId\(\s*(\d+)\s*\)`),
	regexp.MustCompile(`\bspellId:\s*(\d+)`),
}

const maxSubstitutions = 4

// A name a Go file binds to a chain, with `var`, `=` or `:=`. A name bound inside a function is seen
// from its binding to the end of that function's body, as offsets in its file; to is 0 on a name
// bound at package level, which every file of the package sees.
type declaration struct {
	name     string
	chain    *chain
	file     string
	line     int
	from, to int
}

// A Go file as the hover reads it. The AST is partial where the text does not parse, and nil where it
// states no package clause.
type parsedFile struct {
	path  string
	text  string
	file  *ast.File
	tok   *token.File
	decls []declaration
	// The names the file binds to a value that is not a chain, with why.
	unread map[*ast.Ident]error
}

func parseGo(path, text string) *parsedFile {
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, path, text, parser.SkipObjectResolution)
	f := &parsedFile{path: path, text: text, unread: map[*ast.Ident]error{}}
	// ParseFile answers an empty file, not nil, where the text states no package clause.
	if file == nil || !file.Package.IsValid() {
		return f
	}
	f.file, f.tok = file, fset.File(file.Pos())

	var around []ast.Node
	bind := func(name ast.Expr, value ast.Expr) {
		ident, ok := name.(*ast.Ident)
		if !ok || ident.Name == "_" {
			return
		}
		c, err := walkChain(value, f.tok.Offset)
		if err != nil {
			f.unread[ident] = err
			return
		}
		d := declaration{name: ident.Name, chain: c, file: path, line: f.tok.Line(ident.Pos())}
		if body := enclosingBody(around); body != nil {
			d.from, d.to = f.tok.Offset(ident.Pos()), f.tok.Offset(body.End())
		}
		f.decls = append(f.decls, d)
	}
	ast.Inspect(file, func(node ast.Node) bool {
		if node == nil {
			around = around[:len(around)-1]
			return true
		}
		around = append(around, node)
		switch n := node.(type) {
		case *ast.ValueSpec:
			if len(n.Names) == len(n.Values) {
				for i, name := range n.Names {
					bind(name, n.Values[i])
				}
			}
		case *ast.AssignStmt:
			if (n.Tok == token.DEFINE || n.Tok == token.ASSIGN) && len(n.Lhs) == len(n.Rhs) {
				for i, name := range n.Lhs {
					bind(name, n.Rhs[i])
				}
			}
		}
		return true
	})
	return f
}

func enclosingBody(nodes []ast.Node) *ast.BlockStmt {
	for i := len(nodes) - 1; i >= 0; i-- {
		switch fn := nodes[i].(type) {
		case *ast.FuncDecl:
			return fn.Body
		case *ast.FuncLit:
			return fn.Body
		}
	}
	return nil
}

// The nodes that hold pos, outermost first. An end counts as inside, so the cursor just past a name
// is on it.
func (f *parsedFile) enclosing(pos token.Pos) []ast.Node {
	var out []ast.Node
	ast.Inspect(f.file, func(node ast.Node) bool {
		if node == nil || pos < node.Pos() || pos > node.End() {
			return false
		}
		out = append(out, node)
		return true
	})
	return out
}

type workspace struct {
	buffers map[string]string
	folders map[string]map[string][]declaration
	// The files of a cached folder that changed since it was read.
	dirty map[string]bool
	// The modification time of each file as it was read from disk, which is how a change the editor
	// does not report - a checkout, a generator, another session - reaches the cache.
	stamps  map[string]time.Time
	current *parsedFile
}

func newWorkspace() *workspace {
	return &workspace{
		buffers: map[string]string{},
		folders: map[string]map[string][]declaration{},
		dirty:   map[string]bool{},
		stamps:  map[string]time.Time{},
	}
}

// The text an editor holds for a file, or with open false, the file as it stands on disk again.
func (w *workspace) update(uri, text string, open bool) {
	path := uriPath(uri)
	if open {
		w.buffers[path] = text
	} else {
		delete(w.buffers, path)
	}
	if _, cached := w.folders[filepath.Dir(path)]; cached {
		w.dirty[path] = true
	}
}

// The file being hovered, parsed once per text.
func (w *workspace) parse(path, text string) *parsedFile {
	if w.current == nil || w.current.path != path || w.current.text != text {
		w.current = parseGo(path, text)
	}
	return w.current
}

// The declarations of a file as the editor holds it or as it stands on disk; false where it is neither.
func (w *workspace) read(path string) ([]declaration, bool) {
	text, open := w.buffers[path]
	if !open {
		info, err := os.Stat(path)
		if err != nil {
			return nil, false
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, false
		}
		w.stamps[path] = info.ModTime()
		text = string(data)
	}
	return parseGo(path, text).decls, true
}

// Every name the package folder binds at package level, and the names the current file binds in the
// function around the offset at, the last binding before at winning. With no current file there is
// no position to scope by, and every binding of the folder is read.
func (w *workspace) declarations(folder string, current *parsedFile, at int, trace *tracer) map[string]declaration {
	files, cached := w.folders[folder]
	if !cached {
		files = map[string][]declaration{}
		w.folders[folder] = files
	}
	w.markChanged(folder, files, trace)
	read := 0
	for path := range w.dirty {
		if filepath.Dir(path) != folder {
			continue
		}
		delete(w.dirty, path)
		read++
		if decls, ok := w.read(path); ok {
			files[path] = decls
		} else {
			delete(files, path)
		}
	}
	if cached && read == 0 {
		trace.add("declarations %s: cache hit", trace.rel(folder))
	} else {
		trace.add("declarations %s: cache miss, read %d files", trace.rel(folder), read)
	}

	own := current.pathIn(folder)
	paths := slices.Collect(maps.Keys(files))
	if _, listed := files[own]; own != "" && !listed {
		paths = append(paths, own)
	}
	slices.Sort(paths)

	out := map[string]declaration{}
	for _, path := range paths {
		decls := files[path]
		if path == own {
			decls = current.decls
		}
		for _, d := range decls {
			if d.to == 0 || current == nil {
				out[d.name] = d
			}
		}
	}
	if own != "" {
		for _, d := range current.decls {
			if d.to != 0 && d.from <= at && at <= d.to {
				out[d.name] = d
			}
		}
	}
	return out
}

// A file the editor does not hold is read again when its modification time moved or it is gone.
func (w *workspace) markChanged(folder string, files map[string][]declaration, trace *tracer) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		trace.add("declarations %s: %v", trace.rel(folder), err)
	}
	onDisk := map[string]bool{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		path := filepath.Join(folder, entry.Name())
		onDisk[path] = true
		_, read := files[path]
		_, open := w.buffers[path]
		if info, err := entry.Info(); !read || !open && (err != nil || !info.ModTime().Equal(w.stamps[path])) {
			w.dirty[path] = true
		}
	}
	for path := range w.buffers {
		if _, read := files[path]; !read && filepath.Dir(path) == folder && strings.HasSuffix(path, ".go") {
			w.dirty[path] = true
		}
	}
	for path := range files {
		if _, open := w.buffers[path]; !open && !onDisk[path] {
			w.dirty[path] = true
		}
	}
}

func (f *parsedFile) pathIn(folder string) string {
	if f == nil || filepath.Dir(f.path) != folder {
		return ""
	}
	return f.path
}

// What a chain hover evaluates, and the name it answers for where the cursor was on a bound name
// rather than on a part of a chain.
type chainHover struct {
	chain *chain
	name  string
}

type tracer struct {
	root  string
	lines []string
}

func (t *tracer) add(format string, args ...any) {
	t.lines = append(t.lines, fmt.Sprintf(format, args...))
}

func (t *tracer) rel(path string) string {
	if t.root != "" {
		if rel, err := filepath.Rel(t.root, path); err == nil && !strings.HasPrefix(rel, "..") {
			return rel
		}
	}
	return path
}

func (t *tracer) fail(format string, args ...any) (string, []string, bool) {
	t.add("✗ "+format, args...)
	return "", t.lines, false
}

// The hover for a position in a document, as markdown, and the lines saying how it was reached. The
// line and column are 0-based and the column counts UTF-16 code units, as LSP positions do.
func (w *workspace) hover(text string, line, col int, uri string) (string, []string, bool) {
	path := uriPath(uri)
	trace := &tracer{}
	trace.root, _ = moduleRoot()

	start, end, ok := lineBounds(text, line)
	if !ok {
		return trace.fail("%s:%d is past the end of the document", trace.rel(path), line+1)
	}
	lineText := strings.TrimSuffix(text[start:end], "\r")
	column := byteOffsetOfUTF16Column(lineText, col)
	trace.add("%s:%d:%d", trace.rel(path), line+1, col+1)

	if match := matchCovering(lineText, column, spellIDPatterns...); match != nil {
		id, _ := strconv.ParseInt(lineText[match[2]:match[3]], 10, 32)
		trace.add("id %d", id)
		s, err := findSpell(int32(id))
		if err != nil {
			return trace.fail("%v", err)
		}
		trace.add("✓ %s", title(s))
		return idMarkdown(s), trace.lines, true
	}

	if filepath.Ext(path) != ".go" {
		return trace.fail("no spell id under the cursor")
	}
	f := w.parse(path, text)
	if f.file == nil {
		return trace.fail("no spell id, family or name under the cursor")
	}
	folder := filepath.Dir(path)
	pkg := filepath.Base(folder)
	at := start + column
	nodes := f.enclosing(f.tok.Pos(at))

	if call := spellConfigAt(nodes, f.tok.Pos(at)); call != nil {
		trace.add("SpellConfig")
		result, err := evalSpellConfig(call, w.declarations(folder, f, at, trace), pkg, trace)
		if err != nil {
			return trace.fail("%v", err)
		}
		trace.add("✓ %s → %d fields", title(result.spell), len(result.rows))
		return configMarkdown(result), trace.lines, true
	}

	if field := familyAt(nodes); field != "" {
		trace.add("family %s/%s", pkg, field)
		family, err := findFamily(ladderFamilies(), field, pkg)
		if err != nil {
			return trace.fail("%v", err)
		}
		if family.err != nil {
			return trace.fail("%v", family.err)
		}
		trace.add("✓ %s, %d ranks", family.key(), family.ladder.Len())
		return familyMarkdown(family), trace.lines, true
	}

	declarations := w.declarations(folder, f, at, trace)
	hover, ok := chainHoverAt(f, nodes, at, declarations, trace.rel(folder), trace)
	if !ok {
		return "", trace.lines, false
	}

	trace.add("expr %s", hover.chain.text(false))
	result, err := evalExpr(ladderFamilies(), hover.chain, pkg)
	if err != nil {
		return trace.fail("%v", err)
	}
	trace.add("✓ %s", resultSummary(result))
	return exprMarkdown(result, hover), trace.lines, true
}

func spellConfigAt(nodes []ast.Node, pos token.Pos) *ast.CallExpr {
	for _, node := range nodes {
		call, ok := node.(*ast.CallExpr)
		if ok && isSelector(call.Fun, "spelldata", "SpellConfig") && pos >= call.Fun.Pos() && pos <= call.Fun.End() {
			return call
		}
	}
	return nil
}

func familyAt(nodes []ast.Node) string {
	for _, node := range nodes {
		if sel, ok := node.(*ast.SelectorExpr); ok && isSelector(sel, "spellData", sel.Sel.Name) {
			return sel.Sel.Name
		}
	}
	return ""
}

func isSelector(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	return ok && ident.Name == pkg && sel.Sel.Name == name
}

func resultSummary(result *exprResult) string {
	parts := []string{result.card().heading()}
	if result.readEffect > 0 {
		parts = append(parts, fmt.Sprintf("effect %d", result.readEffect))
	}
	if result.kind == kindValue {
		parts = append(parts, result.value)
	}
	return strings.Join(parts, " → ")
}

// The outermost chain under the cursor answers for the accessor the cursor is on, arguments included,
// and a name on its own answers for what the package binds it to.
func chainHoverAt(f *parsedFile, nodes []ast.Node, at int, declarations map[string]declaration, folder string, trace *tracer) (chainHover, bool) {
	for _, node := range nodes {
		expr, ok := node.(ast.Expr)
		if !ok {
			continue
		}
		c, err := walkChain(expr, f.tok.Offset)
		if err != nil {
			continue
		}
		if row := c.root.byID; row != nil && row.covers(at) {
			trace.add("segment %s", row.text(false))
			return resolved(&chain{root: c.root}, "", declarations, trace)
		}
		for i, seg := range c.segments {
			if seg.covers(at) {
				trace.add("segment %s", seg.text(false))
				return resolved(&chain{root: c.root, segments: c.segments[:i+1]}, "", declarations, trace)
			}
		}
		if c.root.name == "" || len(c.segments) == 0 {
			break
		}
		trace.add("ident %s", c.root.name)
		return resolved(&chain{root: c.root}, "", declarations, trace)
	}

	var ident *ast.Ident
	for _, node := range nodes {
		if name, ok := node.(*ast.Ident); ok {
			ident = name
		}
	}
	if ident == nil {
		trace.add("✗ no spell id, family or name under the cursor")
		return chainHover{}, false
	}
	trace.add("ident %s", ident.Name)
	if _, bound := declarations[ident.Name]; bound {
		return resolved(&chain{root: root{name: ident.Name}}, ident.Name, declarations, trace)
	}
	if err, ok := f.unread[ident]; ok {
		trace.add("✗ %s is bound to no chain the evaluator reads: %v", ident.Name, err)
	} else {
		trace.add("✗ no ladder-shaped declaration of %s in %s", ident.Name, folder)
	}
	return chainHover{}, false
}

func resolved(c *chain, name string, declarations map[string]declaration, trace *tracer) (chainHover, bool) {
	resolved, err := resolveChain(c, declarations, trace)
	if err != nil {
		trace.add("✗ %v", err)
		return chainHover{}, false
	}
	return chainHover{chain: resolved, name: name}, true
}

// The chain with every name the package binds substituted by the chain it stands for, down to a
// ladder: at most maxSubstitutions names deep, and never through a name twice.
func resolveChain(c *chain, declarations map[string]declaration, trace *tracer) (*chain, error) {
	return resolveFrom(c, declarations, maxSubstitutions, map[string]bool{}, trace)
}

func resolveFrom(c *chain, declarations map[string]declaration, depth int, seen map[string]bool, trace *tracer) (*chain, error) {
	if c.rooted() {
		return c, nil
	}

	name := c.root.name
	bound, ok := declarations[name]
	if !ok {
		if isFamilyName(name) {
			return &chain{root: root{family: name}, segments: c.segments}, nil
		}
		return nil, fmt.Errorf("no ladder-shaped declaration of %s in the package", name)
	}
	if seen[name] {
		return nil, fmt.Errorf("%s stands on itself", name)
	}
	if depth <= 0 {
		return nil, fmt.Errorf("%s stands on more than %d names", c.text(false), maxSubstitutions)
	}

	seen[name] = true
	trace.add("  %s = %s  (%s:%d)", name, bound.chain.text(false), trace.rel(bound.file), bound.line)
	head, err := resolveFrom(bound.chain, declarations, depth-1, seen, trace)
	if err != nil {
		return nil, err
	}
	return &chain{root: head.root, segments: append(slices.Clip(head.segments), c.segments...)}, nil
}

// The byte offsets of a line, counted from 0, without its newline.
func lineBounds(text string, line int) (int, int, bool) {
	if line < 0 {
		return 0, 0, false
	}
	start := 0
	for ; line > 0; line-- {
		next := strings.IndexByte(text[start:], '\n')
		if next < 0 {
			return 0, 0, false
		}
		start += next + 1
	}
	end := strings.IndexByte(text[start:], '\n')
	if end < 0 {
		return start, len(text), true
	}
	return start, start + end, true
}

// The byte offset of an LSP position, clamped to the end of its line and of the text.
func offsetOf(text string, line, col int) int {
	start, end, ok := lineBounds(text, line)
	if !ok {
		return len(text)
	}
	return start + byteOffsetOfUTF16Column(strings.TrimSuffix(text[start:end], "\r"), col)
}

func matchCovering(lineText string, column int, patterns ...*regexp.Regexp) []int {
	for _, pattern := range patterns {
		for _, match := range pattern.FindAllStringSubmatchIndex(lineText, -1) {
			if column >= match[0] && column <= match[1] {
				return match
			}
		}
	}
	return nil
}

func byteOffsetOfUTF16Column(lineText string, units int) int {
	offset := 0
	for offset < len(lineText) && units > 0 {
		r, size := utf8.DecodeRuneInString(lineText[offset:])
		units -= utf16.RuneLen(r)
		offset += size
	}
	return offset
}

// A Windows URI states the drive after the path's leading slash, file:///c:/x.
func uriPath(uri string) string {
	if u, err := url.Parse(uri); err == nil && u.Scheme == "file" {
		path := u.Path
		if filepath.VolumeName(strings.TrimPrefix(path, "/")) != "" {
			path = strings.TrimPrefix(path, "/")
		}
		return filepath.Clean(filepath.FromSlash(path))
	}
	return filepath.Clean(uri)
}

func pathURI(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	abs = filepath.ToSlash(abs)
	if !strings.HasPrefix(abs, "/") {
		abs = "/" + abs
	}
	return (&url.URL{Scheme: "file", Path: abs}).String()
}
