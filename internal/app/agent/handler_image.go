// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package agent

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sylabs/fuzzball-agent/internal/pkg/cache"
	scs "github.com/sylabs/scs-library-client/client"
)

type image struct {
	URI string
}

func (a *Agent) imageDownloadHandler(subject, reply string, i image) {
	log := logrus.WithFields(logrus.Fields{
		"subject":  subject,
		"reply":    reply,
		"imageURI": i.URI,
	})
	log.Print("handling image download")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("handled image download")
	}(time.Now())

	// Send acknowledgement.
	if err := a.ec.Publish(reply, nil); err != nil {
		log.WithError(err).Warn("failed to acknowledge image download")
	}

	// Returned tag should be a unique identifier of the image
	err := a.libraryImageDownload(i.URI)
	if err != nil {
		log.WithError(err).Warnf("could not download library image: %v", err)
	}

	// Send result.
	res := struct {
		Err error
	}{err}
	if err := a.ec.Publish("image.download", res); err != nil {
		log.WithError(err).Warn("failed to report image download")
	}
}

func (a *Agent) libraryImageDownload(uri string) error {
	// Parse image uri
	r, err := scs.Parse(uri)
	if err != nil {
		return err
	}

	var tag string
	if len(r.Tags) > 0 {
		tag = r.Tags[0]
	}

	// Ensure tag guarantees for reproducable image pulls
	if !scs.IsImageHash(tag) {
		return fmt.Errorf("tag must be the image hash, received %v", tag)
	}

	// Point library client to specific library if included in uri
	var scsConf *scs.Config
	if r.Host != "" {
		scsConf = &scs.Config{BaseURL: "https://" + r.Host}
	}

	// Initialize library client
	client, err := scs.NewClient(scsConf)
	if err != nil {
		return err
	}

	// Get cache entry for image to be downloaded
	entry := a.c.GetEntry(cache.SIFType, tag)

	// Open file descriptor to download image to cache entry
	f, err := os.Create(entry.Path())
	if err != nil {
		return err
	}
	defer f.Close()

	// Download image from library.
	err = client.DownloadImage(context.Background(), f, runtime.GOARCH, r.Path, tag, nil)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) imageCachedHandler(subject, reply string, hash string) {
	log := logrus.WithFields(logrus.Fields{
		"subject": subject,
		"reply":   reply,
		"hash":    hash,
	})
	log.Print("handling image cache check")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("handled image cache check")
	}(time.Now())

	// Send acknowledgement.
	if err := a.ec.Publish(reply, nil); err != nil {
		log.WithError(err).Warn("failed to acknowledge image cache check")
	}

	// Get cache entry for image to be downloaded
	entry := a.c.GetEntry(cache.SIFType, hash)

	log.Infof("entry exists: %v", entry.Exists())
	// Send result.
	res := struct {
		Exists bool
	}{entry.Exists()}
	if err := a.ec.Publish("image.cached", res); err != nil {
		log.WithError(err).Warn("failed to report image cache check")
	}
}
