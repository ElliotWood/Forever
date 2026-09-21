package core

import "fmt"

// Named bits of SpellAuraOptions.ProcTypeMask word 0, under TrinityCore's names for them. Word 1
// carries no bit the sim models.
const (
	ProcFlagKilled              uint32 = 1 << 0  // 00 The caster was killed
	ProcFlagKill                uint32 = 1 << 1  // 01 The caster killed the target
	ProcFlagDealMeleeSwing      uint32 = 1 << 2  // 02 Dealt a melee auto attack
	ProcFlagTakeMeleeSwing      uint32 = 1 << 3  // 03 Took a melee auto attack
	ProcFlagDealMeleeAbility    uint32 = 1 << 4  // 04 Dealt a spell of damage class melee
	ProcFlagTakeMeleeAbility    uint32 = 1 << 5  // 05 Took a spell of damage class melee
	ProcFlagDealRangedAttack    uint32 = 1 << 6  // 06 Dealt a ranged auto attack
	ProcFlagTakeRangedAttack    uint32 = 1 << 7  // 07 Took a ranged auto attack
	ProcFlagDealRangedAbility   uint32 = 1 << 8  // 08 Dealt a spell of damage class ranged
	ProcFlagTakeRangedAbility   uint32 = 1 << 9  // 09 Took a spell of damage class ranged
	ProcFlagDealHelpfulAbility  uint32 = 1 << 10 // 10 Dealt a positive spell of damage class none
	ProcFlagTakeHelpfulAbility  uint32 = 1 << 11 // 11 Took a positive spell of damage class none
	ProcFlagDealHarmfulAbility  uint32 = 1 << 12 // 12 Dealt a negative spell of damage class none
	ProcFlagTakeHarmfulAbility  uint32 = 1 << 13 // 13 Took a negative spell of damage class none
	ProcFlagDealHelpfulSpell    uint32 = 1 << 14 // 14 Dealt a positive spell of damage class magic
	ProcFlagTakeHelpfulSpell    uint32 = 1 << 15 // 15 Took a positive spell of damage class magic
	ProcFlagDealHarmfulSpell    uint32 = 1 << 16 // 16 Dealt a negative spell of damage class magic
	ProcFlagTakeHarmfulSpell    uint32 = 1 << 17 // 17 Took a negative spell of damage class magic
	ProcFlagDealHarmfulPeriodic uint32 = 1 << 18 // 18 Dealt periodic damage
	ProcFlagTakeHarmfulPeriodic uint32 = 1 << 19 // 19 Took periodic damage
	ProcFlagTakeAnyDamage       uint32 = 1 << 20 // 20 Took damage of any kind
	ProcFlagDealHelpfulPeriodic uint32 = 1 << 21 // 21 Dealt periodic healing
	ProcFlagMainHandWeaponSwing uint32 = 1 << 22 // 22 Dealt a main-hand melee attack, auto or ability
	ProcFlagOffHandWeaponSwing  uint32 = 1 << 23 // 23 Dealt an off-hand melee attack, auto or ability
	ProcFlagDeath               uint32 = 1 << 24 // 24 The caster died
	ProcFlagJump                uint32 = 1 << 25 // 25 The caster jumped
	ProcFlagEnterCombat         uint32 = 1 << 27 // 27 The caster entered combat
	ProcFlagEncounterStart      uint32 = 1 << 28 // 28 The encounter started
	ProcFlagCastEnded           uint32 = 1 << 29 // 29 A cast ended, however it ended
	ProcFlagLooted              uint32 = 1 << 30 // 30 The caster looted
)

// The taken bits that name a direct hit arriving on the character. The two helpful takes are out
// because the sim models neither - one comes back as unsupported, the other is ignored - and the
// any-damage take has a rule of its own below.
const procFlagAnyDirectTaken = ProcFlagTakeMeleeSwing |
	ProcFlagTakeMeleeAbility |
	ProcFlagTakeRangedAttack |
	ProcFlagTakeRangedAbility |
	ProcFlagTakeHarmfulAbility |
	ProcFlagTakeHarmfulSpell

// Every direct hit dealt.
const procFlagAnyDirectDealt = ProcFlagDealMeleeSwing |
	ProcFlagDealMeleeAbility |
	ProcFlagDealRangedAttack |
	ProcFlagDealRangedAbility |
	ProcFlagDealHarmfulSpell

// What the mask alone cannot say, read off the spell's tooltip by the caller. The mask states
// which hits reach the listener; the wording around it states the trigger condition and, for the
// helpful bits, whether they mean anything at all.
type ProcHint uint8

