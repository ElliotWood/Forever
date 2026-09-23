package shared

import (
	"fmt"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// An item or enchant proc whose spell slows the enemy it lands on and does nothing else the sim
// models: Frostguard's Chilled 16927. BuffSpellID names the spell where it is not the trigger.
func NewSpellDataSlowProc(cfg SpellDataProc, variants []ItemVariant) {
	forEachSpellDataVariant(cfg, variants, registerSpellDataSlowProc)
}

func registerSpellDataSlowProc(cfg SpellDataProc) {
	source := cfg.effectSource()
	if source.isAlreadyImplemented() {
		return
	}

	trigger := cfg.trigger()
	slow := trigger
	if cfg.BuffSpellID != 0 {
		slow = spelldata.MustFind(cfg.BuffSpellID)
	}

	if !cfg.IsWeaponProc && decodedCallback(trigger) == core.CallbackEmpty {
		return
	}

	source.registerEffect(func(agent core.Agent) {
		applySpellDataSlowProc(agent, cfg, source, trigger, slow)
	})
}

func applySpellDataSlowProc(agent core.Agent, cfg SpellDataProc, source effectSource, trigger *spelldata.Spell, slow *spelldata.Spell) {
	character := agent.GetCharacter()
	slows := slowAuras(character, slow)

	config := spellDataDamageTrigger(character, cfg, source, trigger)
	callback := config.Callback
	config.Handler = func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		slows.Get(procDamageTarget(character, callback, spell, result)).Activate(sim)
	}
	config.TriggerImmediately = true

	source.registerProc(character, character.MakeProcTriggerAura(config), source.eligibleSlots(character))
}

// The row's slow on each enemy, for its duration. The aura is the row's rather than the wearer's, so
// two wearers of one proc refresh a single slow on the target instead of multiplying two. Each slow
// takes its exclusive category, attacks Thunder Clap's and casts Slow's, where only the strongest
// applies.
func slowAuras(character *core.Character, row *spelldata.Spell) core.AuraArray {
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
		return aura
	})
}

// The slow a proc's damage spell puts on each enemy its hit lands on, or nil where the row states
// none. A spell that deals only damage over time has no hit, and its target takes the slow with it.
func slowOnLanding(character *core.Character, row *spelldata.Spell) afterDealt {
	if len(row.SlowEffects()) == 0 {
		return nil
	}

	slows := slowAuras(character, row)
	return func(sim *core.Simulation, _ *core.Spell, target *core.Unit, results core.SpellResultSlice) {
		if len(results) == 0 {
			slows.Get(target).Activate(sim)
			return
		}
		for _, result := range results {
			if result.Landed() {
				slows.Get(result.Target).Activate(sim)
			}
		}
	}
}
