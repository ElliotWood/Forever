package priest

// TODO: To be implemented. Power Word: Barrier is content Forever added; the sim has never
// modelled it, so this is a stub rather than a port.
//
// Spell IDs: 425207
// SpellData: none -- the client ships this unranked, so the generator emits no table
// SpellMask: PriestSpellPowerWordBarrier
//
// Not wired into Initialize() yet -- wiring an unimplemented ability would panic
// on every character of this class. Wire it when it is implemented.
func (priest *Priest) registerPowerWordBarrier() {
	panic("To be implemented")
}
