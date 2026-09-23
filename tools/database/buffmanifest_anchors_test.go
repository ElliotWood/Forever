package database

// The manifest pins the spell each buff reads, so generation needs no client database. What pinned
// it is the rank resolution below, which reads SkillLineAbility, the rune enchantments and the trait
// trees - tables the store does not capture - so this is where a new client is checked against the
// pins. Skips without tools/database/wowsims.db.

import (
	"testing"

	"github.com/wowsims/forever/tools/database/buffmanifest"
)

func TestManifestAnchorsMatchTheClient(t *testing.T) {
	helper := openBuffTestDB(t)

	rows, err := ResolveBuffManifest(helper)
	if err != nil {
		t.Fatalf("resolving the manifest: %v", err)
	}

	for i, row := range rows {
		spec := buffmanifest.Manifest[i]
		if spec.SpellID != row.SpellID {
			t.Errorf("%s: the manifest pins spell %d, the client resolves %d", spec.Field, spec.SpellID, row.SpellID)
		}
		if spec.CastID != 0 && spec.CastID != row.CastSpellID {
			t.Errorf("%s: the manifest pins cast %d, the client resolves %d", spec.Field, spec.CastID, row.CastSpellID)
		}
		if row.Kind == buffmanifest.KindExternalCD && row.CastSpellID != row.SpellID && spec.CastID == 0 {
			t.Errorf("%s: the cast %d times the cooldown and states what the aura %d does not, so the manifest has to pin it",
				spec.Field, row.CastSpellID, row.SpellID)
		}
		if spec.Talent != nil && spec.Talent.SpellID != row.TalentSpellID {
			t.Errorf("%s: the manifest pins talent %d, the owner's tree prices %d", spec.Field, spec.Talent.SpellID, row.TalentSpellID)
		}
	}
}
