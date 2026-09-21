package shaman

import (
	"github.com/wowsims/forever/sim/core"
)

// The totem lays down a pulse spell of its own; the damage and coefficient live on that spell, not on
// the totem, and the pulse's cast time is the interval between pulses.
var searingTotemRank = spellData.SearingTotem.HighestRank()
var searingTotemAttack = spellData.SearingTotemTriggered.HighestRank()
var magmaTotemRank = spellData.MagmaTotem.HighestRank()
var magmaTotemPulse = spellData.MagmaTotemTriggered.BySpellID(10581)

// Forever has no Fire Nova Totem: the totem's ids are gone and the Fire Nova the spellbook teaches in
// its place is the caster-centred nova, on a 10 sec cooldown.
var fireNovaRank = spellData.FireNova.HighestRank()
var fireNovaDamage = spellData.FireNovaTriggered.HighestRank()

func (shaman *Shaman) registerSearingTotemSpell() {
	attack := shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: searingTotemAttack.SpellID},
		SpellSchool:      searingTotemAttack.SpellSchool,
		DefenseType:      searingTotemAttack.DefenseType,
		ProcMask:         core.ProcMaskEmpty,
		Flags:            SpellFlagShamanSpell | core.SpellFlagPassiveSpell,
		ClassSpellMask:   SpellMaskSearingTotem,
		MissileSpeed:     searingTotemAttack.MissileSpeed,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: searingTotemAttack.Direct.BonusCoefficient(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, searingTotemAttack.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	// The pulse's own cast time is the interval between pulses.
	tickLength := searingTotemAttack.CastTime
	duration := searingTotemRank.Duration

	shaman.SearingTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: searingTotemRank.SpellID},
		SpellSchool:    searingTotemRank.SpellSchool,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL | SpellFlagShamanSpell | SpellFlagInstant,
		ClassSpellMask: SpellMaskSearingTotem,
		ManaCost: core.ManaCostOptions{
			FlatCost: searingTotemRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: searingTotemRank.GCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Searing Totem",
			},
			NumberOfTicks: int32(duration / tickLength),
			TickLength:    tickLength,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				attack.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.cancelFireTotems(sim)
			spell.Dot(sim.Encounter.ActiveTargetUnits[0]).Apply(sim)
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + duration
		},
	})
}

func (shaman *Shaman) registerMagmaTotemSpell() {
	duration := magmaTotemRank.Duration
	tickLength := core.DurationFromSeconds(2)

	shaman.MagmaTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: magmaTotemRank.SpellID},
		SpellSchool:    magmaTotemRank.SpellSchool,
		DefenseType:    core.DefenseTypeMagic,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL | SpellFlagShamanSpell | SpellFlagInstant,
		ClassSpellMask: SpellMaskMagmaTotem,
		ManaCost: core.ManaCostOptions{
			FlatCost: magmaTotemRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: magmaTotemRank.GCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Magma Totem",
			},
			NumberOfTicks:    int32(duration / tickLength),
			TickLength:       tickLength,
			BonusCoefficient: magmaTotemPulse.Direct.BonusCoefficient(),

			OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
				dot.Spell.CalcPeriodicAoeDamage(sim, magmaTotemPulse.Direct.Damage(sim), dot.Spell.OutcomeTickMagicHitAndCrit)
				dot.Spell.DealBatchedPeriodicDamage(sim)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.cancelFireTotems(sim)
			spell.AOEDot().Apply(sim)
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + duration
		},
	})
}

func (shaman *Shaman) registerFireNovaSpell() {
	shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: fireNovaRank.SpellID},
		SpellSchool:    fireNovaRank.SpellSchool,
		DefenseType:    fireNovaRank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | SpellFlagShamanSpell | SpellFlagInstant,
		ClassSpellMask: SpellMaskFireNova,
		ManaCost: core.ManaCostOptions{
			FlatCost: fireNovaRank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: fireNovaRank.GCD,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: fireNovaRank.Cooldown,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: fireNovaDamage.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.CalcAoeDamage(sim, fireNovaDamage.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
			spell.DealBatchedAoeDamage(sim)
		},
	})
}

func (shaman *Shaman) cancelFireTotems(sim *core.Simulation) {
	shaman.MagmaTotem.AOEDot().Deactivate(sim)
	if searingTotemDot := shaman.SearingTotem.Dot(shaman.CurrentTarget); searingTotemDot != nil {
		searingTotemDot.Deactivate(sim)
	}
	if shaman.TotemOfWrath != nil {
		shaman.TotemOfWrath.RelatedSelfBuff.Deactivate(sim)
	}
}
