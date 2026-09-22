package database

// Pins what the generator emits for a row it can express. Only battle_shout of
// the real manifest renders code so far - every other row is still a shell - so
// without these synthetic rows most of the supported branch of the templates
// would be untested, and nothing would notice a generated constructor that no
// longer compiles against sim/core/buffs_gen_support.go.
//
// Needs no client database. Set UPDATE_BUFF_FIXTURES=1 to rewrite the fixtures
// after a deliberate change.

import (
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/wowsims/forever/sim/core/stats"
	"github.com/wowsims/forever/tools/database/buffmanifest"
	"github.com/wowsims/forever/tools/database/dbc"
)

// One row per shape the templates have a branch for. The proto field of each row
// is a real one whose compiled type matches the row's declared type, so the
// rendered apply blocks type-check against the sim.
func syntheticBuffRows() []ResolvedBuff {
	return []ResolvedBuff{
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "mana_spring_totem", Scope: buffmanifest.ScopeParty,
				Proto: buffmanifest.ProtoTristate, Kind: buffmanifest.KindStatFlat,
				Go: "SynthManaSpring", Name: "Mana Spring Totem", Category: "ManaSpringTotem",
			},
			SpellID: 10494, Supported: true,
			Stats:         []StatAmount{{Stat: stats.MP5, Amount: 25}},
			TalentCurve:   []float64{25, 26, 27, 28, 30, 31},
			TalentApplies: buffmanifest.TalentScalesValue,
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "greater_blessing_of_kings", Scope: buffmanifest.ScopeIndividual,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindStatPct,
				Go: "SynthBlessingOfKings", Name: "Blessing of Kings",
				Pet: buffmanifest.PetStripWhenSummonedLate,
			},
			SpellID: 20217, DurationMs: 3600000, Supported: true,
			Stats: []StatAmount{
				{Stat: stats.Strength, Amount: 1.1, Multiplicative: true},
				{Stat: stats.Agility, Amount: 1.1, Multiplicative: true},
			},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "battle_shout", Scope: buffmanifest.ScopeParty,
				Proto: buffmanifest.ProtoTristate, Kind: buffmanifest.KindStatFlat,
				Go: "SynthBattleShout", Name: "Battle Shout", Category: "SynthBattleShout",
				SingleAura: true, Driver: true,
			},
			SpellID: 25289, DurationMs: 180000, Supported: true,
			Stats: []StatAmount{{Stat: stats.AttackPower, Amount: 139}},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "devotion_aura", Scope: buffmanifest.ScopeParty,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindResistance,
				Go: "SynthDevotionAura", Name: "Devotion Aura", Category: "DevotionAura",
				SharedCategory: "SynthPaladinAura", SingleAura: true,
			},
			SpellID: 10293, DurationMs: 600000, Supported: true,
			Stats:         []StatAmount{{Stat: stats.Armor, Amount: 735}},
			TalentCurve:   []float64{600000, 900000, 1200000},
			TalentApplies: buffmanifest.TalentScalesDuration,
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "frost_resistance_aura", Scope: buffmanifest.ScopeRaid,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindResistance,
				Go: "SynthFrostResistanceAura", Name: "Frost Resistance Aura",
				Category: "FrostResistanceAura", SharedCategory: "SynthPaladinAura", SingleAura: true,
			},
			SpellID: 19898, Supported: true,
			Stats: []StatAmount{{Stat: stats.FrostResistance, Amount: 60}},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "thunder_clap", Scope: buffmanifest.ScopeDebuff,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindDebuffAtkSpeed,
				Go: "SynthThunderClap", Name: "Thunder Clap", Category: "AtkSpdReduction",
			},
			SpellID: 11581, DurationMs: 30000, Supported: true,
			Pseudo: []PseudoMod{
				{Kind: "MeleeSpeedMultiplier", Amount: 0.8, Multiplicative: true},
			},
			TalentCurve:    []float64{0.8, 0.78, 0.76},
			TalentApplies:  buffmanifest.TalentScalesValue,
			TalentOnPseudo: true,
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "innervates", Scope: buffmanifest.ScopeIndividual,
				Proto: buffmanifest.ProtoInt32, Kind: buffmanifest.KindExternalCD,
				Go: "SynthInnervates", Name: "Innervate", Label: "Innervates",
				Category: "Innervate",
			},
			SpellID: 29166, DurationMs: 20000, CooldownMs: 360000, Supported: true,
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "power_infusions", Scope: buffmanifest.ScopeIndividual,
				Proto: buffmanifest.ProtoInt32, Kind: buffmanifest.KindExternalCD,
				Go: "SynthPowerInfusions", Name: "Power Infusion", Label: "Power Infusions",
				Category: "PowerInfusion",
			},
			SpellID: 10060, DurationMs: 15000, CooldownMs: 180000, Supported: true,
			Pseudo: []PseudoMod{
				{Kind: "SchoolDamageDealtMultiplier", Amount: 1.2, Multiplicative: true, SchoolMask: 126},
				{Kind: "HealingDealtMultiplier", Amount: 1.2, Multiplicative: true},
			},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "totem_twisting", Scope: buffmanifest.ScopeParty,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindAbsent,
				Go: "SynthInheritedNeck", Category: "Inherited Neck",
				Pet: buffmanifest.PetInheritOwnerAura,
			},
			Reason: "the neck has no Item row",
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "atiesh_mage", Scope: buffmanifest.ScopeParty,
				Proto: buffmanifest.ProtoInt32, Kind: buffmanifest.KindItemCount,
				Go: "SynthAtieshMage", Anchor: 28142, Label: "Atiesh - Mage",
			},
			SpellID: 28142, Supported: true,
			Stats: []StatAmount{{Stat: stats.SpellCritPercent, Amount: 2}},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "thorns", Scope: buffmanifest.ScopeRaid,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindDamageShield,
				Go: "SynthThorns", Name: "Thorns", Category: "Thorns",
				Pet: buffmanifest.PetStrip,
			},
			SpellID: 9910, DurationMs: 600000, SchoolMask: 8, Supported: true,
			Effects:       []ResolvedEffect{{Index: 0, Effect: 6, Aura: 15, Value: 22}},
			TalentCurve:   []float64{22, 27, 33},
			TalentApplies: buffmanifest.TalentScalesValue,
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "sunder_armor", Scope: buffmanifest.ScopeDebuff,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindDebuffStacking,
				Go: "SynthSunderArmor", Name: "Sunder Armor", Category: "MajorArmorReduction",
				SingleAura: true, Driver: true,
			},
			SpellID: 11597, DurationMs: 30000, MaxStacks: 5, Supported: true,
			Stats: []StatAmount{{Stat: stats.Armor, Amount: -450}},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "expose_armor", Scope: buffmanifest.ScopeDebuff,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindDebuffStat,
				Go: "SynthExposeArmor", Name: "Expose Armor", Category: "MajorArmorReduction",
				SingleAura: true,
			},
			SpellID: 11198, DurationMs: 30000, Supported: true,
			Note:  "Effect 0 is worth -450.0 per combo point; this is the 5-point finisher.",
			Stats: []StatAmount{{Stat: stats.Armor, Amount: -2250}},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "curse_of_elements", Scope: buffmanifest.ScopeDebuff,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindDebuffDamageTaken,
				Go: "SynthCurseOfElements", Name: "Curse of the Elements", Category: "CurseOfElements",
				SingleAura: true,
			},
			SpellID: 1311680, DurationMs: 300000, Supported: true,
			Stats: []StatAmount{
				{Stat: stats.FireResistance, Amount: -75},
				{Stat: stats.NatureResistance, Amount: -75},
				{Stat: stats.FrostResistance, Amount: -75},
				{Stat: stats.ShadowResistance, Amount: -75},
				{Stat: stats.ArcaneResistance, Amount: -75},
			},
			Pseudo: []PseudoMod{
				{Kind: "SchoolDamageTakenMultiplier", Amount: 1.1, Multiplicative: true, SchoolMask: 126},
			},
		},
	}
}

