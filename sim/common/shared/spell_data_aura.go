package shared

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// The auras a row puts on the wearer, on each of its pets and on an enemy, each carrying the effects
// that land on that unit. The wearer's and the pets' are registered from the config handed in; the
// enemy's is the row's own aura, one per enemy, shared by every character that applies it, so the
// first to register it parses it and the multiplier reaches every attacker's hits once.
type spellDataAuras struct {
	wearer  *core.Aura
	pets    []*core.Aura
	enemies core.AuraArray
}

func newSpellDataAuras(character *core.Character, row *spelldata.Spell, config core.Aura) *spellDataAuras {
	auras := &spellDataAuras{}

	if effects := spelldata.EffectsOn(row, spelldata.AuraOnWearer); len(effects) > 0 {
		auras.wearer = character.RegisterAura(config)
		spelldata.ParseEffects(character, auras.wearer, row, spelldata.Effects(effects...))
	}

	if effects := spelldata.EffectsOn(row, spelldata.AuraOnPet); len(effects) > 0 {
		for _, pet := range character.Pets {
			if pet.IsGuardian() {
				continue
			}
			aura := pet.RegisterAura(config)
			spelldata.ParseEffects(&pet.Character, aura, row, spelldata.Effects(effects...))
			auras.pets = append(auras.pets, aura)
		}
	}

	if effects := spelldata.EffectsOn(row, spelldata.AuraOnEnemy); len(effects) > 0 {
		debuff := spelldata.AuraConfig(row)
		debuff.Duration = config.Duration
		auras.enemies = character.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
			if aura := target.GetAura(debuff.Label); aura != nil {
				return aura
			}
			aura := target.RegisterAura(debuff)
			spelldata.ParseEffects(nil, aura, row, spelldata.Effects(effects...))
			return aura
		})
	}

	return auras
}

func (auras *spellDataAuras) activate(sim *core.Simulation, target *core.Unit) {
	auras.wearer.Activate(sim)
	for _, aura := range auras.pets {
		if aura.Unit.IsEnabled() {
			aura.Activate(sim)
		}
	}
	auras.enemies.Get(target).Activate(sim)
}

// An on-use item whose spell applies auras: the wearer's buff, its pets' and the debuff on the enemy
// it is used on, for the row's duration.
func NewSpellDataAuraOnUse(itemID int32) {
	registerSpellDataOnUse(itemID, core.CooldownTypeDPS, spellDataOnUseAuraSpell)
}

func spellDataOnUseAuraSpell(character *core.Character, row *spelldata.Spell) core.SpellConfig {
	auras := newSpellDataAuras(character, row, spelldata.AuraConfig(row))

	return core.SpellConfig{
		SpellSchool: row.SpellSchool(),
		ProcMask:    core.ProcMaskEmpty,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			auras.activate(sim, target)
		},
	}
}
