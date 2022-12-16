package managed

import (
	"errors"
	"fmt"
	"github.com/carwale/golibraries/workerpool"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/discover/remote"
	"github.com/kiamev/moogle-mod-manager/discover/repo"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"strings"
	"sync"
)

func CheckForUpdates(game config.GameDef, result func(err error)) {
	var (
		dispatcher = workerpool.NewDispatcher(
			fmt.Sprintf("Checker%d", game),
			workerpool.SetMaxWorkers(4))
		wg  = sync.WaitGroup{}
		ucs []updateChecker
	)

	if err := repo.NewGetter(repo.Read).Pull(); err != nil {
		result(err)
		return
	}

	for _, tm := range lookup.GetMods(game) {
		switch tm.Kind() {
		case mods.Hosted:
			wg.Add(1)
			h := &hostedUpdateChecker{tm: tm, wg: &wg}
			ucs = append(ucs, h)
			dispatcher.JobQueue <- h
		case mods.Nexus:
			wg.Add(1)
			n := &remoteUpdateChecker{tm: tm, wg: &wg, client: remote.NewNexusClient()}
			ucs = append(ucs, n)
			dispatcher.JobQueue <- n
		case mods.CurseForge:
			wg.Add(1)
			n := &remoteUpdateChecker{tm: tm, wg: &wg, client: remote.NewCurseForgeClient()}
			ucs = append(ucs, n)
			dispatcher.JobQueue <- n
		default:
			result(fmt.Errorf("unknown mod kind %s", tm.Kind()))
			return
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
	tm  mods.TrackedMod
	wg  *sync.WaitGroup
	err error
}

func (c *hostedUpdateChecker) Process() error {
	defer c.wg.Done()

	remoteMod, err := repo.NewGetter(repo.Read).GetMod(c.tm.Mod())
	if err != nil {
		util.ShowErrorLong(err)
		return nil
	}

	if remoteMod.ID() != c.tm.ID() {
		util.ShowErrorLong(errors.New("Could not download remote version for " + c.tm.DisplayName()))
		return nil
	}
	if isVersionNewer(c.tm.Mod().Version, remoteMod.Version) {
		markForUpdate(c.tm, remoteMod)
	}
	return nil
}

func (c *hostedUpdateChecker) getError() error {
	return c.err
}

type remoteUpdateChecker struct {
	tm     mods.TrackedMod
	wg     *sync.WaitGroup
	client remote.Client
	err    error
}

func (c *remoteUpdateChecker) Process() error {
	defer c.wg.Done()
	_, mod, err := c.client.GetFromMod(c.tm.Mod())
	if err != nil {
		c.err = err
		return nil
	}
	if isVersionNewer(mod.Version, c.tm.Mod().Version) {
		markForUpdate(c.tm, mod)
	}
	return nil
}

func (c *remoteUpdateChecker) getError() error {
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
	return len(newSl) > len(oldSl)
}

func markForUpdate(tm mods.TrackedMod, mod *mods.Mod) {
	tm.SetUpdatedMod(mods.NewModForVersion(tm.Mod(), mod))
}
