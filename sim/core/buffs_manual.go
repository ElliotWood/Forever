// Package-local home of the buff drivers the generator calls by name. A
// generated apply block for a cooldown, proc, uptime, item-count or manual row,
// or for a row the manifest marks as driven, calls drive<Go>(unit, scope);
// declaring that function here is what makes the row compile. A buff-scope row
// hands the driver the *Character the buffs are applied to and the message its
// field lives on, a debuff-scope row the *Unit the debuffs land on plus the
// debuffs and the raid. The driver reads its own field out of the message, and
// any other field it needs: Grace of Air lasts 9 seconds while the party is
// twisting totems.
package core

import (
	"time"

	"github.com/wowsims/forever/sim/core/proto"
)

// Spell 29166 states 100% mana regen while casting (aura 134) and +400% of it
// (aura 110); neither is a stat or a pseudo-stat, so the regen is the driver's.
const innervateSpiritRegenMultiplier = 5.0

// The party's Battle Shout is the external caster's copy, which chains behind
// the player's own shout rather than being up from the start.
func driveBattleShout(char *Character, _ *proto.PartyBuffs) {
	ApplyFixedShoutAura(char, BattleShoutAura(&char.Unit, false, 0), BattleShoutCategory)
}

// The party's Commanding Shout chains the same way its sibling does.
func driveCommandingShout(char *Character, _ *proto.PartyBuffs) {
	ApplyFixedShoutAura(char, CommandingShoutAura(&char.Unit, false, 0), CommandingShoutCategory)
}

// A druid innervates a character who is nearly out of mana, so that every other
// mana cooldown is spent first. The aura forces full spirit regen while it is
// up and the metrics record what the character gains from it.
func driveInnervates(char *Character, individual *proto.IndividualBuffs) {
	aura := InnervatesAura(&char.Unit, false, 0)
	manaMetrics := char.NewManaMetrics(aura.ActionID)

	threshold, expectedMana := 0.0, 0.0
	char.Env.RegisterPostFinalizeEffect(func() {
		threshold = innervateManaThreshold(char)
		expectedMana = char.SpiritManaRegenPerSecond() * innervateSpiritRegenMultiplier * aura.Duration.Seconds()
	})

	const ticks = 10
	aura.ApplyOnGain(func(aura *Aura, sim *Simulation) {
		char.PseudoStats.ForceFullSpiritRegen = true
		char.PseudoStats.SpiritRegenMultiplier *= innervateSpiritRegenMultiplier
		char.UpdateManaRegenRates()

		perTick := expectedMana / ticks
		StartPeriodicAction(sim, PeriodicActionOptions{
			Period:   aura.Duration / ticks,
			NumTicks: ticks,
			OnAction: func(sim *Simulation) {
				manaMetrics.AddEvent(perTick, perTick)
			},
		})
	}).ApplyOnExpire(func(aura *Aura, sim *Simulation) {
		char.PseudoStats.ForceFullSpiritRegen = false
		char.PseudoStats.SpiritRegenMultiplier /= innervateSpiritRegenMultiplier
		char.UpdateManaRegenRates()
	})

	newGeneratedExternalCD(char, aura, GeneratedExternalCD{
		NumSources: individual.Innervates,
		Cooldown:   InnervatesCooldown(),
		Type:       CooldownTypeMana,
		ShouldActivate: func(_ *Simulation, char *Character) bool {
			return char.CurrentMana() <= threshold
		},
	})
}

// A mage burns mana fast enough that waiting for a flat thousand left would
// waste most of the innervate.
func innervateManaThreshold(char *Character) float64 {
	if char.Class == proto.Class_ClassMage {
		return char.MaxMana() * 0.4
	}
	return 1000
}

// The priest has nothing to hold Power Infusion for, so it goes out on cooldown.
func drivePowerInfusions(char *Character, individual *proto.IndividualBuffs) {
	newGeneratedExternalCD(char, PowerInfusionsAura(&char.Unit, false, 0), GeneratedExternalCD{
		NumSources: individual.PowerInfusions,
		Cooldown:   PowerInfusionsCooldown(),
		Type:       CooldownTypeDPS,
	})
}

