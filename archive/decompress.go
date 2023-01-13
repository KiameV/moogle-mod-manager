package archive

import (
	"context"
	"github.com/gen2brain/go-unarr"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/mholt/archiver/v4"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type (
	ExtractedFile struct {
		Name     string
		From     string
		Relative string
	}
	extractor struct {
		extracted               []ExtractedFile
		to                      string
		files                   map[string]bool
		dirsRecursive           []string
		includeBaseDirRecursive bool
		dirs                    map[string]bool
		includeBaseDir          bool
	}
)

func Decompress(from string, to string, continueIfExists bool, ti *mods.ToInstall) (extracted []ExtractedFile, err error) {
	var (
		f  *os.File
		fi os.FileInfo
		a  *unarr.Archive
		e  = newExtractor(to, ti)
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
		if f, err = os.Open(from); err != nil {
			return
		}
		err = archiver.Rar{}.Extract(context.Background(), f, nil, e.extractRar)
	} else { // zip/7z
		if a, err = unarr.NewArchive(from); err != nil {
			return
		}
		defer func() { _ = a.Close() }()

		if err = os.MkdirAll(to, 0777); err != nil {
			return
		}
		err = e.extractArchive(a)
	}
	extracted = e.extracted
	return
}

func newExtractor(to string, ti *mods.ToInstall) *extractor {
	e := &extractor{
		to:    to,
		files: make(map[string]bool),
		dirs:  make(map[string]bool),
	}
	for _, df := range ti.DownloadFiles {
		for _, f := range df.Files {
			e.files[f.From] = true
		}
		for _, d := range df.Dirs {
			if d.From == "." {
				if d.Recursive {
					e.includeBaseDirRecursive = true
					break
				} else {
					e.includeBaseDir = true
				}
			} else {
				if d.Recursive {
					e.dirsRecursive = append(e.dirsRecursive, d.From)
				} else {
					e.dirs[d.From] = true
				}
			}
		}
	}
	return e
}

func (e *extractor) extractRar(_ context.Context, f archiver.File) (err error) {
	if !f.IsDir() {
		var r io.ReadCloser
		if r, err = f.Open(); err != nil {
			return
		}
		defer func() { _ = r.Close() }()

		if e.shouldSkip(f.NameInArchive) {
			return nil
		}

		fp := filepath.Join(e.to, f.NameInArchive)
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

		e.extracted = append(e.extracted, ExtractedFile{
			Name:     strings.ReplaceAll(filepath.Base(fp), "\\", "/"),
			From:     strings.ReplaceAll(fp, "\\", "/"),
			Relative: strings.ReplaceAll(f.NameInArchive, "\\", "/"),
		})
	}
	return
}

func (e *extractor) extractArchive(a *unarr.Archive) (err error) {
	var (
		files []string
		rel   string
	)
	if files, err = a.Extract(e.to); err == nil {
		e.extracted = make([]ExtractedFile, 0, len(files))
		err = filepath.WalkDir(e.to,
			func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if d.IsDir() {
					return nil
				}
				rel, _ = filepath.Rel(e.to, path)
				if e.shouldSkip(rel) {
					_ = os.Remove(path)
					return nil
				}
				if rel, err = filepath.Rel(e.to, path); err != nil {
					return err
				}
				e.extracted = append(e.extracted, ExtractedFile{
					Name:     d.Name(),
					From:     path,
					Relative: rel,
				})
				return nil
			})
	}
	return
}

func (e *extractor) shouldSkip(path string) bool {
	var (
		lowerName = strings.ToLower(filepath.Base(path))
		found     bool
	)
	if strings.HasPrefix(lowerName, "readme") ||
		strings.HasPrefix(lowerName, "license") ||
		strings.HasPrefix(lowerName, ".git") ||
		strings.HasPrefix(lowerName, "__macosx") ||
		strings.HasPrefix(lowerName, ".ds_store") {
		return true
	}

	if e.includeBaseDirRecursive {
		return false
	}

	if _, found = e.files[path]; found {
		return false
	}

	dir := filepath.Dir(path)
	if e.includeBaseDir && dir == filepath.Dir(dir) {
		return false
	}

	if _, found = e.dirs[dir]; found {
		return false
	}
	for _, d := range e.dirsRecursive {
		if strings.HasPrefix(dir, d) {
			return false
		}
	}
	return true
}
