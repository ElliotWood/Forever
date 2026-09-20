package paladin

// Spiritual Attunement (Rank 2, SpellID 33776): Whenever you are healed by another character's spell,
// you regain 10% of the amount healed as mana.
// In the sim, this is modeled as mana return from damage taken (since the healing model offsets damage).
//
// TODO: To be implemented. The Forever client ships no rank ladder the generator can
// read for this ability -- it survives as a single spell with no "Rank N" subtext and
// no ranked SkillLineAbility row -- so there is no data to build the spell from.
func (paladin *Paladin) RegisterSpiritualAttunement() {
	// Registered unconditionally, so this returns instead of panicking -- a panic
	// here would stop the sim from starting at all rather than flagging one ability.
}
