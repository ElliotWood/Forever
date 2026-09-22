package core

import (
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core/stats"
)

// A category's SingleAura flag is whatever its last member registered, so a
// generated buff has to declare the same thing the hand-written members of the
// category declare. Thunder Clap's AtkSpdReduction holds several auras at once;
// a paladin aura's category holds one.
func TestGeneratedAuraKeepsTheCategorySingleAuraFlag(t *testing.T) {
	target := newExclusiveTestTarget()

	newGeneratedDebuff(target, GeneratedBuff{
		Label:    "Generated Thunder Clap",
		ActionID: ActionID{SpellID: 11581},
		Duration: time.Second * 30,
		Category: "AtkSpdReduction",
		Pseudo: []PseudoConfig{
			{Kind: PseudoStatMeleeSpeedMultiplier, Amount: 0.8, IsMultiplicative: true},
		},
	})

	slow := target.ExclusiveEffectManager.GetExclusiveEffectCategory("AtkSpdReduction")
	if slow.SingleAura {
		t.Error("a debuff the manifest does not mark SingleAura turned its category single-aura, " +
			"which would stop the hand-written members from activating")
	}
	if len(slow.effects) != 1 {
		t.Errorf("registered %d effects under the bare category, want 1", len(slow.effects))
	}

	newGeneratedStatAura(target, GeneratedBuff{
		Label:      "Generated Devotion Aura",
		ActionID:   ActionID{SpellID: 10293},
		Duration:   NeverExpires,
		Category:   "DevotionAura",
		SingleAura: true,
		IsPlayer:   true,
		Stats:      []StatConfig{{stats.Armor, 735, false}},
	})

	devotion := target.ExclusiveEffectManager.GetExclusiveEffectCategory("DevotionAura")
	if !devotion.SingleAura {
		t.Error("a buff the manifest marks SingleAura did not turn its category single-aura")
	}
	if len(devotion.effects) != 1 {
		t.Errorf("registered %d effects under the bare category, want 1", len(devotion.effects))
	}
	if priority := devotion.effects[0].Priority; priority != 735 {
		t.Errorf("the category effect bids %v, want the positive magnitude 735", priority)
	}
}

// Every school of a damage-taken debuff has to be applied, which a category
// cannot do with one effect per school: it keeps a single effect per aura.
func TestGeneratedPseudoStatsCoverEverySchool(t *testing.T) {
	sim := &Simulation{}
	target := newExclusiveTestTarget()
	// The bare test unit leaves every multiplier at zero, and a damage-taken
	// multiplier is only meaningful against the 1 a real unit starts at.
	target.PseudoStats = stats.NewPseudoStats()

	aura := newGeneratedDebuff(target, GeneratedBuff{
		Label:    "Generated Curse of the Elements",
		ActionID: ActionID{SpellID: 1311680},
		Duration: time.Minute * 5,
		Category: "CurseOfElements",
		Pseudo: []PseudoConfig{{
			Kind: PseudoStatSchoolDamageTakenMultiplier, Amount: 1.1,
			IsMultiplicative: true, SchoolMask: 126,
		}},
	})

	aura.Activate(sim)

	for _, school := range generatedSchoolIndexes(126) {
		if got := target.PseudoStats.SchoolDamageTakenMultiplier[school]; got != 1.1 {
			t.Errorf("school %d takes %v times the damage, want 1.1", school, got)
		}
	}
	if got := target.PseudoStats.SchoolDamageTakenMultiplier[generatedSchoolIndexes(1)[0]]; got != 1 {
		t.Errorf("physical damage taken is %v, want the mask to have left it alone", got)
	}
}

