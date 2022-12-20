package discover

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/discover/remote"
	"github.com/kiamev/moogle-mod-manager/discover/repo"
	"github.com/kiamev/moogle-mod-manager/mods"
	"golang.org/x/sync/errgroup"
)

type (
	gameMods struct {
		lookup map[config.GameID]mods.ModLookup[*mods.Mod]
	}
)

func (m gameMods) Get(game config.GameDef) (l mods.ModLookup[*mods.Mod], found bool) {
	l, found = m.lookup[game.ID()]
	return
}

func (m gameMods) Set(game config.GameDef, lookup mods.ModLookup[*mods.Mod]) {
	m.lookup[game.ID()] = lookup
}

var (
	utilLookup    = mods.NewModLookup[*mods.Mod]()
	gameModLookup = &gameMods{lookup: make(map[config.GameID]mods.ModLookup[*mods.Mod])}
)

/* TODO REMOVE func GetMods(game config.GameDef) (found []*mods.Mod, lookup mods.ModLookup, err error) {
	if lookup, err = GetModsAsLookup(game); err != nil {
		return
	}

	found = make([]*mods.Mod, 0, len(lookup))
	for _, m := range lookup {
		found = append(found, m)
	}
	return
}*/

func GetModsAsLookup(game config.GameDef) (lookup mods.ModLookup[*mods.Mod], err error) {
	var (
		remoteMods []*mods.Mod
		repoMods   []*mods.Mod
		found      *mods.Mod
		eg         errgroup.Group
		ok         bool
	)

	/* TODO is this cache needed?
	if game == nil {
		lookup = utilLookup
		ok = true
	} else {
		lookup, ok = gameModLookup.Get(game)
	}
	if lookup != nil && lookup.Len() > 0 && ok {
		return
	}*/

	if game != nil {
		eg.Go(func() (e error) {
			remoteMods, e = remote.GetMods(game)
			return
		})
		eg.Go(func() (e error) {
			repoMods, e = repo.NewGetter(repo.Read).GetMods(game)
			return
		})
		if err = eg.Wait(); err != nil {
			return
		}
	} else { // utilities
		if repoMods, err = repo.NewGetter(repo.Read).GetUtilities(); err != nil {
			return
		}
	}

	lookup = mods.NewModLookup[*mods.Mod]()
	for _, m := range repoMods {
		if !lookup.Has(m) {
			lookup.Set(m)
		}
	}
	for _, m := range remoteMods {
		if found, ok = lookup.Get(m); !ok {
			lookup.Set(m)
		} else {
			found.Mod().Merge(*m)
		}
	}
	if game == nil {
		utilLookup = lookup
	} else {
		gameModLookup.Set(game, lookup)
	}
	return
}

func GetDisplayName(game config.GameDef, modID mods.ModID) (string, error) {
	lookup, err := GetModsAsLookup(game)
	if err != nil {
		return "", err
	}
	if mod, ok := lookup.GetByID(modID); ok {
		return string(mod.Mod().Name), nil
	}
	return "", fmt.Errorf("mod [%s] not found", modID)
}
