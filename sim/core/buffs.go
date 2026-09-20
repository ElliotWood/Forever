package core

import (
	"fmt"
	"math"
	"slices"
	"time"

	googleProto "google.golang.org/protobuf/proto"

	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

// Exclusive categories for resistance buffs — ensures only the highest value applies per school.
const (
	ResistanceCategoryArcane = "ResistanceArcane"
	ResistanceCategoryFire   = "ResistanceFire"
	ResistanceCategoryFrost  = "ResistanceFrost"
	ResistanceCategoryNature = "ResistanceNature"
	ResistanceCategoryShadow = "ResistanceShadow"
)

// Exclusive category for flat stat buffs which don't stack with each other
// e.g. Arcane Brilliance vs Scroll of Intellect.
const StatBuffCategory = "StatBuff"

type BuffConfig struct {
	Label             string
	ActionID          ActionID
	Duration          time.Duration
	Stats             []StatConfig
	ExclusiveCategory string
}

type StatConfig struct {
	Stat             stats.Stat
	Amount           float64
	IsMultiplicative bool
}

func makeMultiplierBuff(aura *Aura, stat stats.Stat, value float64) {
	dep := aura.Unit.NewDynamicMultiplyStat(stat, value)
	aura.ApplyOnGain(func(aura *Aura, sim *Simulation) {
		aura.Unit.EnableBuildPhaseStatDep(sim, dep)
	}).ApplyOnExpire(func(aura *Aura, sim *Simulation) {
		aura.Unit.DisableBuildPhaseStatDep(sim, dep)
	})
}

func makeFlatStatBuff(aura *Aura, stat stats.Stat, value float64) {
	aura.ApplyOnGain(func(aura *Aura, sim *Simulation) {
		aura.Unit.AddStatDynamic(sim, stat, value)
	}).ApplyOnExpire(func(aura *Aura, sim *Simulation) {
		aura.Unit.AddStatDynamic(sim, stat, -value)
	})
}

func registerStatEffect(aura *Aura, config []StatConfig) {
	for _, statConfig := range config {
		if statConfig.IsMultiplicative {
			makeMultiplierBuff(aura, statConfig.Stat, statConfig.Amount)
		} else {
			makeFlatStatBuff(aura, statConfig.Stat, statConfig.Amount)
		}
	}
}

func makeExclusiveMultiplierBuff(aura *Aura, stat stats.Stat, value float64, exclusiveCategory string) {
	dep := aura.Unit.NewDynamicMultiplyStat(stat, value)
	aura.NewExclusiveEffect(exclusiveCategory+stat.StatName()+"Mul", false, ExclusiveEffect{
		Priority: value,
		OnGain: func(ee *ExclusiveEffect, s *Simulation) {
			ee.Aura.Unit.EnableBuildPhaseStatDep(s, dep)
		},
		OnExpire: func(ee *ExclusiveEffect, s *Simulation) {
			ee.Aura.Unit.DisableBuildPhaseStatDep(s, dep)
		},
	})
}

func makeExclusiveFlatStatBuff(aura *Aura, stat stats.Stat, value float64, exclusiveCategory string) {
	aura.NewExclusiveEffect(exclusiveCategory+stat.StatName()+"Add", false, ExclusiveEffect{
		Priority: value,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stat, value)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stat, -value)
		},
	})
}

func registerExlusiveEffects(aura *Aura, config []StatConfig, exclusiveCategory string) {
	for _, statConfig := range config {
		if statConfig.IsMultiplicative {
			makeExclusiveMultiplierBuff(aura, statConfig.Stat, statConfig.Amount, exclusiveCategory)
		} else {
			makeExclusiveFlatStatBuff(aura, statConfig.Stat, statConfig.Amount, exclusiveCategory)
		}
	}
}

func makeStatBuff(char *Character, config BuffConfig) *Aura {
	if config.Label == "" {
		panic("Buff without label.")
	}

	if ActionID.IsEmptyAction(config.ActionID) {
		panic("Buff without ActionID")
	}

	if config.ActionID.Tag == 0 {
		config.ActionID = config.ActionID.WithTag(-1)
	}

	baseAura := char.GetOrRegisterAura(Aura{
		Label:      config.Label,
		ActionID:   config.ActionID,
		Duration:   TernaryDuration(config.Duration > 0, config.Duration, NeverExpires),
		BuildPhase: Ternary(config.ActionID.Tag == -1, CharacterBuildPhaseBuffs, CharacterBuildPhaseNone),
	})

	if config.ExclusiveCategory != "" {
		registerExlusiveEffects(baseAura, config.Stats, config.ExclusiveCategory)
	} else {
		registerStatEffect(baseAura, config.Stats)
	}
	return baseAura
}

