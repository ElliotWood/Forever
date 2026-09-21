package spelldata

import (
	"slices"
	"strconv"
	"strings"

	"github.com/wowsims/forever/sim/core"
)

// An addition to the resolved config for what the client does not state: the proc mask, the metrics
// the sim keeps and the flags a row cannot justify.
type SpellOpt func(*core.SpellConfig, *Spell)

// SpellCategories.StartRecoveryCategory of the global cooldown. It is the only category the store
// carries - 1763 rows state 133 and the rest state nothing - so a row outside it spends no GCD.
const globalCooldownCategory int16 = 133

// SpellPower.PowerType. Rage is powerTypeRage in spell.go, where the 0-1000 bar is converted.
const (
	powerTypeMana   int8 = 0
	powerTypeFocus  int8 = 2
	powerTypeEnergy int8 = 3
)

// ImplicitTarget_0 values that name a friendly unit, by the client's Targets enum: TARGET_UNIT_CASTER,
// _PET, _CASTER_AREA_PARTY, _TARGET_ALLY, _SRC_AREA_ALLY, _DEST_AREA_ALLY, _SRC_AREA_PARTY,
// _DEST_AREA_PARTY, _TARGET_PARTY, _TARGET_CHAINHEAL_ALLY, _CASTER_AREA_RAID and _TARGET_RAID.
var helpfulTargets = []uint8{1, 5, 20, 21, 30, 31, 33, 34, 35, 45, 56, 57}

// What the client states about a spell, as the fields core registers it through. The caller adds the
// proc mask, ApplyEffects and anything the client does not carry to the returned value before handing
// it to RegisterSpell: the resolver fills the row's own fields and nothing else.
//
// The unit is needed for the cooldown timers, so a config is built where the sim has a character,
// not at package init. MaxTargets has no SpellConfig field; a caller that needs the row's cap reads
// s.MaxTargets itself.
func SpellConfig(unit *core.Unit, s *Spell, opts ...SpellOpt) core.SpellConfig {
	config := core.SpellConfig{
		ActionID:     core.ActionID{SpellID: s.ID},
		Rank:         rankOf(s),
		SpellSchool:  s.SpellSchool(),
		DefenseType:  s.DefenseTypeCore(),
		ClassFlags:   s.ClassFlags,
		Flags:        rowFlags(s),
		MissileSpeed: float64(s.Speed),
		MinRange:     float64(s.MinRange),
		MaxRange:     float64(s.MaxRange),
		Cast:         castConfig(unit, s),
	}
	applyCost(&config, s)

	// The row first, the caller's options on top, so an option sees what the row filled.
	for _, opt := range opts {
		opt(&config, s)
	}
	return config
}

// The proc mask, the melee metrics bucket and the multipliers an ability needs, for a physical spell.
func Melee(mask core.ProcMask) SpellOpt {
	return func(config *core.SpellConfig, _ *Spell) {
		config.ProcMask = mask
		config.Flags |= core.SpellFlagMeleeMetrics | core.SpellFlagAPL
		config.DamageMultiplier = 1
		config.ThreatMultiplier = 1
		config.Cast.IgnoreHaste = true
	}
}

// The same for a spell that scales with spell power, whose share of it the row states on the effect
// that deals the damage, or heals where the spell has no damaging effect.
func Magic(mask core.ProcMask) SpellOpt {
	return func(config *core.SpellConfig, s *Spell) {
		config.ProcMask = mask
		config.Flags |= core.SpellFlagAPL
		config.DamageMultiplier = 1
		config.ThreatMultiplier = 1
		config.BonusCoefficient = spellPowerCoeff(s)
	}
}

// A spell another spell or an aura casts: it is never in a rotation and does not feed on-cast
// effects. It keeps its metrics, the way the warrior's Deep Wounds, Whirlwind off-hand and Blood
// Craze sub-spells do.
func Proc() SpellOpt {
	return func(config *core.SpellConfig, _ *Spell) {
		config.Flags |= core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete
	}
}

// Flags the client does not state, such as the sim's own metrics and rotation flags.
func Flags(flags core.SpellFlag) SpellOpt {
	return func(config *core.SpellConfig, _ *Spell) {
		config.Flags |= flags
	}
}

