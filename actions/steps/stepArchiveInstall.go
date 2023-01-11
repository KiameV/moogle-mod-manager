package steps

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/files"
	"github.com/kiamev/moogle-mod-manager/mods"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*type zipReader struct {
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
}*/

func installDirectMoveToArchive(state *State, backupDir string) (mods.Result, error) {
	var (
		rel, name, bu   string
		absArch         string
		ai              = newArchiveInjector()
		installDir, err = config.Get().GetDir(state.Game, config.GameDirKind)
		z7              = config.PWD + "\\7z"
	)
	if err != nil {
		return mods.Error, err
	} else if installDir == "" {
		return mods.Error, fmt.Errorf("install directory not found")
	}

	for _, e := range state.ExtractedFiles {
		for _, ti := range e.FilesToInstall() {
			if rel, err = filepath.Rel(installDir, ti.AbsoluteTo); err != nil {
				return mods.Error, err
			}
			rel = filepath.Dir(rel)
			name = filepath.Base(ti.Relative)
			absArch = filepath.Join(installDir, *ti.archive)
			// Check if file already exists in the zip file
			cmd := exec.Command(z7, "l", absArch, fmt.Sprintf("%s/%s", rel, name))
			if err = cmd.Run(); err == nil {
				// Extract file and move to backup directory
				bu = filepath.Join(backupDir, ArchiveAsDir(ti.archive), rel)
				if err = extractFile(z7, absArch, rel, name, bu); err != nil {
					return mods.Error, err
				}
			}
			if err = ai.add(*ti.archive, ti.AbsoluteFrom, rel, name); err != nil {
				return mods.Error, err
			}
		}
	}
	if err = ai.updateArchives(state, z7, installDir, archiveUpdate); err != nil {
		return mods.Error, err
	}
	return mods.Ok, nil
}

func uninstallDirectMoveToArchive(state *State) (mods.Result, error) {
	var (
		z7 = config.PWD + "\\7z"
		ai = newArchiveInjector()

		absBackup string
		gameDir   string
		backupDir string
		rel, name string
		err       error
	)
	if gameDir, err = config.Get().GetDir(state.Game, config.GameDirKind); err != nil {
		return mods.Error, err
	}
	if backupDir, err = config.Get().GetDir(state.Game, config.BackupDirKind); err != nil {
		return mods.Error, err
	}
	for a, i := range files.Archives(state.Game, state.Mod.ID()) {
		for _, f := range i.Keys() {
			absBackup = filepath.Join(backupDir, ArchiveAsDir(&a), f)
			rel = filepath.Dir(f)
			name = filepath.Base(f)
			if err = ai.add(a, absBackup, rel, name); err != nil {
				ai.revertFileMoves()
				return mods.Error, err
			}
		}
	}
	if err = ai.updateArchives(state, z7, gameDir, archiveRestoreBackup); err != nil {
		ai.revertFileMoves()
		return mods.Error, err
	}
	return mods.Ok, nil
}

func extractFile(z7, archive, rel, name string, backupDir string) error {
	// Create the target directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}
	_ = os.Remove(filepath.Join(backupDir, name))
	// Extract the file to the target directory
	cmd := exec.Command(z7, "e", archive, "-o"+backupDir, fmt.Sprintf("%s/%s", rel, name))
	if err := cmd.Run(); err != nil {
		return err
	}
	// Remove the file from the zip file
	//cmd = exec.Command(z7, "d", archive, fmt.Sprintf("%s/%s", rel, name))
	//if err := cmd.Run(); err != nil {
	//	return err
	//}
	return nil
}

type (
	archiveAction byte
	archiveFile   string
	archiveFiles  struct {
		dirToInject string
		files       []string
	}
	archiveInjector struct {
		archives map[archiveFile]*archiveFiles
		files    []string
		renames  []fromTo
	}
	fromTo struct {
		from, to string
	}
)

const (
	_ archiveAction = iota
	archiveUpdate
	archiveRestoreBackup
)

func newArchiveInjector() *archiveInjector {
	return &archiveInjector{
		archives: make(map[archiveFile]*archiveFiles),
	}
}

func (i *archiveInjector) add(archive, absoluteFrom string, rel, name string) (err error) {
	af, ok := i.archives[archiveFile(archive)]
	rel = strings.ReplaceAll(rel, "\\", "/")
	if !ok {
		// Modify File Structure
		dir := filepath.Dir(absoluteFrom)
		dir = filepath.Join(dir, "aexta")
		if err = os.MkdirAll(filepath.Join(dir, rel), 0755); err != nil {
			return
		}
		if rel != "" && rel != "." {
			if strings.Contains(rel, "/") {
				sp := strings.Split(rel, "/")
				dir = filepath.Join(dir, sp[0])
			} else {
				dir = filepath.Join(dir, rel)
			}
		}
		af = &archiveFiles{dirToInject: dir}
		i.archives[archiveFile(archive)] = af
	}

	// Move the file to its new relative location
	to := filepath.Join(filepath.Dir(af.dirToInject), rel, name)
	i.renames = append(i.renames, fromTo{from: absoluteFrom, to: to})
	if err = os.Rename(absoluteFrom, to); err != nil {
		return
	}

	af.files = append(af.files, filepath.Join(rel, name))
	return
}

func (i *archiveInjector) updateArchives(state *State, z7 string, gameDir string, action archiveAction) (err error) {
	// Update the zip file
	for archive, af := range i.archives {
		cmd := exec.Command(z7, "a", filepath.Join(gameDir, string(archive)), af.dirToInject, "-r", "-y")
		if err = cmd.Run(); err != nil {
			return
		}
		if action == archiveRestoreBackup {
			files.RemoveArchiveFiles(state.Game, state.Mod.ID(), string(archive), af.files...)
		} else {
			files.AppendArchiveFiles(state.Game, state.Mod.ID(), string(archive), af.files...)
		}
		_ = os.RemoveAll(af.dirToInject)
	}
	return
}

func (i *archiveInjector) revertFileMoves() {
	for _, ft := range i.renames {
		_ = os.Rename(ft.to, ft.from)
	}
}
