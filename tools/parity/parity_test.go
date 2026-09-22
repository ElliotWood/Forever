// Package parity is the switch gate's DPS check: every arena spec on master and on
// forever-next with the same talents, bonus stats, weapons and target. Run it through
// run.sh, which produces master's numbers first; on its own the test skips.
//
// With ARENA_OUT=<dir> it also writes forever-next's result for each spec as the arena's
// rawResult (tools/arena), one <spec>.json per spec, and runs without master's numbers:
//
//	ARENA_OUT=<dir> PARITY_ITERATIONS=2000 go test --tags=with_db ./tools/parity -run '^TestParity$' -count=1 -v
package parity

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"

	"github.com/wowsims/forever/sim"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"
	"github.com/wowsims/forever/sim/core/stats"
)

type paritySpec struct {
	Name     string  `json:"name"`
	Race     string  `json:"race"`
	Profile  string  `json:"profile"`
	Talents  string  `json:"talents"`
	NextApl  string  `json:"nextApl"`
	MH       int32   `json:"mh"`
	OH       int32   `json:"oh"`
	Ranged   int32   `json:"ranged"`
	Distance float64 `json:"distance"`
	Tank     bool    `json:"tank"`
	Gear     string  `json:"gear"` // the spec's default gear preset
}

type parityResult struct {
	Dps     float64            `json:"dps"`
	Stats   map[string]float64 `json:"stats"`
	Actions map[string]float64 `json:"actions"`
	Casts   map[string]float64 `json:"casts"`
	Oom     float64            `json:"oom"`
	Dtps    float64            `json:"dtps"`
	Error   string             `json:"error,omitempty"`
	// For ARENA_OUT only: damage per spell id, and auto attacks.
	Damage       map[string]float64 `json:"-"`
	WeaponDamage float64            `json:"-"`
}

// tools/arena's rawResult, the shape master's sim/arenalib writes.
type arenaResult struct {
	Spec         string             `json:"spec"`
	Talents      string             `json:"talents"`
	Build        string             `json:"build"`
	Gear         string             `json:"gear"`
	Rotation     string             `json:"rotation"`
	Consumables  string             `json:"consumables"`
	Dps          float64            `json:"dps"`
	Damage       map[string]float64 `json:"damage"`
	WeaponDamage float64            `json:"weaponDamage"`
}

func specOptions(name string) (proto.Class, interface{}, *proto.ConsumesSpec) {
	none := &proto.ConsumesSpec{}
	switch name {
	case "balance_druid":
		return proto.Class_ClassDruid, &proto.Player_BalanceDruid{BalanceDruid: &proto.BalanceDruid{Options: &proto.BalanceDruid_Options{
			ClassOptions: &proto.DruidOptions{}}}}, none
	case "feral_druid":
		return proto.Class_ClassDruid, &proto.Player_FeralCatDruid{FeralCatDruid: &proto.FeralCatDruid{
			Rotation: &proto.FeralCatDruid_Rotation{}, Options: &proto.FeralCatDruid_Options{}}}, none
	case "feral_tank_druid":
		return proto.Class_ClassDruid, &proto.Player_FeralBearDruid{FeralBearDruid: &proto.FeralBearDruid{Options: &proto.FeralBearDruid_Options{
			StartingRage: 20}}}, none
	case "hunter":
		return proto.Class_ClassHunter, &proto.Player_Hunter{Hunter: &proto.Hunter{Options: &proto.Hunter_Options{ClassOptions: &proto.HunterOptions{
			Ammo: proto.HunterOptions_ThoriumHeadedArrow, PetType: proto.HunterOptions_Cat, PetUptime: 1,
			PetAttackSpeed: proto.HunterOptions_OneTwo}}}}, none
	case "mage":
		return proto.Class_ClassMage, &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{ClassOptions: &proto.MageOptions{
			DefaultMageArmor: proto.MageArmor_MageArmorMageArmor}}}}, none
	case "protection_paladin":
		return proto.Class_ClassPaladin, &proto.Player_ProtectionPaladin{ProtectionPaladin: &proto.ProtectionPaladin{Options: &proto.ProtectionPaladin_Options{
			ClassOptions: &proto.PaladinOptions{}}}}, none
	case "retribution_paladin":
		return proto.Class_ClassPaladin, &proto.Player_RetributionPaladin{RetributionPaladin: &proto.RetributionPaladin{Options: &proto.RetributionPaladin_Options{
			ClassOptions: &proto.PaladinOptions{}}}}, none
	case "shadow_priest", "smite_priest":
		return proto.Class_ClassPriest, &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{ClassOptions: &proto.PriestOptions{
			Armor: proto.PriestOptions_InnerFire}}}}, none
	case "rogue":
		return proto.Class_ClassRogue, &proto.Player_Rogue{Rogue: &proto.Rogue{Options: &proto.Rogue_Options{ClassOptions: &proto.RogueOptions{}}}},
			&proto.ConsumesSpec{MhImbueId: 26891, OhImbueId: 27186} // Instant, Deadly Poison
	case "elemental_shaman":
		return proto.Class_ClassShaman, &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{
			ClassOptions: &proto.ShamanOptions{}}}}, none
	case "enhancement_shaman":
		return proto.Class_ClassShaman, &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{
			SyncType: proto.ShamanSyncType_Auto, ImbueOh: proto.ShamanImbue_WindfuryWeapon,
			ClassOptions: &proto.ShamanOptions{ImbueMh: proto.ShamanImbue_WindfuryWeapon}}}}, none
	case "warlock":
		return proto.Class_ClassWarlock, &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{ClassOptions: &proto.WarlockOptions{
			Armor: proto.WarlockOptions_DemonArmor, Summon: proto.WarlockOptions_Succubus}}}}, none
	case "warrior":
		return proto.Class_ClassWarrior, &proto.Player_DpsWarrior{DpsWarrior: &proto.DpsWarrior{Options: &proto.DpsWarrior_Options{ClassOptions: &proto.WarriorOptions{
			DefaultShout: proto.WarriorShout_WarriorShoutBattle, DefaultStance: proto.WarriorStance_WarriorStanceBerserker}}}}, none
	case "tank_warrior":
		return proto.Class_ClassWarrior, &proto.Player_ProtectionWarrior{ProtectionWarrior: &proto.ProtectionWarrior{Options: &proto.ProtectionWarrior_Options{ClassOptions: &proto.WarriorOptions{
			DefaultShout: proto.WarriorShout_WarriorShoutBattle, DefaultStance: proto.WarriorStance_WarriorStanceDefensive}}}}, none
	}
	return proto.Class_ClassUnknown, nil, none
}

