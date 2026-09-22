//go:build with_db

// The proc resolver against the item procs the sim registers today: every live registration in
// sim/common/forever/stat_bonus_procs_auto_gen.go, each compared with the rate the registration
// actually fires at, so that switching the generator over to the resolver can be shown to change
// only what this file says it changes.

package spelldata

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/wowsims/forever/assets/database"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
)

const procParityExpectedPath = "../../../tools/database/proc_parity_expected.txt"

// One live registration as the generated file states it. The listener fields are always the
// literals the generator wrote; where the rate and the internal cooldown are stated depends on the
// constructor, which readsRateFromEntry below says.
//
// buffSpellID is the effect entry's BuffId, which is the spell granting the stats or dealing the
// damage. The spell carrying the proc flags is a different one, and the entry names it in BuffName.
type liveProc struct {
	itemID       int32
	name         string
	buffSpellID  int32
	isWeaponProc bool

	callback           core.AuraCallback
	procMask           core.ProcMask
	outcome            core.HitOutcome
	requireDamageDealt bool
	canProcFromProcs   bool

	// The rate and the lockout the registration fires at today. A damage proc carries both as
	// literals in the generated call; a stat proc reads them off the item effect entry at run time,
	// which is what readsRateFromEntry says.
	generatedChance    float64
	generatedICDMs     int32
	readsRateFromEntry bool
}

const meleeAllMask = core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto |
	core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial

const shieldSpikeMask = meleeAllMask | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial

func liveProcs() []liveProc {
	shieldSpike := func(itemID int32, name string) liveProc {
		return liveProc{itemID: itemID, name: name, buffSpellID: 16782,
			callback: core.CallbackOnSpellHitTaken, procMask: shieldSpikeMask,
			outcome: core.OutcomeLanded, requireDamageDealt: true, generatedChance: 0.05}
	}
	fieryWeapon := func(itemID int32, name string, buff int32) liveProc {
		return liveProc{itemID: itemID, name: name, buffSpellID: buff,
			callback: core.CallbackOnSpellHitDealt, procMask: meleeAllMask,
			outcome: core.OutcomeLanded, requireDamageDealt: true, generatedChance: 1}
	}

	return []liveProc{
		fieryWeapon(12631, "Fiery Plate Gauntlets", 7714),
		fieryWeapon(12632, "Storm Gauntlets", 16614),
		fieryWeapon(17111, "Blazefury Medallion", 7712),
		shieldSpike(18825, "Grand Marshal's Aegis"),
		shieldSpike(18826, "High Warlord's Shield Wall"),
		shieldSpike(234562, "High Warlord's Shield Wall (reissue)"),
		shieldSpike(234588, "Grand Marshal's Aegis (reissue)"),
		shieldSpike(272591, "Premier High Warlord's Shield Wall"),
		shieldSpike(272838, "Premier Grand Marshal's Aegis"),
		{
			itemID: 260205, name: "Highborne Research Tablet", buffSpellID: 1265634,
			callback: core.CallbackOnSpellHitDealt,
			procMask: shieldSpikeMask | core.ProcMaskSpellDamage,
			outcome:  core.OutcomeLanded, generatedChance: 1,
		},
		{
			itemID: 285278, name: "Satchel of Dark Iron Bombs", buffSpellID: 1318121,
			callback: core.CallbackOnSpellHitDealt, procMask: core.ProcMaskRangedAuto,
			outcome: core.OutcomeLanded, requireDamageDealt: true, generatedChance: 1,
		},
		{
			itemID: 12798, name: "Annihilator", buffSpellID: 16928, isWeaponProc: true,
			callback: core.CallbackOnSpellHitDealt, procMask: core.ProcMaskUnknown,
			outcome: core.OutcomeLanded, requireDamageDealt: true, readsRateFromEntry: true,
		},
		{
			itemID: 19288, name: "Darkmoon Card: Blue Dragon", buffSpellID: 23684,
			callback: core.CallbackOnCastComplete, procMask: core.ProcMaskSpellDamage,
			outcome: core.OutcomeEmpty, readsRateFromEntry: true,
		},
		{
			itemID: 21190, name: "Wrath of Cenarius", buffSpellID: 25907,
			callback: core.CallbackOnCastComplete, procMask: core.ProcMaskSpellDamage,
			outcome: core.OutcomeEmpty, canProcFromProcs: true, readsRateFromEntry: true,
		},
		{
			itemID: 249473, name: "Dormant Heart of the Mountain", buffSpellID: 1249119,
			callback: core.CallbackOnHealDealt, procMask: core.ProcMaskSpellHealing,
			outcome: core.OutcomeLanded, canProcFromProcs: true, readsRateFromEntry: true,
		},
		{
			itemID: 275630, name: "Depleted Eye of Influence", buffSpellID: 1297082,
			callback: core.CallbackOnCastComplete, procMask: core.ProcMaskSpellDamage,
			outcome: core.OutcomeEmpty, readsRateFromEntry: true,
		},
	}
}

