package item_sets

import (
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

// Keep these in alphabetical order.

// https://www.wowhead.com/classic/item-set=489/black-dragon-mail
var ItemSetBlackDragonMail = core.NewItemSet(core.ItemSet{
	Name: "Black Dragon Mail",
	ID:   489,
	Bonuses: map[int32]core.ApplyEffect{
		// Improves your chance to hit by 1%.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeHit, 1)
		},
		// Improves your chance to get a critical strike by 2%.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 2*core.CritRatingPerCritChance)
		},
		// +10 Fire Resistance.
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.FireResistance, 10)
		},
	},
})

// https://www.wowhead.com/classic/item-set=491/blue-dragon-mail
var ItemSetBlueDragonMail = core.NewItemSet(core.ItemSet{
	Name: "Blue Dragon Mail",
	ID:   491,
	Bonuses: map[int32]core.ApplyEffect{
		// +4 All Resistances.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddResistances(4)
		},
		// Increases damage and healing done by magical spells and effects by up to 28.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 28)
		},
	},
})

// https://www.wowhead.com/classic/item-set=443/bloodsoul-embrace
var ItemSetBloodsoulEmbrace = core.NewItemSet(core.ItemSet{
	Name: "Bloodsoul Embrace",
	Bonuses: map[int32]core.ApplyEffect{
		// Restores 12 mana per 5 sec.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MP5, 12)
		},
	},
})

// https://www.wowhead.com/classic/item-set=421/bloodvine-garb
var ItemSetBloodvineGarb = core.NewItemSet(core.ItemSet{
	Name: "Bloodvine Garb",
	Bonuses: map[int32]core.ApplyEffect{
		// Improves your chance to get a critical strike with spells by 2%.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			if character.HasProfession(proto.Profession_Tailoring) {
				character.AddStat(stats.SpellCrit, 2*core.SpellCritRatingPerCritChance)
			}
		},
	},
})

// https://www.wowhead.com/classic/item-set=442/blood-tiger-harness
var ItemSetBloodTigerHarness = core.NewItemSet(core.ItemSet{
	Name: "Blood Tiger Harness",
	Bonuses: map[int32]core.ApplyEffect{
		// Improves your chance to get a critical strike by 1%.
		// Improves your chance to get a critical strike with spells by 1%.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance)
			character.AddStat(stats.SpellCrit, 1*core.SpellCritRatingPerCritChance)
		},
	},
})

// https://www.wowhead.com/classic/item-set=143/devilsaur-armor
var ItemSetDevilsaurArmor = core.NewItemSet(core.ItemSet{
	Name: "Devilsaur Armor",
	ID:   143,
	Bonuses: map[int32]core.ApplyEffect{
		// Improves your chance to hit by 2%. Forever's 460230 adds the spell half (aura 55).
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeHit, 2*core.MeleeHitRatingPerHitChance)
			character.AddStat(stats.SpellHit, 2*core.SpellHitRatingPerHitChance)
		},
	},
})

// https://www.wowhead.com/classic/item-set=490/green-dragon-mail
var ItemSetGreenDragonMail = core.NewItemSet(core.ItemSet{
	Name: "Green Dragon Mail",
	ID:   490,
	Bonuses: map[int32]core.ApplyEffect{
		// Restores 3 mana per 5 sec.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MP5, 3)
		},
		// Allows 15% of your Mana regeneration to continue while casting.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.PseudoStats.SpiritRegenRateCasting += .15
		},
	},
})

