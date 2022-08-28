package repo

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/util"
)

var lookup = make([]*discoveredMods, 6)

type discoveredMods struct {
	Mods []*mods.Mod
}

func GetMods(game config.Game) ([]*mods.Mod, error) {
	var (
		d        = lookup[game]
		modFiles []string
		r        repo
		f        string
		err      error
	)
	if d == nil {
		d = &discoveredMods{}
		if modFiles, err = r.GetMods(game); err != nil {
			return nil, err
		}
		for _, f = range modFiles {
			var mod mods.Mod
			if err = util.LoadFromFile(f, &mod); err != nil {
				// TODO log error
				continue
			}
			d.Mods = append(d.Mods, &mod)
		}

		lookup[game] = d
	}
	if d.Mods, err = getNexusMods(game, d.Mods); err != nil {
		return nil, err
	}
	return d.Mods, nil
}
