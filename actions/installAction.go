package actions

import (
	"errors"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/actions/steps"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"github.com/kiamev/moogle-mod-manager/ui/util/working"
	"sync"
	"time"
)

type (
	Done   func()
	Action interface {
		Run() error
	}
	action struct {
		done  Done
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
		steps.PreDownload,
		steps.ShowWorkingDialog,
		steps.Download,
		steps.Extract,
		steps.Conflicts,
		steps.Install,
		steps.EnableMod,
		steps.PostInstall,
	}
	uninstallMoveSteps = []steps.Step{
		steps.VerifyDisable,
		steps.ShowWorkingDialog,
		steps.UninstallMove,
		steps.DisableMod,
	}
	installImmediateDecompressSteps = []steps.Step{
		steps.VerifyEnable,
		steps.PreDownload,
		steps.ShowWorkingDialog,
		steps.Download,
		steps.Extract,
		steps.EnableMod,
		steps.PostInstall,
	}
	updateSteps = []steps.Step{
		steps.VerifyEnable,
		steps.ShowWorkingDialog,
		steps.DisableMod,
		steps.UninstallMove,
		steps.Install,
		steps.PostInstall,
	}
	running = false
	mutex   = sync.Mutex{}
)

func New(kind ActionKind, game config.GameDef, mod mods.TrackedMod, done Done) (Action, error) {
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
		s, err = createInstallSteps(game, mod)
	case Uninstall:
		s, err = createUninstallSteps(game, mod)
	case Update:
		s, err = createUpdateSteps(game, mod)
	}
	return &action{
		done:  done,
		state: steps.NewState(game, mod),
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
		working.HideDialog()
		mutex.Lock()
		running = false
		mutex.Unlock()
		if a.done != nil {
			go func() {
				time.Sleep(100 * time.Millisecond)
				a.done()
			}()
		}
	}()
	var (
		result mods.Result
		err    error
	)
	for _, step := range a.steps {
		if result, err = step(a.state); err != nil {
			working.HideDialog()
			util.ShowErrorLong(err)
			return
		} else if result == mods.Cancel {
			break
		} else if result == mods.Working {
			working.ShowDialog()
		}
	}
}

// TODO Update
// TODO Remove
