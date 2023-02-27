package steps

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/kiamev/moogle-mod-manager/archive"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type (
	Extracted struct {
		ToInstall      *mods.ToInstall
		Files          []archive.ExtractedFile
		filesToInstall []*FileToInstall
	}
	FileToInstall struct {
		Relative     string
		AbsoluteFrom string
		AbsoluteTo   string
		Skip         bool
		archive      *string
	}
)

func ArchiveAsDir(archive *string) string {
	if archive != nil {
		return strings.ReplaceAll(*archive, ".", "_")
	}
	return ""
}

func newFileToInstallFromFile(relToExtracted map[string]archive.ExtractedFile, f *mods.ModFile, installDir string, archive *string) (*FileToInstall, error) {
	af, found := relToExtracted[f.From]
	if !found {
		return nil, fmt.Errorf("file %v not found in extracted files", f)
	}
	return &FileToInstall{
		Relative:     f.To,
		AbsoluteFrom: af.From,
		AbsoluteTo:   filepath.Join(installDir, f.To),
		Skip:         false,
		archive:      archive,
	}, nil
}

func newFileToInstallFromDir(relToExtracted map[string]archive.ExtractedFile, rel string, d *mods.ModDir, installDir string, archive *string) (*FileToInstall, error) {
	var (
		af, found = relToExtracted[rel]
		toRel     = rel
	)
	if !found {
		return nil, fmt.Errorf("dir %v not found in extracted files", d.From)
	}
	if d.From != "." {
		toRel = strings.TrimPrefix(rel, d.From)
	}
	return &FileToInstall{
		Relative:     af.Relative,
		AbsoluteFrom: af.From,
		AbsoluteTo:   filepath.Join(installDir, d.To, toRel),
		Skip:         false,
		archive:      archive,
	}, nil
}

func (e *Extracted) FilesToInstall() []*FileToInstall {
	return e.filesToInstall
}

func (e *Extracted) Compile(game config.GameDef, extractedDir string) (err error) {
	if len(e.filesToInstall) > 0 {
		return
	}

	var (
		fromToExtracted = make(map[string]archive.ExtractedFile)
		installDir      string
		fti             *FileToInstall
		rel             string
	)
	if installDir, err = config.Get().GetDir(game, config.GameDirKind); err != nil {
		return
	}

	for _, f := range e.Files {
		fromToExtracted[strings.ReplaceAll(f.Relative, "\\", "/")] = f
	}
	for _, df := range e.ToInstall.DownloadFiles {
		for _, f := range df.Files {
			if fti, err = newFileToInstallFromFile(fromToExtracted, f, installDir, f.ToArchive); err != nil {
				return
			}
			e.filesToInstall = append(e.filesToInstall, fti)
			/*ex, found := fromToExtracted[f.From]
			if !found {
				return nil, fmt.Errorf("file %v not found", f.From)
			}
			if f.From == f.To {
				result = append(result, ex.Relative)
			} else {
				// Add root directory
				filepath.Base(ex.File)
				result = append(result, filepath.Join(f.To, filepath.Base(ex.File)))
			}*/
		}
		for _, d := range df.Dirs {
			if err = filepath.WalkDir(filepath.Join(extractedDir, d.From), func(path string, de fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if de.IsDir() {
					return nil
				}
				path = filepath.ToSlash(path)
				if rel, err = filepath.Rel(extractedDir, path); err != nil {
					return err
				}
				rel = strings.ReplaceAll(rel, "\\", "/")
				if fti, err = newFileToInstallFromDir(fromToExtracted, rel, d, installDir, d.ToArchive); err != nil {
					return err
				}
				e.filesToInstall = append(e.filesToInstall, fti)
				return nil
			}); err != nil {
				return
			}
		}
	}
	return
}