func TestProcParityWithGeneratedItemProcs(t *testing.T) {
	withGeneratedStore(t)
	expected := loadProcParityExpectations(t)

	for _, live := range liveProcs() {
		effect := procEffectFor(t, live.itemID, live.buffSpellID)
		if effect == nil {
			continue
		}

		trigger := resolvedTrigger(t, live, effect)
		if trigger == nil {
			continue
		}

		proc := effect.GetProc()
		t.Logf("%s (%d) trigger %d: callback %d mask %d outcome %d damage-dealt %v procs %v chance %v icd %v",
			live.name, live.itemID, triggerSpellID(effect), trigger.Callback, trigger.ProcMask,
			trigger.Outcome, trigger.RequireDamageDealt, trigger.CanProcFromProcs, trigger.ProcChance, trigger.ICD)

		report := func(field string, got any, want any) {
			expected.check(t, live, field, fmt.Sprintf("%v, generated %v", got, want))
		}

		if trigger.Callback != live.callback {
			report("Callback", trigger.Callback, live.callback)
		}
		if trigger.ProcMask != live.procMask {
			report("ProcMask", trigger.ProcMask, live.procMask)
		}
		if trigger.Outcome != live.outcome {
			report("Outcome", trigger.Outcome, live.outcome)
		}
		if trigger.RequireDamageDealt != live.requireDamageDealt {
			report("RequireDamageDealt", trigger.RequireDamageDealt, live.requireDamageDealt)
		}
		if trigger.CanProcFromProcs != live.canProcFromProcs {
			report("CanProcFromProcs", trigger.CanProcFromProcs, live.canProcFromProcs)
		}
		chance, icd := live.generatedChance, time.Millisecond*time.Duration(live.generatedICDMs)
		if live.readsRateFromEntry {
			chance, icd = proc.GetProcChance(), time.Millisecond*time.Duration(proc.IcdMs)
		}

		if trigger.ProcChance != chance {
			report("ProcChance", trigger.ProcChance, chance)
		}
		if trigger.ICD != icd {
			report("ICD", trigger.ICD, icd)
		}
		// A rate stated as procs per minute is a manager and no roll at all, so both halves are
		// asserted: a manager with a chance beside it would be gated by the chance and never
		// measured, which is the shape a bare "is it non-nil" check would not catch.
		if proc.GetPpm() > 0 {
			if trigger.DPM == nil {
				report("DPM", "no proc manager", fmt.Sprintf("%v procs per minute", proc.GetPpm()))
			}
			if trigger.ProcChance != 0 {
				report("ProcChance", fmt.Sprintf("%v beside a proc manager", trigger.ProcChance), 0)
			}
		}

		// The buff the proc applies, for the registrations that build one. A damage proc casts a
		// spell instead, so only the stat procs have a duration and a stack count to disagree about.
		if !live.readsRateFromEntry {
			continue
		}

		buff := AuraConfig(Find(effect.BuffId))
		if buff.Duration != time.Millisecond*time.Duration(effect.EffectDurationMs) {
			report("Duration", buff.Duration, time.Millisecond*time.Duration(effect.EffectDurationMs))
		}
		if buff.MaxStacks != effect.MaxCumulativeStacks {
			report("MaxStacks", buff.MaxStacks, effect.MaxCumulativeStacks)
		}
	}

	expected.requireAllUsed(t)
}