const (
	// The tooltip names the cast itself as the trigger: "each time you cast a spell", "chance on
	// spell cast".
	ProcHintCastTrigger ProcHint = 1 << iota
	// The tooltip names a critical strike as the trigger.
	ProcHintCrit
	// The tooltip's trigger clause names healing, or an unrestricted "a spell". Either is the
	// evidence that a helpful-spell bit carries a real trigger rather than a leftover.
	ProcHintHeals
	// The tooltip restricts the trigger to healing spells.
	ProcHintPureHeal
	// The trigger clause names one ability: "Your Shock spells", "Your Moonfire ability".
	ProcHintNamedAbility
	// The trigger clause states an attack outcome the mask has no bit for: a resist, a block, a
	// dodge or a parry.
	ProcHintOutcomeTaken
)

// Returns whether there is any overlap between the given hints.
func (h ProcHint) Matches(other ProcHint) bool {
	return (h & other) != 0
}

// The listener shape a ProcTypeMask describes, in the sim's own terms.
type ProcTypeInfo struct {
	Callback           AuraCallback
	ProcMask           ProcMask
	Outcome            HitOutcome
	RequireDamageDealt bool
	Unsupported        []string // named bits the decoder does not model, empty when supported
}

// DecodeProcTypeMask reads SpellAuraOptions.ProcTypeMask (two 32-bit words).
//
// ProcHintNamedAbility and ProcHintOutcomeTaken are decoded past: a trigger restricted to one
// named ability, or to an outcome the mask has no bit for, is a shape no ProcTypeMask can state,
// so the caller refuses those spells rather than the decoder.
func DecodeProcTypeMask(mask [2]uint32, hint ProcHint) ProcTypeInfo {
	info := ProcTypeInfo{RequireDamageDealt: true}
	word := mask[0]

	if word&ProcFlagDealMeleeSwing != 0 {
		info.ProcMask |= ProcMaskMeleeWhiteHit
	}

	if word&ProcFlagDealMeleeAbility != 0 {
		info.ProcMask |= ProcMaskMeleeSpecial
	}

	if word&ProcFlagDealRangedAttack != 0 {
		info.ProcMask |= ProcMaskRangedAuto
	}

	if word&ProcFlagDealRangedAbility != 0 {
		info.ProcMask |= ProcMaskRangedSpecial
	}

	if word&(ProcFlagDealHarmfulPeriodic|ProcFlagDealHarmfulSpell) != 0 {
		info.ProcMask |= ProcMaskSpellDamage
	}

	if word&procFlagAnyDirectTaken != 0 {
		info.Callback |= CallbackOnSpellHitTaken

		if word&ProcFlagTakeMeleeSwing != 0 {
			info.ProcMask |= ProcMaskMeleeWhiteHit
		}

		if word&ProcFlagTakeMeleeAbility != 0 {
			info.ProcMask |= ProcMaskMeleeSpecial
		}

		if word&ProcFlagTakeRangedAttack != 0 {
			info.ProcMask |= ProcMaskRangedAuto
		}

		if word&ProcFlagTakeRangedAbility != 0 {
			info.ProcMask |= ProcMaskRangedSpecial
		}

		if word&ProcFlagTakeHarmfulSpell != 0 {
			info.ProcMask |= ProcMaskSpellDamage
		}
	}

	if word&ProcFlagTakeHarmfulPeriodic != 0 {
		info.Callback |= CallbackOnPeriodicDamageTaken
	}

	// Damage of any kind, however it arrived, which is both of the taken callbacks. The bit names
	// damage rather than a hit, so a landed hit dealing none does not count - the default this
	// decode starts from, and no client mask pairs this bit with one that clears it.
	if word&ProcFlagTakeAnyDamage != 0 {
		info.Callback |= CallbackOnSpellHitTaken | CallbackOnPeriodicDamageTaken
	}

	// A mask made of nothing but the spell-cast bits. The harmful one has to be present: a
	// helpful-only mask carries no evidence that casting is the trigger at all, and the helpful
	// branch below already demands tooltip evidence before it believes one - the PvP Librams
	// that buff a heal target read "Causes your Flash of Light to increase the target's
	// Resilience" and are neither a self buff nor unrestricted.
	spellCastMask := word&ProcFlagDealHarmfulSpell != 0 &&
		word&^(ProcFlagDealHarmfulSpell|ProcFlagDealHelpfulSpell) == 0

	// Whether the cast itself is the trigger. A mask of only the harmful-spell bit does not care
	// whether the spell landed. Adding the helpful bit settles nothing either way, and the two
	// items that pin it down disagree despite carrying the identical mask: Memento of Tyrande
	// procs off resists in logs, while Band of the Eternal Restorer does not proc on a miss or a
	// full resist. What separates them is that the first names the cast as the trigger and the
	// second does not, so for that pair the tooltip decides.
	castOnly := spellCastMask && (word == ProcFlagDealHarmfulSpell || hint.Matches(ProcHintCastTrigger))

	// A tooltip naming an outcome is the exception to all of it: a crit is only known once the
	// hit resolves, so those stay on hit-dealt.
	if castOnly && !hint.Matches(ProcHintCrit) {
		info.Callback |= CallbackOnCastComplete
		info.RequireDamageDealt = false
	} else if word&procFlagAnyDirectDealt != 0 {
		info.Callback |= CallbackOnSpellHitDealt

		if word&ProcFlagDealHarmfulSpell != 0 {
			info.RequireDamageDealt = false
		}
	}

	if word&ProcFlagDealHarmfulPeriodic != 0 {
		info.Callback |= CallbackOnPeriodicDamageDealt
	}

	if word&ProcFlagDealHelpfulSpell != 0 && hint.Matches(ProcHintHeals) {
		info.RequireDamageDealt = false
		info.ProcMask |= ProcMaskSpellHealing

		// Casting the heal is already the trigger above, so adding heal-dealt on top would
		// proc twice for one heal.
		if !info.Callback.Matches(CallbackOnCastComplete) {
			info.Callback |= CallbackOnHealDealt

			// handle HoTs only with direct heals for now, there are some odd cases with HoT / DoT overlaps
			if word&ProcFlagDealHelpfulPeriodic != 0 {
				info.Callback |= CallbackOnPeriodicHealDealt
			}

			// Check if we have periodic damage flag but only heal paired with it
			// This usually indicates a pure heal proc mask
			if word&procFlagAnyDirectDealt == 0 {
				info.Callback &= ^CallbackOnPeriodicDamageDealt
				info.Callback &= ^CallbackOnSpellHitDealt
				info.ProcMask &= ^ProcMaskSpellDamage
			}
		}
	}

	// A mask naming one hand hears that hand's melee hits only. The bit restricts melee and
	// nothing else: a mask pairing it with a spell or ranged bit keeps those hits, so this drops
	// the other hand's melee bits rather than intersecting the whole proc mask. Naming both hands
	// is the same as naming neither, so only a mask with exactly one of the two restricts
	// anything.
	switch word & (ProcFlagMainHandWeaponSwing | ProcFlagOffHandWeaponSwing) {
	case ProcFlagMainHandWeaponSwing:
		info.ProcMask &= ^ProcMaskMeleeOH
	case ProcFlagOffHandWeaponSwing:
		info.ProcMask &= ^ProcMaskMeleeMH
	}

	// An outcome the listener can be given. Only the crit hint names one: the mask itself states
	// which hits arrive, never how they resolved, so everything else listens to a landed hit. A
	// cast has not resolved into a hit at all and takes no outcome.
	switch {
	case info.Callback.Matches(CallbackOnCastComplete):
		info.Outcome = OutcomeEmpty
	case hint.Matches(ProcHintCrit):
		info.Outcome = OutcomeCrit
	default:
		info.Outcome = OutcomeLanded
	}

	if hint.Matches(ProcHintPureHeal) {
		info.Callback &= ^CallbackOnSpellHitDealt
		info.Callback &= ^CallbackOnPeriodicDamageDealt
	}

	info.Unsupported = unsupportedProcFlags(mask)

	return info
}

