package files

import (
	"errors"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/collections"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	uu "github.com/kiamev/moogle-mod-manager/ui/util"
	"github.com/kiamev/moogle-mod-manager/util"
	"path/filepath"
	"syscall"
)

const file = "filetracker.json"

type (
	gameTracker struct {
		Games map[config.GameID]*modTracker `json:"games"`
	}
	modTracker struct {
		Mods map[mods.ModID]*fileTracker `json:"mods"`
		//Backups collections.Set[string]     `json:"backup"`
	}
	fileTracker struct {
		Files collections.Set[string] `json:"files"`
	}
)

var tracker = &gameTracker{Games: make(map[config.GameID]*modTracker)}

func Initialize() error {
	if err := util.LoadFromFile(filepath.Join(config.PWD, file), tracker); err != nil {
		if !errors.Is(err, syscall.ERROR_FILE_NOT_FOUND) {
			return fmt.Errorf("failed to load file tracker: %v", err)
		}
	}
	return nil
}

func ModTracker(game config.GameDef) *modTracker {
	mt, ok := tracker.Games[game.ID()]
	if !ok {
		mt = &modTracker{
			Mods: make(map[mods.ModID]*fileTracker),
			//Backups: collections.NewSet[string](),
		}
		tracker.Games[game.ID()] = mt
	}
	return mt
}

func modFiles(game config.GameDef, modID mods.ModID) *fileTracker {
	var (
		mt     = ModTracker(game)
		ft, ok = mt.Mods[modID]
	)
	if !ok {
		ft = &fileTracker{
			Files: collections.NewSet[string](),
		}
		mt.Mods[modID] = ft
	}
	return ft
}

func Files(game config.GameDef, modID mods.ModID) collections.Set[string] {
	return modFiles(game, modID).Files
}

func EmptyMods(game config.GameDef) (result []mods.ModID) {
	for id, ft := range ModTracker(game).Mods {
		if ft.Files.Len() == 0 {
			result = append(result, id)
		}
	}
	return
}

//func Backups(game config.GameDef) collections.Set[string] {
//	return ModTracker(game).Backups
//}

func HasFile(game config.GameDef, file string) (modID mods.ModID, found bool) {
	var ft *fileTracker
	for modID, ft = range ModTracker(game).Mods {
		if ft.Files.Contains(file) {
			return modID, true
		}
	}
	return
}

//func HasBackup(game config.GameDef, file string) bool {
//	return ModTracker(game).Backups.Contains(file)
//}

func SetFiles(game config.GameDef, modID mods.ModID, files ...string) {
	var (
		ft = modFiles(game, modID)
	)
	for _, f := range files {
		ft.Files.Set(f)
	}
	tracker.save()
}

/*func SetBackups(game config.GameDef, backups ...string) {
	mt := ModTracker(game)
	for _, f := range backups {
		mt.Backups.Set(f)
	}
	save()
}*/

/*func RemoveBackups(game config.GameDef, backups ...string) {
	mt := ModTracker(game)
	for _, bu := range backups {
		mt.Backups.Remove(bu)
	}
	save()
}*/

func RemoveFiles(game config.GameDef, modID mods.ModID, files ...string) {
	for _, f := range files {
		modFiles(game, modID).Files.Remove(f)
	}
	tracker.save()
}

func (t *gameTracker) save() {
	if err := util.SaveToFile(filepath.Join(config.PWD, file), t); err != nil {
		uu.ShowErrorLong(fmt.Errorf("failed to save file tracker: %v", err))
	}
}
