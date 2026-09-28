package hunter

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

type PetAbilityType int

// Pet AI doesn't use abilities immediately, so model this with a 1.6s GCD.
const PetGCD = time.Millisecond * 1600

const (
	Unknown PetAbilityType = iota

	Bite
	Claw
	LightningBreath
	Screech
	ScorpidPoison
	SavageRend
	Pinch
	Dismember
)

func (hp *HunterPet) NewPetAbility(abilityType PetAbilityType) *core.Spell {
	switch abilityType {
	case Bite:
		return hp.newBite()
	case Claw:
		return hp.newClaw()
	case LightningBreath:
		return hp.newLightningBreath()
	case Screech:
		return hp.newScreech()
	case ScorpidPoison:
		return hp.newScorpidPoison()
	case SavageRend:
		return hp.newSavageRend()
	case Pinch:
		return hp.newPetStrike(spellData.PinchTriggered.Highest())
	case Dismember:
		return hp.newPetStrike(spellData.DismemberTriggered.Highest())
	case Unknown:
		return nil
	default:
		panic("Invalid pet ability type")
	}
}

func (hp *HunterPet) newBite() *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 17261},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: 35,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: 10 * time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, sim.Roll(81, 99), spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

func (hp *HunterPet) newClaw() *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 3009},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, sim.Roll(43, 59), spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

// Beta client: every rank lower than Classic's, and no more growth per level. Rank 6 is 86-98.
func (hp *HunterPet) newLightningBreath() *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 25012},
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMagic,
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskSpellDamage,
		MaxRange:       20,

		FocusCost: core.FocusCostOptions{
			Cost: 50,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, sim.Roll(86, 98), spell.OutcomeMagicHitAndCrit)
		},
	})
}

// Demoralizing Screech in the beta client: new damage, and a 10 sec cooldown Classic did not have.
// The attack power reduction it also applies is left out.
func (hp *HunterPet) newScreech() *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 24582},
		SpellSchool:    core.SpellSchoolPhysical,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: 20,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, sim.Roll(24, 42), spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

// Beta client: 5 a tick at rank 4, down from 8.
func (hp *HunterPet) newScorpidPoison() *core.Spell {
	const baseDamageTick = 5.0

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 24587},
		SpellSchool:    core.SpellSchoolNature,
		DefenseType:    core.DefenseTypeMelee,
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagPassiveSpell | core.SpellFlagPoison,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: 30,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: time.Second * 4,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Scorpid Poison",
				MaxStacks: 5,
				Duration:  time.Second * 10,
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Only the first stack snapshots the multiplier.
				if dot.GetStacks() <= 1 {
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.Index], true)
					dot.SnapshotBaseDamage = 0
				}
				dot.SnapshotBaseDamage += baseDamageTick
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if !result.Landed() {
				return
			}

			dot := spell.Dot(target)
			dot.Apply(sim)
			if dot.GetStacks() < dot.MaxStacks {
				dot.AddStack(sim)
				dot.TakeSnapshot(sim)
			}
		},
	})
}

// Savage Rend is new in Forever and the Raptor's alone (client SkillLineAbility puts it on skill line
// 217, Raptor): a bleed for 50 focus on a 1 min cooldown, everything read off the client row. Beta
// logs (Tynman's and Consumer's raptors, foreverlogs 2650/2669/2673/2674) put a tick at 6.4-7.0 at
// rank 1, which is the 5 base times the pet's happiness and Raptor damage scalars, so no attack power
// share. The 5% more bleed damage the tooltip adds is left out: nothing else of ours bleeds.
func (hp *HunterPet) newSavageRend() *core.Spell {
	rank := spellData.SavageRendTriggered.Highest()
	tick := rank.PeriodicEffect()
	tickLength := tick.Period()

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: rank.Cooldown(),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Savage Rend",
			},
			NumberOfTicks: int32(rank.Duration() / tickLength),
			TickLength:    tickLength,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Average(core.CharacterLevel))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, rank.TickOutcomeHitRolled(dot))
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}

// Pinch (Crab, skill line 214) and Dismember (Crocolisk, 212) are new in Forever: a single melee hit
// with the damage, focus cost and cooldown read off the client row (rank 5: Pinch 95 for 50 focus on
// 30 sec, Dismember 54 for 35 focus on 6 sec, both +-7%). Pinch's snare and Dismember's healing
// reduction are left out. Beta logs: Consumer's crab (foreverlogs 2674) landed rank 1 Pinch (20)
// 17 times at ~17 a hit through level 20 mob armor, so no attack power share.
func (hp *HunterPet) newPetStrike(rank *spelldata.Spell) *core.Spell {
	damage := rank.DamageEffect()

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.ID},
		SpellSchool:    rank.SpellSchool(),
		DefenseType:    rank.DefenseTypeCore(),
		ClassSpellMask: HunterPetDamage,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		MaxRange:       core.MaxMeleeRange,

		FocusCost: core.FocusCostOptions{
			Cost: int32(rank.Cost()),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: rank.Cooldown(),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, damage.Roll(sim, core.CharacterLevel), spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}
