package priest

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

var TalentTreeSizes = [3]int{18, 17, 18}

type Priest struct {
	core.Character
	SelfBuffs
	Talents *proto.PriestTalents

	Latency float64

	ShadowfiendAura *core.Aura
	ShadowfiendPet  *Shadowfiend

	Shadowfiend    *core.Spell
	InnerFocusAura *core.Aura

	VampiricEmbrace *core.Spell
}

type SelfBuffs struct {
	UseShadowfiend bool
	PreShadowform  bool
	Armor          proto.PriestOptions_Armor
}

func (priest *Priest) GetCharacter() *core.Character {
	return &priest.Character
}
func (priest *Priest) GetPriest() *Priest {
	return priest
}

func (priest *Priest) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (priest *Priest) Initialize() {
	mindblastCDTimer := priest.NewTimer()
	shadowWordDeathCDTimer := priest.NewTimer()

	MindBlastRankMap.RegisterAll(func(rank shared.SpellData) {
		priest.registerMindBlastSpell(rank, mindblastCDTimer)
	})
	ShadowWordPainRankMap.RegisterAll(priest.registerShadowWordPainSpell)
	ShadowWordDeathRankMap.RegisterAll(func(rank shared.SpellData) {
		priest.registerShadowWordDeathSpell(rank, shadowWordDeathCDTimer)
	})
	SmiteRankMap.RegisterAll(priest.registerSmiteSpell)
	priest.registerShadowfiendSpell()

	if priest.Race == proto.Race_RaceNightElf {
		starshardsCDTimer := priest.NewTimer()
		StarshardsRankMap.RegisterAll(func(rank shared.SpellData) {
			priest.registerStarshardsSpell(rank, starshardsCDTimer)
		})
	}
	if priest.Race == proto.Race_RaceUndead {
		devouringPlagueCDTimer := priest.NewTimer()
		DevouringPlagueRankMap.RegisterAll(func(rank shared.SpellData) {
			priest.registerDevouringPlagueSpell(rank, devouringPlagueCDTimer)
		})
	}
}

func (priest *Priest) ApplyTalents() {
	// Discipline
	priest.applyInnerFocus()
	priest.applyMeditation()
	priest.applyMentalAgility()
	priest.applyMentalStrength()
	priest.applySilentResolve()
	priest.applyPowerInfusion()

	// Holy
	priest.applyHolyNova()
	priest.applyDivineFury()
	priest.applySearingLight()
	priest.applySpiritualGuidance()
	priest.applySpiritOfRedemption()

	// Shadow
	priest.applyMindFlay()
	priest.applyDarkness()
	priest.applyShadowFocus()
	priest.applyImprovedShadowWordPain()
	priest.applyShadowAffinity()
	priest.applyShadowWeaving()
	priest.applyShadowform()
	priest.applyVampiricEmbrace()
	priest.applyImprovedMindBlast()

	// Forever additions, not yet implemented.
	// Discipline
	priest.applyPowerInLight()
	priest.applyWandSpecialization()
	priest.applyTwinDisciplines()
	priest.applyHolyPrecision()
	priest.applyImprovedPowerWordShield()
	priest.applyMartyrdom()
	priest.applyImprovedInnerFire()
	priest.applySoulWarding()
	priest.applyImprovedManaBurn()
	priest.applyPenance()
	priest.applyRenewedHope()
	priest.applyDivineAegis()
	// Holy
	priest.applyTwilightFocus()
	priest.applyImprovedRenew()
	priest.applyHolySpecialization()
	priest.applySpellWarding()
	priest.applyBlessedRecovery()
	priest.applyInspiration()
	priest.applyHolyReach()
	priest.applyImprovedHealing()
	priest.applyBindingHeal()
	priest.applyLitanyOfLight()
	priest.applySpiritualHealing()
	priest.applyPrayerOfMending()
	// Shadow Magic
	priest.applyBlackout()
	priest.applySpiritTap()
	priest.applyShadowReach()
	priest.applyImprovedPsychicScream()
	priest.applyImprovedMindFlay()
	priest.applyImprovedFade()
	priest.applySilence()
	priest.applyDevouringContagion()
	priest.applyEarlyDemise()
}