// Applies buffs that affect individual players.
func applyBuffEffects(agent Agent, raidBuffs *proto.RaidBuffs, partyBuffs *proto.PartyBuffs, individual *proto.IndividualBuffs) {
	char := agent.GetCharacter()

	applyGeneratedBuffs(char, raidBuffs, partyBuffs, individual)

	// Raid Buffs
	if raidBuffs.Thorns {
		MakePermanent(ThornsAura(char, 0))
	}

	if raidBuffs.Bloodlust {
		registerBloodlustCD(char)
	}

	// Party Buffs
	if partyBuffs.AtieshDruid > 0 {
		MakePermanent(AtieshAura(char, proto.Class_ClassDruid, float64(partyBuffs.AtieshDruid)))
	}

	if partyBuffs.AtieshMage > 0 {
		MakePermanent(AtieshAura(char, proto.Class_ClassMage, float64(partyBuffs.AtieshMage)))
	}

	if partyBuffs.AtieshPriest > 0 {
		MakePermanent(AtieshAura(char, proto.Class_ClassPriest, float64(partyBuffs.AtieshPriest)))
	}

	if partyBuffs.AtieshWarlock > 0 {
		MakePermanent(AtieshAura(char, proto.Class_ClassWarlock, float64(partyBuffs.AtieshWarlock)))
	}

	if partyBuffs.BraidedEterniumChain {
		MakePermanent(BraidedEterniumChainAura(char))
	}

	if partyBuffs.ChainOfTheTwilightOwl {
		MakePermanent(ChainOfTheTwilightOwlAura(char))
	}

	if partyBuffs.EyeOfTheNight {
		MakePermanent(EyeOfTheNightAura(char))
	}

	if partyBuffs.JadePendantOfBlasting {
		MakePermanent(JadePendantOfBlastingAura(char))
	}

	if partyBuffs.DraeneiRacialCaster {
		DraneiRacialAura(char, true)
	}

	if partyBuffs.DraeneiRacialMelee {
		DraneiRacialAura(char, false)
	}

	if partyBuffs.FerociousInspiration > 0 {
		MakePermanent(FerociousInspiration(char, partyBuffs.FerociousInspiration))
	}

	if partyBuffs.GraceOfAirTotem {
		GraceOfAirTotemAura(char, false, partyBuffs.TotemTwisting)
	}

	if partyBuffs.RetributionAura {
		MakePermanent(RetributionAuraBuff(char, false, 0))
	}

	if partyBuffs.TotemOfWrath > 0 {
		MakePermanent(TotemOfWrathAura(char, partyBuffs.TotemOfWrath))
	}

	if partyBuffs.TranquilAirTotem {
		MakePermanent(TranquilAirTotemAura(char))
	}

	if partyBuffs.TrueshotAura {
		MakePermanent(TrueShotAuraBuff(char))
	}

	if partyBuffs.WindfuryTotem {
		WindfuryTotemAura(char, false)
	}

	if partyBuffs.WrathOfAirTotem != proto.TristateEffect_TristateEffectMissing {
		MakePermanent(WrathOfAirTotemAura(char, IsImproved(partyBuffs.WrathOfAirTotem)))
	}
	if partyBuffs.Drums > 0 {
		DrumsBuff(char, partyBuffs.Drums)
	}

	// Individual Buffs
	if individual.BlessingOfSanctuary {
		MakePermanent(BlessingOfSanctuaryAura(char))
	}

	if individual.ShadowPriestDps > 0 {
		MakePermanent(ShadowPriestDPSManaAura(char, float64(individual.ShadowPriestDps)))
	}

	if individual.UnleashedRage {
		MakePermanent(UnleashedRageAura(char, -1, 5))
	}

}

///////////////////////////////////////////////////////////////////////////
//							Raid Buffs
///////////////////////////////////////////////////////////////////////////

func ThornsAura(char *Character, points int32) *Aura {
	actionID := ActionID{SpellID: 26992}

	procSpell := char.RegisterSpell(SpellConfig{
		ActionID:    actionID,
		SpellSchool: SpellSchoolNature,
		Flags:       SpellFlagBinary | SpellFlagPassiveSpell,
		ProcMask:    ProcMaskEmpty,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			baseDamage := 25 * (1 + 0.25*float64(points))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
		},
	})

	return char.MakeProcTriggerAura(ProcTrigger{
		Name:     "Thorns",
		ActionID: actionID,
		Duration: time.Minute * 10,
		Outcome:  OutcomeLanded,
		Callback: CallbackOnSpellHitTaken,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			if spell.SpellSchool.Matches(SpellSchoolPhysical) {
				procSpell.Cast(sim, spell.Unit)
			}
		},
	})
}

// /////////////////////////////////////////////////////////////////////////
//
//	Party Buffs
//
// /////////////////////////////////////////////////////////////////////////

func ApplyFixedShoutAura(char *Character, aura *Aura, category string) {
	aura.ApplyOnInit(func(aura *Aura, sim *Simulation) {
		auras := char.GetAurasWithTag(category)
		var playerAura *Aura
		for _, bsAura := range auras {
			if bsAura.ActionID.Tag == 0 {
				playerAura = bsAura
				break
			}
		}

		if playerAura == nil {
			return
		}

		playerAura.ApplyOnGain(func(_ *Aura, _ *Simulation) {
			pa := sim.GetConsumedPendingActionFromPool()
			pa.NextActionAt = playerAura.ExpiresAt() + char.ReactionTime
			pa.OnAction = func(sim *Simulation) {
				aura.Activate(sim)

				StartPeriodicAction(sim, PeriodicActionOptions{
					Period:   aura.Duration + 1,
					NumTicks: 1,
					OnAction: func(sim *Simulation) {
						aura.Activate(sim)
					},
				})
			}
			sim.AddPendingAction(pa)
		})
	})

	ApplyFixedUptimeAura(aura, 1, aura.Duration+1, -1)
}

