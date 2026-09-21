package paladin

import (
	"fmt"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
	// "github.com/wowsims/forever/sim/core/stats" -- used only by the commented-out implementation below
)

// Package-level state the commented-out implementations used:
// var SealOfJusticeRanks = sealRankMap{

type proc struct {
	spellID int32
	value   float64
	coeff   float64
}

type judge struct {
	spellID   int32
	minDamage float64
	maxDamage float64
	coeff     float64
}

type seal struct {
	rank     int32
	spellID  int32
	manaCost float64
	proc     proc
	judge    judge
}

func (seal seal) GetRank() int32 { return seal.rank }

func (seal seal) GetSpellID() int32 { return seal.spellID }

func (seal seal) GetRankLabel() string {
	return fmt.Sprintf("Rank %d", seal.rank)
}

// Seal of Justice is the one family the client charges nothing for, on either rank.
func sealMana(s seal, mana float64) seal {
	s.manaCost = mana
	return s
}

// Everything the client states, judgement damage included - Righteousness, the Crusader and Command
// all state theirs.
func sealOf(seals, judges shared.SpellDataTable, rank int32, p proc) seal {
	j := judges.ByRank(rank)
	return sealWithJudgement(seals, judges, rank, p, judge{
		minDamage: shared.SpellDataMin(j.Direct),
		maxDamage: shared.SpellDataMax(j.Direct),
		coeff:     shared.SpellDataCoef(j.Direct),
	})
}

// For the three families whose judgement damage has to be supplied. A seal is three spells for one
// rank, so it keeps its own row rather than becoming a SpellData. The proc comes in by hand for every
// family but Righteousness, whose row's Direct holds it.
func sealWithJudgement(seals, judges shared.SpellDataTable, rank int32, p proc, j judge) seal {
	s := seals.ByRank(rank)
	j.spellID = judges.ByRank(rank).SpellID
	return seal{rank: rank, spellID: s.SpellID, manaCost: float64(s.Cost), proc: p, judge: j}
}

// Seal of the Crusader has no separate JudgementOfTheCrusader table - each rank's own row carries the
// judgement's spell ID as an A_DUMMY effect (index 2), the same way every other seal's own effect 0
// points back at its judgement.
func sealOfTheCrusader(rank shared.SpellData) seal {
	return seal{
		rank:     rank.Rank,
		spellID:  rank.SpellID,
		manaCost: float64(rank.Cost),
		proc:     proc{value: rank.Effect(shared.A_MOD_ATTACK_POWER, 0).High()},
		judge:    judge{spellID: int32(rank.Effect(shared.A_DUMMY, 0).Value)},
	}
}

// The seal's rows as the client states them, for the regeneration gate; SealOfRighteousnessRanks is
// what registers.
var SealOfRighteousnessTable = spellData.SealOfRighteousness

// The per-hit proc of one Seal of Righteousness rank: the damage spell it fires, and the number and
// coefficient the row's Direct holds. Those are the judgement's effect 3, a dummy the seal's tooltip
// renders each hit from ("$/87;20286s3" on rank 8) and the generator follows the description to. The
// seal's own effect 0 states the same number on every rank but carries 0.1 for a coefficient on ranks
// 1-7 and nothing on rank 8, where the judgement's dummy says 0.058 on rank 1 rising to 0.2 from rank
// 4. The proc's damage formula scales with spell power on its own terms, so the coefficient is carried
// here rather than applied; the hand values it replaces were half the client's on every rank.
func righteousnessProc(spellID int32, rank int32) proc {
	d := spellData.SealOfRighteousness.ByRank(rank).Direct
	return proc{spellID: spellID, value: shared.SpellDataMin(d), coeff: shared.SpellDataCoef(d)}
}

var SealOfRighteousnessRanks = sealRankMap{
	// The judgement damage matches the hand rows on ranks 2-8; rank 1 was a flat 26 against the
	// client's 25-26, and the client's is taken.
	sealOf(spellData.SealOfRighteousness, spellData.JudgementOfRighteousness, 1, righteousnessProc(25742, 1)),
	sealOf(spellData.SealOfRighteousness, spellData.JudgementOfRighteousness, 2, righteousnessProc(25740, 2)),
	sealOf(spellData.SealOfRighteousness, spellData.JudgementOfRighteousness, 3, righteousnessProc(25739, 3)),
	sealOf(spellData.SealOfRighteousness, spellData.JudgementOfRighteousness, 4, righteousnessProc(25738, 4)),
	sealOf(spellData.SealOfRighteousness, spellData.JudgementOfRighteousness, 5, righteousnessProc(25737, 5)),
	sealOf(spellData.SealOfRighteousness, spellData.JudgementOfRighteousness, 6, righteousnessProc(25736, 6)),
	sealOf(spellData.SealOfRighteousness, spellData.JudgementOfRighteousness, 7, righteousnessProc(25735, 7)),
	sealOf(spellData.SealOfRighteousness, spellData.JudgementOfRighteousness, 8, righteousnessProc(25713, 8)),
}

