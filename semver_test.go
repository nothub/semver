package semver

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"
)

func TestParseValids(t *testing.T) {
	tests := []struct {
		input    string
		expected Version
	}{
		{
			input: "0.0.4",
			expected: Version{
				Major: 0,
				Minor: 0,
				Patch: 4,
			},
		},
		{
			input: "10.20.30",
			expected: Version{
				Major: 10,
				Minor: 20,
				Patch: 30,
			},
		},
		{
			input: "1.1.2-prerelease+meta",
			expected: Version{
				Major:      1,
				Minor:      1,
				Patch:      2,
				PreRelease: []string{"prerelease"},
				Build:      []string{"meta"},
			},
		},
		{
			input: "1.1.2+meta",
			expected: Version{
				Major: 1,
				Minor: 1,
				Patch: 2,
				Build: []string{"meta"},
			},
		},
		{
			input: "1.1.2+meta-valid",
			expected: Version{
				Major: 1,
				Minor: 1,
				Patch: 2,
				Build: []string{"meta-valid"},
			},
		},
		{
			input: "1.0.0-alpha",
			expected: Version{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: []string{"alpha"},
			},
		},
		{
			input: "1.0.0-alpha.beta.1",
			expected: Version{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: []string{"alpha", "beta", "1"},
			},
		},
		{
			input: "1.0.0-alpha.beta",
			expected: Version{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: []string{"alpha", "beta"},
			},
		},
		{
			input: "1.0.0-alpha.0valid",
			expected: Version{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: []string{"alpha", "0valid"},
			},
		},
		{
			input: "1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
			expected: Version{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: []string{"alpha-a", "b-c-somethinglong"},
				Build:      []string{"build", "1-aef", "1-its-okay"},
			},
		},
		{
			input: "10.2.3-DEV-SNAPSHOT",
			expected: Version{
				Major:      10,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"DEV-SNAPSHOT"},
			},
		},
		{
			input: "2.0.0+build.1848",
			expected: Version{
				Major: 2,
				Minor: 0,
				Patch: 0,
				Build: []string{"build", "1848"},
			},
		},
		{
			input: "2.0.1-alpha.1227",
			expected: Version{
				Major:      2,
				Minor:      0,
				Patch:      1,
				PreRelease: []string{"alpha", "1227"},
			},
		},
		{
			input: "1.2.3----RC-SNAPSHOT.12.9.1--.12+788",
			expected: Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"---RC-SNAPSHOT", "12", "9", "1--", "12"},
				Build:      []string{"788"},
			},
		},
		{
			input: "1.2.3----RC-SNAPSHOT.12.9.1--.12",
			expected: Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"---RC-SNAPSHOT", "12", "9", "1--", "12"},
			},
		},
		{
			input: "1.2.3----R-S.12.9.1--.12+meta",
			expected: Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"---R-S", "12", "9", "1--", "12"},
				Build:      []string{"meta"},
			},
		},
		{
			input: "1.0.0+0.build.1-rc.10000aaa-kk-0.1",
			expected: Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
				Build: []string{"0", "build", "1-rc", "10000aaa-kk-0", "1"},
			},
		},
		{
			input: "1.0.0-0A.is.legal",
			expected: Version{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: []string{"0A", "is", "legal"},
			},
		},
		{
			input: "9.8.7-whatever+meta",
			expected: Version{
				Major:      9,
				Minor:      8,
				Patch:      7,
				PreRelease: []string{"whatever"},
				Build:      []string{"meta"},
			},
		},
		{
			input: "9.8.7+whatever-meta",
			expected: Version{
				Major: 9,
				Minor: 8,
				Patch: 7,
				Build: []string{"whatever-meta"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := Parse(test.input)
			if err != nil {
				t.Errorf("error = %v", err)
				return
			}
			if !reflect.DeepEqual(test.expected, result) {
				t.Errorf("unexpected result:\nexpected = %+v\nactual   = %+v", test.expected, result)
			}
			if test.input != result.String() {
				t.Errorf("non equal string format:\nexpected = %s\nactual   = %s", test.input, result.String())
			}
		})
	}
}

