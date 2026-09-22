package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

var TalentTreeSizes = [3]int{18, 16, 18}

type Paladin struct {
	core.Character

	Talents *proto.PaladinTalents

	Forbearance       *core.Aura
	RighteousFuryAura *core.Aura

	Judgement *core.Spell

	// The seal the paladin is under, and the ranks it has cast so far this iteration.
	currentSeal *sealConfig

	// The judgement effects this paladin can put on enemies; its melee strikes refresh them.
	JudgementAuras []core.AuraArray

	// Twist of Light: one Echo per seal that leaves one, keyed by the Echo's spell id.
	echoes map[int32]*sealEcho

	// Consecrated Ground: the targets the talent's bonus is active on, one aura per enemy.
	consecratedGroundAuras core.AuraArray

	// Light's Vigil: one entry per rank, so Holy Shock can find the vigil it consumes.
	lightsVigils []*lightsVigil

	// Timers shared by the ranks of one ability.
	judgementTimer     *core.Timer
	holyStrikeTimer    *core.Timer
	consecrationTimer  *core.Timer
	exorcismTimer      *core.Timer
	hammerOfWrathTimer *core.Timer
	holyWrathTimer     *core.Timer
	holyShockTimer     *core.Timer
	holyShieldTimer    *core.Timer
	lightsVigilTimer   *core.Timer
}

// Implemented by each Paladin spec.
type PaladinAgent interface {
	GetPaladin() *Paladin
}

func (paladin *Paladin) GetCharacter() *core.Character {
	return &paladin.Character
}

func (paladin *Paladin) GetPaladin() *Paladin {
	return paladin
}

func (paladin *Paladin) AddRaidBuffs(_ *proto.RaidBuffs) {
}

func (paladin *Paladin) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (paladin *Paladin) Initialize() {
	paladin.registerForbearance()
	paladin.registerRighteousFury()

	paladin.registerJudgement()
	paladin.registerSeals()
	paladin.registerAuras()

	HolyStrikeRankMap.RegisterAll(paladin.registerHolyStrike)
	paladin.registerHammerOfTheRighteous()
	ConsecrationRankMap.RegisterAll(paladin.registerConsecration)
	ExorcismRankMap.RegisterAll(paladin.registerExorcism)
	HammerOfWrathRankMap.RegisterAll(paladin.registerHammerOfWrath)
	HolyWrathRankMap.RegisterAll(paladin.registerHolyWrath)

	HolyLightRankMap.RegisterAll(paladin.registerHolyLight)
	FlashOfLightRankMap.RegisterAll(paladin.registerFlashOfLight)
	LayOnHandsRankMap.RegisterAll(paladin.registerLayOnHands)
}

func (paladin *Paladin) Reset(_ *core.Simulation) {
	paladin.currentSeal = nil
	for _, echo := range paladin.echoes {
		echo.seal = nil
	}
}

func (paladin *Paladin) OnEncounterStart(_ *core.Simulation) {
}

func (paladin *Paladin) GetMainHandType() proto.HandType {
	mh := paladin.GetMHWeapon()

	if mh != nil && (mh.HandType == proto.HandType_HandTypeTwoHand) {
		return proto.HandType_HandTypeTwoHand
	}

	return proto.HandType_HandTypeOneHand
}

func (paladin *Paladin) sharedTimer(timer **core.Timer) *core.Timer {
	if *timer == nil {
		*timer = paladin.NewTimer()
	}
	return *timer
}

// The cost a row states: a flat number, or a share of base mana for the spells the client prices
// that way (Judgement, Righteous Fury, Seal of Justice).
func manaCost(row shared.SpellData) core.ManaCostOptions {
	if pct := powerCostPct[row.SpellID]; pct > 0 {
		return core.ManaCostOptions{BaseCostPercent: pct}
	}
	return core.ManaCostOptions{FlatCost: row.Cost}
}

