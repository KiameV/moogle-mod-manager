package actions

import (
	"errors"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/actions/steps"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type (
	Action interface {
		Run() error
	}
	action struct {
		state *steps.State
		steps []steps.Step
	}
	ActionKind byte
)

const (
	Install ActionKind = iota
	Uninstall
	Update
)

var (
	installMoveSteps = []steps.Step{
		steps.VerifyEnable,
		steps.Download,
		steps.Extract,
		steps.PreInstall,
		steps.Conflicts,
		steps.Install,
		steps.EnableMod,
	}
	uninstallMoveSteps = []steps.Step{
		steps.VerifyDisable,
		steps.UninstallMove,
	}
	installImmediateDecompressSteps = []steps.Step{
		steps.VerifyEnable,
		steps.Download,
		steps.Extract,
		steps.EnableMod,
	}
	updateSteps = []steps.Step{
		steps.VerifyEnable,
		steps.DisableMod,
		steps.InstallMod,
	}
)

func New(kind ActionKind, game config.GameDef, tm mods.TrackedMod) (Action, error) {
	var (
		s   []steps.Step
		err error
	)
	switch kind {
	case Install:
		s, err = createInstallSteps(game, tm)
	case Uninstall:
		s, err = createUninstallSteps(game, tm)
	case Update:
		s, err = createUpdateSteps(game, tm)
	}
	return &action{
		state: steps.NewState(game, tm),
		steps: s,
	}, err
}

func createInstallSteps(game config.GameDef, tm mods.TrackedMod) (s []steps.Step, err error) {
	switch tm.InstallType(game) {
	case config.Move:
		s = installMoveSteps
	case config.ImmediateDecompress:
		s = installImmediateDecompressSteps
	case config.MoveToArchive:
		err = errors.New("not implemented")
	default:
		err = fmt.Errorf("unknown install %s for mod %s", tm.InstallType(game), tm.Mod().Name)
	}
	return
}

func createUninstallSteps(game config.GameDef, tm mods.TrackedMod) (s []steps.Step, err error) {
	switch tm.InstallType(game) {
	case config.Move, config.ImmediateDecompress:
		s = uninstallMoveSteps
	case config.MoveToArchive:
		err = errors.New("not implemented")
	default:
		err = fmt.Errorf("unknown install %s for mod %s", tm.InstallType(game), tm.Mod().Name)
	}
	return
}

func createUpdateSteps(game config.GameDef, tm mods.TrackedMod) (s []steps.Step, err error) {
	switch tm.InstallType(game) {
	case config.Move, config.ImmediateDecompress:
		s = updateSteps
	case config.MoveToArchive:
		err = errors.New("not implemented")
	default:
		err = fmt.Errorf("unknown install %s for mod %s", tm.InstallType(game), tm.Mod().Name)
	}
	return
}

func (a action) Run() error {
	for _, step := range a.steps {
		if err := step(a.state); err != nil {
			return err
		}
	}
	return nil
}

// TODO Update
// TODO Remove

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
