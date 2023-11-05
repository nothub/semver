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
func Compare(a Version, b Version) int {
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

	if a.IsRelease() && !b.IsRelease() {
		return +1
	}
	if b.IsRelease() && !a.IsRelease() {
		return -1
	}

	for i := 0; i < max(len(a.PreRelease), len(b.PreRelease)); i++ {
		if i < len(a.PreRelease) && i >= len(b.PreRelease) {
			return +1
		}
		if i < len(b.PreRelease) && i >= len(a.PreRelease) {
			return -1
		}
		result := strings.Compare(a.PreRelease[i], b.PreRelease[i])
		if result != 0 {
			return result
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
