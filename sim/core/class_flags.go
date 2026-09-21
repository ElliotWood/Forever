package core

import "fmt"

// ClassFlags is the client's SpellClassOptions: the family (SpellClassSet) and the four
// 32-bit words of SpellClassMask. A talent effect carries the same shape (EffectSpellClassMask)
// naming the spells it modifies.
type ClassFlags struct {
	Family int32
	Mask   [4]uint32
}

// Whether the spell carries no class options at all. The client leaves both the family and every
// mask word at zero on the spells no talent addresses by family.
func (f ClassFlags) IsZero() bool {
	return f.Family == 0 && f.Mask == [4]uint32{}
}

// Whether the two sets name at least one spell in common. Families are separate namespaces, so
// the same bit means a different spell in each and a cross-family overlap is never a match.
func (f ClassFlags) Matches(o ClassFlags) bool {
	if f.Family != o.Family {
		return false
	}

	for i := range f.Mask {
		if f.Mask[i]&o.Mask[i] != 0 {
			return true
		}
	}

	return false
}

// The union of the two sets. A zero set carries no family of its own and takes the other's, which
// is what makes Or usable as an accumulator; any other family mismatch is a bug in the caller,
// since the union of two families cannot be expressed in one ClassFlags.
func (f ClassFlags) Or(o ClassFlags) ClassFlags {
	if f.IsZero() {
		return o
	}

	if o.IsZero() {
		return f
	}

	if f.Family != o.Family {
		panic(fmt.Sprintf("cannot combine ClassFlags of families %d and %d", f.Family, o.Family))
	}

	union := ClassFlags{Family: f.Family}
	for i := range union.Mask {
		union.Mask[i] = f.Mask[i] | o.Mask[i]
	}

	return union
}
