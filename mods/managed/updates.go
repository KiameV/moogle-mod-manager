package managed

import (
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"github.com/carwale/golibraries/workerpool"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	"github.com/kiamev/moogle-mod-manager/mods/nexus"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"strings"
	"sync"
)

const updateAvailableFlag = " *"

func CheckForUpdates(game config.Game, result func(err error)) {
	var (
		dispatcher = workerpool.NewDispatcher(
			fmt.Sprintf("Checker%d", game),
			workerpool.SetMaxWorkers(4))
		wg  = sync.WaitGroup{}
		ucs []updateChecker
	)
	for _, tm := range lookup[game].Mods {
		if tm.Mod.ModKind.Kind == mods.Hosted {
			wg.Add(1)
			h := &hostedUpdateChecker{tm: tm, wg: &wg}
			ucs = append(ucs, h)
			dispatcher.JobQueue <- h
		} else if tm.Mod.ModKind.Kind == mods.Nexus {
			wg.Add(1)
			n := &nexusUpdateChecker{tm: tm, wg: &wg}
			ucs = append(ucs, n)
			dispatcher.JobQueue <- n
		}
	}
	wg.Wait()
	for _, uc := range ucs {
		if uc.getError() != nil {
			result(uc.getError())
			return
		}
	}
	result(nil)
}

type updateChecker interface {
	getError() error
}

type hostedUpdateChecker struct {
	tm  *model.TrackedMod
	wg  *sync.WaitGroup
	err error
}

func (c *hostedUpdateChecker) Process() error {
	defer c.wg.Done()
	var mod mods.Mod
	for _, l := range c.tm.Mod.ModKind.Hosted.ModFileLinks {
		if b, err := browser.DownloadAsBytes(l); err == nil {
			if e := json.Unmarshal(b, &mod); e != nil && mod.ModKind != nil && mod.ModKind.Kind == mods.Hosted && mod.ModKind.Hosted != nil {
				continue
			}
		}
		if mod.ID == "" {
			dialog.ShowError(errors.New("Could not download remote version for "+c.tm.Mod.Name), state.Window)
			return nil
		}
		if isVersionNewer(c.tm.Mod.Version, mod.Version) {
			markForUpdate(c.tm, &mod)
		}
		break
	}
	return nil
}

func (c *hostedUpdateChecker) getError() error {
	return c.err
}

type nexusUpdateChecker struct {
	tm  *model.TrackedMod
	wg  *sync.WaitGroup
	err error
}

func (c *nexusUpdateChecker) Process() error {
	defer c.wg.Done()
	mod, err := nexus.GetModFromNexusForMod(c.tm.Mod)
	if err != nil {
		c.err = err
		return nil
	}
	if isVersionNewer(mod.Version, c.tm.Mod.Version) {
		markForUpdate(c.tm, mod)
	}
	return nil
}

func (c *nexusUpdateChecker) getError() error {
	return c.err
}

func isVersionNewer(new string, old string) bool {
	if new == old {
		return false
	}
	newSl := strings.Split(new, ".")
	oldSl := strings.Split(old, ".")
	for i := 0; i < len(newSl) && i < len(oldSl); i++ {
		if newSl[i] > oldSl[i] {
			return true
		}
	}
	if len(newSl) > len(oldSl) {
		return true
	}
	return false
}

func markForUpdate(tm *model.TrackedMod, mod *mods.Mod) {
	tm.UpdatedMod = mod
	tm.DisplayName = mod.Name + updateAvailableFlag
}