var (
	PaladinAuraCategory     = "PaladinAura"
	RetributionAuraCategory = "RetributionAura"
)

// paladinAuraPriority returns the exclusivity priority used for self/external
// variants of the same paladin aura. Self-cast always wins over party-buff-
// applied (external) versions so they never stack.
func paladinAuraPriority(isPlayer bool) float64 {
	return TernaryFloat64(isPlayer, 1, 0)
}

func FerociousInspiration(char *Character, count int32) *Aura {
	dmgBuff := 0.03 * float64(count)

	return char.GetOrRegisterAura(Aura{
		Label:    "Ferocious Inspiration",
		ActionID: ActionID{SpellID: 34460},
		Duration: time.Second * 10,
	}).AttachMultiplicativePseudoStatBuff(&char.PseudoStats.DamageDealtMultiplier, 1+dmgBuff)
}

func RetributionAuraBuff(char *Character, isPlayer bool, impRetributionAuraRank int32) *Aura {
	actionID := ActionID{SpellID: 27150}.WithTag(TernaryInt32(isPlayer, 0, -1))
	impMultiplier := 1 + 0.25*float64(impRetributionAuraRank)

	procSpell := char.RegisterSpell(SpellConfig{
		ActionID:    actionID.WithTag(2),
		SpellSchool: SpellSchoolHoly,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagBinary | SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			spell.CalcAndDealDamage(sim, target, 26*impMultiplier, spell.OutcomeAlwaysHit)
		},
	})

	aura := char.GetOrRegisterAura(Aura{
		Label:      fmt.Sprintf("Retribution Aura (%s)", Ternary(isPlayer, "Player", "External")),
		ActionID:   actionID,
		Duration:   NeverExpires,
		BuildPhase: Ternary(isPlayer, CharacterBuildPhaseNone, CharacterBuildPhaseBuffs),
	}).AttachProcTrigger(ProcTrigger{
		Name:     "Retribution Aura Damage",
		Callback: CallbackOnSpellHitTaken,
		Outcome:  OutcomeLanded,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			if spell.SpellSchool.Matches(SpellSchoolPhysical) {
				procSpell.Cast(sim, spell.Unit)
			}
		},
	})

	// Self and external share RetributionAuraCategory (SingleAura) so only one
	// variant's proc trigger fires at a time; self wins via higher priority.
	aura.NewExclusiveEffect(RetributionAuraCategory, true, ExclusiveEffect{
		Priority: paladinAuraPriority(isPlayer),
	})

	if isPlayer {
		aura.NewExclusiveEffect(PaladinAuraCategory, true, ExclusiveEffect{})
	}

	return aura
}

func TrueShotAuraBuff(char *Character) *Aura {
	apBuff := 125.0

	return makeStatBuff(char, BuffConfig{
		Label:    "Trueshot Aura",
		ActionID: ActionID{SpellID: 27066},
		Stats: []StatConfig{
			{stats.RangedAttackPower, apBuff, false},
			{stats.AttackPower, apBuff, false},
		},
	})
}

var UnleashedRageCategory = "UnleashedRage"

func UnleashedRageAura(char *Character, casterIdx int32, points int32) *Aura {
	return makeStatBuff(char, BuffConfig{
		Label:    fmt.Sprintf("Unleashed Rage-%d", casterIdx),
		Duration: time.Second * 10,
		ActionID: ActionID{SpellID: 30809}.WithTag(casterIdx),
		Stats: []StatConfig{
			{stats.AttackPower, 1 + 0.02*float64(points), true},
		},
		ExclusiveCategory: UnleashedRageCategory,
	})
}

// //////////////////////////
//
//	Totems
//
// //////////////////////////
var GraceOfAirTotemCategory = "GraceOfAirTotem"

func GraceOfAirTotemAura(char *Character, improved bool, wfActive bool) *Aura {
	agiBuff := 77.0
	if improved {
		// Truncated like the game: 77*1.15=88.55 -> 88.
		agiBuff = math.Floor(agiBuff * 1.15)
	}

	duration := NeverExpires
	if wfActive {
		duration = time.Second * 9
	}

	return makeStatBuff(char, BuffConfig{
		Label:    "Grace of Air Totem",
		ActionID: ActionID{SpellID: 25359},
		Stats: []StatConfig{
			{stats.Agility, agiBuff, false},
		},
		Duration:          duration,
		ExclusiveCategory: GraceOfAirTotemCategory,
	}).ApplyOnReset(func(aura *Aura, sim *Simulation) {
		if wfActive {
			StartPeriodicAction(sim, PeriodicActionOptions{
				Period:   time.Second * 10,
				Priority: ActionPriorityAuto,
				OnAction: func(sim *Simulation) {
					aura.Activate(sim)
				},
			})
		} else {
			aura.Activate(sim)
		}
	})
}

func TotemOfWrathAura(char *Character, count int32) *Aura {
	modValue := 3.0 * float64(count)

	return makeStatBuff(char, BuffConfig{
		Label:    "Totem of Wrath",
		ActionID: ActionID{SpellID: 30706},
		Stats: []StatConfig{
			{stats.SpellCritPercent, modValue, false},
			{stats.SpellHitPercent, modValue, false},
		},
	})
}

