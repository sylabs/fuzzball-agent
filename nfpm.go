// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build mage

package main

import (
	"os"

	"github.com/goreleaser/nfpm"
	_ "github.com/goreleaser/nfpm/deb"
	_ "github.com/goreleaser/nfpm/rpm"
)

// makePackage creates a package based on the supplied suffix.
func makePackage(suffix string) error {
	config, err := nfpm.ParseFile("nfpm.yaml")
	if err != nil {
		return err
	}
	config.Version = version()

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
