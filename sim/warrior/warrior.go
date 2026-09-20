package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

var TalentTreeSizes = [3]int{17, 18, 18}

type WarriorInputs struct {
	DefaultShout  proto.WarriorShout
	DefaultStance proto.WarriorStance

	StartingRage          float64
	QueueDelay            int32
	StanceSnapshot        bool
	HasBsSolarianSapphire bool
	HasBsT2               bool
}

const (
	SpellFlagBleed = core.SpellFlagAgentReserved1
)

const (
	SpellMaskNone int64 = 0
	// Abilities that don't cost rage and aren't attacks
	SpellMaskBattleShout int64 = 1 << iota
	SpellMaskBerserkerRage
	SpellMaskRecklessness
	SpellMaskDeathWish
	SpellMaskRetaliation
	SpellMaskRetaliationHit
	SpellMaskShieldWall
	SpellMaskLastStand
	SpellMaskCharge
	SpellMaskIntercept
	SpellMaskDemoralizingShout

	// Stances
	SpellMaskBattleStance
	SpellMaskBerserkerStance
	SpellMaskDefensiveStance

	// Special attacks
	SpellMaskRend
	SpellMaskDeepWounds
	SpellMaskSweepingStrikes
	SpellMaskSweepingStrikesHit
	SpellMaskSweepingStrikesNormalizedHit
	SpellMaskHeroicStrike
	SpellMaskCleave
	SpellMaskExecute
	SpellMaskOverpower
	SpellMaskRevenge
	SpellMaskSlam
	SpellMaskSunderArmor
	SpellMaskThunderClap
	SpellMaskWhirlwind
	SpellMaskWhirlwindOh
	SpellMaskShieldSlam
	SpellMaskConcussionBlow
	SpellMaskShieldBash
	SpellMaskBloodthirst
	SpellMaskMortalStrike
	SpellMaskShieldBlock
	SpellMaskHamstring
	SpellMaskPummel
	SpellMaskMockingBlow
	SpellMaskChallengingShout
	SpellMaskIntimidatingShout
	SpellMaskDisarm
	SpellMaskTaunt
	SpellMaskVictoryRush

	WarriorSpellLast
	WarriorSpellsAll = WarriorSpellLast<<1 - 1

	SpellMaskDirectDamageSpells = SpellMaskSweepingStrikesHit | SpellMaskSweepingStrikesNormalizedHit |
		SpellMaskCleave | SpellMaskExecute | SpellMaskHeroicStrike | SpellMaskOverpower |
		SpellMaskRevenge | SpellMaskSlam | SpellMaskShieldBash | SpellMaskSunderArmor |
		SpellMaskThunderClap | SpellMaskWhirlwind | SpellMaskWhirlwindOh | SpellMaskShieldSlam |
		SpellMaskBloodthirst | SpellMaskMortalStrike | SpellMaskIntercept | SpellMaskRetaliationHit |
		SpellMaskMockingBlow | SpellMaskVictoryRush

	SpellMaskDamageSpells = SpellMaskDirectDamageSpells | SpellMaskDeepWounds | SpellMaskRend
)

const EnrageTag = "EnrageEffect"

type Warrior struct {
	core.Character

	ClassSpellScaling float64

	Talents *proto.WarriorTalents

	WarriorInputs

	// Current state
	Stance                Stance
	ChargeRageGain        float64
	BerserkerRageRageGain float64

	BattleShout       *core.Spell
	DemoralizingShout *core.Spell
	BattleStance      *core.Spell
	DefensiveStance   *core.Spell
	BerserkerStance   *core.Spell

	Rend                            *core.Spell
	DeepWounds                      *core.Spell
	MortalStrike                    *core.Spell
	SweepingStrikesNormalizedAttack *core.Spell

	HeroicStrike       *core.Spell
	Cleave             *core.Spell
	MockingBlow        *core.Spell
	ChallengingShout   *core.Spell
	IntimidatingShout  *core.Spell
	Disarm             *core.Spell
	Taunt              *core.Spell
	VictoryRush        *core.Spell
	curQueueAura       *core.Aura
	curQueuedAutoSpell *core.Spell

	sharedShoutsCD   *core.Timer
	queuedRealismICD *core.Cooldown

	EnrageAura *core.Aura

	SweepingStrikesAura *core.Aura

	DemoralizingShoutAuras core.AuraArray
	SunderArmorAuras       core.AuraArray
}

func (warrior *Warrior) GetCharacter() *core.Character {
	return &warrior.Character
}

func (warrior *Warrior) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
}

