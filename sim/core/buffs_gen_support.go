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

func (kind PseudoStatKind) Name() string {
	switch kind {
	case PseudoStatThreatMultiplier:
		return "ThreatMultiplier"
	case PseudoStatDamageDealtMultiplier:
		return "DamageDealtMultiplier"
	case PseudoStatDamageTakenMultiplier:
		return "DamageTakenMultiplier"
	case PseudoStatSchoolDamageTakenMultiplier:
		return "SchoolDamageTakenMultiplier"
	case PseudoStatMeleeSpeedMultiplier:
		return "MeleeSpeedMultiplier"
	case PseudoStatPushbackChance:
		return "PushbackChance"
	case PseudoStatBonusPhysicalDamageTaken:
		return "BonusPhysicalDamageTaken"
	case PseudoStatBonusSpellDamageTaken:
		return "BonusSpellDamageTaken"
	case PseudoStatBonusAttackPower:
		return "BonusAttackPower"
	case PseudoStatBonusRangedAttackPower:
		return "BonusRangedAttackPower"
	}
	return "PseudoStatKind(unknown)"
}

type PseudoConfig struct {
	Kind             PseudoStatKind
	Amount           float64
	IsMultiplicative bool

	// The client's school bits (1 physical, 2 holy, 4 fire, 8 nature, 16 frost,
	// 32 shadow, 64 arcane), read by the school-specific kinds only.
	SchoolMask int32
}

type GeneratedBuff struct {
	Label     string
	ActionID  ActionID
	Duration  time.Duration
	MaxStacks int32

	// Three separate roles. StatCategory is what the individual stats compete
	// under, so that a paladin's resistance aura and a shaman's totem of the same
	// school do not both apply. Category is the aura's own, which decides whether
	// a second copy of this buff can sit next to it. SharedCategory is one it
	// joins as a member without an effect of its own, which is how the paladin
	// auras exclude each other across schools.
	StatCategory   string
	Category       string
	SharedCategory string
	SingleAura     bool

	Stats    []StatConfig
	Pseudo   []PseudoConfig
	IsPlayer bool
}

// The aura a generated buff registers on the player. Tag -1 is the external
// caster's copy, which the character build phase has to see so that stat
// dependencies are computed with it; tag 0 is the player's own and is applied
// during the fight. The aura's Tag is the category name, so that a driver
// handed the external copy can find the player's own among the unit's auras.
func newGeneratedStatAura(unit *Unit, config GeneratedBuff) *Aura {
	auraConfig := Aura{
		Label:      config.Label,
		Tag:        config.Category,
		ActionID:   config.ActionID,
		Duration:   TernaryDuration(config.Duration > 0, config.Duration, NeverExpires),
		MaxStacks:  config.MaxStacks,
		BuildPhase: Ternary(config.ActionID.Tag == -1, CharacterBuildPhaseBuffs, CharacterBuildPhaseNone),
	}

	perStack := generatedMagnitude(config)
	var effect *ExclusiveEffect
	if generatedHasCategoryEffect(config, false) {
		attachGeneratedStackPricing(&auraConfig, config, &effect, perStack)
	}

	aura := unit.GetOrRegisterAura(auraConfig)
	effect = registerGeneratedEffects(aura, config, perStack, false)
	joinSharedCategory(aura, config)
	return aura
}

// Where a buff's amounts are applied and what it bids for them.
//
// A buff that states a StatCategory competes stat by stat, the way every
// hand-written resistance source does, and its own category is then only the
// bid that decides whether a second copy may sit next to it. A buff whose
// category holds one aura at a time competes as a whole instead, bidding once
// for everything it applies, which is also what every debuff does. Anything
// else applies its amounts outright.
func registerGeneratedEffects(aura *Aura, config GeneratedBuff, perStack float64, bareWhenCategory bool) *ExclusiveEffect {
	// An aura that bids for everything it applies at once keeps its resistances
	// with the rest: two resistance-reducing debuffs exclude each other in that
	// category, and pulling the schools out of it would put both on the target.
	if !bareWhenCategory {
		config.Stats = registerGeneratedSchoolResistances(aura, config.Stats)
	}

	if config.StatCategory != "" {
		registerExlusiveEffects(aura, config.Stats, config.StatCategory)
		attachGeneratedPseudoStats(aura, config)
		if config.Category == "" {
			return nil
		}
		return aura.NewExclusiveEffect(config.Category, config.SingleAura, ExclusiveEffect{Priority: perStack})
	}

	if generatedHasCategoryEffect(config, bareWhenCategory) {
		return registerGeneratedCategoryEffect(aura, config, perStack)
	}

	if config.Category != "" {
		registerExlusiveEffects(aura, config.Stats, config.Category)
	} else {
		registerStatEffect(aura, config.Stats)
	}
	attachGeneratedPseudoStats(aura, config)
	return nil
}

