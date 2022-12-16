package remote

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"golang.org/x/sync/errgroup"
	"sync"
)

func GetFromUrl(kind mods.Kind, url string) (bool, *mods.Mod, error) {
	var c Client
	switch kind {
	case mods.CurseForge:
		c = NewCurseForgeClient()
	case mods.Nexus:
		c = NewNexusClient()
	default:
		return false, nil, fmt.Errorf("invalid kind to GetFromUrl %v", kind)
	}
	return c.GetFromUrl(url)
}

func GetMods(game config.GameDef) (result []*mods.Mod, err error) {
	var (
		eg = errgroup.Group{}
		m  = sync.Mutex{}
	)

	for _, c := range GetClients() {
		getMods(game, c, &eg, &m, &result)
	}
	err = eg.Wait()
	return
}

func GetClients() []Client {
	return []Client{
		NewNexusClient(),
		NewCurseForgeClient(),
	}
	/*var c []Client
	if config.GetSecrets().NexusApiKey != "" {
		c = append(c, NewNexusClient())
	}
	if config.GetSecrets().CfApiKey != "" {
		NewCurseForgeClient()
	}
	return c*/
}

func getMods(game config.GameDef, c Client, eg *errgroup.Group, m *sync.Mutex, result *[]*mods.Mod) {
	eg.Go(func() error {
		r, e := c.GetMods(game)
		if e != nil {
			return e
		}
		m.Lock()
		*result = append(*result, r...)
		m.Unlock()
		return nil
	})
}
