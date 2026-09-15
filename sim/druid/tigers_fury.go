package druid

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (druid *Druid) registerTigersFurySpell() {
	actionID := core.ActionID{SpellID: map[int32]int32{
		25: 5217,
		40: 6793,
		50: 9845,
		60: 9846,
	}[druid.Level]}

	dmgBonus := map[int32]float64{
		25: 10.0,
		40: 20.0,
		50: 30.0,
		60: 40.0,
	}[druid.Level]

	druid.TigersFuryAura = druid.RegisterAura(core.Aura{
		Label:    "Tiger's Fury Aura",
		ActionID: actionID,
		Duration: 6 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.BonusPhysicalDamage += dmgBonus
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.BonusPhysicalDamage -= dmgBonus
		},
	})

	// Tagged so it does not collide with the metrics the 30 energy cost registers under
	// the same action. Two resource metrics sharing an id and type are indistinguishable
	// in the resources tab, and the concurrency combiner folds them into one.
	energyMetrics := druid.NewEnergyMetrics(actionID.WithTag(1))

	// Forever's King of the Jungle reads "Tiger's Fury now instantly grants you Energy", the
	// Wrath talent word for word, and Wrath's Tiger's Fury had no Energy cost and a 30 second
	// cooldown. Classic's (30 Energy, no cooldown) would let the talent mint Energy, so the
	// Wrath shape is assumed. The +40 damage is Classic's rank 4.
	// TODO: Beta will confirm the cooldown, the cost and the damage bonus.
	forever := druid.Env.IsForever()
	energyCost := core.TernaryFloat64(forever, 0, 30)
	cooldown := core.TernaryDuration(forever, 30*time.Second, time.Second)

	spell := druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost: energyCost,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: cooldown,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if druid.Talents.KingOfTheJungle > 0 {
				druid.AddEnergy(sim, 20*float64(druid.Talents.KingOfTheJungle), energyMetrics)
			}

			druid.TigersFuryAura.Activate(sim)
		},
	})

	druid.TigersFury = spell
}