// A resistance aura competes for the school with every other source of it, and
// separately for its own slot, so the two roles cannot be the same category.
func TestGeneratedResistanceAuraCompetesForTheSchool(t *testing.T) {
	sim := &Simulation{}
	target := newExclusiveTestTarget()
	// Enough of an environment for a stat change: the exclusive effects apply
	// theirs for real, which is what the competition is about.
	target.Env = &Environment{MeasuringStats: true}

	generated := MakePermanent(newGeneratedStatAura(target, GeneratedBuff{
		Label:      "Generated Frost Resistance Aura",
		ActionID:   ActionID{SpellID: 19898},
		Duration:   NeverExpires,
		Category:   FrostResistanceAuraCategory,
		SingleAura: true,
		IsPlayer:   true,
		Stats:      []StatConfig{{stats.FrostResistance, 60, false}},
	}))

	weaker := target.GetOrRegisterAura(Aura{
		Label:    "Frost Resistance Totem",
		ActionID: ActionID{SpellID: 10477},
		Duration: NeverExpires,
	})
	makeExclusiveFlatStatBuff(weaker, stats.FrostResistance, 30, ResistanceCategoryFrost)

	weaker.Activate(sim)
	if got := target.stats[stats.FrostResistance]; got != 30 {
		t.Fatalf("the weaker source applied %v frost resistance, want 30", got)
	}

	generated.Activate(sim)
	if got := target.stats[stats.FrostResistance]; got != 60 {
		t.Errorf("both sources applied %v frost resistance, want only the stronger 60", got)
	}

	school := target.ExclusiveEffectManager.GetExclusiveEffectCategory(
		ResistanceCategoryFrost + stats.FrostResistance.StatName() + "Add")
	if school.SingleAura {
		t.Error("the school category turned single-aura, which would push the other sources' auras off")
	}
	if len(school.effects) != 2 {
		t.Errorf("the school category holds %d effects, want the generated aura and the totem", len(school.effects))
	}

	own := target.ExclusiveEffectManager.GetExclusiveEffectCategory(FrostResistanceAuraCategory)
	if !own.SingleAura {
		t.Error("the aura's own category is not single-aura, so a second copy could sit next to it")
	}
	if !weaker.IsActive() {
		t.Error("the totem's aura was deactivated, which only the aura's own category may do")
	}
}

// A buff that grants several stats at once still competes per school for the
// resistances among them, which is how Gift of the Wild reads next to a
// dedicated resistance buff. Everything else it grants is applied outright.
func TestGeneratedBuffCompetesPerSchoolForItsResistances(t *testing.T) {
	sim := &Simulation{}
	target := newExclusiveTestTarget()
	target.Env = &Environment{MeasuringStats: true}

	giftOfTheWild := MakePermanent(newGeneratedStatAura(target, GeneratedBuff{
		Label:    "Generated Gift of the Wild",
		ActionID: ActionID{SpellID: 21850},
		Duration: time.Hour,
		Stats: []StatConfig{
			{stats.Armor, 385, false},
			{stats.Stamina, 16, false},
			{stats.ShadowResistance, 27, false},
			{stats.FireResistance, 27, false},
		},
	}))

	shadowProtection := MakePermanent(newGeneratedStatAura(target, GeneratedBuff{
		Label:    "Generated Prayer of Shadow Protection",
		ActionID: ActionID{SpellID: 27683},
		Duration: time.Minute * 20,
		Stats:    []StatConfig{{stats.ShadowResistance, 60, false}},
	}))

	giftOfTheWild.Activate(sim)
	shadowProtection.Activate(sim)

	if got := target.stats[stats.ShadowResistance]; got != 60 {
		t.Errorf("shadow resistance is %v, want only the stronger source's 60", got)
	}
	if got := target.stats[stats.FireResistance]; got != 27 {
		t.Errorf("fire resistance is %v, want the buff's own 27", got)
	}
	if got := target.stats[stats.Armor]; got != 385 {
		t.Errorf("armor is %v, want 385 applied outright", got)
	}
	if got := target.stats[stats.Stamina]; got != 16 {
		t.Errorf("stamina is %v, want 16 applied outright", got)
	}

	school := target.ExclusiveEffectManager.GetExclusiveEffectCategory(
		ResistanceCategoryShadow + stats.ShadowResistance.StatName() + "Add")
	if len(school.effects) != 2 {
		t.Errorf("the shadow school holds %d effects, want both sources", len(school.effects))
	}
	if !giftOfTheWild.IsActive() {
		t.Error("losing the school pushed the whole aura off, and it still grants everything else")
	}
}

