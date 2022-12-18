package archive

import (
	"context"
	"github.com/gen2brain/go-unarr"
	"github.com/mholt/archiver/v4"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Decompress(from string, to string, continueIfExists bool) (moved []string, err error) {
	var (
		f  *os.File
		fi os.FileInfo
		a  *unarr.Archive
	)
	if fi, err = os.Stat(to); err == nil && fi.IsDir() {
		var fis []os.DirEntry
		if fis, err = os.ReadDir(to); err == nil && len(fis) > 0 {
			if !continueIfExists {
				return nil, nil
			}
		}
	}
	if filepath.Ext(from) == ".rar" {
		handler := func(ctx context.Context, f archiver.File) (err error) {
			if !f.IsDir() {
				var r io.ReadCloser
				if r, err = f.Open(); err != nil {
					return
				}
				defer func() { _ = r.Close() }()

				fp := filepath.Join(to, f.NameInArchive)
				if err = os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
					return
				}
				buf := new(strings.Builder)
				if _, err = io.Copy(buf, r); err != nil {
					return
				}
				var file *os.File
				if file, err = os.Create(fp); err != nil {
					return
				}
				defer func() { _ = file.Close() }()

				_, err = file.WriteString(buf.String())
				moved = append(moved, fp)
			}
			return
		}

		if f, err = os.Open(from); err != nil {
			return nil, err
		}
		return moved, archiver.Rar{}.Extract(context.Background(), f, nil, handler)
	}

	if a, err = unarr.NewArchive(from); err != nil {
		return nil, err
	}
	defer func() { _ = a.Close() }()

	if err = os.MkdirAll(to, 0777); err != nil {
		return nil, err
	}

	return a.Extract(to)
}
