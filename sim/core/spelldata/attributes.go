package spelldata

import "github.com/wowsims/forever/sim/core/dbcenums"

// A word past the end of the client's 17 reads as unset rather than panicking, so an attribute the
// store does not carry never takes a caller down.
func (s *Spell) HasAttr(word int, bit uint32) bool {
	if word < 0 || word >= len(s.Attr) {
		return false
	}
	return s.Attr[word]&bit != 0
}

func (s *Spell) IsPassive() bool {
	return s.HasAttr(dbcenums.ATTR_INDEX_BASE, dbcenums.ATTR_PASSIVE)
}

func (s *Spell) IsChanneled() bool {
	return s.HasAttr(dbcenums.ATTR_INDEX_EX_1, dbcenums.ATTR_EX_1_IS_CHANNELLED|dbcenums.ATTR_EX_1_IS_SELF_CHANNELLED)
}

func (s *Spell) RefundsOnMiss() bool {
	return s.HasAttr(dbcenums.ATTR_INDEX_EX_1, dbcenums.ATTR_EX_1_DISCOUNT_POWER_ON_MISS)
}

func (s *Spell) PeriodicCanCrit() bool {
	return s.HasAttr(dbcenums.ATTR_INDEX_EX_8, dbcenums.ATTR_EX_8_PERIODIC_CAN_CRIT)
}

func (s *Spell) CanProcFromProcs() bool {
	return s.HasAttr(dbcenums.ATTR_INDEX_EX_3, dbcenums.ATTR_EX_3_CAN_PROC_FROM_PROCS)
}

func (s *Spell) ClassSpellsOnly() bool {
	return s.HasAttr(dbcenums.ATTR_INDEX_EX_12, dbcenums.ATTR_EX_12_ONLY_PROC_FROM_CLASS_ABILITIES)
}

func (s *Spell) SuppressesWeaponProcs() bool {
	return s.HasAttr(dbcenums.ATTR_INDEX_EX_4, dbcenums.ATTR_EX_4_SUPPRESS_WEAPON_PROCS)
}

// The share of the cost a miss refunds, for RageCostOptions.Refund: 80% where the client flags
// Discount Power On Miss, nothing otherwise.
func (s *Spell) MissRefund() float64 {
	if s.RefundsOnMiss() {
		return 0.8
	}
	return 0
}
