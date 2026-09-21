package core

import "github.com/wowsims/forever/sim/core/stats"

// Rage from a landed auto attack is flat on Forever: set by the weapon's speed and nothing
// else. A crit is worth the same as a hit, and a geared warrior the same as a levelling one.
//
// Measured from public beta combat logs by BrawnyBravo (ElliotWood/Forever issue #252): 63
// clean pairs across nine warriors at levels 10-15, consecutive auto-attack snapshots 0-4s
// apart with no ability used and no damage taken between them. Damage ran from 15 to 64 a
// hit, crits included, with no effect on the rage gained.
//
//	~2.1s one-hand   7.2-7.3      2.1 x 3.5 = 7.35
//	~2.5s one-hand   8.6-8.7      2.5 x 3.5 = 8.75
//	~3.2s two-hand  14.4          3.2 x 4.5 = 14.4
//	~3.3s two-hand  14.9          3.3 x 4.5 = 14.85
//	~3.5s two-hand  15.7          3.5 x 4.5 = 15.75
//
// Base weapon speed, not the hasted interval, so haste buys rage. Nobody in the sample had
// haste or dual wielded, so both of those (and an off-hand swing paying the same rate as a
// main-hand one) are assumptions, the first things to re-measure.
const (
	ForeverRagePerSecondOneHand = 3.5
	ForeverRagePerSecondTwoHand = 4.5
)

func ForeverWhiteHitRage(weapon *Weapon) float64 {
	if weapon == nil {
		return 0
	}
	if weapon.TwoHand {
		return weapon.SwingSpeed * ForeverRagePerSecondTwoHand
	}
	return weapon.SwingSpeed * ForeverRagePerSecondOneHand
}

// Forever pays hit and critical strike from gear against every kind of attack rather than
// splitting them into a melee and a spell pool. The client data states most of it as the
// generic rating, which the database reads as melee, so a caster would otherwise get no hit
// from Neltharion's Tear. Summed in percent, then paid into both pools as rating.
func unifyGearHitAndCrit(equipStats stats.Stats) stats.Stats {
	hit := equipStats[stats.MeleeHitRating]/PhysicalHitRatingPerHitPercent + equipStats[stats.SpellHitRating]/SpellHitRatingPerHitPercent
	crit := equipStats[stats.MeleeCritRating]/PhysicalCritRatingPerCritPercent + equipStats[stats.SpellCritRating]/SpellCritRatingPerCritPercent

	equipStats[stats.MeleeHitRating] = hit * PhysicalHitRatingPerHitPercent
	equipStats[stats.SpellHitRating] = hit * SpellHitRatingPerHitPercent
	equipStats[stats.MeleeCritRating] = crit * PhysicalCritRatingPerCritPercent
	equipStats[stats.SpellCritRating] = crit * SpellCritRatingPerCritPercent
	return equipStats
}