// Every proc-carrying item effect the database ships, not just the ones the generator emits a live
// registration for: a generator reading the resolver has to get an answer for all of them, so the
// verdicts are counted here rather than discovered one item at a time.
func TestProcParitySweepOverEveryItemProc(t *testing.T) {
	withGeneratedStore(t)

	counts := map[string]int{}
	effects := 0
	for _, item := range database.Load().Items {
		for _, effect := range item.ItemEffects {
			if effect.GetProc() == nil {
				continue
			}
			effects++

			verdict := sweepVerdict(effect)
			counts[verdict]++
			if verdict != "resolved" {
				t.Logf("%s: item %d %q, trigger %d", verdict, item.Id, item.Name, triggerSpellID(effect))
			}

			// The buff as the two sides describe it. The item effect entry takes the trigger's
			// duration where the buff states none, and the higher of the two stack counts, so a
			// buff read off its own row alone can differ from the entry on both.
			buff := AuraConfig(Find(effect.BuffId))
			if buff.Duration != time.Millisecond*time.Duration(effect.EffectDurationMs) {
				counts["buff duration differs from the entry"]++
				t.Logf("duration: item %d %q buff %d states %v, the entry states %v",
					item.Id, item.Name, effect.BuffId, buff.Duration, time.Millisecond*time.Duration(effect.EffectDurationMs))
			}
			if buff.MaxStacks != effect.MaxCumulativeStacks {
				counts["buff stacks differ from the entry"]++
				t.Logf("stacks: item %d %q buff %d states %d, the entry states %d",
					item.Id, item.Name, effect.BuffId, buff.MaxStacks, effect.MaxCumulativeStacks)
			}
		}
	}

	t.Logf("proc-carrying item effects: %d", effects)
	for verdict, count := range counts {
		t.Logf("%s: %d", verdict, count)
	}
	if effects == 0 {
		t.Fatal("no proc-carrying item effects found, so the sweep proves nothing")
	}
	if counts["no trigger row"] > 0 {
		t.Errorf("%d proc effects name a spell the store does not carry", counts["no trigger row"])
	}
}

// What the resolver makes of one item effect, as the class the sweep counts. A row with no listener
// or no rate is not a failure: those are the effects the generator refuses today as well, and only
// the count is pinned.
func sweepVerdict(effect *proto.ItemEffect) string {
	trigger := Find(triggerSpellID(effect))
	switch {
	case trigger == Nil:
		return "no trigger row"
	case core.DecodeProcTypeMask(trigger.ProcFlags, trigger.ProcHint).Callback == core.CallbackEmpty:
		return "no listener in the mask"
	case !statesARate(trigger, effect.GetProc().GetPpm()):
		return "no rate stated anywhere"
	}

	character := &core.Character{}
	config := ProcTrigger(character, trigger, nil, procRate(effect.GetProc()))
	if len(ProcTriggerUnsupported(character, trigger)) > 0 {
		return "resolved, with unsupported bits"
	}
	if config.ICD > 0 {
		return "resolved, with an internal cooldown"
	}

	return "resolved"
}

// Whether a rate can be found for the row at all. The resolver refuses a row that states none, so
// the sweep has to answer the same question before it asks for a trigger.
func statesARate(s *Spell, ppm float64) bool {
	if ppm > 0 || s.RPPM > 0 {
		return true
	}

	switch s.ProcChanceSource {
	case ProcChancePPM:
		return false
	case ProcChanceColumn:
		return s.ProcChance > 0
	}

	return true
}

// The trigger's own row, which is the spell the item effect names rather than the buff it applies.
// The database writes that spell's ID into BuffName - see makeBaseProto in tools/database/dbc - and
// it is the only place the effect entry keeps it.
var buffNameSpellID = regexp.MustCompile(`\((\d+)\)$`)

func triggerSpellID(effect *proto.ItemEffect) int32 {
	match := buffNameSpellID.FindStringSubmatch(effect.BuffName)
	if match == nil {
		return effect.BuffId
	}

	id, err := strconv.Atoi(match[1])
	if err != nil {
		return effect.BuffId
	}

	return int32(id)
}

