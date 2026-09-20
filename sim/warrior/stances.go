package warrior

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

type Stance uint8

const (
	StanceNone          = 0
	BattleStance Stance = 1 << iota
	DefensiveStance
	BerserkerStance
)

const stanceEffectCategory = "Stance"

func (warrior *Warrior) StanceMatches(other Stance) bool {
	return (warrior.Stance & other) != 0
}

func (warrior *Warrior) makeStanceSpell(stance Stance, mask int64, defenseType core.DefenseType, aura *core.Aura, stanceCD *core.Timer) *core.Spell {
	// Tactical Mastery (1310185) is a baseline passive that keeps 10 rage on a stance change, and
	// Improved Tactical Mastery adds its ladder on top.
	// TODO: Manual review needed -- the passive has one rank and no table; 10 is its effect value.
	maxRetainedRage := 10.0 + spellData.ImprovedTacticalMastery.ValueAt(warrior.Talents.ImprovedTacticalMastery)
	actionID := aura.ActionID
	rageMetrics := warrior.NewRageMetrics(actionID)

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		DefenseType:    defenseType,
		ClassSpellMask: mask,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    stanceCD,
				Duration: time.Second * 1,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.Stance != stance
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if warrior.WarriorInputs.StanceSnapshot {
				// Delayed, so same-GCD casts are affected by the current aura.
				// Alternatively, those casts could just (artificially) happen before the stance change.
				pa := sim.GetConsumedPendingActionFromPool()
				pa.NextActionAt = sim.CurrentTime + 10*time.Millisecond
				pa.OnAction = aura.Activate
				sim.AddPendingAction(pa)
			} else {
				aura.Activate(sim)
			}

			if warrior.CurrentRage() > maxRetainedRage {
				warrior.SpendRage(sim, warrior.CurrentRage()-maxRetainedRage, rageMetrics)
			}

			warrior.Stance = stance
		},

		RelatedSelfBuff: aura,
	})
}

func (warrior *Warrior) registerBattleStanceAura() *core.Aura {
	actionID := core.ActionID{SpellID: 2457}

	aura := warrior.RegisterAura(core.Aura{
		Label:      "Battle Stance",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(warrior.DefaultStance == proto.WarriorStance_WarriorStanceBattle, core.CharacterBuildPhaseBuffs, core.CharacterBuildPhaseNone),
	}).AttachMultiplicativePseudoStatBuff(
		// TODO: Manual review needed -- Battle Stance Passive (21156) states -20% threat.
		&warrior.PseudoStats.ThreatMultiplier, 0.8,
	)

	aura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})

	return aura
}

// TODO: Manual review needed -- Defensive Stance Passive (7376) states -10% damage taken, -10% damage
// done and +30% threat.
func (warrior *Warrior) registerDefensiveStanceAura() *core.Aura {
	actionID := core.ActionID{SpellID: 71}

	aura := warrior.RegisterAura(core.Aura{
		Label:      "Defensive Stance",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(warrior.DefaultStance == proto.WarriorStance_WarriorStanceDefensive, core.CharacterBuildPhaseBuffs, core.CharacterBuildPhaseNone),
	}).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.ThreatMultiplier, 1.3,
	).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageTakenMultiplier, 0.9,
	).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageDealtMultiplier, 0.9,
	)
	if warrior.Talents.Defiance > 0 {
		// Defiance (12792) raises the stance's threat by 5% per rank while a shield is equipped, so
		// the multiplier follows both the stance and the off-hand.
		defiance := spellData.Defiance.Effect(shared.A_MOD_THREAT, 127).MultiplierAt(warrior.Talents.Defiance)
		applied := false
		apply := func(sim *core.Simulation) {
			if !applied && warrior.PseudoStats.CanBlock {
				warrior.PseudoStats.ThreatMultiplier *= defiance
				applied = true
			}
		}
		remove := func(sim *core.Simulation) {
			if applied {
				warrior.PseudoStats.ThreatMultiplier /= defiance
				applied = false
			}
		}
		aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) { apply(sim) })
		aura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) { remove(sim) })
		warrior.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}, func(sim *core.Simulation, slot proto.ItemSlot) {
			if !aura.IsActive() {
				return
			}
			if warrior.PseudoStats.CanBlock {
				apply(sim)
			} else {
				remove(sim)
			}
		})
	}

	aura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})

	return aura
}

// TODO: Manual review needed -- Berserker Stance Passive (7381) states +3% crit, +10% damage taken
// and -20% threat.
func (warrior *Warrior) registerBerserkerStanceAura() *core.Aura {
	actionId := core.ActionID{SpellID: 2458}

	aura := warrior.RegisterAura(core.Aura{
		Label:      "Berserker Stance",
		ActionID:   actionId,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(warrior.DefaultStance == proto.WarriorStance_WarriorStanceBerserker, core.CharacterBuildPhaseBuffs, core.CharacterBuildPhaseNone),
	}).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.ThreatMultiplier, 0.8,
	).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageTakenMultiplier, 1.1,
	).AttachStatBuff(stats.PhysicalCritPercent, 3)

	aura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})

	return aura
}

func (warrior *Warrior) registerStances() {
	stanceCD := warrior.NewTimer()
	battleStanceAura := warrior.registerBattleStanceAura()
	defensiveStanceAura := warrior.registerDefensiveStanceAura()
	berserkerStanceAura := warrior.registerBerserkerStanceAura()
	// DefenseType per stance's SpellCategories row: Battle Stance (2457) and Berserker Stance (2458)
	// are Melee; Defensive Stance (71) has no row (None).
	warrior.BattleStance = warrior.makeStanceSpell(BattleStance, SpellMaskBattleStance, core.DefenseTypeMelee, battleStanceAura, stanceCD)
	warrior.DefensiveStance = warrior.makeStanceSpell(DefensiveStance, SpellMaskDefensiveStance, core.DefenseTypeNone, defensiveStanceAura, stanceCD)
	warrior.BerserkerStance = warrior.makeStanceSpell(BerserkerStance, SpellMaskBerserkerStance, core.DefenseTypeMelee, berserkerStanceAura, stanceCD)

	switch warrior.DefaultStance {
	case proto.WarriorStance_WarriorStanceBattle:
		core.MakePermanent(warrior.BattleStance.RelatedSelfBuff)
	case proto.WarriorStance_WarriorStanceDefensive:
		core.MakePermanent(warrior.DefensiveStance.RelatedSelfBuff)
	case proto.WarriorStance_WarriorStanceBerserker:
		core.MakePermanent(warrior.BerserkerStance.RelatedSelfBuff)
	}
}
