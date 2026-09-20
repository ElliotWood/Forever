package core

// The API sim/core/buffs_auto_gen.go and sim/core/debuffs_auto_gen.go are
// generated against. Everything the generator emits is a data literal plus one
// call into this file, so a generated constructor cannot go wrong in a way the
// compiler does not catch.

import (
	"time"

	"github.com/wowsims/forever/sim/core/stats"
)

// PseudoStatKind names the PseudoStats field a generated buff modifies.
type PseudoStatKind int

const (
	PseudoStatThreatMultiplier PseudoStatKind = iota
	PseudoStatDamageDealtMultiplier
	PseudoStatDamageTakenMultiplier
	PseudoStatSchoolDamageTakenMultiplier
	PseudoStatMeleeSpeedMultiplier
	PseudoStatPushbackChance
	PseudoStatBonusPhysicalDamageTaken
	PseudoStatBonusSpellDamageTaken
	PseudoStatBonusAttackPower
	PseudoStatBonusRangedAttackPower
)

type PseudoConfig struct {
	Kind             PseudoStatKind
	Amount           float64
	IsMultiplicative bool

	// The client's school bits (1 physical, 2 holy, 4 fire, 8 nature, 16 frost,
	// 32 shadow, 64 arcane), read by the school-specific kinds only.
	SchoolMask int32
}

type GeneratedBuff struct {
	Label      string
	ActionID   ActionID
	Duration   time.Duration
	MaxStacks  int32
	Category   string
	SingleAura bool
	Stats      []StatConfig
	Pseudo     []PseudoConfig
	IsPlayer   bool
}

// The aura a generated buff registers on the player. Tag -1 is the external
// caster's copy, which the character build phase has to see so that stat
// dependencies are computed with it; tag 0 is the player's own and is applied
// during the fight.
func newGeneratedStatAura(unit *Unit, config GeneratedBuff) *Aura {
	aura := unit.GetOrRegisterAura(Aura{
		Label:      config.Label,
		ActionID:   config.ActionID,
		Duration:   TernaryDuration(config.Duration > 0, config.Duration, NeverExpires),
		MaxStacks:  config.MaxStacks,
		BuildPhase: Ternary(config.ActionID.Tag == -1, CharacterBuildPhaseBuffs, CharacterBuildPhaseNone),
	})

	if config.Category != "" {
		registerExlusiveEffects(aura, config.Stats, config.Category)
		if config.SingleAura {
			aura.NewExclusiveEffect(config.Category, true, ExclusiveEffect{
				Priority: generatedPriority(config),
			})
		}
	} else {
		registerStatEffect(aura, config.Stats)
	}

	attachGeneratedPseudoStats(aura, config.Pseudo)
	return aura
}

// The aura a generated debuff registers on the target. A stacking debuff prices
// its exclusive effect by what it is worth at full stacks, which is what lets
// the strongest of several armor reductions win before any of them has stacked.
func newGeneratedDebuff(target *Unit, config GeneratedBuff) *Aura {
	aura := target.GetOrRegisterAura(Aura{
		Label:     config.Label,
		ActionID:  config.ActionID,
		Duration:  TernaryDuration(config.Duration > 0, config.Duration, NeverExpires),
		MaxStacks: config.MaxStacks,
	})

	if config.Category != "" {
		registerExlusiveEffects(aura, config.Stats, config.Category)
		if config.SingleAura {
			aura.NewExclusiveEffect(config.Category, true, ExclusiveEffect{
				Priority: generatedPriority(config),
			})
		}
	} else {
		registerStatEffect(aura, config.Stats)
	}

	attachGeneratedPseudoStats(aura, config.Pseudo)
	return aura
}

// The damage a generated damage shield deals back to whoever lands a melee hit.
func newGeneratedDamageShield(unit *Unit, config GeneratedBuff, school SpellSchool, damage float64) *Aura {
	procSpell := unit.RegisterSpell(SpellConfig{
		ActionID:    config.ActionID.WithTag(config.ActionID.Tag + 2),
		SpellSchool: school,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagBinary | SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHit)
		},
	})

	aura := unit.GetOrRegisterAura(Aura{
		Label:      config.Label,
		ActionID:   config.ActionID,
		Duration:   TernaryDuration(config.Duration > 0, config.Duration, NeverExpires),
		BuildPhase: Ternary(config.ActionID.Tag == -1, CharacterBuildPhaseBuffs, CharacterBuildPhaseNone),
	}).AttachProcTrigger(ProcTrigger{
		Name:     config.Label + " Damage",
		Callback: CallbackOnSpellHitTaken,
		Outcome:  OutcomeLanded,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			if spell.SpellSchool.Matches(SpellSchoolPhysical) {
				procSpell.Cast(sim, spell.Unit)
			}
		},
	})

	if config.Category != "" {
		aura.NewExclusiveEffect(config.Category, true, ExclusiveEffect{Priority: generatedPriority(config)})
	}
	return aura
}

