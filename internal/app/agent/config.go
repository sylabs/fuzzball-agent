// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package agent

import (
	"io"

	vol "github.com/sylabs/fuzzball-agent/internal/pkg/volume"
	"gopkg.in/yaml.v3"
)

type rawConfig struct {
	VolumeSupport vol.Config `yaml:"volumeSupport"` // List of available volume types.
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

func (nc *NodeConfig) SetVolumeConfig(vc vol.Config) {
	nc.raw.VolumeSupport = vc
}

func (nc NodeConfig) VolumeConfig() vol.Config {
	return nc.raw.VolumeSupport
}