var SealOfLightRanks = sealRankMap{
	// These derive - the heal is on a same-name spell the client never grants, 20185 to 20267 - but
	// every rule that reaches it also breaks something: by skill line puts Blizzard's tick damage on
	// the channel, by ungranted name gives warrior Enrage a Direct off a creature ability.
	sealWithJudgement(spellData.SealOfLight, spellData.JudgementOfLight, 1, proc{spellID: 20167, value: 39, coeff: 0.0}, judge{minDamage: 25, maxDamage: 25, coeff: 0.0}),
	sealWithJudgement(spellData.SealOfLight, spellData.JudgementOfLight, 2, proc{spellID: 20333, value: 53, coeff: 0.0}, judge{minDamage: 34, maxDamage: 34, coeff: 0.0}),
	sealWithJudgement(spellData.SealOfLight, spellData.JudgementOfLight, 3, proc{spellID: 20334, value: 76, coeff: 0.0}, judge{minDamage: 49, maxDamage: 49, coeff: 0.0}),
	sealWithJudgement(spellData.SealOfLight, spellData.JudgementOfLight, 4, proc{spellID: 20340, value: 94, coeff: 0.0}, judge{minDamage: 61, maxDamage: 61, coeff: 0.0}),
	// TODO: Forever drops Seal of Light rank 5; the row is removed rather than indexing a rank
	// neither the seal nor its judgement table holds.
}

var SealOfWisdomRanks = sealRankMap{
	sealWithJudgement(spellData.SealOfWisdom, spellData.JudgementOfWisdom, 1, proc{spellID: 20168, value: 50, coeff: 0.0}, judge{minDamage: 33, maxDamage: 33, coeff: 0.0}),
	sealWithJudgement(spellData.SealOfWisdom, spellData.JudgementOfWisdom, 2, proc{spellID: 20350, value: 71, coeff: 0.0}, judge{minDamage: 46, maxDamage: 46, coeff: 0.0}),
	sealWithJudgement(spellData.SealOfWisdom, spellData.JudgementOfWisdom, 3, proc{spellID: 20351, value: 90, coeff: 0.0}, judge{minDamage: 59, maxDamage: 59, coeff: 0.0}),
	// TODO: Forever drops Seal of Wisdom rank 4; the row is removed rather than indexing a rank
	// neither the seal nor its judgement table holds.
}

var SealOfTheCrusaderRanks = sealRankMap{
	// TBC shipped 7 ranks; Forever's SealOfTheCrusader table stops at 6, so there is no rank 7 row to
	// index here.
	sealOfTheCrusader(spellData.SealOfTheCrusader.ByRank(1)),
	sealOfTheCrusader(spellData.SealOfTheCrusader.ByRank(2)),
	sealOfTheCrusader(spellData.SealOfTheCrusader.ByRank(3)),
	sealOfTheCrusader(spellData.SealOfTheCrusader.ByRank(4)),
	sealOfTheCrusader(spellData.SealOfTheCrusader.ByRank(5)),
	sealOfTheCrusader(spellData.SealOfTheCrusader.ByRank(6)),
}

var SealOfCommandRanks = sealRankMap{
	// The 70 lives on proc spell 20424, which has no rank subtext and so is in no table. The
	// judgement damage is the client's full number; registerSealOfCommandRank halves it.
	sealOf(spellData.SealOfCommand, spellData.JudgementOfCommand, 1, proc{spellID: 20424, value: 0.70, coeff: 0.29}),
	sealOf(spellData.SealOfCommand, spellData.JudgementOfCommand, 2, proc{spellID: 20424, value: 0.70, coeff: 0.29}),
	sealOf(spellData.SealOfCommand, spellData.JudgementOfCommand, 3, proc{spellID: 20424, value: 0.70, coeff: 0.29}),
	sealOf(spellData.SealOfCommand, spellData.JudgementOfCommand, 4, proc{spellID: 20424, value: 0.70, coeff: 0.29}),
	sealOf(spellData.SealOfCommand, spellData.JudgementOfCommand, 5, proc{spellID: 20424, value: 0.70, coeff: 0.29}),
	// TODO: Forever drops Seal of Command rank 6; the row is removed rather than indexing a rank
	// neither the seal nor its judgement table holds.
}

