package bulk

import (
	"slices"

	"github.com/wowsims/forever/sim/core/proto"
)

func (generator *bulkSimCandidateGenerator) getAllWeaponCombos() [][2]*bulkSimCandidateOption {
	allWeaponCombos := make([][2]*bulkSimCandidateOption, 0)
	all2HWeapons := make([]bulkSimCandidateOption, 0)
	for _, bulkSlot := range []BulkSimItemSlot{BulkSimItemSlotMainHand, BulkSimItemSlotHandWeapon} {
		for _, option := range generator.selectedByBulkSlot[bulkSlot] {
			if option.item.HandType == proto.HandType_HandTypeTwoHand {
				all2HWeapons = append(all2HWeapons, option)
			}
		}
	}

	for i := range all2HWeapons {
		allWeaponCombos = append(allWeaponCombos, [2]*bulkSimCandidateOption{&all2HWeapons[i], nil})
	}

	mhOptions := generator.selectedByBulkSlot[BulkSimItemSlotMainHand]
	ohOptions := generator.selectedByBulkSlot[BulkSimItemSlotOffHand]
	if len(mhOptions) > 0 {
		for i := range mhOptions {
			if optionsContainEquivalent(all2HWeapons, mhOptions[i]) {
				continue
			}
			if len(ohOptions) > 0 {
				for j := range ohOptions {
					allWeaponCombos = append(allWeaponCombos, [2]*bulkSimCandidateOption{&mhOptions[i], &ohOptions[j]})
				}
			} else {
				allWeaponCombos = append(allWeaponCombos, [2]*bulkSimCandidateOption{&mhOptions[i], nil})
			}
		}
	} else if len(ohOptions) > 0 {
		for i := range ohOptions {
			allWeaponCombos = append(allWeaponCombos, [2]*bulkSimCandidateOption{nil, &ohOptions[i]})
		}
	}

	oneHandOptions := generator.selectedByBulkSlot[BulkSimItemSlotHandWeapon]
	if len(oneHandOptions) > 0 {
		unique := make([]bulkSimCandidateOption, 0, len(oneHandOptions))
		for _, option := range oneHandOptions {
			if optionsContainEquivalent(all2HWeapons, option) {
				continue
			}
			if optionsContainEquivalent(unique, option) {
				continue
			}
			unique = append(unique, option)
		}

		for i := range unique {
			iCanMH := unique[i].item.HandType != proto.HandType_HandTypeOffHand
			iCanOH := unique[i].item.HandType != proto.HandType_HandTypeMainHand
			// Only wield the same 1H weapon in both hands when two copies exist and the weapon
			// itself allows it (not unique, no limit category).
			if generator.hasTwoCopies(unique[i], unique[i].item) && iCanMH && iCanOH {
				allWeaponCombos = append(allWeaponCombos, [2]*bulkSimCandidateOption{&unique[i], &unique[i]})
			}
			for j := i + 1; j < len(unique); j++ {
				jCanMH := unique[j].item.HandType != proto.HandType_HandTypeOffHand
				jCanOH := unique[j].item.HandType != proto.HandType_HandTypeMainHand
				if iCanMH && jCanOH {
					allWeaponCombos = append(allWeaponCombos, [2]*bulkSimCandidateOption{&unique[i], &unique[j]})
				}
				if jCanMH && iCanOH {
					allWeaponCombos = append(allWeaponCombos, [2]*bulkSimCandidateOption{&unique[j], &unique[i]})
				}
			}
		}
	}

	filteredCombos := make([][2]*bulkSimCandidateOption, 0, len(allWeaponCombos))
	for _, combo := range allWeaponCombos {
		if generator.weaponComboMatchesSettings(combo[0], combo[1]) {
			filteredCombos = append(filteredCombos, combo)
		}
	}

	return filteredCombos
}

func (generator *bulkSimCandidateGenerator) matchesWeaponTypeFilter(option *bulkSimCandidateOption, slot proto.ItemSlot) bool {
	filter := generator.weaponTypeFilters[slot]
	if len(filter) == 0 {
		return true
	}
	if option == nil {
		return false
	}
	return option.item.WeaponType > proto.WeaponType_WeaponTypeUnknown && slices.Contains(filter, option.item.WeaponType)
}

func (generator *bulkSimCandidateGenerator) weaponComboMatchesSettings(mhItem *bulkSimCandidateOption, ohItem *bulkSimCandidateOption) bool {
	frozenWeaponItem := generator.frozenWeaponItem
	if generator.frozenWeaponSlot == proto.ItemSlot_ItemSlotMainHand && frozenWeaponItem != nil && !candidateOptionEqualsItemPtr(mhItem, frozenWeaponItem) {
		return false
	}
	if generator.frozenWeaponSlot == proto.ItemSlot_ItemSlotOffHand && frozenWeaponItem != nil && !candidateOptionEqualsItemPtr(ohItem, frozenWeaponItem) {
		return false
	}
	return generator.matchesWeaponTypeFilter(mhItem, proto.ItemSlot_ItemSlotMainHand) && generator.matchesWeaponTypeFilter(ohItem, proto.ItemSlot_ItemSlotOffHand)
}
