package priest

// TODO: To be implemented. Increased Fortitude is content Forever added; the sim has never
// modelled it, so this is a stub rather than a port.
//
// Spell IDs: 436951
// SpellData: none -- the client ships this unranked, so the generator emits no table
// SpellMask: PriestSpellIncreasedFortitude
//
// Not wired into Initialize() yet -- wiring an unimplemented ability would panic
// on every character of this class. Wire it when it is implemented.
func (priest *Priest) registerIncreasedFortitude() {
	panic("To be implemented")
}
