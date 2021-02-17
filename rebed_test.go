package rebed_test

import (
	"embed"
	"fmt"
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

func TestTree(t *testing.T) {
	os.RemoveAll(testDir)
	os.Mkdir(testDir, 0777)
	err := os.Chdir(testDir)
	if err != nil {
		t.Error(err)
	}
	err = rebed.Tree(testFS)
	if err != nil {
		t.Error(err)
	}
	// fsFolder := baseFolder(testFS)
	// We search our embedded directories and check if our created filesystem has the entries
	err = rebed.Walk(testFS, ".", func(path string, de fs.DirEntry) error {
		pathToCreated := filepath.Join(path, de.Name())
		if de.IsDir() {
			info, err := os.Stat(pathToCreated)
			if err != nil {
				t.Error(err)
			}
			if !info.IsDir() {
				t.Errorf("folder %q not found", pathToCreated)
			}
		}
		return nil
	})
}

func cleanup() error {
	matcher := filepath.Join(testDir, "*")
	filesToRemove, err := filepath.Glob(matcher)
	if err != nil {
		return err
	}
	for _, f := range filesToRemove {
		os.RemoveAll(f)
	}
	return nil
}

func isIn(path string, fsys embed.FS) bool {
	errFound := fmt.Errorf("found file!")
	errLook := rebed.Walk(fsys, ".", func(embedPath string, de fs.DirEntry) error {
		fullPath := filepath.Join(embedPath, de.Name())
		if fullPath == "." {
			return nil
		}
		if filepathCmp(fullPath, path) {
			return errFound // exits immediately
		}
		return nil
	})
	return errLook == errFound
}

func filepathCmp(a, b string) bool {
	return filepath.Clean(a) == filepath.Clean(b)
}

func baseFolder(fsys embed.FS) string {
	folder, err := fsys.ReadDir(".")
	if err != nil {
		panic(err)
	}
	return folder[0].Name()
}
