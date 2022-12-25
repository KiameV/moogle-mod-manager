package remote

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"golang.org/x/sync/errgroup"
	"sync"
)

var cache = make(map[config.GameID]map[mods.ModID]*mods.Mod)

func GetMod(game config.GameDef, id mods.ModID, rebuildCache bool) (found bool, mod *mods.Mod, err error) {
	if _, err = GetMods(game, rebuildCache); err != nil {
		return
	}
	if c, ok := cache[game.ID()]; ok {
		mod, found = c[id]
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
		c  = cache[game.ID()]
	)

	if c != nil && len(c) > 0 && !rebuildCache {
		// Use the cache
		result = make([]*mods.Mod, 0, len(cache))
		for _, mod := range c {
			result = append(result, mod)
		}
		return
	}

	// Get the mods from the remote sources
	for _, cl := range GetClients() {
		getMods(game, cl, &eg, &m, &result)
	}
	if err = eg.Wait(); err != nil {
		return
	}

	// Build the cache
	l := make(map[mods.ModID]*mods.Mod)
	for _, mod := range result {
		l[mod.ID()] = mod
	}
	cache[game.ID()] = l
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
