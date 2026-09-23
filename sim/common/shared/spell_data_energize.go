package shared

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// An on-use item whose spell restores the wearer's mana, rage or energy, at once or over time, on the
// item's cooldowns with the cast time and global cooldown the row states. A character without that
// bar does not register it. The cooldown manager uses it once the whole gain fits under the bar's
// maximum.
func NewSpellDataEnergizeOnUse(itemID int32) {
	// Soft fail to allow for overrides for bad effects
	if core.HasItemEffect(itemID) {
		return
	}

	core.NewItemEffect(itemID, func(agent core.Agent) {
		character := agent.GetCharacter()

		for _, itemEffect := range onUseEffectsFor(itemID) {
			row := spelldata.MustFind(itemEffect.BuffId)
			effect := row.ProcEnergizeEffect()
			power := dbcenums.PowerType(effect.Misc)
			bar, ok := energizedBarOf(character, power, core.ActionID{ItemID: itemID})
			if !ok {
				continue
			}

			amount := func(sim *core.Simulation) float64 {
				gain := effect.Roll(sim, character.Level)
				if power.InTenths() {
					gain /= 10
				}
				return gain
			}

			config := spelldata.SpellConfig(&character.Unit, row)
			config.ActionID = core.ActionID{ItemID: itemID}

			onUse := onUseCast(character, itemEffect)
			config.Cast = spelldata.Cast(row)
			config.Cast.CD, config.Cast.SharedCD = onUse.CD, onUse.SharedCD
			config.Flags &^= core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete |
				core.SpellFlagNoOnDamageDealt | core.SpellFlagProc

			ticks := 1.0
			if effect.Aura == dbcenums.A_PERIODIC_ENERGIZE {
				config.Hot = spelldata.DotConfig(row, effect)
				config.Hot.SelfOnly = true
				config.Hot.OnTick = func(sim *core.Simulation, _ *core.Unit, _ *core.Dot) {
					bar.add(sim, amount(sim), bar.metrics)
				}
				config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
					spell.SelfHot().Apply(sim)
				}
				ticks = float64(config.Hot.NumberOfTicks)
			} else {
				config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
					bar.add(sim, amount(sim), bar.metrics)
				}
			}

			whole := effect.Max(character.Level) * ticks
			if power.InTenths() {
				whole /= 10
			}

			character.AddMajorCooldown(core.MajorCooldown{
				Spell: character.RegisterSpell(config),
				Type:  bar.cdType,
				ShouldActivate: func(_ *core.Simulation, _ *core.Character) bool {
					return bar.room() >= whole
				},
			})
		}
	})
}

// The bar a gain fills: the room left in it, the character's own gain on it metered under the item,
// and the cooldown type core files a gain of that bar under.
type energizedBar struct {
	room    func() float64
	add     func(*core.Simulation, float64, *core.ResourceMetrics)
	metrics *core.ResourceMetrics
	cdType  core.CooldownType
}

func energizedBarOf(character *core.Character, power dbcenums.PowerType, actionID core.ActionID) (energizedBar, bool) {
	switch {
	case power == dbcenums.POWER_MANA && character.HasManaBar():
		return energizedBar{
			room:    func() float64 { return character.MaxMana() - character.CurrentMana() },
			add:     character.AddMana,
			metrics: character.NewManaMetrics(actionID),
			cdType:  core.CooldownTypeMana,
		}, true
	case power == dbcenums.POWER_RAGE && character.HasRageBar():
		return energizedBar{
			room:    func() float64 { return character.MaximumRage() - character.CurrentRage() },
			add:     character.AddRage,
			metrics: character.NewRageMetrics(actionID),
			cdType:  core.CooldownTypeDPS,
		}, true
	case power == dbcenums.POWER_ENERGY && character.HasEnergyBar():
		return energizedBar{
			room:    func() float64 { return character.MaximumEnergy() - character.CurrentEnergy() },
			add:     character.AddEnergy,
			metrics: character.NewEnergyMetrics(actionID),
			cdType:  core.CooldownTypeDPS,
		}, true
	}
	return energizedBar{}, false
}
