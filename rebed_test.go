package rebed_test

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/soypat/rebed"
)

// defer os.Chdir(strings.Repeat("../", len(filepath.SplitList(testDir))))

//go:embed testFS/*
var testFS embed.FS

// this is where filesystems are created
const testDir = "testdata/testdir"

func setup(path string, t *testing.T) {
	os.RemoveAll(path)
	os.MkdirAll(path, 0777)
	err := os.Chdir(path)
	if err != nil {
		t.Error(err)
	}
}
func TestTree(t *testing.T) {
	tDir := filepath.Join(testDir, t.Name())
	setup(tDir, t)
	defer os.RemoveAll(tDir)
	err := rebed.Tree(testFS, "")
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

func TestCreate(t *testing.T) {
	testFileCreation(rebed.Write, t)
}

func TestPatch(t *testing.T) {
	testFileCreation(rebed.Patch, t)
}

func testFileCreation(rebedder func(embed.FS, string) error, t *testing.T) {
	// shadow testDir
	tDir := filepath.Join(testDir, t.Name())
	setup(tDir, t)
	defer os.RemoveAll(tDir)
	err := rebedder(testFS, "")
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

func cleanup(path string) error {
	matcher := filepath.Join(path, "*")
	filesToRemove, err := filepath.Glob(matcher)
	if err != nil {
		return err
	}
	for _, f := range filesToRemove {
		os.RemoveAll(f)
	}
	return nil
}
