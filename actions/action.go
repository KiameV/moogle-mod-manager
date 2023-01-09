package actions

import (
	"errors"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/actions/steps"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	"github.com/kiamev/moogle-mod-manager/ui/util/working"
	"sync"
	"time"
)

type (
	Done   func(result Result)
	Action interface {
		Run() error
		Close()
	}
	Result struct {
		Status       mods.Result
		Err          error
		RequiredMods []mods.TrackedMod
	}
	action struct {
		done             Done
		state            *steps.State
		steps            []steps.Step
		isInternalAction bool
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
	updateMoveSteps = []steps.Step{
		steps.VerifyEnable,
		steps.ShowWorkingDialog,
		steps.DisableMod,
		steps.UninstallMove,
	}
	running = false
	mutex   = sync.Mutex{}
)

func init() {
	updateMoveSteps = append(updateMoveSteps, installMoveSteps...)
}

func New(kind ActionKind, game config.GameDef, mod mods.TrackedMod, done Done) (Action, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if running {
		return nil, errors.New("Another action is running. Please wait for the current mod to finish installing or uninstalling.")
	}
	a, err := newAction(kind, game, mod, done)
	if err != nil {
		return nil, err
	}
	a.isInternalAction = false
	return a, nil
}

func newAction(kind ActionKind, game config.GameDef, mod mods.TrackedMod, done Done) (*action, error) {
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
		done:             done,
		state:            steps.NewState(game, mod),
		steps:            s,
		isInternalAction: true,
	}, err
}

func createInstallSteps(game config.GameDef, tm mods.TrackedMod) (s []steps.Step, err error) {
	switch tm.InstallType(game) {
	case config.Move, config.MoveToArchive:
		s = installMoveSteps
	case config.ImmediateDecompress:
		s = installImmediateDecompressSteps
	default:
		err = fmt.Errorf("unknown install %s for mod %s", tm.InstallType(game), tm.Mod().Name)
	}
	return
}

func createUninstallSteps(game config.GameDef, tm mods.TrackedMod) (s []steps.Step, err error) {
	switch tm.InstallType(game) {
	case config.Move, config.ImmediateDecompress, config.MoveToArchive:
		s = uninstallMoveSteps
	default:
		err = fmt.Errorf("unknown install %s for mod %s", tm.InstallType(game), tm.Mod().Name)
	}
	return
}

func createUpdateSteps(game config.GameDef, tm mods.TrackedMod) (s []steps.Step, err error) {
	switch tm.InstallType(game) {
	case config.Move, config.ImmediateDecompress, config.MoveToArchive:
		s = updateMoveSteps
	default:
		err = fmt.Errorf("unknown install %s for mod %s", tm.InstallType(game), tm.Mod().Name)
	}
	return
}

func (a *action) Run() (err error) {
	defer a.Close()
	if !a.isInternalAction {
		mutex.Lock()
		if running {
			err = errors.New("Another action is running. Please wait for the current mod to finish installing or uninstalling.")
		} else {
			running = true
		}
		mutex.Unlock()
		if err != nil {
			return
		}
	}
	go func() {
		a.run()
	}()
	return
}

func (a *action) run() {
	var (
		result mods.Result
		err    error
	)
	defer func() {
		if !a.isInternalAction {
			working.HideDialog()
			mutex.Lock()
			running = false
			mutex.Unlock()
		}
		if a.done != nil {
			go func() {
				time.Sleep(100 * time.Millisecond)
				a.done(a.newResult(result, err))
			}()
		}
	}()
	for i := 0; i < len(a.steps); i++ {
		if result, err = a.steps[i](a.state); err != nil {
			return
		} else if result == mods.Cancel {
			break
		} else if result == mods.Working {
			working.ShowDialog()
		} else if result == mods.Repeat {
			i--
			if a.state.Requires != nil {
				if result, err = installRequiredMod(a.state); result == mods.Cancel || result == mods.Error || err != nil {
					return
				}
			} else {
				panic("repeat step without required mod")
			}
		} else if result != mods.Ok {
			result = mods.Error
			err = fmt.Errorf("unknown result %d", result)
			return
		}
	}
}

func installRequiredMod(state *steps.State) (result mods.Result, err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(result *mods.Result, err *error) {
		var (
			a  Action
			tm mods.TrackedMod
			e  error
		)
		if tm, e = managed.AddMod(state.Game, state.Requires); e != nil {
			*result = mods.Error
			*err = e
			return
		}
		if a, e = newAction(Install, state.Game, tm, func(r Result) {
			// Done running install
			*result = r.Status
			*err = r.Err
			if r.Status == mods.Ok {
				state.Added = append(state.Added, tm)
			}
			wg.Done()
		}); e != nil {
			*result = mods.Error
			*err = e
			return
		}
		if e = a.Run(); e != nil {
			*result = mods.Error
			*err = e
			return
		}
	}(&result, &err)
	if result == mods.Cancel || result == mods.Error {
		return
	}
	wg.Wait()
	return
}

func (a *action) newResult(r mods.Result, err error) Result {
	if err != nil {
		r = mods.Error
	}
	return Result{
		Status:       r,
		Err:          err,
		RequiredMods: a.state.Added,
	}
}

func (a *action) Close() {
	a.state.Close()
}

// TODO Update
// TODO Remove
