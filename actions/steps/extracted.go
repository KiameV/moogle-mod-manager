package steps

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/files/archive"
	"github.com/kiamev/moogle-mod-manager/mods"
	"io/fs"
	"path/filepath"
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
	}
)

func newFileToInstallFromFile(fromToExtracted map[string]archive.ExtractedFile, f *mods.ModFile, installDir string) (*FileToInstall, error) {
	af, found := fromToExtracted[f.From]
	if !found {
		return nil, fmt.Errorf("file %v not found in extracted files", f)
	}
	return &FileToInstall{
		Relative:     f.To,
		AbsoluteFrom: af.File,
		AbsoluteTo:   filepath.Join(installDir, f.To),
		Skip:         false,
	}, nil
}

func newFileToInstallFromDir(fromToExtracted map[string]archive.ExtractedFile, d *mods.ModDir, installDir string) (*FileToInstall, error) {
	af, found := fromToExtracted[d.From]
	if !found {
		return nil, fmt.Errorf("file %v not found in extracted files", f)
	}
	if path, err = filepath.Rel(d.From, path); err != nil {
		return err
	}
	return &FileToInstall{
		Relative:     d.To,
		AbsoluteFrom: af.File,
		AbsoluteTo:   filepath.Join(installDir, d.To),
		Skip:         false,
	}, nil
}

func (e *Extracted) FilesRelative(game config.GameDef) (result []*FileToInstall, err error) {
	if len(e.filesToInstall) > 0 {
		return e.filesToInstall, nil
	}

	var (
		fromToExtracted = make(map[string]archive.ExtractedFile)
		installDir      string
		fti             *FileToInstall
	)
	if installDir, err = config.Get().GetDir(game, config.GameDirKind); err != nil {
		return
	}

	for _, f := range e.Files {
		fromToExtracted[f.Relative] = f
	}
	for _, df := range e.ToInstall.DownloadFiles {
		for _, f := range df.Files {
			if fti, err = newFileToInstallFromFile(fromToExtracted, f, installDir); err != nil {
				return
			}
			result = append(result, fti)
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
			if err = filepath.WalkDir(d.From, func(path string, de fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if de.IsDir() {
					return nil
				}
				if fti, err = newFileToInstallFromDir(fromToExtracted, dir, installDir); err != nil {
					return err
				}
				result = append(result, fti)
				return nil
			}); err != nil {
				result = nil
				break
			}
		}
	}
	return
}
