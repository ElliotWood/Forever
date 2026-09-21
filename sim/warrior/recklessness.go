package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
	"github.com/wowsims/forever/sim/core/stats"
)

var recklessnessRank = spellData.Recklessness.Highest()

func (warrior *Warrior) registerRecklessness() {
	recklessnessCritValue := recklessnessRank.Effect(dbcenums.A_MOD_CRIT_PCT, 0).Average(core.CharacterLevel)
	aura := warrior.RegisterAura(spelldata.AuraConfig(recklessnessRank)).AttachStatsBuff(
		stats.Stats{
			stats.PhysicalCritPercent: recklessnessCritValue,
			stats.SpellCritPercent:    recklessnessCritValue,
		},
	).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageTakenMultiplier,
		1+recklessnessRank.Effect(dbcenums.A_MOD_DAMAGE_PERCENT_TAKEN, 127).Percent(),
	).
		AttachFearImmunity()

	config := spelldata.SpellConfig(&warrior.Unit, recklessnessRank,
		spelldata.Flags(core.SpellFlagAPL|core.SpellFlagCastWhileIncapacitated))
	config.ClassSpellMask = SpellMaskRecklessness

	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return warrior.StanceMatches(BerserkerStance)
	}

	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		aura.Activate(sim)
	}

	config.RelatedSelfBuff = aura

	spell := warrior.RegisterSpell(config)

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}
