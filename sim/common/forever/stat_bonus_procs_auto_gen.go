package forever

import (
	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

func RegisterAllProcs() {

	// Procs

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 10 health every 5.0 sec.
	// https://www.wowhead.com/forever/spell=5707
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 833, ItemName: "Lifestone"},
	//	{ItemID: 833, ItemName: "Lifestone"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance to strike your ranged target with a Flaming Cannonball for 49 Fire damage.
	// https://www.wowhead.com/forever/spell=29639
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 2099, ItemName: "Dwarven Hand Cannon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance to strike your target with a Frost Arrow for 45 Frost damage.
	// https://www.wowhead.com/forever/spell=29502
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 2824, ItemName: "Hurricane"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 5 Fire damage to anyone who strikes you with a melee attack.
	// https://www.wowhead.com/forever/spell=21142
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 3475, ItemName: "Cloak of Flames"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 1992 mana over 30s. Must remain seated while drinking.
	// https://www.wowhead.com/forever/spell=1135
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 4696, ItemName: "Lapidis Tankard of Tidesippe"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Blasts enemies in front of you with the power of wind, fire, all that kind of thing!
	// https://www.wowhead.com/forever/spell=14537
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 7734, ItemName: "Six Demon Bag"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Absorbs 600 magical damage. Lasts 2min.
	// https://www.wowhead.com/forever/spell=10618
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 8367, ItemName: "Dragonscale Breastplate"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 5 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21596
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 10659, ItemName: "Shard of the Splithooves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 33 when fighting Demons.
	// https://www.wowhead.com/forever/spell=18079
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 10696, ItemName: "Enchanted Azsharite Felbane Sword"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 33 when fighting Demons.
	// https://www.wowhead.com/forever/spell=18079
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 10697, ItemName: "Enchanted Azsharite Felbane Dagger"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 78 when fighting Demons.
	// https://www.wowhead.com/forever/spell=18087
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 10698, ItemName: "Enchanted Azsharite Felbane Staff"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 30 when fighting Undead.
	// https://www.wowhead.com/forever/spell=18074
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 10805, ItemName: "Eater of the Dead"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Has a 2% chance when struck in combat of protecting you with a holy shield.
	// https://www.wowhead.com/forever/spell=10368
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 11302, ItemName: "Uther's Strength"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 24 when fighting Beasts.
	// https://www.wowhead.com/forever/spell=18201
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 11628, ItemName: "Houndmaster's Bow"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 24 when fighting Beasts.
	// https://www.wowhead.com/forever/spell=18201
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 11629, ItemName: "Houndmaster's Rifle"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat inflicts 3 Arcane damage to the attacker.
	// https://www.wowhead.com/forever/spell=15438
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 11669, ItemName: "Naglering"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Channels 75 health into mana every 1.0 sec for 10s.
	// https://www.wowhead.com/forever/spell=17447
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 11808, ItemName: "Circle of Flame"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat has a 1% chance of reducing all melee damage taken by 25 for 10 sec.
	// https://www.wowhead.com/forever/spell=15595
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 11810, ItemName: "Force of Will"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 30 mana every 1.0 sec for 10s.
	// https://www.wowhead.com/forever/spell=15604
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 11819, ItemName: "Second Wind"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces mana cost of all spells by 100 for 10s.
	// https://www.wowhead.com/forever/spell=15646
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 11832, ItemName: "Burst of Knowledge"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Flings a magical boomerang towards target enemy dealing 187 damage and has a chance to stun or disarm
	// them.
	// https://www.wowhead.com/forever/spell=15712
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 11905, ItemName: "Linken's Boomerang"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 30 when fighting Beasts.
	// https://www.wowhead.com/forever/spell=14565
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 11906, ItemName: "Beastsmasher"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 72 when fighting Beasts.
	// https://www.wowhead.com/forever/spell=18076
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 11907, ItemName: "Beastslayer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=16372
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12344, ItemName: "Seal of Ascension"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 90 damage when you are the victim of a critical melee strike.
	// https://www.wowhead.com/forever/spell=16550
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12588, ItemName: "Bonespike Shoulder"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck has a 3% chance of stealing 120 life from the attacker over 4s.
	// https://www.wowhead.com/forever/spell=16608
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12628, ItemName: "Demon Forged Breastplate"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Disarm duration reduced by 50%.
	// https://www.wowhead.com/forever/spell=43588
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12639, ItemName: "Stronghold Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat has a 5% chance to make you invulnerable to melee damage for 3s. This effect can
	// only occur once every 30 sec.
	// https://www.wowhead.com/forever/spell=16621
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12641, ItemName: "Invulnerable Mail"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 45 when fighting Beasts.
	// https://www.wowhead.com/forever/spell=18067
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 12709, ItemName: "Pip's Skinner"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance to strike your ranged target with a Shadow Bolt for 19 Shadow damage.
	// https://www.wowhead.com/forever/spell=29640
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13040, ItemName: "Heartseeking Crossbow"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 99 when fighting Demons.
	// https://www.wowhead.com/forever/spell=18212
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13044, ItemName: "Demonslayer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// A protective mana shield surrounds the caster absorbing 500 damage. While the shield holds, increases
	// mana regeneration by 22 every 5.0 sec for 30min.
	// https://www.wowhead.com/forever/spell=17252
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13143, ItemName: "Mark of the Dragon Lord"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance to strike your ranged target with a Flaming Shell for 26 Fire damage.
	// https://www.wowhead.com/forever/spell=29647
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13146, ItemName: "Shell Launcher Shotgun"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 125 Fire damage to all targets in a cone in front of the caster.
	// https://www.wowhead.com/forever/spell=17283
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13171, ItemName: "Smokey's Lighter"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 81 when fighting Undead. It also allows the acquisition of Scourgestones on
	// behalf of the Argent Dawn.
	// https://www.wowhead.com/forever/spell=23930
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13209, ItemName: "Seal of the Dawn"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 48 when fighting Beasts.
	// https://www.wowhead.com/forever/spell=17482
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13212, ItemName: "Halycon's Spiked Collar"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Poisons target for 20 Nature damage every 2.0 sec for 20s.
	// https://www.wowhead.com/forever/spell=17330
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13213, ItemName: "Smolderweb's Eye"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 35 damage every time you block.
	// https://www.wowhead.com/forever/spell=17496
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13375, ItemName: "Crest of Retribution"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Causes nearby players to dance.
	// https://www.wowhead.com/forever/spell=18400
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13379, ItemName: "Piccolo of the Flaming Fire"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 5 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21596
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13387, ItemName: "Foresight Girdle"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the effect that healing and mana potions have on the wearer by 40%. This effect does not stack.
	// https://www.wowhead.com/forever/spell=17619
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13503, ItemName: "Alchemist's Stone"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases movement speed and life regeneration rate.
	// https://www.wowhead.com/forever/spell=17625
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13505, ItemName: "Runeblade of Baron Rivendare"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Harness the power of lightning to strike down all enemies around you for 200 Nature damage.
	// https://www.wowhead.com/forever/spell=17668
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13515, ItemName: "Ramstein's Lightning Bolts"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives 20 additional intellect to party members within 30 yards.
	// https://www.wowhead.com/forever/spell=18264
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 13937, ItemName: "Headmaster's Charge"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 25 Fire damage every 5.0 sec to all nearby enemies for 15s.
	// https://www.wowhead.com/forever/spell=18364
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 14134, ItemName: "Cloak of Fire"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 625 mana.
	// https://www.wowhead.com/forever/spell=18385
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 14152, ItemName: "Robe of the Archmage"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Heal your pet for 750.
	// https://www.wowhead.com/forever/spell=18386
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 14153, ItemName: "Robe of the Void"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the cooldown of your Fade ability by -2.0 sec.
	// https://www.wowhead.com/forever/spell=18388
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 14154, ItemName: "Truefaith Vestments"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Protects the wearer from being fully engulfed by Shadow Flame.
	// https://www.wowhead.com/forever/spell=22683
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 15138, ItemName: "Onyxia Scale Cloak"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 33 when fighting Beasts.
	// https://www.wowhead.com/forever/spell=19380
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 15782, ItemName: "Beaststalker Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 33 when fighting Beasts.
	// https://www.wowhead.com/forever/spell=19380
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 15783, ItemName: "Beasthunter Dagger"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance to strike your ranged target with Shadow Shot for 26 Shadow damage.
	// https://www.wowhead.com/forever/spell=29641
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16004, ItemName: "Dark Iron Rifle"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the duration of any Silence or Interrupt effects used against the wearer by 10%. This effect does
	// not stack with other similar effects.
	// https://www.wowhead.com/forever/spell=42184
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16009, ItemName: "Voice Amplification Modulator"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16391, ItemName: "Knight-Lieutenant's Silk Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16392, ItemName: "Knight-Lieutenant's Leather Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of your Arcane Shot by 15.
	// https://www.wowhead.com/forever/spell=23157
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16403, ItemName: "Knight-Lieutenant's Chain Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=22778
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16406, ItemName: "Knight-Lieutenant's Plate Gauntlets"},
	//	{ItemID: 23286, ItemName: "Knight-Lieutenant's Plate Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the Holy damage bonus of your Judgement of the Crusader by 20.
	// https://www.wowhead.com/forever/spell=23300
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16410, ItemName: "Knight-Lieutenant's Lamellar Gauntlets"},
	//	{ItemID: 23274, ItemName: "Knight-Lieutenant's Lamellar Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16440, ItemName: "Marshal's Silk Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16446, ItemName: "Marshal's Leather Footguards"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the Holy damage bonus of your Judgement of the Crusader by 20.
	// https://www.wowhead.com/forever/spell=23300
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16471, ItemName: "Marshal's Lamellar Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=22778
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16484, ItemName: "Marshal's Plate Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16487, ItemName: "Blood Guard's Silk Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16498, ItemName: "Blood Guard's Leather Treads"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=22778
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16510, ItemName: "Blood Guard's Plate Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%. Does not function for players higher than level
	// 60.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16518, ItemName: "Blood Guard's Mail Walkers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of your Arcane Shot by 15.
	// https://www.wowhead.com/forever/spell=23157
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16530, ItemName: "Blood Guard's Chain Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16540, ItemName: "General's Silk Handguards"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=22778
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16548, ItemName: "General's Plate Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16558, ItemName: "General's Leather Treads"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%. Does not function for players higher than level
	// 60.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16573, ItemName: "General's Mail Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 100 health every 1.0 sec for 10s.
	// https://www.wowhead.com/forever/spell=20631
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16768, ItemName: "Furbolg Medicine Pouch"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Disarm duration reduced by 50%.
	// https://www.wowhead.com/forever/spell=43588
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16907, ItemName: "Bloodfang Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 4 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21347
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16928, ItemName: "Nemesis Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 4 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21347
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16929, ItemName: "Nemesis Skullcap"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 4 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21347
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 16932, ItemName: "Nemesis Spaulders"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat inflicts 3 Arcane damage to the attacker.
	// https://www.wowhead.com/forever/spell=15438
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17066, ItemName: "Drillborer Disk"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 16 health per 5 sec.
	// https://www.wowhead.com/forever/spell=23210
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17082, ItemName: "Shard of the Flame"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 150 mana or 20 rage when you kill a target that gives experience; this effect cannot occur more
	// than once every 10 seconds.
	// https://www.wowhead.com/forever/spell=21186
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17104, ItemName: "Spinal Reaper"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Deals 5 Fire damage to anyone who strikes you with a melee attack.
	// https://www.wowhead.com/forever/spell=21142
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17182, ItemName: "Sulfuras, Hand of Ragnaros"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17564, ItemName: "Knight-Lieutenant's Dreadweave Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17577, ItemName: "Blood Guard's Dreadweave Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17584, ItemName: "Marshal's Dreadweave Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17588, ItemName: "General's Dreadweave Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17596, ItemName: "Knight-Lieutenant's Satin Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17608, ItemName: "Marshal's Satin Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17617, ItemName: "Blood Guard's Satin Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17620, ItemName: "General's Satin Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 5 health every 5.0 sec.
	// https://www.wowhead.com/forever/spell=20969
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17743, ItemName: "Resurgence Rod"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Removes 1 poison effect.
	// https://www.wowhead.com/forever/spell=21954
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17744, ItemName: "Heart of Noxxion"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance to strike your ranged target with Keeper's Sting for 21 Nature damage.
	// https://www.wowhead.com/forever/spell=29655
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17753, ItemName: "Verdant Keeper's Aim"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Absorbs 650 physical damage. Lasts 10s.
	// https://www.wowhead.com/forever/spell=21956
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17759, ItemName: "Mark of Resolution"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 2 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21346
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17900, ItemName: "Stormpike Insignia Rank 2"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 2 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21346
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17901, ItemName: "Stormpike Insignia Rank 3"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 5 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21596
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17902, ItemName: "Stormpike Insignia Rank 4"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 7 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21600
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17903, ItemName: "Stormpike Insignia Rank 5"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 8 health per 5 sec.
	// https://www.wowhead.com/forever/spell=20885
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17904, ItemName: "Stormpike Insignia Rank 6"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 2 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21346
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17905, ItemName: "Frostwolf Insignia Rank 2"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 2 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21346
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17906, ItemName: "Frostwolf Insignia Rank 3"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 5 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21596
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17907, ItemName: "Frostwolf Insignia Rank 4"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 7 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21600
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17908, ItemName: "Frostwolf Insignia Rank 5"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 8 health per 5 sec.
	// https://www.wowhead.com/forever/spell=20885
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 17909, ItemName: "Frostwolf Insignia Rank 6"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 7 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21601
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18298, ItemName: "Unbridled Leggings"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 36 when fighting Elementals.
	// https://www.wowhead.com/forever/spell=22836
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18310, ItemName: "Fiendish Machete"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 6 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21598
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18315, ItemName: "Ring of Demonic Potency"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 7 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21601
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18319, ItemName: "Fervent Helm"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat inflicts 3 Arcane damage to the attacker.
	// https://www.wowhead.com/forever/spell=20886
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18326, ItemName: "Razor Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the duration of any Silence or Interrupt effects used against the wearer by 10%. This effect does
	// not stack with other similar effects.
	// https://www.wowhead.com/forever/spell=42184
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18345, ItemName: "Murmuring Ring"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases damage done to Undead by magical spells and effects by up to 35.
	// https://www.wowhead.com/forever/spell=22849
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18346, ItemName: "Threadbare Trousers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage of your Imp's Firebolt spell by 8.
	// https://www.wowhead.com/forever/spell=22855
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18354, ItemName: "Pimgib's Collar"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage done by your pets by 4%.
	// https://www.wowhead.com/forever/spell=22854
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18355, ItemName: "Ferra's Collar"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 4 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21347
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18382, ItemName: "Fluctuating Cloak"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 7 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21601
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18406, ItemName: "Onyxia Blood Talisman"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 4 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21347
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18535, ItemName: "Milli's Shield"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 6 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21348
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18602, ItemName: "Tome of Sacrifice"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reflects Frost spells back at their caster for 5s. Chance to be resisted when used by players over level
	// 60.
	// https://www.wowhead.com/forever/spell=23131
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18634, ItemName: "Gyrofreeze Ice Reflector"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 625 health and mana to a friendly target and attempts to dispel any polymorph effects from them.
	// Reduced effectiveness against polymorph effects from casters of level 61 and higher.
	// https://www.wowhead.com/forever/spell=23064
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18637, ItemName: "Major Recombobulator"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reflects Fire spells back at their caster for 5s. Chance to be resisted when used by players over level
	// 60.
	// https://www.wowhead.com/forever/spell=23097
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18638, ItemName: "Hyper-Radiant Flame Reflector"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reflects Shadow spells back at their caster for 5s. Chance to be resisted when used by players over level
	// 60.
	// https://www.wowhead.com/forever/spell=23132
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18639, ItemName: "Ultra-Flash Shadow Reflector"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 45 when fighting Demons.
	// https://www.wowhead.com/forever/spell=14097
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18715, ItemName: "Lok'delar, Stave of the Ancient Keepers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Disarm duration reduced by 50%.
	// https://www.wowhead.com/forever/spell=43588
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18722, ItemName: "Death Grips"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 45 when fighting Undead.
	// https://www.wowhead.com/forever/spell=18098
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18758, ItemName: "Specter's Blade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 8 health per 5 sec.
	// https://www.wowhead.com/forever/spell=20885
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18760, ItemName: "Necromantic Band"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 6 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21599
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18809, ItemName: "Sash of Whispered Secrets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck in combat inflicts 13 Fire damage to the attacker.
	// https://www.wowhead.com/forever/spell=23266
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18815, ItemName: "Essence of the Pure Flame"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Safely transport yourself to Gadgetzan in Tanaris! Emphasis on Safe! Yup, nothing bad could ever happen
	// while using this device!
	// https://www.wowhead.com/forever/spell=23453
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 18986, ItemName: "Ultrasafe Transporter: Gadgetzan"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Absorbs 1250 damage. Lasts 20s.
	// https://www.wowhead.com/forever/spell=23506
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19024, ItemName: "Arena Grand Master"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Sometimes heals bearer of 180 damage when damaging an enemy in melee.
	// https://www.wowhead.com/forever/spell=23682
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19287, ItemName: "Darkmoon Card: Heroism"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// 2% chance on successful spellcast to allow 100% of your Mana regeneration to continue while casting for
	// 15s.
	// https://www.wowhead.com/forever/spell=23684
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19288, ItemName: "Darkmoon Card: Blue Dragon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance to strike your melee target with lightning for 300 Nature damage.
	// https://www.wowhead.com/forever/spell=23687
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19289, ItemName: "Darkmoon Card: Maelstrom"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives the wearer a 10% chance of being able to resurrect with 20% health and mana.
	// https://www.wowhead.com/forever/spell=23701
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19290, ItemName: "Darkmoon Card: Twisting Nether"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 6 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21348
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19302, ItemName: "Darkmoon Ring"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Infuses you with Arcane energy, causing your next Arcane Shot fired within 10s to detonate at the target.
	// The Arcane Detonation will deal 15215 damage to enemies near the target.
	// https://www.wowhead.com/forever/spell=23721
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19336, ItemName: "Arcane Infused Gem"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Decreases the mana cost of all Druid shapeshifting forms by 550 for 20s.
	// https://www.wowhead.com/forever/spell=23724
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19340, ItemName: "Rune of Metamorphosis"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your maximum health by 1500 for 20 sec. After the effects wear off, you will lose the extra
	// maximum health.
	// https://www.wowhead.com/forever/spell=23725
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19341, ItemName: "Lifegiving Gem"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage dealt by Instant Poison by 65 and the periodic damage dealt by Deadly Poison by 22
	// for 20s.
	// https://www.wowhead.com/forever/spell=23726
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19342, ItemName: "Venomous Totem"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the cost of your Hamstring ability by -2.0 rage points.
	// https://www.wowhead.com/forever/spell=24428
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19577, ItemName: "Rage of Mugamba"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of Hammer of Justice by 0.5 sec.
	// https://www.wowhead.com/forever/spell=24188
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19588, ItemName: "Hero's Brand"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the amount of damage absorbed by Power Word: Shield by 35.
	// https://www.wowhead.com/forever/spell=24191
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19594, ItemName: "The All-Seeing Eye of Zuldazar"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the cooldown of Counterspell by -2.0 sec.
	// https://www.wowhead.com/forever/spell=24429
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19601, ItemName: "Jewel of Kajaro"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the radius of Rain of Fire and Hellfire by 1 yard.
	// https://www.wowhead.com/forever/spell=24430
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19605, ItemName: "Kezan's Unstoppable Taint"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Decreases the mana cost of your Healing Stream and Mana Spring totems by 20.
	// https://www.wowhead.com/forever/spell=24436
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19609, ItemName: "Unmarred Vision of Voodress"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Decreases the cooldown of Kick by -0.5 sec.
	// https://www.wowhead.com/forever/spell=24434
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19617, ItemName: "Zandalarian Shadow Mastery Talisman"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Decreases the cooldown of Feign Death by -2.0 sec.
	// https://www.wowhead.com/forever/spell=24432
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19621, ItemName: "Maelstrom's Wrath"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases damage done to Undead by magical spells and effects by up to 48. It also allows the acquisition
	// of Scourgestones on behalf of the Argent Dawn.
	// https://www.wowhead.com/forever/spell=24198
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19812, ItemName: "Rune of the Dawn"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance to decapitate the target on a melee swing, causing 676 damage.
	// https://www.wowhead.com/forever/spell=24241
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19874, ItemName: "Halberd of Smiting"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 60 mana every 5.0 sec for 30s.
	// https://www.wowhead.com/forever/spell=24268
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19930, ItemName: "Mar'li's Eye"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 60 when fighting Beasts.
	// https://www.wowhead.com/forever/spell=18207
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19946, ItemName: "Tigule's Harpoon"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your melee and ranged damage by 40 for 20s. Every time you hit a target, this bonus is reduced
	// by 2.
	// https://www.wowhead.com/forever/spell=24661
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19949, ItemName: "Zandalarian Hero Medallion"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your spell damage by up to 204 and your healing by up to 408 for 20s. Every time you cast a
	// spell, the bonus is reduced by 17 spell damage and 34 healing.
	// https://www.wowhead.com/forever/spell=24658
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19950, ItemName: "Zandalarian Hero Charm"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Instantly increases your rage by 30.
	// https://www.wowhead.com/forever/spell=24571
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19951, ItemName: "Gri'lek's Charm of Might"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Instantly clears the cooldowns of Aimed Shot, Multishot, and Volley.
	// https://www.wowhead.com/forever/spell=24531
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19953, ItemName: "Renataki's Charm of Beasts"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Instantly increases your energy by 60.
	// https://www.wowhead.com/forever/spell=24532
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19954, ItemName: "Renataki's Charm of Trickery"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage dealt by your Lightning Shield spell by 305 for 20s.
	// https://www.wowhead.com/forever/spell=24499
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19956, ItemName: "Wushoolay's Charm of Spirits"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 48 when fighting Dragonkin.
	// https://www.wowhead.com/forever/spell=24291
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19961, ItemName: "Gri'lek's Grinder"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 117 when fighting Dragonkin.
	// https://www.wowhead.com/forever/spell=24292
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19962, ItemName: "Gri'lek's Carver"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 117 when fighting Demons.
	// https://www.wowhead.com/forever/spell=14098
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19963, ItemName: "Pitchfork of Madness"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces damage from falling.
	// https://www.wowhead.com/forever/spell=24350
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19982, ItemName: "Duskbat Drape"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your pet's next attack is guaranteed to critically strike if that attack is capable of striking critically.
	// https://www.wowhead.com/forever/spell=24353
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 19992, ItemName: "Devilsaur Tooth"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your Frostbolt spells have a 6% chance to restore 50 mana when cast.
	// https://www.wowhead.com/forever/spell=24392
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20035, ItemName: "Glacial Spike"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20048, ItemName: "Highlander's Plate Greaves"},
	//	{ItemID: 20127, ItemName: "Highlander's Plate Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20049, ItemName: "Highlander's Lamellar Greaves"},
	//	{ItemID: 20109, ItemName: "Highlander's Lamellar Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20050, ItemName: "Highlander's Chain Greaves"},
	//	{ItemID: 20091, ItemName: "Highlander's Chain Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20051, ItemName: "Highlander's Mail Greaves"},
	//	{ItemID: 20121, ItemName: "Highlander's Mail Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20052, ItemName: "Highlander's Leather Boots"},
	//	{ItemID: 20112, ItemName: "Highlander's Leather Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20053, ItemName: "Highlander's Lizardhide Boots"},
	//	{ItemID: 20100, ItemName: "Highlander's Lizardhide Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20054, ItemName: "Highlander's Cloth Boots"},
	//	{ItemID: 20094, ItemName: "Highlander's Cloth Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Absorbs 605 physical damage. Lasts 15s.
	// https://www.wowhead.com/forever/spell=23991
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20071, ItemName: "Talisman of Arathor"},
	//	{ItemID: 21117, ItemName: "Talisman of Arathor"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Absorbs 605 physical damage. Lasts 15s.
	// https://www.wowhead.com/forever/spell=23991
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20072, ItemName: "Defiler's Talisman"},
	//	{ItemID: 21115, ItemName: "Defiler's Talisman"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20154, ItemName: "Defiler's Chain Greaves"},
	//	{ItemID: 20155, ItemName: "Defiler's Chain Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20159, ItemName: "Defiler's Cloth Boots"},
	//	{ItemID: 20160, ItemName: "Defiler's Cloth Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20167, ItemName: "Defiler's Lizardhide Boots"},
	//	{ItemID: 20170, ItemName: "Defiler's Lizardhide Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20181, ItemName: "Defiler's Lamellar Greaves"},
	//	{ItemID: 20185, ItemName: "Defiler's Lamellar Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20186, ItemName: "Defiler's Leather Boots"},
	//	{ItemID: 20189, ItemName: "Defiler's Leather Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20199, ItemName: "Defiler's Mail Greaves"},
	//	{ItemID: 20202, ItemName: "Defiler's Mail Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Run speed increased slightly.
	// https://www.wowhead.com/forever/spell=23990
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20208, ItemName: "Defiler's Plate Greaves"},
	//	{ItemID: 20211, ItemName: "Defiler's Plate Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 6 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21598
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20219, ItemName: "Tattered Hakkari Cape"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 45 when fighting Demons.
	// https://www.wowhead.com/forever/spell=14097
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20487, ItemName: "Lok'delar, Stave of the Ancient Keepers DEP"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 400 mana over 10s.
	// https://www.wowhead.com/forever/spell=24884
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20525, ItemName: "Earthen Sigil"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Casts your Summon Voidwalker spell with no mana or Soul Shard requirements.
	// https://www.wowhead.com/forever/spell=25112
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20534, ItemName: "Abyss Shard"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 4 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21595
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 20579, ItemName: "Green Dragonskin Cloak"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces your threat to enemy targets within 30 yards, making them less likely to attack you.
	// https://www.wowhead.com/forever/spell=25892
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21181, ItemName: "Grace of Earth"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Spikes sprout from you causing 25 Nature damage to attackers when hit. Lasts 30s.
	// https://www.wowhead.com/forever/spell=26168
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21488, ItemName: "Fetish of Chitinous Spikes"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the duration of any Silence or Interrupt effects used against the wearer by 20%. This effect does
	// not stack with other similar effects.
	// https://www.wowhead.com/forever/spell=35126
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21517, ItemName: "Gnomish Turban of Psychic Might"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Allows underwater breathing.
	// https://www.wowhead.com/forever/spell=11789
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21526, ItemName: "Band of Icy Depths"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your magical heals provide the target with a shield that absorbs damage equal to 15% of the amount healed
	// for 30s.
	// https://www.wowhead.com/forever/spell=26467
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21625, ItemName: "Scarab Brooch"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the threat you generate by 70% for 20s.
	// https://www.wowhead.com/forever/spell=26400
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21647, ItemName: "Fetish of the Sand Reaver"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your spell resistances by 100 for 1min. Every time a hostile spell lands on you, this bonus
	// is reduced by 10 resistance.
	// https://www.wowhead.com/forever/spell=26463
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21685, ItemName: "Petrified Scarab"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 5 health per 5 sec.
	// https://www.wowhead.com/forever/spell=21596
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21687, ItemName: "Ukko's Ring of Darkness"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Heals for 18 damage when you get a critical hit.
	// https://www.wowhead.com/forever/spell=26606
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21780, ItemName: "Blood Crown"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Calls down a meteor, burning all enemies within the area for 442 total Fire damage.
	// https://www.wowhead.com/forever/spell=26789
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 21891, ItemName: "Shard of the Fallen Star"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your pet's maximum health by 3%.
	// https://www.wowhead.com/forever/spell=27038
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22013, ItemName: "Beastmaster's Cap"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your pet's critical strike chance by 2%.
	// https://www.wowhead.com/forever/spell=27043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22015, ItemName: "Beastmaster's Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases your pet's armor by 10%.
	// https://www.wowhead.com/forever/spell=27225
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22060, ItemName: "Beastmaster's Tunic"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases damage dealt by your pet by 3%.
	// https://www.wowhead.com/forever/spell=27206
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22061, ItemName: "Beastmaster's Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Spell Damage received is reduced by 10.
	// https://www.wowhead.com/forever/spell=27518
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22191, ItemName: "Obsidian Mail Tunic"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// On successful melee or ranged attack gain 8 mana and if possible drain 8 mana from the target.
	// https://www.wowhead.com/forever/spell=0
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22194, ItemName: "Black Grasp of the Destroyer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck by a non-periodic damage spell you have a 30% chance of getting a 6s spell shield that absorbs
	// 500 of that school of damage.
	// https://www.wowhead.com/forever/spell=27539
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22196, ItemName: "Thick Obsidian Breastplate"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// When struck by a harmful spell, the caster of that spell has a 5% chance to be silenced for 3s.
	// https://www.wowhead.com/forever/spell=27559
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22198, ItemName: "Jagged Obsidian Shield"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance to bathe your melee target in flames for 180 Fire damage.
	// https://www.wowhead.com/forever/spell=27655
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22321, ItemName: "Heart of Wyrmthalak"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the cooldown of Reincarnation by -10.0 minutes.
	// https://www.wowhead.com/forever/spell=27797
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22345, ItemName: "Totem of Rebirth"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases damage done by Earth Shock, Flame Shock, and Frost Shock by up to 30.
	// https://www.wowhead.com/forever/spell=27859
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22395, ItemName: "Totem of Rage"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases healing done by Lesser Healing Wave by up to 80.
	// https://www.wowhead.com/forever/spell=27855
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22396, ItemName: "Totem of Life"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage of your Claw and Rake abilites by 20.
	// https://www.wowhead.com/forever/spell=27851
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22397, ItemName: "Idol of Ferocity"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases healing done by Rejuvenation by up to 50.
	// https://www.wowhead.com/forever/spell=27853
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22398, ItemName: "Idol of Rejuvenation"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the amount healed by Healing Touch by 100.
	// https://www.wowhead.com/forever/spell=27846
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22399, ItemName: "Idol of Health"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the armor from your Devotion Aura by 110.
	// https://www.wowhead.com/forever/spell=27850
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22400, ItemName: "Libram of Truth"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of your Cleanse spell by 25.
	// https://www.wowhead.com/forever/spell=27847
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22402, ItemName: "Libram of Grace"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=31796
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22589, ItemName: "Atiesh, Greatstaff of the Guardian"},
	//	{ItemID: 22630, ItemName: "Atiesh, Greatstaff of the Guardian"},
	//	{ItemID: 22631, ItemName: "Atiesh, Greatstaff of the Guardian"},
	//	{ItemID: 22632, ItemName: "Atiesh, Greatstaff of the Guardian"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Inflicts the will of the Ashbringer upon the wielder.
	// https://www.wowhead.com/forever/spell=28282
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22691, ItemName: "Corrupted Ashbringer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22856, ItemName: "Blood Guard's Leather Walkers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%. Does not function for players higher than level
	// 60.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22857, ItemName: "Blood Guard's Mail Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22865, ItemName: "Blood Guard's Dreadweave Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Hamstring Rage cost reduced by -3.0.
	// https://www.wowhead.com/forever/spell=22778
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22868, ItemName: "Blood Guard's Plate Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22869, ItemName: "Blood Guard's Satin Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 22870, ItemName: "Blood Guard's Silk Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the threat you generate by 35% for 20s.
	// https://www.wowhead.com/forever/spell=28862
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23001, ItemName: "Eye of Diminution"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gain up to 25 mana each time you cast Healing Touch.
	// https://www.wowhead.com/forever/spell=28847
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23004, ItemName: "Idol of Longevity"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Regain up to 10 mana each time you cast Lesser Healing Wave.
	// https://www.wowhead.com/forever/spell=28849
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23005, ItemName: "Totem of Flowing Water"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases healing done by Flash of Light by up to 83.
	// https://www.wowhead.com/forever/spell=28851
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23006, ItemName: "Libram of Light"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 500 mana.
	// https://www.wowhead.com/forever/spell=28760
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23027, ItemName: "Warmth of Forgiveness"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Disarm duration reduced by 50%.
	// https://www.wowhead.com/forever/spell=43588
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23072, ItemName: "Fists of the Unrelenting"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 60 when fighting Undead.
	// https://www.wowhead.com/forever/spell=28870
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23078, ItemName: "Gauntlets of Undead Slaying"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 60 when fighting Undead.
	// https://www.wowhead.com/forever/spell=28870
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23081, ItemName: "Handwraps of Undead Slaying"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 60 when fighting Undead.
	// https://www.wowhead.com/forever/spell=28870
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23082, ItemName: "Handguards of Undead Slaying"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases damage done to Undead by magical spells and effects by up to 35.
	// https://www.wowhead.com/forever/spell=22849
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23084, ItemName: "Gloves of Undead Cleansing"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases damage done to Undead by magical spells and effects by up to 48.
	// https://www.wowhead.com/forever/spell=24197
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23085, ItemName: "Robe of Undead Cleansing"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 81 when fighting Undead.
	// https://www.wowhead.com/forever/spell=17319
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23087, ItemName: "Breastplate of Undead Slaying"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 81 when fighting Undead.
	// https://www.wowhead.com/forever/spell=17319
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23088, ItemName: "Chestguard of Undead Slaying"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 81 when fighting Undead.
	// https://www.wowhead.com/forever/spell=17319
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23089, ItemName: "Tunic of Undead Slaying"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 45 when fighting Undead.
	// https://www.wowhead.com/forever/spell=18098
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23090, ItemName: "Bracers of Undead Slaying"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases damage done to Undead by magical spells and effects by up to 26.
	// https://www.wowhead.com/forever/spell=28876
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23091, ItemName: "Bracers of Undead Cleansing"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 45 when fighting Undead.
	// https://www.wowhead.com/forever/spell=18098
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23092, ItemName: "Wristguards of Undead Slaying"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 45 when fighting Undead.
	// https://www.wowhead.com/forever/spell=18098
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23093, ItemName: "Wristwraps of Undead Slaying"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases healing done by Lesser Healing Wave by up to 53.
	// https://www.wowhead.com/forever/spell=28856
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23200, ItemName: "Totem of Sustaining"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases healing done by Flash of Light by up to 53.
	// https://www.wowhead.com/forever/spell=28853
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23201, ItemName: "Libram of Divinity"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Chance to discharge electricity causing 150 Nature damage to your target.
	// https://www.wowhead.com/forever/spell=29151
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23221, ItemName: "Misplaced Servo Arm"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Searing Pain.
	// https://www.wowhead.com/forever/spell=23046
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23282, ItemName: "Knight-Lieutenant's Dreadweave Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the duration of your Sprint ability by 3.0 sec.
	// https://www.wowhead.com/forever/spell=23049
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23285, ItemName: "Knight-Lieutenant's Leather Walkers"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Gives you a 50% chance to avoid interruption caused by damage while casting Mind Blast.
	// https://www.wowhead.com/forever/spell=23043
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23288, ItemName: "Knight-Lieutenant's Satin Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage absorbed by your Mana Shield by 285.
	// https://www.wowhead.com/forever/spell=23037
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23290, ItemName: "Knight-Lieutenant's Silk Handwraps"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Absorbs 900 damage. Lasts 20s.
	// https://www.wowhead.com/forever/spell=29506
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23558, ItemName: "The Burrower's Shell"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reach into the hat for a drink.
	// https://www.wowhead.com/forever/spell=29830
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23587, ItemName: "Mirren's Drinking Hat"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Engage the rocket boots to greatly increase your speed... most of the time.
	// https://www.wowhead.com/forever/spell=51582
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23824, ItemName: "Rocket Boots Xtreme"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Protects you with a shield of force that stops 4000 damage for 8s.
	// https://www.wowhead.com/forever/spell=30458
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23825, ItemName: "Nigh Invulnerability Belt"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Turns the target into a chicken for 15s. Well, that is assuming the transmogrification polarity has not
	// been reversed...
	// https://www.wowhead.com/forever/spell=30507
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23835, ItemName: "Gnomish Poultryizer"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Fire a powerful rocket at the enemy that does 1440 damage and stuns them for 3 sec. This thing has quite
	// a kick though...
	// https://www.wowhead.com/forever/spell=46567
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 23836, ItemName: "Goblin Rocket Launcher"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases healing done by Rejuvenation by up to 86.
	// https://www.wowhead.com/forever/spell=32402
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 25643, ItemName: "Harold's Rejuvenating Broach"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases healing done by Flash of Light by up to 79.
	// https://www.wowhead.com/forever/spell=32403
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 25644, ItemName: "Blessed Book of Nagrand"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases healing done by Lesser Healing Wave by up to 79.
	// https://www.wowhead.com/forever/spell=32401
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 25645, ItemName: "Totem of the Plains"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces your threat to enemy targets within 30 yards, making them less likely to attack you.
	// https://www.wowhead.com/forever/spell=32599
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 25786, ItemName: "Hypnotist's Watch"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 200 mana when you kill a target that gives experience or honor. This effect cannot occur more
	// than once every 10 seconds.
	// https://www.wowhead.com/forever/spell=33743
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 28108, ItemName: "Power Infused Mushroom"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 200 health when you kill a target that gives experience or honor. This effect cannot occur more
	// than once every 10 seconds.
	// https://www.wowhead.com/forever/spell=33758
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 28109, ItemName: "Essence Infused Mushroom"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 39 when fighting Demons.
	// https://www.wowhead.com/forever/spell=35175
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 29398, ItemName: "Circle of Banishing"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%. Does not function for players higher than level
	// 60.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 29594, ItemName: "Knight-Lieutenant's Mail Greaves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the Holy damage bonus of your Judgement of the Crusader by 20.
	// https://www.wowhead.com/forever/spell=23300
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 29600, ItemName: "Blood Guard's Lamellar Gauntlets"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the speed of your Ghost Wolf ability by 15%. Does not function for players higher than level
	// 60.
	// https://www.wowhead.com/forever/spell=22801
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 29606, ItemName: "Marshal's Mail Boots"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the Holy damage bonus of your Judgement of the Crusader by 20.
	// https://www.wowhead.com/forever/spell=23300
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 29613, ItemName: "General's Lamellar Gloves"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Imbue your weapon with power, increasing attack power against undead and demons by 150. Lasts 5 min.
	// https://www.wowhead.com/forever/spell=37360
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 30696, ItemName: "Scourgebane"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases damage done to Demons by magical spells and effects by up to 185.
	// https://www.wowhead.com/forever/spell=37649
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 30787, ItemName: "Illidari-Bane Mageblade"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 93 when fighting Demons.
	// https://www.wowhead.com/forever/spell=37651
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 30788, ItemName: "Illidari-Bane Broadsword"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 150 when fighting Demons.
	// https://www.wowhead.com/forever/spell=37652
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 30789, ItemName: "Illidari-Bane Claymore"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Launch yourself from Outland to the stars. For the safety of others, please clear the launching platform
	// before use.
	// https://www.wowhead.com/forever/spell=37896
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 30847, ItemName: "X-52 Rocket Helmet"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of Stormstrike by 22.
	// https://www.wowhead.com/forever/spell=37762
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 31031, ItemName: "Stormfury Totem"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases attack power by 93 when fighting Demons.
	// https://www.wowhead.com/forever/spell=37651
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 31745, ItemName: "Illidari-Bane Dagger"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases speed in Flight Form and Swift Flight Form by 10%.
	// https://www.wowhead.com/forever/spell=48403
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 32481, ItemName: "Charm of Swift Flight"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Flash of Light and Holy Light have a 15% chance to grant your target 0 healing over 12s, and your Judgements
	// have a 50% chance to inflict 0 damage on their target over 8s.
	// https://www.wowhead.com/forever/spell=40470
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 32489, ItemName: "Ashtongue Talisman of Zeal"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Attach a lure to your equipped fishing pole, increasing Fishing by 75 for 10 min.
	// https://www.wowhead.com/forever/spell=43699
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 33820, ItemName: "Weather-Beaten Fishing Hat"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// An extremely potent alcoholic beverage.
	// https://www.wowhead.com/forever/spell=44540
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 34140, ItemName: "Dark Iron Tankard"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Engage the rocket boots to greatly increase your speed... most of the time.
	// https://www.wowhead.com/forever/spell=51582
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 35581, ItemName: "Rocket Boots Xtreme Lite"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the effect that healing and mana potions have on the wearer by 40%. This effect does not stack.
	// https://www.wowhead.com/forever/spell=17619
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 35748, ItemName: "Guardian's Alchemist Stone"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the effect that healing and mana potions have on the wearer by 40%. This effect does not stack.
	// https://www.wowhead.com/forever/spell=17619
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 35749, ItemName: "Sorcerer's Alchemist Stone"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the effect that healing and mana potions have on the wearer by 40%. This effect does not stack.
	// https://www.wowhead.com/forever/spell=17619
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 35750, ItemName: "Redeemer's Alchemist Stone"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the effect that healing and mana potions have on the wearer by 40%. This effect does not stack.
	// https://www.wowhead.com/forever/spell=17619
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 35751, ItemName: "Assassin's Alchemist Stone"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// https://www.wowhead.com/forever/spell=52172
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 38506, ItemName: "Don Carlos' Famous Hat"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Restores 600 health.
	// https://www.wowhead.com/forever/spell=352340
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 185988, ItemName: "Communal Stone of Stoicism"},
	//	{ItemID: 185988, ItemName: "Communal Stone of Stoicism"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage of your Moonfire spell by up to 10.
	// https://www.wowhead.com/forever/spell=352504
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 186052, ItemName: "Communal Idol of Wrath"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the damage of your Claw and Rake abilites by 5.
	// https://www.wowhead.com/forever/spell=352574
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 186053, ItemName: "Communal Idol of the Wild"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the periodic healing of your Rejuvenation by up to 15.
	// https://www.wowhead.com/forever/spell=352508
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 186054, ItemName: "Communal Idol of Life"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases healing done by Flash of Light by up to 10.
	// https://www.wowhead.com/forever/spell=352512
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 186065, ItemName: "Communal Book of Healing"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the armor from your Devotion Aura by 50.
	// https://www.wowhead.com/forever/spell=352513
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 186066, ItemName: "Communal Book of Protection"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the base mana cost of your Seal spells by 5.
	// https://www.wowhead.com/forever/spell=352580
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 186067, ItemName: "Communal Book of Righteousness"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases healing done by Lesser Healing Wave by up to 10.
	// https://www.wowhead.com/forever/spell=352517
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 186072, ItemName: "Communal Totem of Restoration"},
	// })

	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Reduces the mana cost of Stormstrike by 5.
	// https://www.wowhead.com/forever/spell=352522
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback:           core.CallbackEmpty,
	//	ProcMask:           core.ProcMaskUnknown,
	//	Outcome:            core.OutcomeEmpty,
	//	RequireDamageDealt: false
	// }, []shared.ItemVariant{
	//	{ItemID: 186073, ItemName: "Communal Totem of the Storm"},
	// })

	// When struck in combat has a 1% chance of inflicting 50 Frost damage to the attacker and freezing them
	// for 5s.
	// https://www.wowhead.com/forever/spell=18798
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:      942,
		SpellID:     18798,
		School:      core.SpellSchoolFrost,
		DefenseType: core.DefenseTypeMagic,
		MinDmg:      50,
		MaxDmg:      50,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagNoOnDamageDealt | core.SpellFlagProc,
		Trigger: core.ProcTrigger{
			Name:               "Freezing Band",
			ActionID:           core.ActionID{ItemID: 942},
			Callback:           core.CallbackOnSpellHitTaken,
			ProcMask:           core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial,
			Outcome:            core.OutcomeLanded,
			RequireDamageDealt: true,
			ProcChance:         0.01,
		},
	})

	// Adds 4 fire damage to your weapon attack.
	// https://www.wowhead.com/forever/spell=7714
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:      12631,
		SpellID:     7714,
		School:      core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		MinDmg:      4,
		MaxDmg:      4,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagNoOnDamageDealt | core.SpellFlagProc,
		Trigger: core.ProcTrigger{
			Name:               "Fiery Plate Gauntlets",
			ActionID:           core.ActionID{ItemID: 12631},
			Callback:           core.CallbackOnSpellHitDealt,
			ProcMask:           core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial,
			Outcome:            core.OutcomeLanded,
			RequireDamageDealt: true,
			ProcChance:         1,
		},
	})

	// Adds 3 Lightning damage to your melee attacks.
	// https://www.wowhead.com/forever/spell=16614
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:      12632,
		SpellID:     16614,
		School:      core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		MinDmg:      3,
		MaxDmg:      3,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagNoOnDamageDealt | core.SpellFlagProc,
		Trigger: core.ProcTrigger{
			Name:               "Storm Gauntlets",
			ActionID:           core.ActionID{ItemID: 12632},
			Callback:           core.CallbackOnSpellHitDealt,
			ProcMask:           core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial,
			Outcome:            core.OutcomeLanded,
			RequireDamageDealt: true,
			ProcChance:         1,
		},
	})

	// Reduces an enemy's armor by 200. Stacks up to 3 times.
	// https://www.wowhead.com/forever/spell=16928
	shared.NewStackingStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskUnknown,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
		IsWeaponProc:       true,
	}, []shared.ItemVariant{
		{ItemID: 12798, ItemName: "Annihilator"},
	})

	// Has a 1% chance when struck in combat of increasing block rating by 250 for 10 sec.
	// https://www.wowhead.com/forever/spell=17351
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback:           core.CallbackOnSpellHitTaken,
		ProcMask:           core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
	}, []shared.ItemVariant{
		{ItemID: 13243, ItemName: "Argent Defender"},
	})

	// When struck in combat has a 1% chance of increasing all party member's armor by 250 for 30s.
	// https://www.wowhead.com/forever/spell=18946
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback:           core.CallbackOnSpellHitTaken,
		ProcMask:           core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
	}, []shared.ItemVariant{
		{ItemID: 14557, ItemName: "The Lion Horn of Stormwind"},
	})

	// Adds 2 fire damage to your melee attacks.
	// https://www.wowhead.com/forever/spell=7712
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:      17111,
		SpellID:     7712,
		School:      core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		MinDmg:      2,
		MaxDmg:      2,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagNoOnDamageDealt | core.SpellFlagProc,
		Trigger: core.ProcTrigger{
			Name:               "Blazefury Medallion",
			ActionID:           core.ActionID{ItemID: 17111},
			Callback:           core.CallbackOnSpellHitDealt,
			ProcMask:           core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial,
			Outcome:            core.OutcomeLanded,
			RequireDamageDealt: true,
			ProcChance:         1,
		},
	})

	// Chance on landing a direct damage spell to deal 100 Shadow damage and restore 100 mana to you.
	// https://www.wowhead.com/forever/spell=27860
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:      17780,
		SpellID:     27860,
		School:      core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeNone,
		MinDmg:      100,
		MaxDmg:      100,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagNoOnDamageDealt | core.SpellFlagProc,
		Trigger: core.ProcTrigger{
			Name:               "Blade of Eternal Darkness",
			ActionID:           core.ActionID{ItemID: 17780},
			Callback:           core.CallbackOnCastComplete,
			ProcMask:           core.ProcMaskSpellDamage,
			Outcome:            core.OutcomeEmpty,
			RequireDamageDealt: false,
			ProcChance:         0.1,
		},
	})

	// When struck in combat has a 5% chance of inflicting 65 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=16782
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:      18825,
		SpellID:     16782,
		School:      core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		MinDmg:      35,
		MaxDmg:      65,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagNoOnDamageDealt | core.SpellFlagProc,
		Trigger: core.ProcTrigger{
			Name:               "Grand Marshal's Aegis",
			ActionID:           core.ActionID{ItemID: 18825},
			Callback:           core.CallbackOnSpellHitTaken,
			ProcMask:           core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
			Outcome:            core.OutcomeLanded,
			RequireDamageDealt: true,
			ProcChance:         0.05,
		},
	})

	// When struck in combat has a 5% chance of inflicting 65 Nature damage to the attacker.
	// https://www.wowhead.com/forever/spell=16782
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:      18826,
		SpellID:     16782,
		School:      core.SpellSchoolNature,
		DefenseType: core.DefenseTypeMagic,
		MinDmg:      35,
		MaxDmg:      65,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagNoOnDamageDealt | core.SpellFlagProc,
		Trigger: core.ProcTrigger{
			Name:               "High Warlord's Shield Wall",
			ActionID:           core.ActionID{ItemID: 18826},
			Callback:           core.CallbackOnSpellHitTaken,
			ProcMask:           core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
			Outcome:            core.OutcomeLanded,
			RequireDamageDealt: true,
			ProcChance:         0.05,
		},
	})

	// Adds 2 fire damage to your melee attacks.
	// https://www.wowhead.com/forever/spell=7712
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:      19968,
		SpellID:     7712,
		School:      core.SpellSchoolFire,
		DefenseType: core.DefenseTypeMagic,
		MinDmg:      2,
		MaxDmg:      2,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagNoOnDamageDealt | core.SpellFlagProc,
		Trigger: core.ProcTrigger{
			Name:               "Fiery Retributer",
			ActionID:           core.ActionID{ItemID: 19968},
			Callback:           core.CallbackOnSpellHitDealt,
			ProcMask:           core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial,
			Outcome:            core.OutcomeLanded,
			RequireDamageDealt: true,
			ProcChance:         1,
		},
	})

	// Gives a chance when your harmful spells land to reduce the magical resistances of your spell targets by
	// 50 for 8s.
	// https://www.wowhead.com/forever/spell=25768
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback:           core.CallbackOnCastComplete,
		ProcMask:           core.ProcMaskSpellDamage,
		Outcome:            core.OutcomeEmpty,
		RequireDamageDealt: false,
	}, []shared.ItemVariant{
		{ItemID: 21128, ItemName: "Staff of the Qiraji Prophets"},
	})

	// Gives a chance when your harmful spells land to increase the damage of your spells and effects by 132
	// for 10s.
	// https://www.wowhead.com/forever/spell=25907
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback:           core.CallbackOnCastComplete,
		ProcMask:           core.ProcMaskSpellDamage,
		Outcome:            core.OutcomeEmpty,
		RequireDamageDealt: false,
		CanProcFromProcs:   true,
	}, []shared.ItemVariant{
		{ItemID: 21190, ItemName: "Wrath of Cenarius"},
	})

	// When struck in combat has a chance of increasing your armor by 800 for 10s.
	// https://www.wowhead.com/forever/spell=35078
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback:           core.CallbackOnSpellHitTaken,
		ProcMask:           core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
	}, []shared.ItemVariant{
		{ItemID: 29297, ItemName: "Band of the Eternal Defender"},
	})

	// Chance on hit to increase your attack power by 160 for 10 seconds.
	// https://www.wowhead.com/forever/spell=35081
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback:           core.CallbackOnSpellHitDealt,
		ProcMask:           core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
	}, []shared.ItemVariant{
		{ItemID: 29301, ItemName: "Band of the Eternal Champion"},
	})

	// Your offensive spells have a chance on hit to increase your spell damage by 95 for 10 secs.
	// https://www.wowhead.com/forever/spell=35084
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback:           core.CallbackOnCastComplete,
		ProcMask:           core.ProcMaskSpellDamage,
		Outcome:            core.OutcomeEmpty,
		RequireDamageDealt: false,
	}, []shared.ItemVariant{
		{ItemID: 29305, ItemName: "Band of the Eternal Sage"},
	})

	// Your healing and damage spells have a chance to increase your healing by up to 175 and damage by up to
	// 59 for 10 secs.
	// https://www.wowhead.com/forever/spell=35087
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback:           core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:           core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: false,
	}, []shared.ItemVariant{
		{ItemID: 29309, ItemName: "Band of the Eternal Restorer"},
	})

	// Have a 2% chance when struck in combat of increasing armor by 350 for 15s.
	// https://www.wowhead.com/forever/spell=10342
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback:           core.CallbackOnSpellHitTaken,
		ProcMask:           core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial,
		Outcome:            core.OutcomeLanded,
		RequireDamageDealt: true,
	}, []shared.ItemVariant{
		{ItemID: 185986, ItemName: "Communal Stone of Durability"},
	})

	// Skipped
	// Not simulated: Mask of Thero-shan: "Stealth 5" (17746) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=17746
	// Not simulated: Nightscape Boots: "Stealth 5" (17746) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=17746
	// Not simulated: Catseye Ultra Goggles: "Stealth Detection" (12418) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=12418
	// Not simulated: Carrot on a Stick: "Mount Speed" (48777) - ignored aura type 130
	// https://www.wowhead.com/forever/spell=48777
	// Not simulated: Pip's Skinner: "Pip's Skinner" (16718) - ignored aura type 30
	// https://www.wowhead.com/forever/spell=16718
	// Not simulated: Spectral Essence: "Visions of the Past" (17623) - ignored aura type 19
	// https://www.wowhead.com/forever/spell=17623
	// Not simulated: Knight-Lieutenant's Dragonhide Gloves: "Stealth Detection" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Marshal's Dragonhide Gauntlets: "Stealth Detection" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Blood Guard's Dragonhide Gauntlets: "Stealth Detection" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: General's Dragonhide Gloves: "Stealth Detection" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Shard of the Defiler: "Echo of Archimonde" (21079) - ignored aura type 56
	// https://www.wowhead.com/forever/spell=21079
	// Not simulated: Mark of Resolution: "Stout Heart" (21958) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=21958
	// Not simulated: The Eye of Divinity: "Eye of Divinity" (23101) - ignored aura type 19
	// https://www.wowhead.com/forever/spell=23101
	// Not simulated: Nat Pagle's Extreme Angler FC-5000: "Fishing Skill +25" (8082) - ignored aura type 30
	// https://www.wowhead.com/forever/spell=8082
	// Not simulated: Zulian Slicer: "Improved Skinning" (24591) - ignored aura type 30
	// https://www.wowhead.com/forever/spell=24591
	// Not simulated: Nat Pagle's Extreme Anglin' Boots: "Fishing Skill +5" (7823) - ignored aura type 30
	// https://www.wowhead.com/forever/spell=7823
	// Not simulated: Arcanite Fishing Pole: "Fishing Skill +35" (24301) - ignored aura type 30
	// https://www.wowhead.com/forever/spell=24301
	// Not simulated: Lucky Fishing Hat: "Fishing Skill +5" (7823) - ignored aura type 30
	// https://www.wowhead.com/forever/spell=7823
	// Not simulated: Pirate's Eye Patch: "Fear Resistance 4" (24351) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=24351
	// Not simulated: Bloodvine Lens: "Stealth Detection" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Whisperwalk Boots: "Stealth 5" (17746) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=17746
	// Not simulated: Figurine - Black Pearl Panther: "Stealth 3" (26578) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=26578
	// Not simulated: Darkmantle Boots: "Stealth +8" (27037) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=27037
	// Not simulated: Blood Guard's Dragonhide Grips: "Stealth Detection" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Boots of Displacement: "Stealth +8" (27037) - ignored aura type 154
	// https://www.wowhead.com/forever/spell=27037
	// Not simulated: Stygian Buckler: "Stygian Grasp" (29164) - ignored aura type 33
	// https://www.wowhead.com/forever/spell=29164
	// Not simulated: Knight-Lieutenant's Dragonhide Grips: "Stealth Detection" (23217) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=23217
	// Not simulated: Ultra-Spectropic Detection Goggles: "Gas Cloud Tracking" (30645) - ignored aura type 44
	// https://www.wowhead.com/forever/spell=30645
	// Not simulated: Foreman's Enchanted Helmet: "Increased Stun Resist +10%" (40386) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=40386
	// Not simulated: Foreman's Reinforced Helmet: "Foreman's Reinforced Helmet" (30519) - ignored aura type 117
	// https://www.wowhead.com/forever/spell=30519
	// Not simulated: Seth's Graphite Fishing Pole: "Fishing Skill +20" (7826) - ignored aura type 30
	// https://www.wowhead.com/forever/spell=7826
	// Not simulated: Blade of the Unyielding: "Unyielding Knights" (38162) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=38162
	// Not simulated: Rod of the Unyielding: "Unyielding Knights" (38162) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=38162
	// Not simulated: Evoker's Helmet of Second Sight: "Spectrecles" (39841) - ignored aura type 19
	// https://www.wowhead.com/forever/spell=39841
	// Not simulated: Overlord's Helmet of Second Sight: "Spectrecles" (39841) - ignored aura type 19
	// https://www.wowhead.com/forever/spell=39841
	// Not simulated: Stalker's Helmet of Second Sight: "Spectrecles" (39841) - ignored aura type 19
	// https://www.wowhead.com/forever/spell=39841
	// Not simulated: Shamanistic Helmet of Second Sight: "Spectrecles" (39841) - ignored aura type 19
	// https://www.wowhead.com/forever/spell=39841
	// Not simulated: Stealther's Helmet of Second Sight: "Spectrecles" (39841) - ignored aura type 19
	// https://www.wowhead.com/forever/spell=39841
	// Not simulated: Druidic Helmet of Second Sight: "Spectrecles" (39841) - ignored aura type 19
	// https://www.wowhead.com/forever/spell=39841
	// Not simulated: Skybreaker Whip: "Mount Speed" (48776) - ignored aura type 172
	// https://www.wowhead.com/forever/spell=48776
	// Not simulated: Weather-Beaten Fishing Hat: "Fishing Skill +5" (7823) - ignored aura type 30
	// https://www.wowhead.com/forever/spell=7823
}
