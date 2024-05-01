package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nothub/semver"
)

const usage string = `semver - utilities for semantic versioning

Options:
    -h, -help    Show help

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
	flag.Usage = func() {
		log.Print(usage)
	}
	flag.Parse()

	switch strings.ToLower(flag.Arg(0)) {
	case "next":
		next()
	case "strip":
		strip()
	case "valid":
		valid()
	case "help":
		fmt.Print(usage)
		os.Exit(0)
	default:
		log.Print(usage)
		os.Exit(1)
	}
}

func next() {
	if len(flag.Args()) < 3 {
		log.Print(usage)
		os.Exit(1)
	}

	ver, err := semver.Parse(flag.Arg(2))
	if err != nil {
		log.Fatalln(err)
	}

	switch strings.ToLower(flag.Arg(1)) {
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

func strip() {
	if len(flag.Args()) < 3 {
		log.Print(usage)
		os.Exit(1)
	}

	ver, err := semver.Parse(flag.Arg(2))
	if err != nil {
		log.Fatalln(err)
	}

	switch strings.ToLower(flag.Arg(1)) {
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

func valid() {
	if len(flag.Args()) < 2 {
		log.Print(usage)
		os.Exit(1)
	}

	_, err := semver.Parse(flag.Arg(1))
	if err != nil {
		log.Fatalln(err)
	}
}