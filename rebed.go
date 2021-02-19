// Package rebed brings simple embedded file functionality
// to Go's new embed directive.
//
// It can recreate the directory structure
// from the embed.FS type with or without
// the files it contains. This is useful to
// expose the filesystem to the end user so they
// may see and modify the files.
//
// It also provides basic directory walking functionality for
// the embed.FS type.
package rebed

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// folderPerm MkdirAll is called with this permission to prevent restricted folders
// from being created.  0755=rwxr-xr-x
const folderPerm os.FileMode = 0755

// Tree creates the target filesystem folder structure.
func Tree(fsys embed.FS, outputPath string) error {
	return Walk(fsys, ".", func(dirpath string, de fs.DirEntry) error {
		fullpath := filepath.Join(outputPath, dirpath, de.Name())
		if de.IsDir() {
			return os.MkdirAll(fullpath, folderPerm)
		}
		return nil
	})
}

// Touch creates the target filesystem folder structure in the binary's
// current working directory with empty files. Does not modify
// already existing files.
func Touch(fsys embed.FS, outputPath string) error {
	return Walk(fsys, ".", func(dirpath string, de fs.DirEntry) error {
		fullpath := filepath.Join(outputPath, dirpath, de.Name())
		if de.IsDir() {
			return os.MkdirAll(fullpath, folderPerm)
		}
		// unsure how IsNotExist works. this could be improved
		_, err := os.Stat(fullpath)
		if os.IsNotExist(err) {
			_, err = os.Create(fullpath)
		}
		return err
	})
}

// Write overwrites files of same path/name
// in binaries current working directory or
// creates new ones if not exist.
func Write(fsys embed.FS, outputPath string) error {
	return Walk(fsys, ".", func(dirpath string, de fs.DirEntry) error {
		embedPath := filepath.Join(dirpath, de.Name())
		fullpath := filepath.Join(outputPath, embedPath)
		if de.IsDir() {
			return os.MkdirAll(fullpath, folderPerm)
		}
		return embedCopyToFile(fsys, embedPath, fullpath)
	})
}

// Patch creates files which are missing in
// FS filesystem. Does not modify existing files
func Patch(fsys embed.FS, outputPath string) error {
	return Walk(fsys, ".", func(dirpath string, de fs.DirEntry) error {
		embedPath := filepath.Join(dirpath, de.Name())
		fullpath := filepath.Join(outputPath, embedPath)
		if de.IsDir() {
			return os.MkdirAll(fullpath, folderPerm)
		}
		_, err := os.Stat(fullpath)
		if os.IsNotExist(err) {
			err = embedCopyToFile(fsys, embedPath, fullpath)
		}
		return err
	})
}

// Walk expects a relative path within fsys.
// f called on every file/directory found recursively.
//
// f's first argument is the relative/absolute path to directory being scanned.
// "." as startPath will scan all files and folders.
//
// Any error returned by f will cause Walk to return said error immediately.
func Walk(fsys embed.FS, startPath string, f func(path string, de fs.DirEntry) error) error {
	folders := make([]string, 0) // buffer of folders to process
	err := WalkDir(fsys, startPath, func(dirpath string, de fs.DirEntry) error {
		if de.IsDir() {
			folders = append(folders, filepath.Join(dirpath, de.Name()))
		}
		return f(dirpath, de)
	})
	if err != nil {
		return err
	}
	n := len(folders)
	for n != 0 {
		for i := 0; i < n; i++ {
			err = WalkDir(fsys, folders[i], func(dirpath string, de fs.DirEntry) error {
				if de.IsDir() {
					folders = append(folders, filepath.Join(dirpath, de.Name()))
				}
				return f(dirpath, de)
			})
			if err != nil {
				return err
			}
		}
		// we process n folders at a time, add new folders while
		//processing n folders, then discard those n folders once finished
		// and resume with a new n list of folders
		var newFolders int = len(folders) - n
		folders = folders[n : n+newFolders] // if found 0 new folders, end
		n = len(folders)
	}
	return nil
}

// WalkDir applies f to every file/folder in embedded directory fsys.
//
// f's first argument is the relative/absolute path to directory being scanned.
func WalkDir(fsys embed.FS, startPath string, f func(path string, de fs.DirEntry) error) error {
	items, err := fsys.ReadDir(startPath)
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := f(startPath, item); err != nil {
			return err
		}
	}
	return nil
}

// embedCopyToFile copies an embedded file's contents
// to a file in same relative path
func embedCopyToFile(fsys embed.FS, embedPath, path string) error {
	fi, err := fsys.Open(embedPath)
	if err != nil {
		return err
	}
	fo, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(fo, fi)
	return err
}
