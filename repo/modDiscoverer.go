package repo

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/util"
)

var lookup = make([]*discoveredMods, 6)

type discoveredMods struct {
	Mods []*mods.Mod `json:"-"`
}

func GetMods(game config.Game) ([]*mods.Mod, error) {
	var (
		d     = lookup[game]
		files []string
		r     repo
		err   error
	)
	if d == nil {
		d = &discoveredMods{}
		lookup[game] = d

		if files, err = r.GetMods(game); err != nil {
			return nil, err
		}
		for _, f := range files {
			var mod mods.Mod
			if err = util.LoadFromFile(f, &mod); err != nil {
				// TODO log error
				continue
			}
			d.Mods = append(d.Mods, &mod)
		}
	}
	return d.Mods, nil
}
