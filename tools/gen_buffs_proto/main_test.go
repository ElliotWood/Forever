package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/wowsims/forever/tools/database/buffmanifest"
)

type protoField struct {
	Type   string
	Number int32
}

var messageNames = []string{"RaidBuffs", "PartyBuffs", "IndividualBuffs", "Debuffs"}

var scopeOfMessage = map[string]buffmanifest.BuffScope{
	"RaidBuffs":       buffmanifest.ScopeRaid,
	"PartyBuffs":      buffmanifest.ScopeParty,
	"IndividualBuffs": buffmanifest.ScopeIndividual,
	"Debuffs":         buffmanifest.ScopeDebuff,
}

func TestRenderMatchesCommittedFile(t *testing.T) {
	committed, err := os.ReadFile("../../proto/buffs.proto")
	if err != nil {
		t.Fatalf("read buffs.proto: %v", err)
	}

	rendered := Render(buffmanifest.Manifest)
	if bytes.Equal(rendered, committed) {
		return
	}

	t.Errorf("proto/buffs.proto does not match the manifest; run `go run ./tools/gen_buffs_proto`")
	committedLines := strings.Split(string(committed), "\n")
	for i, line := range strings.Split(string(rendered), "\n") {
		if i >= len(committedLines) {
			t.Errorf("rendered line %d %q is past the end of the committed file", i+1, line)
			return
		}
		if line != committedLines[i] {
			t.Errorf("line %d: rendered %q, committed %q", i+1, line, committedLines[i])
			return
		}
	}
}

// The next free number is past every live field and every reserved one, so that a
// new row cannot take a number a pre-17 payload still carries.
func TestRenderNextIndex(t *testing.T) {
	rendered := string(Render(buffmanifest.Manifest))

	for _, message := range messageNames {
		var max int32
		for _, field := range parseProtoMessage(t, []byte(rendered), message) {
			if field.Number > max {
				max = field.Number
			}
		}
		for _, number := range parseReserved(t, rendered, message) {
			if number > max {
				max = number
			}
		}
		want := fmt.Sprintf("// Next index: %d\nmessage %s {\n", max+1, message)
		if !strings.Contains(rendered, want) {
			t.Errorf("%s is not preceded by %q", message, want)
		}
	}
}

func TestRenderReservesEveryRetiredNumber(t *testing.T) {
	rendered := string(Render(buffmanifest.Manifest))

	for message, scope := range scopeOfMessage {
		numbers := parseReserved(t, rendered, message)
		wantNumbers := buffmanifest.RetiredNumbers(scope)
		if fmt.Sprint(numbers) != fmt.Sprint(wantNumbers) {
			t.Errorf("%s reserves the numbers %v, want %v", message, numbers, wantNumbers)
		}

		names := parseReservedNames(t, rendered, message)
		wantNames := buffmanifest.RetiredFields(scope)
		if fmt.Sprint(names) != fmt.Sprint(wantNames) {
			t.Errorf("%s reserves the names %v, want %v", message, names, wantNames)
		}

		for name, field := range parseProtoMessage(t, []byte(rendered), message) {
			if slices.Contains(wantNumbers, field.Number) {
				t.Errorf("%s declares a field with the reserved number %d", message, field.Number)
			}
			if slices.Contains(wantNames, name) {
				t.Errorf("%s declares a field with the reserved name %q", message, name)
			}
		}
	}
}

// The 25 fields api version 17 typed bool where they were a TristateEffect. The
// list is the one ui/sim/proto/buff_field_migration.ts rewrites before parsing, so
// the two are checked against each other: a name on one side only would leave a
// saved payload either unconverted or converted into the wrong type.
var retypedFields = map[string][]string{
	"RaidBuffs":       {"power_word_fortitude", "divine_spirit", "gift_of_the_wild", "thorns"},
	"PartyBuffs":      {"blood_pact", "moonkin_aura", "leader_of_the_pack", "devotion_aura", "retribution_aura", "concentration_aura", "grace_of_air_totem", "strength_of_earth_totem", "windfury_totem", "battle_shout", "commanding_shout"},
	"IndividualBuffs": {"blessing_of_wisdom", "blessing_of_might"},
	"Debuffs":         {"improved_seal_of_the_crusader", "curse_of_elements", "expose_armor", "faerie_fire", "hunters_mark", "demoralizing_roar", "demoralizing_shout", "thunder_clap"},
}

func TestRetypedFieldsAreBool(t *testing.T) {
	rendered := Render(buffmanifest.Manifest)

	count := 0
	for message, names := range retypedFields {
		fields := parseProtoMessage(t, rendered, message)
		for _, name := range names {
			count++
			field, ok := fields[name]
			if !ok {
				t.Errorf("%s.%s is retyped but not declared", message, name)
				continue
			}
			if field.Type != "bool" {
				t.Errorf("%s.%s is %s, want bool", message, name, field.Type)
			}
		}
	}
	if count != 25 {
		t.Errorf("the retyped list names %d fields, want the 25 api version 17 declares", count)
	}
}

func TestRetypedFieldsMatchTheMigration(t *testing.T) {
	raw := readMigration(t)

	for message, names := range retypedFields {
		for _, name := range names {
			want := "'" + tsFieldName(name) + "'"
			if !strings.Contains(raw, want) {
				t.Errorf("%s.%s is retyped but the migration does not name %s", message, name, want)
			}
		}
	}
}