func TestParseInvalids(t *testing.T) {
	tests := []string{
		"1",
		"1.2",
		"1.2.3-0123",
		"1.2.3-0123.0123",
		"1.1.2+.123",
		"+invalid",
		"-invalid",
		"-invalid+invalid",
		"-invalid.01",
		"alpha",
		"alpha.beta",
		"alpha.beta.1",
		"alpha.1",
		"alpha+beta",
		"alpha_beta",
		"alpha.",
		"alpha..",
		"beta",
		"1.0.0-alpha_beta",
		"-alpha.",
		"1.0.0-alpha..",
		"1.0.0-alpha..1",
		"1.0.0-alpha...1",
		"1.0.0-alpha....1",
		"1.0.0-alpha.....1",
		"1.0.0-alpha......1",
		"1.0.0-alpha.......1",
		"01.1.1",
		"1.01.1",
		"1.1.01",
		"1.2",
		"1.2.3.DEV",
		"1.2-SNAPSHOT",
		"1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788",
		"1.2-RC-SNAPSHOT",
		"-1.0.3-gamma+b7718",
		"+justmeta",
		"9.8.7+meta+meta",
		"9.8.7-whatever+meta+meta",
		"99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12",
		"1.0.0-0A..is..illegal",
		"1.0.0-a.01",
		"1.0.0-a.010",
	}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			_, err := Parse(test)
			if err == nil {
				t.Error("should error")
				return
			}
			if !errors.Is(err, ErrInvalid) {
				t.Errorf("unexpected error = %s", err)
				return
			}
		})
	}
}

func TestParseOverflow(t *testing.T) {
	var test = fmt.Sprintf("0.0.%s", strconv.Itoa(math.MaxInt)+"0")
	_, err := Parse(test)
	if err == nil {
		t.Error("should error")
		return
	}
	if !errors.Is(err, strconv.ErrRange) {
		t.Errorf("unexpected error = %s", err)
		return
	}
}

func TestCountPreReleaseData(t *testing.T) {
	tests := []struct {
		v Version
		l int
	}{
		{v: MustParse("0.1.0-A"), l: 1},
		{v: MustParse("0.1.0-A-B-C"), l: 1},
		{v: MustParse("0.1.0-A.B.C"), l: 3},
	}
	for _, tt := range tests {
		l := len(tt.v.PreRelease)
		if l != tt.l {
			t.Errorf("unexpected data count:\nexpected = %v\nactual   = %v", tt.l, l)
		}
	}
}

func TestCountBuildData(t *testing.T) {
	tests := []struct {
		v Version
		l int
	}{
		{v: MustParse("0.1.0+A"), l: 1},
		{v: MustParse("0.1.0+A-B"), l: 1},
		{v: MustParse("0.1.0+A.B.C"), l: 3},
	}
	for _, tt := range tests {
		l := len(tt.v.Build)
		if l != tt.l {
			t.Errorf("unexpected data count:\nexpected = %v\nactual   = %v", tt.l, l)
		}
	}
}

