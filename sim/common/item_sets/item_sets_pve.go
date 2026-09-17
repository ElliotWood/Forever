package item_sets

import (
	"time"

	"github.com/wowsims/classic/sim/common/guardians"
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                                 Cloth
///////////////////////////////////////////////////////////////////////////

// The four Scholomance and Stratholme sets below were rebuilt by Forever the same way the dungeon
// sets were, read from ItemSetSpell in beta client 1.60.1.69893. Their 3 piece is now a proc whose
// chance the client does not store (SpellAuraOptions holds 100, which means unset), so those are
// described and left unapplied rather than given an invented rate. Everything else is a plain value.

var ItemSetNecropileRaiment = core.NewItemSet(core.ItemSet{
	Name: "Necropile Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		// Improves your chance to hit by 0.5% (1299734, 5 rating).
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeHit, 0.5*core.MeleeHitRatingPerHitChance)
			character.AddStat(stats.SpellHit, 0.5*core.SpellHitRatingPerHitChance)
		},
		// Necrophile Drain (1299737): harmful spell casts have a chance to steal 28 life. The client
		// stores no proc chance for it, so it is left out rather than guessed.
		3: func(_ core.Agent) {},
		// +5 All Resistances (18676, was 15).
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddResistances(5)
		},
		// Increases damage and healing done by magical spells and effects by up to 23.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 23)
		},
	},
})

var ItemSetIronweaveBattlesuit = core.NewItemSet(core.ItemSet{
	Name: "Ironweave Battlesuit",
	Bonuses: map[int32]core.ApplyEffect{
		// +200 Armor, moved down from eight pieces.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Armor, 200)
		},
		// Decreases the magical resistances of your spell targets by 5.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPenetration, 5)
		},
		// Increases your chance to resist Silence and Interrupt effects by 10%.
		4: func(_ core.Agent) {},
		// Increases damage and healing done by magical spells and effects by up to 23.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 23)
		},
		// Reduces damage taken while Stunned by 15%; a raid boss encounter never stuns.
		6: func(_ core.Agent) {},
	},
})

var ItemSetThePostmaster = core.NewItemSet(core.ItemSet{
	Name: "The Postmaster",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases run speed by 8%.
		2: func(_ core.Agent) {},
		// Increases damage and healing done by magical spells and effects by up to 23, where
		// Classic gave 12 at four pieces.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 23)
		},
		// Return to Sender: reflects the next spell cast on you after dropping below 25% health.
		4: func(_ core.Agent) {},
		// Improves your chance to hit by 1% (432639).
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeHit, 1*core.MeleeHitRatingPerHitChance)
			character.AddStat(stats.SpellHit, 1*core.SpellHitRatingPerHitChance)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Leather
///////////////////////////////////////////////////////////////////////////

var ItemSetCadaverousGarb = core.NewItemSet(core.ItemSet{
	Name: "Cadaverous Garb",
	Bonuses: map[int32]core.ApplyEffect{
		// +10 Attack Power, and movement impairing effects 10% shorter. Forever has no 2 piece.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.AttackPower, 10)
			character.AddStat(stats.RangedAttackPower, 10)
		},
		// +5 All Resistances (18676, was 15).
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddResistances(5)
		},
		// Improves your chance to hit by 2%, now spells as well as melee (460230). The sim was
		// adding 2 hit rating rather than 2%.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeHit, 2*core.MeleeHitRatingPerHitChance)
			character.AddStat(stats.SpellHit, 2*core.SpellHitRatingPerHitChance)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Mail
///////////////////////////////////////////////////////////////////////////

var ItemSetBloodmailRegalia = core.NewItemSet(core.ItemSet{
	Name: "Bloodmail Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		// +10 Attack Power.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.AttackPower, 10)
			character.AddStat(stats.RangedAttackPower, 10)
		},
		// Bloodmail (1299740): melee attacks have a chance to bleed the target for 7 a second over
		// 10 sec. The client stores no proc chance, so it is left out rather than guessed.
		3: func(_ core.Agent) {},
		// +5 All Resistances (18676, was 15).
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddResistances(5)
		},
		// Improves your chance to crit by 1.5% (1299738, 21 rating), where Classic gave 1% parry.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1.5*core.CritRatingPerCritChance)
			character.AddStat(stats.SpellCrit, 1.5*core.SpellCritRatingPerCritChance)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Plate
///////////////////////////////////////////////////////////////////////////

var ItemSetDeathboneGuardian = core.NewItemSet(core.ItemSet{
	Name: "Deathbone Guardian",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +3.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Defense, 3)
		},
		// Deathbone Surge: 15% chance when struck to gain 70 spell damage and healing for 10 sec,
		// once a minute.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 1299743}
			procAura := character.NewTemporaryStatsAura("Deathbone Surge", core.ActionID{SpellID: 1299746}, stats.Stats{stats.SpellPower: 70}, time.Second*10)
			icd := core.Cooldown{Timer: character.NewTimer(), Duration: time.Minute}

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Deathbone Surge",
				Callback:   core.CallbackOnSpellHitTaken,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 0.15,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					if icd.IsReady(sim) {
						icd.Use(sim)
						procAura.Activate(sim)
					}
				},
			})
		},
		// +5 All Resistances (18676, was 15).
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddResistances(5)
		},
		// Reduces the chance for your attacks to be dodged or parried by 2% (1213289), where
		// Classic gave 1% parry.
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Expertise, 2)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Other
///////////////////////////////////////////////////////////////////////////

