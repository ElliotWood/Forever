package buffs

import (
	"fmt"
	"math"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// What applying an aura does to a fresh unit, in every respect a raid buff is compared on: how it is
// registered, what it bids and at what, and what its stats and pseudo-stats are worth at each stack.
type auraShape struct {
	Label      string
	ActionID   core.ActionID
	Tag        string
	Duration   time.Duration
	MaxStacks  int32
	BuildPhase core.CharacterBuildPhase
	Bids       []string
	Stacks     []unitDelta
	Shield     string
}

type unitDelta struct {
	Stats  map[string]float64
	Pseudo map[string]float64
}

// Every stat starts above zero, so that a multiplier has something to move.
const shapeBaseStat = 1000.0

func shapeCharacter() *core.Character {
	character := core.NewCharacter(&core.Party{}, 0, &proto.Player{
		Name:      "Buff Tester",
		Race:      proto.Race_RaceOrc,
		Class:     proto.Class_ClassWarrior,
		Spec:      &proto.Player_ProtectionWarrior{ProtectionWarrior: &proto.ProtectionWarrior{}},
		Equipment: &proto.EquipmentSpec{Items: []*proto.ItemSpec{}},
	})
	character.Env = &core.Environment{MeasuringStats: true}
	character.AddStats(shapeBaseStats())
	return &character
}

func shapeTarget() *core.Unit {
	target := core.NewTarget(&proto.Target{Level: 63}, 0)
	target.Env = &core.Environment{MeasuringStats: true}
	target.AddStats(shapeBaseStats())
	return &target.Unit
}

func shapeBaseStats() stats.Stats {
	var s stats.Stats
	for i := range s {
		s[i] = shapeBaseStat
	}
	return s
}

// Registers the aura on a fresh unit, applies it at every stack it can hold, and checks that
// expiring it hands every stat and pseudo-stat back.
func shapeOf(t *testing.T, onTarget bool, build func(*core.Unit) *core.Aura) auraShape {
	t.Helper()
	unit := &shapeCharacter().Unit
	if onTarget {
		unit = shapeTarget()
	}
	sim := &core.Simulation{}

	measure := func() (stats.Stats, stats.PseudoStats) {
		return unit.SortAndApplyStatDependencies(unit.GetStats()), unit.PseudoStats
	}
	beforeStats, beforePseudo := measure()

	aura := build(unit)
	shape := auraShape{
		Label:      aura.Label,
		ActionID:   aura.ActionID,
		Tag:        aura.Tag,
		Duration:   aura.Duration,
		MaxStacks:  aura.MaxStacks,
		BuildPhase: aura.BuildPhase,
	}

	aura.Activate(sim)
	for stacks := int32(1); ; stacks++ {
		if aura.MaxStacks > 0 {
			aura.SetStacks(sim, stacks)
		}
		afterStats, afterPseudo := measure()
		shape.Stacks = append(shape.Stacks, unitDelta{statsDiff(beforeStats, afterStats), pseudoDiff(beforePseudo, afterPseudo)})
		if stacks >= aura.MaxStacks {
			break
		}
	}

	for _, ee := range aura.ExclusiveEffects {
		shape.Bids = append(shape.Bids, fmt.Sprintf("%s single=%v priority=%v", ee.Category.Name, ee.Category.SingleAura, round6(ee.Priority)))
	}
	slices.Sort(shape.Bids)

	for _, spell := range unit.Spellbook {
		if spell.ActionID.SpellID == aura.ActionID.SpellID {
			shape.Shield = fmt.Sprintf("%v school %v", spell.ActionID, spell.SpellSchool)
		}
	}

	aura.Deactivate(sim)
	afterStats, afterPseudo := measure()
	if d := statsDiff(beforeStats, afterStats); d != nil {
		t.Errorf("%s leaves %v behind on expiry", aura.Label, d)
	}
	if d := pseudoDiff(beforePseudo, afterPseudo); d != nil {
		t.Errorf("%s leaves %v behind on expiry", aura.Label, d)
	}
	return shape
}

func round6(f float64) float64 {
	return math.Round(f*1e6) / 1e6
}

func statsDiff(before, after stats.Stats) map[string]float64 {
	diff := map[string]float64{}
	for i := range before {
		if d := round6(after[i] - before[i]); d != 0 {
			diff[stats.Stat(i).StatName()] = d
		}
	}
	if len(diff) == 0 {
		return nil
	}
	return diff
}

// Every float field of PseudoStats, the per-school arrays included.
func pseudoDiff(before, after stats.PseudoStats) map[string]float64 {
	diff := map[string]float64{}
	bv, av := reflect.ValueOf(before), reflect.ValueOf(after)
	for i := 0; i < bv.NumField(); i++ {
		name := bv.Type().Field(i).Name
		bf, af := bv.Field(i), av.Field(i)
		switch bf.Kind() {
		case reflect.Float64:
			if d := round6(af.Float() - bf.Float()); d != 0 {
				diff[name] = d
			}
		case reflect.Array:
			for j := 0; j < bf.Len(); j++ {
				if bf.Index(j).Kind() != reflect.Float64 {
					continue
				}
				if d := round6(af.Index(j).Float() - bf.Index(j).Float()); d != 0 {
					diff[fmt.Sprintf("%s[%d]", name, j)] = d
				}
			}
		}
	}
	if len(diff) == 0 {
		return nil
	}
	return diff
}

// A Meta's constructor builds what the constructor generated for the same row builds, on both
// copies: the player's own and the external caster's.
func TestMetaConstructorsBuildWhatTheGeneratedRowsBuild(t *testing.T) {
	find := spelldata.MustFind
	type ctor func(unit *core.Unit, isPlayer bool, talentPoints int32) *core.Aura

	itemCount := func(m *Meta) ctor {
		return func(unit *core.Unit, isPlayer bool, _ int32) *core.Aura {
			return newItemCountBuff(unit, m, isPlayer, 2)
		}
	}
	generatedItemCount := func(f func(*core.Unit, bool, int32, float64) *core.Aura) ctor {
		return func(unit *core.Unit, isPlayer bool, talentPoints int32) *core.Aura {
			return f(unit, isPlayer, talentPoints, 2)
		}
	}
	buff := func(m *Meta) ctor {
		return func(unit *core.Unit, isPlayer bool, tp int32) *core.Aura { return newBuff(unit, m, isPlayer, tp) }
	}
	debuff := func(m *Meta) ctor {
		return func(unit *core.Unit, isPlayer bool, tp int32) *core.Aura { return newDebuff(unit, m, isPlayer, tp) }
	}
	shield := func(m *Meta) ctor {
		return func(unit *core.Unit, isPlayer bool, tp int32) *core.Aura {
			return newDamageShield(unit, m, isPlayer, tp)
		}
	}

	cases := []struct {
		name      string
		onTarget  bool
		ranks     int32
		meta      ctor
		generated ctor
	}{
		{"BloodPact", false, 0,
			buff(&Meta{Label: "Blood Pact", Spell: find(11767)}), BloodPactAura},
		{"BattleShout", false, 0,
			buff(&Meta{Label: "Battle Shout", Spell: find(25289), Category: "BattleShout", SingleAura: true}),
			BattleShoutAura},
		{"LeaderOfThePack", false, 0,
			buff(&Meta{Label: "Leader of the Pack", Spell: find(24932), Category: "DruidCritAura", SingleAura: true}),
			LeaderOfThePackAura},
		{"ManaSpringTotem", false, 5,
			buff(&Meta{Label: "Mana Spring Totem", Spell: find(10494), Category: "ManaSpringTotem",
				Talent: spelldata.Talent(16187, 5), TalentEffect: 1}),
			ManaSpringTotemAura},
		{"ArcaneBrilliance", false, 0,
			buff(&Meta{Label: "Arcane Brilliance", Spell: find(23028), Category: "StatBuff"}), ArcaneBrillianceAura},
		{"GreaterBlessingOfKings", false, 0,
			buff(&Meta{Label: "Greater Blessing of Kings", Spell: find(25898)}), GreaterBlessingOfKingsAura},
		{"GiftOfTheWild", false, 0,
			buff(&Meta{Label: "Gift of the Wild", Spell: find(21850)}), GiftOfTheWildAura},
		{"FireResistanceAura", false, 0,
			buff(&Meta{Label: "Fire Resistance Aura", Spell: find(19900), Category: "FireResistanceAura",
				SharedCategory: "PaladinAura", SingleAura: true, SkipAuras: paladinAuraSkips}),
			FireResistanceAuraAura},
		{"ConcentrationAura", false, 0,
			buff(&Meta{Label: "Concentration Aura", Spell: find(19746), Category: "ConcentrationAura",
				SharedCategory: "PaladinAura", SingleAura: true, SkipAuras: paladinAuraSkips}),
			ConcentrationAuraAura},
		{"GreaterBlessingOfSalvation", false, 0,
			buff(&Meta{Label: "Greater Blessing of Salvation", Spell: find(25895)}), GreaterBlessingOfSalvationAura},
		{"AtieshWarlock", false, 0,
			itemCount(&Meta{Label: "Atiesh - Warlock", Spell: find(28143)}), generatedItemCount(AtieshWarlockAura)},
		{"RetributionAura", false, 0,
			shield(&Meta{Label: "Retribution Aura", Spell: find(10301), Category: "RetributionAura",
				SharedCategory: "PaladinAura", SingleAura: true, SkipAuras: paladinAuraSkips}),
			RetributionAuraAura},
		{"Thorns", false, 0,
			shield(&Meta{Label: "Thorns", Spell: find(9910), Category: "Thorns", SingleAura: true}), ThornsAura},
		{"CurseOfElements", true, 0,
			debuff(&Meta{Label: "Curse of the Elements", Spell: find(1311680), Category: "CurseOfElements", SingleAura: true}),
			CurseOfElementsAura},
		{"SunderArmor", true, 0,
			debuff(&Meta{Label: "Sunder Armor", Spell: find(11597), Category: "MajorArmorReduction", SingleAura: true}),
			SunderArmorAura},
		{"ExposeArmor", true, 0,
			debuff(&Meta{Label: "Expose Armor", Spell: find(11198), Category: "MajorArmorReduction", SingleAura: true,
				FullComboPoints: true}),
			ExposeArmorAura},
		{"ThunderClap", true, 0,
			debuff(&Meta{Label: "Thunder Clap", Spell: find(11581), Category: "AtkSpdReduction"}), ThunderClapAura},
		{"InsectSwarm", true, 0,
			debuff(&Meta{Label: "Insect Swarm", Spell: find(24977)}), InsectSwarmAura},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for points := int32(0); points <= c.ranks; points++ {
				for _, isPlayer := range []bool{true, false} {
					got := shapeOf(t, c.onTarget, func(u *core.Unit) *core.Aura { return c.meta(u, isPlayer, points) })
					want := shapeOf(t, c.onTarget, func(u *core.Unit) *core.Aura { return c.generated(u, isPlayer, points) })

					// The generated debuffs register without a tag; a Meta's aura carries its category.
					if c.onTarget {
						want.Tag = got.Tag
					}
					if !reflect.DeepEqual(got, want) {
						t.Errorf("isPlayer %v, %d points:\n got %+v\nwant %+v", isPlayer, points, got, want)
					}
				}
			}
		})
	}
}

// A single-aura buff bids its whole value under its category, which is what AddGeneratedFlatBonus
// raises: the improved copy grants the bonus while it holds the category.
func TestMetaBuffTakesAFlatBonus(t *testing.T) {
	meta := &Meta{Label: "Battle Shout", Spell: spelldata.MustFind(25289), Category: "BattleShout", SingleAura: true}
	character := shapeCharacter()
	sim := &core.Simulation{}
	before := character.GetStat(stats.AttackPower)

	aura := newBuff(&character.Unit, meta, false, 0)
	core.AddGeneratedFlatBonus(aura, stats.AttackPower, meta.Value(0), BattleShoutT2Bonus)
	aura.Activate(sim)

	if got, want := character.GetStat(stats.AttackPower)-before, meta.Value(0)+BattleShoutT2Bonus; got != want {
		t.Errorf("improved Battle Shout grants %v attack power, want %v", got, want)
	}
}