// The curve has to scale the amount the client states and convert afterwards: a
// percentage aura reaches the sim as 1 + value/100, and spending talent points
// on that number instead of on the client's -20 would price the buff at 0.9 + n.
func TestTalentCurveScalesTheClientAmount(t *testing.T) {
	row := ResolvedBuff{Effects: []ResolvedEffect{
		{Effect: dbc.E_APPLY_AURA, Aura: dbc.A_MOD_MELEE_HASTE_3, Value: -20},
	}}

	target, onPseudo, ok := row.talentTarget()
	if !ok || !onPseudo {
		t.Fatalf("talentTarget() = %v, onPseudo %v, ok %v; want the pseudo-stat effect", target, onPseudo, ok)
	}

	if got := row.convertedAmount(target, onPseudo); math.Abs(got-0.8) > 1e-9 {
		t.Errorf("untalented amount = %v, want 0.8", got)
	}

	scaled := target
	scaled.Value = math.Trunc(applyTalentPoints(target.Value, 10, dbc.A_ADD_PCT_MODIFIER))
	if scaled.Value != -22 {
		t.Errorf("scaled client amount = %v, want -22", scaled.Value)
	}
	if got := row.convertedAmount(scaled, onPseudo); math.Abs(got-0.78) > 1e-9 {
		t.Errorf("talented amount = %v, want 0.78", got)
	}
}

