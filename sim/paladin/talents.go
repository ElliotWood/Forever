package paladin

import (
	"github.com/wowsims/forever/sim/common/shared"
)

func (paladin *Paladin) registerTalentSpells() {
	// Holy Tree
	if paladin.Talents.DivineFavor {
		paladin.registerDivineFavor()
	}
	if paladin.Talents.HolyShock {
		paladin.registerHolyShock(shared.SpellData{})
	}
	// Divine Illumination: Forever drops the talent; see registerDivineIllumination

	// Protection Tree
	if paladin.Talents.HolyShield {
		HolyShieldRankMap.RegisterAll(paladin.registerHolyShield)
	}
	// Avenger's Shield: Forever drops the talent; see registerAvengersShield

	// Retribution Tree
	if paladin.Talents.SealOfCommand {
		SealOfCommandRanks.RegisterAll(paladin.registerSealOfCommandRank)
	}
	// Sanctity Aura: Forever drops the talent; see registerSanctityAura
	// if paladin.Talents.Repentance {
	// 	paladin.registerRepentance()
	// }
	// Crusader Strike: Forever drops the talent; see registerCrusaderStrike
}

func (paladin *Paladin) ApplyTalents() {
	paladin.registerTalentSpells()

	paladin.registerHolyTalents()
	paladin.registerProtectionTalents()
	paladin.registerRetributionTalents()
}