func TranquilAirTotemAura(char *Character) *Aura {
	return char.GetOrRegisterAura(Aura{
		Label:    "Tranquil Air Totem",
		ActionID: ActionID{SpellID: 25909},
	}).AttachMultiplicativePseudoStatBuff(&char.PseudoStats.ThreatMultiplier, 0.8)
}

var WindfuryTotemCategory = "WindfuryTotem"

func WindfuryTotemAura(char *Character, isImpoved bool) *Aura {
	apBonus := 445.0
	if isImpoved {
		// Truncated like the game: 445*1.3=578.5 -> 578.
		apBonus = math.Floor(apBonus * 1.3)
	}
	// Chance on MH Auto Attack to instantly attack with another AA with apBonus.
	// AP bonus lingers until 2 auto attacks are performed.
	// If procced from a normal auto, this consumes the buff almost instantly (server tick rate applies)
	// If procced from a "Next Auto" special (HS/Cleave), this can result in the aura lasting the entire 1.5s duration.

	wfProcAura := char.NewTemporaryStatsAura("Windfury Totem Proc", ActionID{SpellID: 25584}, stats.Stats{stats.AttackPower: apBonus}, time.Millisecond*1500)
	wfProcAura.MaxStacks = 2
	wfProcAura.AttachProcTrigger(ProcTrigger{
		Name:     "Windfury Attack",
		Callback: CallbackOnSpellHitDealt,
		ProcMask: ProcMaskMeleeMHAuto | ProcMaskMeleeOHAuto,
		// TriggerImmediately ommited for improved UI clarity (the timeline tick would be near invisible for MHAuto procs)
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			if wfProcAura.IsActive() && !spell.ProcMask.Matches(ProcMaskMeleeSpecial) {
				wfProcAura.RemoveStack(sim)
			}
		},
	})

	var windfurySpell *Spell
	wfProcTrigger := char.MakeProcTriggerAura(ProcTrigger{
		Name:               "Windfury Totem Trigger",
		MetricsActionID:    ActionID{SpellID: 25580, Tag: -1},
		IsWeaponProc:       true,
		ProcChance:         0.2,
		Duration:           NeverExpires,
		Outcome:            OutcomeLanded,
		Callback:           CallbackOnSpellHitDealt,
		ProcMask:           ProcMaskMeleeMHAuto,
		ICD:                time.Millisecond * 1500,
		TriggerImmediately: true,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			wfProcAura.Activate(sim)
			if spell.ProcMask == ProcMaskMeleeMHAuto {
				wfProcAura.SetStacks(sim, 1)
			} else {
				wfProcAura.SetStacks(sim, 2)
			}
			char.AutoAttacks.MaybeReplaceMHSwing(sim, windfurySpell).Cast(sim, result.Target)
		},
	})

	wfAura := char.GetOrRegisterAura(Aura{
		Label:    "Windfury Totem",
		ActionID: ActionID{SpellID: 25587, Tag: -1},
		Duration: time.Second * 10,
	}).ApplyOnInit(func(aura *Aura, sim *Simulation) {
		config := *char.AutoAttacks.MHConfig()
		config.ActionID = config.ActionID.WithTag(25584)
		windfurySpell = char.GetOrRegisterSpell(config)
	}).ApplyOnReset(func(aura *Aura, sim *Simulation) {
		aura.Activate(sim)
		StartPeriodicAction(sim, PeriodicActionOptions{
			Period:   time.Second * 5,
			Priority: ActionPriorityAuto,
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
			},
		})
	})

	wfAura.NewExclusiveEffect(WindfuryTotemCategory, false, ExclusiveEffect{
		Priority: apBonus,
		OnGain: func(_ *ExclusiveEffect, sim *Simulation) {
			wfProcTrigger.Activate(sim)
		},
		OnExpire: func(_ *ExclusiveEffect, sim *Simulation) {
			wfProcTrigger.Deactivate(sim)
			wfAura.Deactivate(sim)
		},
	})

	char.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}, func(sim *Simulation, slot proto.ItemSlot) {
		wfAura.Deactivate(sim)
	})

	return wfProcTrigger
}

const (
	WrathOfAirTotemCategory      = "WrathOfAirTotem"
	WrathOfAirTotemBaseValue     = 101.0
	WrathOfAirTotemImprovedValue = 20.0
)

func WrathOfAirTotemValue(improved bool) float64 {
	return WrathOfAirTotemBaseValue + TernaryFloat64(improved, WrathOfAirTotemImprovedValue, 0)
}

func WrathOfAirTotemAura(char *Character, improved bool) *Aura {
	buff := WrathOfAirTotemValue(improved)

	return makeStatBuff(char, BuffConfig{
		Label:    "Wrath of Air Totem",
		ActionID: ActionID{SpellID: 3738},
		Stats: []StatConfig{
			{stats.SpellDamage, buff, false},
			{stats.HealingPower, buff, false},
		},
		ExclusiveCategory: WrathOfAirTotemCategory,
	})
}

////////////////////////////
//	Item Buffs
////////////////////////////