// A_REDUCE_PUSHBACK states how much pushback is taken away, while the sim's
// PushbackChance is the chance of being pushed back and starts at 1.
func TestGeneratedConcentrationAuraReducesPushback(t *testing.T) {
	sim := &Simulation{}
	target := newExclusiveTestTarget()
	target.PseudoStats = stats.NewPseudoStats()

	aura := MakePermanent(newGeneratedStatAura(target, GeneratedBuff{
		Label:      "Generated Concentration Aura",
		ActionID:   ActionID{SpellID: 19746},
		Duration:   NeverExpires,
		Category:   "ConcentrationAura",
		SingleAura: true,
		Pseudo:     []PseudoConfig{{Kind: PseudoStatPushbackChance, Amount: -0.35}},
	}))

	aura.Activate(sim)

	if got := target.PseudoStats.PushbackChance; got != 0.65 {
		t.Errorf("pushback chance is %v, want 1 reduced by 35%% to 0.65", got)
	}
}

// The branch a buff with a category but no SingleAura takes: the pseudo-stats
// are registered by attachGeneratedPseudoStats rather than folded into one
// category-wide effect.
func TestGeneratedBuffPseudoStatsCoverEveryField(t *testing.T) {
	sim := &Simulation{}
	target := newExclusiveTestTarget()
	target.PseudoStats = stats.NewPseudoStats()

	aura := MakePermanent(newGeneratedStatAura(target, GeneratedBuff{
		Label:    "Generated School Shield",
		ActionID: ActionID{SpellID: 1311680},
		Duration: NeverExpires,
		Category: "GeneratedSchoolShield",
		Pseudo: []PseudoConfig{{
			Kind: PseudoStatSchoolDamageTakenMultiplier, Amount: 1.1,
			IsMultiplicative: true, SchoolMask: 126,
		}},
	}))

	aura.Activate(sim)

	for _, school := range generatedSchoolIndexes(126) {
		if got := target.PseudoStats.SchoolDamageTakenMultiplier[school]; got != 1.1 {
			t.Errorf("school %d takes %v times the damage, want 1.1", school, got)
		}
	}
	if got := target.PseudoStats.SchoolDamageTakenMultiplier[generatedSchoolIndexes(1)[0]]; got != 1 {
		t.Errorf("physical damage taken is %v, want the mask to have left it alone", got)
	}

	category := target.ExclusiveEffectManager.GetExclusiveEffectCategory(
		"GeneratedSchoolShield" + PseudoStatSchoolDamageTakenMultiplier.Name() + "Mul")
	if len(category.effects) != 1 {
		t.Errorf("the pseudo-stat category holds %d effects, want one covering every school", len(category.effects))
	}
}

// A buff the client states for a mask of schools raises what the unit deals in
// those schools only; mask 126 is every school but physical, so a melee swing
// is untouched while the six magic schools are multiplied.
func TestGeneratedBuffRaisesTheSchoolsItsMaskNames(t *testing.T) {
	sim := &Simulation{}
	unit := newExclusiveTestTarget()
	unit.PseudoStats = stats.NewPseudoStats()

	aura := MakePermanent(newGeneratedStatAura(unit, GeneratedBuff{
		Label:    "Generated Power Infusion",
		ActionID: ActionID{SpellID: 10060},
		Duration: NeverExpires,
		Pseudo: []PseudoConfig{{
			Kind: PseudoStatSchoolDamageDealtMultiplier, Amount: 1.2,
			IsMultiplicative: true, SchoolMask: 126,
		}},
	}))

	aura.Activate(sim)

	for _, school := range generatedSchoolIndexes(126) {
		if got := unit.PseudoStats.SchoolDamageDealtMultiplier[school]; got != 1.2 {
			t.Errorf("school %d deals %v times the damage, want 1.2", school, got)
		}
	}
	physical := generatedSchoolIndexes(1)[0]
	if got := unit.PseudoStats.SchoolDamageDealtMultiplier[physical]; got != 1 {
		t.Errorf("physical damage dealt is %v, want the mask to have left it alone", got)
	}
	if got := unit.PseudoStats.DamageDealtMultiplier; got != 1 {
		t.Errorf("the all-school multiplier is %v, want a masked buff to stay out of it", got)
	}

	aura.Deactivate(sim)
	for _, school := range generatedSchoolIndexes(127) {
		if got := unit.PseudoStats.SchoolDamageDealtMultiplier[school]; got != 1 {
			t.Errorf("school %d still deals %v times the damage once the buff is gone", school, got)
		}
	}
}

