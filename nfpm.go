// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/goreleaser/nfpm"
	_ "github.com/goreleaser/nfpm/deb"
	_ "github.com/goreleaser/nfpm/rpm"
)

// getPackageInfo returns the target based on suffix and c.
func getPackageInfo(c nfpm.Config, format string) (*nfpm.Info, error) {
	d, err := describeHead()
	if err != nil {
		return nil, err
	}

	v, err := getVersion(d)
	if err != nil {
		return nil, err
	}
	c.Version = v.String()

	info, err := c.Get(format)
	if err != nil {
		return nil, err
	}
	info = nfpm.WithDefaults(info)

	switch format {
	case "deb":
		// Ref: https://www.debian.org/doc/manuals/debian-faq/ch-pkg_basics.en.html#s-pkgname
		info.Target = fmt.Sprintf("%v_%v-%v_%v.%v",
			info.Name,
			info.Version,
			info.Release,
			info.Arch,
			format)
	case "rpm":
		// Ref: http://ftp.rpm.org/max-rpm/ch-rpm-file-format.html
		info.Target = fmt.Sprintf("%v-%v-%v.%v.%v",
			info.Name,
			info.Version,
			info.Release,
			info.Arch,
			format)
	default:
		return nil, fmt.Errorf("unknown package format: %v", format)
	}

	if err = nfpm.Validate(info); err != nil {
		return nil, err
	}
	return info, nil
}

// makePackage creates a package based on the supplied format.
func makePackage(format string) error {
	config, err := nfpm.ParseFile("nfpm.yaml")
	if err != nil {
		return err
	}

	info, err := getPackageInfo(config, format)
	if err != nil {
		return err
	}

	p, err := nfpm.Get(format)
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
