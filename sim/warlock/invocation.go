package warlock

// TODO: To be implemented. Invocation is content Forever added; the sim has never
// modelled it, so this is a stub rather than a port.
//
// Spell IDs: 426241, 426245, 426246, 426247, 426331
// SpellData: none -- the client ships this unranked, so the generator emits no table
// SpellMask: WarlockSpellInvocation
//
// Not wired into Initialize() yet -- wiring an unimplemented ability would panic
// on every character of this class. Wire it when it is implemented.
func (warlock *Warlock) registerInvocation() {
	panic("To be implemented")
}
