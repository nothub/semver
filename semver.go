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

func (v Version) IsRelease() bool {
	return len(v.PreRelease) == 0
}

// Newer returns true if a is newer than b.
// Build metadata is ignored in this comparison.
func (a Version) Newer(b Version) bool {
	return Compare(a, b) == +1
}

// Older returns true if a is older than b.
// Build metadata is ignored in this comparison.
func (a Version) Older(b Version) bool {
	return Compare(a, b) == -1
}

// Same returns true if a and b are equal.
// Build metadata is ignored in this comparison.
func (a Version) Same(b Version) bool {
	return Compare(a, b) == 0
}

// Compare returns an integer comparing two Version objects.
//
//	a older than b: -1
//	a newer than b: +1
//	a is same as b:  0
//
// Build metadata is ignored in this comparison.
func Compare(a Version, b Version) int {
	// compare version core
	for _, result := range []int{
		strings.Compare(a.Major, b.Major),
		strings.Compare(a.Minor, b.Minor),
		strings.Compare(a.Patch, b.Patch),
	} {
		switch result {
		case +1:
			return +1
		case -1:
			return -1
		}
	}

	// release versions have precedence
	if a.IsRelease() && !b.IsRelease() {
		return +1
	}
	if b.IsRelease() && !a.IsRelease() {
		return -1
	}

	for i := 0; i < max(len(a.PreRelease), len(b.PreRelease)); i++ {
		// a larger set of pre-release data has a higher precedence
		// (if all the preceding identifiers are equal)
		if i < len(a.PreRelease) && i >= len(b.PreRelease) {
			return +1
		}
		if i < len(b.PreRelease) && i >= len(a.PreRelease) {
			return -1
		}

		if isDigits(a.PreRelease[i]) && isDigits(b.PreRelease[i]) {
			// identifiers consisting only of digits are compared numerically

			// if a digit string is longer, it has precedence
			if len(a.PreRelease[i]) > len(b.PreRelease[i]) {
				return +1
			}
			if len(a.PreRelease[i]) < len(b.PreRelease[i]) {
				return -1
			}

			// for same length digit strings, compare strings digit by digit
			for j := range a.PreRelease[i] {
				result := strings.Compare(a.PreRelease[i][j:j+1], b.PreRelease[i][j:j+1])
				if result != 0 {
					return result
				}
			}
		} else {
			// identifiers with letters or hyphens are compared lexically in ASCII sort order
			result := strings.Compare(a.PreRelease[i], b.PreRelease[i])
			if result != 0 {
				return result
			}
		}
	}

	// versions are equal
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

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
