package semver

import (
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Version struct {
	// incompatible API changes
	Major int
	// backward compatible functionality
	Minor int
	// backward compatible bug fixes
	Patch int

	// pre-release metadata
	PreRelease []string
	// build metadata
	Build []string
}

var ErrInvalid = errors.New("invalid semver string")

// https://regex101.com/r/vkijKf/1/
var regex = regexp.MustCompile("^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$")

// Parse will attempts to convert a string to a semver.Version struct.
//
// Parse might return semver.ErrInvalid, strconv.ErrRange or strconv.ErrSyntax.
func Parse(str string) (Version, error) {
	var ver Version

	m := regex.FindStringSubmatch(str)
	if m == nil {
		return ver, ErrInvalid
	}

	var err error
	ver.Major, err = strconv.Atoi(m[1])
	if err != nil {
		return ver, err
	}
	ver.Minor, err = strconv.Atoi(m[2])
	if err != nil {
		return ver, err
	}
	ver.Patch, err = strconv.Atoi(m[3])
	if err != nil {
		return ver, err
	}

	if len(m[4]) > 0 {
		ver.PreRelease = strings.Split(m[4], ".")
		ver.PreRelease = removeEmpty(ver.PreRelease)
	}
	if len(m[5]) > 0 {
		ver.Build = strings.Split(m[5], ".")
		ver.Build = removeEmpty(ver.Build)
	}

	return ver, nil
}

func MustParse(str string) Version {
	vers, err := Parse(str)
	if err != nil {
		panic(err)
	}
	return vers
}

func ParseAll(strs []string) ([]Version, error) {
	var vers []Version
	for _, str := range strs {
		ver, err := Parse(str)
		if err != nil {
			return nil, err
		}
		vers = append(vers, ver)
	}
	return vers, nil
}

func MustParseAll(strs []string) []Version {
	vers, err := ParseAll(strs)
	if err != nil {
		panic(err)
	}
	return vers
}

// IsRelease returns true if Version contains no pre-release metadata.
func (v *Version) IsRelease() bool {
	return len(v.PreRelease) == 0
}

// Newer returns true if a is newer than b.
// Build metadata is ignored in this comparison.
func (a *Version) Newer(b Version) bool {
	return Compare(*a, b) == +1
}

// Older returns true if a is older than b.
// Build metadata is ignored in this comparison.
func (a *Version) Older(b Version) bool {
	return Compare(*a, b) == -1
}

// Same returns true if a and b are equal.
// Build metadata is ignored in this comparison.
func (a *Version) Same(b Version) bool {
	return Compare(*a, b) == 0
}

// Compare returns an integer comparing two Version objects.
//
//	a, _ := Parse("1.2.3")
//	b, _ := Parse("2.0.0")
//	Compare(a, a) ->  0 (a older than b)
//	Compare(a, b) -> -1 (a newer than b)
//	Compare(b, a) -> +1 (a is same as b)
//
// Build metadata is ignored in this comparison.
func Compare(a Version, b Version) int {
	result := compareVersionCore([]int{
		compareInt(a.Major, b.Major),
		compareInt(a.Minor, b.Minor),
		compareInt(a.Patch, b.Patch),
	})
	if result != 0 {
		return result
	}
	return comparePreRelease(a, b)
}

func compareInt(a int, b int) int {
	if a > b {
		return +1
	}
	if a < b {
		return -1
	}
	return 0
}

func compareVersionCore(core []int) int {
	for _, result := range core {
		switch result {
		case +1:
			return +1
		case -1:
			return -1
		}
	}
	return 0
}

func comparePreRelease(a Version, b Version) int {
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
	return 0
}

// String will build and return the string representation of Version.
func (v *Version) String() string {
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(v.Major))
	sb.WriteString(".")
	sb.WriteString(strconv.Itoa(v.Minor))
	sb.WriteString(".")
	sb.WriteString(strconv.Itoa(v.Patch))
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

func removeEmpty(arr []string) []string {
	res := make([]string, 0, len(arr))
	for _, s := range arr {
		if strings.TrimSpace(s) != "" {
			res = append(res, s)
		}
	}
	return res
}

func SortAsc(vers []Version) []Version {
	sort.SliceStable(vers, func(i, j int) bool {
		return vers[i].Older(vers[j])
	})
	return vers
}

func SortDesc(vers []Version) []Version {
	sort.SliceStable(vers, func(i, j int) bool {
		return vers[i].Newer(vers[j])
	})
	return vers
}