// SpellPower.PowerCostPct for the paladin spells priced as a share of base mana. forever-next's
// generator does not carry the column yet (master's and their #48 stack's do), so the rows read
// Cost 0; these are the values those generators emit.
// ponytail: delete once gen_spelldata writes PowerCostPct into spell_data_auto_gen.go.
var powerCostPct = map[int32]float64{
	1044: 10, 20217: 8, 10278: 7, 1038: 8, 4987: 8, 20216: 4, 25895: 16, 407632: 6, 20271: 6,
	1152: 8, 7328: 80, 10322: 80, 10324: 80, 20772: 80, 20773: 80, 25780: 30, 20164: 13,
}

// A direct-damage row's roll. The client rolls these between a min and a max, and forever-next's
// generator stores only the centre of the range, truncated (the damage-range-as-mean gap), so the
// ranges our beta-client values give stay here, at each rank's max level; every centre below
// truncates to the table's number.
// ponytail: drop once gen_spelldata carries the client's damage variance.
var damageRanges = map[int32][2]float64{
	// Judgement of Righteousness
	20187: {24.8, 26.8}, 20280: {36.4, 38.4}, 20281: {53.4, 57.4}, 20282: {73.8, 79.8}, 20283: {96.6, 104.6}, 20284: {124.8, 134.8}, 20285: {155.6, 167.6}, 20286: {170.2, 186.2},
	// Judgement of Command, keyed by the judgement dummies the rows hang off (20425 casts 20467 and
	// so on); the stunned number, which the judgement halves otherwise
	20425: {137.8, 145.8}, 20962: {194.8, 208.8}, 20961: {248.8, 268.8}, 20967: {309.8, 335.8}, 20968: {339, 373},
	// Hammer of Wrath
	24275: {285, 315}, 24274: {382, 421}, 24239: {473, 523},
	// Exorcism
	879: {79, 91}, 5614: {140, 158}, 5615: {199, 225}, 10312: {285, 321}, 10313: {376, 420}, 10314: {474, 530},
	// Holy Wrath
	2812: {368.4, 434.4}, 10318: {490, 576},
	// Holy Shock's damage spells (the hand table in holy_shock.go holds the centres)
	1311604: {128, 140}, 25912: {175, 189}, 25911: {248, 268}, 25902: {334, 362},
}

func directDamage(sim *core.Simulation, row shared.SpellData) float64 {
	if r, ok := damageRanges[row.SpellID]; ok {
		return sim.Roll(r[0], r[1])
	}
	return row.Direct.Damage(sim)
}

// The effect at the client's EffectIndex, for the effects Effect(aura, misc) cannot name: the
// weapon-damage effects carry no aura.
func effectAt(row shared.SpellData, index int32) shared.SpellDataEffect {
	for _, e := range row.Effects {
		if e.Index == index {
			return e
		}
	}
	panic("spell has no effect at the index")
}

func NewPaladin(character *core.Character, talentsStr string, _ *proto.PaladinOptions) *Paladin {
	paladin := &Paladin{
		Character: *character,
		Talents:   &proto.PaladinTalents{},
	}

	core.FillTalentsProto(paladin.Talents.ProtoReflect(), talentsStr, TalentTreeSizes)

	// The attack table already holds the base 5% parry and block (sim/core/target.go); adding them
	// here too gave the paladin 10% of each. Base dodge 0.7% and 1% dodge per 19.8 Agility are the
	// Classic paladin's, the same as its crit, as on master.
	paladin.PseudoStats.CanParry = true
	paladin.PseudoStats.BaseDodgeChance += 0.007

	paladin.EnableManaBar()

	paladin.EnableAutoAttacks(paladin, core.AutoAttackOptions{
		MainHand:       paladin.WeaponFromMainHand(),
		AutoSwingMelee: true,
	})

	paladin.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	paladin.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[character.Class])
	paladin.AddStatDependency(stats.Agility, stats.DodgeRating, core.CritPerAgiMaxLevel[character.Class]*core.DodgeRatingPerDodgePercent)
	paladin.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	return paladin
}
