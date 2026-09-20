package hunter

// TODO: To be implemented. Flanking Strike is content Forever added; the sim has never
// modelled it, so this is a stub rather than a port.
//
// Spell IDs: 415320
// SpellData: none -- the client ships this unranked, so the generator emits no table
// SpellMask: HunterSpellFlankingStrike
//
// Not wired into Initialize() yet -- wiring an unimplemented ability would panic
// on every character of this class. Wire it when it is implemented.
func (hunter *Hunter) registerFlankingStrike() {
	panic("To be implemented")
}
