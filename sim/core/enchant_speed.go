package core

import (
	"fmt"

	"github.com/wowsims/forever/sim/core/proto"
)

// An enchant's haste pseudo stats multiply the wearer's melee, ranged and cast speed while the
// enchanted item is equipped. Each enchanted item applies its own.
func (character *Character) registerEnchantSpeedAuras() {
	for slot := proto.ItemSlot(0); slot < NumItemSlots; slot++ {
		enchants := []Enchant{character.Equipment[slot].Enchant}
		if character.ItemSwap.IsEnabled() {
			enchants = append(enchants, character.ItemSwap.unEquippedItems[slot].Enchant)
		}

		for idx, enchant := range enchants {
			pseudoStat := func(pseudoStat proto.PseudoStat) float64 {
				if int(pseudoStat) < len(enchant.PseudoStats) {
					return enchant.PseudoStats[pseudoStat]
				}
				return 0
			}
			melee := pseudoStat(proto.PseudoStat_PseudoStatMeleeHastePercent)
			ranged := pseudoStat(proto.PseudoStat_PseudoStatRangedHastePercent)
			cast := pseudoStat(proto.PseudoStat_PseudoStatSpellHastePercent)
			if melee == 0 && ranged == 0 && cast == 0 {
				continue
			}
			if idx > 0 && enchant.EffectID == enchants[0].EffectID {
				continue
			}

			activeAtStart := idx == 0
			aura := character.GetOrRegisterAura(Aura{
				Label:      fmt.Sprintf("Enchant %d Speed (%s)", enchant.EffectID, slot),
				BuildPhase: Ternary(activeAtStart, CharacterBuildPhaseGear, CharacterBuildPhaseNone),
				Duration:   NeverExpires,
			})
			if activeAtStart {
				aura = MakePermanent(aura)
			}
			if melee != 0 {
				aura.AttachMultiplyMeleeSpeed(1 + melee/100)
			}
			if ranged != 0 {
				aura.AttachMultiplyRangedSpeed(1 + ranged/100)
			}
			if cast != 0 {
				aura.AttachMultiplyCastSpeed(1 + cast/100)
			}

			character.ItemSwap.RegisterEnchantProcWithSlots(enchant.EffectID, aura, []proto.ItemSlot{slot})
		}
	}
}
