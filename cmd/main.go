package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nothub/semver"
)

const usage string = `semver - utilities for semantic versioning

Options:
    -h, --help   Show help

Commands:
    next - Bump to the next version
    Usage: semver [opts...] next (major|minor|patch) <version>

    strip - Remove pre-release or build metadata
    Usage: semver [opts...] strip (all|pre|build) <version>

    valid - Check input for conformity
    Usage: semver [opts...] valid <version>
`

func main() {
	log.SetFlags(0)

	checkArgs(1)
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Print(usage)
		os.Exit(0)
	}

	cmd := strings.ToLower(os.Args[1])
	switch strings.ToLower(cmd) {
	case "next":
		checkArgs(2)
		next(os.Args[2], os.Args[3])
	case "strip":
		checkArgs(2)
		strip(os.Args[2], os.Args[3])
	case "valid":
		checkArgs(1)
		valid(os.Args[2])
	default:
		log.Print(usage)
		os.Exit(1)
	}
}

func checkArgs(min int) {
	if len(os.Args) < (min + 1) {
		log.Print(usage)
		os.Exit(1)
	}
}

func next(mode string, str string) {
	ver, err := semver.Parse(str)
	if err != nil {
		log.Fatalln(err)
	}

	switch strings.ToLower(mode) {
	case "major":
		ver.Major = ver.Major + 1
	case "minor":
		ver.Minor = ver.Minor + 1
	case "patch":
		ver.Patch = ver.Patch + 1
	default:
		log.Print(usage)
		os.Exit(1)
	}

	fmt.Print(ver.String())
}

func strip(mode string, str string) {
	ver, err := semver.Parse(str)
	if err != nil {
		log.Fatalln(err)
	}

	switch strings.ToLower(mode) {
	case "all":
		ver.PreRelease = nil
		ver.Build = nil
	case "pre":
		ver.PreRelease = nil
	case "build":
		ver.Build = nil
	default:
		log.Print(usage)
		os.Exit(1)
	}

	fmt.Print(ver.String())
}

func valid(str string) {
	_, err := semver.Parse(str)
	if err != nil {
		log.Fatalln(err)
	}
}
