package buffmanifest

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestUniqueScopeField(t *testing.T) {
	seen := map[string]bool{}
	for _, spec := range Manifest {
		key := fmt.Sprintf("%s/%s", spec.Scope, spec.Field)
		if seen[key] {
			t.Errorf("duplicate (scope, field): %s", key)
		}
		seen[key] = true
	}
}

func TestUniqueScopeNumber(t *testing.T) {
	seen := map[string]string{}
	for _, spec := range Manifest {
		key := fmt.Sprintf("%s/%d", spec.Scope, spec.Number)
		if other, ok := seen[key]; ok {
			t.Errorf("duplicate (scope, number) %s: %s and %s", key, other, spec.Field)
		}
		seen[key] = spec.Field
	}
}

func TestUniqueGoStem(t *testing.T) {
	seen := map[string]bool{}
	for _, spec := range Manifest {
		if spec.Go == "" {
			t.Errorf("%s has no Go stem", spec.Field)
			continue
		}
		if seen[spec.Go] {
			t.Errorf("duplicate Go stem: %s", spec.Go)
		}
		seen[spec.Go] = true
	}
}

func TestShellRowsHaveNotes(t *testing.T) {
	for _, spec := range Manifest {
		switch spec.Kind {
		case KindManual, KindAbsent, KindFlag:
			if spec.Notes == "" {
				t.Errorf("%s is %s and needs Notes", spec.Field, spec.Kind)
			}
		}
	}
}

func TestResolvableRowsHaveAnchor(t *testing.T) {
	for _, spec := range Manifest {
		switch spec.Kind {
		case KindAbsent, KindFlag, KindEnum:
			continue
		}
		if spec.Name == "" && spec.Anchor == 0 {
			t.Errorf("%s (%s) has neither Name nor Anchor", spec.Field, spec.Kind)
		}
	}
}

func TestProtoTypeMatchesKind(t *testing.T) {
	for _, spec := range Manifest {
		switch spec.Kind {
		case KindFlag:
			if spec.Proto != ProtoBool {
				t.Errorf("%s is KindFlag and must be %s, got %s", spec.Field, ProtoBool, spec.Proto)
			}
		case KindEnum:
			if spec.Proto != ProtoEnumDrums {
				t.Errorf("%s is KindEnum and must be %s, got %s", spec.Field, ProtoEnumDrums, spec.Proto)
			}
		case KindDebuffUptime:
			if spec.Proto != ProtoDouble {
				t.Errorf("%s is KindDebuffUptime and must be %s, got %s", spec.Field, ProtoDouble, spec.Proto)
			}
		case KindItemCount, KindExternalCD:
			if spec.Proto != ProtoInt32 {
				t.Errorf("%s is %s and must be %s, got %s", spec.Field, spec.Kind, ProtoInt32, spec.Proto)
			}
		}
		if spec.Proto == ProtoEnumDrums && spec.Kind != KindEnum {
			t.Errorf("%s is %s but only KindEnum may be %s", spec.Field, spec.Kind, ProtoEnumDrums)
		}
	}
}

func TestTalentImpliesTristate(t *testing.T) {
	for _, spec := range Manifest {
		if spec.Talent != nil && spec.Proto != ProtoTristate {
			t.Errorf("%s carries talent %q but is %s", spec.Field, spec.Talent.Name, spec.Proto)
		}
	}
}

func TestFieldNaming(t *testing.T) {
	cases := []struct {
		field, goName, tsName string
	}{
		{"battle_shout", "BattleShout", "battleShout"},
		{"soe_enhancement_2pt4", "SoeEnhancement_2Pt4", "soeEnhancement2Pt4"},
		{"joc_retribution_2pt4", "JocRetribution_2Pt4", "jocRetribution2Pt4"},
		{"snapshot_bs_t2", "SnapshotBsT2", "snapshotBsT2"},
		{"snapshot_bs_booming_voice_rank", "SnapshotBsBoomingVoiceRank", "snapshotBsBoomingVoiceRank"},
		{"isb_uptime", "IsbUptime", "isbUptime"},
		{"expose_weakness_hunter_agility", "ExposeWeaknessHunterAgility", "exposeWeaknessHunterAgility"},
	}
	for _, c := range cases {
		spec := BuffSpec{Field: c.field}
		if got := spec.GoField(); got != c.goName {
			t.Errorf("GoField(%q) = %q, want %q", c.field, got, c.goName)
		}
		if got := spec.TSField(); got != c.tsName {
			t.Errorf("TSField(%q) = %q, want %q", c.field, got, c.tsName)
		}
	}
}

func TestFieldNamesRoundTrip(t *testing.T) {
	for _, spec := range Manifest {
		goName := spec.GoField()
		if goName == "" {
			t.Errorf("%s has an empty GoField", spec.Field)
			continue
		}
		want := strings.ReplaceAll(goName, "_", "")
		want = strings.ToLower(want[:1]) + want[1:]
		if got := spec.TSField(); got != want {
			t.Errorf("%s: TSField %q does not match GoField %q", spec.Field, got, goName)
		}
	}
}

func TestCensusMatchesProto(t *testing.T) {
	messages := map[string]BuffScope{
		"RaidBuffs":       ScopeRaid,
		"PartyBuffs":      ScopeParty,
		"IndividualBuffs": ScopeIndividual,
		"Debuffs":         ScopeDebuff,
	}

	want := map[string]int32{}
	for message, scope := range messages {
		for field, number := range parseProtoMessage(t, message) {
			want[fmt.Sprintf("%s/%s", scope, field)] = number
		}
	}

	got := map[string]int32{}
	for _, spec := range Manifest {
		got[fmt.Sprintf("%s/%s", spec.Scope, spec.Field)] = spec.Number
	}

	for key, number := range want {
		if other, ok := got[key]; !ok {
			t.Errorf("%s is in common.proto but not in the manifest", key)
		} else if other != number {
			t.Errorf("%s has number %d in the manifest and %d in common.proto", key, other, number)
		}
	}
	for key := range got {
		if _, ok := want[key]; !ok {
			t.Errorf("%s is in the manifest but not in common.proto", key)
		}
	}
}

var protoFieldRE = regexp.MustCompile(`^\s*(?:[A-Za-z][\w.]*)\s+([a-z][a-z0-9_]*)\s*=\s*(\d+)\s*;`)

func parseProtoMessage(t *testing.T, message string) map[string]int32 {
	t.Helper()

	raw, err := os.ReadFile("../../../proto/common.proto")
	if err != nil {
		t.Fatalf("read common.proto: %v", err)
	}

	fields := map[string]int32{}
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
		number, err := strconv.Atoi(match[2])
		if err != nil {
			t.Fatalf("field %s of %s: %v", match[1], message, err)
		}
		fields[match[1]] = int32(number)
	}

	if len(fields) == 0 {
		t.Fatalf("no fields parsed for message %s", message)
	}
	return fields
}
