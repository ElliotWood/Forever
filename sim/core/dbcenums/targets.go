package dbcenums

// ImplicitTarget_0/_1, under TrinityCore's Targets names, for the values the store's rows and
// resolver use.
type ImplicitTarget uint8

const (
	TARGET_UNIT_CASTER                    ImplicitTarget = 1
	TARGET_UNIT_PET                       ImplicitTarget = 5
	TARGET_UNIT_TARGET_ENEMY              ImplicitTarget = 6
	TARGET_UNIT_SRC_AREA_ENTRY            ImplicitTarget = 7
	TARGET_DEST_HOME                      ImplicitTarget = 9
	TARGET_UNIT_SRC_AREA_ENEMY            ImplicitTarget = 15
	TARGET_UNIT_DEST_AREA_ENEMY           ImplicitTarget = 16
	TARGET_DEST_DB                        ImplicitTarget = 17
	TARGET_DEST_CASTER                    ImplicitTarget = 18
	TARGET_UNIT_CASTER_AREA_PARTY         ImplicitTarget = 20
	TARGET_UNIT_TARGET_ALLY               ImplicitTarget = 21
	TARGET_SRC_CASTER                     ImplicitTarget = 22
	TARGET_GAMEOBJECT_TARGET              ImplicitTarget = 23
	TARGET_UNIT_CONE_ENEMY_24             ImplicitTarget = 24
	TARGET_UNIT_TARGET_ANY                ImplicitTarget = 25
	TARGET_UNIT_MASTER                    ImplicitTarget = 27
	TARGET_DEST_DYNOBJ_ENEMY              ImplicitTarget = 28
	TARGET_UNIT_SRC_AREA_ALLY             ImplicitTarget = 30
	TARGET_UNIT_DEST_AREA_ALLY            ImplicitTarget = 31
	TARGET_DEST_CASTER_SUMMON             ImplicitTarget = 32
	TARGET_UNIT_SRC_AREA_PARTY            ImplicitTarget = 33
	TARGET_UNIT_DEST_AREA_PARTY           ImplicitTarget = 34
	TARGET_UNIT_TARGET_PARTY              ImplicitTarget = 35
	TARGET_UNIT_NEARBY_ENTRY              ImplicitTarget = 38
	TARGET_DEST_CASTER_FRONT_RIGHT        ImplicitTarget = 41
	TARGET_DEST_CASTER_BACK_RIGHT         ImplicitTarget = 42
	TARGET_DEST_CASTER_BACK_LEFT          ImplicitTarget = 43
	TARGET_DEST_CASTER_FRONT_LEFT         ImplicitTarget = 44
	TARGET_UNIT_TARGET_CHAINHEAL_ALLY     ImplicitTarget = 45
	TARGET_DEST_CASTER_FRONT              ImplicitTarget = 47
	TARGET_DEST_CASTER_RIGHT              ImplicitTarget = 49
	TARGET_GAMEOBJECT_DEST_AREA           ImplicitTarget = 52
	TARGET_DEST_TARGET_ENEMY              ImplicitTarget = 53
	TARGET_DEST_CASTER_FRONT_LEAP         ImplicitTarget = 55
	TARGET_UNIT_CASTER_AREA_RAID          ImplicitTarget = 56
	TARGET_UNIT_TARGET_RAID               ImplicitTarget = 57
	TARGET_DEST_TARGET_BACK               ImplicitTarget = 65
	TARGET_DEST_CASTER_RANDOM             ImplicitTarget = 72
	TARGET_DEST_TARGET_RANDOM             ImplicitTarget = 74
	TARGET_DEST_CHANNEL_TARGET            ImplicitTarget = 76
	TARGET_UNIT_CHANNEL_TARGET            ImplicitTarget = 77
	TARGET_DEST_DEST_RIGHT                ImplicitTarget = 80
	TARGET_DEST_DEST_FRONT_RIGHT          ImplicitTarget = 82
	TARGET_DEST_DEST                      ImplicitTarget = 87
	TARGET_UNIT_CONE_CASTER_TO_DEST_ENEMY ImplicitTarget = 104
	TARGET_UNIT_CASTER_AND_SUMMONS        ImplicitTarget = 120
	TARGET_CORPSE_TARGET_ALLY             ImplicitTarget = 121
	TARGET_DEST_TARGET_ALLY               ImplicitTarget = 132
)