// A restoration shaman drops Mana Tide once the party has mana to refill, which
// is 40 seconds in, or halfway through a fight shorter than that.
func driveManaTideTotems(char *Character, party *proto.PartyBuffs) {
	initialDelay := time.Duration(0)
	char.Env.RegisterPostFinalizeEffect(func() {
		initialDelay = min(char.Env.BaseDuration/2, time.Second*40)
	})

	newGeneratedExternalCD(char, ManaTideTotemsAura(&char.Unit, false, 0), GeneratedExternalCD{
		NumSources: party.ManaTideTotems,
		Cooldown:   ManaTideTotemsCooldown(),
		Type:       CooldownTypeMana,
		ShouldActivate: func(sim *Simulation, _ *Character) bool {
			return sim.CurrentTime >= initialDelay
		},
	})
}

// The totem's aura is the attack power a windfury proc grants; what the client
// does not state is the proc itself, so the driver keeps the 20% chance on a
// main-hand swing, the 1.5 second internal cooldown and the extra attack the
// proc lands, and the totem aura that holds the category.
func driveWindfuryTotem(char *Character, _ *proto.PartyBuffs) {
	procAura := WindfuryTotemAura(&char.Unit, false, 0)

	var windfurySpell *Spell
	procTrigger := char.MakeProcTriggerAura(ProcTrigger{
		Name:               "Windfury Totem Trigger",
		MetricsActionID:    ActionID{SpellID: 25580, Tag: -1},
		IsWeaponProc:       true,
		ProcChance:         0.2,
		Duration:           NeverExpires,
		Outcome:            OutcomeLanded,
		Callback:           CallbackOnSpellHitDealt,
		ProcMask:           ProcMaskMeleeMHAuto,
		ICD:                time.Millisecond * 1500,
		TriggerImmediately: true,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			procAura.Activate(sim)
			char.AutoAttacks.MaybeReplaceMHSwing(sim, windfurySpell).Cast(sim, result.Target)
		},
	})

	// The totem stands for 10 seconds and the shaman drops a new one every 5,
	// so the aura that holds the category is simply refreshed.
	totemAura := char.GetOrRegisterAura(Aura{
		Label:    "Windfury Totem",
		ActionID: ActionID{SpellID: 25587, Tag: -1},
		Duration: time.Second * 10,
	}).ApplyOnInit(func(aura *Aura, sim *Simulation) {
		config := *char.AutoAttacks.MHConfig()
		config.ActionID = config.ActionID.WithTag(25584)
		windfurySpell = char.GetOrRegisterSpell(config)
	}).ApplyOnReset(func(aura *Aura, sim *Simulation) {
		aura.Activate(sim)
		StartPeriodicAction(sim, PeriodicActionOptions{
			Period:   time.Second * 5,
			Priority: ActionPriorityAuto,
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
			},
		})
	})

	totemAura.NewExclusiveEffect(WindfuryTotemCategory, false, ExclusiveEffect{
		Priority: WindfuryTotemValue(0),
		OnGain: func(_ *ExclusiveEffect, sim *Simulation) {
			procTrigger.Activate(sim)
		},
		OnExpire: func(_ *ExclusiveEffect, sim *Simulation) {
			procTrigger.Deactivate(sim)
			totemAura.Deactivate(sim)
		},
	})

	char.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}, func(sim *Simulation, slot proto.ItemSlot) {
		totemAura.Deactivate(sim)
	})
}

