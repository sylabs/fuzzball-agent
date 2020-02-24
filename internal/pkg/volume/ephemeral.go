// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package volume

import (
	"io/ioutil"
	"os"
)

type ephemeralDriver struct {
	baseDir string
}

func (ed ephemeralDriver) new() handler {
	return &ephemeral{baseDir: ed.baseDir}
}

// ephemeral represents a short lived volume that does not retain data.
type ephemeral struct {
	baseDir string
	path    string
}

func (e *ephemeral) create(id string) (err error) {
	e.path, err = ioutil.TempDir(e.baseDir, id)
	return err
}

func (e ephemeral) delete() error {
	if err := os.RemoveAll(e.path); err != nil {
		return err
	}

	return nil
}

func (e ephemeral) handle() string {
	return e.path
}
