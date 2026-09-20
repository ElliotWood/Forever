package mage

// TODO: To be implemented. spellData.FireBlast holds the seven trainer ranks, 2136 to 10199.
// The client also carries 400616 to 400623, the copies the Season of Discovery rune passive
// Overheat (400615) swaps onto the action bar. Overheat is an Engrave grant with no place in
// Forever, and the generator drops its stand-ins.
func (mage *Mage) registerFireBlastSpell() {
	// Registered unconditionally, so this returns instead of panicking -- a panic
	// here would stop the sim from starting at all rather than flagging one ability.
}
