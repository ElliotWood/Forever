package warrior

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/spelldata"
)

func (warrior *Warrior) registerFuryTalents() {
	// Tier 1
	// Booming Voice: shouts.go
	warrior.registerCruelty()

	// Tier 2
	warrior.registerIronWill()
	warrior.registerUnbridledWrath()

	// Tier 3
	warrior.registerImprovedCleave()
	warrior.registerPiercingHowl()
	warrior.registerBloodCraze()
	// Boundless Rage: warrior.go, when it enables the rage bar

	// Tier 4
	warrior.registerDualWieldSpecialization()
	warrior.registerRagingBlows()
	warrior.registerEnrage()
	warrior.registerImprovedExecute()

	// Tier 5
	warrior.registerPrecision()
	warrior.registerDeathWish()
	warrior.registerImprovedIntercept()

	// Tier 6
	// Improved Berserker Rage: berserker_rage.go
	warrior.registerFlurry()

	// Tier 7
	warrior.registerBloodthirst()
}

func (warrior *Warrior) registerCruelty() {
	if warrior.Talents.Cruelty == 0 {
		return
	}

	spelldata.ParseStatic(&warrior.Character, spellData.Cruelty.Rank(warrior.Talents.Cruelty))
}

var unbridledWrathRank = spellData.UnbridledWrathTriggered.Highest()

// The energize is on the client's 0-1000 rage bar.
var unbridledWrathRage = unbridledWrathRank.EnergizeEffect().Tenths()

func (warrior *Warrior) registerUnbridledWrath() {
	if warrior.Talents.UnbridledWrath == 0 {
		return
	}

	rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: unbridledWrathRank.ID})

	// The tooltip of 12322 doubles the rage for a two-handed weapon.
	rageGain := unbridledWrathRage
	twoHanded := func() {
		rageGain = unbridledWrathRage * core.TernaryFloat64(warrior.GetMainHandType() == proto.HandType_HandTypeTwoHand, 2, 1)
	}
	twoHanded()
	warrior.RegisterItemSwapCallback(core.AllMeleeWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
		twoHanded()
	})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Unbridled Wrath",
		ProcMask:           core.ProcMaskMeleeWhiteHit,
		ProcChance:         spellData.UnbridledWrath.FractionAt(warrior.Talents.UnbridledWrath),
		RequireDamageDealt: true,
		Outcome:            core.OutcomeLanded,
		Callback:           core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.AddRage(sim, rageGain, rageMetrics)
		},
	})
}

func (warrior *Warrior) registerDualWieldSpecialization() {
	if warrior.Talents.DualWieldSpecialization == 0 {
		return
	}

	rank := spellData.DualWieldSpecialization.Rank(warrior.Talents.DualWieldSpecialization)

	// Effect 3 states A_MOD_HIT_CHANCE with nothing narrowing it, which the parse reads as hit on
	// every attack; the tooltip states the chance to hit with off-hand attacks, so it stays a mod
	// on the off-hand's own hits.
	spelldata.ParseStatic(&warrior.Character, rank, spelldata.SkipEffects(3))

	warrior.AddStaticMod(core.SpellModConfig{
		ProcMask:   core.ProcMaskMeleeOH,
		Kind:       core.SpellMod_BonusHit_Percent,
		FloatValue: rank.EffectN(3).BaseValue(),
	})

	// The off-hand rage the tooltip's $m2 states is the dummy at index 2, which the parse skips.
	warrior.SetOffHandRageMultiplier(1 + rank.EffectN(2).Percent())
}

// Iron Will (12962) shortens the stuns and fears the warrior suffers; the client files the fear
// ladder under mechanic 1 and the stun ladder under mechanic 12.
func (warrior *Warrior) registerIronWill() {
	if warrior.Talents.IronWill == 0 {
		return
	}
	spelldata.ParseStatic(&warrior.Character, spellData.IronWill.Rank(warrior.Talents.IronWill))
}

func (warrior *Warrior) registerImprovedExecute() {
	if warrior.Talents.ImprovedExecute == 0 {
		return
	}

	spelldata.ParseStatic(&warrior.Character, spellData.ImprovedExecute.Rank(warrior.Talents.ImprovedExecute))
}

