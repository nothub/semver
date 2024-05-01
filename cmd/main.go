package main

import (
	"errors"
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

var errUsage = errors.New("invalid usage")

func main() {
	log.SetFlags(0)

	args := os.Args[1:]

	mustLen(args, 1)
	if args[0] == "-h" || args[0] == "--help" {
		fmt.Print(usage)
		os.Exit(0)
	}

	cmd := strings.ToLower(args[0])
	args = args[1:]

	var out string
	var err error

	switch cmd {
	case "next":
		mustLen(args, 2)
		out, err = next(args[0], args[1])
	case "strip":
		mustLen(args, 2)
		out, err = strip(args[0], args[1])
	case "valid":
		mustLen(args, 1)
		err = valid(args[0])
	default:
		err = errUsage
	}

	if err != nil {
		if errors.Is(err, errUsage) {
			log.Print(usage)
			os.Exit(1)
		} else {
			log.Fatalln(err.Error())
		}
	}

	fmt.Print(out)
}

func mustLen(args []string, minLen int) {
	if len(args) < minLen {
		log.Print(usage)
		os.Exit(1)
	}
}

func next(mode string, str string) (string, error) {
	ver, err := semver.Parse(str)
	if err != nil {
		return "", err
	}

	switch strings.ToLower(mode) {
	case "major":
		ver.Major = ver.Major + 1
	case "minor":
		ver.Minor = ver.Minor + 1
	case "patch":
		ver.Patch = ver.Patch + 1
	default:
		return "", errUsage
	}

	return ver.String(), nil
}

func strip(mode string, str string) (string, error) {
	ver, err := semver.Parse(str)
	if err != nil {
		return "", err
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
		return "", errUsage
	}

	return ver.String(), nil
}

func valid(str string) error {
	_, err := semver.Parse(str)
	return err
}
