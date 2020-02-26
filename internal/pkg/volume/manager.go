// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package volume

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	// TypeEphemeral represents a short lived volume that does not retain data.
	TypeEphemeral = "EPHEMERAL"
	// TypePersistent represents a volume that leaves data intact during creation and removal.
	TypePersistent = "PERSISTENT"
)

// Config describes volume manager configuration.
type Config map[string]Spec

// Spec defines a local resource to use as a volume.
type Spec struct {
	Location string `yaml:"location"`
}

// driver is specific to a volume type and generates handlers
// for individual instances of a volume.
type driver interface {
	new() handler
}

// handler manages a single instance of a volume for a workflow.
type handler interface {
	create(string) error
	delete() error
	handle() string
}

// Manager creates and tracks volumes in use.
type Manager struct {
	m       sync.Mutex
	support map[string]driver
	volumes map[string]handler
}

// NewManager creates a new Manager based on the supplied volume configuration.
func NewManager(c Config) (*Manager, error) {
	var m Manager
	m.support = make(map[string]driver)
	m.volumes = make(map[string]handler)

	// read from config and register driver for different types.
	for t, v := range c {
		t = strings.ToUpper(t)
		var d driver
		switch t {
		case TypeEphemeral:
			d = &ephemeralDriver{baseDir: v.Location}
		case TypePersistent:
			d = &persistentDriver{path: v.Location}
		default:
			return nil, fmt.Errorf("unsupported volume type: %s", t)
		}
		m.support[t] = d
		logrus.WithFields(logrus.Fields{
			"driver":   t,
			"location": v.Location,
		}).Infof("registered volume driver")
	}

	return &m, nil
}

// Purge will call delete() on every volume and remove it from the manager.
// Any errors will be logged with logrus.
func (m *Manager) Purge() {
	m.m.Lock()
	defer m.m.Unlock()

	for id, vol := range m.volumes {
		log := logrus.WithFields(logrus.Fields{
			"volumeID": id,
		})

		log.Infof("deleting volume")
		err := vol.delete()
		if err != nil {
			log.WithError(err).Warn("failed to delete volume")
		}
		log.Infof("volume deleted")
	}
	return
}

// Create adds a volume to the manager and preforms any required
// setup based on the volume type.
func (m *Manager) Create(id, t string) error {
	h, err := m.create(id, t)
	if err != nil {
		return err
	}

	return h.create(id)
}

// create registers a volume handler with the manager in a thread safe manner.
func (m *Manager) create(id, t string) (handler, error) {
	m.m.Lock()
	defer m.m.Unlock()

	d, ok := m.support[t]
	if !ok {
		return nil, fmt.Errorf("unsupported volume type: %s", t)
	}

	if _, ok := m.volumes[id]; ok {
		return nil, fmt.Errorf("volume %s already exists", id)
	}

	h := d.new()
	m.volumes[id] = h
	return h, nil
}

// Delete removes the volume from the manager and cleans up
// the filesystem when required by the volume type.
func (m *Manager) Delete(id string) error {
	h, err := m.delete(id)
	if err != nil {
		return err
	}

	return h.delete()
}

// remove deletes a volume handler from the manager in a thread safe manner.
func (m *Manager) delete(id string) (handler, error) {
	m.m.Lock()
	defer m.m.Unlock()

	h, ok := m.volumes[id]
	if !ok {
		return nil, fmt.Errorf("volume %s does not exist", id)
	}
	delete(m.volumes, id)

	return h, nil
}

// GetHandle returns the filesystem location of the volume.
func (m *Manager) GetHandle(id string) (string, error) {
	h, err := m.getHandler(id)
	if err != nil {
		return "", err
	}

	return h.handle(), nil
}

// getHandler accesses a volume handler in the manager in a thread safe manner.
func (m *Manager) getHandler(id string) (handler, error) {
	m.m.Lock()
	defer m.m.Unlock()
	h, ok := m.volumes[id]
	if !ok {
		return nil, fmt.Errorf("volume %s does not exist", id)
	}

	return h, nil
}
