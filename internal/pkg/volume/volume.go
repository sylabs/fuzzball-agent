// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package volume

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Manager struct {
	// maps volumeID to fs path
	v map[string]string
}

func NewManager() *Manager {
	return &Manager{v: make(map[string]string)}
}

func (m *Manager) Create(id, t string) error {
	if t != "EPHEMERAL" {
		return fmt.Errorf("unknown volume type: %s", t)
	}

	if _, ok := m.v[id]; ok {
		return fmt.Errorf("volume %s already exists", id)
	}

	handle, err := ioutil.TempDir("", id)
	if err != nil {
		return err
	}

	m.v[id] = handle
	return nil
}

func (m *Manager) Delete(id string) error {
	if _, ok := m.v[id]; !ok {
		return fmt.Errorf("volume %s does not exist", id)
	}

	if err := os.RemoveAll(m.v[id]); err != nil {
		return err
	}

	return nil
}

func (m *Manager) GetHandle(id string) (string, error) {
	if _, ok := m.v[id]; !ok {
		return "", fmt.Errorf("volume %s does not exist", id)
	}

	return m.v[id], nil
}
