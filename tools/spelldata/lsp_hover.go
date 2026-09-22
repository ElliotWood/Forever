package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/wowsims/forever/sim/core/spelldata"
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

var familyPattern = regexp.MustCompile(`\bspellData\.([A-Za-z_]\w*)`)

var identPattern = regexp.MustCompile(`[A-Za-z_]\w*`)

// `name = value`, `var name = value` or `name := value`; a value starting with `=` is the right half of
// a comparison.
var declarationPattern = regexp.MustCompile(`^(\s*(?:var\s+)?)([A-Za-z_]\w*)(\s*(?::=|=)\s*)(\S.*)$`)

var goKeywords = map[string]bool{}

func init() {
	for _, keyword := range strings.Fields("break case chan const continue default defer else fallthrough " +
		"for func go goto if import interface map package range return select struct switch type var") {
		goKeywords[keyword] = true
	}
}

const maxSubstitutions = 4

type declaration struct {
	name      string
	nameStart int
	nameEnd   int
	chain     *chain
	file      string
	line      int
}

func declarationOnLine(lineText string) (*declaration, error) {
	if strings.HasPrefix(strings.TrimSpace(lineText), "//") {
		return nil, nil
	}
	code, _, _ := strings.Cut(lineText, "//")
	match := declarationPattern.FindStringSubmatchIndex(code)
	if match == nil || strings.HasPrefix(code[match[8]:], "=") {
		return nil, nil
	}
	if goKeywords[code[match[4]:match[5]]] {
		return nil, nil
	}

	d := &declaration{name: code[match[4]:match[5]], nameStart: match[4], nameEnd: match[5]}
	c, err := parseChain(strings.TrimRight(code[match[8]:match[9]], " \t"), match[8])
	if err != nil {
		return d, err
	}
	d.chain = c
	return d, nil
}

func declarationsIn(text, file string) []declaration {
	var out []declaration
	for i, line := range strings.Split(text, "\n") {
		if d, err := declarationOnLine(line); err == nil && d != nil {
			d.file, d.line = file, i+1
			out = append(out, *d)
		}
	}
	return out
}

type workspace struct {
	mu      sync.Mutex
	buffers map[string]string
	folders map[string]map[string][]declaration
}

func newWorkspace() *workspace {
	return &workspace{buffers: map[string]string{}, folders: map[string]map[string][]declaration{}}
}

var defaultWorkspace = newWorkspace()

// The hover for a position in a document, as markdown, and the lines saying how it was reached. The
// line and column are 0-based and the column counts UTF-16 code units, as LSP positions do.
func Hover(text string, line, col int, uri string) (string, []string, bool) {
	return defaultWorkspace.hover(text, line, col, uri)
}

func (w *workspace) setBuffer(uri, text string) {
	path := uriPath(uri)
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buffers[path] = text
	delete(w.folders, filepath.Dir(path))
}

func (w *workspace) dropBuffer(uri string) {
	path := uriPath(uri)
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.buffers, path)
	delete(w.folders, filepath.Dir(path))
}

func (w *workspace) invalidate(uri string) {
	path := uriPath(uri)
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.folders, filepath.Dir(path))
}

func (w *workspace) buffer(uri string) (string, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	text, ok := w.buffers[uriPath(uri)]
	return text, ok
}

func (w *workspace) declarations(folder, current, text string, trace *tracer) map[string]declaration {
	files, hit := w.folders[folder]
	if hit {
		trace.add("declarations %s: cache hit", trace.rel(folder))
	} else {
		files = map[string][]declaration{}
		entries, err := os.ReadDir(folder)
		if err != nil {
			trace.add("declarations %s: %v", trace.rel(folder), err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
				continue
			}
			path := filepath.Join(folder, entry.Name())
			source, open := w.buffers[path]
			if !open {
				data, err := os.ReadFile(path)
				if err != nil {
					continue
				}
				source = string(data)
			}
			files[path] = declarationsIn(source, path)
		}
		for path, source := range w.buffers {
			if _, read := files[path]; !read && filepath.Dir(path) == folder && strings.HasSuffix(path, ".go") {
				files[path] = declarationsIn(source, path)
			}
		}
		w.folders[folder] = files
		trace.add("declarations %s: cache miss, read %d files", trace.rel(folder), len(files))
	}

	paths := make([]string, 0, len(files)+1)
	for path := range files {
		paths = append(paths, path)
	}
	if _, listed := files[current]; !listed {
		paths = append(paths, current)
	}
	sort.Strings(paths)

	out := map[string]declaration{}
	for _, path := range paths {
		decls := files[path]
		if path == current {
			decls = declarationsIn(text, path)
		}
		for _, d := range decls {
			out[d.name] = d
		}
	}
	return out
}