func (paladin *Paladin) registerSeals() {
	SealOfRighteousnessRanks.RegisterAll(paladin.registerSealOfRighteousness)
	SealOfLightRanks.RegisterAll(paladin.registerSealOfLight)
	SealOfWisdomRanks.RegisterAll(paladin.registerSealOfWisdom)
	paladin.registerSealOfJustice(seal{})
	SealOfTheCrusaderRanks.RegisterAll(paladin.registerSealOfTheCrusader)
}

// Seal Twist
const TwistTag = "Twistable"

// Command -> Righteousness
// Command -> Wisdom
// Command -> Light
// Command -> Justice

// Righteous -> Command
// Righteous -> Wisdom
// Righteous -> Light
// Righteous -> Justice

// Wisdom -> X

// Light -> X

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Justice -> X
func (paladin *Paladin) applySeal(newSeal *core.Aura, sealSpell *core.Spell, judgement *core.Spell, sim *core.Simulation) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// if paladin.CurrentSeal != nil {
	// 	newSealLabel := newSeal.ActionID.SpellID
	// 	if newSealLabel == 0 {
	// 		newSealLabel = newSeal.ActionIDForProc.SpellID
	// 	}
	//
	// 	currentSealLabel := paladin.CurrentSeal.ActionID.SpellID
	// 	if currentSealLabel == 0 {
	// 		currentSealLabel = paladin.CurrentSeal.ActionIDForProc.SpellID
	// 	}
	// 	// If they are recasting the same seal, reactivate or refresh
	// 	if newSealLabel == currentSealLabel {
	// 		paladin.CurrentSeal.Activate(sim)
	// 		return
	// 	}
	// }
	//
	// // Twisting only occurs when current seal is Command or Righteousness
	// if paladin.CurrentSeal.IsActive() {
	// 	if paladin.CurrentSeal.Tag == TwistTag {
	// 		paladin.PreviousSealSpell = sealSpell
	// 		paladin.PreviousSeal = paladin.CurrentSeal
	// 		paladin.PreviousJudgement = paladin.CurrentJudgement
	// 		pendingAction := core.NewDelayedAction(core.DelayedActionOptions{
	// 			DoAt:     sim.CurrentTime + (time.Millisecond * 399),
	// 			Priority: core.ActionPriorityLow,
	// 			OnAction: func(sim *core.Simulation) {
	// 				paladin.PreviousSeal.Deactivate(sim)
	// 			},
	// 		})
	// 		sim.AddPendingAction(pendingAction)
	// 	} else {
	// 		paladin.CurrentSeal.Deactivate(sim)
	// 	}
	// }
	//
	// paladin.CurrentSealSpell = sealSpell
	// paladin.CurrentSeal = newSeal
	// paladin.CurrentJudgement = judgement
	// paladin.CurrentSeal.Activate(sim)
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Seal of Righteousness
// https://www.wowhead.com/forever/spell=21084
//
// Fills the Paladin with divine spirit for 30 sec, granting each melee attack
// additional Holy damage. Only one Seal can be active on the Paladin at any one time.
//
// Unleashing this Seal's energy will judge an enemy, instantly causing Holy damage.
func (paladin *Paladin) registerSealOfRighteousness(seal seal) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// // ~~~~~~~~~ SEASON OF DISCOVERY DESCRIPTION, INFO SHOULD BE VERIFIED ~~~~~~~~~
	//
	// /*
	// * Seal of Righteousness is a Spell/Aura that when active makes the paladin capable of procing
	// * two different SpellIDs depending on a paladin's casted spell or melee swing.
	// *
	// * (Judgement of Righteousness):
	// *   - Deals flat damage that is affected by Improved SoR talent, and
	// *     has a spellpower scaling that is unaffected by that talent.
	// *   - Targets magic defense and rolls to hit and crit.
	// *
	// * (Seal of Righteousness):
	// *   - Procs from white hits.
	// *   - Cannot miss or be dodged/parried/blocked if the underlying white hit lands.
	// *   - Deals damage that is a function of weapon speed, and spellpower.
	// *   - Calculates damage including spellpower scaling but ignoring damage multipliers,
	// *      then feeds that value as base damage into the proc spell.
	//  */
	//
	// judgeSpell := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: seal.judge.spellID},
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagBinary,
	// 	ClassSpellMask: SpellMaskJudgementOfRighteousness,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	BonusCoefficient: seal.judge.coeff,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		flags := spell.Flags
	// 		baseDamage := sim.Roll(seal.judge.minDamage, seal.judge.maxDamage)
	// 		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	//
	// 		action := core.NewDelayedAction(core.DelayedActionOptions{
	// 			DoAt:     sim.CurrentTime + core.SpellBatchWindow,
	// 			Priority: core.ActionPriorityLow,
	// 			OnAction: func(sim *core.Simulation) {
	// 				currentFlags := spell.Flags
	// 				spell.Flags = flags
	// 				spell.DealDamage(sim, result)
	// 				spell.Flags = currentFlags
	// 			},
	// 		})
	//
	// 		sim.AddPendingAction(action)
	// 	},
	// })
	//
	// // Canonical Seal of Righteousness proc formula (Maintankadin / EJ, matches in-game testing):
	// //   1H: damage = (0.85 * SoRcoef * Speed) - (QualityModifier * Speed * 0.03) + (0.03 * AvgWeaponDmg) + (0.092 * Speed * SP)
	// //   2H: damage = (1.20 * SoRcoef * Speed) - (QualityModifier * Speed * 0.03) + (0.03 * AvgWeaponDmg) + (0.108 * Speed * SP)
	// sorCoef := seal.proc.value * 1.2 * 1.03 / 100
	// procSpell := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:    core.ActionID{SpellID: seal.proc.spellID},
	// 	SpellSchool: core.SpellSchoolHoly,
	// 	DefenseType: core.DefenseTypeMelee,
	// 	ProcMask:    core.ProcMaskMeleeMHSpecial,
	// 	// The damage spells (25713 .. 27156) carry Suppress Weapon Procs and are procs: weapon procs
	// 	// and auras without Can Proc From Procs never hear them.
	// 	Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagProc | core.SpellFlagSuppressWeaponProcs,
	// 	ClassSpellMask: SpellMaskSealOfRighteousness,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	//
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		mh := paladin.MainHand()
	//
	// 		baseCoef := 0.85
	// 		spCoef := 0.092
	// 		if mh.HandType == proto.HandType_HandTypeTwoHand {
	// 			baseCoef = 1.2
	// 			spCoef = 0.108
	// 		}
	//
	// 		speed := mh.SwingSpeed
	// 		spell.BonusCoefficient = spCoef * speed
	//
	// 		avgWeaponDmg := paladin.AutoAttacks.MH().AverageDamage()
	// 		flatDamage := baseCoef*sorCoef*speed - mh.QualityModifier*speed*0.03 + 0.03*avgWeaponDmg
	// 		result := spell.CalcDamage(sim, target, flatDamage, spell.OutcomeAlwaysHit)
	//
	// 		action := core.NewDelayedAction(core.DelayedActionOptions{
	// 			DoAt:     sim.CurrentTime + core.SpellBatchWindow,
	// 			Priority: core.ActionPriorityLow,
	// 			OnAction: func(sim *core.Simulation) {
	// 				spell.DealDamage(sim, result)
	// 			},
	// 		})
	//
	// 		sim.AddPendingAction(action)
	// 	},
	// })
	//
	// aura := paladin.MakeProcTriggerAura(core.ProcTrigger{
	// 	Name:            "Seal of Righteousness" + paladin.Label + " " + seal.GetRankLabel(),
	// 	ActionID:        core.ActionID{SpellID: seal.spellID},
	// 	MetricsActionID: core.ActionID{SpellID: seal.spellID},
	// 	Duration:        time.Second * 30,
	// 	Outcome:         core.OutcomeLanded,
	// 	Callback:        core.CallbackOnSpellHitDealt,
	// 	ProcMask:        core.ProcMaskMeleeWhiteHit,
	// 	Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 		procSpell.Cast(sim, result.Target)
	// 	},
	// })
	// aura.Tag = TwistTag
	//
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: seal.spellID},
	// 	ClassSpellMask: SpellMaskSealOfRighteousness,
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskEmpty,
	// 	Flags:          core.SpellFlagAPL,
	// 	Rank:           seal.rank,
	//
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	//
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: int32(seal.manaCost),
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			GCD: core.GCDDefault,
	// 		},
	// 	},
	//
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		paladin.applySeal(aura, spell, judgeSpell, sim)
	// 	},
	// })
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Seal of Light
// https://www.wowhead.com/forever/spell=20165
//
// Fills the Paladin with divine light for 30 sec, giving each melee attack
// a chance to heal the Paladin. Only one Seal can be active on the Paladin
// at any one time.
//
// Unleashing this Seal's energy will judge an enemy for 40 sec, granting
// attacks against the judged enemy a chance to heal the attacker.
func (paladin *Paladin) registerSealOfLight(seal seal) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// judgementOfLightAuras := paladin.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
	// 	return core.JudgementOfLightAura(target, true, 0)
	// })
	// paladin.JudgementAuras = append(paladin.JudgementAuras, judgementOfLightAuras)
	//
	// judgeSpell := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:         core.ActionID{SpellID: seal.judge.spellID},
	// 	SpellSchool:      core.SpellSchoolHoly,
	// 	DefenseType:      core.DefenseTypeMagic,
	// 	ProcMask:         core.ProcMaskEmpty,
	// 	Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagBinary,
	// 	ClassSpellMask:   SpellMaskJudgementOfLight,
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
	// 		judgementOfLightAuras.Get(target).Activate(sim)
	// 	},
	// })
	//
	// procSpell := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:         core.ActionID{SpellID: seal.proc.spellID},
	// 	ClassSpellMask:   SpellMaskSealOfLight,
	// 	SpellSchool:      core.SpellSchoolHoly,
	// 	ProcMask:         core.ProcMaskSpellHealing,
	// 	Flags:            core.SpellFlagHelpful | core.SpellFlagPassiveSpell | core.SpellFlagProc, // 20167, 20333, 20334, 20340, 27161 lack Not a Proc.
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.CalcAndDealHealing(sim, target, seal.proc.value, spell.OutcomeAlwaysHit)
	// 	},
	// })
	//
	// aura := paladin.MakeProcTriggerAura(core.ProcTrigger{
	// 	Name:            "Seal of Light" + paladin.Label + " " + seal.GetRankLabel(),
	// 	ActionID:        core.ActionID{SpellID: seal.spellID},
	// 	MetricsActionID: core.ActionID{SpellID: seal.spellID},
	// 	Duration:        time.Second * 30,
	// 	Outcome:         core.OutcomeLanded,
	// 	Callback:        core.CallbackOnSpellHitDealt,
	// 	ProcMask:        core.ProcMaskMeleeWhiteHit,
	// 	Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 		procSpell.Cast(sim, result.Target)
	// 	},
	// })
	//
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: seal.spellID},
	// 	ClassSpellMask: SpellMaskSealOfLight,
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskEmpty,
	// 	Flags:          core.SpellFlagAPL,
	// 	Rank:           seal.rank,
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: int32(seal.manaCost),
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{GCD: core.GCDDefault},
	// 	},
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		paladin.applySeal(aura, spell, judgeSpell, sim)
	// 	},
	// })
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Seal of Wisdom
// https://www.wowhead.com/forever/spell=20166
//
// Fills the Paladin with divine wisdom for 30 sec, giving each melee attack
// a chance to restore mana to the Paladin. Only one Seal can be active on
// the Paladin at any one time.
//
// Unleashing this Seal's energy will judge an enemy for 40 sec, granting
// attacks against the judged enemy a chance to restore mana to the attacker.
func (paladin *Paladin) registerSealOfWisdom(seal seal) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// judgementOfWisdomAuras := paladin.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
	// 	return core.JudgementOfWisdomAura(target, true, 0)
	// })
	// paladin.JudgementAuras = append(paladin.JudgementAuras, judgementOfWisdomAuras)
	//
	// judgeSpell := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:         core.ActionID{SpellID: seal.judge.spellID},
	// 	SpellSchool:      core.SpellSchoolHoly,
	// 	DefenseType:      core.DefenseTypeMagic,
	// 	ProcMask:         core.ProcMaskEmpty,
	// 	Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagBinary,
	// 	ClassSpellMask:   SpellMaskJudgementOfWisdom,
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
	// 		judgementOfWisdomAuras.Get(target).Activate(sim)
	// 	},
	// })
	// sealManaMetrics := paladin.Unit.NewManaMetrics(core.ActionID{SpellID: seal.proc.spellID})
	// procSpell := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:         core.ActionID{SpellID: seal.proc.spellID},
	// 	ClassSpellMask:   SpellMaskSealOfWisdom,
	// 	SpellSchool:      core.SpellSchoolHoly,
	// 	ProcMask:         core.ProcMaskEmpty,
	// 	Flags:            core.SpellFlagHelpful | core.SpellFlagPassiveSpell | core.SpellFlagProc, // 20168, 20350, 20351, 27167 lack Not a Proc.
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		if spell.Unit.HasManaBar() {
	// 			spell.Unit.AddMana(sim, seal.proc.value, sealManaMetrics)
	// 		}
	// 	},
	// })
	// aura := paladin.MakeProcTriggerAura(core.ProcTrigger{
	// 	Name:            "Seal of Wisdom" + paladin.Label + " " + seal.GetRankLabel(),
	// 	ActionID:        core.ActionID{SpellID: seal.spellID},
	// 	MetricsActionID: core.ActionID{SpellID: seal.spellID},
	// 	Duration:        time.Second * 30,
	// 	Outcome:         core.OutcomeLanded,
	// 	Callback:        core.CallbackOnSpellHitDealt,
	// 	ProcMask:        core.ProcMaskMeleeWhiteHit,
	// 	Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 		procSpell.Cast(sim, result.Target)
	// 	},
	// })
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: seal.spellID},
	// 	ClassSpellMask: SpellMaskSealOfWisdom,
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskEmpty,
	// 	Flags:          core.SpellFlagAPL,
	// 	Rank:           seal.rank,
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: int32(seal.manaCost),
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{GCD: core.GCDDefault},
	// 	},
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		paladin.applySeal(aura, spell, judgeSpell, sim)
	// 	},
	// })
}

