package core

import (
	"github.com/wowsims/classic/sim/core/proto"
	"github.com/wowsims/classic/sim/core/stats"
)

// Forever gives the gathering professions a combat passive. Only the two with a
// published number are modelled: Skinning's 5% damage to Beasts and Dragonkin, and
// Mining's 5% health. Herbalism's is a heal over time, which the sim does not measure.
func applyProfessionEffects(agent Agent) {
	character := agent.GetCharacter()
	if !character.Env.IsForever() {
		return
	}

	if character.HasProfession(proto.Profession_Skinning) {
		character.mobTypeDamageAura(proto.MobType_MobTypeBeast, 1.05)
		character.mobTypeDamageAura(proto.MobType_MobTypeDragonkin, 1.05)
	}

	if character.HasProfession(proto.Profession_Mining) {
		character.MultiplyStat(stats.Health, 1.05)
	}
}
