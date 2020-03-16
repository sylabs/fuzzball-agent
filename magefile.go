// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

//+build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// ldFlags returns standard linker flags to pass to various Go commands.
func ldFlags() string {
	return fmt.Sprintf("-X main.version=%s", version())
}

// Build creates a binary in the current directory.
func Build() error {
	return sh.Run("go", "build", "-ldflags", ldFlags(), "./cmd/fuzzball-agent")
}

// Install installs the agent binary in $GOBIN.
func Install() error {
	return sh.Run("go", "install", "-ldflags", ldFlags(), "./cmd/fuzzball-agent")
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
