package core

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

func newGeneratedDebuffTestTarget() *Unit {
	target := &Unit{
		Type:        EnemyUnit,
		Index:       0,
		Level:       63,
		auraTracker: newAuraTracker(),
		Env:         &Environment{MeasuringStats: true},
	}
	target.PseudoStats = stats.NewPseudoStats()
	return target
}

// The raid config registers its debuffs on every enemy and each is permanent,
// so the fight starts with all of them activated; the reset that does it in the
// sim needs a whole environment, which this stands in for.
func applyGeneratedTestDebuffs(target *Unit, debuffs *proto.Debuffs) *Simulation {
	sim := &Simulation{}
	applyDebuffEffects(target, 0, debuffs, &proto.Raid{})
	for _, aura := range target.auras {
		aura.Activate(sim)
	}
	return sim
}

func TestGeneratedCurseOfTheElementsHitsEverySchoolAndResistance(t *testing.T) {
	target := newGeneratedDebuffTestTarget()

	applyGeneratedTestDebuffs(target, &proto.Debuffs{CurseOfElements: true})

	for _, school := range generatedSchoolIndexes(126) {
		if got := target.PseudoStats.SchoolDamageTakenMultiplier[school]; got != 1.1 {
			t.Errorf("school %d takes %v times the damage, want the client's 1.1", school, got)
		}
	}
	if got := target.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical]; got != 1 {
		t.Errorf("physical damage taken is %v, want the client's mask to have left it alone", got)
	}

	for _, stat := range []stats.Stat{
		stats.ArcaneResistance, stats.FireResistance, stats.FrostResistance,
		stats.NatureResistance, stats.ShadowResistance,
	} {
		if got := target.stats[stat]; got != -75 {
			t.Errorf("%s is %v, want the client's -75", stat.StatName(), got)
		}
	}

	aura := target.GetAura("Curse of the Elements (External)")
	if aura == nil {
		t.Fatalf("no aura is labelled %q; the target has %v", "Curse of the Elements (External)", targetAuraLabels(target))
	}
	if want := (ActionID{SpellID: 1311680, Tag: -1}); aura.ActionID != want {
		t.Errorf("the curse is %v, want %v", aura.ActionID, want)
	}
	if aura.Duration != NeverExpires {
		t.Errorf("the curse lasts %v, want the raid config's copy to be up all fight", aura.Duration)
	}
	if got := CurseOfElementsDuration(0); got != 5*time.Minute {
		t.Errorf("the curse is %v long before it is made permanent, want the client's 5 minutes", got)
	}
}

// Expose Armor states its armor per combo point, and the raid config's debuff
// is the five-point finisher.
func TestGeneratedExposeArmorIsTheFivePointFinisher(t *testing.T) {
	target := newGeneratedDebuffTestTarget()

	applyGeneratedTestDebuffs(target, &proto.Debuffs{ExposeArmor: true})

	if got := target.stats[stats.Armor]; got != -2250 {
		t.Errorf("armor is %v, want five times the client's -450", got)
	}

	category := target.ExclusiveEffectManager.GetExclusiveEffectCategory(ExposeArmorCategory)
	if !category.SingleAura {
		t.Error("the armor category is not single-aura, so a second armor reduction could sit next to it")
	}
	if len(category.effects) != 1 || category.effects[0].Priority != 2250 {
		t.Errorf("the category holds %d effects, first bid %v; want one bidding the magnitude 2250",
			len(category.effects), category.effects[0].Priority)
	}
}

func targetAuraLabels(target *Unit) []string {
	labels := make([]string, 0, len(target.auras))
	for _, aura := range target.auras {
		labels = append(labels, aura.Label)
	}
	return labels
}