var ItemSetSpidersKiss = core.NewItemSet(core.ItemSet{
	Name: "Spider's Kiss",
	Bonuses: map[int32]core.ApplyEffect{
		// Chance on Hit: Immobilizes the target and lowers their armor by 100 for 10 sec.
		// Unsure about exlusive effects with this aura also looks like it might be lowering the characters armor instead of the enemy?
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			procAura := character.NewTemporaryStatsAura("Spider's Kiss", core.ActionID{SpellID: 17333}, stats.Stats{stats.Armor: -100}, time.Second*10)
			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{SpellID: 17333},
				Name:       "Spider's Kiss",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 0.05,
				Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					procAura.Activate(sim)
				},
			})
		},
	},
})

var ItemSetDalRendsArms = core.NewItemSet(core.ItemSet{
	Name: "Dal'Rend's Arms",
	Bonuses: map[int32]core.ApplyEffect{
		// +50 Attack Power.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.AttackPower, 50)
			character.AddStat(stats.RangedAttackPower, 50)
		},
	},
})

var ItemSetShardOfTheGods = core.NewItemSet(core.ItemSet{
	Name: "Shard of the Gods",
	Bonuses: map[int32]core.ApplyEffect{
		// +10 All Resistances.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddResistances(15)
		},
	},
})

var ItemSetSpiritOfEskhandar = core.NewItemSet(core.ItemSet{
	Name: "Spirit of Eskhandar",
	Bonuses: map[int32]core.ApplyEffect{
		// +30 Attack Power against Humanoids (1298478), new in Forever.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			matchingTargets := core.FilterSlice(
				character.Env.Encounter.TargetUnits,
				func(unit *core.Unit) bool { return unit.MobType == proto.MobType_MobTypeHumanoid },
			)

			core.MakePermanent(character.GetOrRegisterAura(core.Aura{
				Label: "Spirit of Eskhandar - Humanoid Slaying",
				OnGain: func(_ *core.Aura, _ *core.Simulation) {
					for _, target := range matchingTargets {
						for _, at := range character.AttackTables[target.UnitIndex] {
							at.BonusAttackPowerTaken += 30
						}
					}
				},
				OnExpire: func(_ *core.Aura, _ *core.Simulation) {
					for _, target := range matchingTargets {
						for _, at := range character.AttackTables[target.UnitIndex] {
							at.BonusAttackPowerTaken -= 30
						}
					}
				},
			}))
		},
		// Improves your chance to crit by 1% (1314828), new in Forever. The client doubles it at
		// night, which the sim has no clock for.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance)
			character.AddStat(stats.SpellCrit, 1*core.SpellCritRatingPerCritChance)
		},
		// 1% chance on a melee hit to call forth the spirit of Eskhandar to protect you in battle for 2 min.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Call of Eskhandar Trigger",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 1,
				ICD:        time.Minute * 1,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					for _, petAgent := range character.PetAgents {
						if eskhandar, ok := petAgent.(*guardians.Eskhandar); ok {
							eskhandar.EnableWithTimeout(sim, eskhandar, time.Minute*2)
							break
						}
					}
				},
			})
		},
	},
})

var ItemSetRegaliaOfUndeadCleansing = core.NewItemSet(core.ItemSet{
	Name: "Regalia of Undead Cleansing",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases your damage against undead by 2%.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			delta := 1.02

			character.Env.RegisterPostFinalizeEffect(func() {
				for _, target := range character.Env.Encounter.TargetUnits {
					if target.MobType != proto.MobType_MobTypeUndead {
						continue
					}

					for _, at := range character.AttackTables[target.UnitIndex] {
						at.DamageDealtMultiplier *= delta
						at.CritMultiplier *= delta
					}
				}
			})
		},
	},
})

var ItemSetUndeadSlayersArmor = core.NewItemSet(core.ItemSet{
	Name: "Undead Slayer's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases your damage against undead by 2%.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			delta := 1.02

			character.Env.RegisterPostFinalizeEffect(func() {
				for _, target := range character.Env.Encounter.TargetUnits {
					if target.MobType != proto.MobType_MobTypeUndead {
						continue
					}

					for _, at := range character.AttackTables[target.UnitIndex] {
						at.DamageDealtMultiplier *= delta
						at.CritMultiplier *= delta
					}
				}
			})
		},
	},
})

var ItemSetGarbOfTheUndeadSlayer = core.NewItemSet(core.ItemSet{
	Name: "Garb of the Undead Slayer",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases your damage against undead by 2%.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			delta := 1.02

			character.Env.RegisterPostFinalizeEffect(func() {
				for _, target := range character.Env.Encounter.TargetUnits {
					if target.MobType != proto.MobType_MobTypeUndead {
						continue
					}

					for _, at := range character.AttackTables[target.UnitIndex] {
						at.DamageDealtMultiplier *= delta
						at.CritMultiplier *= delta
					}
				}
			})
		},
	},
})

var ItemSetBattlegearOfUndeadSlaying = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Undead Slaying",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases your damage against undead by 2%.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			delta := 1.02

			character.Env.RegisterPostFinalizeEffect(func() {
				for _, target := range character.Env.Encounter.TargetUnits {
					if target.MobType != proto.MobType_MobTypeUndead {
						continue
					}

					for _, at := range character.AttackTables[target.UnitIndex] {
						at.DamageDealtMultiplier *= delta
						at.CritMultiplier *= delta
					}
				}
			})
		},
	},
})
