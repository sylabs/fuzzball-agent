// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

//+build mage

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/sirupsen/logrus"
)

// ldFlags returns standard linker flags to pass to various Go commands.
func ldFlags() string {
	vals := []string{
		fmt.Sprintf("-X main.builtAt=%v", time.Now().UTC().Format(time.RFC3339)),
	}

	// Attempt to get git details.
	d, err := describeHead()
	if err == nil {
		vals = append(vals, fmt.Sprintf("-X main.gitCommit=%v", d.ref.Hash().String()))

		if d.isClean {
			vals = append(vals, "-X main.gitTreeState=clean")
		} else {
			vals = append(vals, "-X main.gitTreeState=dirty")
		}

		if v, err := getVersion(d); err != nil {
			logrus.WithError(err).Warn("failed to get version from git description")
		} else {
			vals = append(vals, fmt.Sprintf("-X main.gitVersion=%v", v.String()))
		}
	}

	return strings.Join(vals, " ")
}

// Build creates a binary in the current directory.
func Build() error {
	return sh.RunV(mg.GoCmd(), "build", "-ldflags", ldFlags(), "./cmd/fuzzball-agent")
}

// Install installs the agent binary in $GOBIN.
func Install() error {
	return sh.RunV(mg.GoCmd(), "install", "-ldflags", ldFlags(), "./cmd/fuzzball-agent")
}

// Run runs the Agent using `go run`.
func Run() error {
	return sh.RunV(mg.GoCmd(), "run", "-ldflags", ldFlags(), "./cmd/fuzzball-agent/")
}

// Test runs unit tests using `go test`.
func Test() error {
	return sh.RunV(mg.GoCmd(), "test", "-ldflags", ldFlags(), "-cover", "-race", "./...")
}

// Deb builds a deb package.
func Deb() error {
	mg.Deps(Build)
	return makePackage("deb")
}

// RPM builds a RPM package.
func RPM() error {
	mg.Deps(Build)
	return makePackage("rpm")
}