func AtieshAura(char *Character, class proto.Class, numStaves float64) *Aura {
	switch class {
	case proto.Class_ClassDruid:
		return makeStatBuff(char, BuffConfig{
			Label:    "Power of the Guardian - Druid",
			ActionID: ActionID{SpellID: 28145},
			Stats: []StatConfig{
				{stats.MP5, 11 * numStaves, false},
			},
		})
	case proto.Class_ClassMage:
		return makeStatBuff(char, BuffConfig{
			Label:    "Power of the Guardian - Mage",
			ActionID: ActionID{SpellID: 28142},
			Stats: []StatConfig{
				{stats.SpellCritRating, 28 * numStaves, false},
			},
		})
	case proto.Class_ClassPriest:
		return makeStatBuff(char, BuffConfig{
			Label:    "Power of the Guardian - Priest",
			ActionID: ActionID{SpellID: 28144},
			Stats: []StatConfig{
				{stats.HealingPower, 62 * numStaves, false},
			},
		})
	default: // Use warlock as default to satisfy compiler
		return makeStatBuff(char, BuffConfig{
			Label:    "Power of the Guardian - Warlock",
			ActionID: ActionID{SpellID: 28143},
			Stats: []StatConfig{
				{stats.SpellDamage, 33 * numStaves, false},
				{stats.HealingPower, 33 * numStaves, false},
			},
		})
	}

}

const (
	BraidedEterniumChainAuraLabel  = "Braided Eternium Chain"
	ChainOfTheTwilightOwlAuraLabel = "Chain of the Twilight Owl"
	EyeOfTheNightAuraLabel         = "Eye of the Night"
	JadePendantOfBlastingAuraLabel = "Jade Pendant of Blasting"
)

func BraidedEterniumChainAura(char *Character) *Aura {
	return makeStatBuff(char, BuffConfig{
		Label:             BraidedEterniumChainAuraLabel,
		ActionID:          ActionID{SpellID: 31025},
		ExclusiveCategory: BraidedEterniumChainAuraLabel,
		Stats: []StatConfig{
			{stats.MeleeCritRating, 28, false},
		},
	})
}

func ChainOfTheTwilightOwlAura(char *Character) *Aura {
	return makeStatBuff(char, BuffConfig{
		Label:             ChainOfTheTwilightOwlAuraLabel,
		ActionID:          ActionID{SpellID: 31035},
		ExclusiveCategory: ChainOfTheTwilightOwlAuraLabel,
		Stats: []StatConfig{
			{stats.SpellCritPercent, 2, false},
		},
	})
}

func EyeOfTheNightAura(char *Character) *Aura {
	return makeStatBuff(char, BuffConfig{
		Label:             EyeOfTheNightAuraLabel,
		ActionID:          ActionID{SpellID: 31033},
		ExclusiveCategory: EyeOfTheNightAuraLabel,
		Stats: []StatConfig{
			{stats.SpellDamage, 34, false},
		},
	})
}

func JadePendantOfBlastingAura(char *Character) *Aura {
	return makeStatBuff(char, BuffConfig{
		Label:             JadePendantOfBlastingAuraLabel,
		ActionID:          ActionID{SpellID: 25607},
		ExclusiveCategory: JadePendantOfBlastingAuraLabel,
		Stats: []StatConfig{
			{stats.SpellDamage, 15, false},
		},
	})
}

func DraneiRacialAura(char *Character, caster bool) *Aura {
	alliance := []proto.Race{
		proto.Race_RaceDraenei,
		proto.Race_RaceDwarf,
		proto.Race_RaceGnome,
		proto.Race_RaceHuman,
		proto.Race_RaceNightElf,
	}
	if !slices.Contains(alliance, char.Race) {
		return nil
	}
	var aura *Aura
	if caster {
		aura = makeStatBuff(char, BuffConfig{
			Label:    "Inspiring Presence",
			ActionID: ActionID{SpellID: 28878},
			Stats: []StatConfig{
				{stats.SpellHitPercent, 1, false},
			},
			ExclusiveCategory: "Inspiring Presence",
		})
	} else {
		aura = makeStatBuff(char, BuffConfig{
			Label:    "Heroic Presence",
			ActionID: ActionID{SpellID: 6562},
			Stats: []StatConfig{
				{stats.PhysicalHitPercent, 1, false},
			},
			ExclusiveCategory: "Heroic Presence",
		})
	}

	return MakePermanent(aura)
}

const TinnitusAuraLabel = "Tinnitus"

