package actions

import (
	"errors"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/actions/steps"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"sync"
)

type (
	Action interface {
		Run() error
	}
	action struct {
		state         *steps.State
		steps         []steps.Step
		workingDialog WorkingDialog
	}
	WorkingDialog struct {
		Show func()
		Hide func()
	}
	Params struct {
		Game          config.GameDef
		Mod           mods.TrackedMod
		WorkingDialog WorkingDialog
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
		steps.PreDownload,
		steps.ShowWorkingDialog,
		steps.Download,
		steps.Extract,
		steps.PreInstall,
		steps.Conflicts,
		steps.Install,
		steps.EnableMod,
	}
	uninstallMoveSteps = []steps.Step{
		steps.VerifyDisable,
		steps.ShowWorkingDialog,
		steps.UninstallMove,
		steps.RestoreBackups,
		steps.DisableMod,
	}
	installImmediateDecompressSteps = []steps.Step{
		steps.VerifyEnable,
		steps.ShowWorkingDialog,
		steps.Download,
		steps.Extract,
		steps.EnableMod,
	}
	updateSteps = []steps.Step{
		steps.VerifyEnable,
		steps.ShowWorkingDialog,
		steps.DisableMod,
		steps.UninstallMove,
		steps.RestoreBackups,
		steps.InstallMod,
	}
	running = false
	mutex   = sync.Mutex{}
)

func New(kind ActionKind, params Params) (Action, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if running {
		return nil, errors.New("another action is running")
	}

	var (
		s   []steps.Step
		err error
	)
	switch kind {
	case Install:
		s, err = createInstallSteps(params.Game, params.Mod)
	case Uninstall:
		s, err = createUninstallSteps(params.Game, params.Mod)
	case Update:
		s, err = createUpdateSteps(params.Game, params.Mod)
	}
	return &action{
		state:         steps.NewState(params.Game, params.Mod),
		steps:         s,
		workingDialog: params.WorkingDialog,
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

func (a action) Run() (err error) {
	mutex.Lock()
	if running {
		err = errors.New("another action is running")
	} else {
		running = true
	}
	mutex.Unlock()
	if err != nil {
		return
	}

	go func() {
		a.run()
	}()
	return
}

func (a action) run() {
	defer func() {
		a.workingDialog.Hide()
		mutex.Lock()
		running = false
		mutex.Unlock()
	}()
	var (
		result mods.Result
		err    error
	)
	for _, step := range a.steps {
		if result, err = step(a.state); err != nil {
			util.ShowErrorLong(err)
			return
		} else if result == mods.Cancel {
			break
		} else if result == mods.Working {
			a.workingDialog.Show()
		}
	}
}

// TODO Update
// TODO Remove
