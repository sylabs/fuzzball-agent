// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package volume

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestPersistent(t *testing.T) {
	testID := "TestID"
	testContent := []byte("testfile")
	testFilename := "test-file"

	// create tmpdir as base for volume
	path, err := ioutil.TempDir("", "test-persistent-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)
	testPath := filepath.Join(path, testFilename)

	// write to a test file in the volume
	if err := ioutil.WriteFile(testPath, testContent, 0644); err != nil {
		t.Fatalf("failed to write to test file in volume")
	}

	p := persistent{
		path: path,
	}

	// NOTE: testID is ignored by the persistent file handler
	err = p.create(testID)
	if err != nil {
		t.Fatal(err)
	}

	// ensure test file appears in volume
	handlePath := p.handle()
	content, err := ioutil.ReadFile(filepath.Join(handlePath, testFilename))
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(content, testContent) != 0 {
		t.Fatalf("want %s, got %s", string(testContent), string(content))
	}

	err = p.delete()
	if err != nil {
		t.Fatal(err)
	}

	// ensure test file remains after volume deletion
	content, err = ioutil.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(content, testContent) != 0 {
		t.Fatalf("want %s, got %s", string(testContent), string(content))
	}
}