func (warrior *Warrior) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (warrior *Warrior) Initialize() {
	warrior.registerRecklessness()
	warrior.registerShieldWall()
	warrior.registerRetaliation()

	warrior.registerBerserkerRage()
	warrior.registerBloodrage()
	warrior.registerCharge()
	warrior.registerIntercept()
	warrior.registerPummel()
	warrior.registerHamstring()

	warrior.registerRend()
	warrior.registerSunderArmor()
	warrior.registerHeroicStrike()
	warrior.registerCleave()
	warrior.registerOverpower()
	warrior.registerSlam()
	warrior.registerWhirlwind()
	warrior.registerExecute()
	warrior.registerThunderClap()
	warrior.registerRevenge()
	warrior.registerShieldBlock()
	warrior.registerShieldBash()

	warrior.registerStances()
	warrior.registerShouts()
	warrior.registerForeverAbilities()
}

func (warrior *Warrior) Reset(_ *core.Simulation) {
	warrior.curQueueAura = nil
	warrior.curQueuedAutoSpell = nil

	warrior.ChargeRageGain = 15
	warrior.BerserkerRageRageGain = 0

	switch warrior.DefaultStance {
	case proto.WarriorStance_WarriorStanceBattle:
		warrior.Stance = BattleStance
	case proto.WarriorStance_WarriorStanceDefensive:
		warrior.Stance = DefensiveStance
	case proto.WarriorStance_WarriorStanceBerserker:
		warrior.Stance = BerserkerStance
	}
}

func (warrior *Warrior) OnEncounterStart(sim *core.Simulation) {}

func (warrior *Warrior) GetMainHandType() proto.HandType {
	mh := warrior.GetMHWeapon()

	if mh != nil && (mh.HandType == proto.HandType_HandTypeTwoHand) {
		return proto.HandType_HandTypeTwoHand
	}

	return proto.HandType_HandTypeOneHand
}

func NewWarrior(character *core.Character, options *proto.WarriorOptions, talents string, inputs WarriorInputs) *Warrior {
	warrior := &Warrior{
		Character:     *character,
		Talents:       &proto.WarriorTalents{},
		WarriorInputs: inputs,
	}
	core.FillTalentsProto(warrior.Talents.ProtoReflect(), talents, TalentTreeSizes)

	warrior.EnableRageBar(core.RageBarOptions{
		// Boundless Rage (1310236) raises the cap by 10 per rank.
		MaxRage:            100 + spellData.BoundlessRage.ValueAt(warrior.Talents.BoundlessRage)/10,
		BaseRageMultiplier: 1,
		StartingRage:       inputs.StartingRage,
	})

	warrior.EnableAutoAttacks(warrior, core.AutoAttackOptions{
		MainHand:       warrior.WeaponFromMainHand(),
		OffHand:        warrior.WeaponFromOffHand(),
		AutoSwingMelee: true,
		ReplaceMHSwing: warrior.TryHSOrCleave,
	})

	warrior.PseudoStats.CanParry = true
	warrior.PseudoStats.BaseDodgeChance += 0.0075
	warrior.PseudoStats.BaseParryChance += 0.05
	warrior.PseudoStats.BaseBlockChance += 0.05

	warrior.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	warrior.AddStatDependency(stats.Strength, stats.BlockValue, 1/20.0)
	warrior.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[character.Class])
	warrior.AddStatDependency(stats.Agility, stats.DodgeRating, 1/30.0*core.DodgeRatingPerDodgePercent)
	warrior.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	warrior.sharedShoutsCD = warrior.NewTimer()
	warrior.ChargeRageGain = 15
	warrior.BerserkerRageRageGain = 0
	// The sim often re-enables heroic strike in an unrealistic amount of time.
	// This can cause an unrealistic immediate double-hit around wild strikes procs
	warrior.queuedRealismICD = &core.Cooldown{
		Timer:    warrior.NewTimer(),
		Duration: time.Millisecond * time.Duration(warrior.WarriorInputs.QueueDelay),
	}

	return warrior
}

func (warrior *Warrior) CastNormalizedSweepingStrikesAttack(results core.SpellResultSlice, sim *core.Simulation) {
	if warrior.SweepingStrikesAura != nil && warrior.SweepingStrikesAura.IsActive() {
		for _, result := range results {
			if result.Landed() {
				warrior.SweepingStrikesNormalizedAttack.Cast(sim, warrior.Env.NextActiveTargetUnit(result.Target))
				warrior.SweepingStrikesAura.RemoveStack(sim)
				break
			}
		}
	}
}

// Agent is a generic way to access underlying warrior on any of the agents.
type WarriorAgent interface {
	GetWarrior() *Warrior
}
