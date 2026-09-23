package core

import (
	"fmt"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// The haste pseudo stats of an item and of its enchant multiply the wearer's melee, ranged and cast
// speed while the item is equipped. Each item and each enchanted item applies its own.
func (character *Character) registerEquipSpeedAuras() {
	for slot := proto.ItemSlot(0); slot < NumItemSlots; slot++ {
		items := []Item{character.Equipment[slot]}
		if character.ItemSwap.IsEnabled() {
			items = append(items, character.ItemSwap.unEquippedItems[slot])
		}

		for idx, item := range items {
			activeAtStart := idx == 0
			if idx == 0 || item.ID != items[0].ID {
				label := fmt.Sprintf("Item %d Speed (%s)", item.ID, slot)
				if aura := character.registerSpeedAura(label, item.PseudoStats, activeAtStart); aura != nil {
					character.ItemSwap.RegisterProcWithSlots(item.ID, aura, []proto.ItemSlot{slot})
				}
			}
			if idx == 0 || item.Enchant.EffectID != items[0].Enchant.EffectID {
				label := fmt.Sprintf("Enchant %d Speed (%s)", item.Enchant.EffectID, slot)
				if aura := character.registerSpeedAura(label, item.Enchant.PseudoStats, activeAtStart); aura != nil {
					character.ItemSwap.RegisterEnchantProcWithSlots(item.Enchant.EffectID, aura, []proto.ItemSlot{slot})
				}
			}
		}
	}
}

func (character *Character) registerSpeedAura(label string, pseudoStats []float64, activeAtStart bool) *Aura {
	if melee, ranged, cast := hastePercents(pseudoStats); melee == 0 && ranged == 0 && cast == 0 {
		return nil
	}

	aura := character.GetOrRegisterAura(Aura{
		Label:      label,
		BuildPhase: Ternary(activeAtStart, CharacterBuildPhaseGear, CharacterBuildPhaseNone),
		Duration:   NeverExpires,
	})
	if activeAtStart {
		aura = MakePermanent(aura)
	}
	return aura.AttachHastePseudoStats(pseudoStats)
}

func hastePercents(pseudoStats []float64) (melee, ranged, cast float64) {
	return stats.PseudoStatValue(pseudoStats, proto.PseudoStat_PseudoStatMeleeHastePercent),
		stats.PseudoStatValue(pseudoStats, proto.PseudoStat_PseudoStatRangedHastePercent),
		stats.PseudoStatValue(pseudoStats, proto.PseudoStat_PseudoStatSpellHastePercent)
}

// Multiplies the unit's melee, ranged and cast speed by the haste percents pseudoStats states, indexed
// by proto.PseudoStat, while the aura is up.
func (aura *Aura) AttachHastePseudoStats(pseudoStats []float64) *Aura {
	melee, ranged, cast := hastePercents(pseudoStats)
	if melee != 0 {
		aura.AttachMultiplyMeleeSpeed(1 + melee/100)
	}
	if ranged != 0 {
		aura.AttachMultiplyRangedSpeed(1 + ranged/100)
	}
	if cast != 0 {
		aura.AttachMultiplyCastSpeed(1 + cast/100)
	}
	return aura
}
