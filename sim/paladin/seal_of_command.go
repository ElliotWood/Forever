package paladin

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

// Seal of Command (talent)
// https://www.wowhead.com/forever/spell=20920
//
// Gives the Paladin a chance to deal additional Holy damage equal to 70% of normal weapon damage.
// Only one Seal can be active on the Paladin at any one time. Lasts 30 sec.
//
// Unleashing this Seal's energy will judge an enemy, instantly causing Holy damage, double if the
// target is stunned or incapacitated.
//
// The client states the stunned number and halves it otherwise. The proc spell 20424 has no
// rank subtext and so is in no table: its 70% of weapon damage is read from the trigger row, and
// its 0.29 coefficient and 7 procs per minute are the numbers the Classic sim carries.
func (paladin *Paladin) registerSealOfCommand(row shared.SpellData) {
	judgeRow := spellData.JudgementOfCommand.ByRank(row.Rank)

	judgement := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: judgeRow.SpellID},
		SpellSchool: core.SpellSchoolHoly,
		// The judgement is a dummy that triggers the damage spell, and that one is Melee in
		// SpellCategories: it crits on the melee table.
		DefenseType:    core.DefenseTypeMelee,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskJudgementOfCommand,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: judgeRow.Direct.BonusCoefficient(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := judgeRow.Direct.Damage(sim)
			if !target.PseudoStats.Stunned {
				baseDamage /= 2
			}
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	procRow := spellData.SealOfCommandTriggered.ByRank(1)
	weaponPercent := effectAt(procRow, 0).Value / 100
	procSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: procRow.SpellID},
		SpellSchool: core.SpellSchoolHoly,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		// 20424 carries Not a Proc: it is an ability hit to every listener, weapon procs included.
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		ClassSpellMask: SpellMaskSealOfCommandProc,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.29,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower(target)) * weaponPercent
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			dealAfterBatch(sim, spell, result)
		},
	})

	// The seal and its Echo roll the same chance: an Echo of Command "empowers your next melee
	// attack with a chance to activate Seal of Command".
	dpm := paladin.NewLegacyPPMManager(7, core.ProcMaskMeleeWhiteHit)
	icd := core.Cooldown{Timer: paladin.NewTimer(), Duration: time.Second}
	tryProc := func(sim *core.Simulation, target *core.Unit) {
		if icd.IsReady(sim) && dpm.Proc(sim, core.ProcMaskMeleeMHAuto, "Seal of Command") {
			icd.Use(sim)
			procSpell.Cast(sim, target)
		}
	}

	aura := paladin.makeSealExclusive(paladin.RegisterAura(core.Aura{
		Label:    sealLabel("Seal of Command", paladin, row),
		ActionID: core.ActionID{SpellID: row.SpellID},
		Duration: sealDuration,
	}).AttachProcTrigger(core.ProcTrigger{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeWhiteHit,
		Outcome:  core.OutcomeLanded,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			tryProc(sim, result.Target)
		},
	}))

	paladin.registerSealSpell(&sealConfig{
		row:       row,
		classMask: SpellMaskSealOfCommand,
		aura:      aura,
		judgement: judgement,
		echoID:    echoOfCommandID,
		echo:      tryProc,
	})
}
