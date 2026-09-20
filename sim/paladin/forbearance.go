package paladin

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Forbearance
// https://www.wowhead.com/forever/spell=25771
//
// Cannot be made invulnerable by Divine Shield, Divine Protection, Blessing of Protection or be affected by Avenging Wrath.
func (paladin *Paladin) registerForbearance() {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// paladin.Forbearance = paladin.RegisterAura(core.Aura{
	// 	Label:    "Forbearance",
	// 	ActionID: core.ActionID{SpellID: 25771},
	// 	Duration: time.Minute,
	// })
}
