package files

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
)

func AddModFiles(enabler *mods.ModEnabler, files []*mods.DownloadFiles, done mods.DoneCallback) {
	var (
		game    = enabler.Game
		mmf, ok = managed[game]
	)
	if !ok {
		mmf = &managedModsAndFiles{
			Mods: make(map[mods.ModID]*modFiles),
		}
		managed[game] = mmf
	}

	ResolveConflicts(enabler, mmf.Mods, files, func(result mods.Result, cr conflictResult, err ...error) {
		if result == mods.Cancel {
			done(result)
		} else if result == mods.Error {
			done(result, err...)
		} else {
			switch enabler.Game {
			case config.ChronoCross:
				// TODO Chrono Cross
			case config.BofIII, config.BofIV:
				// TODO BoF
			default:
				if e := addModFiles(enabler, mmf, files, cr); e != nil {
					done(mods.Error, e)
				} else {
					done(mods.Ok)
				}
			}
		}
	})
}

func RemoveModFiles(game config.Game, tm *mods.TrackedMod) (err error) {
	var (
		mmf, ok = managed[game]
		mf      *modFiles
	)
	if !ok {
		return fmt.Errorf("%s is not enabled", tm.Mod.Name)
	}
	if mf, ok = mmf.Mods[tm.GetModID()]; !ok {
		return fmt.Errorf("%s is not enabled", tm.Mod.Name)
	}

	switch game {
	case config.ChronoCross:
		// TODO Chrono Cross
	case config.BofIII, config.BofIV:
		// TODO BoF
	default:
		err = removeModFiles(mf, mmf, tm)
	}
	return
}