func drumsSpellConfig(character *Character, drum proto.Drums, isExternal bool) SpellConfig {
	var drumLabel string
	var drumStats stats.Stats
	var duration time.Duration
	var actionID ActionID
	switch drum {
	case proto.Drums_GreaterDrumsOfBattle, proto.Drums_LesserDrumsOfBattle:
		drumLabel = "Drums of Battle"
		drumStats = stats.Stats{stats.MeleeHasteRating: 80, stats.SpellHasteRating: 80}
		duration = time.Second * 30
		actionID = ActionID{SpellID: 35476}
	case proto.Drums_GreaterDrumsOfWar, proto.Drums_LesserDrumsOfWar:
		drumLabel = "Drums of War"
		drumStats = stats.Stats{stats.AttackPower: 60, stats.RangedAttackPower: 60, stats.SpellDamage: 30}
		duration = time.Second * 30
		actionID = ActionID{SpellID: 35475}
	case proto.Drums_GreaterDrumsOfRestoration, proto.Drums_LesserDrumsOfRestoration:
		drumLabel = "Drums of Restoration"
		drumStats = stats.Stats{stats.MP5: 200}
		duration = time.Second * 15
		actionID = ActionID{SpellID: 35478}
	}

	if isExternal {
		actionID = actionID.WithTag(-1)
		drumLabel = drumLabel + " (External)"
	}

	aura := character.NewTemporaryStatsAura(drumLabel, actionID, drumStats, duration)

	tinnitus := character.GetOrRegisterAura(Aura{
		Label:    TinnitusAuraLabel,
		ActionID: ActionID{SpellID: 369770},
		Duration: time.Minute * 2,
	})

	aura.ApplyOnGain(func(_ *Aura, sim *Simulation) {
		tinnitus.Activate(sim)
	})

	spellConfig := SpellConfig{
		ActionID: actionID,
		Flags:    SpellFlagNoOnCastComplete,
		ProcMask: ProcMaskEmpty,
		ExtraCastCondition: func(sim *Simulation, target *Unit) bool {
			if !character.HasActiveAura(TinnitusAuraLabel) {
				return true
			}
			return false
		},
		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			if !character.HasActiveAura(TinnitusAuraLabel) {
				aura.Activate(sim)
			}
		},

		RelatedSelfBuff: aura.Aura,
	}

	return spellConfig
}

func DrumsBuff(character *Character, drum proto.Drums) {
	config := drumsSpellConfig(character, drum, true)
	config.Cast = CastConfig{
		CD: Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Minute * 2,
		},
	}
	spell := character.RegisterSpell(config)

	character.AddMajorCooldown(MajorCooldown{
		Spell:    spell,
		Type:     CooldownTypeDPS,
		Priority: CooldownPriorityDrums,
	})
}

///////////////////////////////////////////////////////////////////////////
//							Individual Buffs
///////////////////////////////////////////////////////////////////////////

func AmplifyMagicAura(char *Character, improved bool) *Aura {
	baseMod := 120.0
	if improved {
		baseMod *= 1.50
	}
	return char.GetOrRegisterAura(Aura{
		Label:    "Amplify Magic",
		ActionID: ActionID{SpellID: 33946},
		Duration: time.Minute * 10,

		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusHealingTaken += baseMod * 2
			aura.Unit.PseudoStats.BonusPhysicalDamageTaken += baseMod
		},

		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusHealingTaken -= baseMod * 2
			aura.Unit.PseudoStats.BonusPhysicalDamageTaken -= baseMod
		},
	})
}

func DampenMagicAura(char *Character, improved bool) *Aura {
	baseMod := 120.0
	if improved {
		baseMod *= 1.50
	}
	return char.GetOrRegisterAura(Aura{
		Label:    "Amplify Magic",
		ActionID: ActionID{SpellID: 33946},
		Duration: time.Minute * 10,

		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusHealingTaken -= baseMod * 2
			aura.Unit.PseudoStats.BonusSpellDamageTaken -= baseMod
		},

		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.BonusHealingTaken += baseMod * 2
			aura.Unit.PseudoStats.BonusSpellDamageTaken += baseMod
		},
	})
}

// func BlessingOfLight(char *Character) *Aura {
// 	return char.GetOrRegisterAura(Aura{
// 		Label:    "Blessing of Light",
// 		ActionID: ActionID{SpellID: 27145},
// 		Duration: time.Minute * 30,

// 		OnApplyEffects: func(aura *Aura, sim *Simulation, target *Unit, spell *Spell) {
// 			if spell.ProcMask != ProcMaskSpellHealing {
// 				return
// 			}

// 			if spell.Unit.ownerClass != proto.Class_ClassPaladin {
// 				return
// 			}

// 			// Keep an eye on if this changes in paladin.go
// 			// FlashOfLight = 2
// 			// HolyLight = 3
// 			if spell.ClassSpellMask != 2 || spell.ClassSpellMask != 3 {
// 				return
// 			}

// 			if spell.ClassSpellMask == 2 {
// 				spell.BonusSpellDamage += 185
// 			} else {
// 				spell.BonusSpellDamage += 580
// 			}
// 		},
// 	})
// }

func BlessingOfSanctuaryAura(char *Character) *Aura {
	actionID := ActionID{SpellID: 27169}

	procSpell := char.RegisterSpell(SpellConfig{
		ActionID:    actionID,
		SpellSchool: SpellSchoolHoly,
		Flags:       SpellFlagBinary | SpellFlagPassiveSpell,
		ProcMask:    ProcMaskEmpty,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			spell.CalcAndDealDamage(sim, target, 46, spell.OutcomeMagicHit)
		},
	})

	return char.MakeProcTriggerAura(ProcTrigger{
		Name:     "Blessing of Sanctuary",
		ActionID: actionID,
		Duration: time.Minute * 10,
		Outcome:  OutcomeBlock,
		Callback: CallbackOnSpellHitTaken,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			procSpell.Cast(sim, spell.Unit)
		},
	}).AttachMultiplicativePseudoStatBuff(&char.PseudoStats.BonusPhysicalDamageTaken, -80)
}

////////////////////////////
//  Individual Buffs
////////////////////////////

