// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package agent

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sylabs/fuzzball-agent/internal/pkg/cache"
)

type job struct {
	ID      string
	Name    string
	Image   string
	Command []string
	Volumes []volumeRequirement
	Cached  bool
	Hash    string
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

	// Create stream for output.
	s := stream{j.ID, a.nc}

	// Run Job.
	rc, err := a.runJob(context.TODO(), *j, s) // TODO: use context for cancellation?
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
func (a *Agent) runJob(ctx context.Context, j job, s stream) (int, error) {
	// Locate Singularity in PATH.
	path, err := exec.LookPath("singularity")
	if err != nil {
		return 0, err
	}

	// Generate bind path args for volumes
	var bindPaths []string
	for _, v := range j.Volumes {
		h, err := a.vm.GetHandle(v.VolumeID)
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

	image := j.Image
	if j.Cached {
		// Lookup image in cache and ensure it exists
		entry := a.c.GetEntry(cache.SIFType, j.Hash)
		if !entry.Exists() {
			return 0, fmt.Errorf("expected cached image does not exist in cache")
		}
		image = entry.Path()
	}

	args = append(args, image)
	args = append(args, j.Command...)

	// Run Singularity.
	state, err := runCommand(ctx, path, args, []string{}, "", nil, s, s)
	if err != nil {
		if state != nil {
			return state.ExitCode(), err
		}
		return 0, err
	}
	return state.ExitCode(), nil
}
