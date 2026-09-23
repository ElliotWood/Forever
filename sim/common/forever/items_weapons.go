package forever

import (
	"time"

	"github.com/wowsims/forever/sim/common/itemhelpers"
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/stats"
)

func init() {
	// The Lobotomizer: Forever's Brain Damage (client 1290950) wounds for 250 +-40% (200 to 300)
	// and slows the target's casting. 0.4 PPM and the magic hit table as on master.
	itemhelpers.CreateWeaponProcSpell(itemhelpers.WeaponProcSpell{
		ItemID: 19324,
		Name:   "The Lobotomizer",
		PPM:    0.4,
		Spell: func(character *core.Character) *core.Spell {
			return character.GetOrRegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 1290950},
				ProcMask:    core.ProcMaskEmpty,
				SpellSchool: core.SpellSchoolPhysical,
				DefenseType: core.DefenseTypeMelee,
				Flags:       core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, sim.Roll(200, 300), spell.OutcomeMagicHitAndCrit)
				},
			})
		},
	})

	// Bonereaver's Edge
	itemhelpers.CreateWeaponProcTrigger(itemhelpers.WeaponProcTrigger{
		ItemID: 17076,
		Name:   "Bonereaver's Edge",
		PPM:    2,
		Handler: func(character *core.Character) core.ProcHandler {
			arpAura := core.MakeStackingAura(
				character,
				core.StackingStatAura{
					Aura: core.Aura{
						Label:     "Bonereaver's Edge",
						ActionID:  core.ActionID{SpellID: 21153},
						Duration:  time.Second * 10,
						MaxStacks: 3,
					},
					BonusPerStack: stats.Stats{
						stats.ArmorPenetration: 700,
					},
				},
			)

			return func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				arpAura.Activate(sim)
				arpAura.AddStack(sim)
			}
		},
	})

	// Bashguuder, Rivenspike: chance on hit, Puncture Armor (client 17315: -100 armor, 3 stacks,
	// 30s). The client has no proc rate; 2 PPM is master's (Armaments Discord).
	for itemID, name := range map[int32]string{13204: "Bashguuder", 13286: "Rivenspike"} {
		itemhelpers.CreateWeaponProcTrigger(itemhelpers.WeaponProcTrigger{
			ItemID: itemID,
			Name:   name,
			PPM:    2,
			Handler: func(character *core.Character) core.ProcHandler {
				auras := character.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
					return target.GetOrRegisterAura(core.Aura{
						Label:     "Puncture Armor",
						ActionID:  core.ActionID{SpellID: 17315},
						Duration:  time.Second * 30,
						MaxStacks: 3,
						OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
							aura.Unit.AddStatDynamic(sim, stats.Armor, -100*float64(newStacks-oldStacks))
						},
					})
				})

				return func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					aura := auras.Get(result.Target)
					aura.Activate(sim)
					aura.AddStack(sim)
				}
			},
		})
	}

}
