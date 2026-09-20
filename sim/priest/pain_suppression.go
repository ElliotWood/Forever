package priest

// TODO: To be implemented. Pain Suppression is content Forever added; the sim has never
// modelled it, so this is a stub rather than a port.
//
// Spell IDs: 402004
// SpellData: none -- the client ships this unranked, so the generator emits no table
// SpellMask: PriestSpellPainSuppression
//
// Not wired into Initialize() yet -- wiring an unimplemented ability would panic
// on every character of this class. Wire it when it is implemented.
func (priest *Priest) registerPainSuppression() {
	panic("To be implemented")
}