// A shaman twisting totems keeps Grace of Air up for 9 seconds out of every 10,
// because the air slot is holding another totem the rest of the time; a shaman
// who is not twisting leaves it standing.
func driveGraceOfAirTotem(char *Character, party *proto.PartyBuffs) {
	aura := GraceOfAirTotemAura(&char.Unit, false, 0)

	if !party.TotemTwisting {
		MakePermanent(aura)
		return
	}

	aura.Duration = time.Second * 9
	aura.ApplyOnReset(func(aura *Aura, sim *Simulation) {
		StartPeriodicAction(sim, PeriodicActionOptions{
			Period:          time.Second * 10,
			TickImmediately: true,
			Priority:        ActionPriorityAuto,
			OnAction: func(sim *Simulation) {
				aura.Activate(sim)
			},
		})
	})
}

// The staff's aura is worth its amounts once per Atiesh in the party.
func driveAtieshDruid(char *Character, party *proto.PartyBuffs) {
	MakePermanent(AtieshDruidAura(&char.Unit, false, 0, float64(party.AtieshDruid)))
}

func driveAtieshMage(char *Character, party *proto.PartyBuffs) {
	MakePermanent(AtieshMageAura(&char.Unit, false, 0, float64(party.AtieshMage)))
}

func driveAtieshPriest(char *Character, party *proto.PartyBuffs) {
	MakePermanent(AtieshPriestAura(&char.Unit, false, 0, float64(party.AtieshPriest)))
}

func driveAtieshWarlock(char *Character, party *proto.PartyBuffs) {
	MakePermanent(AtieshWarlockAura(&char.Unit, false, 0, float64(party.AtieshWarlock)))
}

// The judgement the paladin leaves on the target heals whoever strikes it; the
// client's trigger spell 5373 is a dummy, so how much and how often is the
// driver's.
func driveJudgementOfLight(target *Unit, _ *proto.Debuffs, _ *proto.Raid) {
	healthMetrics := target.NewHealthMetrics(ActionID{SpellID: 20346})

	MakePermanent(JudgementOfLightAura(target, false, 0)).AttachProcTrigger(ProcTrigger{
		Name:     "Judgement of Light - Heal",
		Callback: CallbackOnSpellHitTaken,
		ProcMask: ProcMaskMelee,
		Outcome:  OutcomeLanded,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			if sim.Proc(0.5, "Judgement of Light - Heal") {
				spell.Unit.GainHealth(sim, 95.0, healthMetrics)
			}
		},
	})
}

// Judgement of Wisdom returns mana to whoever strikes the target, on the same
// terms: 1826 is a dummy, so the amount and the chance stay here. Melee claim
// it returns mana on a miss as well.
func driveJudgementOfWisdom(target *Unit, _ *proto.Debuffs, _ *proto.Raid) {
	actionID := ActionID{SpellID: 20355}

	MakePermanent(JudgementOfWisdomAura(target, false, 0)).AttachProcTrigger(ProcTrigger{
		Name:       "Judgement of Wisdom",
		ActionID:   actionID,
		ProcChance: 0.5,
		ProcMask:   ProcMaskDirect,
		Callback:   CallbackOnSpellHitTaken,
		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			if !spell.ProcMask.Matches(ProcMaskMeleeOrRanged) && !result.Landed() {
				return
			}

			unit := spell.Unit
			if !unit.HasManaBar() {
				return
			}
			if unit.JowManaMetrics == nil {
				unit.JowManaMetrics = unit.NewManaMetrics(actionID)
			}
			unit.AddMana(sim, 74.0, unit.JowManaMetrics)
		},
	})
}

// A stack of Sunder Armor is worth nothing until it is on the target, so the
// raid's copy is ramped to five over the first five global cooldowns, which is
// how long a warrior takes to stack it.
func driveSunderArmor(target *Unit, _ *proto.Debuffs, _ *proto.Raid) {
	aura := MakePermanent(SunderArmorAura(target, false, 0))

	ScheduledAura(aura, PeriodicActionOptions{
		Period:          GCDDefault,
		NumTicks:        5,
		TickImmediately: true,
		Priority:        ActionPriorityDOT,
		OnAction: func(sim *Simulation) {
			aura.Activate(sim)
			if aura.IsActive() {
				aura.AddStack(sim)
			}
		},
	})
}
