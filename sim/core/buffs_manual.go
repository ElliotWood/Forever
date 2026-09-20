// Package-local home of the buff drivers the generator calls by name. A
// generated apply block for a cooldown, proc, uptime, item-count or manual row,
// or for a row the manifest marks as driven, calls drive<Go>(unit, scope);
// declaring that function here is what makes the row compile. A buff-scope row
// hands the driver the *Character the buffs are applied to and the message its
// field lives on, a debuff-scope row the *Unit the debuffs land on plus the
// debuffs and the raid. The driver reads its own field out of the message, and
// any other field it needs: Grace of Air lasts 9 seconds while the party is
// twisting totems.
package core

import "github.com/wowsims/forever/sim/core/proto"

// The party's Battle Shout is the external caster's copy, which chains behind
// the player's own shout rather than being up from the start.
func driveBattleShout(char *Character, _ *proto.PartyBuffs) {
	ApplyFixedShoutAura(char, BattleShoutAura(&char.Unit, false, 0), BattleShoutCategory)
}

// The party's Commanding Shout chains the same way its sibling does.
func driveCommandingShout(char *Character, _ *proto.PartyBuffs) {
	ApplyFixedShoutAura(char, CommandingShoutAura(&char.Unit, false, 0), CommandingShoutCategory)
}

// A stack of Sunder Armor is worth nothing until it is on the target, so the
// raid's copy is ramped to five over the first five global cooldowns, which is
// how long a warrior takes to stack it.
func driveSunderArmor(target *Unit, _ *proto.Debuffs, _ *proto.Raid) {
	aura := MakePermanent(SunderArmorAura(target, false, 0))

	ScheduledAura(aura, PeriodicActionOptions{
		Period:          GCDDefault,
		NumTicks:        5,
		TickImmediately: true,
		Priority:        ActionPriorityDOT,
		OnAction: func(sim *Simulation) {
			aura.Activate(sim)
			if aura.IsActive() {
				aura.AddStack(sim)
			}
		},
	})
}
