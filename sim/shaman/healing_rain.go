package shaman

// TODO: To be implemented. Healing Rain is content Forever added; the sim has never
// modelled it, so this is a stub rather than a port.
//
// Spell IDs: 415236, 415242
// SpellData: none -- the client ships this unranked, so the generator emits no table
// SpellMask: SpellMaskHealingRain
//
// Not wired into Initialize() yet -- wiring an unimplemented ability would panic
// on every character of this class. Wire it when it is implemented.
func (shaman *Shaman) registerHealingRain() {
	panic("To be implemented")
}
