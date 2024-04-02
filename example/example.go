//go:build exclude

package example

import (
	"fmt"
	"git.sr.ht/~hub/semver"
)

func main() {

	// parse semver from string
	v1, _ := semver.Parse("1.2.3")

	// or use Version struct directly
	v2 := semver.Version{Major: "2", Minor: "0", Patch: "0"}
	v3 := semver.Version{Major: "3", Minor: "0", Patch: "0"}
	v3pr := semver.Version{Major: "3", Minor: "0", Patch: "0", PreRelease: []string{"alpha"}}
	v3b := semver.Version{Major: "3", Minor: "0", Patch: "0", Build: []string{"foobar"}}

	fmt.Println("\ncomparing by version core:")
	fmt.Printf("is %q older then %q? %v\n", v1, v2, v1.Older(v2))
	fmt.Printf("compare %q and %q -> %v\n", v1, v2, semver.Compare(v1, v2))
	fmt.Printf("is %q older then %q? %v\n", v2, v1, v2.Older(v1))
	fmt.Printf("is %q newer then %q? %v\n", v2, v1, v2.Newer(v1))
	fmt.Printf("compare %q and %q -> %v\n", v2, v1, semver.Compare(v2, v1))
	fmt.Printf("is %q newer then %q? %v\n", v1, v1, v1.Newer(v1))
	fmt.Printf("compare %q and %q -> %v\n", v1, v1, semver.Compare(v1, v1))

	fmt.Println("\ncomparing by pre-release data:")
	fmt.Printf("is %q newer then %q? %v\n", v3pr, v3, v3pr.Newer(v3))
	fmt.Printf("is %q older then %q? %v\n", v3pr, v3, v3pr.Older(v3))
	fmt.Printf("compare %q and %q -> %v\n", v3pr, v3, semver.Compare(v3pr, v3))

	fmt.Println("\ncomparison ignores build data:")
	fmt.Printf("is %q newer then %q? %v\n", v3, v3b, v3.Newer(v3b))
	fmt.Printf("is %q older then %q? %v\n", v3, v3b, v3.Older(v3b))
	fmt.Printf("compare %q and %q -> %v\n", v3, v3b, semver.Compare(v3, v3b))
}
