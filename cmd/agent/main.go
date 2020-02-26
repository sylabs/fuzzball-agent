// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/sylabs/compute-agent/internal/app/agent"
	"github.com/sylabs/compute-agent/internal/pkg/volume"
)

const (
	org  = "Sylabs"
	name = "Compute Agent"
)

var (
	version = "unknown"

	natsURIs   = flag.String("nats_uris", nats.DefaultURL, "Comma-separated list of NATS server URIs")
	configPath = flag.String("config_path", "/etc/fuzzball/config.yaml", "Path to agent configuration on node")
)

// signalHandler catches SIGINT/SIGTERM to perform an orderly shutdown.
func signalHandler(a agent.Agent) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	logrus.WithFields(logrus.Fields{
		"signal": (<-c).String(),
	}).Info("shutting down due to signal")

	a.Stop()
}

// defaultNodeConfig will return a configuration where only ephemeral volumes are enabled
// and will be located at the temporary directory of the system.
func defaultNodeConfig() *agent.NodeConfig {
	var nc agent.NodeConfig
	vc := volume.Config{
		volume.TypeEphemeral: volume.Spec{
			Location: os.TempDir(),
		},
	}
	nc.SetVolumeConfig(vc)
	return &nc
}

// parseNodeConfig parses the node configuration at the specified path.
func parseNodeConfig() (*agent.NodeConfig, error) {
	// Parse node configuration.
	f, err := os.Open(*configPath)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.WithFields(logrus.Fields{
				"config": *configPath,
			}).Warnf("node config not found, using default")
			return defaultNodeConfig(), nil
		}
		return nil, err
	}
	defer f.Close()

	nodeConfig, err := agent.Read(f)
	if err != nil {
		return nil, err
	}
	return nodeConfig, nil
}

func main() {
	flag.Parse()

	log := logrus.WithFields(logrus.Fields{
		"org":     org,
		"name":    name,
		"version": version,
	})
	log.Info("starting")
	defer log.Info("stopped")

	nodeConfig, err := parseNodeConfig()
	if err != nil {
		logrus.WithError(err).Error("failed to parse node configuration")
		return
	}
	// Spin up agent.
	c := agent.Config{
		NodeConfig:  nodeConfig,
		NATSServers: strings.Split(*natsURIs, ","),
	}
	a, err := agent.New(c)
	if err != nil {
		logrus.WithError(err).Error("failed to create agent")
		return
	}

	// Spin off signal handler to do graceful shutdown.
	go signalHandler(a)

	// Main agent routine.
	a.Run()
}
