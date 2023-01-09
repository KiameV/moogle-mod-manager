package steps

import (
	"archive/zip"
	"errors"
	"fmt"
	aio "github.com/kiamev/moogle-mod-manager/actions/archive-io"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/util"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func backupArchivedFiles(state *State, backupDir string) error {
	var (
		zio   aio.ZipIO
		found bool
		err   error
	)
	for _, e := range state.ExtractedFiles {
		for _, ti := range e.FilesToInstall() {
			if ti.archive == nil {
				return fmt.Errorf("archive not set for file %s", ti.AbsoluteTo)
			}
			// Read all archives
			if ti.Skip {
				continue
			}
			if zio, found = state.ZipIO[ti.archive.Path]; !found {
				zio = aio.NewZipIO(ti.archive.Path)
				state.ZipIO[ti.archive.Path] = zio
			}
			if err = backupArchiveFile(zio, backupDir, ti); err != nil {
				return err
			}
		}
	}
	return nil
}

func backupArchiveFile(zio aio.ZipIO, backupDir string, ti *FileToInstall) error {
	var (
		err   = zio.LoadFiles()
		zf    *zip.File
		zfr   io.ReadCloser
		buf   *os.File
		found bool
	)
	if err != nil {
		return err
	}
	if zf, found = zio.HasFile(ti.AbsoluteTo); found {
		dir := filepath.Join(backupDir, ti.archive.Name, filepath.Dir(ti.AbsoluteTo))
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		if zfr, err = zf.Open(); err != nil {
			return err
		}
		defer func() { _ = zfr.Close() }()

		file := filepath.Join(dir, ti.Relative)
		if buf, err = os.Create(file); err != nil {
			return err
		}
		defer func() { _ = buf.Close() }()
		if _, err = io.Copy(buf, zfr); err != nil {
			return err
		}
	}
	return nil
}

func installDirectMoveToArchive(state *State, backupDir string) (mods.Result, error) {
	var (
		archiveDirs = make(map[string]string)
		tmp         = filepath.Join(config.PWD, "tmp")
		to          string
		b           []byte
		err         error
	)
	_ = os.RemoveAll(tmp)
	_ = os.MkdirAll(tmp, 0755)
	defer func() {
		_ = os.RemoveAll(tmp)
	}()

	for _, e := range state.ExtractedFiles {
		for _, ti := range e.FilesToInstall() {
			if ti.Skip {
				continue
			}
			if ti.archive != nil {
				to = filepath.Join(tmp, util.CreateFileName(ti.archive.Name))
				archiveDirs[ti.archive.Path] = to

				to = filepath.Join(to, strings.TrimRight(ti.AbsoluteTo, ti.Relative))
				if err = os.MkdirAll(to, 0777); err != nil {
					return mods.Error, err
				}
				to = filepath.Join(to, ti.Relative)
				if err = os.Rename(ti.AbsoluteFrom, to); err != nil {
					return mods.Error, err
				}
			}
		}
	}

	if err = backupArchivedFiles(state, backupDir); err != nil {
		return mods.Error, err
	}

	for archive, path := range archiveDirs {
		if b, err = exec.Command("cmd", "/C", config.PWD, "7z.exe", "u", archive, path).Output(); err != nil {
			if len(b) > 0 {
				err = errors.New(string(b))
			}
			return mods.Error, err
		}
	}

	return mods.Ok, nil
}

/*
// ZipFiles compresses one or many files into a single zip archive file.
// The name of the file will be the first file's name with the ".zip" extension.
func ZipFiles(filename string, files []string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func updateFileInZip(zipFile *zip.File, fileToUpdateWith string) error {
	// Open the file to update
	f, err := os.Open(fileToUpdateWith)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	// Get the file information
	info, err := f.Stat()
	if err != nil {
		return err
	}

	// Create a new header for the file
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Method = zip.Deflate

	// Open the file in the zip archive for writing
	w, err := zipFile.(header)
	if err != nil {
		return err
	}

	// Write the updated file to the zip archive
	_, err = io.Copy(w, f)
	return err
}
*/
