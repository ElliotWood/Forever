package warlock

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

const ConflagrateRanks = 4

func (warlock *Warlock) getConflagrateConfig(rank int) core.SpellConfig {
	spellId := [ConflagrateRanks + 1]int32{0, 17962, 18930, 18931, 18932}[rank]
	// TODO: The Forever tooltip puts rank 1 at 109 to 132, less than half Classic's 249 to 316,
	// and Incinerate at 125 to 140. Neither spell's higher ranks were shown, so both keep their
	// Classic tables until the beta lists them.
	baseDamageMin := [ConflagrateRanks + 1]float64{0, 249, 319, 395, 447}[rank]
	baseDamageMax := [ConflagrateRanks + 1]float64{0, 316, 400, 491, 557}[rank]
	manaCost := [ConflagrateRanks + 1]float64{0, 165, 200, 230, 255}[rank]
	level := [ConflagrateRanks + 1]int{0, 0, 48, 54, 60}[rank]

	spCoeff := 0.429

	// 20% per point, so at 5/5 Conflagrate stops consuming Immolate altogether. The demo
	// only showed rank 1 and the tree repeated its 20% at every rank, which is why this was
	// held flat; the beta has since confirmed 80% at rank 4 and 100% at rank 5.
	keepImmolateChance := 0.2 * float64(warlock.Talents.ShadowAndFlame)

	return core.SpellConfig{
		SpellCode:     SpellCode_WarlockConflagrate,
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagAPL | WarlockFlagDestruction,
		Rank:          rank,
		RequiredLevel: level,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.getActiveImmolateSpell(target) != nil
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageMin, baseDamageMax)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			immoSpell := warlock.getActiveImmolateSpell(target)
			if immoSpell != nil && !sim.Proc(keepImmolateChance, "Shadow and Flame") {
				immoSpell.Dot(target).Deactivate(sim)
			}
		},
	}
}

func (warlock *Warlock) registerConflagrateSpell() {
	if !warlock.Talents.Conflagrate {
		return
	}

	warlock.Conflagrate = make([]*core.Spell, 0)
	for rank := 1; rank <= ConflagrateRanks; rank++ {
		config := warlock.getConflagrateConfig(rank)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.Conflagrate = append(warlock.Conflagrate, warlock.GetOrRegisterSpell(config))
		}
	}
}