var enrageBuff = spellData.EnrageTriggered.Highest()

func (warrior *Warrior) registerEnrage() {
	if warrior.Talents.Enrage == 0 {
		return
	}

	warrior.EnrageAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Enrage",
		ActionID: core.ActionID{SpellID: enrageBuff.ID},
		Duration: enrageBuff.Duration(),
	}).AttachSpellMod(core.SpellModConfig{
		School:     core.SpellSchoolPhysical,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: spellData.Enrage.FractionAt(warrior.Talents.Enrage),
	})

	warrior.EnrageAura.NewExclusiveEffect("Enrage", true, core.ExclusiveEffect{Priority: spellData.Enrage.ValueAt(warrior.Talents.Enrage)})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Enrage - Trigger",
		Callback:           core.CallbackOnSpellHitTaken,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
		ProcChance:         spellData.Enrage.ProcChanceAt(warrior.Talents.Enrage),
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.EnrageAura.Activate(sim)
		},
	})
}

var flurryBuff = spellData.FlurryTriggered.Highest()

func (warrior *Warrior) registerFlurry() {
	if warrior.Talents.Flurry == 0 {
		return
	}

	// TODO: Ingame test needed: the talent ladder gives 5% per point (25% at rank 5) while the
	// applied buff 12966 carries a flat 30%.
	flurryAura := warrior.RegisterAura(core.Aura{
		Label:     "Flurry",
		ActionID:  core.ActionID{SpellID: flurryBuff.ID},
		Duration:  flurryBuff.Duration(),
		MaxStacks: int32(flurryBuff.ProcCharges),
	}).AttachMultiplyMeleeSpeed(spellData.Flurry.MultiplierAt(warrior.Talents.Flurry))

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Flurry - Trigger",
		ActionID:           core.ActionID{SpellID: 12319},
		ProcMask:           core.ProcMaskMelee,
		TriggerImmediately: true,
		Callback:           core.CallbackOnSpellHitDealt,
		Outcome:            core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(SpellMaskWhirlwindOh) {
				return
			}

			if result.Outcome.Matches(core.OutcomeCrit) {
				flurryAura.Activate(sim)
				flurryAura.SetStacks(sim, int32(flurryBuff.ProcCharges))
				return
			}

			if flurryAura.IsActive() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				flurryAura.RemoveStack(sim)
			}
		},
	})
}

func (warrior *Warrior) registerPrecision() {
	if warrior.Talents.Precision == 0 {
		return
	}

	spelldata.ParseStatic(&warrior.Character, spellData.Precision.Rank(warrior.Talents.Precision))
}

var bloodthirstRank = spellData.Bloodthirst.ByID(23894)

func (warrior *Warrior) registerBloodthirst() {
	if !warrior.Talents.Bloodthirst {
		return
	}

	// The attack power share sits on the second effect; the first is the flat damage added to it.
	apShare := bloodthirstRank.EffectN(2).Percent()

	config := spelldata.SpellConfig(&warrior.Unit, bloodthirstRank, spelldata.Melee(core.ProcMaskMeleeMHSpecial))
	config.ClassSpellMask = SpellMaskBloodthirst

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := spell.MeleeAttackPower(target)*apShare + bloodthirstRank.DamageEffect().Average(core.CharacterLevel)
		result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		if !result.Landed() {
			spell.IssueRefund(sim)
		}
	}

	warrior.RegisterSpell(config)
}

var piercingHowlRank = spellData.PiercingHowl.Highest()

// TODO: The daze itself is not modelled; the encounter's targets do not move, so the -50% movement
// speed 12323 applies for 6 seconds within 10 yards has nothing to act on.
func (warrior *Warrior) registerPiercingHowl() {
	if !warrior.Talents.PiercingHowl {
		return
	}

	config := spelldata.SpellConfig(&warrior.Unit, piercingHowlRank, spelldata.Flags(core.SpellFlagAPL))
	config.ProcMask = core.ProcMaskEmpty

	warrior.RegisterSpell(config)
}

var bloodCrazeHot = spellData.BloodCrazeTriggered.Highest()

