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
