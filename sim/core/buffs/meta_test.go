package buffs

import (
	"math"
	"testing"
	"time"

	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

// The paladin auras state a healing-taken row of 0, and Concentration Aura two mechanic rows of 0, none
// of which the raid buff applies.
var paladinAuraSkips = []dbcenums.EffectAuraType{dbcenums.A_MOD_HEALING_PCT, dbcenums.A_MECHANIC_DURATION_MOD}

// A Meta reads the same numbers off the store as the constructors generated for each row, which is what
// lets the generator write a Meta in their place.
func TestMetaReadsWhatTheGeneratedRowsRead(t *testing.T) {
	find := spelldata.MustFind
	cases := []struct {
		name     string
		meta     *Meta
		value    func(int32) float64
		duration func(int32) time.Duration
		ranks    int32
	}{
		{"BloodPact", &Meta{Spell: find(11767)}, BloodPactValue, BloodPactDuration, 0},
		{"BattleShout", &Meta{Spell: find(25289)}, BattleShoutValue, BattleShoutDuration, 0},
		{"DevotionAura", &Meta{Spell: find(10293), SkipAuras: paladinAuraSkips}, DevotionAuraValue, DevotionAuraDuration, 0},
		{"LeaderOfThePack", &Meta{Spell: find(24932)}, LeaderOfThePackValue, LeaderOfThePackDuration, 0},
		{"ManaSpringTotem", &Meta{Spell: find(10494), Talent: spelldata.Talent(16187, 5), TalentEffect: 1},
			ManaSpringTotemValue, ManaSpringTotemDuration, 5},
		{"ManaTideTotems", &Meta{Spell: find(17360), Cast: find(17359)}, ManaTideTotemsValue, ManaTideTotemsDuration, 0},
		{"RetributionAura", &Meta{Spell: find(10301), SkipAuras: paladinAuraSkips}, RetributionAuraValue, RetributionAuraDuration, 0},
		{"ConcentrationAura", &Meta{Spell: find(19746), SkipAuras: paladinAuraSkips}, ConcentrationAuraValue, ConcentrationAuraDuration, 0},
		{"TrueshotAura", &Meta{Spell: find(20906)}, TrueshotAuraValue, TrueshotAuraDuration, 0},
		{"AtieshWarlock", &Meta{Spell: find(28143)}, AtieshWarlockValue, AtieshWarlockDuration, 0},
		{"AtieshDruid", &Meta{Spell: find(28145)}, AtieshDruidValue, AtieshDruidDuration, 0},
		{"WindfuryTotem", &Meta{Spell: find(10610)}, WindfuryTotemValue, WindfuryTotemDuration, 0},
		{"GreaterBlessingOfKings", &Meta{Spell: find(25898)}, GreaterBlessingOfKingsValue, GreaterBlessingOfKingsDuration, 0},
		{"GreaterBlessingOfWisdom", &Meta{Spell: find(25918)}, GreaterBlessingOfWisdomValue, GreaterBlessingOfWisdomDuration, 0},
		{"GreaterBlessingOfSalvation", &Meta{Spell: find(25895)}, GreaterBlessingOfSalvationValue, GreaterBlessingOfSalvationDuration, 0},
		{"GiftOfTheWild", &Meta{Spell: find(21850)}, GiftOfTheWildValue, GiftOfTheWildDuration, 0},
		{"Thorns", &Meta{Spell: find(9910)}, ThornsValue, ThornsDuration, 0},
		{"FireResistanceAura", &Meta{Spell: find(19900), SkipAuras: paladinAuraSkips}, FireResistanceAuraValue, FireResistanceAuraDuration, 0},
		{"PowerInfusions", &Meta{Spell: find(10060)}, PowerInfusionsValue, PowerInfusionsDuration, 0},
		{"HuntersMark", &Meta{Spell: find(14325)}, HuntersMarkValue, HuntersMarkDuration, 0},
		{"CurseOfElements", &Meta{Spell: find(1311680)}, CurseOfElementsValue, CurseOfElementsDuration, 0},
		{"CurseOfRecklessness", &Meta{Spell: find(11717)}, CurseOfRecklessnessValue, CurseOfRecklessnessDuration, 0},
		{"ExposeArmor", &Meta{Spell: find(11198), FullComboPoints: true}, ExposeArmorValue, ExposeArmorDuration, 0},
		{"SunderArmor", &Meta{Spell: find(11597)}, SunderArmorValue, SunderArmorDuration, 0},
		{"GiftOfArthas", &Meta{Spell: find(11374)}, GiftOfArthasValue, GiftOfArthasDuration, 0},
		{"DemoralizingShout", &Meta{Spell: find(11556)}, DemoralizingShoutValue, DemoralizingShoutDuration, 0},
		{"ThunderClap", &Meta{Spell: find(11581)}, ThunderClapValue, ThunderClapDuration, 0},
		{"InsectSwarm", &Meta{Spell: find(24977)}, InsectSwarmValue, InsectSwarmDuration, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for points := int32(0); points <= c.ranks; points++ {
				if got, want := c.meta.Value(points), c.value(points); math.Abs(got-want) > 1e-9 {
					t.Errorf("Value(%d) = %v, the generated row reads %v", points, got, want)
				}
				if got, want := c.meta.Duration(points), c.duration(points); got != want {
					t.Errorf("Duration(%d) = %v, the generated row reads %v", points, got, want)
				}
			}
		})
	}

	cooldowns := []struct {
		name     string
		meta     *Meta
		cooldown func() time.Duration
	}{
		{"ManaTideTotems", &Meta{Spell: find(17360), Cast: find(17359)}, ManaTideTotemsCooldown},
		{"Innervates", &Meta{Spell: find(29166)}, InnervatesCooldown},
		{"PowerInfusions", &Meta{Spell: find(10060)}, PowerInfusionsCooldown},
	}
	for _, c := range cooldowns {
		if got, want := c.meta.Cooldown(), c.cooldown(); got != want {
			t.Errorf("%s: Cooldown() = %v, the generated row reads %v", c.name, got, want)
		}
	}
}

// The raid's Retribution Aura is a damage shield and a healing-taken row of 0: left in, the 0 would be
// the first amount the parse attaches, a multiplier of 1.
func TestMetaValueNeedsThePaladinSkips(t *testing.T) {
	if got := (&Meta{Spell: spelldata.MustFind(10301)}).Value(0); got != 1 {
		t.Errorf("Retribution Aura without the skips reads %v, want the healing-taken row's 1", got)
	}
}
