package repo

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/util"
	"path/filepath"
)

var lookup = make([]*discoveredMods, 6)

type discoveredMods struct {
	Mods      []*mods.Mod
	Overrides []*mods.Override
}

func GetMods(game config.Game) ([]*mods.Mod, []*mods.Override, error) {
	var (
		d             = lookup[game]
		modFiles      []string
		overrideFiles []string
		r             repo
		f             string
		err           error
	)
	if d == nil {
		d = &discoveredMods{}
		lookup[game] = d

		if modFiles, overrideFiles, err = r.GetMods(game); err != nil {
			return nil, nil, err
		}
		for _, f = range modFiles {
			var mod mods.Mod
			if err = util.LoadFromFile(f, &mod); err != nil {
				// TODO log error
				continue
			}
			d.Mods = append(d.Mods, &mod)
		}
		for _, f = range overrideFiles {
			var override mods.Override
			if err = util.LoadFromFile(f, &override); err != nil {
				// TODO log error
				continue
			}
			d.Overrides = append(d.Overrides, &override)
		}
	}
	if d.Mods, err = getNexusMods(game, d.Mods, d.Overrides); err != nil {
		return nil, nil, err
	}
	return d.Mods, d.Overrides, nil
}

func repoDir() string {
	return filepath.Join(config.PWD, "remote")
}

func repoGameDir(game config.Game) string {
	return filepath.Join(repoDir(), config.String(game))
}