// Every source of a school's resistance bids under that school, whatever else
// the buff does, so that Gift of the Wild's 27 and a resistance aura's 60 do
// not both land on the unit. The stats that are left are returned for the
// caller to place. Armor is not a school and is one of them.
func registerGeneratedSchoolResistances(aura *Aura, config []StatConfig) []StatConfig {
	var rest []StatConfig
	for _, statConfig := range config {
		category := resistanceCategoryOfStat(statConfig.Stat)
		if category == "" || statConfig.IsMultiplicative {
			rest = append(rest, statConfig)
			continue
		}
		makeExclusiveFlatStatBuff(aura, statConfig.Stat, statConfig.Amount, category)
	}
	return rest
}

func resistanceCategoryOfStat(stat stats.Stat) string {
	switch stat {
	case stats.ArcaneResistance:
		return ResistanceCategoryArcane
	case stats.FireResistance:
		return ResistanceCategoryFire
	case stats.FrostResistance:
		return ResistanceCategoryFrost
	case stats.NatureResistance:
		return ResistanceCategoryNature
	case stats.ShadowResistance:
		return ResistanceCategoryShadow
	}
	return ""
}

// Whether the aura will bid under its own category, which is what the stack
// pricing re-prices.
func generatedHasCategoryEffect(config GeneratedBuff, bareWhenCategory bool) bool {
	return config.Category != "" && (bareWhenCategory || config.SingleAura || config.StatCategory != "")
}

// The aura a generated debuff registers on the target. Its exclusive effect
// carries the bare category name, which is the one the hand-written debuffs
// compete in, and holds every stat and pseudo-stat the debuff applies so that
// only the strongest of several armor reductions is on the target at a time.
func newGeneratedDebuff(target *Unit, config GeneratedBuff) *Aura {
	auraConfig := Aura{
		Label:     config.Label,
		ActionID:  config.ActionID,
		Duration:  TernaryDuration(config.Duration > 0, config.Duration, NeverExpires),
		MaxStacks: config.MaxStacks,
	}

	if config.Category == "" {
		aura := target.GetOrRegisterAura(auraConfig)
		registerStatEffect(aura, config.Stats)
		attachGeneratedPseudoStats(aura, config)
		return aura
	}

	perStack := generatedMagnitude(config)
	var effect *ExclusiveEffect
	attachGeneratedStackPricing(&auraConfig, config, &effect, perStack)

	aura := target.GetOrRegisterAura(auraConfig)
	effect = registerGeneratedEffects(aura, config, perStack, true)
	joinSharedCategory(aura, config)
	return aura
}

// A stacking debuff is worth nothing until it has a stack and re-prices itself
// on every one; SetPriority re-applies the amounts as it goes.
func attachGeneratedStackPricing(auraConfig *Aura, config GeneratedBuff, effect **ExclusiveEffect, perStack float64) {
	if config.MaxStacks <= 0 {
		return
	}
	auraConfig.OnStacksChange = func(_ *Aura, sim *Simulation, _ int32, newStacks int32) {
		(*effect).SetPriority(sim, perStack*float64(newStacks))
	}
}

// One exclusive effect under the bare category name - the name the hand-written
// buffs and debuffs compete under - holding every amount the aura applies, so
// that only the strongest member of the category is on the unit at a time.
// Whether that means the other member's aura is kicked off entirely is the
// manifest's SingleAura, which has to agree with what the hand-written members
// of the same category register.
func registerGeneratedCategoryEffect(aura *Aura, config GeneratedBuff, perStack float64) *ExclusiveEffect {
	// A stat the aura multiplies rather than adds goes through a dependency, and
	// a dependency is either on or off: stacks scale the flat amounts only.
	multipliers := make([]*stats.StatDependency, len(config.Stats))
	for i, statConfig := range config.Stats {
		if statConfig.IsMultiplicative {
			multipliers[i] = aura.Unit.NewDynamicMultiplyStat(statConfig.Stat, statConfig.Amount)
		}
	}

	priority := perStack
	if config.MaxStacks > 0 {
		priority = 0
	}

	return aura.NewExclusiveEffect(config.Category, config.SingleAura, ExclusiveEffect{
		Priority: priority,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			applyGeneratedAmounts(ee.Aura.Unit, sim, config, multipliers, generatedStackFactor(ee, perStack))
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			applyGeneratedAmounts(ee.Aura.Unit, sim, config, multipliers, -generatedStackFactor(ee, perStack))
		},
	})
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

	// The shield has no stat to apply or remove: what the category decides is
	// which aura keeps its proc trigger, so the damage is the whole bid.
	if config.Category != "" {
		aura.NewExclusiveEffect(config.Category, config.SingleAura, ExclusiveEffect{Priority: damage})
	}
	joinSharedCategory(aura, config)
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

