// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package agent

import (
	"io"

	"github.com/sylabs/fuzzball-agent/internal/pkg/cache"
	vol "github.com/sylabs/fuzzball-agent/internal/pkg/volume"
	"gopkg.in/yaml.v3"
)

type rawConfig struct {
	NATSServers   []string     `yaml:"natsServers"`   // Array of nats server endpopints.
	VolumeSupport vol.Config   `yaml:"volumeSupport"` // List of available volume types.
	CacheConfig   cache.Config `yaml:"cacheConfig"`   // Description of fs location to store temporary data.
}

// NodeConfig represents a configuration.
type NodeConfig struct {
	raw rawConfig
}

// Read reads a config from the specified reader.
func Read(r io.Reader) (*NodeConfig, error) {
	var c NodeConfig
	if err := yaml.NewDecoder(r).Decode(&c.raw); err != nil {
		return nil, err
	}
	return &c, nil
}

func (nc *NodeConfig) SetNATSServers(uris []string) {
	nc.raw.NATSServers = uris
}

func (nc NodeConfig) NATSServers() []string {
	return nc.raw.NATSServers
}

func (nc *NodeConfig) SetVolumeConfig(vc vol.Config) {
	nc.raw.VolumeSupport = vc
}

func (nc NodeConfig) VolumeConfig() vol.Config {
	return nc.raw.VolumeSupport
}

func (nc *NodeConfig) SetCacheConfig(cc cache.Config) {
	nc.raw.CacheConfig = cc
}

func (nc NodeConfig) CacheConfig() cache.Config {
	return nc.raw.CacheConfig
}
