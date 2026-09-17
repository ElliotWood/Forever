package hunter

import (
	"testing"

	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

// Forever pairs Aimed Shot's cooldown with Multi-Shot rather than Classic's Arcane Shot,
// and both tooltips say so outright: "Aimed Shot shares its cooldown with Multi-Shot" and
// the reverse on Multi-Shot's own. The pairing is wiring rather than a number, so nothing
// in the results files would notice if it were undone - the DPS would simply drift.
func TestAimedShotSharesMultiShotsCooldownUnderForever(t *testing.T) {
	for _, tc := range []struct {
		ruleset    proto.Ruleset
		wantShared bool
	}{
		{proto.Ruleset_RulesetForever, true},
		{proto.Ruleset_RulesetClassic, false},
	} {
		t.Run(tc.ruleset.String(), func(t *testing.T) {
			options := googleProto.Clone(core.DefaultSimTestOptions).(*proto.SimOptions)
			options.Ruleset = tc.ruleset

			raid := core.SinglePlayerRaidProto(&proto.Player{
				Class:         proto.Class_ClassHunter,
				Race:          proto.Race_RaceOrc,
				Equipment:     core.GetGearSet("../../ui/hunter/gear_sets", "p0.bis").GearSet,
				TalentsString: P1Talents,
				Rotation:      core.GetAplRotation("../../ui/hunter/apls", "p1").Rotation,
				Spec:          P1PlayerOptions,
			}, nil, nil, nil)

			environment, _, _ := core.NewEnvironment(raid, core.MakeSingleTargetEncounter(0), tc.ruleset, false)
			character := environment.Raid.Parties[0].Players[0].GetCharacter()

			// The highest rank of each is the one the level 60 rotation casts; Forever has one Multi-Shot.
			aimed := character.GetSpell(core.ActionID{SpellID: 20904})
			multi := character.GetSpell(core.ActionID{SpellID: 2643})
			arcane := character.GetSpell(core.ActionID{SpellID: 14287})
			if aimed == nil || multi == nil {
				t.Fatal("Aimed Shot or Multi-Shot is not registered")
			}

			if shared := aimed.CD.Timer == multi.CD.Timer; shared != tc.wantShared {
				t.Errorf("Aimed Shot shares Multi-Shot's timer = %v, want %v", shared, tc.wantShared)
			}
			if arcane != nil {
				if shared := aimed.CD.Timer == arcane.CD.Timer; shared == tc.wantShared {
					t.Errorf("Aimed Shot shares Arcane Shot's timer = %v, want %v", shared, !tc.wantShared)
				}
			}
		})
	}
}
