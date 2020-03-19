// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	cacheDirName = "fuzzball"
	SIFType      = "sif"
)

var (
	defaultCacheList = []string{
		SIFType,
	}
)

type Config struct {
	CacheDir string `yaml:"cacheDir"`
}

// Entry simply represents a filesystem location to store data.
type Entry struct {
	path string
}

type Cache struct {
	// baseDir is one level below the directory specified in the cache configuration
	// to ensure that cache operations are occuring on a directory only the agent controls
	baseDir string
}

func New(c Config) (*Cache, error) {
	var cache Cache
	cache.baseDir = filepath.Join(c.CacheDir, cacheDirName)
	if err := ensureDir(cache.baseDir); err != nil {
		return nil, err
	}

	for _, t := range defaultCacheList {
		if err := ensureDir(cache.cachePath(t)); err != nil {
			return nil, err
		}
	}

	return &cache, nil
}

// ensureDir check if a directory exists, if not it will create it and it's parents
func ensureDir(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("couldn't create cache directory %v: %v", dir, err)
		}
	} else if err != nil {
		return err
	}

	// Re-stat path to ensure we are working with a directory
	fi, err := os.Stat(dir)
	if !fi.IsDir() {
		return fmt.Errorf("path %s is not a directory", dir)
	}
	return nil
}

func (c *Cache) cachePath(cacheType string) string {
	return filepath.Join(c.baseDir, cacheType)
}

func (c *Cache) entryPath(cacheType, hash string) string {
	return filepath.Join(c.cachePath(cacheType), hash)
}

func (e *Entry) Exists() bool {
	_, err := os.Stat(e.path)
	return !os.IsNotExist(err)
}

func (e *Entry) Path() string {
	return e.path
}

func (c *Cache) GetEntry(cacheType, hash string) *Entry {
	return &Entry{c.entryPath(cacheType, hash)}
}
