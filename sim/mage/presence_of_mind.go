package mage

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// The next spell with a cast time is instant. The row states no duration, so the buff holds until
// that spell is cast.
func (mage *Mage) registerPresenceOfMindSpell() {
	if !mage.Talents.PresenceOfMind {
		return
	}

	pomRank := spellData.PresenceOfMind.HighestRank()
	hasCastTime := MageSpellsAll &^ (MageSpellInstantCast | MageSpellBlizzard | MageSpellEvocation)

	pomMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  hasCastTime,
		FloatValue: pomRank.Effect(shared.A_ADD_PCT_MODIFIER, shared.SPELLMOD_CASTING_TIME).Value / 100,
		Kind:       core.SpellMod_CastTime_Pct,
	})

	var pomSpell *core.Spell
	mage.PresenceOfMindAura = mage.RegisterAura(core.Aura{
		Label:    "Presence of Mind",
		ActionID: core.ActionID{SpellID: pomRank.SpellID},
		Duration: core.NeverExpires,
		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			pomMod.Activate()
		},
		OnExpire: func(_ *core.Aura, sim *core.Simulation) {
			pomMod.Deactivate()
			pomSpell.CD.Use(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(hasCastTime) && spell.DefaultCast.CastTime > 0 {
				aura.Deactivate(sim)
			}
		},
	})

	pomSpell = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: pomRank.SpellID},
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: MageSpellPresenceOfMind,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: pomRank.Cooldown,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, _ *core.Unit) bool {
			return mage.GCD.IsReady(sim)
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.PresenceOfMindAura.Activate(sim)
		},
		RelatedSelfBuff: mage.PresenceOfMindAura,
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: pomSpell,
		Type:  core.CooldownTypeDPS,
	})
}
