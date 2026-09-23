package database

import (
	"testing"

	"github.com/wowsims/forever/sim/core"
)

// A trigger restricted to one named ability is a shape no ProcTypeMask states, so the wording is the
// only evidence there is. The rows below are the client's own, raw and rendered: the store reads the
// raw description and the item generator the rendered tooltip, and both go through this matcher.
func TestNamedAbilityWording(t *testing.T) {
	for _, tc := range []struct {
		named bool
		text  string
		why   string
	}{
		{true, "Your casts of $?s2060[Greater Heal]?s5185[Healing Touch][Greater Heal, Healing Touch] in combat grant up to $1249119s1 increased healing.",
			"Eternal Power 1249118, as the client writes it"},
		{true, "Your casts of Greater Heal in combat grant up to 40 increased healing and up to 13 increased damage for 15s.",
			"the same tooltip, rendered"},
		{true, "Your Backstab has a $m1% chance to cause your next Ambush to not require Stealth.",
			`"has" after a named ability`},
		{true, "Heals from your Earth Shield have a $s1% chance to make your next cast time heal instant cast.",
			`"have" after a named ability`},
		{true, "When you cast Chain Heal or Riptide you gain $s1 charges of Tidal Waves.",
			"a cast of one named ability"},
		{true, "Your Shock spells have a chance to deal extra damage.", "the wording the matcher already read"},
		{false, "When you cast a healing spell, gain Mana equal to $m1% of the base cost of the spell.",
			"no ability is named"},
		{false, "20% chance to regain 100 mana when you cast a Judgement.",
			"a lowercase article keeps the class of spell general"},
		{false, "Multiple casts of Divine Light do not accumulate this shield.",
			"not the caster's own casts, and not a trigger clause"},
		{false, "Your melee attacks have a chance to deal extra damage.",
			"no ability is named"},
		{false, "2% chance on successful spellcast to increase your Spirit by 150 for 15s.",
			"Darkmoon Card: Blue Dragon, which fires off any cast"},
	} {
		if got := procTooltipHints(tc.text).Matches(core.ProcHintNamedAbility); got != tc.named {
			t.Errorf("named ability = %v, want %v for %s:\n  %s", got, tc.named, tc.why, tc.text)
		}
	}
}

// The wearer's own attack dodged or parried is a trigger only in a condition clause. The same words
// state a magnitude on every expertise row, and "when you parry" is an attack the wearer takes.
func TestAttackAvoidedWording(t *testing.T) {
	const avoided = core.ProcHintAttackDodged | core.ProcHintAttackParried

	for _, tc := range []struct {
		want core.ProcHint
		text string
		why  string
	}{
		{avoided, "Permanently enchant a Melee Weapon to trigger Recovery when you are Parried or Dodged, healing you for 5% of your maximum health. Cannot occur more often than once every 10 sec.",
			"Recovery's grant 1248760, rendered"},
		{core.ProcHintAttackParried, "Whenever your melee attacks are parried, gain 10 rage.",
			"one outcome named, one bit"},
		{0, "Reduces the chance for your attacks to be dodged or parried by $s1%.",
			"Increased Expertise 1213288, a magnitude"},
		{0, "Your Taunt ability never misses, and your chance to be Dodged or Parried is reduced by $s1%.",
			"the Naxxramas tank 2P 1219540, a magnitude"},
		{0, "Instantly overpower the enemy, causing weapon damage plus $s1.  Only useable after the target dodges.  The Overpower cannot be blocked, dodged or parried.",
			"Overpower 7384, which states what it cannot be"},
		{0, "Your Shield Slam deals $s1% increased threat and its cooldown is reset if it is Dodged, Parried, or Blocked.",
			"TAQ tank 4P 1214162, one named ability rather than the wearer's attacks"},
		{0, "When you parry an attack, gain 10 rage.",
			"the wearer parrying an attack it takes"},
	} {
		if got := procTooltipHints(tc.text) & avoided; got != tc.want {
			t.Errorf("hint = %q, want %q for %s:\n  %s", formatProcHint(got), formatProcHint(tc.want), tc.why, tc.text)
		}
	}
}
