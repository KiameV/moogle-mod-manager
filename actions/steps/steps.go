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
	ui "github.com/kiamev/moogle-mod-manager/ui/state"
	uiu "github.com/kiamev/moogle-mod-manager/ui/util"
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

func VerifyEnable(state *State) (result mods.Result, err error) {
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
	return
}

func VerifyDisable(state *State) (result mods.Result, err error) {
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
		uiu.ShowErrorLong(err)
		return
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
		}); err == nil {
			ui.ShowScreen(ui.ConfigInstaller)
			wg.Wait()
			time.Sleep(100 * time.Millisecond)
		}
		// Failed to set up config installer screen
		return
	}

	if len(state.ToInstall) == 0 {
		return mods.Error, errors.New("no files to install")
	}

	wg.Add(1)
	// Confirm Download
	if err = confirm.NewConfirmer(confirm.NewParams(state.Game, state.Mod, state.ToInstall)).Downloads(func(r mods.Result) {
		result = r
		wg.Done()
	}); err == nil {
		wg.Wait()
		time.Sleep(100 * time.Millisecond)
	}
	return
}

func Download(state *State) (result mods.Result, err error) {
	if err = downloads.Download(state.Game, state.Mod, state.ToInstall); err != nil {
		result = mods.Error
	} else {
		result = mods.Ok
	}
	return
}

func Extract(state *State) (result mods.Result, err error) {
	var (
		to       string
		override bool
		ef       []archive.ExtractedFile
	)
	for _, ti := range state.ToInstall {
		if to, err = filepath.Abs(*ti.Download.DownloadedArchiveLocation); err != nil {
			return
		}

		if state.Mod.InstallType(state.Game) == config.ImmediateDecompress {
			if to, err = config.Get().GetDir(state.Game, config.GameDirKind); err != nil {
				return
			}
			override = true
		} else {
			//override = false
			override = true
		}

		if ef, err = archive.Decompress(*ti.Download.DownloadedArchiveLocation, to, override); err != nil {
			return
		}

		if state.Mod.InstallType(state.Game) == config.ImmediateDecompress {
			installed := make([]string, len(ef))
			for i, f := range ef {
				installed[i] = f.File
			}
			files.SetFiles(state.Game, state.Mod.ID(), installed...)
		} else {
			state.ExtractedFiles = append(state.ExtractedFiles, Extracted{ti, files})
		}
	}
	return mods.Error, nil
}

func Conflicts(state *State) (result mods.Result, err error) {
	var (
		backupDir string
		rels      []string
	)
	if backupDir, err = config.Get().GetDir(state.Game, config.BackupDirKind); err != nil {
		return
	}
	for _, ef := range state.ExtractedFiles {
		if rels, err = ef.FilesRelative(); err != nil {

		}
	}
	for _, ti := range state.ToInstall {
		if ti.Conflicts, err = files.FindConflicts(state.Game); err != nil {
			return
		}
	}
}

func Install(state *State) (result mods.Result, err error) {
	return mods.Error, nil
}

func InstallMod(state *State) (result mods.Result, err error) {
	return mods.Error, nil
}

func UninstallMove(state *State) (mods.Result, error) {
	i := files.Files(state.Game, state.Mod.ID())
	for _, f := range i.Keys() {
		_ = os.Remove(f)
	}
	return mods.Error, nil
}

func RestoreBackups(state *State) (result mods.Result, err error) {
	return mods.Error, nil
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

func ShowWorkingDialog(_ *State) (result mods.Result, err error) {
	return mods.Working, nil
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