// Seal of Justice
// https://www.wowhead.com/forever/spell=20164
//
// Fills the Paladin with the spirit of justice for 30 sec, giving each melee
// attack a chance to stun the target for 2 sec. Only one Seal can be active
// on the Paladin at any one time.
//
// Unleashing this Seal's energy will judge an enemy for 20 sec, preventing
// them from fleeing.
//
// TODO: To be implemented. The Forever client ships this as a single unranked class spell:
// it has a SkillLineAbility row but no "Rank N" subtext, so no ladder can be built for it.
func (paladin *Paladin) registerSealOfJustice(seal seal) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// registerJoJDebuff := func(target *core.Unit) *core.Aura {
	// 	return target.GetOrRegisterAura(core.Aura{
	// 		Label:    "Judgement of Justice",
	// 		ActionID: core.ActionID{SpellID: seal.judge.spellID},
	// 		Tag:      JudgementAuraTag,
	// 		Duration: time.Second * 20,
	// 		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
	// 				aura.Refresh(sim)
	// 			}
	// 		},
	// 	})
	// }
	//
	// judgementOfJusticeAuras := paladin.NewEnemyAuraArray(registerJoJDebuff)
	// paladin.JudgementAuras = append(paladin.JudgementAuras, judgementOfJusticeAuras)
	//
	// judgeSpell := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:         core.ActionID{SpellID: seal.judge.spellID},
	// 	SpellSchool:      core.SpellSchoolHoly,
	// 	DefenseType:      core.DefenseTypeMagic,
	// 	ProcMask:         core.ProcMaskEmpty,
	// 	Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagBinary,
	// 	ClassSpellMask:   SpellMaskJudgementOfJustice,
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
	// 		judgementOfJusticeAuras.Get(target).Activate(sim)
	// 	},
	// })
	// procSpell := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:         core.ActionID{SpellID: seal.proc.spellID},
	// 	ClassSpellMask:   SpellMaskSealOfJustice,
	// 	SpellSchool:      core.SpellSchoolHoly,
	// 	DefenseType:      core.DefenseTypeMagic,
	// 	ProcMask:         core.ProcMaskEmpty,
	// 	Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagProc, // 20170 lacks Not a Proc.
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
	// 	},
	// })
	// aura := paladin.MakeProcTriggerAura(core.ProcTrigger{
	// 	Name:            "Seal of Justice" + paladin.Label + " " + seal.GetRankLabel(),
	// 	ActionID:        core.ActionID{SpellID: seal.spellID},
	// 	MetricsActionID: core.ActionID{SpellID: seal.spellID},
	// 	Duration:        time.Second * 30,
	// 	Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 		procSpell.Cast(sim, result.Target)
	// 	},
	// })
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: seal.spellID},
	// 	ClassSpellMask: SpellMaskSealOfJustice,
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskEmpty,
	// 	Flags:          core.SpellFlagAPL,
	// 	Rank:           seal.rank,
	// 	ManaCost: core.ManaCostOptions{
	// 		BaseCostPercent: seal.manaCost,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{GCD: core.GCDDefault},
	// 	},
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		paladin.applySeal(aura, spell, judgeSpell, sim)
	// 	},
	// })
}

