package discover

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/discover/remote"
	"github.com/kiamev/moogle-mod-manager/discover/repo"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"golang.org/x/sync/errgroup"
)

var (
	utilLookup    = make(map[string]*mods.Mod)
	gameModLookup = [6]map[string]*mods.Mod{}
)

func GetMods(game *config.Game) (found []*mods.Mod, lookup map[string]*mods.Mod, err error) {
	if lookup, err = GetModsAsLookup(game); err != nil {
		return
	}

	found = make([]*mods.Mod, 0, len(lookup))
	for _, m := range lookup {
		found = append(found, m)
	}
	return
}

func GetModsAsLookup(game *config.Game) (lookup map[string]*mods.Mod, err error) {
	if game == nil {
		lookup = utilLookup
	} else {
		lookup = gameModLookup[*game]
	}
	if len(lookup) > 0 {
		return
	}

	var (
		remoteMods []*mods.Mod
		repoMods   []*mods.Mod
		found      *mods.Mod
		eg         errgroup.Group
		ok         bool
	)
	eg.Go(func() (e error) {
		remoteMods, e = remote.GetNexusMods(game)
		return
	})
	eg.Go(func() (e error) {
		repoMods, e = repo.NewGetter().GetMods(*state.CurrentGame)
		return
	})
	if err = eg.Wait(); err != nil {
		return
	}

	lookup = make(map[string]*mods.Mod)
	for _, m := range repoMods {
		if _, ok = lookup[m.ModUniqueID(*state.CurrentGame)]; !ok {
			lookup[m.ModUniqueID(*state.CurrentGame)] = m
		}
	}
	for _, m := range remoteMods {
		if found, ok = lookup[m.ModUniqueID(*state.CurrentGame)]; !ok {
			lookup[m.ModUniqueID(*state.CurrentGame)] = m
		} else {
			found.Merge(m)
		}
	}
	if game == nil {
		utilLookup = lookup
	} else {
		gameModLookup[*game] = lookup
	}
	return
}