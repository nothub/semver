package semver

import (
	"errors"
	"regexp"
	"strings"
)

type Version struct {
	// incompatible API changes
	Major string
	// backward compatible functionality
	Minor string
	// backward compatible bug fixes
	Patch string

	// pre-release metadata
	PreRelease []string
	// build metadata
	Build []string
}

var ErrInvalid = errors.New("invalid semver string")

// https://regex101.com/r/vkijKf/1/
var re = regexp.MustCompile("^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$")

func Parse(str string) (Version, error) {
	var ver Version

	m := re.FindStringSubmatch(str)
	if m == nil {
		return ver, ErrInvalid
	}

	ver.Major = m[1]
	ver.Minor = m[2]
	ver.Patch = m[3]

	if len(m[4]) > 0 {
		ver.PreRelease = strings.Split(m[4], ".")
	}
	if len(m[5]) > 0 {
		ver.Build = strings.Split(m[5], ".")
	}

	return ver, nil
}

// Newer returns true if a is newer than b.
// Build metadata is ignored in this comparison.
func (a Version) Newer(b Version) bool {
	if a.Major != b.Major {
		return false
	}
	if a.Minor != b.Minor {
		return false
	}
	if a.Patch != b.Patch {
		return false
	}

	strings.Compare(strings.Join(a.PreRelease, ""), strings.Join(b.PreRelease, ""))

	return true
}

// Older returns true if a is older than b.
// Build metadata is ignored in this comparison.
func (a Version) Older(b Version) bool {
	if a.Major < b.Major {
		return true
	}
	if a.Minor < b.Minor {
		return true
	}
	if a.Patch < b.Patch {
		return true
	}

	// TODO:
	//   Precedence for two pre-release versions with the same major, minor, and patch version MUST be determined by comparing each dot separated identifier from left to right until a difference is found as follows:
	//    1. Identifiers consisting of only digits are compared numerically.
	//    2. Identifiers with letters or hyphens are compared lexically in ASCII sort order.
	//    3. Numeric identifiers always have lower precedence than non-numeric identifiers.
	//    4. A larger set of pre-release fields has a higher precedence than a smaller set, if all of the preceding identifiers are equal.

	return false
}

// Same returns true if a and b are equal.
// Build metadata is ignored in this comparison.
func (a Version) Same(b Version) bool {
	if a.Major != b.Major {
		return false
	}
	if a.Minor != b.Minor {
		return false
	}
	if a.Patch != b.Patch {
		return false
	}

	if len(a.PreRelease) != len(b.PreRelease) {
		return false
	}
	for i := range a.PreRelease {
		if a.PreRelease[i] != b.PreRelease[i] {
			return false
		}
	}

	return true
}

// Compare returns an integer comparing two Version objects.
//
//	a is same as b:  0
//	a older than b: -1
//	a newer than b: +1
func Compare(a Version, b Version) int {
	if a.Older(b) {
		return -1
	}

	if a.Newer(b) {
		return 1
	}

	// implies: a.Same(b) == true
	return 0
}

func (v Version) String() string {
	sb := strings.Builder{}
	sb.WriteString(v.Major)
	sb.WriteString(".")
	sb.WriteString(v.Minor)
	sb.WriteString(".")
	sb.WriteString(v.Patch)
	if len(v.PreRelease) > 0 {
		sb.WriteString("-")
		sb.WriteString(strings.Join(v.PreRelease, "."))
	}
	if len(v.Build) > 0 {
		sb.WriteString("+")
		sb.WriteString(strings.Join(v.Build, "."))
	}
	return sb.String()
}
