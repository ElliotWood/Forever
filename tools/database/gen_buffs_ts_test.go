package database

// Re-renders the settings-UI buff inputs from the client database and asserts the
// committed file agrees, plus a database-free check of the shapes the renderer
// chooses per proto type and scope.

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/tools/database/buffmanifest"
)

func TestGeneratedBuffsDebuffsTS(t *testing.T) {
	helper := openBuffTestDB(t)

	rows, err := ResolveBuffManifest(helper)
	if err != nil {
		t.Fatalf("resolving the manifest: %v", err)
	}
	rendered, err := RenderBuffsDebuffsTS(rows)
	if err != nil {
		t.Fatalf("rendering %s: %v", buffsDebuffsTSFile, err)
	}

	root, err := repoRoot()
	if err != nil {
		t.Fatalf("finding the repository root: %v", err)
	}
	committed, err := os.ReadFile(filepath.Join(root, buffsDebuffsTSFile))
	if err != nil {
		t.Fatalf("reading %s: %v", buffsDebuffsTSFile, err)
	}
	if !bytes.Equal(committed, rendered) {
		t.Errorf("%s does not match the client database, regenerate with `make spelldata`:\n%s",
			buffsDebuffsTSFile, unifiedBuffDiff(string(committed), string(rendered)))
	}
}

func TestRenderBuffsDebuffsTSShapes(t *testing.T) {
	rows := []ResolvedBuff{
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "shadow_protection", Scope: buffmanifest.ScopeRaid, Proto: buffmanifest.ProtoBool,
				Kind: buffmanifest.KindResistance, Go: "ShadowProtection", Owner: proto.Class_ClassPriest,
				Stats: []proto.Stat{proto.Stat_StatShadowResistance, proto.Stat_StatStamina},
			},
			SpellID: 10958, DBName: "Shadow Protection",
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "hunters_mark", Scope: buffmanifest.ScopeDebuff, Proto: buffmanifest.ProtoTristate,
				Kind: buffmanifest.KindDebuffStat, Go: "HuntersMark", Owner: proto.Class_ClassHunter,
				Stats: []proto.Stat{proto.Stat_StatRangedAttackPower},
			},
			SpellID: 14325, DBName: "Hunter's Mark",
			TalentCurve: []float64{71, 110}, TalentSpellID: 19425,
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "mana_tide_totems", Scope: buffmanifest.ScopeParty, Proto: buffmanifest.ProtoInt32,
				Kind: buffmanifest.KindExternalCD, Go: "ManaTideTotems", Owner: proto.Class_ClassShaman,
				Stats: []proto.Stat{proto.Stat_StatMP5}, Label: "Mana Tide Totem",
			},
			SpellID: 17359, DBName: "Mana Tide Totem",
		},
	}

	rendered, err := RenderBuffsDebuffsTS(rows)
	if err != nil {
		t.Fatalf("rendering three rows: %v", err)
	}
	out := string(rendered)

	want := []string{
		`export const ShadowProtection = makeBooleanRaidBuffInput({
	actionId: ActionId.fromSpellId(10958),
	fieldName: 'shadowProtection',
	label: 'Shadow Protection',
});`,
		`export const HuntersMark = makeTristateDebuffInput({
	actionId: ActionId.fromSpellId(14325),
	impId: ActionId.fromSpellId(19425),
	fieldName: 'huntersMark',
	label: "Hunter's Mark",
});`,
		`export const ManaTideTotems = makeMultistatePartyBuffInput({
	actionId: ActionId.fromSpellId(17359),
	numStates: 5,
	fieldName: 'manaTideTotems',
	label: 'Mana Tide Totem',
});`,
		`export const GENERATED_RAID_BUFFS_CONFIG: GeneratedStatOption[] = [
	{
		config: ShadowProtection,
		stats: [Stat.StatShadowResistance, Stat.StatStamina],
		ownerClass: Class.ClassPriest,
	},
];`,
		`export const GENERATED_PARTY_BUFFS_CONFIG: GeneratedStatOption[] = [
	{
		config: ManaTideTotems,
		stats: [Stat.StatMP5],
		ownerClass: Class.ClassShaman,
	},
];`,
		`export const GENERATED_INDIVIDUAL_BUFFS_CONFIG: GeneratedStatOption[] = [
];`,
		`export const GENERATED_DEBUFFS_CONFIG: GeneratedStatOption[] = [
	{
		config: HuntersMark,
		stats: [Stat.StatRangedAttackPower],
		ownerClass: Class.ClassHunter,
	},
];`,
	}
	for _, block := range want {
		if !strings.Contains(out, block) {
			t.Errorf("the rendered file is missing:\n%s\ngot:\n%s", block, out)
		}
	}
}

func TestRenderBuffsDebuffsTSSkips(t *testing.T) {
	rows := []ResolvedBuff{
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "drums", Scope: buffmanifest.ScopeParty, Proto: buffmanifest.ProtoEnumDrums,
				Kind: buffmanifest.KindEnum, Go: "Drums",
				Notes: "the drum items are not in the client.",
			},
			Reason: "the drum items are not in the client.",
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "blessing_of_salvation", Scope: buffmanifest.ScopeIndividual, Proto: buffmanifest.ProtoBool,
				Kind: buffmanifest.KindPseudoMult, Go: "BlessingOfSalvation", Owner: proto.Class_ClassPaladin,
			},
			SpellID: 1038, DBName: "Blessing of Salvation",
		},
	}

	rendered, err := RenderBuffsDebuffsTS(rows)
	if err != nil {
		t.Fatalf("rendering two skipped rows: %v", err)
	}
	out := string(rendered)

	if strings.Contains(out, "export const Drums") || strings.Contains(out, "export const BlessingOfSalvation") {
		t.Errorf("a row with no input rendered one:\n%s", out)
	}
	for _, comment := range []string{
		"// drums: the drum items are not in the client.",
		"// blessing_of_salvation: " + manualBuffInputs["blessing_of_salvation"],
	} {
		if !strings.Contains(out, comment) {
			t.Errorf("the rendered file is missing %q:\n%s", comment, out)
		}
	}
}

func TestRenderBuffsDebuffsTSRejectsAnUncountedInt32(t *testing.T) {
	rows := []ResolvedBuff{{
		BuffSpec: buffmanifest.BuffSpec{
			Field: "totem_of_wrath", Scope: buffmanifest.ScopeParty, Proto: buffmanifest.ProtoInt32,
			Kind: buffmanifest.KindItemCount, Go: "TotemOfWrath", Owner: proto.Class_ClassShaman,
		},
		SpellID: 30706, DBName: "Totem of Wrath",
	}}

	if _, err := RenderBuffsDebuffsTS(rows); err == nil {
		t.Error("an int32 row with no numStates rendered without an error")
	}
}