// Seal of the Crusader
// https://www.wowhead.com/forever/spell=21082
//
// Fills the Paladin with the spirit of a crusader for 30 sec, increasing
// attack speed but reducing damage caused by each weapon hit. The
// Paladin also causes additional threat. Only one Seal can be active on
// the Paladin at any one time.
//
// Unleashing this Seal's energy will judge an enemy for 20 sec, increasing
// Holy damage taken from all sources.
// TODO: To be implemented. Seal of the Crusader exists in Forever with six ranks (TBC had
// seven); SealOfTheCrusaderRanks above reads the judgement spell id out of each rank's
// A_DUMMY effect, the way every other seal in this file does. The implementation below is
// the TBC one on that data. It stays commented until it has been reviewed.
func (paladin *Paladin) registerSealOfTheCrusader(seal seal) {
	panic("To be implemented")

	// percentBonus := core.Ternary(paladin.CouldHaveSetBonus(ItemSetJusticarBattlegear, 2), 1.15, 1.0)
	//
	// judgementOfTheCrusaderAuras := paladin.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
	// 	// TODO: core.ImprovedSealOfTheCrusaderAura hardcodes 219.0 as the TBC rank-7 holy damage
	// 	// bonus (219 = "Max Rank Seal Of Crusader (Rank 7)" per its own comment). Forever's
	// 	// SealOfTheCrusader table tops out at rank 6, and spellData does not state what a
	// 	// rank-6-capped version of this aura's bonus should be, so the TBC rank-7 number stays.
	// 	// TODO: the librams that add flat holy damage taken (23203 Libram of Fervor, 27949 and
	// 	// 27983 Libram of Zeal) need a parameter for it; the aura takes the percent bonus alone.
	// 	return core.ImprovedSealOfTheCrusaderAura(target, 1, percentBonus)
	// })
	//
	// paladin.JudgementAuras = append(paladin.JudgementAuras, judgementOfTheCrusaderAuras)
	//
	// judgeSpell := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:         core.ActionID{SpellID: seal.judge.spellID},
	// 	SpellSchool:      core.SpellSchoolHoly,
	// 	DefenseType:      core.DefenseTypeMagic,
	// 	ProcMask:         core.ProcMaskEmpty,
	// 	Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagBinary,
	// 	ClassSpellMask:   SpellMaskJudgementOfTheCrusader,
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
	// 		judgementOfTheCrusaderAuras.Get(target).Activate(sim)
	// 	},
	// })
	//
	// aura := paladin.RegisterAura(core.Aura{
	// 	Label:    "Seal of the Crusader" + paladin.Label + " " + seal.GetRankLabel(),
	// 	ActionID: core.ActionID{SpellID: seal.spellID},
	// 	Duration: time.Second * 30,
	// }).
	// 	AttachMultiplyMeleeSpeed(1.4).
	// 	AttachSpellMod(core.SpellModConfig{
	// 		ProcMask:   core.ProcMaskMeleeMHAuto,
	// 		Kind:       core.SpellMod_DamageDone_Flat,
	// 		FloatValue: -0.4,
	// 	}).
	// 	AttachStatBuff(stats.AttackPower, seal.proc.value)
	//
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:         aura.ActionID,
	// 	ClassSpellMask:   SpellMaskSealOfTheCrusader,
	// 	SpellSchool:      core.SpellSchoolHoly,
	// 	DefenseType:      core.DefenseTypeMagic,
	// 	ProcMask:         core.ProcMaskEmpty,
	// 	Flags:            core.SpellFlagAPL,
	// 	Rank:             seal.rank,
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: int32(seal.manaCost),
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{GCD: core.GCDDefault},
	// 	},
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		paladin.applySeal(aura, spell, judgeSpell, sim)
	// 	},
	// 	RelatedSelfBuff: aura,
	// })
}

