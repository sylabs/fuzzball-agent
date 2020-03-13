// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package agent

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

// connectionFields returns Fields representing nc suitable for use with WithFields.
func connectionFields(nc *nats.Conn) logrus.Fields {
	f := logrus.Fields{}
	if v := nc.ConnectedAddr(); v != "" {
		f["connectedAddr"] = v
	}
	if v := nc.ConnectedServerId(); v != "" {
		f["connectedServerID"] = v
	}
	if id, err := nc.GetClientID(); err == nil {
		f["clientID"] = id
	}
	return f
}

// connect opens a connection to NATS.
func connect(c Config) (nc *nats.Conn, ec *nats.EncodedConn, err error) {
	o := nats.GetDefaultOptions()
	o.Servers = c.NodeConfig.NATSServers()

	// Log disconnect/reconnect events.
	o.DisconnectedErrCB = func(c *nats.Conn, err error) {
		log := logrus.WithFields(connectionFields(c))
		if err != nil {
			log.WithError(err).Warn("messaging system disconnected")
		} else {
			log.Info("messaging system disconnected")
		}
	}
	o.ReconnectedCB = func(c *nats.Conn) {
		logrus.WithFields(connectionFields(c)).Info("messaging system reconnected")
	}

	// Connect to messaging system.
	logrus.WithFields(logrus.Fields{
		"natsServers": o.Servers,
	}).Info("connecting to messaging system")
	nc, err = o.Connect()
	if err != nil {
		return nil, nil, err
	}
	defer func(nc *nats.Conn) {
		if err != nil {
			nc.Close()
		}
	}(nc)
	logrus.WithFields(connectionFields(nc)).Print("messaging system connected")

	// Wrap connection with JSON encoder.
	if ec, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER); err != nil {
		return nil, nil, err
	}
	return nc, ec, nil
}

// subscribe expresses interest in subjects that are relevant to the Agent.
func (a Agent) subscribe() error {
	subs := []struct {
		subject string
		handler nats.Handler
	}{
		{fmt.Sprintf("node.%s.job.start", a.id), a.jobStartHandler},
		{fmt.Sprintf("node.%s.volume.create", a.id), a.volumeCreateHandler},
		{fmt.Sprintf("node.%s.volume.delete", a.id), a.volumeDeleteHandler},
	}
	for _, s := range subs {
		if _, err := a.ec.Subscribe(s.subject, s.handler); err != nil {
			logrus.WithField("subject", s.subject).WithError(err).Warn("failed to subscribe")
			return err
		}
		logrus.WithField("subject", s.subject).Info("subscribed")
	}
	return nil
}
