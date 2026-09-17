package druid

import (
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/stats"
)

var ItemSetFeralheartRaiment = core.NewItemSet(core.ItemSet{
	Name: "Feralheart Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// (2) Set : +8 All Resistances.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// (3) Set : Restores 8 mana per 5 sec.
		3: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.MP5, 8)
		},
		// (4) Set : Nature's Bounty, cut from 300 mana / 40 energy / 10 rage to 200 / 20 / 10. The
		// client stores no proc chance, so Classic's 2% is kept.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450608}
			manaMetrics := c.NewManaMetrics(actionID)
			energyMetrics := c.NewEnergyMetrics(actionID)
			rageMetrics := c.NewRageMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Nature's Bounty (Mana)",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.02,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					c.AddMana(sim, 200, manaMetrics)
				},
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Nature's Bounty (Energy)",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMeleeWhiteHit,
				ProcChance: 0.02,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					if c.HasEnergyBar() {
						c.AddEnergy(sim, 20, energyMetrics)
					}
				},
			})
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Nature's Bounty (Rage)",
				Callback:   core.CallbackOnSpellHitTaken,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 0.02,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					if c.HasRageBar() {
						c.AddRage(sim, 10, rageMetrics)
					}
				},
			})
		},
		// (5) Set : Wild Heart, a movement speed burst when struck.
		5: func(_ core.Agent) {},
		// (6) Set : +23 damage and healing done by magical spells and effects, and +40 Attack Power.
		// Classic gave 15 and 26, at eight and six pieces respectively.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.SpellDamage:       23,
				stats.HealingPower:      23,
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
	},
})

var ItemSetCenarionRaiment = core.NewItemSet(core.ItemSet{
	Name: "Cenarion Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// (3) Set : Damage dealt by Thorns increased by 4 and duration increased by 50%.
		3: func(agent core.Agent) {
			// Nothing to do
		},
		// (5) Set : Improves your chance to get a critical strike with spells by 2%.
		5: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellCrit, 2*core.SpellCritRatingPerCritChance)
		},
		// (8) Set : Reduces the cooldown of your Tranquility and Hurricane spells by 50%.
		8: func(agent core.Agent) {
			// Nothing to do
		},
	},
})

var ItemSetStormrageRaiment = core.NewItemSet(core.ItemSet{
	Name: "Stormrage Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// (3) Set : Allows 15% of your Mana regeneration to continue while casting.
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.PseudoStats.SpiritRegenRateCasting += .15
		},
		// (5) Set : Reduces the casting time of your Regrowth spell by 0.2 sec.
		5: func(agent core.Agent) {
			// Nothing to do.
		},
		// (8) Set : Increases the duration of your Rejuvenation spell by 3 sec.
		8: func(agent core.Agent) {
			// Nothing to do.
		},
	},
})

var ItemSetSymbolsOfUnendingLife = core.NewItemSet(core.ItemSet{
	Name: "Symbols of Unending Life",
	Bonuses: map[int32]core.ApplyEffect{
		// (3) Set : Your finishing moves now refund 30 energy on a Miss, Dodge, Block, or Parry.
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 26107}
			energyMetrics := c.NewEnergyMetrics(actionID)
			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				Name:     "Symbols of Unending Life Finisher Bonus",
				Callback: core.CallbackOnSpellHitDealt,
				Outcome:  core.OutcomeMiss | core.OutcomeDodge | core.OutcomeBlock | core.OutcomeParry,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if spell.SpellCode == SpellCode_DruidFerociousBite || spell.SpellCode == SpellCode_DruidRip {
						c.AddEnergy(sim, 30, energyMetrics)
					}
				},
			})
		},
	},
})
