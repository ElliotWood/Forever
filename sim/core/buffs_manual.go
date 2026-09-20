// Package-local home of the buff drivers the generator calls by name. A
// generated apply block for a cooldown, proc, uptime, item-count or manual row,
// or for a row the manifest marks as driven, calls drive<Go>(unit, value);
// declaring that function here is what makes the row compile. A buff-scope row
// hands the driver the *Character the buffs are applied to, a debuff-scope row
// the *Unit the debuffs land on, and value is the proto field: a bool, a count
// or an uptime.
package core

// The party's Battle Shout is the external caster's copy, which chains behind
// the player's own shout rather than being up from the start.
func driveBattleShout(char *Character, _ bool) {
	ApplyFixedShoutAura(char, BattleShoutAura(&char.Unit, false, 0), BattleShoutCategory)
}

// The party's Commanding Shout chains the same way its sibling does.
func driveCommandingShout(char *Character, _ bool) {
	ApplyFixedShoutAura(char, CommandingShoutAura(&char.Unit, false, 0), CommandingShoutCategory)
}

// A stack of Sunder Armor is worth nothing until it is on the target, so the
// raid's copy is ramped to five over the first five global cooldowns, which is
// how long a warrior takes to stack it.
func driveSunderArmor(target *Unit, _ bool) {
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
