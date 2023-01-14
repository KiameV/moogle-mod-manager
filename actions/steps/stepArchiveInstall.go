package steps

import (
	"fmt"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/files"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

const (
	z7url = "https://www.7-zip.org/download.html"
	z7cmd = "7z"
)

func checkFor7zip() (mods.Result, error) {
	if _, err := exec.Command("where", z7cmd).Output(); err != nil {
		wg := sync.WaitGroup{}
		wg.Add(1)
		d := dialog.NewCustom(
			"7-Zip not found",
			"Ok",
			container.NewCenter(widget.NewRichTextFromMarkdown(fmt.Sprintf(
				"Please download 7-Zip from [%s](%s) and install it.\n\n"+
					"Make sure to include it on the system path when installing.\n\n"+
					"Restart Moogle Mod Manager once 7-zip is installed.", z7url, z7url),
			)),
			ui.ActiveWindow())
		d.SetOnClosed(func() {
			wg.Done()
		})
		d.Show()
		wg.Wait()
		return mods.Cancel, nil
	}
	return mods.Ok, nil
}

func installDirectMoveToArchive(state *State, backupDir string) (mods.Result, error) {
	var (
		rel, name, bu string
		absArch       string
		installDir    string
		b             []byte
		ai            = newArchiveInjector()
		r, err        = checkFor7zip()
	)
	if r != mods.Ok {
		return r, err
	}

	if installDir, err = config.Get().GetDir(state.Game, config.GameDirKind); err != nil {
		return mods.Error, err
	} else if installDir == "" {
		return mods.Error, fmt.Errorf("install directory not found")
	}

	for _, e := range state.ExtractedFiles {
		for _, ti := range e.FilesToInstall() {
			absArch = filepath.Join(installDir, *ti.archive)
			if _, err = os.Stat(absArch); err != nil {
				return mods.Error, fmt.Errorf("archive not found: %s", absArch)
			}
			if rel, err = filepath.Rel(installDir, ti.AbsoluteTo); err != nil {
				return mods.Error, err
			}
			rel = filepath.Dir(rel)
			name = filepath.Base(ti.Relative)
			// Check if file already exists in the zip file
			cmd := exec.Command("l", absArch, fmt.Sprintf("%s/%s", rel, name))
			b, err = cmd.Output()
			if err == nil && !strings.Contains(string(b), "0 files") {
				// Extract file and move to backup directory
				bu = filepath.Join(backupDir, ArchiveAsDir(ti.archive), rel)
				if err = extractFile(absArch, rel, name, bu); err != nil {
					return mods.Error, err
				}
			}
			if err = ai.add(*ti.archive, ti.AbsoluteFrom, rel, name); err != nil {
				return mods.Error, err
			}
		}
	}
	if err = ai.updateArchives(state, installDir, archiveUpdate); err != nil {
		return mods.Error, err
	}
	return mods.Ok, nil
}

func uninstallDirectMoveToArchive(state *State) (mods.Result, error) {
	var (
		absBackup string
		gameDir   string
		backupDir string
		rel, name string
		ai        = newArchiveInjector()
		r, err    = checkFor7zip()
	)
	if r != mods.Ok {
		return r, err
	}
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
			_ = ai.add(a, absBackup, rel, name)
			// Ignore this error, in this case the file was not overridden
			// TODO May need to change archive files as whether they were added or removed
		}
	}
	if err = ai.updateArchives(state, gameDir, archiveRestoreBackup); err != nil {
		ai.revertFileMoves()
		return mods.Error, err
	}
	return mods.Ok, nil
}

func extractFile(archive, rel, name string, backupDir string) error {
	// Create the target directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}
	_ = os.Remove(filepath.Join(backupDir, name))
	// Extract the file to the target directory
	cmd := exec.Command(z7cmd, "e", archive, "-o"+backupDir, fmt.Sprintf("%s/%s", rel, name))
	if b, err := cmd.Output(); err != nil {
		return fmt.Errorf("%s: %s", err, b)
	}
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
		dir := absoluteFrom
		for !strings.HasSuffix(dir, "extracted") {
			d := dir
			dir = filepath.Dir(dir)
			if d == dir {
				return fmt.Errorf("could not find [extracted] directory")
			}
		}
		dir = filepath.Join(dir, asDir(archive))
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

func (i *archiveInjector) updateArchives(state *State, gameDir string, action archiveAction) (err error) {
	// Update the zip file
	var (
		cmd *exec.Cmd
		b   []byte
	)
	for archive, af := range i.archives {
		cmd = exec.Command(z7cmd, "a", filepath.Join(gameDir, string(archive)), af.dirToInject, "-r", "-y")
		if b, err = cmd.Output(); err != nil {
			err = fmt.Errorf("%s: %s", err, b)
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

func asDir(s string) string {
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "/", "_")
	return strings.ReplaceAll(s, "\\", "_")
}
