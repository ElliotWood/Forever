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
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/tools/database/buffmanifest"
	"github.com/wowsims/forever/tools/database/dbc"
)

// One row per shape the templates have a branch for. The proto field of each row
// is a real one whose compiled type matches the row's declared type, so the
// rendered apply blocks type-check against the sim. Each row states its effects
// and goes through the same mapping a resolved row does.
func syntheticBuffRows() []ResolvedBuff {
	aura := func(a dbc.EffectAuraType, misc int32, value float64) ResolvedEffect {
		return ResolvedEffect{Effect: dbcenums.E_APPLY_AURA, Aura: a, Misc: misc, Value: value}
	}
	rows := []ResolvedBuff{
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "mana_spring_totem", Scope: buffmanifest.ScopeParty,
				Proto: buffmanifest.ProtoTristate, Kind: buffmanifest.KindStatFlat,
				Go: "SynthManaSpring", Name: "Mana Spring Totem", Category: "ManaSpringTotem",
			},
			SpellID: 10494, CastSpellID: 10494, Supported: true,
			Effects:       []ResolvedEffect{{Effect: dbcenums.E_APPLY_AURA, Aura: dbcenums.A_PERIODIC_ENERGIZE, Value: 10, PeriodMs: 2000}},
			TalentRanks:   5,
			TalentApplies: buffmanifest.TalentScalesValue,
			TalentSpellID: 16187, TalentPosition: 1,
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "greater_blessing_of_kings", Scope: buffmanifest.ScopeIndividual,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindStatPct,
				Go: "SynthBlessingOfKings", Name: "Blessing of Kings",
				Pet: buffmanifest.PetStripWhenSummonedLate,
			},
			SpellID: 20217, CastSpellID: 20217, DurationMs: 3600000, Supported: true,
			Effects: []ResolvedEffect{
				aura(dbcenums.A_MOD_TOTAL_STAT_PERCENTAGE, 0, 10),
				aura(dbcenums.A_MOD_TOTAL_STAT_PERCENTAGE, 1, 10),
			},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "battle_shout", Scope: buffmanifest.ScopeParty,
				Proto: buffmanifest.ProtoTristate, Kind: buffmanifest.KindStatFlat,
				Go: "SynthBattleShout", Name: "Battle Shout", Category: "SynthBattleShout",
				SingleAura: true, Driver: true,
			},
			SpellID: 25289, CastSpellID: 25289, DurationMs: 180000, Supported: true,
			Effects: []ResolvedEffect{aura(dbcenums.A_MOD_ATTACK_POWER, 0, 139)},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "devotion_aura", Scope: buffmanifest.ScopeParty,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindResistance,
				Go: "SynthDevotionAura", Name: "Devotion Aura", Category: "DevotionAura",
				SharedCategory: "SynthPaladinAura", SingleAura: true,
			},
			SpellID: 10293, CastSpellID: 10293, DurationMs: 600000, Supported: true,
			Effects:       []ResolvedEffect{aura(dbcenums.A_MOD_RESISTANCE, 1, 735)},
			TalentRanks:   2,
			TalentApplies: buffmanifest.TalentScalesDuration,
			TalentSpellID: 20140, TalentPosition: 2,
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "frost_resistance_aura", Scope: buffmanifest.ScopeRaid,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindResistance,
				Go: "SynthFrostResistanceAura", Name: "Frost Resistance Aura",
				Category: "FrostResistanceAura", SharedCategory: "SynthPaladinAura", SingleAura: true,
			},
			SpellID: 19898, CastSpellID: 19898, Supported: true,
			Effects: []ResolvedEffect{aura(dbcenums.A_MOD_RESISTANCE, 16, 60)},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "thunder_clap", Scope: buffmanifest.ScopeDebuff,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindDebuffAtkSpeed,
				Go: "SynthThunderClap", Name: "Thunder Clap", Category: "AtkSpdReduction",
			},
			SpellID: 11581, CastSpellID: 11581, DurationMs: 30000, Supported: true,
			Effects:        []ResolvedEffect{aura(dbcenums.A_MOD_MELEE_HASTE_3, 0, -20)},
			TalentRanks:    2,
			TalentApplies:  buffmanifest.TalentScalesValue,
			TalentOnPseudo: true,
			TalentSpellID:  12287, TalentPosition: 1,
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "innervates", Scope: buffmanifest.ScopeIndividual,
				Proto: buffmanifest.ProtoInt32, Kind: buffmanifest.KindExternalCD,
				Go: "SynthInnervates", Name: "Innervate", Label: "Innervates",
				Category: "Innervate",
			},
			SpellID: 29166, CastSpellID: 29166, DurationMs: 20000, CooldownMs: 360000, Supported: true,
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "power_infusions", Scope: buffmanifest.ScopeIndividual,
				Proto: buffmanifest.ProtoInt32, Kind: buffmanifest.KindExternalCD,
				Go: "SynthPowerInfusions", Name: "Power Infusion", Label: "Power Infusions",
				Category: "PowerInfusion",
			},
			SpellID: 10060, CastSpellID: 10060, DurationMs: 15000, CooldownMs: 180000, Supported: true,
			Effects: []ResolvedEffect{
				aura(dbcenums.A_MOD_DAMAGE_PERCENT_DONE, 126, 20),
				aura(dbcenums.A_MOD_HEALING_DONE_PERCENT, 0, 20),
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
				Go: "SynthAtieshMage", Label: "Atiesh - Mage",
			},
			SpellID: 28142, CastSpellID: 28142, Supported: true,
			Effects: []ResolvedEffect{aura(dbcenums.A_MOD_SPELL_CRIT_CHANCE, 0, 2)},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "thorns", Scope: buffmanifest.ScopeRaid,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindDamageShield,
				Go: "SynthThorns", Name: "Thorns", Category: "Thorns",
				Pet: buffmanifest.PetStrip,
			},
			SpellID: 9910, CastSpellID: 9910, DurationMs: 600000, SchoolMask: 8, Supported: true,
			Effects:       []ResolvedEffect{aura(dbcenums.A_DAMAGE_SHIELD, 0, 22)},
			TalentRanks:   2,
			TalentApplies: buffmanifest.TalentScalesValue,
			TalentSpellID: 16836, TalentPosition: 1,
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "sunder_armor", Scope: buffmanifest.ScopeDebuff,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindDebuffStacking,
				Go: "SynthSunderArmor", Name: "Sunder Armor", Category: "MajorArmorReduction",
				SingleAura: true, Driver: true,
			},
			SpellID: 11597, CastSpellID: 11597, DurationMs: 30000, MaxStacks: 5, Supported: true,
			Effects: []ResolvedEffect{aura(dbcenums.A_MOD_RESISTANCE, 1, -450)},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "expose_armor", Scope: buffmanifest.ScopeDebuff,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindDebuffStat,
				Go: "SynthExposeArmor", Name: "Expose Armor", Category: "MajorArmorReduction",
				SingleAura: true,
			},
			SpellID: 11198, CastSpellID: 11198, DurationMs: 30000, Supported: true,
			Effects: []ResolvedEffect{{Effect: dbcenums.E_APPLY_AURA, Aura: dbcenums.A_MOD_RESISTANCE, Misc: 1, PerResource: -450}},
		},
		{
			BuffSpec: buffmanifest.BuffSpec{
				Field: "curse_of_elements", Scope: buffmanifest.ScopeDebuff,
				Proto: buffmanifest.ProtoBool, Kind: buffmanifest.KindDebuffDamageTaken,
				Go: "SynthCurseOfElements", Name: "Curse of the Elements", Category: "CurseOfElements",
				SingleAura: true,
			},
			SpellID: 1311680, CastSpellID: 1311680, DurationMs: 300000, Supported: true,
			Effects: []ResolvedEffect{
				aura(dbcenums.A_MOD_RESISTANCE, 124, -75),
				aura(dbcenums.A_MOD_DAMAGE_PERCENT_TAKEN, 126, 10),
			},
		},
	}
	for i := range rows {
		if rows[i].Supported {
			setEffectRefs(&rows[i])
			mapEffects(&rows[i])
		}
	}
	return rows
}

