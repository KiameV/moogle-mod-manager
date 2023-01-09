package archive_io

import (
	"archive/zip"
	"os"
)

type (
	ZipIO interface {
		Writer(zipFile string) (*zip.Writer, error)
		HasFile(file string) (*zip.File, bool)
		LoadFiles() error
		Close()
	}
	zipIO struct {
		zipFile string
		file    *os.File
		files   map[string]*zip.File
		writer  *zip.Writer
		reader  *zip.ReadCloser
	}
)

func NewZipIO(zipFile string) ZipIO {
	return &zipIO{zipFile: zipFile}
}

func (z *zipIO) LoadFiles() error {
	r, err := zip.OpenReader(z.zipFile)
	if err != nil {
		return err
	}
	z.reader = r
	z.files = make(map[string]*zip.File)
	for _, f := range r.File {
		z.files[f.Name] = f
	}
	return nil
}

func (z *zipIO) HasFile(file string) (f *zip.File, found bool) {
	f, found = z.files[file]
	return
}

func (z *zipIO) Writer(zipFile string) (*zip.Writer, error) {
	var err error
	if z.writer == nil {
		if z.file, err = os.OpenFile(zipFile, os.O_RDWR, 0644); err == nil {
			z.writer = zip.NewWriter(z.file)
		}
	}
	return z.writer, err
}

func (z *zipIO) Close() {
	if z.reader != nil {
		_ = z.reader.Close()
		z.reader = nil
	}
	if z.writer != nil {
		_ = z.writer.Close()
		z.writer = nil
	}
	if z.file != nil {
		_ = z.file.Close()
		z.file = nil
	}
}
