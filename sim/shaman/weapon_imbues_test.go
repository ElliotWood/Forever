package shaman

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
)

// Rank 6's dummy is hundredths of damage per second of weapon speed (16344: 2810 at 60).
func TestFlametongueDamagePerSecond(t *testing.T) {
	if got := flametongueImbue.EffectN(1).Average(core.CharacterLevel) / 100; got != 28.1 {
		t.Fatalf("Flametongue rank 6 deals %v a second of weapon speed, want 28.1", got)
	}
}
