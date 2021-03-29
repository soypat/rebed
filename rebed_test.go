package rebed_test

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/chengziqing/rebed"
)

// defer os.Chdir(strings.Repeat("../", len(filepath.SplitList(testDir))))

//go:embed testFS/*
var testFS embed.FS

// this is where filesystems are created
const testDir = "testdata/testdir"

func setup(path string, t *testing.T) {
	os.RemoveAll(path)
	os.MkdirAll(path, 0777)
}
func TestTree(t *testing.T) {
	tDir := filepath.Join(testDir, t.Name())
	setup(tDir, t)
	defer os.RemoveAll(tDir)
	err := rebed.Tree(testFS, tDir)
	if err != nil {
		t.Error(err)
	}
	// We search our embedded directories and check if our created filesystem has the entries
	err = rebed.Walk(testFS, ".", func(path string, de fs.DirEntry) error {
		pathToCreated := filepath.Join(path, de.Name())
		if de.IsDir() {
			info, err := os.Stat(pathToCreated)
			if os.IsNotExist(err) {
				t.Errorf("folder %q not found", pathToCreated)
			}
			if !info.IsDir() {
				t.Errorf("expected a folder %q, got file", pathToCreated)
			}
			if err != nil {
				t.Error(err)
			}
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestTouch(t *testing.T) {
	testFileCreation(rebed.Touch, t)
}

func TestWrite(t *testing.T) {
	testFileCreation(rebed.Write, t)
}

func TestPatch(t *testing.T) {
	testFileCreation(rebed.Patch, t)
}

func TestCreate(t *testing.T) {
	testFileCreation(rebed.Create, t)
}

func TestCreateError(t *testing.T) {
	tDir := filepath.Join(testDir, t.Name())
	setup(tDir, t)
	defer os.RemoveAll(tDir)
	err := rebed.Create(testFS, tDir)
	if err != nil {
		t.Errorf("Create failed with new directory %s", err)
	}
	err = rebed.Create(testFS, tDir)
	if err == nil || !os.IsExist(err) {
		t.Errorf("Create should have failed during conflicting filesystem creation: %s", err)
	}
}

func TestWalkError(t *testing.T) {
	theError := fmt.Errorf("oops")
	err := rebed.Walk(testFS, ".", func(path string, de fs.DirEntry) error {
		return theError
	})
	if err != theError {
		t.Errorf("Expected specific error, got %s", err)
	}
}

func TestWalkDirError(t *testing.T) {
	err := rebed.Walk(testFS, "nonexistent", func(path string, de fs.DirEntry) error {
		return nil
	})
	if err == nil {
		t.Errorf("expected error while walking non-existent directory")
	}
}

func testFileCreation(rebedder func(embed.FS, string) error, t *testing.T) {
	tDir := filepath.Join(testDir, t.Name())
	setup(tDir, t)
	defer os.RemoveAll(tDir)
	err := rebedder(testFS, tDir)
	if err != nil {
		t.Error(err)
	}
	err = rebed.Walk(testFS, ".", func(path string, de fs.DirEntry) error {
		pathToCreated := filepath.Join(path, de.Name())
		info, err := os.Stat(pathToCreated)
		if err != nil {
			return err
		}
		if de.IsDir() != info.IsDir() {
			t.Errorf("expected folder/file got file/folder %q", pathToCreated)
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
