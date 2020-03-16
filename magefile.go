// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

//+build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build creates a binary in the current directory.
func Build() error {
	return sh.Run("go", "build", "./cmd/fuzzball-agent")
}

// Install installs the agent binary in $GOBIN.
func Install() error {
	return sh.Run("go", "install", "./cmd/fuzzball-agent")
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
