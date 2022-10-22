package files

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/conflict"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/managed"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/mover"
)

func AddModFiles(enabler *mods.ModEnabler, files []*mods.DownloadFiles, done mods.DoneCallback) {
	mmf := managed.GetModsWithManagedFiles(enabler.Game)
	conflict.ResolveConflicts(enabler, mmf.Mods, files, func(result mods.Result, cr conflict.Result, err ...error) {
		if result == mods.Cancel {
			done(result)
		} else if result == mods.Error {
			done(result, err...)
		} else {
			fm := mover.NewFileMover(enabler.Game)
			if e := fm.AddModFiles(enabler, mmf, files, cr); e != nil {
				done(mods.Error, e)
			} else {
				done(mods.Ok)
			}
		}
	})
}

func RemoveModFiles(game config.Game, tm *mods.TrackedMod) error {
	var (
		mmf    = managed.GetModsWithManagedFiles(game)
		mf, ok = mmf.Mods[tm.GetModID()]
	)
	if !ok {
		return fmt.Errorf("%s is not enabled", tm.Mod.Name)
	}
	return mover.NewFileMover(game).RemoveModFiles(mf, mmf, tm)
}