// The sim holds flat damage taken in a physical field and a spell field, and
// the client states the aura with a school mask. A mask naming some spell
// schools and not others - Judgement of the Crusader is holy alone - fits
// neither, and the spell field would raise what every school does to the target.
func TestSchoolMaskedDamageTakenHasNoFieldToLandOn(t *testing.T) {
	physical := ResolvedEffect{Aura: dbc.A_MOD_DAMAGE_TAKEN, Misc: 1, Value: 8}
	if mods, ok := pseudoModsOf(physical); !ok || mods[0].Kind != "BonusPhysicalDamageTaken" || mods[0].Amount != 8 {
		t.Errorf("a physical mask maps to %v, ok %v; want 8 BonusPhysicalDamageTaken", mods, ok)
	}

	everySchool := ResolvedEffect{Aura: dbc.A_MOD_DAMAGE_TAKEN, Misc: 126, Value: 40}
	if mods, ok := pseudoModsOf(everySchool); !ok || mods[0].Kind != "BonusSpellDamageTaken" || mods[0].Amount != 40 {
		t.Errorf("a mask of every spell school maps to %v, ok %v; want 40 BonusSpellDamageTaken", mods, ok)
	}

	holy := ResolvedEffect{Aura: dbc.A_MOD_DAMAGE_TAKEN, Misc: 2, Value: 161}
	if mods, ok := pseudoModsOf(holy); ok {
		t.Errorf("a holy-only mask maps to %v, want the row to stay a shell", mods)
	}
}

var syntheticFixtures = map[string]string{
	buffsGenFile:   filepath.Join("testdata", "buffs_synthetic_auto_gen.go"),
	debuffsGenFile: filepath.Join("testdata", "debuffs_synthetic_auto_gen.go"),
}

func renderSyntheticBuffFiles(t *testing.T) map[string][]byte {
	t.Helper()

	files, err := renderBuffFiles(syntheticBuffRows())
	if err != nil {
		t.Fatalf("rendering the synthetic rows: %v", err)
	}
	return files
}

func TestRenderedBuffFilesMatchTheFixtures(t *testing.T) {
	files := renderSyntheticBuffFiles(t)

	for name, rendered := range files {
		fixture := syntheticFixtures[name]
		if os.Getenv("UPDATE_BUFF_FIXTURES") != "" {
			if err := os.WriteFile(fixture, rendered, 0644); err != nil {
				t.Fatalf("writing %s: %v", fixture, err)
			}
			continue
		}

		committed, err := os.ReadFile(fixture)
		if err != nil {
			t.Fatalf("reading %s: %v", fixture, err)
		}
		if string(committed) != string(rendered) {
			t.Errorf("%s no longer matches what the generator emits, "+
				"rewrite it with `UPDATE_BUFF_FIXTURES=1 go test ./tools/database/`:\n%s",
				fixture, unifiedBuffDiff(string(committed), string(rendered)))
		}
	}
}

