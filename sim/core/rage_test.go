package core

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/simsignals"
)

func init() {
	RegisterAgentFactory(
		proto.Player_DpsWarrior{},
		proto.Spec_SpecDpsWarrior,
		NewFakeRageWarrior,
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DpsWarrior)
			if !ok {
				panic("Invalid spec value for Dps Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

const (
	fakeMHSwingSpeed = 2.6
	fakeOHSwingSpeed = 1.8
)

type FakeRageWarrior struct {
	Character
}

func (fw *FakeRageWarrior) GetCharacter() *Character { return &fw.Character }

func (fw *FakeRageWarrior) Initialize()                    {}
func (fw *FakeRageWarrior) ApplyTalents()                  {}
func (fw *FakeRageWarrior) Reset(_ *Simulation)            {}
func (fw *FakeRageWarrior) OnGCDReady(_ *Simulation)       {}
func (fw *FakeRageWarrior) OnEncounterStart(_ *Simulation) {}

func NewFakeRageWarrior(char *Character, _ *proto.Player, _ *proto.Raid) Agent {
	fw := &FakeRageWarrior{
		Character: *char,
	}

	fw.EnableRageBar(RageBarOptions{
		MaxRage:            100,
		BaseRageMultiplier: 1,
		StartingRage:       0,
	})

	fw.EnableAutoAttacks(fw, AutoAttackOptions{
		MainHand:       Weapon{SwingSpeed: fakeMHSwingSpeed},
		OffHand:        Weapon{SwingSpeed: fakeOHSwingSpeed},
		AutoSwingMelee: true,
	})

	return fw
}

func SetupFakeRageSim() *Simulation {
	sim := NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{
			RandomSeed: 100,
		},
		Raid: &proto.Raid{
			Parties: []*proto.Party{
				{
					Players: []*proto.Player{
						{
							Name:      "Warrior",
							Class:     proto.Class_ClassWarrior,
							Buffs:     &proto.IndividualBuffs{},
							Spec:      &proto.Player_DpsWarrior{},
							Equipment: &proto.EquipmentSpec{},
							// AddRage() pokes the rotation, so it needs to be non-nil.
							Rotation: &proto.APLRotation{Type: proto.APLRotation_TypeAPL},
						},
					},
					Buffs: &proto.PartyBuffs{},
				},
			},
		},
		Encounter: &proto.Encounter{
			Targets: []*proto.Target{
				{Name: "target", Level: 63, MobType: proto.MobType_MobTypeHumanoid},
			},
			Duration: 180,
		},
	}, simsignals.CreateSignals())
	sim.Reset()

	return sim
}

// Feeds a hand-crafted auto attack result to the Rage bar and returns the Rage gained from it.
// PostArmorAndResistanceMultiplier is always populated, because the sim fills it in before the
// outcome is applied, even for swings which end up doing no damage.
func rageFromAutoAttack(sim *Simulation, fw *FakeRageWarrior, spell *Spell, outcome HitOutcome, preOutcomeDamage float64) float64 {
	result := &SpellResult{
		Target:                           sim.Encounter.ActiveTargetUnits[0],
		Outcome:                          outcome,
		PostArmorAndResistanceMultiplier: preOutcomeDamage,
	}
	if outcome.Matches(OutcomeLanded) {
		result.Damage = preOutcomeDamage
	}

	rageBefore := fw.CurrentRage()
	fw.Unit.OnSpellHitDealt(sim, spell, result)

	return fw.CurrentRage() - rageBefore
}

func TestAutoAttackRageGeneration(t *testing.T) {
	// Forever pays a flat amount per landed swing, set by the weapon's speed:
	// a 2.6 speed one-hander gives 2.6 * 3.5 = 9.1, crit or not, whatever the damage.
	// Dodges, parries and misses give nothing.
	const swingDamage = 500.0

	tests := []struct {
		name     string
		outcome  HitOutcome
		wantRage float64
	}{
		{name: "hit", outcome: OutcomeHit, wantRage: 9.1},
		{name: "crit", outcome: OutcomeCrit, wantRage: 9.1},
		{name: "dodge", outcome: OutcomeDodge, wantRage: 0},
		{name: "parry", outcome: OutcomeParry, wantRage: 0},
		{name: "miss", outcome: OutcomeMiss, wantRage: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sim := SetupFakeRageSim()
			fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)

			gained := rageFromAutoAttack(sim, fw, fw.AutoAttacks.MHAuto(), test.outcome, swingDamage)

			if !WithinToleranceFloat64(test.wantRage, gained, 0.01) {
				t.Fatalf("Incorrect Rage generated on %s: Expected: %0.3f, Actual: %0.3f", test.name, test.wantRage, gained)
			}
		})
	}

	sim := SetupFakeRageSim()
	fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)
	if small, big := rageFromAutoAttack(sim, fw, fw.AutoAttacks.MHAuto(), OutcomeHit, 20),
		rageFromAutoAttack(sim, fw, fw.AutoAttacks.MHAuto(), OutcomeHit, 2000); small != big {
		t.Fatalf("Rage depends on damage: %0.3f for a 20 damage swing, %0.3f for 2000", small, big)
	}

	// Stances and talents still scale it through MultiplyAutoAttackRageGen.
	fw.MultiplyAutoAttackRageGen(1.5)
	if gained := rageFromAutoAttack(sim, fw, fw.AutoAttacks.MHAuto(), OutcomeHit, swingDamage); !WithinToleranceFloat64(13.65, gained, 0.01) {
		t.Fatalf("Multiplied MH hit: Expected: 13.650, Actual: %0.3f", gained)
	}
}

func TestOffHandAutoAttackRageGeneration(t *testing.T) {
	// OH swings pay the same rate on their own speed: 1.8 * 3.5 = 6.3.
	const swingDamage = 500.0

	sim := SetupFakeRageSim()
	fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)
	ohAuto := fw.AutoAttacks.OHAuto()

	hitRage := rageFromAutoAttack(sim, fw, ohAuto, OutcomeHit, swingDamage)
	if !WithinToleranceFloat64(6.3, hitRage, 0.01) {
		t.Fatalf("Incorrect Rage generated on OH hit: Expected: %0.3f, Actual: %0.3f", 6.3, hitRage)
	}

	if dodgeRage := rageFromAutoAttack(sim, fw, ohAuto, OutcomeDodge, swingDamage); dodgeRage != 0 {
		t.Fatalf("Dodged OH swing generated %0.3f Rage, expected none", dodgeRage)
	}
}

func TestOffHandRageMultiplier(t *testing.T) {
	// Dual Wield Specialization scales only the OH swing's flat rage: 1.8 * 3.5 * 2 = 12.6.
	sim := SetupFakeRageSim()
	fw := sim.Raid.Parties[0].Players[0].(*FakeRageWarrior)
	fw.SetOffHandRageMultiplier(2)

	if gained := rageFromAutoAttack(sim, fw, fw.AutoAttacks.OHAuto(), OutcomeHit, 500); !WithinToleranceFloat64(12.6, gained, 0.01) {
		t.Fatalf("Multiplied OH hit: Expected: 12.600, Actual: %0.3f", gained)
	}
	if gained := rageFromAutoAttack(sim, fw, fw.AutoAttacks.MHAuto(), OutcomeHit, 500); !WithinToleranceFloat64(9.1, gained, 0.01) {
		t.Fatalf("MH hit changed with the OH multiplier: Expected: 9.100, Actual: %0.3f", gained)
	}
}
