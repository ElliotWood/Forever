package forever

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

func RegisterAllEnchants() {

	// Enchants

	// Permanently enchant a melee weapon to often inflict a curse on the target reducing their melee damage.
	// https://www.wowhead.com/forever/spell=27102
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:               "Enchant Weapon - Unholy Weapon",
		EnchantID:          1899,
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskUnknown,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
		IsWeaponProc:       true,
	})
}