// The trigger the resolver builds for a live registration, under the two options
// shared.NewSpellDataProc gives it: the on-hit shape a weapon proc takes instead of a proc mask, and
// the rate an item's procs per minute reach the sim through.
//
// The two options are copied here rather than called: sim/common/shared imports this package, so a
// test inside it cannot import shared back. Their rules are pinned by the shared package's own
// tests; what this file is for is the numbers the pair of them ends up producing.
func resolvedTrigger(t *testing.T, live liveProc, effect *proto.ItemEffect) *core.ProcTrigger {
	t.Helper()

	trigger := Find(triggerSpellID(effect))
	if trigger == Nil {
		t.Errorf("%s (%d): the store carries no row for trigger spell %d",
			live.name, live.itemID, triggerSpellID(effect))
		return nil
	}

	config := ProcTrigger(&core.Character{}, trigger, nil, weaponProcShape(live.isWeaponProc), procRate(effect.GetProc()))

	// The buff-category fallback NewSpellDataProc applies, for the procs that state their internal
	// cooldown on the buff rather than on the trigger.
	if config.ICD == 0 && effect.BuffId != trigger.ID {
		config.ICD = Find(effect.BuffId).CategoryCooldown()
	}

	return &config
}

func weaponProcShape(isWeaponProc bool) ProcOpt {
	return func(_ *core.Character, trigger *core.ProcTrigger) {
		if !isWeaponProc {
			return
		}

		trigger.Callback = core.CallbackOnSpellHitDealt
		trigger.ProcMask = core.ProcMaskUnknown
		trigger.RequireDamageDealt = true
		trigger.CanProcFromProcs = false
	}
}

func procRate(proc *proto.ProcEffect) ProcOpt {
	return func(character *core.Character, trigger *core.ProcTrigger) {
		if proc.GetPpm() <= 0 {
			return
		}

		trigger.ProcChance = 0
		trigger.DPM = character.NewLegacyPPMManager(proc.GetPpm(), trigger.ProcMask)
	}
}

func procEffectFor(t *testing.T, itemID int32, buffID int32) *proto.ItemEffect {
	t.Helper()

	item := core.GetItemByID(itemID)
	if item == nil {
		t.Errorf("item %d is not in the database", itemID)
		return nil
	}

	for _, effect := range item.ItemEffects {
		if effect.BuffId == buffID && effect.GetProc() != nil {
			return effect
		}
	}

	t.Errorf("item %d carries no proc effect for buff %d", itemID, buffID)
	return nil
}

// The differences the resolver is expected to have, by item and field. A difference the file does
// not name fails the test; a line the test never hits fails it too, so the file cannot go stale.
type procParityExpectations struct {
	reasons map[string]string
	used    map[string]bool
}

func (e procParityExpectations) check(t *testing.T, live liveProc, field string, detail string) {
	t.Helper()

	key := fmt.Sprintf("%d %s", live.itemID, field)
	reason, ok := e.reasons[key]
	if !ok {
		t.Errorf("%s (%d) %s: resolver says %s - add it to %s with a reason if it is intended",
			live.name, live.itemID, field, detail, procParityExpectedPath)
		return
	}

	e.used[key] = true
	t.Logf("%s (%d) %s: resolver says %s - expected: %s", live.name, live.itemID, field, detail, reason)
}

func (e procParityExpectations) requireAllUsed(t *testing.T) {
	t.Helper()

	for key := range e.reasons {
		if !e.used[key] {
			t.Errorf("%s expects a difference for %q that the resolver does not have any more - remove the line",
				procParityExpectedPath, key)
		}
	}
}

func loadProcParityExpectations(t *testing.T) procParityExpectations {
	t.Helper()

	file, err := os.Open(procParityExpectedPath)
	if err != nil {
		t.Fatalf("cannot read the expected differences: %v", err)
	}
	defer file.Close()

	expectations := procParityExpectations{reasons: map[string]string{}, used: map[string]bool{}}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		itemAndField, reason, found := strings.Cut(line, ":")
		if !found {
			t.Fatalf("expected difference %q is not in the \"<item id> <field>: <reason>\" shape", line)
		}

		expectations.reasons[strings.TrimSpace(itemAndField)] = strings.TrimSpace(reason)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("cannot read the expected differences: %v", err)
	}

	return expectations
}