// Builds sim/core with the rendered files overlaid onto it, which is the only
// check that a generated constructor still names an identifier the support API
// declares. The apply functions are renamed because the real generated files
// already declare them, and the drivers the rows call are stubbed here the way
// sim/core/buffs_manual.go would declare them.
func TestRenderedBuffFilesCompile(t *testing.T) {
	goTool := findGoTool(t)
	root, err := repoRoot()
	if err != nil {
		t.Fatalf("finding the repository root: %v", err)
	}

	dir := t.TempDir()
	overlay := map[string]string{}
	for name, rendered := range renderSyntheticBuffFiles(t) {
		body := strings.ReplaceAll(string(rendered), "func applyGenerated", "func synthApplyGenerated")
		path := filepath.Join(dir, filepath.Base(name))
		if err := os.WriteFile(path, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
		overlay[filepath.Join(root, "sim", "core", "zz_synthetic_"+filepath.Base(name))] = path
	}

	drivers := filepath.Join(dir, "drivers.go")
	if err := os.WriteFile(drivers, []byte("package core\n\n"+
		"import \"github.com/wowsims/forever/sim/core/proto\"\n\n"+
		"func driveSynthInnervates(char *Character, individual *proto.IndividualBuffs) {\n"+
		"\tnewGeneratedExternalCD(char, SynthInnervatesAura(&char.Unit, false, 0),"+
		" GeneratedExternalCD{NumSources: individual.Innervates,"+
		" Cooldown: SynthInnervatesCooldown(), Type: CooldownTypeMana})\n}\n\n"+
		"func driveSynthBattleShout(char *Character, _ *proto.PartyBuffs) {\n"+
		"\tApplyFixedShoutAura(char, SynthBattleShoutAura(&char.Unit, false, 0),"+
		" SynthBattleShoutCategory)\n}\n\n"+
		"func driveSynthSunderArmor(target *Unit, _ *proto.Debuffs, _ *proto.Raid) {\n"+
		"\tMakePermanent(SynthSunderArmorAura(target, false, 0))\n}\n\n"+
		"func driveSynthAtieshMage(char *Character, party *proto.PartyBuffs) {\n"+
		"\tMakePermanent(SynthAtieshMageAura(&char.Unit, false, 0,"+
		" float64(party.AtieshMage)))\n}\n\n"+
		"func driveSynthPowerInfusions(char *Character, individual *proto.IndividualBuffs) {\n"+
		"\tnewGeneratedExternalCD(char, SynthPowerInfusionsAura(&char.Unit, false, 0),"+
		" GeneratedExternalCD{NumSources: individual.PowerInfusions,"+
		" Cooldown: SynthPowerInfusionsCooldown(), Type: CooldownTypeDPS})\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	overlay[filepath.Join(root, "sim", "core", "zz_synthetic_drivers.go")] = drivers

	overlayPath := filepath.Join(dir, "overlay.json")
	if err := os.WriteFile(overlayPath, []byte(overlayJSON(overlay)), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(goTool, "build", "-overlay", overlayPath, "./sim/core/")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("the generated constructors do not compile against sim/core:\n%s", out)
	}
}

func overlayJSON(replace map[string]string) string {
	var b strings.Builder
	b.WriteString(`{"Replace":{`)
	first := true
	for from, to := range replace {
		if !first {
			b.WriteString(",")
		}
		first = false
		b.WriteString(`"` + filepath.ToSlash(from) + `":"` + filepath.ToSlash(to) + `"`)
	}
	b.WriteString("}}")
	return b.String()
}

func findGoTool(t *testing.T) string {
	t.Helper()

	if path := filepath.Join(runtime.GOROOT(), "bin", "go"); runtime.GOROOT() != "" {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	path, err := exec.LookPath("go")
	if err != nil {
		t.Skip("no go tool on PATH, so the generated files cannot be type-checked")
	}
	return path
}
