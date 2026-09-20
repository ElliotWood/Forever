package warlock

// TODO: To be implemented. The Forever client ships no rank ladder the generator can
// read for this ability -- it survives as a single spell with no "Rank N" subtext and
// no ranked SkillLineAbility row -- so there is no data to build the spell from.
//
// Soul Siphon rides on this spell and nothing else, so stubbing Drain Life left
// Talents.SoulSiphon the one talent field in the repo that no code reads. Re-apply it
// here when the spell comes back: spellData.SoulSiphon holds the per-rank value, and
// upstream scaled Drain Life and Drain Soul by it per Affliction effect on the target.
func (warlock *Warlock) registerDrainLife() {
	// Registered unconditionally, so this returns instead of panicking -- a panic
	// here would stop the sim from starting at all rather than flagging one ability.
}
