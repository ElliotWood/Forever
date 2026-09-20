package warlock

// TODO: To be implemented. Demon Charge is content Forever added; the sim has never
// modelled it, so this is a stub rather than a port.
//
// Spell IDs: 412788
// SpellData: none -- the client ships this unranked, so the generator emits no table
// SpellMask: WarlockSpellDemonCharge
//
// Not wired into Initialize() yet -- wiring an unimplemented ability would panic
// on every character of this class. Wire it when it is implemented.
func (warlock *Warlock) registerDemonCharge() {
	panic("To be implemented")
}