// TODO: To be implemented. TBC body below needs no porting; kept commented until this class's port is reviewed.
//
// Seal of Command
// https://www.wowhead.com/forever/spell=20375
//
// Gives the Paladin a chance to deal additional Holy damage equal to 70%
// of normal weapon damage. Only one Seal can be active on the Paladin at
// any one time. Lasts 30 sec.
//
// Unleashing this Seal's energy will judge an enemy, instantly causing
// 228 to 252 Holy damage, 456 to 504 if the target is stunned or incapacitated.
//
// The table carries 456-504, which is what the client states; the half is applied at cast time
// against PseudoStats.Stunned. No encounter stuns the boss, so the halved branch is the only one
// any preset takes - the stunned path ships untested.
func (paladin *Paladin) registerSealOfCommandRank(seal seal) {
	panic("To be implemented")

	// The TBC implementation, kept for the port:
	// minDamage := seal.judge.minDamage
	// maxDamage := seal.judge.maxDamage
	// judgeSpell := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:    core.ActionID{SpellID: seal.judge.spellID},
	// 	SpellSchool: core.SpellSchoolHoly,
	// 	// The judgement spell itself (20425 .. 27172) is a dummy that triggers the damage spell
	// 	// (20467, 20963 .. 20966, 27171), and that one is Melee in SpellCategories.
	// 	DefenseType:      core.DefenseTypeMelee,
	// 	ProcMask:         core.ProcMaskMeleeMHSpecial,
	// 	Flags:            core.SpellFlagMeleeMetrics,
	// 	ClassSpellMask:   SpellMaskJudgementOfCommand,
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	BonusCoefficient: seal.judge.coeff,
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		// The client states the stunned number and the game halves it otherwise, so the roll comes
	// 		// first and the halving second - rolling a halved range is not the same distribution.
	// 		baseDamage := sim.Roll(minDamage, maxDamage)
	// 		if !target.PseudoStats.Stunned {
	// 			baseDamage /= 2
	// 		}
	// 		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
	// 		action := core.NewDelayedAction(core.DelayedActionOptions{
	// 			DoAt:     sim.CurrentTime + core.SpellBatchWindow,
	// 			Priority: core.ActionPriorityLow,
	// 			OnAction: func(sim *core.Simulation) {
	// 				spell.DealDamage(sim, result)
	// 			},
	// 		})
	// 		sim.AddPendingAction(action)
	// 	},
	// })
	//
	// procSpell := paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: seal.proc.spellID},
	// 	ClassSpellMask: SpellMaskSealOfCommand,
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMelee,
	// 	ProcMask:       core.ProcMaskMeleeMHSpecial,
	// 	// 20424 carries Not a Proc: it is an ability hit to every listener, weapon procs included.
	// 	Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
	// 	DamageMultiplier: 1,
	// 	ThreatMultiplier: 1,
	// 	BonusCoefficient: seal.proc.coeff,
	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower(target)) * seal.proc.value
	// 		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
	// 		action := core.NewDelayedAction(core.DelayedActionOptions{
	// 			DoAt:     sim.CurrentTime + core.SpellBatchWindow,
	// 			Priority: core.ActionPriorityLow,
	// 			OnAction: func(sim *core.Simulation) {
	// 				spell.DealDamage(sim, result)
	// 			},
	// 		})
	// 		sim.AddPendingAction(action)
	// 	},
	// })
	//
	// aura := paladin.RegisterAura(core.Aura{
	// 	Label:    "Seal of Command" + paladin.Label + " " + seal.GetRankLabel(),
	// 	ActionID: core.ActionID{SpellID: seal.spellID},
	// 	Duration: time.Second * 30,
	// 	Tag:      TwistTag,
	// }).AttachProcTrigger(core.ProcTrigger{
	// 	Outcome:  core.OutcomeLanded,
	// 	ProcMask: core.ProcMaskMeleeWhiteHit,
	// 	Callback: core.CallbackOnSpellHitDealt,
	// 	ICD:      time.Second * 1,
	// 	DPM:      paladin.NewLegacyPPMManager(7, core.ProcMaskMeleeWhiteHit),
	// 	Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 		procSpell.Cast(sim, result.Target)
	// 	},
	// })
	//
	// paladin.RegisterSpell(core.SpellConfig{
	// 	ActionID:       aura.ActionID,
	// 	ClassSpellMask: SpellMaskSealOfCommand,
	// 	SpellSchool:    core.SpellSchoolHoly,
	// 	DefenseType:    core.DefenseTypeMagic,
	// 	ProcMask:       core.ProcMaskEmpty,
	// 	Flags:          core.SpellFlagAPL,
	// 	Rank:           seal.rank,
	// 	ManaCost: core.ManaCostOptions{
	// 		FlatCost: int32(seal.manaCost),
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{GCD: core.GCDDefault},
	// 	},
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
	// 		paladin.applySeal(aura, spell, judgeSpell, sim)
	// 	},
	// })
}
