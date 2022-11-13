package remote

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"golang.org/x/sync/errgroup"
	"sync"
)

func GetMods(game *config.Game) (result []*mods.Mod, err error) {
	var (
		eg = errgroup.Group{}
		m  = sync.Mutex{}
	)

	for _, c := range []discoverClient{
		newNexusClient(),
		newCurseForgeClient(),
	} {
		getMods(game, c, &eg, &m, &result)
	}
	err = eg.Wait()
	return
}

func getMods(game *config.Game, c discoverClient, eg *errgroup.Group, m *sync.Mutex, result *[]*mods.Mod) {
	eg.Go(func() error {
		r, e := c.GetMods(game)
		if e != nil {
			return e
		}
		m.Lock()
		defer m.Unlock()
		*result = append(*result, r...)
		return nil
	})
}
