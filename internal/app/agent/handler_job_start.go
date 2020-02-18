// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package agent

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	vol "github.com/sylabs/compute-agent/internal/pkg/volume"
)

type job struct {
	ID      string
	Name    string
	Image   string
	Command []string
	Volumes []volumeRequirement
}

type volumeRequirement struct {
	VolumeID string
	Location string
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
	var out strings.Builder
	rc, err := runJob(context.TODO(), *j, a.vm, &out) // TODO: use context for cancellation?
	// Send result.
	status := "COMPLETED"
	if err != nil {
		status = "FAILED"
	}
	res := struct {
		Status string
		RC     int
		Out    string
	}{status, rc, out.String()}
	if err := a.ec.Publish(fmt.Sprintf("job.%v.finished", j.ID), res); err != nil {
		log.WithError(err).Warn("failed to report job finished")
	}
}

// runJob runs the specified job, returning the process exitCode.
func runJob(ctx context.Context, j job, vm *vol.Manager, output io.Writer) (int, error) {
	// Locate Singularity in PATH.
	path, err := exec.LookPath("singularity")
	if err != nil {
		return 0, err
	}

	// Generate bind path args for volumes
	var bindPaths []string
	for _, v := range j.Volumes {
		h, err := vm.GetHandle(v.VolumeID)
		if err != nil {
			return 0, err
		}

		bp := h + ":" + v.Location
		bindPaths = append(bindPaths, bp)
	}

	// Build up arguments to Singularity.
	args := []string{
		"exec",
	}

	// Bind paths are a comma separated list.
	if bindPaths != nil {
		args = append(args, "--bind", strings.Join(bindPaths, ","))
	}

	args = append(args, j.Image)
	args = append(args, j.Command...)

	// Run Singularity.

	s, err := runCommand(ctx, path, args, []string{}, "", nil, output, output)
	if err != nil {
		if s != nil {
			return s.ExitCode(), err
		}
		return 0, err
	}
	return s.ExitCode(), nil
}
