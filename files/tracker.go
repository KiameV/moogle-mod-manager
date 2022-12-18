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
	}
	fileTracker struct {
		Files  collections.Set[string] `json:"files"`
		Backup collections.Set[string] `json:"backup"`
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
		mt = &modTracker{Mods: make(map[mods.ModID]*fileTracker)}
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
			Files:  collections.NewSet[string](),
			Backup: collections.NewSet[string](),
		}
		mt.Mods[modID] = ft
	}
	return ft
}

func Files(game config.GameDef, modID mods.ModID) collections.Set[string] {
	return modFiles(game, modID).Files
}

func Backups(game config.GameDef, modID mods.ModID) collections.Set[string] {
	return modFiles(game, modID).Backup
}

func HasBackup(game config.GameDef, file string) (modID mods.ModID, found bool) {
	var ft *fileTracker
	for modID, ft = range ModTracker(game).Mods {
		if ft.Backup.Contains(file) {
			return modID, true
		}
	}
	return
}

func SetFiles(game config.GameDef, modID mods.ModID, files ...string) {
	var (
		ft = modFiles(game, modID)
	)
	for _, f := range files {
		ft.Files.Set(f)
	}
	save()
}

func SetBackups(game config.GameDef, modID mods.ModID, backups ...string) {
	var (
		ft = modFiles(game, modID)
	)
	for _, f := range backups {
		ft.Backup.Set(f)
	}
	save()
}

func RemoveBackups(game config.GameDef, modID mods.ModID, backups ...string) {
	var (
		ft = modFiles(game, modID)
	)
	for _, bu := range backups {
		ft.Backup.Remove(bu)
	}
	save()
}

func save() {
	if err := util.SaveToFile(filepath.Join(config.PWD, file), tracker); err != nil {
		uu.ShowErrorLong(fmt.Errorf("failed to save file tracker: %v", err))
	}
}
