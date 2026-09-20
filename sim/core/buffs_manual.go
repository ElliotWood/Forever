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
