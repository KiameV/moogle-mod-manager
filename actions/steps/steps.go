package steps

import (
	"errors"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/downloads"
	"github.com/kiamev/moogle-mod-manager/files"
	"github.com/kiamev/moogle-mod-manager/files/archive"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	ci "github.com/kiamev/moogle-mod-manager/ui/config-installer"
	"github.com/kiamev/moogle-mod-manager/ui/confirm"
	uic "github.com/kiamev/moogle-mod-manager/ui/conflicts"
	ui "github.com/kiamev/moogle-mod-manager/ui/state"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type (
	State struct {
		Game           config.GameDef
		Mod            mods.TrackedMod
		Downloaded     []string
		ToInstall      []*mods.ToInstall
		ExtractedFiles []Extracted
	}
	Step func(state *State) (result mods.Result, err error)
)

func NewState(game config.GameDef, mod mods.TrackedMod) *State {
	return &State{
		Game: game,
		Mod:  mod,
	}
}

func VerifyEnable(state *State) (mods.Result, error) {
	var (
		tm      = state.Mod
		c       = tm.Mod().ModCompatibility
		mc      *mods.ModCompat
		mod     mods.TrackedMod
		found   bool
		enabled bool
	)
	if c != nil {
		if len(c.Forbids) > 0 {
			for _, mc = range c.Forbids {
				if mod, found, enabled = managed.IsModEnabled(state.Game, mc.ModID()); found && enabled {
					return mods.Error, fmt.Errorf("[%s] cannot be enabled because [%s] is enabled", tm.DisplayName(), mod.DisplayName())
				}
			}
		}
		if len(c.Requires) > 0 {
			for _, mc = range c.Requires {
				mod, found, enabled = managed.IsModEnabled(state.Game, mc.ModID())
				if !found {
					return mods.Error, fmt.Errorf("[%s] cannot be enabled because [%s] is not enabled", tm.DisplayName(), mc.ModID())
				} else if !enabled {
					return mods.Error, fmt.Errorf("[%s] cannot be enabled because [%s] is not enabled", tm.DisplayName(), mod.DisplayName())
				}
			}
		}
	}
	return mods.Ok, nil
}

func VerifyDisable(state *State) (mods.Result, error) {
	var (
		tm  = state.Mod
		mod = tm.Mod()
		c   = mod.ModCompatibility
		id  = mod.ModID
		mc  *mods.ModCompat
	)
	for _, m := range managed.GetMods(state.Game) {
		if m.ID() != tm.ID() && m.Enabled() {
			if c = m.Mod().ModCompatibility; c != nil {
				if len(c.Requires) > 0 {
					for _, mc = range c.Requires {
						if mc.ModID() == id {
							return mods.Error, fmt.Errorf("[%s] cannot be disabled because [%s] is enabled", tm.DisplayName(), m.DisplayName())
						}
					}
				}
			}
		}
	}
	return mods.Ok, nil
}

func PreDownload(state *State) (result mods.Result, err error) {
	var (
		mod = state.Mod.Mod()
		wg  sync.WaitGroup
	)
	if state.ToInstall, err = mods.NewToInstallForMod(mod.Kind(), mod, mod.AlwaysDownload); err != nil {
		return mods.Error, err
	}
	// Handle any mod configurations
	if len(mod.Configurations) > 0 {
		wg.Add(1)
		modPath := filepath.Join(config.Get().GetModsFullPath(state.Game), mod.ID().AsDir())
		if err = ui.GetScreen(ui.ConfigInstaller).(ci.ConfigInstaller).Setup(mod, modPath, func(r mods.Result, ti []*mods.ToInstall) error {
			result = r
			if len(ti) > 0 {
				state.ToInstall = append(state.ToInstall, ti...)
			}
			wg.Done()
			return nil
		}); err != nil {
			// Failed to set up config installer screen
			return mods.Error, err
		}
		ui.ShowScreen(ui.ConfigInstaller)
		wg.Wait()
		time.Sleep(100 * time.Millisecond)

	}

	if len(state.ToInstall) == 0 {
		return mods.Error, errors.New("no files to install")
	}

	wg.Add(1)
	// Confirm Download
	confirmer := confirm.NewConfirmer(confirm.NewParams(state.Game, state.Mod, state.ToInstall))
	if err = confirmer.Downloads(func(r mods.Result) {
		result = r
		wg.Done()
	}); err == nil {
		wg.Wait()
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		return mods.Error, err
	}
	return result, nil
}

func Download(state *State) (result mods.Result, err error) {
	if err = downloads.Download(state.Game, state.Mod, state.ToInstall); err != nil {
		result = mods.Error
	} else {
		result = mods.Ok
	}
	return
}

func Extract(state *State) (mods.Result, error) {
	var (
		to       string
		override bool
		ef       []archive.ExtractedFile
		err      error
	)
	for _, ti := range state.ToInstall {
		if state.Mod.InstallType(state.Game) == config.ImmediateDecompress {
			if to, err = config.Get().GetDir(state.Game, config.GameDirKind); err != nil {
				return mods.Error, err
			}
		} else {
			to = ti.Download.DownloadedArchiveLocation.ExtractDir()
		}
		if state.Mod.InstallType(state.Game) == config.ImmediateDecompress {
			if to, err = config.Get().GetDir(state.Game, config.GameDirKind); err != nil {
				return mods.Error, err
			}
			override = true
		} else {
			//override = false
			override = true
		}

		if ef, err = archive.Decompress(string(*ti.Download.DownloadedArchiveLocation), to, override); err != nil {
			return mods.Error, err
		}

		if state.Mod.InstallType(state.Game) == config.ImmediateDecompress {
			installed := make([]string, len(ef))
			for i, f := range ef {
				installed[i] = f.From
			}
			files.SetFiles(state.Game, state.Mod.ID(), installed...)
		} else {
			e := Extracted{
				ToInstall: ti,
				Files:     ef,
			}
			if err = e.Compile(state.Game, to); err != nil {
				return mods.Error, err
			}
			state.ExtractedFiles = append(state.ExtractedFiles, e)
		}
	}
	return mods.Ok, nil
}

func Conflicts(state *State) (result mods.Result, err error) {
	var (
		mod            = state.Mod.Mod()
		tos            []string
		tosToToInstall = make(map[string]*FileToInstall)
		wg             sync.WaitGroup
		ti             *FileToInstall
		found          bool
	)
	for _, e := range state.ExtractedFiles {
		for _, ti = range e.FilesToInstall() {
			tos = append(tos, ti.AbsoluteTo)
			tosToToInstall[ti.AbsoluteTo] = ti
		}
	}

	conflicts := files.FindConflicts(state.Game, tos)

	result = mods.Ok
	if len(conflicts) > 0 {
		wg.Add(1)
		uic.ShowConflicts(state.Mod.Mod(), conflicts, func(r mods.Result) {
			result = r
			if result == mods.Ok {
				for _, c := range conflicts {
					if c.Selection != mod {
						// Use other mod
						if ti, found = tosToToInstall[c.Path]; found {
							ti.Skip = true
						}
					} else {
						// Use this mod
						files.RemoveFiles(state.Game, c.Owner.ID(), c.Path)
					}
				}
			}
			wg.Done()
		})
		wg.Wait()
	}
	if err != nil {
		return mods.Error, err
	}
	return result, nil
}

func Install(state *State) (mods.Result, error) {
	var (
		backupDir string
		err       error
	)
	if backupDir, err = config.Get().GetDir(state.Game, config.BackupDirKind); err != nil {
		return mods.Error, err
	}
	switch state.Mod.InstallType(state.Game) {
	case config.Move:
		for _, e := range state.ExtractedFiles {
			for _, ti := range e.FilesToInstall() {
				if ti.Skip {
					continue
				}

				dir := filepath.Dir(ti.AbsoluteTo)
				if _, err = os.Stat(dir); err != nil {
					// Create the directory structure
					if err = os.MkdirAll(dir, 0755); err != nil {
						return mods.Error, err
					}
				} else if _, err = os.Stat(ti.AbsoluteTo); err == nil {
					// File Exists
					// See if there's a file backup
					absBackup := filepath.Join(backupDir, ti.Relative)
					if _, err = os.Stat(absBackup); err == nil {
						// Backup Exists
						if err = os.Remove(ti.AbsoluteTo); err != nil {
							return mods.Error, err
						}
					} else {
						// No Backup
						if err = os.MkdirAll(filepath.Dir(absBackup), 0755); err != nil {
							return mods.Error, err
						}
						if err = os.Rename(ti.AbsoluteTo, absBackup); err != nil {
							return mods.Error, err
						}
					}
				}

				// Install the file
				if err = os.Rename(ti.AbsoluteFrom, ti.AbsoluteTo); err != nil {
					return mods.Error, err
				}
				files.SetFiles(state.Game, state.Mod.ID(), ti.AbsoluteTo)
			}
		}
	case config.MoveToArchive:
		// TODO
		panic("not implemented")
	default:
		return mods.Error, fmt.Errorf("unknown install type: %v", state.Mod.InstallType(state.Game))
	}

	return mods.Ok, nil
}

func UninstallMove(state *State) (mods.Result, error) {
	var (
		i         = files.Files(state.Game, state.Mod.ID())
		gameDir   string
		backupDir string
		rel       string
		err       error
	)
	if gameDir, err = config.Get().GetDir(state.Game, config.GameDirKind); err != nil {
		return mods.Error, err
	}
	if backupDir, err = config.Get().GetDir(state.Game, config.BackupDirKind); err != nil {
		return mods.Error, err
	}
	for _, f := range i.Keys() {
		if err = os.Remove(f); err != nil {
			return mods.Error, err
		}
		files.RemoveFiles(state.Game, state.Mod.ID(), f)

		if rel, err = filepath.Rel(gameDir, f); err != nil {
			return mods.Error, err
		}

		absBackup := filepath.Join(backupDir, rel)
		if _, err = os.Stat(absBackup); err == nil {
			if err = os.Rename(absBackup, f); err != nil {
				return mods.Error, err
			}
		}
	}
	return mods.Ok, nil
}

func EnableMod(state *State) (result mods.Result, err error) {
	state.Mod.Enable()
	result = mods.Ok
	if err = managed.Save(); err != nil {
		result = mods.Error
	}
	return
}

func DisableMod(state *State) (result mods.Result, err error) {
	state.Mod.Disable()
	result = mods.Ok
	if err = managed.Save(); err != nil {
		result = mods.Error
	}
	return
}

func ShowWorkingDialog(_ *State) (mods.Result, error) {
	return mods.Working, nil
}

func PostInstall(state *State) (mods.Result, error) {
	for _, id := range files.EmptyMods(state.Game) {
		if m, found := managed.TryGetMod(state.Game, id); m != nil && found {
			m.Disable()
		}
	}
	for _, ti := range state.ToInstall {
		l := ti.Download.DownloadedArchiveLocation
		if l != nil {
			_ = os.RemoveAll(ti.Download.DownloadedArchiveLocation.ExtractDir())
			if config.Get().DeleteDownloadAfterInstall {
				_ = os.Remove(string(*l))
			}
		}
	}
	return mods.Ok, nil
}

/*
func enableMod(enabler *mods.ModEnabler, err error) {
	var (
		game       = enabler.Game
		tm         = enabler.TrackedMod
		tis        = enabler.ToInstall
		to         string
		kind       mods.Kind
		movedFiles []string
		modPath    = filepath.Join(config.Get().GetModsFullPath(game), tm.ID().AsDir())
	)
	if err != nil {
		tm.Disable()
		enabler.DoneCallback(mods.Error, err)
		return
	}
	enabler.ShowWorking()

	for _, ti := range tis {
		if tm.Mod().InstallType(enabler.Game) {
			if to, err = config.Get().GetDir(game, config.GameDirKind); err != nil {
				tm.Disable()
				enabler.DoneCallback(mods.Error, err)
				return
			}
			to = filepath.Join(to, enabler.TrackedMod.Mod().AlwaysDownload[0].Dirs[0].To)
			if movedFiles, err = decompress(*ti.Download.DownloadedArchiveLocation, to, true); err != nil {
				tm.Disable()
				enabler.DoneCallback(mods.Error, err)
				return
			}
			for i, f := range movedFiles {
				movedFiles[i] = filepath.Join(to, f)
			}
			managed.GetManagedFiles(game, tm.ID())
			tm.Enabled()
			enabler.DoneCallback(mods.Ok, nil)
		} else {
			to = filepath.Join(modPath, ti.Download.Name)
			if _, err = decompress(*ti.Download.DownloadedArchiveLocation, to, false); err != nil {
				tm.Disable()
				enabler.DoneCallback(mods.Error, err)
				return
			}
		}

		kind = tm.Kind()
		if kind == mods.Nexus || kind == mods.CurseForge {
			var fi os.FileInfo
			sa := filepath.Join(to, "StreamingAssets")
			if fi, err = os.Stat(sa); err == nil && fi.IsDir() {
				newTo := filepath.Join(to, string(game.BaseDir()))
				_ = os.MkdirAll(newTo, 0777)
				_ = os.Rename(sa, filepath.Join(newTo, "StreamingAssets"))
			} else if !tm.Mod().IsManuallyCreated {
				dir := filepath.Join(to, string(game.BaseDir()))
				if _, err = os.Stat(dir); err != nil {
					tm.Disable()
					enabler.DoneCallback(mods.Error, errors.New("unsupported nexus mod"))
					return
				}
			}
		}
	}

	for _, ti := range tis {
		files.AddModFiles(enabler, ti.DownloadFiles, func(result mods.Result, err ...error) {
			if result == mods.Error || result == mods.Cancel {
				tm.Disable()
			} else {
				tm.Enable()
				// Find any mods that are now disabled because all the files have been replaced by other mods
				for _, mod := range GetEnabledMods(enabler.Game) {
					if !managed.HasManagedFiles(enabler.Game, mod.ID()) {
						mod.Disable()
					}
				}
				_ = saveToJson()
			}
			enabler.DoneCallback(result, err...)
		})
	}
}
*/
/*
	func UpdateMod(game config.GameDef, tm mods.TrackedMod) (err error) {
		if tm.UpdatedMod() == nil {
			return errors.New("no update available")
		}

		if err = tm.Mod().Supports(game); err != nil {
			return
		}

		if tm.Enabled() {
			if err = DisableMod(game, tm); err != nil {
				return
			}
		}

		tm.SetMod(tm.UpdatedMod())
		if err = saveMoogle(tm); err != nil {
			return
		}

		tm.SetUpdatedMod(nil)
		return saveToJson()
	}
*/

/*func (b *enableBind) EnableMod() (err error) {
	var (
		tm  = b.tm
		tis []*mods.ToInstall
	)

	if len(b.tm.Mod().Configurations) > 0 {
		err = b.enableModWithConfig()
	} else {
		tis, err = mods.NewToInstallForMod(tm.Kind(), tm.Mod(), tm.Mod().AlwaysDownload)
		if err == nil {
			// Success
			err = b.enableMod(tis)
		}
	}
	return
}

func (b *enableBind) enableModWithConfig() (err error) {
	modPath := filepath.Join(config.Get().GetModsFullPath(state.CurrentGame), b.tm.ID().AsDir())
	if err = state.GetScreen(state.ConfigInstaller).(ci.ConfigInstaller).Setup(b.tm.Mod(), modPath, b.enableMod); err != nil {
		// Failed to set up config installer screen
		return
	}
	state.ShowScreen(state.ConfigInstaller)
	return
}

func (b *enableBind) enableMod(toInstall []*mods.ToInstall) (err error) {
	return managed.EnableMod(mods.NewModEnabler(state.CurrentGame, b.tm, toInstall, b.OnConflict, b.showWorking, func(result mods.Result, err ...error) {
		_ = b.Set(b.tm.Enabled())
		b.done(result, err...)
	}))
}

func (b *enableBind) DisableMod() error {
	return managed.DisableMod(state.CurrentGame, b.tm)
}

func (b *enableBind) OnConflict(conflicts []*mods.FileConflict, confirmationCallback mods.ConflictChoiceCallback) {
	f := widget.NewForm()
	for _, c := range conflicts {
		var name string
		if m, ok := managed.TryGetMod(state.CurrentGame, c.CurrentModID); !ok {
			name = string(c.CurrentModID)
		} else {
			name = m.DisplayName()
		}
		f.Items = append(f.Items, widget.NewFormItem(
			filepath.Base(c.File),
			widget.NewSelect([]string{name, string(b.tm.Mod().Name)}, c.OnChange)))
	}
	d := dialog.NewCustomConfirm("Conflicts", "ok", "cancel", container.NewVScroll(f), func(ok bool) {
		r := mods.Ok
		if !ok {
			r = mods.Cancel
		}
		confirmationCallback(r, conflicts)
	}, ui.Window)
	d.Resize(fyne.NewSize(400, 400))
	d.Show()
}
*/