// The sim holds flat damage taken in a physical field and a spell field, and
// the client states the aura with a school mask. A mask naming some spell
// schools and not others - Judgement of the Crusader is holy alone - fits
// neither, and the spell field would raise what every school does to the target.
func TestSchoolMaskedDamageTakenHasNoFieldToLandOn(t *testing.T) {
	physical := ResolvedEffect{Aura: dbcenums.A_MOD_DAMAGE_TAKEN, Misc: 1, Value: 8}
	if mods, ok := pseudoModsOf(physical); !ok || mods[0].Kind != "BonusPhysicalDamageTaken" || mods[0].Amount != 8 {
		t.Errorf("a physical mask maps to %v, ok %v; want 8 BonusPhysicalDamageTaken", mods, ok)
	}

	everySchool := ResolvedEffect{Aura: dbcenums.A_MOD_DAMAGE_TAKEN, Misc: 126, Value: 40}
	if mods, ok := pseudoModsOf(everySchool); !ok || mods[0].Kind != "BonusSpellDamageTaken" || mods[0].Amount != 40 {
		t.Errorf("a mask of every spell school maps to %v, ok %v; want 40 BonusSpellDamageTaken", mods, ok)
	}

	holy := ResolvedEffect{Aura: dbcenums.A_MOD_DAMAGE_TAKEN, Misc: 2, Value: 161}
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

// Builds sim/core/buffs with the rendered files overlaid onto it, which is the
// only check that a generated constructor still names an identifier the support
// API declares. The apply functions are renamed because the real generated files
// already declare them, and the drivers the rows call are stubbed here the way
// sim/core/buffs/drivers.go would declare them.
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
		overlay[filepath.Join(root, "sim", "core", "buffs", "zz_synthetic_"+filepath.Base(name))] = path
	}

	drivers := filepath.Join(dir, "drivers.go")
	if err := os.WriteFile(drivers, []byte("package buffs\n\n"+
		"import (\n\t\"github.com/wowsims/forever/sim/core\"\n\t\"github.com/wowsims/forever/sim/core/proto\"\n)\n\n"+
		"func driveSynthInnervates(char *core.Character, individual *proto.IndividualBuffs) {\n"+
		"\tcore.NewGeneratedExternalCD(char, SynthInnervatesAura(&char.Unit, false, 0),"+
		" core.GeneratedExternalCD{NumSources: individual.Innervates,"+
		" Cooldown: SynthInnervatesCooldown(), Type: core.CooldownTypeMana})\n}\n\n"+
		"func driveSynthBattleShout(char *core.Character, _ *proto.PartyBuffs) {\n"+
		"\tcore.ApplyFixedShoutAura(char, SynthBattleShoutAura(&char.Unit, false, 0),"+
		" SynthBattleShoutCategory)\n}\n\n"+
		"func driveSynthSunderArmor(target *core.Unit, _ *proto.Debuffs, _ *proto.Raid) {\n"+
		"\tcore.MakePermanent(SynthSunderArmorAura(target, false, 0))\n}\n\n"+
		"func driveSynthAtieshMage(char *core.Character, party *proto.PartyBuffs) {\n"+
		"\tcore.MakePermanent(SynthAtieshMageAura(&char.Unit, false, 0,"+
		" float64(party.AtieshMage)))\n}\n\n"+
		"func driveSynthPowerInfusions(char *core.Character, individual *proto.IndividualBuffs) {\n"+
		"\tcore.NewGeneratedExternalCD(char, SynthPowerInfusionsAura(&char.Unit, false, 0),"+
		" core.GeneratedExternalCD{NumSources: individual.PowerInfusions,"+
		" Cooldown: SynthPowerInfusionsCooldown(), Type: core.CooldownTypeDPS})\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	overlay[filepath.Join(root, "sim", "core", "buffs", "zz_synthetic_drivers.go")] = drivers

	overlayPath := filepath.Join(dir, "overlay.json")
	if err := os.WriteFile(overlayPath, []byte(overlayJSON(overlay)), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(goTool, "build", "-overlay", overlayPath, "./sim/core/buffs/")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("the generated constructors do not compile against sim/core/buffs:\n%s", out)
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
