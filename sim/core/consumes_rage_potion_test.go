package core

import "testing"

// A rage potion on a class that has no rage bar reached rageBar.AddRage through a nil unit and
// crashed the whole sim. The site offers every potion to every character, so this was one
// dropdown away for any rogue or paladin.
func TestRagePotionOnAClassWithoutRageDoesNothing(t *testing.T) {
	character := &Character{}
	if character.HasRageBar() {
		t.Fatal("a bare character should not have a rage bar")
	}
	// Registering used to panic here rather than return an empty cooldown.
	if mcd := makeRageConsumableMCD(13442, character, nil); mcd.Spell != nil {
		t.Error("a rage potion registered a cooldown on a class that cannot hold rage")
	}
}