func (warrior *Warrior) registerBloodCraze() {
	if warrior.Talents.BloodCraze == 0 {
		return
	}

	healthFraction := spellData.BloodCraze.EffectAt(1).FractionAt(warrior.Talents.BloodCraze)
	hitThreshold := spellData.BloodCraze.EffectAt(2).FractionAt(warrior.Talents.BloodCraze)
	tick := bloodCrazeHot.PeriodicEffect()

	config := spelldata.SpellConfig(&warrior.Unit, bloodCrazeHot, spelldata.Proc())
	config.ProcMask = core.ProcMaskSpellHealing
	config.DamageMultiplier = 1
	config.ThreatMultiplier = 1

	config.Hot = spelldata.DotConfig(bloodCrazeHot, tick)
	config.Hot.SelfOnly = true
	// The resolver's callbacks snapshot the multipliers when the dot is applied and tick the
	// effect's own damage; these ticks heal a share of maximum health with the multipliers in force
	// at the tick. The tick count, the period and BonusCoefficient stay the row's, and this effect
	// states no spell power coefficient - one that did would put spell power on every tick.
	config.Hot.OnSnapshot = nil
	config.Hot.OnTick = func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
		healPerTick := warrior.MaxHealth() * healthFraction / float64(dot.ExpectedTickCount())
		dot.Spell.CalcAndDealPeriodicHealing(sim, target, healPerTick, dot.OutcomeTick)
	}

	bloodCraze := warrior.RegisterSpell(config)

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Blood Craze - Damage Taken",
		Callback:           core.CallbackOnSpellHitTaken,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return result.Outcome.Matches(core.OutcomeCrit) || result.Damage > warrior.MaxHealth()*hitThreshold
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bloodCraze.SelfHot().Apply(sim)
		},
	})

	warrior.MakeProcTriggerAura(core.ProcTrigger{
		Name:               "Blood Craze - Bloodthirst",
		Callback:           core.CallbackOnSpellHitDealt,
		ClassSpellMask:     SpellMaskBloodthirst,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bloodCraze.SelfHot().Apply(sim)
		},
	})
}

// Raging Blows (1310315): the Cleave discount here, the off-hand Whirlwind strike in whirlwind.go.
func (warrior *Warrior) registerRagingBlows() {
	if !warrior.Talents.RagingBlows {
		return
	}

	spelldata.ParseStatic(&warrior.Character, spellData.RagingBlows.Rank(1))
}

var deathWishRank = spellData.DeathWish.Highest()

func (warrior *Warrior) registerDeathWish() {
	if !warrior.Talents.DeathWish {
		return
	}

	// The damage done effect carries the physical school mask, the damage taken one all schools.
	deathWishAura := warrior.RegisterAura(spelldata.AuraConfig(deathWishRank))
	spelldata.ParseEffects(&warrior.Character, deathWishAura, deathWishRank)

	// Grants immunity to Fear effects, which the row states as A_MECHANIC_IMMUNITY and the parse
	// skips.
	deathWishAura.AttachFearImmunity()

	config := spelldata.SpellConfig(&warrior.Unit, deathWishRank,
		spelldata.Flags(core.SpellFlagCastWhileIncapacitated))
	config.ClassSpellMask = SpellMaskDeathWish

	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		deathWishAura.Activate(sim)
		warrior.WaitUntil(sim, sim.CurrentTime+core.GCDDefault)
	}

	config.RelatedSelfBuff = deathWishAura

	deathWishSpell := warrior.RegisterSpell(config)

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: deathWishSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (warrior *Warrior) registerImprovedIntercept() {
	if warrior.Talents.ImprovedIntercept == 0 {
		return
	}

	spelldata.ParseStatic(&warrior.Character, spellData.ImprovedIntercept.Rank(warrior.Talents.ImprovedIntercept))
}

// Improved Cleave (12329) states only a rage discount on Cleave.
func (warrior *Warrior) registerImprovedCleave() {
	if warrior.Talents.ImprovedCleave == 0 {
		return
	}
	spelldata.ParseStatic(&warrior.Character, spellData.ImprovedCleave.Rank(warrior.Talents.ImprovedCleave))
}