func ShadowPriestDPSManaAura(char *Character, dps float64) *Aura {
	manaMetrics := char.NewManaMetrics(ActionID{SpellID: 34914})

	manaGain := dps * 0.05

	return char.GetOrRegisterAura(Aura{
		Label:    "Vampiric Touch",
		ActionID: ActionID{SpellID: 34914},
		OnGain: func(aura *Aura, sim *Simulation) {
			StartPeriodicAction(sim, PeriodicActionOptions{
				Period: DurationFromSeconds(1),
				OnAction: func(s *Simulation) {
					char.AddMana(sim, manaGain, manaMetrics)
				},
			})
		},
	})
}

////////////////////////////
//  Cooldowns
////////////////////////////

func multiplyCastSpeedEffect(aura *Aura, multiplier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("MultiplyCastSpeed", false, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyCastSpeed(sim, multiplier)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyCastSpeed(sim, 1/multiplier)
		},
	})
}

func InspirationAura(unit *Unit, points int32) *Aura {
	multiplier := 1 + []float64{0, .08, .16, .25}[points]

	armorDep := unit.NewDynamicMultiplyStat(stats.Armor, multiplier)

	return unit.GetOrRegisterAura(Aura{
		Label:    "Inspiration",
		ActionID: ActionID{SpellID: 15363},
		Duration: time.Second * 15,
	}).AttachStatDependency(armorDep)
}

func ApplyInspiration(character *Character, uptime float64) {
	if uptime <= 0 {
		return
	}
	uptime = min(1, uptime)

	inspirationAura := InspirationAura(&character.Unit, 3)

	ApplyFixedUptimeAura(inspirationAura, uptime, time.Millisecond*2500, 1)
}

// Applies buffs to pets.
func applyPetBuffEffects(petAgent PetAgent, raidBuffs *proto.RaidBuffs, partyBuffs *proto.PartyBuffs, individualBuffs *proto.IndividualBuffs) {
	// Summoned pets, like Mage Water Elemental, aren't around to receive raid buffs.
	if petAgent.GetPet().IsGuardian() {
		return
	}

	// We need to modify the buffs a bit because some things are applied to pets by
	// the owner during combat (Bloodlust) or don't make sense for a pet.
	raidBuffs = googleProto.Clone(raidBuffs).(*proto.RaidBuffs)
	raidBuffs.Bloodlust = false
	raidBuffs.Thorns = false

	partyBuffs = googleProto.Clone(partyBuffs).(*proto.PartyBuffs)
	// Pets can't get extra attacks, doh!
	partyBuffs.WindfuryTotem = false
	// Neck auras are automatically inherited when a unit gets into range (40 yds)
	partyBuffs.ChainOfTheTwilightOwl = partyBuffs.ChainOfTheTwilightOwl || petAgent.GetPet().Owner.HasAura(ChainOfTheTwilightOwlAuraLabel)
	partyBuffs.EyeOfTheNight = partyBuffs.EyeOfTheNight || petAgent.GetPet().Owner.HasAura(EyeOfTheNightAuraLabel)
	partyBuffs.BraidedEterniumChain = partyBuffs.BraidedEterniumChain || petAgent.GetPet().Owner.HasAura(BraidedEterniumChainAuraLabel)
	partyBuffs.JadePendantOfBlasting = partyBuffs.JadePendantOfBlasting || petAgent.GetPet().Owner.HasAura(JadePendantOfBlastingAuraLabel)

	individualBuffs = googleProto.Clone(individualBuffs).(*proto.IndividualBuffs)
	individualBuffs.Innervates = 0
	individualBuffs.PowerInfusions = 0

	partyBuffs.Drums = proto.Drums_DrumsUnknown

	if !petAgent.GetPet().enabledOnStart {
		// Auras etc still apply, but not targeted buffs (usually)
		// Strip targeted buffs that require presence at fight start
		raidBuffs.ArcaneBrilliance = false
		raidBuffs.DivineSpirit = false
		raidBuffs.GiftOfTheWild = false
		raidBuffs.PowerWordFortitude = false
		raidBuffs.ShadowProtection = false
		raidBuffs.Thorns = false
		individualBuffs.BlessingOfMight = false
		individualBuffs.BlessingOfKings = false
		individualBuffs.BlessingOfWisdom = false

		// Only individual buff that would apply is Unleashed Rage.
		unleashedRage := individualBuffs.UnleashedRage
		individualBuffs = &proto.IndividualBuffs{}
		individualBuffs.UnleashedRage = unleashedRage
	}

	applyBuffEffects(petAgent, raidBuffs, partyBuffs, individualBuffs)
}

// Used for approximating cooldowns applied by other players to you, such as
// bloodlust, innervate, power infusion, etc. This is specifically for buffs
// which can be consecutively applied multiple times to a single player.
type externalConsecutiveCDApproximation struct {
	ActionID         ActionID
	AuraTag          string
	CooldownPriority int32
	Type             CooldownType
	AuraDuration     time.Duration
	AuraCD           time.Duration

	// Callback for extra activation conditions.
	ShouldActivate CooldownActivationCondition

	// Applies the buff.
	AddAura           CooldownActivation
	RelatedSelfBuff   *Aura             // Used to attach the aura to the generic spell
	RelatedAuraArrays LabeledAuraArrays // Used to attach the aura to the generic spell
}

