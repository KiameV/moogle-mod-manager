package steps

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/archive"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"io/fs"
	"path/filepath"
	"strings"
)

type (
	Extracted struct {
		ToInstall      *mods.ToInstall
		Files          []archive.ExtractedFile
		filesToInstall []*FileToInstall
	}
	Archive struct {
		Name string
		Path string
	}
	FileToInstall struct {
		Relative     string
		AbsoluteFrom string
		AbsoluteTo   string
		Skip         bool
		archive      *Archive
	}
)

func newFileToInstallFromFile(relToExtracted map[string]archive.ExtractedFile, f *mods.ModFile, installDir string, archive *string) (*FileToInstall, error) {
	af, found := relToExtracted[f.From]
	if !found {
		return nil, fmt.Errorf("file %v not found in extracted files", f)
	}
	if archive != nil {
		// TODO
	}
	return &FileToInstall{
		Relative:     f.To,
		AbsoluteFrom: af.From,
		AbsoluteTo:   filepath.Join(installDir, f.To),
		Skip:         false,
		archive:      nil,
	}, nil
}

func newFileToInstallFromDir(relToExtracted map[string]archive.ExtractedFile, fromRel string, d *mods.ModDir, installDir string, archive *string) (*FileToInstall, error) {
	var (
		af, found = relToExtracted[fromRel]
		to        string
		a         *Archive
		err       error
	)
	if !found {
		return nil, fmt.Errorf("dir %v not found in extracted files", d.From)
	}
	if d.From != "." {
		if to, err = filepath.Rel(d.From, fromRel); err != nil {
			return nil, err
		}
		to = filepath.Join(d.To, to)
	}
	if archive == nil {
		to = filepath.Join(installDir, to)
	} else {
		a = &Archive{
			Name: *archive,
			Path: strings.ReplaceAll(filepath.Join(installDir, *archive), "\\", "/"),
		}
	}
	return &FileToInstall{
		Relative:     af.Relative,
		AbsoluteFrom: af.From,
		AbsoluteTo:   strings.ReplaceAll(to, "\\", "/"),
		Skip:         false,
		archive:      a,
	}, nil
}

func (e *Extracted) FilesToInstall() []*FileToInstall {
	return e.filesToInstall
}

func (e *Extracted) Compile(game config.GameDef, it config.InstallType, extractedDir string) (err error) {
	if it == config.MoveToArchive {
		return e.compileForMoveToArchive(game, extractedDir)
	}
	return e.compileForMove(game, extractedDir)
}

func (e *Extracted) compileForMove(game config.GameDef, extractedDir string) (err error) {
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
		}
		for _, d := range df.Dirs {
			if err = filepath.WalkDir(filepath.Join(extractedDir, d.From), func(p string, de fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if de.IsDir() {
					return nil
				}
				path := filepath.ToSlash(p)
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

func (e *Extracted) compileForMoveToArchive(game config.GameDef, extractedDir string) (err error) {
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