// The profile is in percent for crit and hit; this engine takes rating.
func bonusStats(p map[string]float64) stats.Stats {
	return stats.Stats{
		stats.Strength:          p["str"],
		stats.Agility:           p["agi"],
		stats.Stamina:           p["sta"],
		stats.Intellect:         p["int"],
		stats.Spirit:            p["spi"],
		stats.AttackPower:       p["ap"],
		stats.RangedAttackPower: p["rap"],
		stats.SpellDamage:       p["sd"],
		stats.MeleeCritRating:   p["crit"] * core.PhysicalCritRatingPerCritPercent,
		stats.SpellCritRating:   p["crit"] * core.SpellCritRatingPerCritPercent,
		stats.MeleeHitRating:    p["hit"] * core.PhysicalHitRatingPerHitPercent,
		stats.SpellHitRating:    p["hit"] * core.SpellHitRatingPerHitPercent,
	}
}

func TestParity(t *testing.T) {
	masterFile := os.Getenv("PARITY_MASTER")
	arenaOut := os.Getenv("ARENA_OUT")
	if masterFile == "" && arenaOut == "" {
		t.Skip("run through tools/parity/run.sh")
	}
	sim.RegisterAll()
	var file struct {
		Profiles map[string]map[string]float64 `json:"profiles"`
		Specs    []paritySpec                  `json:"specs"`
	}
	mustReadJSON(t, "specs.json", &file)
	master := map[string]parityResult{}
	if masterFile != "" {
		mustReadJSON(t, masterFile, &master)
	}
	iterations, _ := strconv.Atoi(os.Getenv("PARITY_ITERATIONS"))

	fmt.Printf("\n%-20s %9s %9s %8s   %s\n", "spec", "master", "next", "gap", "final stats master | next: ap sd crit%(melee/spell) hit% str/agi/int")
	for _, spec := range file.Specs {
		m := master[spec.Name]
		n := runSpec(spec, file.Profiles[spec.Profile], int32(iterations))
		if arenaOut != "" {
			writeArena(t, arenaOut, spec, n)
		}
		gap := "-"
		if m.Error == "" && n.Error == "" && m.Dps > 0 {
			gap = fmt.Sprintf("%+.1f%%", (n.Dps/m.Dps-1)*100)
		}
		detail := statsLine(m) + " | " + statsLine(n)
		if m.Error != "" {
			detail = "master: " + m.Error
		}
		if n.Error != "" {
			detail = "next: " + n.Error
		}
		fmt.Printf("%-20s %9.1f %9.1f %8s   %s\n", spec.Name, m.Dps, n.Dps, gap, detail)
		if os.Getenv("PARITY_DETAIL") != "" {
			printActions(m, n)
		}
	}
}

