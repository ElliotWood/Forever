package spelldata

// The SpellMisc.Attributes bits the accessors below read, named after the word they live in the way
// tools/database/dbc/enums.go names them: ATTR_EX_1 is a bit of Attributes_1.
const (
	// The spell is always on rather than cast: a talent's passive aura, a set bonus, an item effect.
	ATTR_PASSIVE uint32 = 0x40

	// The two bits the client marks a channel with. Arcane Missiles and Blizzard carry the first,
	// Evocation and Tranquility only the second, so a channel check has to read both.
	ATTR_EX_1_IS_CHANNELLED      uint32 = 0x4
	ATTR_EX_1_IS_SELF_CHANNELLED uint32 = 0x40

	// The server refunds 80% of the power cost when the spell misses: every rage special that pays up
	// front, Heroic Strike and Rend among them. Cleave and Whirlwind lack it.
	ATTR_EX_1_DISCOUNT_POWER_ON_MISS uint32 = 0x8000000

	// The listener aura also fires from hits of triggered spells.
	ATTR_EX_3_CAN_PROC_FROM_PROCS uint32 = 0x4000000

	// Weapon procs ignore hits of this spell: the Seal of Blood, Righteousness and Martyr damage
	// spells, plus Gouge, Sap, Scatter Shot and Maim.
	ATTR_EX_4_SUPPRESS_WEAPON_PROCS uint32 = 0x800000

	// The ticks of this spell's periodic effect roll a critical strike. Rend and Corruption carry it,
	// Deep Wounds does not.
	ATTR_EX_8_PERIODIC_CAN_CRIT uint32 = 0x200

	// The listener aura only fires from class abilities rather than from any hit.
	ATTR_EX_12_ONLY_PROC_FROM_CLASS_ABILITIES uint32 = 0x80000000
)

// The Attributes word each bit above belongs to.
const (
	ATTR_INDEX_0     = 0
	ATTR_INDEX_EX_1  = 1
	ATTR_INDEX_EX_3  = 3
	ATTR_INDEX_EX_4  = 4
	ATTR_INDEX_EX_8  = 8
	ATTR_INDEX_EX_12 = 12
)

// A word past the end of the client's 17 reads as unset rather than panicking, so an attribute the
// store does not carry never takes a caller down.
func (s *Spell) HasAttr(word int, bit uint32) bool {
	if word < 0 || word >= len(s.Attr) {
		return false
	}
	return s.Attr[word]&bit != 0
}

func (s *Spell) IsPassive() bool {
	return s.HasAttr(ATTR_INDEX_0, ATTR_PASSIVE)
}

func (s *Spell) IsChanneled() bool {
	return s.HasAttr(ATTR_INDEX_EX_1, ATTR_EX_1_IS_CHANNELLED|ATTR_EX_1_IS_SELF_CHANNELLED)
}

func (s *Spell) RefundsOnMiss() bool {
	return s.HasAttr(ATTR_INDEX_EX_1, ATTR_EX_1_DISCOUNT_POWER_ON_MISS)
}

func (s *Spell) PeriodicCanCrit() bool {
	return s.HasAttr(ATTR_INDEX_EX_8, ATTR_EX_8_PERIODIC_CAN_CRIT)
}

func (s *Spell) CanProcFromProcs() bool {
	return s.HasAttr(ATTR_INDEX_EX_3, ATTR_EX_3_CAN_PROC_FROM_PROCS)
}

func (s *Spell) ClassSpellsOnly() bool {
	return s.HasAttr(ATTR_INDEX_EX_12, ATTR_EX_12_ONLY_PROC_FROM_CLASS_ABILITIES)
}

func (s *Spell) SuppressesWeaponProcs() bool {
	return s.HasAttr(ATTR_INDEX_EX_4, ATTR_EX_4_SUPPRESS_WEAPON_PROCS)
}

// The share of the cost a miss refunds, for RageCostOptions.Refund: 80% where the client flags
// Discount Power On Miss, nothing otherwise.
func (s *Spell) MissRefund() float64 {
	if s.RefundsOnMiss() {
		return 0.8
	}
	return 0
}
