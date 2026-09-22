package dbcenums

// SpellPower.PowerType, under TrinityCore's names. The store's rows carry these six.
type PowerType int8

const (
	POWER_HEALTH       PowerType = -2
	POWER_MANA         PowerType = 0
	POWER_RAGE         PowerType = 1
	POWER_FOCUS        PowerType = 2
	POWER_ENERGY       PowerType = 3
	POWER_COMBO_POINTS PowerType = 4
)
