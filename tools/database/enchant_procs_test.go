package database

import (
	"slices"
	"strings"
	"testing"

	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/tools/database/dbc"
)

// Which slots of a real enchant reach the generator, in which shape and at which chance. Every
// entry names the reason it is refused, or none where it registers.
func TestEnchantProcRouting(t *testing.T) {
	inRepositoryRoot(t)
	instance := dbc.GetDBC()
	grants := enchantGrantEffects(instance.SpellEffectsById)

	type want struct {
		trigger   int
		damage    bool
		chancePct int
		reason    string
	}

	for _, tc := range []struct {
		effectID int
		name     string
		want     []want
	}{
		{36, "Fiery Blaze: Effect 1 states 15, the tooltip's 15% chance", []want{{6297, true, 15, ""}}},
		{803, "Fiery Weapon: Effect 1 states 0", []want{{13897, true, 0, spelldata.ReasonStatesNoRate}}},
		{1899, "Unholy Weapon: 'often' is a rate the client does not carry", []want{{20006, false, 0, spelldata.ReasonStatesNoRate}}},
		{1900, "Crusader: the combat spell only, the aura 458112 applies the same Holy Strength", []want{{20007, false, 0, spelldata.ReasonStatesNoRate}}},
		{7940, "Grand Crusader: the same shape as Crusader", []want{{1231124, false, 0, spelldata.ReasonStatesNoRate}}},
		{7210, "Dismantle: two auras, 100 beside 'sometimes' and 30%, both hitting mechanicals only", []want{
			{435467, true, 0, spelldata.ReasonStatesNoRate},
			{442206, true, 0, "mechanicals"},
		}},
		{7223, "Retricutioner: a damage shield", []want{{435901, false, 0, "A_DAMAGE_SHIELD"}}},
		{7941, "Grand Arcanist: spell power and healing register, the mana beside them is not a stat", []want{{1231152, false, 0, ""}}},
		{8217, "Revelation: 100 beside 'a chance to trigger'", []want{{1248806, false, 0, spelldata.ReasonStatesNoRate}}},
		{8721, "Recovery: 100 beside a cooldown, which is a rate; the heal is not a buff", []want{{1248761, false, 0, "E_HEAL_PCT"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			enchant, ok := instance.EnchantsByEffectId[tc.effectID]
			if !ok {
				t.Fatalf("enchant %d is not in the enchant inputs", tc.effectID)
			}
			grant, ok := grants[tc.effectID]
			if !ok {
				t.Fatalf("no spell grants enchant %d", tc.effectID)
			}

			got := routeEnchantProcs(enchant.ProcSlots(), instance, renderSpellTooltip(instance, grant.SpellID))
			if len(got) != len(tc.want) {
				t.Fatalf("%d routings, want %d", len(got), len(tc.want))
			}

			for i, w := range tc.want {
				r := got[i]
				if r.TriggerSpellID != w.trigger || r.Damage != w.damage || r.ProcChancePct != w.chancePct {
					t.Errorf("routing %d: trigger %d, damage %v, chance %d%%; want %d, %v, %d%%",
						i, r.TriggerSpellID, r.Damage, r.ProcChancePct, w.trigger, w.damage, w.chancePct)
				}
				if w.reason == "" && !r.Supported() {
					t.Errorf("routing %d is refused (%s), want it registered", i, r.Reason())
				}
				if w.reason != "" && !strings.Contains(r.Reason(), w.reason) {
					t.Errorf("routing %d: reason %q, want it to name %q", i, r.Reason(), w.reason)
				}
			}
		})
	}
}

// A PPM override answers "states no rate" on the spell the slot is routed through - the combat spell
// of an Effect 1 slot, the equip aura of an Effect 3 one - and leaves every other refusal in place.
func TestEnchantProcRoutingTakesAPPMOverride(t *testing.T) {
	inRepositoryRoot(t)
	instance := dbc.GetDBC()
	grants := enchantGrantEffects(instance.SpellEffectsById)

	for _, tc := range []struct {
		effectID   int
		name       string
		overrideOn int32
	}{
		{803, "Fiery Weapon: Effect 1, a damage combat spell", 13897},
		{1899, "Unholy Weapon: Effect 1, a combat spell granting a buff", 20006},
		{8217, "Revelation: Effect 3, the aura's 100 beside 'a chance to trigger'", 1248806},
	} {
		t.Run(tc.name, func(t *testing.T) {
			enchant := instance.EnchantsByEffectId[tc.effectID]
			route := func() *ProcRouting {
				t.Helper()
				got := routeEnchantProcs(enchant.ProcSlots(), instance,
					renderSpellTooltip(instance, grants[tc.effectID].SpellID))
				if len(got) != 1 || got[0].TriggerSpellID != int(tc.overrideOn) {
					t.Fatalf("routings %v, want one through %d", got, tc.overrideOn)
				}
				return got[0]
			}

			before := route()
			if !statesNoRate(before) {
				t.Fatalf("reason %q names no missing rate to answer", before.Reason())
			}

			withPPMOverride(t, tc.overrideOn, 2)

			after := route()
			want := slices.DeleteFunc(slices.Clone(before.Unsupported), func(reason string) bool {
				return reason == spelldata.ReasonStatesNoRate
			})
			if !slices.Equal(after.Unsupported, want) {
				t.Errorf("unsupported = %v with the override, want %v", after.Unsupported, want)
			}
		})
	}
}

// The store row as a PPM override writes it, restored when the test ends.
func withPPMOverride(t *testing.T, spellID int32, ppm float32) {
	t.Helper()
	row := spelldata.Find(spellID)
	if row == spelldata.Nil {
		t.Fatalf("spell %d is not in the store", spellID)
	}

	original := *row
	overridden := original
	overridden.RPPM = ppm
	overridden.ProcChanceSource, overridden.ProcChanceEffect = spelldata.ProcChancePPM, 0
	*row = overridden
	t.Cleanup(func() { *row = original })
}

// Recovery's tooltip states a cooldown with "more often than", which is not the "often" of a rate.
func TestEnchantTooltipRateWording(t *testing.T) {
	for tooltip, want := range map[string]bool{
		"Permanently enchant a melee weapon to often strike for 40 additional fire damage.":                   true,
		"Permanently enchant a Weapon to cause all spells and attacks to sometimes deal 76 additional damage": true,
		"Permanently enchant a Melee Weapon to have a chance to trigger Revelation when a non-periodic spell": true,
		"healing you for 5% of your maximum health. Cannot occur more often than once every 10 sec.":          false,
	} {
		if got := enchantTooltipStatesAnUnknownRate(tooltip); got != want {
			t.Errorf("%q: %v, want %v", tooltip, got, want)
		}
	}
}
