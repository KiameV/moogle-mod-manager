package decompressor

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func handleArchive(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()
	if err = os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	// Closure to address file descriptors issue with all the deferred .Close() methods
	for _, f := range r.File {
		if err = extractArchiveFile(dest, f); err != nil {
			return err
		}
	}
	return nil
}

func extractArchiveFile(dest string, f *zip.File) (err error) {
	var (
		rc   io.ReadCloser
		file *os.File
		path string
	)
	if rc, err = f.Open(); err != nil {
		return
	}
	defer func() { _ = rc.Close() }()

	path = filepath.Join(dest, f.Name)
	// Check for ZipSlip (Directory traversal)
	if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
		err = fmt.Errorf("illegal file path: %s", path)
		return
	}

	if f.FileInfo().IsDir() {
		if err = os.MkdirAll(path, f.Mode()); err != nil {
			return
		}
	} else {
		if err = os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
			return
		}
		if file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode()); err != nil {
			return
		}
		defer func() { _ = file.Close() }()
		if _, err = io.Copy(file, rc); err != nil {
			return
		}
	}
	return
}
