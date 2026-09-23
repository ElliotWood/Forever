package core

import (
	"fmt"
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
// e.g. Arcane Brilliance vs Greater Arcane Elixir.
const StatBuffCategory = "StatBuff"

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

// Applies buffs that affect individual players.
func applyBuffEffects(agent Agent, raidBuffs *proto.RaidBuffs, partyBuffs *proto.PartyBuffs, individual *proto.IndividualBuffs) {
	char := agent.GetCharacter()

	applyGeneratedBuffs(char, raidBuffs, partyBuffs, individual)
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

// One rank of a paladin aura: the spell the caster used and the number its row states (armor for
// Devotion, damage for Retribution, a percentage for Concentration, resistance for the three
// resistance auras). The paladin registers an aura per rank; each joins the same categories as the
// generated party-buff copy, which is the top rank.
type PaladinAuraRank struct {
	SpellID int32
	Rank    int32
	Value   float64
}

var RetributionAuraMaxRank = PaladinAuraRank{SpellID: 10301, Value: RetributionAuraValue(0)}

func paladinAuraLabel(name string, isPlayer bool, rank PaladinAuraRank) string {
	label := fmt.Sprintf("%s (%s)", name, Ternary(isPlayer, "Player", "External"))
	if rank.Rank > 0 {
		label += fmt.Sprintf(" Rank %d", rank.Rank)
	}
	return label
}

func paladinAuraBuff(name string, category string, isPlayer bool, rank PaladinAuraRank) GeneratedBuff {
	return GeneratedBuff{
		Label:          paladinAuraLabel(name, isPlayer, rank),
		ActionID:       ActionID{SpellID: rank.SpellID}.WithTag(TernaryInt32(isPlayer, 0, -1)),
		Duration:       NeverExpires,
		Category:       category,
		SharedCategory: PaladinAuraCategory,
		SingleAura:     true,
		IsPlayer:       isPlayer,
	}
}

func DevotionAuraBuff(char *Character, isPlayer bool, rank PaladinAuraRank) *Aura {
	config := paladinAuraBuff("Devotion Aura", DevotionAuraCategory, isPlayer, rank)
	config.Stats = []StatConfig{{stats.Armor, rank.Value, false}}
	return newGeneratedStatAura(&char.Unit, config)
}

func ConcentrationAura(char *Character, isPlayer bool, rank PaladinAuraRank) *Aura {
	config := paladinAuraBuff("Concentration Aura", ConcentrationAuraCategory, isPlayer, rank)
	config.Pseudo = []PseudoConfig{{PseudoStatPushbackChance, -rank.Value / 100, false, 0}}
	return newGeneratedStatAura(&char.Unit, config)
}

func FireResistanceAura(char *Character, isPlayer bool, rank PaladinAuraRank) *Aura {
	config := paladinAuraBuff("Fire Resistance Aura", FireResistanceAuraCategory, isPlayer, rank)
	config.Stats = []StatConfig{{stats.FireResistance, rank.Value, false}}
	return newGeneratedStatAura(&char.Unit, config)
}

func FrostResistanceAura(char *Character, isPlayer bool, rank PaladinAuraRank) *Aura {
	config := paladinAuraBuff("Frost Resistance Aura", FrostResistanceAuraCategory, isPlayer, rank)
	config.Stats = []StatConfig{{stats.FrostResistance, rank.Value, false}}
	return newGeneratedStatAura(&char.Unit, config)
}

func ShadowResistanceAura(char *Character, isPlayer bool, rank PaladinAuraRank) *Aura {
	config := paladinAuraBuff("Shadow Resistance Aura", ShadowResistanceAuraCategory, isPlayer, rank)
	config.Stats = []StatConfig{{stats.ShadowResistance, rank.Value, false}}
	return newGeneratedStatAura(&char.Unit, config)
}

// Retribution Aura scales with the casting paladin's Holy spell power in Forever even though its
// client row carries no coefficient (every rank and Thorns are the same: EffectBonusCoefficient 0,
// and the damage still moves with spell power in game). The coefficient is the 1.5 s cast-time
// floor over 3.5, the AoE divisor because the shield hits every attacker, and the 0.95 penalty for
// the aura effect. Confirmed at level 20: with 80 spell power rank 1 (base 7) hits for 17 to 18,
// mostly 18, which is the 17.86 this coefficient predicts; the 0.95² variant (0.129) would have
// shown mostly 17.
const RetributionAuraSpellPowerCoefficient = 1.5 / 3.5 / 3 * 0.95

// RetributionAuraBuff is the aura on the unit the shield protects. The self-cast variant reads
// the paladin's own Holy spell power through the proc spell; the external (party-buff) variant
// cannot see the providing paladin, so externalSpellPower stands in for it and the recipient's
// own stats stay out of the damage.
func RetributionAuraBuff(char *Character, isPlayer bool, rank PaladinAuraRank, externalSpellPower float64) *Aura {
	config := paladinAuraBuff("Retribution Aura", RetributionAuraCategory, isPlayer, rank)
	if char.HasAura(config.Label) {
		return char.GetAura(config.Label)
	}

	if isPlayer {
		return newDamageShield(&char.Unit, config, SpellSchoolHoly, rank.Value, RetributionAuraSpellPowerCoefficient)
	}
	return newDamageShield(&char.Unit, config, SpellSchoolHoly, rank.Value+RetributionAuraSpellPowerCoefficient*externalSpellPower, 0)
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
// innervate, power infusion, etc. This is specifically for buffs
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
// E.g. the number of other druids in the group casting innervate.
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
