// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package agent

import (
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	vol "github.com/sylabs/fuzzball-agent/internal/pkg/volume"
)

// Config describes agent configuration.
type Config struct {
	NodeConfig  *NodeConfig
	NATSServers []string
}

// Agent contains the state of the agent.
type Agent struct {
	nc *nats.Conn
	ec *nats.EncodedConn
	vm *vol.Manager
	id string
}

// New returns a new Agent.
func New(c Config) (a Agent, err error) {
	a = Agent{
		id: "1", // TODO
	}

	if a.vm, err = vol.NewManager(c.NodeConfig.VolumeConfig()); err != nil {
		return Agent{}, err
	}

	if a.nc, a.ec, err = connect(c); err != nil {
		return Agent{}, err
	}

	return a, nil
}

// Run is the main routine for the Agent.
func (a Agent) Run() error {
	// Use a WaitGroup to wait for messaging connection to drain.
	// A closed connection indicates the agent has stopped
	// and system resources should be released.
	wg := sync.WaitGroup{}
	wg.Add(1)
	a.nc.SetClosedHandler(func(c *nats.Conn) {
		logrus.WithFields(connectionFields(c)).Print("messaging system connection closed")

		// Clean up volumes after connection has been closed.
		a.vm.Purge()

		wg.Done()
	})

	// Subscribe to relevant topics.
	if err := a.subscribe(); err != nil {
		logrus.WithError(err).Warn("failed to subscribe")
		return err
	}

	// Wait for messaging connection to close.
	wg.Wait()

	return nil
}

// Stop is used to gracefully stop the Agent.
func (a Agent) Stop() {
	if err := a.nc.Drain(); err == nats.ErrConnectionReconnecting {
		logrus.Info("forcefully closing messaging system connection")
		a.nc.Close()
	} else if err != nil {
		logrus.WithError(err).Warn("failed to drain")
	}
}
