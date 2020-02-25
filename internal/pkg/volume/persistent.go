// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package volume

type persistentDriver struct {
	path string
}

func (pd persistentDriver) new() handler {
	return &persistent{path: pd.path}
}

type persistent struct {
	path string
}

func (persistent) create(id string) error {
	return nil
}

func (persistent) delete() error {
	return nil
}

func (p persistent) handle() string {
	return p.path
}