// Per-action DPS side by side, biggest first. Spell ids can differ between the two
// engines' clients, so one ability can show up as two rows.
func printActions(mr, nr parityResult) {
	m, n := mr.Actions, nr.Actions
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	for k := range n {
		if _, ok := m[k]; !ok {
			keys = append(keys, k)
		}
	}
	sort.Slice(keys, func(i, j int) bool { return m[keys[i]]+n[keys[i]] > m[keys[j]]+n[keys[j]] })
	for _, k := range keys {
		fmt.Printf("    %-28s %9.1f %9.1f   casts %6.1f %6.1f\n", k, m[k], n[k], mr.Casts[k], nr.Casts[k])
	}
}

func statsLine(r parityResult) string {
	s := r.Stats
	return fmt.Sprintf("%.0f %.0f %.1f/%.1f %.1f/%.1f %.0f/%.0f/%.0f oom %.0fs dtps %.0f", math.Max(s["ap"], s["rap"]), s["sd"],
		s["mcrit"], s["scrit"], s["mhit"], s["shit"], s["str"], s["agi"], s["int"], r.Oom, r.Dtps)
}

func runSpec(spec paritySpec, profile map[string]float64, iterations int32) parityResult {
	items := make([]*proto.ItemSpec, len(proto.ItemSlot_name))
	for i := range items {
		items[i] = &proto.ItemSpec{}
	}
	items[proto.ItemSlot_ItemSlotMainHand].Id = spec.MH
	items[proto.ItemSlot_ItemSlotOffHand].Id = spec.OH
	items[proto.ItemSlot_ItemSlotRanged].Id = spec.Ranged
	return runSpecWithGear(spec, profile, iterations, &proto.EquipmentSpec{Items: items})
}

func runSpecWithGear(spec paritySpec, profile map[string]float64, iterations int32, equipment *proto.EquipmentSpec) (res parityResult) {
	defer func() {
		if r := recover(); r != nil {
			res = parityResult{Error: fmt.Sprintf("panic: %.160v", r)}
		}
	}()
	class, options, consumables := specOptions(spec.Name)
	aplPath := filepath.Join("..", "..", spec.NextApl)
	if _, err := os.Stat(aplPath + ".apl.json"); err != nil {
		return parityResult{Error: "no APL " + spec.NextApl}
	}

	distance := spec.Distance
	if distance == 0 {
		distance = 5
	}
	player := core.WithSpec(&proto.Player{
		Class:              class,
		Race:               proto.Race(proto.Race_value[spec.Race]),
		Equipment:          equipment,
		Consumables:        consumables,
		BonusStats:         &proto.UnitStats{Stats: bonusStats(profile).ToProtoArray()},
		TalentsString:      spec.Talents,
		Rotation:           core.GetAplRotation(filepath.Dir(aplPath), filepath.Base(aplPath)).Rotation,
		DistanceFromTarget: distance,
		ReactionTimeMs:     150,
		ChannelClipDelayMs: 50,
	}, options)
	raid := core.SinglePlayerRaidProto(player, &proto.PartyBuffs{}, &proto.RaidBuffs{}, &proto.Debuffs{})
	if spec.Tank {
		raid.Tanks = []*proto.UnitReference{{Type: proto.UnitReference_Player, Index: 0}}
	}
	// Our Forever "Level 60" preset: master's target, stated in full so neither engine's
	// test-dummy defaults leak in.
	encounter := &proto.Encounter{
		Duration:             180,
		ExecuteProportion_20: 0.2,
		ExecuteProportion_25: 0.25,
		ExecuteProportion_35: 0.35,
		// This engine walks 90 -> 45 -> 35 -> 25 -> 20; an unset 90 stalls it at the start.
		ExecuteProportion_45: 0.45,
		ExecuteProportion_90: 0.9,
		Targets: []*proto.Target{{
			Level:         63,
			MobType:       proto.MobType_MobTypeUnknown,
			Stats:         stats.Stats{stats.Armor: 3731, stats.AttackPower: 805}.ToProtoArray(),
			SwingSpeed:    2,
			MinBaseDamage: 3000,
			DamageSpread:  0.3333,
			ParryHaste:    true,
		}},
	}

	final := core.ComputeStats(&proto.ComputeStatsRequest{Raid: raid, Encounter: encounter})
	if final.ErrorResult != "" {
		return parityResult{Error: fmt.Sprintf("%.160s", final.ErrorResult)}
	}
	s := stats.FromUnitStatsProto(final.RaidStats.Parties[0].Players[0].FinalStats)

	result := core.RunRaidSim(&proto.RaidSimRequest{
		Raid:       raid,
		Encounter:  encounter,
		SimOptions: &proto.SimOptions{Iterations: iterations, IsTest: true, RandomSeed: 101},
	})
	if result.Error != nil {
		return parityResult{Error: fmt.Sprintf("%.160s", result.Error.Message)}
	}
	if debugHook != nil {
		debugHook(result)
	}
	actions, casts := actionDps(result.RaidMetrics.Parties[0].Players[0], result.RaidMetrics.Dps.Avg, iterations)
	damage, weapon := spellDamage(result.RaidMetrics.Parties[0].Players[0])
	return parityResult{
		Damage:       damage,
		WeaponDamage: weapon,
		Dps:          result.RaidMetrics.Dps.Avg,
		Actions:      actions,
		Casts:        casts,
		Oom:          result.RaidMetrics.Parties[0].Players[0].SecondsOomAvg,
		Dtps:         result.RaidMetrics.Parties[0].Players[0].Dtps.Avg,
		Stats: map[string]float64{
			"ap": s[stats.AttackPower], "rap": s[stats.RangedAttackPower], "sd": s[stats.SpellDamage],
			"mcrit": s[stats.PhysicalCritPercent], "scrit": s[stats.SpellCritPercent],
			"mhit": s[stats.PhysicalHitPercent], "shit": s[stats.SpellHitPercent],
			"int": s[stats.Intellect], "agi": s[stats.Agility], "str": s[stats.Strength], "spi": s[stats.Spirit],
		},
	}
}