// Every member of the attack-speed category bids how far from 1 its multiplier
// is, so that the strongest slow on the target is the one that applies whatever
// form its source states it in. The generated one changes the melee speed
// through the unit's own method, because the swing timers are computed from it.
func TestGeneratedAttackSpeedDebuffBidsAgainstTheHandWrittenOnes(t *testing.T) {
	sim := &Simulation{}
	target := newExclusiveTestTarget()
	target.Env = &Environment{MeasuringStats: true}
	target.PseudoStats = stats.NewPseudoStats()

	clap := MakePermanent(newGeneratedDebuff(target, GeneratedBuff{
		Label:    "Generated Thunder Clap",
		ActionID: ActionID{SpellID: 11581},
		Duration: time.Second * 30,
		Category: "AtkSpdReduction",
		Pseudo: []PseudoConfig{
			{Kind: PseudoStatMeleeSpeedMultiplier, Amount: 0.8, IsMultiplicative: true},
		},
	}))

	clap.Activate(sim)
	if got := target.PseudoStats.MeleeSpeedMultiplier; got != 0.8 {
		t.Fatalf("the melee speed multiplier is %v, want the client's 0.8", got)
	}
	if got := target.TotalMeleeHasteMultiplier(); got != 0.8 {
		t.Errorf("the target swings at %v times its speed, want the slow to have reached the swing timers", got)
	}

	// A hand-written slow states the multiplier the target's speed is divided
	// by, so 1.5 is a third slower and outbids the clap's fifth.
	slower := target.GetOrRegisterAura(Aura{
		Label:    "Hand-written Slow",
		ActionID: ActionID{SpellID: 27648},
		Duration: time.Second * 12,
	})
	// Through a variable, so that the test does the same float64 arithmetic the
	// helper does rather than Go's exact constant arithmetic.
	divisor := 1.5
	stronger := AtkSpeedReductionEffect(slower, divisor)
	if want := 1 - 1/divisor; stronger.Priority != want {
		t.Errorf("a slow that divides by %v bids %v, want the magnitude %v", divisor, stronger.Priority, want)
	}

	slower.Activate(sim)

	if got := target.PseudoStats.MeleeSpeedMultiplier; got != 1 {
		t.Errorf("the melee speed multiplier is %v, want the weaker slow taken back", got)
	}
	if got, want := target.PseudoStats.AttackSpeedMultiplier, 1/divisor; got != want {
		t.Errorf("the attack speed multiplier is %v, want the stronger slow's %v", got, want)
	}
	if !clap.IsActive() {
		t.Error("losing the category pushed the weaker aura off, which only a single-aura category may do")
	}
}

