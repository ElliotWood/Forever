package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wowsims/forever/sim/core/spelldata"
)

func cell(text string) string {
	return strings.ReplaceAll(text, "|", `\|`)
}

func codeCell(text string) string {
	return "`" + cell(text) + "`"
}

func idMarkdown(s *spelldata.Spell) string {
	var md strings.Builder
	spellCard(&md, s, "**"+title(s)+"**", 0)
	return md.String()
}

func spellCard(md *strings.Builder, s *spelldata.Spell, heading string, read int) {
	fmt.Fprintf(md, "%s\n\n", heading)
	for _, ref := range ladderRefs(s.ID) {
		fmt.Fprintf(md, "`%s`  \n", ref)
	}

	rows := headerFields(s)
	if proc := procSummary(s); proc != "" {
		rows = append(rows, headerField{"proc", proc})
	}
	if refs := refList(s); len(refs) > 0 {
		rows = append(rows, headerField{"refs", strings.Join(refs, ", ")})
	}
	if len(rows) > 0 {
		md.WriteString("\n| | |\n|---|---|\n")
		for _, row := range rows {
			fmt.Fprintf(md, "| %s | %s |\n", row.label, cell(row.value))
		}
	}

	effectTable(md, effectLines(s), read, 0)
	fmt.Fprintf(md, "\n[Wowhead](%s)\n", wowheadURL(s.ID))
}

func effectTable(md *strings.Builder, effects []Line, read, onlyEffect int) {
	if len(effects) == 0 {
		return
	}
	md.WriteString("\n| # | effect | client row |\n|---|---|---|\n")
	for i, effect := range effects {
		n := i + 1
		if onlyEffect > 0 && n != onlyEffect {
			continue
		}
		marker := fmt.Sprint(n)
		if n == read {
			marker += " ▶"
		}
		fmt.Fprintf(md, "| %s | %s | %s |\n", marker, cell(effect.Human), codeCell(effect.Literal))
	}
}

func familyMarkdown(f *ladderFamily) string {
	var md strings.Builder
	fmt.Fprintf(&md, "**%s**\n\n", f.key())
	rows := familyRows(f)
	if len(rows) > 0 && rows[0].Value != "" {
		md.WriteString("| id | name | rank | call | value |\n|---|---|---|---|---|\n")
		for _, row := range rows {
			fmt.Fprintf(&md, "| %d | %s | %s | %s | %s |\n", row.ID, cell(row.Name), cell(row.Rank), codeCell(row.Accessor), cell(row.Value))
		}
	} else {
		md.WriteString("| id | name | rank | call |\n|---|---|---|---|\n")
		for _, row := range rows {
			fmt.Fprintf(&md, "| %d | %s | %s | %s |\n", row.ID, cell(row.Name), cell(row.Rank), codeCell(row.Accessor))
		}
	}
	md.WriteString("\n")
	if highest, err := familyHighest(f); err == nil {
		spellCard(&md, highest.spell, "**"+highest.title()+"**", 0)
	}
	return md.String()
}

func exprMarkdown(result *exprResult, hover chainHover) string {
	var md strings.Builder
	s := result.spell
	called := lastCall(result.trail)

	switch result.kind {
	case kindSpell:
		spellCard(&md, s, fmt.Sprintf("`%s` = %s", hover.label, result.title()), 0)

	case kindEffect:
		fmt.Fprintf(&md, "`%s` = **%s** of %s\n\n", called, result.value, result.title())
		fmt.Fprintf(&md, "`%s`\n", result.trail)
		if result.readEffect > 0 {
			effectTable(&md, effectLines(s), result.readEffect, result.readEffect)
		}
		if len(result.accessors) > 0 {
			md.WriteString("\n")
		}
		for _, accessor := range result.accessors {
			fmt.Fprintf(&md, "`%s`  \n", accessor)
		}
		fmt.Fprintf(&md, "\n[Wowhead](%s)\n", wowheadURL(s.ID))

	default:
		label := hover.label
		if hover.segment {
			label = called
		}
		fmt.Fprintf(&md, "`%s` = **%s**\n\n", label, result.value)
		fmt.Fprintf(&md, "`%s`\n\n", result.trail)
		if hover.segment && result.doc != "" {
			fmt.Fprintf(&md, "%s\n\n", result.doc)
		}
		spellCard(&md, s, "**"+result.title()+"**", result.readEffect)
	}
	return md.String()
}

var lastCallPattern = regexp.MustCompile(`\.([A-Za-z_]\w*\([^()]*\))$`)

// `Average(60)` out of a trail ending `.EffectN(1).Average(60)`.
func lastCall(trail string) string {
	if match := lastCallPattern.FindStringSubmatch(trail); match != nil {
		return match[1]
	}
	return trail
}
