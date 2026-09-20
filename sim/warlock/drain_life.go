package warlock

// TODO: To be implemented. spellData.DrainLife holds the six trainer ranks, 689 to 11700.
// The client also carries 403677 to 403689, the copies the Season of Discovery rune passive
// Master Channeler (403668) swaps onto the action bar. That is an Engrave grant with no place
// in Forever, and the generator drops its stand-ins.
//
// Soul Siphon rides on this spell and nothing else, so stubbing Drain Life left
// Talents.SoulSiphon the one talent field in the repo that no code reads. Re-apply it
// here when the spell comes back: spellData.SoulSiphon holds the per-rank value, and
// upstream scaled Drain Life and Drain Soul by it per Affliction effect on the target.
func (warlock *Warlock) registerDrainLife() {
	// Registered unconditionally, so this returns instead of panicking -- a panic
	// here would stop the sim from starting at all rather than flagging one ability.
}
