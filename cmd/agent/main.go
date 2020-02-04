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
)

const (
	org  = "Sylabs"
	name = "Compute Agent"
)

var (
	version = "unknown"

	natsURIs = flag.String("nats_uris", nats.DefaultURL, "Comma-separated list of NATS server URIs")
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

func main() {
	flag.Parse()

	log := logrus.WithFields(logrus.Fields{
		"org":     org,
		"name":    name,
		"version": version,
	})
	log.Info("starting")
	defer log.Info("stopped")

	// Spin up agent.
	c := agent.Config{
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