// The migration deletes the retired keys before `fromJson` sees them, so it has to
// name exactly the fields the manifest retired: one it misses throws on load, and
// one it names that is live would silently drop a setting.
func TestRetiredFieldsMatchTheMigration(t *testing.T) {
	raw := readMigration(t)

	groups := map[buffmanifest.BuffScope]string{
		buffmanifest.ScopeRaid:       "raidBuffs",
		buffmanifest.ScopeParty:      "partyBuffs",
		buffmanifest.ScopeIndividual: "individualBuffs",
		buffmanifest.ScopeDebuff:     "debuffs",
	}

	declared := migrationRetiredFields(t, raw)
	total := 0
	for scope, group := range groups {
		want := buffmanifest.RetiredFields(scope)
		total += len(want)
		got := declared[group]
		slices.Sort(got)
		sorted := slices.Clone(want)
		slices.Sort(sorted)
		if fmt.Sprint(got) != fmt.Sprint(sorted) {
			t.Errorf("the migration retires %s %v, want %v", group, got, sorted)
		}
	}
	if total != 33 {
		t.Errorf("the manifest retires %d fields, want the 33 api version 17 declares", total)
	}
}

func readMigration(t *testing.T) string {
	t.Helper()

	raw, err := os.ReadFile("../../ui/sim/proto/buff_field_migration.ts")
	if err != nil {
		t.Fatalf("read buff_field_migration.ts: %v", err)
	}
	return string(raw)
}

var migrationGroupRE = regexp.MustCompile(`(?s)export const retiredBuffFields = \{(.*?)\n\} as const;`)

// The quoted proto names under each group of `retiredBuffFields`.
func migrationRetiredFields(t *testing.T, raw string) map[string][]string {
	t.Helper()

	body := migrationGroupRE.FindStringSubmatch(raw)
	if body == nil {
		t.Fatal("buff_field_migration.ts declares no retiredBuffFields")
	}

	out := map[string][]string{}
	group := ""
	for _, line := range strings.Split(body[1], "\n") {
		trimmed := strings.TrimSpace(line)
		if name, rest, found := strings.Cut(trimmed, ":"); found && !strings.HasPrefix(trimmed, "'") {
			group = name
			if _, ok := out[group]; !ok {
				out[group] = []string{}
			}
			trimmed = rest
		}
		for _, quoted := range strings.Split(trimmed, ",") {
			quoted = strings.TrimSpace(strings.Trim(strings.TrimSpace(quoted), "[]"))
			if len(quoted) < 2 || !strings.HasPrefix(quoted, "'") {
				continue
			}
			out[group] = append(out[group], strings.Trim(quoted, "'"))
		}
	}
	return out
}

// The property protobuf-ts generates for a proto field name, which is how the
// migration's retyped list spells them. None of the 25 holds a digit, so the
// digit rule the retired list needs does not apply here.
func tsFieldName(protoName string) string {
	parts := strings.Split(protoName, "_")
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, "")
}

var reservedRE = regexp.MustCompile(`^\s*reserved\s+([\d,\s]+);`)

var reservedNameRE = regexp.MustCompile(`^\s*reserved\s+("[^;]+);`)

func parseReservedNames(t *testing.T, rendered string, message string) []string {
	t.Helper()

	var names []string
	for _, line := range messageLines(rendered, message) {
		match := reservedNameRE.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		for _, part := range strings.Split(match[1], ",") {
			name, err := strconv.Unquote(strings.TrimSpace(part))
			if err != nil {
				t.Fatalf("reserved name %q of %s: %v", part, message, err)
			}
			names = append(names, name)
		}
	}
	return names
}

// The lines between `message <name> {` and its closing brace.
func messageLines(rendered string, message string) []string {
	var lines []string
	inMessage := false
	for _, line := range strings.Split(rendered, "\n") {
		if !inMessage {
			inMessage = strings.HasPrefix(strings.TrimSpace(line), "message "+message+" {")
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "}") {
			break
		}
		lines = append(lines, line)
	}
	return lines
}

func parseReserved(t *testing.T, rendered string, message string) []int32 {
	t.Helper()

	var numbers []int32
	for _, line := range messageLines(rendered, message) {
		match := reservedRE.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		for _, part := range strings.Split(match[1], ",") {
			number, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				t.Fatalf("reserved number %q of %s: %v", part, message, err)
			}
			numbers = append(numbers, int32(number))
		}
	}
	return numbers
}

func TestRenderHeader(t *testing.T) {
	rendered := string(Render(buffmanifest.Manifest))
	want := `// Code generated by tools/gen_buffs_proto from tools/database/buffmanifest. DO NOT EDIT.
syntax = "proto3";
package proto;
option go_package = "./proto";

import "common.proto";
`
	if !strings.HasPrefix(rendered, want) {
		t.Errorf("rendered header is\n%s\nwant\n%s", rendered[:len(want)], want)
	}
	for _, line := range strings.Split(rendered, "\n") {
		if strings.HasPrefix(line, " ") {
			t.Errorf("line %q is space-indented; common.proto uses tabs", line)
		}
	}
}

var protoFieldRE = regexp.MustCompile(`^\s*([A-Za-z][\w.]*)\s+([a-z][a-z0-9_]*)\s*=\s*(\d+)\s*;`)

func parseProtoMessage(t *testing.T, raw []byte, message string) map[string]protoField {
	t.Helper()

	fields := map[string]protoField{}
	inMessage := false
	for _, line := range strings.Split(string(raw), "\n") {
		if !inMessage {
			inMessage = strings.HasPrefix(strings.TrimSpace(line), "message "+message+" {")
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "}") {
			break
		}
		match := protoFieldRE.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		number, err := strconv.Atoi(match[3])
		if err != nil {
			t.Fatalf("field %s of %s: %v", match[2], message, err)
		}
		fields[match[2]] = protoField{Type: match[1], Number: int32(number)}
	}

	if len(fields) == 0 {
		t.Fatalf("no fields parsed for message %s", message)
	}
	return fields
}
