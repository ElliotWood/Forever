package priest

import (
	"fmt"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// The Smite build's opener and the trigger Power in Light and Searing Light both read.
// The client's rank 5 dot (13 a tick) is larger than rank 6's (10 a tick); the table is taken as it is.
var HolyFireRankMap = spellData.HolyFire

func (priest *Priest) registerHolyFireSpell(rank shared.SpellData) {
	tick := rank.Periodic.(shared.SpellDataPeriodic)

	priest.HolyFire = append(priest.HolyFire, priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: rank.SpellID},
		SpellSchool:    rank.SpellSchool,
		DefenseType:    rank.DefenseType,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellHolyFire,
		Rank:           rank.Rank,
		MaxRange:       rank.MaxRange,

		ManaCost: core.ManaCostOptions{
			FlatCost: rank.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      rank.GCD,
				CastTime: rank.CastTime,
			},
		},

		DamageMultiplier: 1,
		BonusCoefficient: rank.Direct.BonusCoefficient(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("HolyFire-%d", rank.Rank),
			},
			NumberOfTicks:       tick.NumberOfTicks,
			TickLength:          tick.TickLength,
			AffectedByCastSpeed: false,
			BonusCoefficient:    tick.Coef,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Snapshot(target, tick.Tick)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, priestTickOutcome(rank, dot))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, rank.Direct.Damage(sim), spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealDamage(sim, result)
		},
	}))
}

// Whether this priest's Holy Fire is burning the target, which is what Power in Light asks.
func (priest *Priest) hasActiveHolyFire(target *core.Unit) bool {
	for _, spell := range priest.HolyFire {
		if spell.Dot(target).IsActive() {
			return true
		}
	}
	return false
}
