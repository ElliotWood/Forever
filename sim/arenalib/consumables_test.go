package arenalib

import (
	"testing"

	"github.com/wowsims/classic/sim/core/proto"
)

// Two ways this goes wrong and neither would show up as a failure anywhere downstream - the
// arena would just publish a number for a shaman who had been quietly disarmed, or for a spec
// that inherited somebody else's weapon buff.
func TestClassImbuesOverrideTheRoleListWithoutLeaking(t *testing.T) {
	plain := consumesFor(Melee, ClassImbues{})
	if plain.Consumes.MainHandImbue != proto.WeaponImbue_Windfury {
		t.Fatalf("melee should start from the role list, got %v", plain.Consumes.MainHandImbue)
	}

	shaman := consumesFor(Melee, ClassImbues{
		MainHand: proto.WeaponImbue_WindfuryWeapon,
		OffHand:  proto.WeaponImbue_WindfuryWeapon,
	})
	if shaman.Consumes.MainHandImbue != proto.WeaponImbue_WindfuryWeapon ||
		shaman.Consumes.OffHandImbue != proto.WeaponImbue_WindfuryWeapon {
		t.Error("a class imbue did not override the role list")
	}
	// Everything else must survive the override.
	if shaman.Consumes.Flask != plain.Consumes.Flask || shaman.Consumes.Food != plain.Consumes.Food {
		t.Error("overriding an imbue dropped the rest of the list")
	}

	// A half override leaves the other hand alone.
	rogue := consumesFor(Melee, ClassImbues{OffHand: proto.WeaponImbue_InstantPoison})
	if rogue.Consumes.MainHandImbue != proto.WeaponImbue_Windfury {
		t.Error("overriding the off-hand changed the main hand")
	}

	// And none of it may have written through to the shared list.
	if again := consumesFor(Melee, ClassImbues{}); again.Consumes.MainHandImbue != proto.WeaponImbue_Windfury ||
		again.Consumes.OffHandImbue != proto.WeaponImbue_ElementalSharpeningStone {
		t.Error("a class imbue leaked into the shared role list")
	}

	// The ranged list exists to dodge the sharpening stone's ranged crit penalty.
	if consumesFor(Ranged, ClassImbues{}).Consumes.OffHandImbue != proto.WeaponImbue_WeaponImbueUnknown {
		t.Error("the hunter list must not carry the off-hand sharpening stone")
	}
}
