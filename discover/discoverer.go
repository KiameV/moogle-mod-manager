package discover

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/discover/remote"
	"github.com/kiamev/moogle-mod-manager/discover/repo"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"golang.org/x/sync/errgroup"
)

var (
	utilLookup    = mods.NewModLookup()
	gameModLookup = make(map[config.GameID]mods.ModLookup[*mods.TrackedMod])
)

/* TODO REMOVE func GetMods(game *config.GameDef) (found []*mods.Mod, lookup mods.ModLookup, err error) {
	if lookup, err = GetModsAsLookup(game); err != nil {
		return
	}

	found = make([]*mods.Mod, 0, len(lookup))
	for _, m := range lookup {
		found = append(found, m)
	}
	return
}*/

func GetModsAsLookup(game *config.GameDef) (lookup mods.ModLookup, err error) {
	if game == nil {
		lookup = utilLookup
	} else {
		lookup = gameModLookup[game.ID]
	}
	if lookup.Len() > 0 {
		return
	}

	var (
		remoteMods []*mods.Mod
		repoMods   []*mods.Mod
		found      mods.IMod
		eg         errgroup.Group
		ok         bool
	)
	eg.Go(func() (e error) {
		remoteMods, e = remote.GetMods(*game)
		return
	})
	eg.Go(func() (e error) {
		repoMods, e = repo.NewGetter().GetMods(*state.CurrentGame)
		return
	})
	if err = eg.Wait(); err != nil {
		return
	}

	lookup = mods.NewModLookup()
	for _, m := range repoMods {
		if lookup.Has(m) {
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
		gameModLookup[game.ID] = lookup
	}
	return
}

func GetDisplayName(game config.GameDef, modID mods.ModID) (string, error) {
	lookup, err := GetModsAsLookup(&game)
	if err != nil {
		return "", err
	}
	if mod, ok := lookup.GetByID(modID); ok {
		return mod.Mod().Name, nil
	}
	return "", fmt.Errorf("mod [%s] not found", modID)
}
