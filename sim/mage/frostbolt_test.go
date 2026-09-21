package mage

import "testing"

// Our Frostbolt damage is kept beside the client table (see frostbolt.go); it must not drift from it.
func TestFrostboltDamageContainsClientValue(t *testing.T) {
	for rank := 1; rank <= FrostboltRanks; rank++ {
		low, high := spellData.Frostbolt.ByRank(int32(rank)).Direct.Range()
		if r := FrostboltBaseDamage[rank]; low < r[0] || high > r[1] {
			t.Errorf("rank %d: table %v-%v outside our %v-%v", rank, low, high, r[0], r[1])
		}
	}
}
