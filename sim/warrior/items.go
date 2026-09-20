package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

// The abilities spell 21887's class mask lists.
const warriorsWrathSpells = SpellMaskHeroicStrike | SpellMaskRend | SpellMaskShieldBash |
	SpellMaskCleave | SpellMaskDisarm | SpellMaskWhirlwind | SpellMaskSunderArmor | SpellMaskSlam |
	SpellMaskHamstring | SpellMaskExecute | SpellMaskPummel | SpellMaskRevenge | SpellMaskOverpower |
	SpellMaskThunderClap | SpellMaskMockingBlow | SpellMaskMortalStrike | SpellMaskConcussionBlow |
	SpellMaskShieldSlam | SpellMaskRetaliation | SpellMaskIntercept | SpellMaskBloodthirst

var ItemSetBattlegearOfMight = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Might",
	ID:   209,
	Bonuses: map[int32]core.ApplySetBonus{
		3: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 23562 states 30 block value.
			setBonusAura.AttachStatBuff(stats.BlockValue, 30)
		},
		5: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 21838 states a 20% chance to generate an additional rage point whenever
			// damage is dealt to you; spell 29478 energizes 10, which is 1 rage.
			warrior := agent.(WarriorAgent).GetWarrior()
			rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 29478})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:               "Battlegear of Might - 5PC",
				ActionID:           core.ActionID{SpellID: 21838},
				Callback:           core.CallbackOnSpellHitTaken | core.CallbackOnPeriodicDamageTaken,
				Outcome:            core.OutcomeLanded,
				RequireDamageDealt: true,
				ProcChance:         0.2,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					warrior.AddRage(sim, 1, rageMetrics)
				},
			})
		},
		8: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 23561 states 15% more Sunder Armor threat. Sunder Armor deals no damage, so
			// all of its threat is the flat bonus, which the threat multiplier mods do not reach.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask:  SpellMaskSunderArmor,
				Kind:       core.SpellMod_Custom,
				FloatValue: 0.15,
				ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
					spell.FlatThreatBonus *= 1 + mod.GetFloatValue()
				},
				RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {
					spell.FlatThreatBonus /= 1 + mod.GetFloatValue()
				},
			})
		},
	},
})

var ItemSetBattlegearOfWrath = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Wrath",
	ID:   218,
	Bonuses: map[int32]core.ApplySetBonus{
		3: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 23563 states 30 attack power on Battle Shout. shouts.go's HasBsT2 option
			// stands in for this bonus, so nothing is applied here.
		},
		5: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 21890 states a 20% chance after an offensive ability requiring rage that the
			// next one costs less rage; spell 21887 states -50, which is 5 rage, for 10 seconds
			// and one charge.
			warrior := agent.(WarriorAgent).GetWarrior()

			var buff *core.Aura
			buff = warrior.RegisterAura(core.Aura{
				Label:    "Warrior's Wrath",
				ActionID: core.ActionID{SpellID: 21887},
				Duration: time.Second * 10,
			}).AttachSpellMod(core.SpellModConfig{
				ClassMask: warriorsWrathSpells,
				Kind:      core.SpellMod_PowerCost_Flat,
				IntValue:  -5,
			}).AttachProcTrigger(core.ProcTrigger{
				Name:               "Warrior's Wrath - Consume",
				ClassSpellMask:     warriorsWrathSpells,
				Callback:           core.CallbackOnCastComplete,
				TriggerImmediately: true,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					buff.Deactivate(sim)
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Battlegear of Wrath - 5PC",
				ActionID:       core.ActionID{SpellID: 21890},
				ClassSpellMask: warriorsWrathSpells,
				Callback:       core.CallbackOnSpellHitDealt,
				Outcome:        core.OutcomeLanded,
				ProcChance:     0.2,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					buff.Activate(sim)
				},
			})
		},
		8: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 23548 states a 4% chance to parry the next attack after a block; spell 23547
			// states 100% parry, lasts until an attack is taken and holds one charge.
			warrior := agent.(WarriorAgent).GetWarrior()

			var parry *core.Aura
			parry = warrior.RegisterAura(core.Aura{
				Label:    "Battlegear of Wrath Parry",
				ActionID: core.ActionID{SpellID: 23547},
				Duration: core.NeverExpires,
			}).AttachStatBuff(stats.ParryRating, 100*core.ParryRatingPerParryPercent).
				AttachProcTrigger(core.ProcTrigger{
					Name:               "Battlegear of Wrath - 8PC Consume",
					ProcMask:           core.ProcMaskMelee,
					Callback:           core.CallbackOnSpellHitTaken,
					TriggerImmediately: true,
					Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
						parry.Deactivate(sim)
					},
				})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:       "Battlegear of Wrath - 8PC",
				ActionID:   core.ActionID{SpellID: 23548},
				Callback:   core.CallbackOnSpellHitTaken,
				Outcome:    core.OutcomeBlock,
				ProcChance: 0.04,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					parry.Activate(sim)
				},
			})
		},
	},
})

var ItemSetConquerorsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Conqueror's Battlegear",
	ID:   496,
	Bonuses: map[int32]core.ApplySetBonus{
		3: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 26109 states -35% rage cost on all warrior shouts, and its class mask lists
			// Battle Shout, Demoralizing Shout, Intimidating Shout and Challenging Shout.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask:  SpellMaskBattleShout | SpellMaskDemoralizingShout | SpellMaskIntimidatingShout | SpellMaskChallengingShout,
				Kind:       core.SpellMod_PowerCost_Pct_Add,
				FloatValue: -0.35,
			})
		},
		5: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 26110 states 50% on Thunder Clap's slow effect and damage.
			// TODO: only the damage half is modelled. Thunder Clap's slow is core.ThunderClapAura,
			// which takes the Improved Thunder Clap rank and no set bonus.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask:  SpellMaskThunderClap,
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.5,
			})
		},
	},
})

var ItemSetDreadnaughtsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Dreadnaught's Battlegear",
	ID:   523,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 28844 states 75 Revenge damage.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask:  SpellMaskRevenge,
				Kind:       core.SpellMod_BaseDamage_Flat,
				FloatValue: 75,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 28843 states 5% chance to hit with Taunt and Challenging Shout.
			// TODO: both spells resolve on OutcomeAlwaysHit, so the bonus changes nothing until
			// they roll against the spell hit table.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask:  SpellMaskTaunt | SpellMaskChallengingShout,
				Kind:       core.SpellMod_BonusHit_Percent,
				FloatValue: 5,
			})
		},
		6: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 28842 states 5% chance to hit with Sunder Armor, Heroic Strike, Revenge and
			// Shield Slam.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask:  SpellMaskSunderArmor | SpellMaskHeroicStrike | SpellMaskRevenge | SpellMaskShieldSlam,
				Kind:       core.SpellMod_BonusHit_Percent,
				FloatValue: 5,
			})
		},
		8: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 28845 states that below 20% health, healing spells cast on you gain up to 160
			// healing (spell 28846) for 5 seconds.
			// TODO: not modelled. The sim has no hook on healing received, and the warrior takes
			// no incoming heals to apply it to.
		},
	},
})