type chainHover struct {
	label   string
	chain   *chain
	segment bool
}

type tracer struct {
	root  string
	here  string
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

func (w *workspace) hover(text string, line, col int, uri string) (string, []string, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	path := uriPath(uri)
	trace := &tracer{}
	trace.root, _ = moduleRoot()

	lines := strings.Split(text, "\n")
	if line < 0 || line >= len(lines) {
		return trace.fail("%s:%d is past the end of the document", trace.rel(path), line+1)
	}
	lineText := strings.TrimSuffix(lines[line], "\r")
	column := byteOffsetOfUTF16Column(lineText, col)
	trace.here = fmt.Sprintf("%s:%d", trace.rel(path), line+1)
	trace.add("%s:%d", trace.here, col+1)

	if match := matchCovering(spellIDPatterns, lineText, column); match != nil {
		id, _ := strconv.ParseInt(lineText[match[2]:match[3]], 10, 32)
		trace.add("id %d", id)
		s := spelldata.Find(int32(id))
		if s == spelldata.Nil {
			return trace.fail("spell %d is not in the store", id)
		}
		trace.add("✓ %s", title(s))
		return idMarkdown(s), trace.lines, true
	}

	if filepath.Ext(path) != ".go" {
		return trace.fail("no spell id under the cursor")
	}
	folder := filepath.Dir(path)
	pkg := filepath.Base(folder)

	if match := matchCovering([]*regexp.Regexp{spellConfigPattern}, lineText, column); match != nil {
		trace.add("SpellConfig")
		start := match[0]
		for _, previous := range lines[:line] {
			start += len(previous) + 1
		}
		call, err := callText(text, start)
		if err != nil {
			return trace.fail("%v", err)
		}
		result, err := evalSpellConfig(call, w.declarations(folder, path, text, trace), pkg, trace)
		if err != nil {
			return trace.fail("%v", err)
		}
		trace.add("✓ %s → %d fields", title(result.spell), len(result.rows))
		return configMarkdown(result), trace.lines, true
	}

	if match := matchCovering([]*regexp.Regexp{familyPattern}, lineText, column); match != nil {
		field := lineText[match[2]:match[3]]
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

	declarations := w.declarations(folder, path, text, trace)
	hover, ok := chainHoverAt(lineText, column, declarations, trace.rel(folder), trace)
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

func resultSummary(result *exprResult) string {
	parts := []string{result.title()}
	if result.readEffect > 0 {
		parts = append(parts, fmt.Sprintf("effect %d", result.readEffect))
	}
	if result.kind == kindValue {
		parts = append(parts, result.value)
	}
	return strings.Join(parts, " → ")
}

func chainHoverAt(lineText string, column int, declarations map[string]declaration, folder string, trace *tracer) (chainHover, bool) {
	declared, declErr := declarationOnLine(lineText)
	if declErr == nil && declared != nil {
		if hover, found, ok := hoverInDeclaration(declared, column, declarations, trace); found {
			return hover, ok
		}
	}

	name := identifierAt(lineText, column)
	if name == "" {
		trace.add("✗ no spell id, family or name under the cursor")
		return chainHover{}, false
	}
	trace.add("ident %s", name)
	if _, bound := declarations[name]; bound {
		return resolved(name, &chain{head: name}, false, declarations, trace)
	}
	if declErr != nil && name == declared.name {
		trace.add("✗ %s is bound to no chain the evaluator reads: %v", name, declErr)
	} else {
		trace.add("✗ no ladder-shaped declaration of %s in %s", name, folder)
	}
	return chainHover{}, false
}

func hoverInDeclaration(declared *declaration, column int, declarations map[string]declaration, trace *tracer) (chainHover, bool, bool) {
	c := declared.chain
	if column >= declared.nameStart && column <= declared.nameEnd {
		trace.add("ident %s", declared.name)
		trace.add("  %s = %s  (%s)", declared.name, c.text(false), trace.here)
		hover, ok := resolved(declared.name, c, false, declarations, trace)
		return hover, true, ok
	}
	for i, seg := range c.segments {
		if column >= seg.start && column <= seg.end {
			trace.add("segment %s", seg.text(false))
			prefix := &chain{head: c.head, segments: c.segments[:i+1]}
			hover, ok := resolved(seg.text(false), prefix, true, declarations, trace)
			return hover, true, ok
		}
	}
	if column >= c.start && column <= c.end {
		trace.add("ident %s", c.head)
		hover, ok := resolved(c.head, &chain{head: c.head}, true, declarations, trace)
		return hover, true, ok
	}
	return chainHover{}, false, false
}

func resolved(label string, c *chain, segment bool, declarations map[string]declaration, trace *tracer) (chainHover, bool) {
	resolved, err := resolveChain(c, declarations, trace)
	if err != nil {
		trace.add("✗ %v", err)
		return chainHover{}, false
	}
	return chainHover{label: label, chain: resolved, segment: segment}, true
}

// The chain with every name the package binds substituted by the chain it stands for, down to a
// ladder: at most maxSubstitutions names deep, and never through a name twice.
func resolveChain(c *chain, declarations map[string]declaration, trace *tracer) (*chain, error) {
	return resolveFrom(c, declarations, maxSubstitutions, map[string]bool{}, trace)
}

func resolveFrom(c *chain, declarations map[string]declaration, depth int, seen map[string]bool, trace *tracer) (*chain, error) {
	if strings.HasPrefix(c.head, "spellData.") {
		if len(c.segments) == 0 {
			return nil, fmt.Errorf("%s names a family, not a rank: follow it with Highest(), Rank(n) or ByID(id)", c.text(false))
		}
		return c, nil
	}

	bound, ok := declarations[c.head]
	if !ok {
		if len(c.segments) > 0 && isFamilyName(c.head) {
			return &chain{head: "spellData." + c.head, segments: c.segments}, nil
		}
		return nil, fmt.Errorf("no ladder-shaped declaration of %s in the package", c.head)
	}
	if seen[c.head] {
		return nil, fmt.Errorf("%s stands on itself", c.head)
	}
	if depth <= 0 {
		return nil, fmt.Errorf("%s stands on more than %d names", c.text(false), maxSubstitutions)
	}

	seen[c.head] = true
	trace.add("  %s = %s  (%s:%d)", c.head, bound.chain.text(false), trace.rel(bound.file), bound.line)
	head, err := resolveFrom(bound.chain, declarations, depth-1, seen, trace)
	if err != nil {
		return nil, err
	}
	return &chain{head: head.head, segments: append(slices.Clip(head.segments), c.segments...)}, nil
}

func identifierAt(lineText string, column int) string {
	match := matchCovering([]*regexp.Regexp{identPattern}, lineText, column)
	if match == nil {
		return ""
	}
	name := lineText[match[0]:match[1]]
	if goKeywords[name] {
		return ""
	}
	return name
}

func matchCovering(patterns []*regexp.Regexp, lineText string, column int) []int {
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
		units -= len(utf16.Encode([]rune{r}))
		offset += size
	}
	return offset
}

func uriPath(uri string) string {
	if u, err := url.Parse(uri); err == nil && u.Scheme == "file" {
		return filepath.Clean(u.Path)
	}
	return filepath.Clean(uri)
}

func pathURI(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return (&url.URL{Scheme: "file", Path: abs}).String()
}
