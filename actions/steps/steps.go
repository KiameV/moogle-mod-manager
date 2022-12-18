package steps

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
)

type (
	State struct {
		Game       config.GameDef
		Mod        mods.TrackedMod
		Downloaded []string
		Extracted  []string
		ToInstall  []string
	}
	Step func(state *State) error
)

func NewState(game config.GameDef, mod mods.TrackedMod) *State {
	return &State{
		Game: game,
		Mod:  mod,
	}
}

func VerifyEnable(state *State) error {
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
					return fmt.Errorf("[%s] cannot be enabled because [%s] is enabled", tm.DisplayName(), mod.DisplayName())
				}
			}
		}
		if len(c.Requires) > 0 {
			for _, mc = range c.Requires {
				mod, found, enabled = managed.IsModEnabled(state.Game, mc.ModID())
				if !found {
					return fmt.Errorf("[%s] cannot be enabled because [%s] is not enabled", tm.DisplayName(), mc.ModID())
				} else if !enabled {
					return fmt.Errorf("[%s] cannot be enabled because [%s] is not enabled", tm.DisplayName(), mod.DisplayName())
				}
			}
		}
	}
	return nil
}

func VerifyDisable(state *State) error {
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
							return fmt.Errorf("[%s] cannot be disabled because [%s] is enabled", tm.DisplayName(), m.DisplayName())
						}
					}
				}
			}
		}
	}
	return nil
}

func Download(state *State) error {
	return nil
}

func Extract(state *State) error {
	return nil
}

func PreInstall(state *State) error {
	return nil
}

func Conflicts(state *State) error {
	return nil
}

func Install(state *State) error {
	return nil
}

func InstallMod(state *State) error {
	return nil
}

func UninstallMove(state *State) error {
	return nil
}

func RestoreBackups(state *State) error {
	return nil
}

func EnableMod(state *State) error {
	state.Mod.Enable()
	return managed.Save()
}

func DisableMod(state *State) error {
	state.Mod.Disable()
	return managed.Save()
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
