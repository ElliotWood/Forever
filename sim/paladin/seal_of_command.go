package paladin

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

// The spell power coefficient of the Seal of Command proc. The client carries none on 20424; this
// is the number the Classic sim settled on.
const sealOfCommandCoefficient = 0.29

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
//
// The proc is a weapon-percent effect, and the percent applies to the paladin's own spell power
// as well as the weapon: at 27 spell power the Forever beta measured 5.5 of each proc coming from
// spell power, 20.4%, against the 29% the coefficient states. Improved Seals raises the percent itself (the tooltip
// reads 80% with the talent, 70 * 1.15). The Holy damage the target takes extra, Judgement of the
// Crusader and the damage-against-mob-type gear, is added after the percent at the full
// coefficient, so it scales at 29% whatever the talent.
func (paladin *Paladin) registerSealOfCommand(row shared.SpellData) {
	judgeRow := spellData.JudgementOfCommand.ByRank(row.Rank)

	// Melee in SpellCategories with No Active Defense: hit and crit on the melee table, never
	// dodged, parried or blocked.
	judgement := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: judgeRow.SpellID},
		SpellSchool:    judgeRow.SpellSchool,
		DefenseType:    judgeRow.DefenseType,
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
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})

	procRow := spellData.SealOfCommandTriggered.ByRank(1)
	weaponPercent := effectAt(procRow, 0).Value / 100 * spellData.ImprovedSeals.MultiplierAt(paladin.Talents.ImprovedSeals)
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
		// No BonusCoefficient: the spell power goes through the weapon percent, so the proc adds
		// it by hand below rather than letting CalcDamage add it on top.

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weaponDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower(target))
			baseDamage := (weaponDamage + sealOfCommandCoefficient*spell.SpellDamage(target)) * weaponPercent

			targetBonus := target.PseudoStats.SchoolBonusSpellDamage[stats.SchoolIndexHoly] +
				spell.Unit.AttackTables[target.UnitIndex].MobTypeBonusStats[target.MobType][stats.SpellDamage]
			baseDamage += sealOfCommandCoefficient * targetBonus

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
