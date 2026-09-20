package core

import (
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

// The caller states the action tag: -1 for a copy the build phase measures.
func makeStatBuff(char *Character, config BuffConfig) *Aura {
	if config.Label == "" {
		panic("Buff without label.")
	}

	if ActionID.IsEmptyAction(config.ActionID) {
		panic("Buff without ActionID")
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

	if raidBuffs.Bloodlust {
		registerBloodlustCD(char)
	}
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
			ActionID: ActionID{SpellID: 28878}.WithTag(-1),
			Stats: []StatConfig{
				{stats.SpellHitPercent, 1, false},
			},
			ExclusiveCategory: "Inspiring Presence",
		})
	} else {
		aura = makeStatBuff(char, BuffConfig{
			Label:    "Heroic Presence",
			ActionID: ActionID{SpellID: 6562}.WithTag(-1),
			Stats: []StatConfig{
				{stats.PhysicalHitPercent, 1, false},
			},
			ExclusiveCategory: "Heroic Presence",
		})
	}

	return MakePermanent(aura)
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

// Applies buffs to pets. Which of the owner's buffs a pet is given is each
// manifest row's pet policy, so the copies below are only there to keep
// applyGeneratedPetBuffs from writing into the raid's own configuration.
func applyPetBuffEffects(petAgent PetAgent, raidBuffs *proto.RaidBuffs, partyBuffs *proto.PartyBuffs, individualBuffs *proto.IndividualBuffs) {
	// Summoned pets, like Mage Water Elemental, aren't around to receive raid buffs.
	if petAgent.GetPet().IsGuardian() {
		return
	}

	raidBuffs = googleProto.Clone(raidBuffs).(*proto.RaidBuffs)
	partyBuffs = googleProto.Clone(partyBuffs).(*proto.PartyBuffs)
	individualBuffs = googleProto.Clone(individualBuffs).(*proto.IndividualBuffs)

	applyGeneratedPetBuffs(petAgent.GetPet(), raidBuffs, partyBuffs, individualBuffs)

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
