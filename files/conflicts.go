package files

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
)

type (
	Conflict struct {
		Owner *mods.Mod
		File  string
	}
)

func (c *Conflict) FindConflicts(game config.GameDef, files []string) (conflicts []Conflict) {
	var (
		owner mods.ModID
		tm    mods.TrackedMod
		found bool
	)
	for _, f := range files {
		if owner, found = HasBackup(game, f); found {
			if tm, found = managed.TryGetMod(game, owner); found {
				conflicts = append(conflicts, Conflict{
					Owner: tm.Mod(),
					File:  f,
				})
			}
		}
	}
	return
}