var debugHook func(*proto.RaidSimResult)

// One spec's result as the arena merger reads it. The build is the parity profile, not a gear
// set: bonus stats and statless weapons, no consumables beyond what the class grants itself.
func writeArena(t *testing.T, dir string, spec paritySpec, r parityResult) {
	consumables := "No consumables"
	if spec.Name == "rogue" {
		consumables += "+class"
	}
	row := arenaResult{
		Spec:         spec.Name,
		Talents:      spec.Talents,
		Build:        "parity build",
		Gear:         "parity " + spec.Profile + " profile",
		Rotation:     filepath.Base(spec.NextApl),
		Consumables:  consumables,
		Dps:          r.Dps,
		Damage:       r.Damage,
		WeaponDamage: r.WeaponDamage,
	}
	if row.Damage == nil {
		row.Damage = map[string]float64{}
	}
	body, err := json.MarshalIndent([]arenaResult{row}, "", "	")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, spec.Name+".json"), body, 0o644); err != nil {
		t.Fatal(err)
	}
}

// Damage per spell, summed over targets, the player's and its pets', as master's arenalib
// collects it: auto attacks are weapon damage, a spell-less item effect is spell "0".
func spellDamage(unit *proto.UnitMetrics) (map[string]float64, float64) {
	damage, weapon := map[string]float64{}, 0.0
	for _, u := range append([]*proto.UnitMetrics{unit}, unit.Pets...) {
		for _, action := range u.Actions {
			total := 0.0
			for _, target := range action.Targets {
				total += target.Damage
			}
			if total <= 0 {
				continue
			}
			if id := action.Id.GetSpellId(); id != 0 {
				damage[fmt.Sprint(id)] += total
			} else if action.Id.GetOtherId() != proto.OtherAction_OtherActionNone {
				weapon += total
			} else {
				damage["0"] += total
			}
		}
	}
	return damage, weapon
}

func mustReadJSON(t *testing.T, path string, v interface{}) {
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, v); err != nil {
		t.Fatal(err)
	}
}

// DPS and casts a fight per action, the player's and its pets', keyed so both engines
// name them alike.
func actionDps(unit *proto.UnitMetrics, dps float64, iterations int32) (map[string]float64, map[string]float64) {
	out := map[string]float64{}
	casts := map[string]float64{}
	total := 0.0
	units := append([]*proto.UnitMetrics{unit}, unit.Pets...)
	for i, u := range units {
		for _, action := range u.Actions {
			damage := 0.0
			for _, target := range action.Targets {
				damage += target.Damage
			}
			if damage <= 0 {
				continue
			}
			key := fmt.Sprintf("spell %d", action.Id.GetSpellId())
			if action.Id.GetSpellId() == 0 {
				key = action.Id.GetOtherId().String()
			}
			if action.Id.GetTag() != 0 {
				key += fmt.Sprintf(" tag %d", action.Id.GetTag())
			}
			if i > 0 {
				key = "pet " + key
			}
			out[key] += damage
			total += damage
			for _, target := range action.Targets {
				casts[key] += float64(target.Casts) / float64(iterations)
			}
		}
	}
	for k := range out {
		out[k] *= dps / total
	}
	return out, casts
}