// Thunderfury's Cyclone states 1.2, which divides the target's speed by 1.2 and
// is a 16.67% slow, so the clap's 20% has to keep the category when the proc
// lands. The two forms are what made the scales easy to confuse: 1.2 - 1 is
// bit-for-bit the clap's 1 - 0.8.
func TestGeneratedThunderClapOutbidsThunderfurysCyclone(t *testing.T) {
	sim := &Simulation{}
	target := newExclusiveTestTarget()
	target.Env = &Environment{MeasuringStats: true}
	target.PseudoStats = stats.NewPseudoStats()

	clap := MakePermanent(newGeneratedDebuff(target, GeneratedBuff{
		Label:    "Generated Thunder Clap",
		ActionID: ActionID{SpellID: 11581},
		Duration: time.Second * 30,
		Category: "AtkSpdReduction",
		Pseudo: []PseudoConfig{
			{Kind: PseudoStatMeleeSpeedMultiplier, Amount: 0.8, IsMultiplicative: true},
		},
	}))
	clap.Activate(sim)

	cyclone := target.GetOrRegisterAura(Aura{
		Label:    "Cyclone",
		ActionID: ActionID{SpellID: 27648},
		Duration: time.Second * 12,
	})
	weaker := AtkSpeedReductionEffect(cyclone, 1.2)
	if clapBid := clap.ExclusiveEffects[0].Priority; weaker.Priority >= clapBid {
		t.Fatalf("the proc bids %v against the clap's %v, want the 16.67%% slow to be worth less",
			weaker.Priority, clapBid)
	}

	cyclone.Activate(sim)

	if got := target.PseudoStats.MeleeSpeedMultiplier; got != 0.8 {
		t.Errorf("the melee speed multiplier is %v, want the clap's 0.8 to have stayed", got)
	}
	if got := target.PseudoStats.AttackSpeedMultiplier; got != 1 {
		t.Errorf("the attack speed multiplier is %v, want the weaker proc to have applied nothing", got)
	}
	if got := target.TotalMeleeHasteMultiplier(); got != 0.8 {
		t.Errorf("the target swings at %v times its speed, want the stronger slow's 0.8", got)
	}
}

// A debuff bids once for everything it applies, its resistances included: the
// school-by-school competition is what keeps two sources of a resistance from
// adding up on a player, while two resistance-reducing curses exclude each
// other in their own category instead.
func TestGeneratedDebuffKeepsItsResistancesInItsOwnCategory(t *testing.T) {
	sim := &Simulation{}
	target := newExclusiveTestTarget()
	target.Env = &Environment{MeasuringStats: true}
	target.PseudoStats = stats.NewPseudoStats()

	curse := MakePermanent(newGeneratedDebuff(target, GeneratedBuff{
		Label:      "Generated Curse of the Elements",
		ActionID:   ActionID{SpellID: 1311680},
		Duration:   time.Minute * 5,
		Category:   "CurseOfElements",
		SingleAura: true,
		Stats: []StatConfig{
			{stats.FireResistance, -75, false},
			{stats.ShadowResistance, -75, false},
		},
		Pseudo: []PseudoConfig{{
			Kind: PseudoStatSchoolDamageTakenMultiplier, Amount: 1.1,
			IsMultiplicative: true, SchoolMask: 126,
		}},
	}))

	stronger := target.GetOrRegisterAura(Aura{
		Label:    "Curse of Shadow",
		ActionID: ActionID{SpellID: 11722},
		Duration: time.Minute * 5,
	})
	stronger.NewExclusiveEffect("CurseOfElements", true, ExclusiveEffect{
		Priority: 100,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stats.ShadowResistance, -100)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stats.ShadowResistance, 100)
		},
	})

	curse.Activate(sim)
	if got := target.stats[stats.ShadowResistance]; got != -75 {
		t.Fatalf("the curse alone reduced shadow resistance by %v, want the client's -75", got)
	}

	stronger.Activate(sim)

	if got := target.stats[stats.ShadowResistance]; got != -100 {
		t.Errorf("shadow resistance is %v, want only the stronger curse's -100", got)
	}
	if got := target.stats[stats.FireResistance]; got != 0 {
		t.Errorf("fire resistance is %v, want the weaker curse to have taken its -75 back", got)
	}
	if curse.IsActive() {
		t.Error("both curses are on the target, so two resistance reductions apply at once")
	}

	school := target.ExclusiveEffectManager.GetExclusiveEffectCategory(
		ResistanceCategoryShadow + stats.ShadowResistance.StatName() + "Add")
	if len(school.effects) != 0 {
		t.Errorf("the curse put %d effects in the shadow school, where nothing can outbid them",
			len(school.effects))
	}
}

