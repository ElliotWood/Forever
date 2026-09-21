package core

import (
	"testing"
)

const fakeThreatSpellClassMask int64 = 1 << 0

var fakeThreatSpellActionID = ActionID{SpellID: 11597}

// Sunder Armor deals no damage, so all of its threat is the flat bonus.
func (fw *FakeRageWarrior) registerFakeThreatSpell() {
	fw.RegisterSpell(SpellConfig{
		ActionID:        fakeThreatSpellActionID,
		ClassSpellMask:  fakeThreatSpellClassMask,
		FlatThreatBonus: 100,
	})
}

func TestFlatThreatBonusPctMod(t *testing.T) {
	sim := SetupFakeRageSim()
	fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)
	spell := fw.GetSpell(fakeThreatSpellActionID)
	attackTable := fw.AttackTables[sim.Encounter.ActiveTargetUnits[0].UnitIndex]

	threat := spell.ThreatFromDamage(sim, OutcomeHit, 0, attackTable)
	if !WithinToleranceFloat64(100, threat, 0.001) {
		t.Fatalf("Incorrect threat without a mod: Expected: %0.3f, Actual: %0.3f", 100.0, threat)
	}

	mod := fw.AddDynamicMod(SpellModConfig{
		ClassMask:  fakeThreatSpellClassMask,
		Kind:       SpellMod_FlatThreatBonus_Pct,
		FloatValue: 0.15,
	})

	mod.Activate()
	threat = spell.ThreatFromDamage(sim, OutcomeHit, 0, attackTable)
	if !WithinToleranceFloat64(115, threat, 0.001) {
		t.Fatalf("Incorrect threat with the mod active: Expected: %0.3f, Actual: %0.3f", 115.0, threat)
	}

	mod.Deactivate()
	threat = spell.ThreatFromDamage(sim, OutcomeHit, 0, attackTable)
	if !WithinToleranceFloat64(100, threat, 0.001) {
		t.Fatalf("Incorrect threat after the mod was removed: Expected: %0.3f, Actual: %0.3f", 100.0, threat)
	}
}
