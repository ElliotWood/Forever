package mage

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

const ArcaneBlastMaxStacks = 4

func (mage *Mage) registerArcaneBlastSpell() {
	if !mage.Talents.ArcaneBlast {
		return
	}

	// Level 60 values aren't datamined yet. The demo tooltip reads 95 to 104, which matches a rank
	// 1 of the spell, so these are picked to give the same damage per cast second as the level 60
	// Frostbolt the mage already casts instead of being scaled straight off the tooltip.
	baseDamage := []float64{435, 505}
	spellCoeff := .714
	manaCost := 240.0
	castTime := time.Millisecond * 2500

	actionID := core.ActionID{SpellID: 30451}

	// Arcane Blast buffs the mage's other spells rather than itself, and the next one of them
	// spends the stacks.
	var affectedSpells []*core.Spell
	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagMage) && spell.SpellCode != SpellCode_MageArcaneBlast {
			affectedSpells = append(affectedSpells, spell)
		}
	})

	mage.ArcaneBlastAura = mage.RegisterAura(core.Aura{
		Label:     "Arcane Blast",
		ActionID:  actionID,
		Duration:  time.Second * 8,
		MaxStacks: ArcaneBlastMaxStacks,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			delta := newStacks - oldStacks
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive += .10 * float64(delta)
			}
			mage.ArcaneBlast.Cost.Multiplier += 175 * delta
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Flags.Matches(SpellFlagMage) || !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}

			// OnCastComplete runs after the damage is rolled, so the spell that drops the stacks is
			// still buffed by them. Arcane Missiles spends them when the channel ends instead, so
			// that every missile benefits.
			switch spell.SpellCode {
			case SpellCode_MageArcaneBlast, SpellCode_MageArcaneMissiles, SpellCode_MageArcaneMissilesTick:
				return
			}

			aura.Deactivate(sim)
		},
	})

	mage.ArcaneBlast = mage.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_MageArcaneBlast,
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		RequiredLevel: 60,
		Rank:          1,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: castTime,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := sim.Roll(baseDamage[0], baseDamage[1])
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

			mage.ArcaneBlastAura.Activate(sim)
			mage.ArcaneBlastAura.AddStack(sim)
		},
	})
}

// Arcane Missiles holds on to the Arcane Blast stacks for the whole channel and spends them when
// the last missile lands.
func (mage *Mage) spendArcaneBlastStacks(sim *core.Simulation) {
	if mage.ArcaneBlastAura != nil {
		mage.ArcaneBlastAura.Deactivate(sim)
	}
}
