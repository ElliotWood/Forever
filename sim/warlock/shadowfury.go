package warlock

// TODO: uncalled -- Forever drops the Shadowfury talent that gated this spell, so
// nothing calls this any more. Kept rather than deleted because "the talent is gone"
// and "the spell is gone" are not the same claim, and the client data does not
// distinguish them. Re-gate before wiring it back up.
// TODO: To be implemented. The Forever client ships no rank ladder the generator can
// read for this ability -- it survives as a single spell with no "Rank N" subtext and
// no ranked SkillLineAbility row -- so there is no data to build the spell from.
func (warlock *Warlock) registerShadowfury() {
	panic("To be implemented")
}
