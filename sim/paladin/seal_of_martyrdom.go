package paladin

// TODO: To be implemented. Seal of Martyrdom is content Forever added; the sim has never
// modelled it, so this is a stub rather than a port.
//
// Spell IDs: 407798, 407799
// SpellData: none -- the client ships this unranked, so the generator emits no table
// SpellMask: NONE -- paladin's mask block is full. SpellMaskNone takes the iota 0 slot
// without using bit 0, so the block spans bits 1-62 and every one is allocated. Giving
// this ability a mask needs bit 0 reclaimed or a second mask word.
//
// Not wired into Initialize() yet -- wiring an unimplemented ability would panic
// on every character of this class. Wire it when it is implemented.
func (paladin *Paladin) registerSealOfMartyrdom() {
	panic("To be implemented")
}
