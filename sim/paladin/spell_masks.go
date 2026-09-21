package paladin

// TBC Paladin Spell Masks
// These are used for spell identification, talent modifiers, and proc triggers

const (
	SpellMaskNone int64 = 0

	// Core Abilities
	SpellMaskJudgement int64 = 1 << iota
	SpellMaskConsecration
	SpellMaskExorcism
	SpellMaskHolyLight
	SpellMaskFlashOfLight
	SpellMaskLayOnHands
	SpellMaskHammerOfJustice
	SpellMaskCleanse
	SpellMaskDivineShield
	SpellMaskDivineProtection
	SpellMaskBlessingOfProtection
	SpellMaskHammerOfWrath
	SpellMaskHolyWrath

	// Seals
	SpellMaskSealOfRighteousness
	SpellMaskSealOfCommand
	SpellMaskSealOfLight
	SpellMaskSealOfWisdom
	SpellMaskSealOfJustice
	SpellMaskSealOfTheCrusader

	// Judgement Effects (different from the Judgement spell itself)
	SpellMaskJudgementOfRighteousness
	SpellMaskJudgementOfCommand
	SpellMaskJudgementOfLight
	SpellMaskJudgementOfWisdom
	SpellMaskJudgementOfJustice
	SpellMaskJudgementOfTheCrusader

	// Auras
	SpellMaskDevotionAura
	SpellMaskRetributionAura
	SpellMaskConcentrationAura
	SpellMaskFireResistanceAura
	SpellMaskFrostResistanceAura
	SpellMaskShadowResistanceAura

	// Blessings
	SpellMaskBlessingOfMight
	SpellMaskBlessingOfWisdom
	SpellMaskBlessingOfKings
	SpellMaskBlessingOfSalvation
	SpellMaskBlessingOfSanctuary

	// Talent Abilities
	SpellMaskDivineFavor
	SpellMaskHolyShock
	SpellMaskHolyShield
	SpellMaskHolyShieldProc
	SpellMaskRepentance
	SpellMaskRighteousFury
)

// Composite masks
const (
	SpellMaskAllSeals = SpellMaskSealOfRighteousness |
		SpellMaskSealOfCommand |
		SpellMaskSealOfLight |
		SpellMaskSealOfWisdom |
		SpellMaskSealOfJustice |
		SpellMaskSealOfTheCrusader

	SpellMaskAllJudgements = SpellMaskJudgementOfRighteousness |
		SpellMaskJudgementOfCommand |
		SpellMaskJudgementOfLight |
		SpellMaskJudgementOfWisdom |
		SpellMaskJudgementOfJustice |
		SpellMaskJudgementOfTheCrusader

	SpellMaskAllAuras = SpellMaskDevotionAura |
		SpellMaskRetributionAura |
		SpellMaskConcentrationAura |
		SpellMaskFireResistanceAura |
		SpellMaskFrostResistanceAura |
		SpellMaskShadowResistanceAura

	SpellMaskAllBlessings = SpellMaskBlessingOfMight |
		SpellMaskBlessingOfWisdom |
		SpellMaskBlessingOfKings |
		SpellMaskBlessingOfSalvation |
		SpellMaskBlessingOfSanctuary

	SpellMaskHealingSpells = SpellMaskHolyLight |
		SpellMaskFlashOfLight |
		SpellMaskLayOnHands |
		SpellMaskHolyShock

	// Spells that can trigger Seal of Command
	SpellMaskCanTriggerSealOfCommand = SpellMaskJudgement

	SpellMaskCanProcTome = SpellMaskAllAuras |
		SpellMaskAllBlessings |
		SpellMaskAllSeals |
		SpellMaskAllJudgements |
		SpellMaskConsecration |
		SpellMaskDivineFavor |
		SpellMaskExorcism |
		SpellMaskHammerOfWrath |
		SpellMaskHealingSpells |
		SpellMaskHolyShield |
		SpellMaskHolyShock |
		SpellMaskHolyWrath |
		SpellMaskRighteousFury
)
