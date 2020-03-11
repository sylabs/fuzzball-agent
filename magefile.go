// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

//+build mage

package main

import (
	"fmt"
	"os"

	"github.com/goreleaser/nfpm"
	_ "github.com/goreleaser/nfpm/deb"
	_ "github.com/goreleaser/nfpm/rpm"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build creates a binary in the current directory.
func Build() error {
	fmt.Println("Building Fuzzball agent")
	return sh.Run("go", "build", "./cmd/fuzzball-agent")
}

// Install installs the agent binary in $GOBIN.
func Install() error {
	fmt.Println("Installing Fuzzball agent")
	return sh.Run("go", "install", "./cmd/fuzzball-agent")
}

// pack creates a package based on the supplied suffix.
func pack(suffix string) error {
	mg.Deps(Build)

	fmt.Printf("Packaging Fuzzball agent as %s\n", suffix)
	config, err := nfpm.ParseFile("nfpm.yaml")
	if err != nil {
		return err
	}

	info, err := config.Get(suffix)
	if err != nil {
		return err
	}
	info.Target = "fuzzball-agent." + suffix

	info = nfpm.WithDefaults(info)

	if err = nfpm.Validate(info); err != nil {
		return err
	}

	p, err := nfpm.Get(suffix)
	if err != nil {
		return err
	}

	f, err := os.Create(info.Target)
	if err != nil {
		return err
	}
	defer f.Close()

	return p.Package(info, f)
}

// Deb builds a deb package.
func Deb() error {
	return pack("deb")
}

// RPM builds a RPM package.
func RPM() error {
	return pack("rpm")
}