// The second category the aura joins without an effect of its own. Only the
// player's own copy joins it: the external copy has to be able to sit next to
// the one the player casts.
func joinSharedCategory(aura *Aura, config GeneratedBuff) {
	if config.SharedCategory == "" || !config.IsPlayer {
		return
	}
	aura.NewExclusiveEffect(config.SharedCategory, true, ExclusiveEffect{})
}

// How strong the buff is, as a positive number. A multiplier is judged by how
// far from 1 it is, so a 20% attack speed reduction (0.8) outbids a 10% one.
func generatedMagnitude(config GeneratedBuff) float64 {
	magnitude := func(amount float64, multiplicative bool) float64 {
		if multiplicative {
			amount -= 1
		}
		if amount < 0 {
			return -amount
		}
		return amount
	}

	if len(config.Stats) > 0 {
		return magnitude(config.Stats[0].Amount, config.Stats[0].IsMultiplicative)
	}
	if len(config.Pseudo) > 0 {
		return magnitude(config.Pseudo[0].Amount, config.Pseudo[0].IsMultiplicative)
	}
	return 0
}

// How many stacks the effect is currently priced at.
func generatedStackFactor(effect *ExclusiveEffect, perStack float64) float64 {
	if perStack == 0 {
		return 1
	}
	return effect.Priority / perStack
}

// Applies every stat and pseudo-stat the buff holds, scaled by the stack count
// the exclusive effect is priced at. A negative factor takes them away again.
func applyGeneratedAmounts(unit *Unit, sim *Simulation, config GeneratedBuff, multipliers []*stats.StatDependency, factor float64) {
	for i, statConfig := range config.Stats {
		if multipliers[i] != nil {
			// The build-phase variants, because a buff applied during
			// applyBuildPhaseAuras must not recompute the unit's stats yet; they
			// delegate to the dynamic ones outside it.
			if factor > 0 {
				unit.EnableBuildPhaseStatDep(sim, multipliers[i])
			} else {
				unit.DisableBuildPhaseStatDep(sim, multipliers[i])
			}
			continue
		}
		unit.AddStatDynamic(sim, statConfig.Stat, statConfig.Amount*factor)
	}
	for _, pseudoConfig := range config.Pseudo {
		for _, field := range generatedPseudoStatFields(unit, pseudoConfig) {
			if pseudoConfig.IsMultiplicative {
				applyGeneratedMultiplier(field, pseudoConfig.Amount, factor)
			} else {
				*field += pseudoConfig.Amount * factor
			}
		}
	}
}

// A multiplier is applied once per stack and divided back out the same way.
func applyGeneratedMultiplier(field *float64, amount float64, factor float64) {
	for i := float64(0); i < factor; i++ {
		*field *= amount
	}
	for i := float64(0); i > factor; i-- {
		*field /= amount
	}
}

// The pseudo-stats a generated buff modifies. Under a category each config
// registers one effect that walks every field it names - a damage-taken debuff
// names one per school - so that two buffs modifying the same fields do not both
// apply, while the fields of one buff always move together.
func attachGeneratedPseudoStats(aura *Aura, config GeneratedBuff) {
	category := config.Category
	if config.StatCategory != "" {
		category = config.StatCategory
	}

	for _, pseudoConfig := range config.Pseudo {
		fields := generatedPseudoStatFields(aura.Unit, pseudoConfig)
		if category == "" {
			for _, field := range fields {
				if pseudoConfig.IsMultiplicative {
					aura.AttachMultiplicativePseudoStatBuff(field, pseudoConfig.Amount)
				} else {
					aura.AttachAdditivePseudoStatBuff(field, pseudoConfig.Amount)
				}
			}
			continue
		}

		// One effect for the whole config rather than one per field: a category
		// keeps a single effect per aura, so a per-field effect would leave every
		// school of a damage-taken debuff but the first unapplied.
		suffix := "Add"
		if pseudoConfig.IsMultiplicative {
			suffix = "Mul"
		}
		amount := pseudoConfig.Amount
		multiplicative := pseudoConfig.IsMultiplicative
		aura.NewExclusiveEffect(category+pseudoConfig.Kind.Name()+suffix, false, ExclusiveEffect{
			Priority: generatedMagnitude(GeneratedBuff{Pseudo: []PseudoConfig{pseudoConfig}}),
			OnGain: func(_ *ExclusiveEffect, _ *Simulation) {
				for _, field := range fields {
					if multiplicative {
						*field *= amount
					} else {
						*field += amount
					}
				}
			},
			OnExpire: func(_ *ExclusiveEffect, _ *Simulation) {
				for _, field := range fields {
					if multiplicative {
						*field /= amount
					} else {
						*field -= amount
					}
				}
			},
		})
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