// https://www.wowhead.com/classic/item-set=321/imperial-plate
var ItemSetImperialPlate = core.NewItemSet(core.ItemSet{
	Name: "Imperial Plate",
	// Rebuilt in Forever: 2/4/6 became 2/3/4/5/6 and every bonus changed (ItemSetSpell, beta client
	// 1.60.1.69893). Classic's armour, attack power and stamina are all gone.
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +7 (13385).
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Defense, 7)
		},
		// Improves your chance to hit by 1% (1251990).
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeHit, 1*core.MeleeHitRatingPerHitChance)
			character.AddStat(stats.SpellHit, 1*core.SpellHitRatingPerHitChance)
		},
		// Imperial March: movement speed cannot drop below 80%.
		4: func(_ core.Agent) {},
		// +20 Strength (1251984).
		5: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Strength, 20)
		},
		// Improves your chance to crit by 1% (1251991).
		6: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.MeleeCrit, 1*core.CritRatingPerCritChance)
			character.AddStat(stats.SpellCrit, 1*core.SpellCritRatingPerCritChance)
		},
	},
})

// https://www.wowhead.com/classic/item-set=144/ironfeather-armor
var ItemSetIronfeatherArmor = core.NewItemSet(core.ItemSet{
	Name: "Ironfeather Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 20.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.SpellPower, 20)
		},
	},
})

// https://www.wowhead.com/classic/item-set=142/stormshroud-armor
var ItemSetStormshroudArmor = core.NewItemSet(core.ItemSet{
	Name: "Stormshroud Armor",
	Bonuses: map[int32]core.ApplyEffect{
		// 5% chance of dealing 15 to 25 Nature damage on a successful melee attack.
		2: func(a core.Agent) {
			char := a.GetCharacter()
			proc := char.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 18980},
				SpellSchool: core.SpellSchoolNature,
				DefenseType: core.DefenseTypeMagic,
				ProcMask:    core.ProcMaskEmpty,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, sim.Roll(15, 25), spell.OutcomeMagicHitAndCrit)
				},
			})
			char.RegisterAura(core.Aura{
				Label:    "Lightning",
				ActionID: core.ActionID{SpellID: 18979},
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
					if sim.RandomFloat("Lightning") < 0.05 {
						proc.Cast(sim, result.Target)
					}
				},
			})
		},
		// 2% chance on melee attack of restoring 30 energy.
		3: func(a core.Agent) {
			char := a.GetCharacter()
			if !char.HasEnergyBar() {
				return
			}
			metrics := char.NewEnergyMetrics(core.ActionID{SpellID: 23863})
			proc := char.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 23864},
				SpellSchool: core.SpellSchoolNature,
				ApplyEffects: func(sim *core.Simulation, u *core.Unit, spell *core.Spell) {
					char.AddEnergy(sim, 30, metrics)
				},
			})
			char.RegisterAura(core.Aura{
				Label:    "Revitalize",
				ActionID: core.ActionID{SpellID: 18979},
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
					if sim.RandomFloat("Revitalize") < 0.02 {
						proc.Cast(sim, result.Target)
					}
				},
			})

		},
		// +14 Attack Power.
		4: func(a core.Agent) {
			a.GetCharacter().AddStat(stats.AttackPower, 14)
			a.GetCharacter().AddStat(stats.RangedAttackPower, 14)
		},
	},
})

// https://www.wowhead.com/classic/item-set=444/the-darksoul
var ItemSetTheDarksoul = core.NewItemSet(core.ItemSet{
	Name: "The Darksoul",
	Bonuses: map[int32]core.ApplyEffect{
		// Increased Defense +20.
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStat(stats.Defense, 20)
		},
	},
})

// https://www.wowhead.com/classic/item-set=141/volcanic-armor
var ItemSetVolcanicArmor = core.NewItemSet(core.ItemSet{
	Name: "Volcanic Armor",
	ID:   141,
	Bonuses: map[int32]core.ApplyEffect{
		// 5% chance of dealing 15 to 25 Fire damage on a successful melee attack.
		3: func(agent core.Agent) {
			character := agent.GetCharacter()
			procSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 9057},
				SpellSchool: core.SpellSchoolFire,
				DefenseType: core.DefenseTypeMagic,
				ProcMask:    core.ProcMaskEmpty,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, sim.Roll(15, 25), spell.OutcomeMagicHitAndCrit)
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Firebolt Trigger (Volcanic Armor)",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: .05,
				Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					procSpell.Cast(sim, result.Target)
				},
			})
		},
	},
})
