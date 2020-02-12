// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package agent

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type volume struct {
	ID   string
	Name string
	Type string
}

func (a *Agent) volumeCreateHandler(subject, reply string, v volume) {
	log := logrus.WithFields(logrus.Fields{
		"subject": subject,
		"reply":   reply,
	})
	log.Print("handling volume creation")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("handled volume creation")
	}(time.Now())

	// Send acknowledgement.
	if err := a.ec.Publish(reply, nil); err != nil {
		log.WithError(err).Warn("failed to acknowledge volume creation")
	}

	// Create volume.
	err := a.vm.Create(v.ID, v.Type) // TODO: use context for cancellation?

	// Send result.
	res := struct {
		err error
	}{err}
	if err := a.ec.Publish(fmt.Sprintf("volume.%v.create", v.ID), res); err != nil {
		log.WithError(err).Warn("failed to report volume creation")
	}
}

func (a *Agent) volumeDeleteHandler(subject, reply string, v volume) {
	log := logrus.WithFields(logrus.Fields{
		"subject": subject,
		"reply":   reply,
	})
	log.Print("handling volume deletion")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("handled volume deletion")
	}(time.Now())

	// Send acknowledgement.
	if err := a.ec.Publish(reply, nil); err != nil {
		log.WithError(err).Warn("failed to acknowledge volume deletion")
	}

	// Delete volume.
	err := a.vm.Delete(v.ID) // TODO: use context for cancellation?

	// Send result.
	res := struct {
		err error
	}{err}
	if err := a.ec.Publish(fmt.Sprintf("volume.%v.delete", v.ID), res); err != nil {
		log.WithError(err).Warn("failed to report volume deletion")
	}
}