func (priest *Priest) Reset(_ *core.Simulation) {

}

func (priest *Priest) OnEncounterStart(sim *core.Simulation) {
}

func New(char *core.Character, selfBuffs SelfBuffs, talents string) *Priest {
	priest := &Priest{
		Character: *char,
		SelfBuffs: selfBuffs,
		Talents:   &proto.PriestTalents{},
	}

	core.FillTalentsProto(priest.Talents.ProtoReflect(), talents, TalentTreeSizes)
	priest.EnableManaBar()
	if selfBuffs.Armor == proto.PriestOptions_InnerFire {
		// Inner Fire rank 7: +1580 armor. The charges never run out on a caster that is not hit.
		priest.AddStat(stats.Armor, 1580)
	}
	priest.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[char.Class])
	priest.ShadowfiendPet = priest.NewShadowfiend()

	return priest
}

// Agent is a generic way to access underlying priest on any of the agents.
type PriestAgent interface {
	GetPriest() *Priest
}

func NewPriest(character *core.Character, options *proto.Player) *Priest {
	classOptions := options.GetDpsPriest().GetOptions().GetClassOptions()
	selfBuffs := SelfBuffs{
		UseShadowfiend: true,
		PreShadowform:  classOptions.GetPreShadowform(),
		Armor:          classOptions.GetArmor(),
	}

	basePriest := New(character, selfBuffs, options.TalentsString)
	basePriest.Latency = float64(basePriest.ChannelClipDelay.Milliseconds())

	return basePriest
}

func RegisterPriest() {
	core.RegisterAgentFactory(
		proto.Player_DpsPriest{},
		proto.Spec_SpecDpsPriest,
		func(character *core.Character, options *proto.Player, _ *proto.Raid) core.Agent {
			return NewPriest(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DpsPriest)
			if !ok {
				panic("Invalid spec value for Priest!")
			}
			player.Spec = playerSpec
		},
	)
}

const (
	PriestSpellFlagNone        int64 = 0
	PriestSpellDevouringPlague int64 = 1 << iota
	PriestSpellDevouringPlagueDoT
	PriestSpellDevouringPlagueHeal
	PriestSpellHolyNova
	PriestSpellHolyFire
	PriestSpellMindBlast
	PriestSpellMindFlay
	PriestSpellPowerInfusion
	PriestSpellStarshards
	PriestSpellShadowform
	PriestSpellShadowWordDeath
	PriestSpellShadowWordPain
	PriestSpellShadowFiend
	PriestSpellVampiricEmbrace
	PriestSpellVampiricTouch
	PriestSpellFade
	PriestSpellSmite

	PriestSpellLast
	PriestSpellsAll    = PriestSpellLast<<1 - 1
	PriestSpellDoT     = PriestSpellDevouringPlague | PriestSpellHolyFire | PriestSpellMindFlay | PriestSpellShadowWordPain | PriestSpellVampiricTouch | PriestSpellStarshards
	PriestSpellInstant = PriestSpellDevouringPlague |
		PriestSpellFade |
		PriestSpellHolyNova |
		PriestSpellPowerInfusion |
		PriestSpellShadowWordDeath |
		PriestSpellShadowWordPain |
		PriestSpellVampiricEmbrace |
		PriestSpellShadowFiend |
		PriestSpellStarshards |
		PriestSpellShadowform |
		PriestSpellPowerInfusion
	PriestShadowSpells = PriestSpellDevouringPlague |
		PriestSpellShadowWordDeath |
		PriestSpellShadowform |
		PriestSpellShadowWordPain |
		PriestSpellMindFlay |
		PriestSpellMindBlast |
		PriestSpellVampiricTouch |
		PriestSpellShadowFiend |
		PriestSpellVampiricEmbrace
	PriestHolySpells = PriestSpellSmite | PriestSpellHolyFire | PriestSpellHolyNova
)
