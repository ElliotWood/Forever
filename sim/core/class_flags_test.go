package core

import "testing"

func TestClassFlagsIsZero(t *testing.T) {
	cases := []struct {
		name  string
		flags ClassFlags
		want  bool
	}{
		{"empty", ClassFlags{}, true},
		{"family only", ClassFlags{Family: 11}, false},
		{"mask only", ClassFlags{Mask: [4]uint32{0, 0, 0, 0x8}}, false},
		{"family and mask", ClassFlags{Family: 11, Mask: [4]uint32{0x1}}, false},
	}

	for _, tc := range cases {
		if got := tc.flags.IsZero(); got != tc.want {
			t.Errorf("%s: IsZero() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestClassFlagsMatches(t *testing.T) {
	cases := []struct {
		name string
		a    ClassFlags
		b    ClassFlags
		want bool
	}{
		{"same word overlaps", ClassFlags{Family: 11, Mask: [4]uint32{0x3}}, ClassFlags{Family: 11, Mask: [4]uint32{0x2}}, true},
		{"later word overlaps", ClassFlags{Family: 11, Mask: [4]uint32{0, 0, 0x40}}, ClassFlags{Family: 11, Mask: [4]uint32{0, 0, 0x60}}, true},
		{"same family, disjoint bits", ClassFlags{Family: 11, Mask: [4]uint32{0x1}}, ClassFlags{Family: 11, Mask: [4]uint32{0x2}}, false},
		{"same bits, other family", ClassFlags{Family: 11, Mask: [4]uint32{0x3}}, ClassFlags{Family: 6, Mask: [4]uint32{0x3}}, false},
		{"bits in different words", ClassFlags{Family: 11, Mask: [4]uint32{0x1}}, ClassFlags{Family: 11, Mask: [4]uint32{0, 0x1}}, false},
		{"against a zero set", ClassFlags{Family: 11, Mask: [4]uint32{0x3}}, ClassFlags{}, false},
	}

	for _, tc := range cases {
		if got := tc.a.Matches(tc.b); got != tc.want {
			t.Errorf("%s: Matches() = %v, want %v", tc.name, got, tc.want)
		}
		if got := tc.b.Matches(tc.a); got != tc.want {
			t.Errorf("%s: reversed Matches() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestClassFlagsOr(t *testing.T) {
	cases := []struct {
		name string
		a    ClassFlags
		b    ClassFlags
		want ClassFlags
	}{
		{
			"same family",
			ClassFlags{Family: 11, Mask: [4]uint32{0x1, 0, 0x40}},
			ClassFlags{Family: 11, Mask: [4]uint32{0x2, 0x8}},
			ClassFlags{Family: 11, Mask: [4]uint32{0x3, 0x8, 0x40}},
		},
		{
			"zero left side takes the other family",
			ClassFlags{},
			ClassFlags{Family: 11, Mask: [4]uint32{0x2}},
			ClassFlags{Family: 11, Mask: [4]uint32{0x2}},
		},
		{
			"zero right side",
			ClassFlags{Family: 11, Mask: [4]uint32{0x2}},
			ClassFlags{},
			ClassFlags{Family: 11, Mask: [4]uint32{0x2}},
		},
		{"both zero", ClassFlags{}, ClassFlags{}, ClassFlags{}},
	}

	for _, tc := range cases {
		if got := tc.a.Or(tc.b); got != tc.want {
			t.Errorf("%s: Or() = %+v, want %+v", tc.name, got, tc.want)
		}
	}
}

func TestClassFlagsOrPanicsOnFamilyMismatch(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Or() of two families did not panic")
		}
	}()

	ClassFlags{Family: 11, Mask: [4]uint32{0x1}}.Or(ClassFlags{Family: 6, Mask: [4]uint32{0x1}})
}
