// Package buffmanifest is the checked-in census of every raid, party, individual
// and enemy-debuff proto field the sim knows about. It is a leaf package: it may
// import the standard library and sim/core/proto and nothing else, so the buff
// generator and the proto emitter can both read it while the generated proto and
// sim/core files are stale.
package buffmanifest

import (
	"strings"

	"github.com/wowsims/forever/sim/core/proto"
)

type BuffScope int

const (
	ScopeRaid BuffScope = iota
	ScopeParty
	ScopeIndividual
	ScopeDebuff
)

func (s BuffScope) String() string {
	switch s {
	case ScopeRaid:
		return "ScopeRaid"
	case ScopeParty:
		return "ScopeParty"
	case ScopeIndividual:
		return "ScopeIndividual"
	case ScopeDebuff:
		return "ScopeDebuff"
	}
	return "BuffScope(unknown)"
}

type BuffProtoType int

const (
	ProtoBool BuffProtoType = iota
	ProtoTristate
	ProtoInt32
	ProtoDouble
	ProtoEnumDrums
)

func (t BuffProtoType) String() string {
	switch t {
	case ProtoBool:
		return "ProtoBool"
	case ProtoTristate:
		return "ProtoTristate"
	case ProtoInt32:
		return "ProtoInt32"
	case ProtoDouble:
		return "ProtoDouble"
	case ProtoEnumDrums:
		return "ProtoEnumDrums"
	}
	return "BuffProtoType(unknown)"
}

type BuffKind int

const (
	KindStatFlat BuffKind = iota
	KindStatPct
	KindResistance
	KindPseudoMult
	KindDamageShield
	KindDebuffStat
	KindDebuffStacking
	KindDebuffDamageTaken
	KindDebuffAtkSpeed
	KindDebuffUptime
	KindExternalCD
	KindProc
	KindItemCount
	KindManual
	KindFlag
	KindEnum
	KindAbsent
)

func (k BuffKind) String() string {
	switch k {
	case KindStatFlat:
		return "KindStatFlat"
	case KindStatPct:
		return "KindStatPct"
	case KindResistance:
		return "KindResistance"
	case KindPseudoMult:
		return "KindPseudoMult"
	case KindDamageShield:
		return "KindDamageShield"
	case KindDebuffStat:
		return "KindDebuffStat"
	case KindDebuffStacking:
		return "KindDebuffStacking"
	case KindDebuffDamageTaken:
		return "KindDebuffDamageTaken"
	case KindDebuffAtkSpeed:
		return "KindDebuffAtkSpeed"
	case KindDebuffUptime:
		return "KindDebuffUptime"
	case KindExternalCD:
		return "KindExternalCD"
	case KindProc:
		return "KindProc"
	case KindItemCount:
		return "KindItemCount"
	case KindManual:
		return "KindManual"
	case KindFlag:
		return "KindFlag"
	case KindEnum:
		return "KindEnum"
	case KindAbsent:
		return "KindAbsent"
	}
	return "BuffKind(unknown)"
}

type TalentApplies int

const (
	TalentScalesValue TalentApplies = iota
	TalentScalesDuration
	TalentAddsStat
)

func (a TalentApplies) String() string {
	switch a {
	case TalentScalesValue:
		return "TalentScalesValue"
	case TalentScalesDuration:
		return "TalentScalesDuration"
	case TalentAddsStat:
		return "TalentAddsStat"
	}
	return "TalentApplies(unknown)"
}

type PetPolicy int

const (
	PetNormal PetPolicy = iota
	PetStrip
	PetInheritOwnerAura
	PetCapAtRegular
	// PetStripWhenSummonedLate drops the buff only for a pet that is not enabled
	// at fight start.
	PetStripWhenSummonedLate
)

func (p PetPolicy) String() string {
	switch p {
	case PetNormal:
		return "PetNormal"
	case PetStrip:
		return "PetStrip"
	case PetInheritOwnerAura:
		return "PetInheritOwnerAura"
	case PetCapAtRegular:
		return "PetCapAtRegular"
	case PetStripWhenSummonedLate:
		return "PetStripWhenSummonedLate"
	}
	return "PetPolicy(unknown)"
}

type TalentMod struct {
	Name string
	// Effect is the index of the modifying effect inside the talent spell.
	Effect  int32
	Applies TalentApplies
}

type ActionRef struct {
	SpellID int32
	ItemID  int32
}

type BuffSpec struct {
	Field          string // proto field name, snake_case, owns the number
	Number         int32  // existing proto field number, preserved
	Scope          BuffScope
	Proto          BuffProtoType // bool when the improving talent is not in the live Trait tree
	Kind           BuffKind
	Go             string      // identifier stem: "BattleShout"
	Name           string      // SpellName.Name_lang of the castable family ("" for rows with no spell)
	AuraName       string      // aura family when the cast is a summon or dummy (totems: "Strength of Earth")
	Anchor         int32       // explicit spell id; 0 = resolve Name via SkillLineAbility
	Owner          proto.Class // class that casts it; ClassUnknown for none
	Talent         *TalentMod  // improving talent family, nil when none
	Category       string      // exclusive-effect category value, "" = none
	SharedCategory string      // second exclusive category the aura also joins, "" = none
	SingleAura     bool
	Driver         bool // apply block hands the field to drive<Go>; the aura is not simply always up
	Pet            PetPolicy
	// StatOverride names the sim stats the aura's value lands on, spelled the
	// way stats.Stat.StatName() spells them, for an aura the client states
	// without saying which stat it means: A_MOD_CRIT_PCT carries no school, so
	// Leader of the Pack and Moonkin Aura name both crit stats here. The names
	// are checked while the row is resolved, and a row that states one may have
	// only one aura effect, since every effect would otherwise land on the same
	// stats.
	StatOverride []string
	Stats        []proto.Stat // UI relevance tags
	ImpAction    *ActionRef   // improved-icon override when it is an item
	Label        string       // UI label override, "" = DB name
	Notes        string       // reason for manual/absent rows; emitted as a comment
}

