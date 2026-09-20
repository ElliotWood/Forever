// Package-local home of the buff drivers the generator calls by name. A
// generated apply block for a cooldown, proc, uptime, item-count or manual row
// calls drive<Go>(unit, value); declaring that function here is what makes the
// row compile. A buff-scope row hands the driver the *Character the buffs are
// applied to, a debuff-scope row the *Unit the debuffs land on, and value is
// the proto field: a bool, a count or an uptime.
package core
