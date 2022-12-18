package movers

import (
	"fmt"
	"os"
	"path/filepath"
)

type directMover struct{}

func (*directMover) CompileFilesToMove(dir string) (files []string, err error) {
	// Walk the destination directory tree
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}
		// Add the file to the list of files to move
		files = append(files, path)
		return nil
	})
	return
}

func (d *directMover) MoveFiles(files []string, from, destination, backupDir string) (backups []string, err error) {
	var moved []string
	if backups, err = d.backupFiles(files, from, destination, filepath.Join(destination, backupDir)); err != nil {
		return
	}
	if moved, err = d.moveFiles(files, from, destination); err != nil {
		d.removeFiles(moved)
		if _, err = d.moveFiles(backups, backupDir, destination); err != nil {
			err = fmt.Errorf("failed to restore files after failed mod install: %v", err)
		}
	}
	return
}

func (*directMover) moveFiles(files []string, from, destination string) (moved []string, err error) {
	var relPath string
	for _, f := range files {
		// Get the relative path of the file
		if relPath, err = filepath.Rel(from, f); err != nil {
			return
		}
		// Create the directory structure in the destination directory
		if err = os.MkdirAll(filepath.Join(destination, filepath.Dir(relPath)), 0755); err != nil {
			return
		}
		// Move the file to the destination directory
		if err = os.Rename(f, filepath.Join(destination, relPath)); err != nil {
			return
		}
		moved = append(moved, f)
	}
	return
}

func (*directMover) backupFiles(files []string, from, destination, backupDir string) (backedUp []string, err error) {
	var (
		relPath string
		orig    string
	)
	// Create the backup directory
	if err = os.MkdirAll(backupDir, 0755); err != nil {
		return
	}
	for _, f := range files {
		// Get the relative path of the file
		if relPath, err = filepath.Rel(from, f); err != nil {
			return
		}
		// Create the directory structure in the backup directory
		if err = os.MkdirAll(filepath.Join(backupDir, filepath.Dir(relPath)), 0755); err != nil {
			return
		}

		// See if a file already exists
		orig = filepath.Join(destination, relPath)
		if _, err = os.Stat(orig); err == nil {
			// Move the file to the backup directory
			if err = os.Rename(orig, filepath.Join(backupDir, relPath)); err != nil {
				return
			}
			backedUp = append(backedUp, relPath)
		}
	}
	return
}

func (d *directMover) removeFiles(files []string) {
	for _, f := range files {
		_ = os.Remove(f)
	}
}
