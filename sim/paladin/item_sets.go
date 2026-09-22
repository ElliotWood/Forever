package paladin

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

// Classic tier and PvP sets, ported from our master sim. Bonuses master does not model (aura
// radius, Judgement of Light chance, blessing duration, Holy Light cast time, Judgement duration,
// Hammer of Justice cooldown, the Lawbringer party heal) are left out here too.

var ItemSetLawbringerArmor = core.NewItemSet(core.ItemSet{
	Name: "Lawbringer Armor",
	Bonuses: map[int32]core.ApplySetBonus{
		// Improves your chance to get a critical strike with spells by 1%, and with melee by 1%.
		5: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatsBuff(stats.Stats{
				stats.MeleeCritRating: 1 * core.PhysicalCritRatingPerCritPercent,
				stats.SpellCritRating: 1 * core.SpellCritRatingPerCritPercent,
			})
		},
	},
})

var ItemSetFreethinkersArmor = core.NewItemSet(core.ItemSet{
	Name: "Freethinker's Armor",
	Bonuses: map[int32]core.ApplySetBonus{
		// Restores 4 mana per 5 sec.
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.MP5, 4)
		},
	},
})

var ItemSetAvengersBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Avenger's Battlegear",
	Bonuses: map[int32]core.ApplySetBonus{
		// Increases damage and healing done by magical spells and effects by up to 71.
		5: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatsBuff(stats.Stats{stats.SpellDamage: 71, stats.HealingPower: 71})
		},
	},
})

var ItemSetBattlegearOfEternalJustice = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Eternal Justice",
	Bonuses: map[int32]core.ApplySetBonus{
		// 20% chance to regain 100 mana when you cast a Judgement (master: Command and Righteousness).
		3: func(agent core.Agent, setBonusAura *core.Aura) {
			character := agent.GetCharacter()
			manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 26135})
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Battlegear of Eternal Justice - 3PC",
				ClassSpellMask: SpellMaskJudgementOfCommand | SpellMaskJudgementOfRighteousness,
				Callback:       core.CallbackOnCastComplete,
				ProcChance:     0.2,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					if character.HasManaBar() {
						character.AddMana(sim, 100, manaMetrics)
					}
				},
			})
		},
	},
})

var ItemSetLieutenantCommandersRedoubt = core.NewItemSet(core.ItemSet{
	Name: "Lieutenant Commander's Redoubt",
	Bonuses: map[int32]core.ApplySetBonus{
		// Increases damage and healing done by magical spells and effects by up to 23.
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatsBuff(stats.Stats{stats.SpellDamage: 23, stats.HealingPower: 23})
		},
		// +20 Stamina.
		6: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.Stamina, 20)
		},
	},
})

var ItemSetFieldMarshalsAegis = core.NewItemSet(core.ItemSet{
	Name: "Field Marshal's Aegis",
	Bonuses: map[int32]core.ApplySetBonus{
		// +20 Stamina.
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.Stamina, 20)
		},
		// Increases healing done by up to 44 and damage done by up to 15 for all magical spells and
		// effects (467550).
		6: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatsBuff(stats.Stats{stats.HealingPower: 44, stats.SpellDamage: 15})
		},
	},
})