// Splits one spell id into several actions, for a spell the sim registers more than once.
func Tag(tag int32) SpellOpt {
	return func(config *core.SpellConfig, _ *Spell) {
		config.ActionID.Tag = tag
	}
}

func spellPowerCoeff(s *Spell) float64 {
	if e := s.DamageEffect(); e != NilEffect {
		return e.Coeff()
	}
	return s.HealEffect().Coeff()
}

// Spell.NameSubtext_lang, which is "Rank 4" on a ranked spell and a word like "Passive" or
// "Shapeshift" on a handful of others. A spell the client shows no rank on answers 0, which is what
// the APL UI reads as "this spell has no ranks".
func rankOf(s *Spell) int32 {
	digits, ranked := strings.CutPrefix(s.Rank, "Rank ")
	if !ranked {
		return 0
	}
	rank, err := strconv.Atoi(digits)
	if err != nil {
		return 0
	}
	return int32(rank)
}

// The flags the row's attributes and targets state. Everything else is the caller's: a flag the
// client does not carry is not invented here.
func rowFlags(s *Spell) core.SpellFlag {
	var flags core.SpellFlag
	if s.IsPassive() {
		flags |= core.SpellFlagPassiveSpell
	}
	if s.IsChanneled() {
		flags |= core.SpellFlagChanneled
	}
	if s.SuppressesWeaponProcs() {
		flags |= core.SpellFlagSuppressWeaponProcs
	}
	// Helpful decides who the APL casts the spell on, so it follows the first effect's target. An
	// attack whose first effect is a self side-effect reads as helpful here and the caller clears it.
	if slices.Contains(helpfulTargets, s.EffectN(1).Target[0]) {
		flags |= core.SpellFlagHelpful
	}
	return flags
}

func castConfig(unit *core.Unit, s *Spell) core.CastConfig {
	cast := core.CastConfig{
		DefaultCast: core.Cast{CastTime: s.CastTime()},
		IgnoreHaste: s.DefenseTypeCore() != core.DefenseTypeMagic,
	}

	if s.StartRecoveryCategory == globalCooldownCategory {
		cast.DefaultCast.GCD = s.GCD()
	}

	// The category cooldown is the one a whole category of spells shares, so it runs off the unit's
	// timer for that category rather than off a timer of this spell's own.
	switch {
	case s.CooldownMs > 0 && s.CategoryCooldownMs > 0:
		cast.CD = core.Cooldown{Timer: unit.NewTimer(), Duration: s.Cooldown()}
		cast.SharedCD = core.Cooldown{Timer: unit.CategoryTimer(int32(s.Category)), Duration: s.CategoryCooldown()}
	case s.CooldownMs > 0:
		cast.CD = core.Cooldown{Timer: unit.NewTimer(), Duration: s.Cooldown()}
	case s.CategoryCooldownMs > 0:
		cast.CD = core.Cooldown{Timer: unit.CategoryTimer(int32(s.Category)), Duration: s.CategoryCooldown()}
	}

	return cast
}

// A cast that spends a resource is not an empty one even where the row states no GCD and no cast
// time: core reads an empty DefaultCast as a proc, which skips the cost. A row that spends nothing
// is left empty on purpose, so a passive or a proc spell keeps the cast path core gives those.
func applyCost(config *core.SpellConfig, s *Spell) {
	if len(s.Powers) == 0 {
		return
	}

	powerType := s.Powers[0].Type
	cost := int32(s.PowerCost(powerType))

	switch powerType {
	case powerTypeMana:
		config.ManaCost = core.ManaCostOptions{FlatCost: cost}
		if s.Powers[0].CostPct > 0 {
			config.ManaCost.BaseCostPercent = float64(s.Powers[0].CostPct)
		}
	case powerTypeRage:
		config.RageCost = core.RageCostOptions{Cost: cost, Refund: s.MissRefund()}
	case powerTypeEnergy:
		config.EnergyCost = core.EnergyCostOptions{Cost: cost, Refund: s.MissRefund()}
	case powerTypeFocus:
		config.FocusCost = core.FocusCostOptions{Cost: cost, Refund: s.MissRefund()}
	}

	if cost > 0 && config.Cast.DefaultCast.GCD == 0 && config.Cast.DefaultCast.CastTime == 0 {
		config.Cast.DefaultCast.NonEmpty = true
	}
}
