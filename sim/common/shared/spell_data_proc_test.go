package shared

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
)

// Eternal Power, the Dormant Heart of the Mountain proc: 1249118 carries the proc flags and 1249119
// is the buff it applies. The item IDs below are the test's own, so that registering them cannot
// collide with a real effect.
const (
	healProcTrigger int32 = 1249118
	healProcBuff    int32 = 1249119
)

func TestSpellDataProcRegistersEveryVariant(t *testing.T) {
	NewSpellDataProc(SpellDataProc{TriggerSpellID: healProcTrigger, BuffSpellID: healProcBuff},
		[]ItemVariant{
			{ItemID: 990001, ItemName: "Reissued Test Trinket"},
			{ItemID: 990002, ItemName: "Test Trinket"},
		})

	if !core.HasItemEffect(990001) || !core.HasItemEffect(990002) {
		t.Error("a variant was left unregistered")
	}

	// Only the highest ID goes into the test suite, so a dozen re-issues of one trinket do not each
	// get their own fixture entry.
	if core.HasItemEffectForTest(990001) {
		t.Error("the lower-ID variant was added to the test suite as well")
	}
	if !core.HasItemEffectForTest(990002) {
		t.Error("the highest-ID variant was kept out of the test suite")
	}
	if !core.AddEffectsToTest {
		t.Error("the variant loop left effects out of the test suite for whatever registers next")
	}
}

// A row that names no hits the sim hears would register a listener that can never fire. 1249119 is
// the buff half of the pair above: it carries the stats and no proc flags at all.
func TestSpellDataProcSkipsARowWithNoListener(t *testing.T) {
	NewSpellDataProc(SpellDataProc{TriggerSpellID: healProcBuff},
		[]ItemVariant{{ItemID: 990003, ItemName: "Listenerless Test Trinket"}})

	if core.HasItemEffect(990003) {
		t.Error("an effect whose row states no listener was registered anyway")
	}
}

// An enchant states its own name and registers through the enchant registry instead.
func TestSpellDataProcRegistersAnEnchant(t *testing.T) {
	NewSpellDataProc(SpellDataProc{
		Name:           "Test Enchant",
		EnchantID:      990004,
		TriggerSpellID: healProcTrigger,
		BuffSpellID:    healProcBuff,
	}, nil)

	if !core.HasEnchantEffect(990004) {
		t.Error("an enchant proc with no variants was left unregistered")
	}
	if core.HasItemEffect(990004) {
		t.Error("an enchant proc was registered as an item effect")
	}
}
