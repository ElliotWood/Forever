package buffs

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// What the generator states about one manifest row: the spell the aura reads, how the sim names it
// and what it bids in. Which effects are which stats, what they are worth and how they stack is read
// off the row at runtime by spelldata.ParseEffects, with the options Options answers.
type Meta struct {
	// The buff's name in the UI. The aura's label adds " (Player)" or " (External)".
	Label string
	// The row the aura's numbers are read from.
	Spell *spelldata.Spell
	// The cast, for a row whose aura states no duration of its own and whose cooldown is the cast's.
	// Nil where the spell's own row times the buff.
	Cast *spelldata.Spell

	// The aura's own exclusive category, which is also its Tag; "" for none.
	Category string
	// A second category the player's own copy joins without an effect of its own.
	SharedCategory string
	// The category holds one aura at a time, and the aura bids in it as a whole.
	SingleAura bool

	// The improving talent, and the effect of it that prices the buff, counted the way EffectN
	// counts. TalentScalesDuration says the talent prices the duration rather than the amount.
	Talent               spelldata.Ladder
	TalentEffect         int32
	TalentScalesDuration bool

	// The level the amounts are priced at; 0 is core.CharacterLevel.
	Level int32
	// The aura effects the buff leaves out, by aura.
	SkipAuras []dbcenums.EffectAuraType
	// The finisher the raid config puts on the target is cast at full combo points.
	FullComboPoints bool
}

// The parse options the row states, for a caster with these talent points.
func (m *Meta) Options(talentPoints int32) []spelldata.ParseOpt {
	level := m.Level
	if level == 0 {
		level = core.CharacterLevel
	}

	opts := []spelldata.ParseOpt{spelldata.Level(level), spelldata.BuffAuras()}
	if len(m.SkipAuras) > 0 {
		opts = append(opts, spelldata.SkipAuras(m.SkipAuras...))
	}
	if !m.TalentScalesDuration {
		opts = append(opts, spelldata.ScaledBy(m.talentMod(talentPoints)))
	}
	if m.FullComboPoints {
		opts = append(opts, spelldata.FullComboPoints())
	}
	return opts
}

// The first amount the aura applies, in the sim's units: the stat, the multiplier or the pseudo-stat
// change it attaches. A damage shield's amount is the damage it deals.
func (m *Meta) Value(talentPoints int32) float64 {
	if applied := spelldata.DryRun(m.Spell, true, m.Options(talentPoints)...).Applied; len(applied) > 0 {
		return applied[0].Value
	}

	for i := range m.Spell.Effects {
		if e := &m.Spell.Effects[i]; e.Aura == dbcenums.A_DAMAGE_SHIELD {
			value := e.Average(m.level())
			if mod := m.talentMod(talentPoints); mod != spelldata.NilEffect && !m.TalentScalesDuration {
				value = talentScaled(value, mod)
			}
			return value
		}
	}
	return 0
}

func (m *Meta) Duration(talentPoints int32) time.Duration {
	if m.TalentScalesDuration {
		return talentScaledDuration(m.Spell, m.talentMod(talentPoints))
	}
	if m.Spell.DurationMs <= 0 && m.Cast != nil {
		return auraDuration(m.Cast)
	}
	return auraDuration(m.Spell)
}

func (m *Meta) Cooldown() time.Duration {
	if m.Cast != nil {
		return cooldown(m.Cast)
	}
	return cooldown(m.Spell)
}

func (m *Meta) level() int32 {
	if m.Level == 0 {
		return core.CharacterLevel
	}
	return m.Level
}

// The talent's modifier at a rank, NilEffect for an untaken talent or a row no talent prices.
func (m *Meta) talentMod(talentPoints int32) *spelldata.Effect {
	return m.Talent.Rank(talentPoints).EffectN(int(m.TalentEffect))
}

// Tag 0 is the player's own copy and -1 the external caster's, which the character build phase has
// to see so that stat dependencies are computed with it.
func (m *Meta) label(isPlayer bool) string {
	return m.Label + " (" + core.Ternary(isPlayer, "Player", "External") + ")"
}

func (m *Meta) actionID(isPlayer bool) core.ActionID {
	return core.ActionID{SpellID: m.Spell.ID}.WithTag(core.TernaryInt32(isPlayer, 0, -1))
}

// The aura a generated buff registers on a player: labelled and tagged by label and actionID, with
// Tag set to the category, MaxStacks the row's own stacks (not its charges), the external copy in
// CharacterBuildPhaseBuffs, and the effects read by spelldata.ParseEffects with Options plus
// SchoolResistances, and where Category is set, Exclusive(Category, true) for a single-aura one or
// ExclusivePerStat(Category) for any other. Only the player's own copy joins SharedCategory.
func newBuff(unit *core.Unit, m *Meta, isPlayer bool, talentPoints int32) *core.Aura {
	panic("not implemented")
}

// A buff worth its amounts once per item in the party: newBuff with Count(count).
func newItemCountBuff(unit *core.Unit, m *Meta, isPlayer bool, count float64) *core.Aura {
	panic("not implemented")
}

// The aura a generated debuff registers on the target: labelled and tagged like newBuff, never in a
// build phase, with Exclusive(Category, SingleAura) where it names a category and no school
// resistance categories of its own.
func newDebuff(target *core.Unit, m *Meta, isPlayer bool, talentPoints int32) *core.Aura {
	panic("not implemented")
}

// A damage shield, which the parse does not read: core.NewGeneratedDamageShield with the spell's
// school and Value as its damage.
func newDamageShield(unit *core.Unit, m *Meta, isPlayer bool, talentPoints int32) *core.Aura {
	panic("not implemented")
}