// A buff other players cast on this one on their own cooldown, approximated by
// numSources casters taking turns.
func newGeneratedExternalCD(char *Character, config GeneratedBuff, numSources int32, cooldown time.Duration, shouldActivate CooldownActivationCondition) {
	if numSources == 0 {
		return
	}

	aura := newGeneratedStatAura(&char.Unit, config)
	registerExternalConsecutiveCDApproximation(char, externalConsecutiveCDApproximation{
		ActionID:         config.ActionID,
		AuraTag:          config.Label,
		CooldownPriority: CooldownPriorityDefault,
		Type:             CooldownTypeDPS,
		AuraDuration:     config.Duration,
		AuraCD:           cooldown,
		ShouldActivate:   shouldActivate,
		AddAura:          func(sim *Simulation, _ *Character) { aura.Activate(sim) },
		RelatedSelfBuff:  aura,
	}, numSources)
}

// The priority an exclusive effect of a whole aura carries: what the buff is
// worth at full stacks, so the larger of two mutually exclusive versions wins
// before either has stacked up.
func generatedPriority(config GeneratedBuff) float64 {
	if len(config.Stats) == 0 {
		return 0
	}
	priority := config.Stats[0].Amount
	if config.MaxStacks > 0 {
		priority *= float64(config.MaxStacks)
	}
	return priority
}

func attachGeneratedPseudoStats(aura *Aura, configs []PseudoConfig) {
	for _, config := range configs {
		for _, field := range generatedPseudoStatFields(aura.Unit, config) {
			if config.IsMultiplicative {
				aura.AttachMultiplicativePseudoStatBuff(field, config.Amount)
			} else {
				aura.AttachAdditivePseudoStatBuff(field, config.Amount)
			}
		}
	}
}

func generatedPseudoStatFields(unit *Unit, config PseudoConfig) []*float64 {
	switch config.Kind {
	case PseudoStatThreatMultiplier:
		return []*float64{&unit.PseudoStats.ThreatMultiplier}
	case PseudoStatDamageDealtMultiplier:
		return []*float64{&unit.PseudoStats.DamageDealtMultiplier}
	case PseudoStatDamageTakenMultiplier:
		return []*float64{&unit.PseudoStats.DamageTakenMultiplier}
	case PseudoStatMeleeSpeedMultiplier:
		return []*float64{&unit.PseudoStats.MeleeSpeedMultiplier}
	case PseudoStatPushbackChance:
		return []*float64{&unit.PseudoStats.PushbackChance}
	case PseudoStatBonusPhysicalDamageTaken:
		return []*float64{&unit.PseudoStats.BonusPhysicalDamageTaken}
	case PseudoStatBonusSpellDamageTaken:
		return []*float64{&unit.PseudoStats.BonusSpellDamageTaken}
	case PseudoStatBonusAttackPower:
		return []*float64{&unit.PseudoStats.BonusAttackPower}
	case PseudoStatBonusRangedAttackPower:
		return []*float64{&unit.PseudoStats.BonusRangedAttackPower}
	case PseudoStatSchoolDamageTakenMultiplier:
		var fields []*float64
		for _, school := range generatedSchoolIndexes(config.SchoolMask) {
			fields = append(fields, &unit.PseudoStats.SchoolDamageTakenMultiplier[school])
		}
		return fields
	}
	return nil
}

// The sim's school indexes a client school mask names.
func generatedSchoolIndexes(mask int32) []stats.SchoolIndex {
	all := []struct {
		Bit    int32
		School stats.SchoolIndex
	}{
		{1, stats.SchoolIndexPhysical},
		{2, stats.SchoolIndexHoly},
		{4, stats.SchoolIndexFire},
		{8, stats.SchoolIndexNature},
		{16, stats.SchoolIndexFrost},
		{32, stats.SchoolIndexShadow},
		{64, stats.SchoolIndexArcane},
	}

	var schools []stats.SchoolIndex
	for _, entry := range all {
		if mask&entry.Bit != 0 {
			schools = append(schools, entry.School)
		}
	}
	return schools
}
