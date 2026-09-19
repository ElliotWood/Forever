package forever

import (
	"time"

	"github.com/wowsims/forever/sim/common/shared"
	"github.com/wowsims/forever/sim/core"
)

func RegisterAllOnUseCds() {

	//
	// shared.NewSimpleStatActive(833) // Lifestone - https://www.wowhead.com/forever/spell=17712
	// shared.NewSimpleStatActive(4696) // Lapidis Tankard of Tidesippe - https://www.wowhead.com/forever/spell=1135
	// shared.NewSimpleStatActive(7734) // Six Demon Bag - https://www.wowhead.com/forever/spell=14537
	// shared.NewSimpleStatActive(8367) // Dragonscale Breastplate - https://www.wowhead.com/forever/spell=10618
	// shared.NewSimpleStatActive(11808) // Circle of Flame - https://www.wowhead.com/forever/spell=17447
	// shared.NewSimpleStatActive(11819) // Second Wind - https://www.wowhead.com/forever/spell=15604
	// shared.NewSimpleStatActive(11832) // Burst of Knowledge - https://www.wowhead.com/forever/spell=15646
	// shared.NewSimpleStatActive(11905) // Linken's Boomerang - https://www.wowhead.com/forever/spell=15712
	// shared.NewSimpleStatActive(13143) // Mark of the Dragon Lord - https://www.wowhead.com/forever/spell=17252
	// shared.NewSimpleStatActive(13171) // Smokey's Lighter - https://www.wowhead.com/forever/spell=17283
	// shared.NewSimpleStatActive(13213) // Smolderweb's Eye - https://www.wowhead.com/forever/spell=17330
	// shared.NewSimpleStatActive(13379) // Piccolo of the Flaming Fire - https://www.wowhead.com/forever/spell=18400
	// shared.NewSimpleStatActive(13515) // Ramstein's Lightning Bolts - https://www.wowhead.com/forever/spell=17668
	// shared.NewSimpleStatActive(13937) // Headmaster's Charge - https://www.wowhead.com/forever/spell=18264
	// shared.NewSimpleStatActive(14134) // Cloak of Fire - https://www.wowhead.com/forever/spell=18364
	// shared.NewSimpleStatActive(14152) // Robe of the Archmage - https://www.wowhead.com/forever/spell=18385
	// shared.NewSimpleStatActive(14153) // Robe of the Void - https://www.wowhead.com/forever/spell=18386
	// shared.NewSimpleStatActive(16768) // Furbolg Medicine Pouch - https://www.wowhead.com/forever/spell=20631
	// shared.NewSimpleStatActive(17744) // Heart of Noxxion - https://www.wowhead.com/forever/spell=21954
	// shared.NewSimpleStatActive(17759) // Mark of Resolution - https://www.wowhead.com/forever/spell=21956
	// shared.NewSimpleStatActive(18634) // Gyrofreeze Ice Reflector - https://www.wowhead.com/forever/spell=23131
	// shared.NewSimpleStatActive(18637) // Major Recombobulator - https://www.wowhead.com/forever/spell=23064
	// shared.NewSimpleStatActive(18638) // Hyper-Radiant Flame Reflector - https://www.wowhead.com/forever/spell=23097
	// shared.NewSimpleStatActive(18639) // Ultra-Flash Shadow Reflector - https://www.wowhead.com/forever/spell=23132
	// shared.NewSimpleStatActive(18986) // Ultrasafe Transporter: Gadgetzan - https://www.wowhead.com/forever/spell=23453
	// shared.NewSimpleStatActive(19024) // Arena Grand Master - https://www.wowhead.com/forever/spell=23506
	// shared.NewSimpleStatActive(19336) // Arcane Infused Gem - https://www.wowhead.com/forever/spell=23721
	// shared.NewSimpleStatActive(19340) // Rune of Metamorphosis - https://www.wowhead.com/forever/spell=23724
	// shared.NewSimpleStatActive(19341) // Lifegiving Gem - https://www.wowhead.com/forever/spell=23725
	// shared.NewSimpleStatActive(19342) // Venomous Totem - https://www.wowhead.com/forever/spell=23726
	// shared.NewSimpleStatActive(19930) // Mar'li's Eye - https://www.wowhead.com/forever/spell=24268
	// shared.NewSimpleStatActive(19949) // Zandalarian Hero Medallion - https://www.wowhead.com/forever/spell=24661
	// shared.NewSimpleStatActive(19950) // Zandalarian Hero Charm - https://www.wowhead.com/forever/spell=24658
	// shared.NewSimpleStatActive(19951) // Gri'lek's Charm of Might - https://www.wowhead.com/forever/spell=24571
	// shared.NewSimpleStatActive(19953) // Renataki's Charm of Beasts - https://www.wowhead.com/forever/spell=24531
	// shared.NewSimpleStatActive(19954) // Renataki's Charm of Trickery - https://www.wowhead.com/forever/spell=24532
	// shared.NewSimpleStatActive(19956) // Wushoolay's Charm of Spirits - https://www.wowhead.com/forever/spell=24499
	// shared.NewSimpleStatActive(19992) // Devilsaur Tooth - https://www.wowhead.com/forever/spell=24353
	// shared.NewSimpleStatActive(20071) // Talisman of Arathor - https://www.wowhead.com/forever/spell=23991
	// shared.NewSimpleStatActive(20072) // Defiler's Talisman - https://www.wowhead.com/forever/spell=23991
	// shared.NewSimpleStatActive(20525) // Earthen Sigil - https://www.wowhead.com/forever/spell=24884
	// shared.NewSimpleStatActive(20534) // Abyss Shard - https://www.wowhead.com/forever/spell=25112
	// shared.NewSimpleStatActive(21115) // Defiler's Talisman - https://www.wowhead.com/forever/spell=25746
	// shared.NewSimpleStatActive(21117) // Talisman of Arathor - https://www.wowhead.com/forever/spell=25746
	// shared.NewSimpleStatActive(21181) // Grace of Earth - https://www.wowhead.com/forever/spell=25892
	// shared.NewSimpleStatActive(21488) // Fetish of Chitinous Spikes - https://www.wowhead.com/forever/spell=26168
	// shared.NewSimpleStatActive(21625) // Scarab Brooch - https://www.wowhead.com/forever/spell=26467
	// shared.NewSimpleStatActive(21647) // Fetish of the Sand Reaver - https://www.wowhead.com/forever/spell=26400
	// shared.NewSimpleStatActive(21685) // Petrified Scarab - https://www.wowhead.com/forever/spell=26463
	// shared.NewSimpleStatActive(21891) // Shard of the Fallen Star - https://www.wowhead.com/forever/spell=26789
	// shared.NewSimpleStatActive(23001) // Eye of Diminution - https://www.wowhead.com/forever/spell=28862
	// shared.NewSimpleStatActive(23027) // Warmth of Forgiveness - https://www.wowhead.com/forever/spell=28760
	// shared.NewSimpleStatActive(23558) // The Burrower's Shell - https://www.wowhead.com/forever/spell=29506
	// shared.NewSimpleStatActive(23587) // Mirren's Drinking Hat - https://www.wowhead.com/forever/spell=29830
	// shared.NewSimpleStatActive(23824) // Rocket Boots Xtreme - https://www.wowhead.com/forever/spell=51582
	// shared.NewSimpleStatActive(23825) // Nigh Invulnerability Belt - https://www.wowhead.com/forever/spell=30458
	// shared.NewSimpleStatActive(23835) // Gnomish Poultryizer - https://www.wowhead.com/forever/spell=30507
	// shared.NewSimpleStatActive(23836) // Goblin Rocket Launcher - https://www.wowhead.com/forever/spell=46567
	// shared.NewSimpleStatActive(25786) // Hypnotist's Watch - https://www.wowhead.com/forever/spell=32599
	// shared.NewSimpleStatActive(30696) // Scourgebane - https://www.wowhead.com/forever/spell=37360
	// shared.NewSimpleStatActive(30847) // X-52 Rocket Helmet - https://www.wowhead.com/forever/spell=37896
	// shared.NewSimpleStatActive(33820) // Weather-Beaten Fishing Hat - https://www.wowhead.com/forever/spell=43699
	// shared.NewSimpleStatActive(34140) // Dark Iron Tankard - https://www.wowhead.com/forever/spell=44540
	// shared.NewSimpleStatActive(35581) // Rocket Boots Xtreme Lite - https://www.wowhead.com/forever/spell=51582
	// shared.NewSimpleStatActive(185988) // Communal Stone of Stoicism - https://www.wowhead.com/forever/spell=352340

	// Agility / Stamina / Strength
	shared.NewSimpleStatActive(15873) // Ragged John's Neverending Cup - https://www.wowhead.com/forever/spell=20587

	// ArcaneDamage
	shared.NewSimpleStatActive(19959) // Hazza'rah's Charm of Magic - https://www.wowhead.com/forever/spell=24544

	// ArcaneResistance / FireResistance / FrostResistance / NatureResistance / ShadowResistance
	shared.NewSimpleStatActive(15867) // Prismcharm - https://www.wowhead.com/forever/spell=19638
	shared.NewSimpleStatActive(23042) // Loatheb's Reflection - https://www.wowhead.com/forever/spell=28778

	// Armor
	shared.NewSimpleStatActive(12532) // Spire of the Stoneshaper - https://www.wowhead.com/forever/spell=16470
	shared.NewSimpleStatActive(19345) // Aegis of Preservation - https://www.wowhead.com/forever/spell=23780

	// Armor / DefenseRating
	// Zandalarian Hero Badge - https://www.wowhead.com/forever/spell=24574
	shared.NewStackingStatBonusCD(shared.StackingStatBonusCD{
		Name:                  "Zandalarian Hero Badge",
		ID:                    19948,
		CD:                    time.Millisecond * 120000,
		Callback:              core.CallbackOnSpellHitTaken,
		ProcMask:              core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:               core.OutcomeLanded,
		RequireDamageDealt:    true,
		TrinketLimitsDuration: true,
	})

	// ArmorPenetration
	// Badge of the Swarmguard - https://www.wowhead.com/forever/spell=26480
	shared.NewStackingStatBonusCD(shared.StackingStatBonusCD{
		Name:                  "Badge of the Swarmguard",
		ID:                    21670,
		CD:                    time.Millisecond * 180000,
		Callback:              core.CallbackOnSpellHitDealt,
		ProcMask:              core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:               core.OutcomeLanded,
		RequireDamageDealt:    true,
		TrinketLimitsDuration: true,
	})

	// AttackPower / MeleeHitRating / RangedAttackPower
	shared.NewSimpleStatActive(19991) // Devilsaur Eye - https://www.wowhead.com/forever/spell=24352

	// AttackPower / RangedAttackPower
	shared.NewSimpleStatActive(14554) // Cloudkeeper Legplates - https://www.wowhead.com/forever/spell=18787
	shared.NewSimpleStatActive(21180) // Earthstrike - https://www.wowhead.com/forever/spell=25891
	shared.NewSimpleStatActive(23041) // Slayer's Crest - https://www.wowhead.com/forever/spell=28777
	shared.NewSimpleStatActive(25628) // Ogre Mauler's Badge - https://www.wowhead.com/forever/spell=32362
	shared.NewSimpleStatActive(25633) // Uniting Charm - https://www.wowhead.com/forever/spell=32362
	shared.NewSimpleStatActive(25937) // Terokkar Tablet of Precision - https://www.wowhead.com/forever/spell=39200
	shared.NewSimpleStatActive(25994) // Rune of Force - https://www.wowhead.com/forever/spell=32955
	shared.NewSimpleStatActive(28041) // Bladefist's Breadth - https://www.wowhead.com/forever/spell=33667
	shared.NewSimpleStatActive(29776) // Core of Ar'kelos - https://www.wowhead.com/forever/spell=35733
	shared.NewSimpleStatActive(31617) // Ancient Draenei War Talisman - https://www.wowhead.com/forever/spell=33667
	shared.NewSimpleStatActive(32654) // Crystalforged Trinket - https://www.wowhead.com/forever/spell=40724

	// BlockRating
	shared.NewSimpleStatActive(30300) // Dabiri's Enigma - https://www.wowhead.com/forever/spell=36372

	// BlockValue
	shared.NewSimpleStatActive(23040) // Glyph of Deflection - https://www.wowhead.com/forever/spell=28773

	// DefenseRating
	shared.NewSimpleStatActive(25996) // Emblem of Perseverance - https://www.wowhead.com/forever/spell=32957

	// DodgeRating
	shared.NewSimpleStatActive(25787) // Charm of Alacrity - https://www.wowhead.com/forever/spell=32600

	// FireDamage
	shared.NewSimpleStatActive(20036) // Fire Ruby - https://www.wowhead.com/forever/spell=24389

	// FireResistance
	shared.NewSimpleStatActive(13164) // Heart of the Scale - https://www.wowhead.com/forever/spell=17275

	// HealingPower / SpellDamage
	shared.NewSimpleStatActive(18820) // Talisman of Ephemeral Power - https://www.wowhead.com/forever/spell=23271
	shared.NewSimpleStatActive(19344) // Natural Alignment Crystal - https://www.wowhead.com/forever/spell=23734
	shared.NewSimpleStatActive(19990) // Blessed Prayer Beads - https://www.wowhead.com/forever/spell=24354
	shared.NewSimpleStatActive(20636) // Hibernation Crystal - https://www.wowhead.com/forever/spell=24998
	shared.NewSimpleStatActive(22268) // Draconic Infused Emblem - https://www.wowhead.com/forever/spell=27675
	// Talisman of Ascendance - https://www.wowhead.com/forever/spell=28200
	shared.NewStackingStatBonusCD(shared.StackingStatBonusCD{
		Name:                  "Talisman of Ascendance",
		ID:                    22678,
		CD:                    time.Millisecond * 60000,
		Callback:              core.CallbackOnSpellHitDealt,
		ProcMask:              core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage,
		Outcome:               core.OutcomeLanded,
		RequireDamageDealt:    false,
		TrinketLimitsDuration: true,
	})
	shared.NewSimpleStatActive(23046) // The Restrained Essence of Sapphiron - https://www.wowhead.com/forever/spell=28779
	shared.NewSimpleStatActive(23047) // Eye of the Dead - https://www.wowhead.com/forever/spell=28780
	shared.NewSimpleStatActive(25619) // Glowing Crystal Insignia - https://www.wowhead.com/forever/spell=32355
	shared.NewSimpleStatActive(25620) // Ancient Crystal Talisman - https://www.wowhead.com/forever/spell=32355
	shared.NewSimpleStatActive(25634) // Oshu'gun Relic - https://www.wowhead.com/forever/spell=32367
	shared.NewSimpleStatActive(25936) // Terokkar Tablet of Vim - https://www.wowhead.com/forever/spell=39201
	shared.NewSimpleStatActive(25995) // Star of Sha'naar - https://www.wowhead.com/forever/spell=32956
	shared.NewSimpleStatActive(28040) // Vengeance of the Illidari - https://www.wowhead.com/forever/spell=33662
	shared.NewSimpleStatActive(30293) // Heavenly Inspiration - https://www.wowhead.com/forever/spell=36347
	shared.NewSimpleStatActive(31615) // Ancient Draenei Arcane Relic - https://www.wowhead.com/forever/spell=33662

	// Health
	shared.NewSimpleStatActive(28042) // Regal Protectorate - https://www.wowhead.com/forever/spell=33668

	// MeleeCritRating / SpellCritRating
	shared.NewSimpleStatActive(20512) // Sanctified Orb - https://www.wowhead.com/forever/spell=24865

	// MeleeHasteRating
	shared.NewSimpleStatActive(22954) // Kiss of the Spider - https://www.wowhead.com/forever/spell=28866

	// MeleeHasteRating / SpellHasteRating
	shared.NewSimpleStatActive(19343) // Scrolls of Blinding Light - https://www.wowhead.com/forever/spell=23733

	// Skipped
	// Not simulated: Lei of Lilies: "Conjure Lily Root" (18831) - ignored effect type 24
	// https://www.wowhead.com/forever/spell=18831
	// Not simulated: Staff of Conjuring: "Conjure Food" (8736) - ignored effect type 24
	// https://www.wowhead.com/forever/spell=8736
	// Not simulated: Orb of Deception: "Orb of Deception" (16739) - ignored aura type 56
	// https://www.wowhead.com/forever/spell=16739
	// Not simulated: Spider Belt: "Freedom" (9774) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=9774
	// Not simulated: Enchanted Moonstalker Cloak: "Form of the Moonstalker" (6298) - ignored aura type 56
	// https://www.wowhead.com/forever/spell=6298
	// Not simulated: Ornate Mithril Boots: "Freedom" (9774) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=9774
	// Not simulated: Chained Essence of Eranikus: "Poison Cloud" (12766) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=12766
	// Not simulated: Mithril Mechanical Dragonling: "Mithril Mechanical Dragonling" (12749) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=12749
	// Not simulated: Bloodsail Admiral's Hat: "Summon Blood Parrot" (17567) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=17567
	// Not simulated: Book of the Dead: "Summon Skeleton" (17490) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=17490
	// Not simulated: Cannonball Runner: "Summon Crimson Cannon" (6251) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=6251
	// Not simulated: Barov Peasant Caller: "Death by Peasant" (18307) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=18307
	// Not simulated: Barov Peasant Caller: "Death by Peasant" (18308) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=18308
	// Not simulated: Arcanite Dragonling: "Arcanite Dragonling" (19804) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=19804
	// Not simulated: Ancient Cornerstone Grimoire: "Summon Skeleton" (17490) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=17490
	// Not simulated: Frostwolf Insignia Rank 1: "Recall" (22563) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=22563
	// Not simulated: Stormpike Insignia Rank 1: "Recall" (22564) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=22564
	// Not simulated: Stormpike Insignia Rank 2: "Recall" (22564) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=22564
	// Not simulated: Stormpike Insignia Rank 3: "Recall" (22564) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=22564
	// Not simulated: Stormpike Insignia Rank 4: "Recall" (22564) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=22564
	// Not simulated: Stormpike Insignia Rank 5: "Recall" (22564) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=22564
	// Not simulated: Frostwolf Insignia Rank 2: "Recall" (22563) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=22563
	// Not simulated: Frostwolf Insignia Rank 3: "Recall" (22563) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=22563
	// Not simulated: Frostwolf Insignia Rank 4: "Recall" (22563) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=22563
	// Not simulated: Frostwolf Insignia Rank 5: "Recall" (22563) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=22563
	// Not simulated: Dimensional Ripper - Everlook: "Everlook Transporter" (23442) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=23442
	// Not simulated: Hook of the Master Angler: "Master Angler" (24347) - ignored aura type 56
	// https://www.wowhead.com/forever/spell=24347
	// Not simulated: Enamored Water Spirit: "Mana Spring Totem" (24854) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=24854
	// Not simulated: Defender of the Timbermaw: "Defender of the Timbermaw" (26066) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=26066
	// Not simulated: Vanquished Tentacle of C'Thun: "Tentacle Call" (26391) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=26391
	// Not simulated: Figurine - Jade Owl: "Jade Owl" (26551) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=26551
	// Not simulated: Figurine - Golden Hare: "Golden Hare" (26571) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=26571
	// Not simulated: Figurine - Black Pearl Panther: "Black Pearl Panther" (26576) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=26576
	// Not simulated: Figurine - Truesilver Crab: "Truesilver Crab" (26581) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=26581
	// Not simulated: Figurine - Truesilver Boar: "Truesilver Boar" (26593) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=26593
	// Not simulated: Figurine - Ruby Serpent: "Ruby Serpent" (26599) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=26599
	// Not simulated: Figurine - Emerald Owl: "Emerald Owl" (26600) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=26600
	// Not simulated: Figurine - Black Diamond Crab: "Black Diamond Crab" (26609) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=26609
	// Not simulated: Figurine - Dark Iron Scorpid: "Dark Iron Scorpid" (26614) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=26614
	// Not simulated: Carved Ogre Idol: "Red Ogre Costume" (30167) - ignored aura type 56
	// https://www.wowhead.com/forever/spell=30167
	// Not simulated: Hyper-Vision Goggles: "Hyper Vision" (30249) - ignored aura type 17
	// https://www.wowhead.com/forever/spell=30249
	// Not simulated: Everlasting Underspore Frond: "Everlasting Underspore Fronds" (33770) - ignored effect type 24
	// https://www.wowhead.com/forever/spell=33770
	// Not simulated: Insignia of the Alliance: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Dimensional Ripper - Area 52: "Area52 Transporter" (36890) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=36890
	// Not simulated: Ultrasafe Transporter: Toshley's Station: "Toshley's Station Transporter" (36941) - ignored effect type
	// Not simulated: 252
	// https://www.wowhead.com/forever/spell=36941
	// Not simulated: Overseer's Badge: "Netherwing Ally" (40811) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=40811
	// Not simulated: Captain's Badge: "Netherwing Ally" (40815) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=40815
	// Not simulated: Blessed Medallion of Karabor: "Teleport: Black Temple" (41234) - ignored effect type 252
	// https://www.wowhead.com/forever/spell=41234
	// Not simulated: Commander's Badge: "Netherwing Ally" (40815) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=40815
	// Not simulated: Don Carlos' Famous Hat: "Summon Coyote Spirit" (51149) - ignored effect type 28
	// https://www.wowhead.com/forever/spell=51149
	// Not simulated: Insignia of the Alliance: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Alliance: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Alliance: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Alliance: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Alliance: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Alliance: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Alliance: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Alliance: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Horde: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Horde: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Horde: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Horde: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Horde: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Horde: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Horde: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292
	// Not simulated: Insignia of the Horde: "PvP Trinket" (42292) - ignored aura type 77
	// https://www.wowhead.com/forever/spell=42292

	// SpellCritRating
	shared.NewSimpleStatActive(19952) // Gri'lek's Charm of Valor - https://www.wowhead.com/forever/spell=24498
	shared.NewSimpleStatActive(19957) // Hazza'rah's Charm of Destruction - https://www.wowhead.com/forever/spell=24543

	// SpellDamage
	shared.NewSimpleStatActive(30340) // Starkiller's Bauble - https://www.wowhead.com/forever/spell=36432

	// SpellDamage / SpellPenetration
	shared.NewSimpleStatActive(21473) // Eye of Moam - https://www.wowhead.com/forever/spell=26166

	// SpellHasteRating
	shared.NewSimpleStatActive(19339) // Mind Quickening Gem - https://www.wowhead.com/forever/spell=23723
	shared.NewSimpleStatActive(19955) // Wushoolay's Charm of Nature - https://www.wowhead.com/forever/spell=24542
	shared.NewSimpleStatActive(19958) // Hazza'rah's Charm of Healing - https://www.wowhead.com/forever/spell=24546

	// SpellHitRating
	shared.NewSimpleStatActive(19947) // Nat Pagle's Broken Reel - https://www.wowhead.com/forever/spell=24610

	// Strength
	shared.NewSimpleStatActive(20130) // Diamond Flask - https://www.wowhead.com/forever/spell=24427
}
