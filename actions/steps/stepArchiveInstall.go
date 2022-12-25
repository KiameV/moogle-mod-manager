package steps

import (
	"archive/zip"
	"github.com/kiamev/moogle-mod-manager/mods"
	"io"
	"os"
	"path/filepath"
)

type zipReader struct {
	Reader *zip.ReadCloser
	Files  map[string]*zip.File
}

func newZipReader(s string) (z *zipReader, err error) {
	z = &zipReader{Files: make(map[string]*zip.File)}
	if z.Reader, err = zip.OpenReader(s); err != nil {
		return
	}
	for _, f := range z.Reader.File {
		z.Files[f.Name] = f
	}
	return
}

func (z *zipReader) close() {
	if z.Reader != nil {
		_ = z.Reader.Close()
	}
}

type zipWriter struct {
	File   *os.File
	Writer *zip.Writer
}

func newZipWriter(s string) (z *zipWriter, err error) {
	z = &zipWriter{}
	if z.File, err = os.OpenFile(s, os.O_RDWR, 0644); err == nil {
		z.Writer = zip.NewWriter(z.File)
	}
	return
}

func (z *zipWriter) close() {
	if z.Writer != nil {
		_ = z.Writer.Close()
	}
	if z.File != nil {
		_ = z.File.Close()
	}
}

func backupArchivedFiles(state *State, backupDir string) error {
	var (
		zr    map[string]*zipReader
		r     *zipReader
		found bool
		err   error
	)
	defer func() {
		for _, r = range zr {
			r.close()
		}
	}()
	for _, e := range state.ExtractedFiles {
		for _, ti := range e.FilesToInstall() {
			// Read all archives
			if ti.Skip {
				continue
			}
			if r, found = zr[*ti.archive]; !found {
				if r, err = newZipReader(*ti.archive); err != nil {
					return err
				}
				zr[*ti.archive] = r
			}
			if err = backupArchiveFile(r, backupDir, ti); err != nil {
				return err
			}
		}
	}
	return nil
}

func backupArchiveFile(r *zipReader, backupDir string, ti *FileToInstall) error {
	var (
		zf    *zip.File
		zfr   io.ReadCloser
		buf   *os.File
		found bool
		err   error
	)
	if zf, found = r.Files[ti.AbsoluteTo]; found {
		dir := filepath.Join(backupDir, *ti.archive, filepath.Dir(ti.AbsoluteTo))
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		if zfr, err = zf.Open(); err != nil {
			return err
		}
		defer func() { _ = zfr.Close() }()
		if buf, err = os.Create(dir); err != nil {
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
	/*	var (
			dir   string
			zw    map[string]*zipWriter
			w     *zipWriter
			found bool
			err   error
		)
		defer func() {
			for _, w = range zw {
				w.close()
			}
		}()

		if err = backupArchivedFiles(state, backupDir); err != nil {
			return mods.Error, err
		}

		for _, e := range state.ExtractedFiles {
			for _, ti := range e.FilesToInstall() {
				if ti.Skip {
					continue
				}
				if w, found = zw[*ti.archive]; !found {
					if w, err = newZipWriter(*ti.archive); err != nil {
						return mods.Error, err
					}
					zw[*ti.archive] = w
				}
				if err = copyFileToArchive(w, ti); err != nil {
					return
				}
			}
		}*/
	return mods.Ok, nil
}