// A flat bonus the client's buff data does not state raises two numbers that a
// non-stacking buff keeps apart: what the aura applies and what it bids for its
// category. A stacking buff prices the one off the other, so it refuses one.
func TestAddGeneratedFlatBonusRaisesTheBidAndTheAmount(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	aura := newGeneratedStatAura(&char.Unit, GeneratedBuff{
		Label:      "Generated Battle Shout",
		ActionID:   ActionID{SpellID: 25289}.WithTag(-1),
		Duration:   NeverExpires,
		Category:   "GeneratedBattleShout",
		SingleAura: true,
		Stats:      []StatConfig{{stats.AttackPower, 139, false}},
	})

	AddGeneratedFlatBonus(aura, stats.AttackPower, 139, 30)
	// A second party member wearing the same set asks for the same total.
	AddGeneratedFlatBonus(aura, stats.AttackPower, 139, 30)

	if priority := aura.ExclusiveEffects[0].Priority; priority != 169 {
		t.Errorf("the category effect bids %v, want the client's 139 plus the set's 30", priority)
	}

	MakePermanent(aura)
	char.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)

	if got := char.stats[stats.AttackPower]; got != 169 {
		t.Errorf("the aura applied %v attack power, want the client's 139 plus the set's 30", got)
	}
}

func TestAddGeneratedFlatBonusRefusesAnAuraThatCannotCarryOne(t *testing.T) {
	char := newGeneratedBuffTestCharacter()

	stacking := newGeneratedDebuff(&char.Unit, GeneratedBuff{
		Label:     "Generated Sunder Armor",
		ActionID:  ActionID{SpellID: 25225},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		Category:  "GeneratedArmorReduction",
		Stats:     []StatConfig{{stats.Armor, -520, false}},
	})
	assertPanics(t, "a stacking aura", func() {
		AddGeneratedFlatBonus(stacking, stats.Armor, -2600, -30)
	})

	uncontested := newGeneratedStatAura(&char.Unit, GeneratedBuff{
		Label:    "Generated Blood Pact",
		ActionID: ActionID{SpellID: 11767}.WithTag(-1),
		Duration: NeverExpires,
		Stats:    []StatConfig{{stats.Stamina, 54, false}},
	})
	assertPanics(t, "an aura with no category", func() {
		AddGeneratedFlatBonus(uncontested, stats.Stamina, 54, 30)
	})

	contested := newGeneratedStatAura(&char.Unit, GeneratedBuff{
		Label:      "Generated Devotion Aura",
		ActionID:   ActionID{SpellID: 10293}.WithTag(-1),
		Duration:   NeverExpires,
		Category:   "GeneratedDevotionAura",
		SingleAura: true,
		Stats:      []StatConfig{{stats.Armor, 735, false}},
	})
	assertPanics(t, "a base that is not what the buff bids", func() {
		AddGeneratedFlatBonus(contested, stats.Armor, 620, 30)
	})

	// A category that holds several auras at once never deactivates the one it
	// outbids, so an amount attached to the aura would apply while the effect
	// holding it did not. No generated row has that shape; the guard is for a
	// hand-written one that does.
	shared := char.GetOrRegisterAura(Aura{
		Label:    "Hand-Written Shout",
		Tag:      "GeneratedSharedCategory",
		ActionID: ActionID{SpellID: 25289}.WithTag(-2),
		Duration: NeverExpires,
	})
	shared.NewExclusiveEffect("GeneratedSharedCategory", false, ExclusiveEffect{Priority: 139})
	assertPanics(t, "a category that holds more than one aura", func() {
		AddGeneratedFlatBonus(shared, stats.AttackPower, 139, 30)
	})
}

func assertPanics(t *testing.T, what string, call func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s took a flat bonus, which it cannot carry", what)
		}
	}()
	call()
}