// numSources is the number of other players assigned to apply the buff to this player.
// E.g. the number of other shaman in the group using bloodlust.
func registerExternalConsecutiveCDApproximation(char *Character, config externalConsecutiveCDApproximation, numSources int32) {
	if numSources == 0 {
		panic("Need at least 1 source!")
	}

	var nextExternalIndex int

	externalTimers := make([]*Timer, numSources)
	for i := 0; i < int(numSources); i++ {
		externalTimers[i] = char.NewTimer()
	}
	sharedTimer := char.NewTimer()

	spell := char.RegisterSpell(SpellConfig{
		ActionID: config.ActionID,
		Flags:    SpellFlagNoOnCastComplete | SpellFlagNoMetrics | SpellFlagNoLogs,

		Cast: CastConfig{
			CD: Cooldown{
				Timer:    sharedTimer,
				Duration: config.AuraDuration, // Assumes that multiple buffs are different sources.
			},
		},
		ExtraCastCondition: func(sim *Simulation, target *Unit) bool {
			if !externalTimers[nextExternalIndex].IsReady(sim) {
				return false
			}

			if char.HasActiveAuraWithTag(config.AuraTag) {
				return false
			}

			return true
		},

		ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
			config.AddAura(sim, char)
			externalTimers[nextExternalIndex].Set(sim.CurrentTime + config.AuraCD)

			nextExternalIndex = (nextExternalIndex + 1) % len(externalTimers)

			if externalTimers[nextExternalIndex].IsReady(sim) {
				sharedTimer.Set(sim.CurrentTime + config.AuraDuration)
			} else {
				sharedTimer.Set(sim.CurrentTime + externalTimers[nextExternalIndex].TimeToReady(sim))
			}
		},
		RelatedSelfBuff:   config.RelatedSelfBuff,
		RelatedAuraArrays: config.RelatedAuraArrays,
	})

	char.AddMajorCooldown(MajorCooldown{
		Spell:    spell,
		Priority: config.CooldownPriority,
		Type:     config.Type,

		ShouldActivate: config.ShouldActivate,
	})
}

var BloodlustActionID = ActionID{SpellID: 2825}

const SatedAuraLabel = "Sated"
const BloodlustAuraTag = "Bloodlust"
const BloodlustDuration = time.Second * 40
const BloodlustCD = time.Minute * 10

func registerBloodlustCD(character *Character) {
	bloodlustAura := BloodlustAura(character, -1)

	spell := character.RegisterSpell(SpellConfig{
		ActionID: bloodlustAura.ActionID,
		Flags:    SpellFlagAPL | SpellFlagNoOnCastComplete | SpellFlagNoMetrics | SpellFlagNoLogs,

		Cast: CastConfig{
			CD: Cooldown{
				Timer:    character.NewTimer(),
				Duration: BloodlustCD,
			},
		},

		ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
			if !character.HasActiveAura(SatedAuraLabel) {
				bloodlustAura.Activate(sim)
			}
		},

		RelatedSelfBuff: bloodlustAura,
	})

	character.AddMajorCooldown(MajorCooldown{
		Spell:    spell,
		Priority: CooldownPriorityBloodlust,
		Type:     CooldownTypeDPS,
		ShouldActivate: func(sim *Simulation, character *Character) bool {
			return !character.HasActiveAura(SatedAuraLabel)
		},
	})
}

func BloodlustAura(character *Character, actionTag int32) *Aura {
	actionID := BloodlustActionID.WithTag(actionTag)

	sated := character.GetOrRegisterAura(Aura{
		Label:    SatedAuraLabel,
		ActionID: ActionID{SpellID: 57724},
		Duration: time.Minute * 10,
	})

	for _, pet := range character.Pets {
		if !pet.IsGuardian() {
			BloodlustAura(&pet.Character, actionTag)
		}
	}

	aura := character.GetOrRegisterAura(Aura{
		Label:    "Bloodlust-" + actionID.String(),
		Tag:      BloodlustAuraTag,
		ActionID: actionID,
		Duration: BloodlustDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1.3)
			for _, pet := range character.Pets {
				if pet.IsEnabled() && !pet.IsGuardian() {
					pet.GetAura(aura.Label).Activate(sim)
				}
			}
			sated.Activate(sim)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1/1.3)
		},
	})
	multiplyCastSpeedEffect(aura, 1.3)
	return aura
}

var PainSuppressionAuraTag = "PainSuppression"

const PainSuppressionDuration = time.Second * 8
const PainSuppressionCD = time.Minute * 3

func registerPainSuppressionCD(char *Character, numPainSuppressions int32) {
	if numPainSuppressions == 0 {
		return
	}

	psAura := PainSuppressionAura(char, -1)

	registerExternalConsecutiveCDApproximation(
		char,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 33206, Tag: -1},
			AuraTag:          PainSuppressionAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			RelatedSelfBuff:  psAura,
			AuraDuration:     PainSuppressionDuration,
			AuraCD:           PainSuppressionCD,
			Type:             CooldownTypeSurvival,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) {
				psAura.Activate(sim)
			},
		},
		numPainSuppressions)
}

func PainSuppressionAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 33206, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "PainSuppression-" + actionID.String(),
		Tag:      PainSuppressionAuraTag,
		ActionID: actionID,
		Duration: PainSuppressionDuration,
	}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.DamageTakenMultiplier, 0.6)
}