// GoField is the field name protoc-gen-go generates for Field. It ports
// protobuf-go's strs.GoCamelCase, including the branch that keeps the underscore
// in front of a digit ("soe_enhancement_2pt4" gives "SoeEnhancement_2Pt4").
func (s BuffSpec) GoField() string {
	var b []byte
	for i := 0; i < len(s.Field); i++ {
		c := s.Field[i]
		switch {
		case c == '_' && i+1 < len(s.Field) && isASCIILower(s.Field[i+1]):
		case isASCIIDigit(c):
			b = append(b, c)
		default:
			if isASCIILower(c) {
				c -= 'a' - 'A'
			}
			b = append(b, c)
			for ; i+1 < len(s.Field) && isASCIILower(s.Field[i+1]); i++ {
				b = append(b, s.Field[i+1])
			}
		}
	}
	return string(b)
}

// TSField is the property name protobuf-ts generates for Field: lower camel case
// in which a digit also capitalises the character that follows it
// ("soe_enhancement_2pt4" gives "soeEnhancement2Pt4").
func (s BuffSpec) TSField() string {
	var b strings.Builder
	capNext := false
	for i := 0; i < len(s.Field); i++ {
		c := s.Field[i]
		switch {
		case c == '_':
			capNext = true
		case isASCIIDigit(c):
			b.WriteByte(c)
			capNext = true
		case capNext:
			if isASCIILower(c) {
				c -= 'a' - 'A'
			}
			b.WriteByte(c)
			capNext = false
		case i == 0:
			if !isASCIILower(c) {
				c += 'a' - 'A'
			}
			b.WriteByte(c)
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

type RetiredField struct {
	Number int32
	Field  string
}

// Retired holds the fields api version 17 gave up, per scope in ascending number
// order, so that no later row takes a number or a name back: a saved payload or a
// share link from before the bump still carries the old field under both, and a new
// field reusing either would read the old value. gen_buffs_proto turns each entry
// into a `reserved` number and a `reserved` name, and
// ui/sim/proto/buff_field_migration.ts drops the same names from an older payload.
var Retired = map[BuffScope][]RetiredField{
	ScopeRaid: {
		{7, "bloodlust"},
	},
	ScopeParty: {
		{1, "ferocious_inspiration"},
		{5, "sanctity_aura"},
		{9, "draenei_racial_melee"},
		{10, "draenei_racial_caster"},
		{11, "drums"},
		{14, "braided_eternium_chain"},
		{15, "eye_of_the_night"},
		{16, "chain_of_the_twilight_owl"},
		{17, "jade_pendant_of_blasting"},
		{20, "totem_of_wrath"},
		{21, "wrath_of_air_totem"},
		{22, "snapshot_improved_wrath_of_air_totem"},
		{25, "snapshot_improved_strength_of_earth_totem"},
		{26, "tranquil_air_totem"},
		{29, "bs_solarian_sapphire"},
		{30, "snapshot_bs_solarian_sapphire"},
		{32, "snapshot_bs_booming_voice_rank"},
		{37, "soe_enhancement_2pt4"},
	},
	ScopeIndividual: {
		{3, "blessing_of_sanctuary"},
		{7, "unleashed_rage"},
	},
	ScopeDebuff: {
		{4, "misery"},
		{6, "isb_uptime"},
		{7, "shadow_weaving"},
		{8, "improved_scorch"},
		{9, "winters_chill"},
		{10, "blood_frenzy"},
		{18, "expose_weakness_uptime"},
		{19, "expose_weakness_hunter_agility"},
		{25, "shadow_embrace"},
		{26, "screech"},
		{27, "hemorrhage_uptime"},
		{28, "joc_retribution_2pt4"},
	},
}

func RetiredNumbers(scope BuffScope) []int32 {
	out := make([]int32, 0, len(Retired[scope]))
	for _, retired := range Retired[scope] {
		out = append(out, retired.Number)
	}
	return out
}

func RetiredFields(scope BuffScope) []string {
	out := make([]string, 0, len(Retired[scope]))
	for _, retired := range Retired[scope] {
		out = append(out, retired.Field)
	}
	return out
}

func ByScope(scope BuffScope) []BuffSpec {
	var out []BuffSpec
	for _, spec := range Manifest {
		if spec.Scope == scope {
			out = append(out, spec)
		}
	}
	return out
}

func isASCIILower(c byte) bool { return c >= 'a' && c <= 'z' }

func isASCIIDigit(c byte) bool { return c >= '0' && c <= '9' }
