// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package volume

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestEphemeral(t *testing.T) {
	testID := "TestID"
	// create tmpdir as base for volume
	baseDir, err := ioutil.TempDir("", "test-ephemeral-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)

	e := ephemeral{
		baseDir: baseDir,
	}

	err = e.create(testID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(e.path)
	if os.IsNotExist(err) {
		t.Fatalf("failed to create ephemeral volume location")
	}

	// write to a test file in the volume
	testData := []byte("testfile")
	testPath := filepath.Join(e.handle(), "test-file")
	if err := ioutil.WriteFile(testPath, testData, 0644); err != nil {
		t.Fatalf("failed to write to test file in volume")
	}

	err = e.delete()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(testPath)
	if !os.IsNotExist(err) {
		t.Fatalf("failed to remove ephemeral volume location")
	}
}
