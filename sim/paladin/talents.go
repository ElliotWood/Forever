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

	// Protection Tree
	if paladin.Talents.HolyShield {
		HolyShieldRankMap.RegisterAll(paladin.registerHolyShield)
	}

	// Retribution Tree
	if paladin.Talents.SealOfCommand {
		SealOfCommandRanks.RegisterAll(paladin.registerSealOfCommandRank)
	}
	// if paladin.Talents.Repentance {
	// 	paladin.registerRepentance()
	// }
}

func (paladin *Paladin) ApplyTalents() {
	paladin.registerTalentSpells()

	paladin.registerHolyTalents()
	paladin.registerProtectionTalents()
	paladin.registerRetributionTalents()
}
