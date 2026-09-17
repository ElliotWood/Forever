package core

import "testing"

// A heal has no DefenseType, so CritMultiplier() used to panic on it. Ticking HoTs like
// Blood Craze (spellID 16488) crit under the Forever ruleset and took the sim down.
func TestCritMultiplierHelpfulWithoutDefenseType(t *testing.T) {
	at := &AttackTable{CritMultiplier: 1}
	heal := &Spell{Flags: SpellFlagHelpful, CritDamageBonus: 1}

	if got := heal.CritMultiplier(at); got != 1.5 {
		t.Fatalf("crit heal multiplier = %v, want 1.5", got)
	}

	defer func() {
		if recover() == nil {
			t.Fatal("damage spell with no DefenseType should still panic")
		}
	}()
	(&Spell{CritDamageBonus: 1}).CritMultiplier(at)
}
