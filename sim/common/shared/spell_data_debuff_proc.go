package shared

import (
	"fmt"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// An item or enchant proc whose spell puts a debuff on the enemy it lands on and does nothing else
// the sim models: Frostguard's Chilled 16927 slows its attacks, Annihilator's Armor Shatter 16928
// takes its armor. BuffSpellID names the spell where it is not the trigger.
func NewSpellDataDebuffProc(cfg SpellDataProc, variants []ItemVariant) {
	forEachSpellDataVariant(cfg, variants, registerSpellDataDebuffProc)
}

func registerSpellDataDebuffProc(cfg SpellDataProc) {
	source := cfg.effectSource()
	if source.isAlreadyImplemented() {
		return
	}

	trigger := cfg.trigger()
	debuff := trigger
	if cfg.BuffSpellID != 0 {
		debuff = spelldata.MustFind(cfg.BuffSpellID)
	}

	if !cfg.IsWeaponProc && decodedCallback(trigger) == core.CallbackEmpty {
		return
	}

	source.registerEffect(func(agent core.Agent) {
		applySpellDataDebuffProc(agent, cfg, source, trigger, debuff)
	})
}

func applySpellDataDebuffProc(agent core.Agent, cfg SpellDataProc, source effectSource, trigger *spelldata.Spell, debuff *spelldata.Spell) {
	character := agent.GetCharacter()
	debuffs := debuffAuras(character, debuff)

	config := spellDataDamageTrigger(character, cfg, source, trigger)
	callback := config.Callback
	config.Handler = func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		applyDebuff(sim, debuffs.Get(procDamageTarget(character, callback, spell, result)))
	}
	config.TriggerImmediately = true

	source.registerProc(character, character.MakeProcTriggerAura(config), source.eligibleSlots(character))
}

// A re-application refreshes the duration and, on a row that stacks, adds a stack up to its count.
func applyDebuff(sim *core.Simulation, aura *core.Aura) {
	aura.Activate(sim)
	if aura.MaxStacks > 0 {
		aura.AddStack(sim)
	}
}

// The row's debuff on each enemy, for its duration. The aura is the row's rather than the wearer's, so
// two wearers of one proc refresh a single debuff on the target instead of applying two. Each slow
// takes its exclusive category, attacks Thunder Clap's and casts Slow's, where only the strongest
// applies; a stat change applies to the target as the row states it, per stack.
func debuffAuras(character *core.Character, row *spelldata.Spell) core.AuraArray {
	label := fmt.Sprintf("%s %d", row.Name, row.ID)
	return character.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		if aura := target.GetAura(label); aura != nil {
			return aura
		}

		aura := target.RegisterAura(spelldata.AuraConfig(row, spelldata.Label(label)))
		for _, i := range row.SlowEffects() {
			effect := row.EffectN(int(i))
			slowedTime := core.SlowedTimeMultiplier(effect.Average(character.Level))
			if effect.Aura == dbcenums.A_MOD_CASTING_SPEED_NOT_STACK {
				core.CastSpeedReductionEffect(aura, slowedTime)
			} else {
				core.AtkSpeedReductionEffect(aura, slowedTime)
			}
		}
		if stats := row.StatDebuffEffects(); len(stats) > 0 {
			spelldata.ParseEffects(nil, aura, row, spelldata.Effects(stats...))
		}
		return aura
	})
}

// The debuff a proc's damage spell puts on each enemy its hit lands on, or nil where the row states
// none. A spell that deals only damage over time has no hit, and its target takes the debuff with it.
func debuffOnLanding(character *core.Character, row *spelldata.Spell) afterDealt {
	if !row.DebuffsTheTarget() {
		return nil
	}

	debuffs := debuffAuras(character, row)
	return func(sim *core.Simulation, _ *core.Spell, target *core.Unit, results core.SpellResultSlice) {
		if len(results) == 0 {
			applyDebuff(sim, debuffs.Get(target))
			return
		}
		for _, result := range results {
			if result.Landed() {
				applyDebuff(sim, debuffs.Get(result.Target))
			}
		}
	}
}