var unsupportedProcFlagNames = map[uint32]string{
	ProcFlagKilled:           "KILLED",
	ProcFlagKill:             "KILL",
	ProcFlagTakeHelpfulSpell: "TAKE_HELPFUL_SPELL",
	ProcFlagDeath:            "DEATH",
	ProcFlagJump:             "JUMP",
	ProcFlagEnterCombat:      "ENTER_COMBAT",
	ProcFlagEncounterStart:   "ENCOUNTER_START",
	ProcFlagCastEnded:        "CAST_ENDED",
	ProcFlagLooted:           "LOOTED",
}

// The bits the shape above says nothing about. The three damage-class-none bits (0x400 and 0x1000
// dealt, 0x800 taken) are left out on purpose: next to real bits they change nothing the decode
// says, and on their own they leave the callback empty, which the caller refuses anyway.
func unsupportedProcFlags(mask [2]uint32) []string {
	// Everything from the death bit up is a state change rather than a hit, and word 1 names
	// nothing the sim models.
	unsupported := mask[0] & (ProcFlagKilled | ProcFlagKill | ProcFlagTakeHelpfulSpell | ^(ProcFlagDeath - 1))

	var names []string
	for bit := 0; bit < 32; bit++ {
		if unsupported&(1<<bit) == 0 {
			continue
		}

		if name, ok := unsupportedProcFlagNames[1<<bit]; ok {
			names = append(names, name)
		} else {
			names = append(names, fmt.Sprintf("bit %d", bit))
		}
	}

	for bit := 0; bit < 32; bit++ {
		if mask[1]&(1<<bit) != 0 {
			names = append(names, fmt.Sprintf("bit %d", 32+bit))
		}
	}

	return names
}
