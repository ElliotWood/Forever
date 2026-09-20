package warlock

// TODO: To be implemented. Menace is content Forever added; the sim has never
// modelled it, so this is a stub rather than a port.
//
// Spell IDs: 403828
// SpellData: none -- the client ships this unranked, so the generator emits no table
// SpellMask: WarlockSpellMenace
//
// Not wired into Initialize() yet -- wiring an unimplemented ability would panic
// on every character of this class. Wire it when it is implemented.
func (warlock *Warlock) registerMenace() {
	panic("To be implemented")
}
