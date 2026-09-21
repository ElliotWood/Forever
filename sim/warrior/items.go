package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

// The rows ItemSetSpell names for each set's thresholds, and the buffs those rows grant. None of
// them is a warrior class spell, so the generated table carries no reference to them.
var (
	mightBlockValue   = spelldata.MustFind(23562)
	mightRageProc     = spelldata.MustFind(21838)
	mightRageEnergize = spelldata.MustFind(29478)
	mightSunderThreat = spelldata.MustFind(23561)

	wrathDiscountProc = spelldata.MustFind(21890)
	wrathDiscountBuff = spelldata.MustFind(21887)
	wrathParryProc    = spelldata.MustFind(23548)
	wrathParryBuff    = spelldata.MustFind(23547)

	conquerorShoutCost   = spelldata.MustFind(26109)
	conquerorThunderClap = spelldata.MustFind(26110)

	dreadnaughtRevenge        = spelldata.MustFind(28844)
	dreadnaughtTauntHit       = spelldata.MustFind(28843)
	dreadnaughtAbilityHit     = spelldata.MustFind(28842)
	dreadnaughtCheatDeath     = spelldata.MustFind(28845)
	dreadnaughtCheatDeathBuff = spelldata.MustFind(28846)
)

var ItemSetBattlegearOfMight = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Might",
	ID:   209,
	Bonuses: map[int32]core.ApplySetBonus{
		3: func(agent core.Agent, setBonusAura *core.Aura) {
			// A_MOD_BLOCK_VALUE_FLAT has no sim kind in the parse table, so only the amount comes
			// off the row.
			setBonusAura.AttachStatBuff(stats.BlockValue, mightBlockValue.EffectN(1).BaseValue())
		},
		5: func(agent core.Agent, setBonusAura *core.Aura) {
			warrior := agent.(WarriorAgent).GetWarrior()
			rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: mightRageEnergize.ID})
			rage := mightRageEnergize.EnergizeEffect().Tenths()

			// The rate is the proc chance column, which the tooltip's $h% says is a real roll.
			trigger := spelldata.ProcTrigger(&warrior.Character, mightRageProc,
				func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					warrior.AddRage(sim, rage, rageMetrics)
				})
			trigger.Name = "Battlegear of Might - 5PC"

			setBonusAura.AttachProcTrigger(trigger)
		},
		8: func(agent core.Agent, setBonusAura *core.Aura) {
			// The row states A_ADD_PCT_MODIFIER on the threat, which the parse maps to
			// SpellMod_ThreatMultiplier_Pct: that scales the threat a spell's damage makes, and
			// Sunder Armor deals none - all of its threat is the flat bonus. So only the amount
			// and the spells it names come off the row.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassFlags: mightSunderThreat.EffectN(1).ClassFlags,
				Kind:       core.SpellMod_FlatThreatBonus_Pct,
				FloatValue: mightSunderThreat.EffectN(1).Percent(),
			})
		},
	},
})

var ItemSetBattlegearOfWrath = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Wrath",
	ID:   218,
	Bonuses: map[int32]core.ApplySetBonus{
		3: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 23563 states 30 attack power on Battle Shout, which battle_shout.go adds through
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

			// The abilities the discount names, which are also the ones the tooltip fires the
			// proc off: the buff's own modifier effect states them and 21890's proc effect
			// states none.
			discounted := wrathDiscountBuff.EffectN(1).ClassFlags

			buff := warrior.RegisterAura(spelldata.AuraConfig(wrathDiscountBuff))
			spelldata.ParseEffects(&warrior.Character, buff, wrathDiscountBuff)

			// The discount is spent at cast, so the buff goes when an ability is cast rather than
			// when one lands, which is what the row's proc flags state.
			buff.AttachProcTrigger(core.ProcTrigger{
				Name:               "Warrior's Wrath - Consume",
				ClassFlags:         discounted,
				Callback:           core.CallbackOnCastComplete,
				TriggerImmediately: true,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					buff.Deactivate(sim)
				},
			})

			// The rate is the proc chance column, which the tooltip's $h% says is a real roll.
			trigger := spelldata.ProcTrigger(&warrior.Character, wrathDiscountProc,
				func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					buff.Activate(sim)
				})
			trigger.Name = "Battlegear of Wrath - 5PC"
			trigger.ClassFlags = discounted

			setBonusAura.AttachProcTrigger(trigger)
		},
		8: func(agent core.Agent, setBonusAura *core.Aura) {
			warrior := agent.(WarriorAgent).GetWarrior()

			parry := warrior.RegisterAura(spelldata.AuraConfig(wrathParryBuff,
				spelldata.Label("Battlegear of Wrath Parry")))
			spelldata.ParseEffects(&warrior.Character, parry, wrathParryBuff)

			// The attack that spends the buff is the one it parries, which is an outcome no proc
			// mask states and one that deals no damage: both of the row's defaults go.
			consume := spelldata.ProcTrigger(&warrior.Character, wrathParryBuff,
				func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					parry.Deactivate(sim)
				})
			consume.Name = "Battlegear of Wrath - 8PC Consume"
			consume.Outcome = core.OutcomeEmpty
			consume.RequireDamageDealt = false
			consume.TriggerImmediately = true
			parry.AttachProcTrigger(consume)

			// The rate is the proc chance column, which the tooltip's $h% says is a real roll; the
			// block it fires on is an outcome no proc mask states.
			trigger := spelldata.ProcTrigger(&warrior.Character, wrathParryProc,
				func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					parry.Activate(sim)
				})
			trigger.Name = "Battlegear of Wrath - 8PC"
			trigger.Outcome = core.OutcomeBlock
			trigger.RequireDamageDealt = false

			setBonusAura.AttachProcTrigger(trigger)
		},
	},
})

var ItemSetConquerorsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Conqueror's Battlegear",
	ID:   496,
	Bonuses: map[int32]core.ApplySetBonus{
		3: func(agent core.Agent, setBonusAura *core.Aura) {
			warrior := agent.(WarriorAgent).GetWarrior()
			spelldata.ParseEffects(&warrior.Character, setBonusAura, conquerorShoutCost)
		},
		5: func(agent core.Agent, setBonusAura *core.Aura) {
			warrior := agent.(WarriorAgent).GetWarrior()
			spelldata.ParseEffects(&warrior.Character, setBonusAura, conquerorThunderClap)

			// The one amount the row states raises the slow the tooltip names alongside the
			// damage, which thunder_clap.go reads.
			slow := conquerorThunderClap.EffectN(1).Percent()
			setBonusAura.ApplyOnGain(func(_ *core.Aura, _ *core.Simulation) {
				warrior.thunderClapEffectBonus += slow
			})
			setBonusAura.ApplyOnExpire(func(_ *core.Aura, _ *core.Simulation) {
				warrior.thunderClapEffectBonus -= slow
			})
		},
	},
})

var ItemSetDreadnaughtsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Dreadnaught's Battlegear",
	ID:   523,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// A_ADD_FLAT_MODIFIER on the damage has no sim kind in the parse table, so only the
			// amount and the spell it names come off the row.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassFlags: dreadnaughtRevenge.EffectN(1).ClassFlags,
				Kind:       core.SpellMod_BaseDamage_Flat,
				FloatValue: dreadnaughtRevenge.EffectN(1).BaseValue(),
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// TODO: both spells resolve on OutcomeAlwaysHit, so the bonus changes nothing until
			// they roll against the spell hit table.
			warrior := agent.(WarriorAgent).GetWarrior()
			spelldata.ParseEffects(&warrior.Character, setBonusAura, dreadnaughtTauntHit)
		},
		6: func(agent core.Agent, setBonusAura *core.Aura) {
			warrior := agent.(WarriorAgent).GetWarrior()
			spelldata.ParseEffects(&warrior.Character, setBonusAura, dreadnaughtAbilityHit)
		},
		8: func(agent core.Agent, setBonusAura *core.Aura) {
			// Spell 28845 states that below 20% health, healing spells cast on you gain up to 160
			// healing (spell 28846) for 5 seconds. The modelled incoming healing of the tank sim
			// bypasses the healing bonus; heals a healer unit casts on the warrior take it.
			warrior := agent.(WarriorAgent).GetWarrior()

			cheatDeath := warrior.RegisterAura(spelldata.AuraConfig(dreadnaughtCheatDeathBuff))
			spelldata.ParseEffects(&warrior.Character, cheatDeath, dreadnaughtCheatDeathBuff)

			trigger := spelldata.ProcTrigger(&warrior.Character, dreadnaughtCheatDeath,
				func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					cheatDeath.Activate(sim)
				})
			trigger.Name = "Cheat Death - Trigger"
			// 28845 sits in class family 5 and its proc effect names a spell mask there, which
			// would restrict the listener to another class's spells; the tooltip restricts it to
			// nothing but the health threshold, which the row does not state either.
			trigger.ClassFlags = core.ClassFlags{}
			trigger.ExtraCondition = func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool {
				return warrior.CurrentHealthPercent() < 0.2
			}

			setBonusAura.AttachProcTrigger(trigger)
		},
	},
})