func TestVersion_Equals(t *testing.T) {
	tests := []struct {
		a      Version
		b      Version
		equals bool
	}{
		{
			a: Version{
				Major: 0,
				Minor: 1,
				Patch: 0,
			},
			b: Version{
				Major: 0,
				Minor: 1,
				Patch: 0,
			},
			equals: true,
		},
		{
			a: Version{
				Major: 0,
				Minor: 1,
				Patch: 1,
			},
			b: Version{
				Major: 0,
				Minor: 1,
				Patch: 0,
			},
			equals: false,
		},
	}
	for _, test := range tests {
		t.Run(test.a.String(), func(t *testing.T) {
			if test.a.Same(test.b) != test.equals {
				t.Error("should equal")
				return
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		expected int
	}{
		{
			a:        "2.0.0",
			b:        "1.0.0",
			expected: +1,
		},
		{
			a:        "1.0.0",
			b:        "2.0.0",
			expected: -1,
		},
		{
			a:        "0.1.1",
			b:        "0.1.0",
			expected: +1,
		},
		{
			a:        "0.1.0",
			b:        "0.1.1",
			expected: -1,
		},
		{
			a:        "1.1.1",
			b:        "1.1.1",
			expected: 0,
		},
		{
			a:        "1.0.0",
			b:        "1.0.0-alpha",
			expected: +1,
		},
		{
			a:        "1.0.0-alpha",
			b:        "1.0.0",
			expected: -1,
		},
		{
			a:        "1.0.0-a",
			b:        "1.0.0-b",
			expected: -1,
		},
		{
			a:        "1.0.0-c",
			b:        "1.0.0-c",
			expected: 0,
		},
		{
			a:        "1.0.0-a.b",
			b:        "1.0.0-a.b.c",
			expected: -1,
		},
		{
			a:        "1.0.0-a.b.c",
			b:        "1.0.0-a.b",
			expected: +1,
		},
		{
			a:        "1.0.0",
			b:        "1.0.0",
			expected: 0,
		},
		{
			a:        "1.0.0-a.b.c+foo",
			b:        "1.0.0-a.b",
			expected: +1,
		},
		{
			a:        "1.0.0-a.b.c",
			b:        "1.0.0-a.b+foo.bar",
			expected: +1,
		},
		{
			a:        "1.0.0-alpha",
			b:        "1.0.0-alpha.1",
			expected: -1,
		},
		{
			a:        "1.0.0-alpha.1",
			b:        "1.0.0-alpha.beta",
			expected: -1,
		},
		{
			a:        "1.0.0-alpha.beta",
			b:        "1.0.0-beta",
			expected: -1,
		},
		{
			a:        "1.0.0-beta",
			b:        "1.0.0-beta.2",
			expected: -1,
		},
		{
			a:        "1.0.0-beta.2",
			b:        "1.0.0-beta.11",
			expected: -1,
		},
		{
			a:        "1.0.0-beta.11",
			b:        "1.0.0-rc.1",
			expected: -1,
		},
		{
			a:        "1.0.0-rc.1",
			b:        "1.0.0",
			expected: -1,
		},
		{
			a:        "1.0.0-a.9999999.B",
			b:        "1.0.0-a.A",
			expected: -1,
		},
		{
			a:        "420.69.1337-a.99999999999999999999999",
			b:        "420.69.1337-a.888888888888888888888888",
			expected: -1,
		},
		{
			a:        "420.69.1337-a.99999999999999999999998",
			b:        "420.69.1337-a.99999999999999999999999",
			expected: -1,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s %s", test.a, test.b), func(t *testing.T) {
			a, err := Parse(test.a)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			b, err := Parse(test.b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			result := Compare(a, b)
			if test.expected != result {
				t.Errorf("unexpected result:\nexpected = %+v\nactual   = %+v", test.expected, result)
				return
			}
			switch result {
			case +1:
				if !a.Newer(b) {
					t.Error("util wrapper Newer() is broken")
				}
			case -1:
				if !a.Older(b) {
					t.Error("util wrapper Older() is broken")
				}
			case 0:
				if !a.Same(b) {
					t.Error("util wrapper Same() is broken")
				}
			}
		})
	}
}

func TestVersion_IsRelease(t *testing.T) {
	tests := []struct {
		name   string
		preRel []string
		want   bool
	}{
		{
			name:   "pre release (snapshot)",
			preRel: []string{"SNAPSHOT"},
			want:   false,
		},
		{
			name:   "pre release (empty string meta)",
			preRel: []string{""},
			want:   false,
		},
		{
			name:   "full release (no meta)",
			preRel: []string{},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				PreRelease: tt.preRel,
			}
			if got := v.IsRelease(); got != tt.want {
				t.Errorf("IsRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortAsc(t *testing.T) {
	var input = MustParseAll(
		[]string{
			"0.1.0",
			"1.0.0",
			"0.0.2-prerelease+meta",
			"0.0.1",
			"0.0.2-prerelease",
			"1.0.0-a.b.c+foo",
			"0.0.2+meta",
		})

	var expected = MustParseAll(
		[]string{
			"0.0.1",
			"0.0.2-prerelease+meta",
			"0.0.2-prerelease",
			"0.0.2+meta",
			"0.1.0",
			"1.0.0-a.b.c+foo",
			"1.0.0",
		})

	var sorted = SortAsc(input)
	if len(sorted) != len(expected) {
		t.Error("len not equal")
	}
	for i := range sorted {
		if sorted[i].String() != (expected[i]).String() {
			t.Error("versions not equal")
		}
	}
}

func TestSortDesc(t *testing.T) {
	var input = MustParseAll(
		[]string{
			"0.1.0",
			"1.0.0",
			"0.0.2-prerelease+meta",
			"0.0.1",
			"0.0.2-prerelease",
			"1.0.0-a.b.c+foo",
			"0.0.2+meta",
		})

	var expected = MustParseAll(
		[]string{
			"1.0.0",
			"1.0.0-a.b.c+foo",
			"0.1.0",
			"0.0.2+meta",
			"0.0.2-prerelease+meta",
			"0.0.2-prerelease",
			"0.0.1",
		})

	var sorted = SortDesc(input)
	if len(sorted) != len(expected) {
		t.Error("len not equal")
	}
	for i := range sorted {
		if sorted[i].String() != (expected[i]).String() {
			t.Error("versions not equal")
		}
	}
}
