// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package agent

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
)

type job struct {
	ID      string
	Name    string
	Image   string
	Command []string
}

func (a *Agent) jobStartHandler(subject, reply string, j *job) {
	log := logrus.WithFields(logrus.Fields{
		"subject": subject,
		"reply":   reply,
		"jobID":   j.ID,
	})
	log.Print("handling job start")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("handled job start")
	}(time.Now())

	// Send acknowledgement.
	if err := a.ec.Publish(reply, nil); err != nil {
		log.WithError(err).Warn("failed to acknowledge job start")
	}

	// Run Job.
	rc, err := runJob(context.TODO(), *j) // TODO: use context for cancellation?

	// Send result.
	status := "COMPLETED"
	if err != nil {
		status = "FAILED"
	}
	res := struct {
		Status string
		RC     int
	}{status, rc}
	if err := a.ec.Publish(fmt.Sprintf("job.%v.finished", j.ID), res); err != nil {
		log.WithError(err).Warn("failed to report job finished")
	}
}

// runJob runs the specified job, returning the process exitCode.
func runJob(ctx context.Context, j job) (int, error) {
	// Locate Singularity in PATH.
	path, err := exec.LookPath("singularity")
	if err != nil {
		return 0, err
	}

	// Build up arguments to Singularity.
	args := []string{
		"exec",
		j.Image,
	}
	args = append(args, j.Command...)

	// Run Singularity.
	s, err := runCommand(ctx, path, args, []string{}, "", nil, nil, nil)
	if err != nil {
		if s != nil {
			return s.ExitCode(), err
		}
		return 0, err
	}
	return s.ExitCode(), nil
}
