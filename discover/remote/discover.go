package remote

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"golang.org/x/sync/errgroup"
	"sync"
)

func GetMod(game config.GameDef, id mods.ModID, rebuildCache bool) (found bool, mod *mods.Mod, err error) {
	var result []*mods.Mod
	if result, err = GetMods(game, rebuildCache); err != nil {
		return
	}
	for _, mod = range result {
		if found = mod.ID() == id; found {
			return
		}
	}
	return
}

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

func GetMods(game config.GameDef, rebuildCache bool) (result []*mods.Mod, err error) {
	var (
		eg = errgroup.Group{}
		m  = sync.Mutex{}
	)

	// Get the mods from the remote sources
	for _, cl := range GetClients() {
		getMods(game, cl, &eg, &m, &result, rebuildCache)
	}
	if err = eg.Wait(); err != nil {
		return
	}
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

func getMods(game config.GameDef, c Client, eg *errgroup.Group, m *sync.Mutex, result *[]*mods.Mod, rebuildCache bool) {
	eg.Go(func() error {
		r, e := c.GetMods(game, rebuildCache)
		if e != nil {
			return e
		}
		m.Lock()
		*result = append(*result, r...)
		m.Unlock()
		return nil
	})
}
