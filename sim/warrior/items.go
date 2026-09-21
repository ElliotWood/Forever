package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

var ItemSetBattlegearOfMight = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Might",
	ID:   209,
	Bonuses: map[int32]core.ApplySetBonus{
		3: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.BlockValue, 30)
		},
		5: func(agent core.Agent, setBonusAura *core.Aura) {
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
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask:  SpellMaskSunderArmor,
				Kind:       core.SpellMod_FlatThreatBonus_Pct,
				FloatValue: 0.15,
			})
		},
	},
})

var ItemSetBattlegearOfWrath = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Wrath",
	ID:   218,
	Bonuses: map[int32]core.ApplySetBonus{
		3: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 23563 states 30 attack power on Battle Shout, which shouts.go adds through
			// the same flag the HasBsT2 option sets. The set aura toggles it so an item swap
			// that removes the pieces takes the bonus with them.
			warrior := agent.(WarriorAgent).GetWarrior()
			fromOptions := warrior.HasBsT2
			setBonusAura.
				ApplyOnGain(func(_ *core.Aura, _ *core.Simulation) {
					warrior.HasBsT2 = true
				}).
				ApplyOnExpire(func(_ *core.Aura, _ *core.Simulation) {
					warrior.HasBsT2 = fromOptions
				})
		},
		5: func(agent core.Agent, setBonusAura *core.Aura) {
			warrior := agent.(WarriorAgent).GetWarrior()

			var buff *core.Aura
			buff = warrior.RegisterAura(core.Aura{
				Label:    "Warrior's Wrath",
				ActionID: core.ActionID{SpellID: 21887},
				Duration: time.Second * 10,
			}).
				AttachSpellMod(core.SpellModConfig{
					ClassMask: SpellMaskOffensiveAbilities,
					Kind:      core.SpellMod_PowerCost_Flat,
					IntValue:  -5,
				}).
				AttachProcTrigger(core.ProcTrigger{
					Name:               "Warrior's Wrath - Consume",
					ClassSpellMask:     SpellMaskOffensiveAbilities,
					Callback:           core.CallbackOnCastComplete,
					TriggerImmediately: true,
					Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
						buff.Deactivate(sim)
					},
				})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Battlegear of Wrath - 5PC",
				ActionID:       core.ActionID{SpellID: 21890},
				ClassSpellMask: SpellMaskOffensiveAbilities,
				Callback:       core.CallbackOnSpellHitDealt,
				Outcome:        core.OutcomeLanded,
				ProcChance:     0.2,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					buff.Activate(sim)
				},
			})
		},
		8: func(agent core.Agent, setBonusAura *core.Aura) {
			warrior := agent.(WarriorAgent).GetWarrior()

			var parry *core.Aura
			parry = warrior.RegisterAura(core.Aura{
				Label:    "Battlegear of Wrath Parry",
				ActionID: core.ActionID{SpellID: 23547},
				Duration: core.NeverExpires,
			}).
				AttachStatBuff(stats.ParryRating, 100*core.ParryRatingPerParryPercent).
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
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask:  SpellMaskShouts,
				Kind:       core.SpellMod_PowerCost_Pct_Add,
				FloatValue: -0.35,
			})
		},
		5: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 26110 states 50% on all of Thunder Clap's effects: the damage here, the slow
			// through the bonus thunder_clap.go reads when its aura lands.
			warrior := agent.(WarriorAgent).GetWarrior()
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask:  SpellMaskThunderClap,
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.5,
			})
			setBonusAura.ApplyOnGain(func(_ *core.Aura, _ *core.Simulation) {
				warrior.thunderClapEffectBonus += 0.5
			})
			setBonusAura.ApplyOnExpire(func(_ *core.Aura, _ *core.Simulation) {
				warrior.thunderClapEffectBonus -= 0.5
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
			// healing (spell 28846) for 5 seconds. The modelled incoming healing of the tank sim
			// bypasses the healing bonus; heals a healer unit casts on the warrior take it.
			warrior := agent.(WarriorAgent).GetWarrior()
			const cheatDeathHealing = 160
			cheatDeath := warrior.RegisterAura(core.Aura{
				Label:    "Cheat Death",
				ActionID: core.ActionID{SpellID: 28846},
				Duration: time.Second * 5,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warrior.PseudoStats.BonusHealingTaken += cheatDeathHealing
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warrior.PseudoStats.BonusHealingTaken -= cheatDeathHealing
				},
			})
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:     "Cheat Death - Trigger",
				Callback: core.CallbackOnSpellHitTaken | core.CallbackOnPeriodicDamageTaken,
				Outcome:  core.OutcomeLanded,
				ExtraCondition: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool {
					return warrior.CurrentHealthPercent() < 0.2
				},
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					cheatDeath.Activate(sim)
				},
			})
		},
	},
})
