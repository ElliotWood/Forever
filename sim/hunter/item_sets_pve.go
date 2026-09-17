package hunter

import (
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/stats"
	"time"
)

///////////////////////////////////////////////////////////////////////////
//                            Phase 1 Item Sets - Molten Core
///////////////////////////////////////////////////////////////////////////

var ItemSetGiantStalkers = core.NewItemSet(core.ItemSet{
	Name: "Giantstalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// (3) Set: Increases the range of your Mend Pet spell by 50% and the effect by 10%. Also reduces the cost by 30%.
		3: func(agent core.Agent) {
			// Not implemented in sim
		},
		// (5) Set: Increases your pet's stamina by 30 and all spell resistances by 40.
		5: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			if hunter.pet == nil {
				return
			}
			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "Nature's Ally",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					hunter.pet.AddStatDynamic(sim, stats.Stamina, 30)
					hunter.pet.AddResistancesDynamic(sim, 40)
				},
			}))
		},
		// (8) Set: Increases the damage of Multi-shot and Volley by 15%.
		8: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.RegisterAura(core.Aura{
				Label: "Improved Volley and Multishot",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					hunter.Volley.BaseDamageMultiplierAdditive += 0.15
					hunter.MultiShot.BaseDamageMultiplierAdditive += 0.15
				},
			})
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            Phase 2 Item Sets - Dire Maul
///////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////
//                            Phase 3 Item Sets - BWL
///////////////////////////////////////////////////////////////////////////

var ItemSetDragonstalkersArmor = core.NewItemSet(core.ItemSet{
	Name: "Dragonstalker Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// (3) Set: Increases the Ranged Attack Power bonus of your Aspect of the Hawk by 20%.
		3: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "Improved Aspect of the Hawk",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					hunter.AspectOfTheHawkAPMultiplier += 0.25
				},
			}))
		},
		// (5) Set: Increases your pet's stamina by 40 and all spell resistances by 60.
		5: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			if hunter.pet == nil {
				return
			}

			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "Nature's Ally",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					hunter.pet.AddStatDynamic(sim, stats.Stamina, 40)
					hunter.pet.AddResistancesDynamic(sim, 60)
				},
			}))
		},
		// (8) Set: You have a chance whenever you deal ranged damage to apply an Expose Weakness effect to the target. Expose Weakness increases the Ranged Attack Power of all attackers against that target by 450 for 7 sec.
		8: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()

			debuffAuras := hunter.NewEnemyAuraArray(core.ExposeWeaknessAura)

			core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
				Name:     "T2 - Hunter - Ranged 8P Bonus Trigger",
				Callback: core.CallbackOnSpellHitDealt,
				Outcome:  core.OutcomeLanded,
				ProcMask: core.ProcMaskRanged,
				PPM:      0.5,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					debuffAuras.Get(result.Target).Activate(sim)
				},
			})
		},
	},
})

var ItemSetPredatorsArmor = core.NewItemSet(core.ItemSet{
	Name: "Predator's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// (2) Set: +20 Attack Power.
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.AddStats(stats.Stats{
				stats.AttackPower:       20,
				stats.RangedAttackPower: 20,
			})
		},
		// (3) Set: Decreases the cooldown of Concussive Shot by 1 sec.
		3: func(agent core.Agent) {
			// Concussive Shot not implemented in sim
		},
		// (5) Set: Increases the duration of Serpent Sting by 3 sec.
		5: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()

			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "Improved Serpent Sting",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					for _, dot := range hunter.SerpentSting.Dots() {
						if dot != nil {
							dot.NumberOfTicks += 1
							dot.RecomputeAuraDuration()
						}
					}
				},
			}))
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            Phase 4 Item Sets - AQ
///////////////////////////////////////////////////////////////////////////

// hhttps://www.wowhead.com/classic/item-set=515/beastmaster-armor
var ItemSetBeastmasterArmor = core.NewItemSet(core.ItemSet{
	Name: "Beastmaster Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// +8 All Resistances.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddResistances(8)
		},
		// Restores 8 mana per 5 sec.
		3: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.MP5, 8)
		},
		// Melee and ranged autoattacks have a 5% chance of restoring 200 mana. Classic had it
		// ranged only, at 4%.
		4: func(agent core.Agent) {
			c := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 450577}
			manaMetrics := c.NewManaMetrics(actionID)

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   actionID,
				Name:       "Hunter Armor Energize",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskWhiteHit,
				ProcChance: 0.05,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					if c.HasManaBar() {
						c.AddMana(sim, 200, manaMetrics)
					}
				},
			})
		},
		// Beast Unleashed: breaks a Root when struck.
		5: func(_ core.Agent) {},
		// +40 Attack Power.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStats(stats.Stats{
				stats.AttackPower:       40,
				stats.RangedAttackPower: 40,
			})
		},
	},
})

// https://www.wowhead.com/classic/item-set=509/strikers-garb
var ItemSetStrikersGarb = core.NewItemSet(core.ItemSet{
	Name: "Striker's Garb",
	Bonuses: map[int32]core.ApplyEffect{
		// (3) Set : Reduces the cost of your Arcane Shots by 10%.
		3: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "Striker's Arcane Shot Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					if hunter.ArcaneShot != nil {
						hunter.ArcaneShot.Cost.Multiplier -= 10.0
					}
				},
			}))
		},
		// (5) Set : Reduces the cooldown of your Rapid Fire ability by 2 minutes.
		5: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			core.MakePermanent(hunter.RegisterAura(core.Aura{
				Label: "Striker's Rapid Bonus",
				OnInit: func(aura *core.Aura, sim *core.Simulation) {
					hunter.RapidFire.CD.Duration -= time.Minute * 2
				},
			}))
		},
	},
})
