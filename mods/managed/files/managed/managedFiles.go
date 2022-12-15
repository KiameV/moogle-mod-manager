package managed

import (
	"encoding/json"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"io/ioutil"
	"path/filepath"
)

const (
	managedXmlName = "managed.json"
)

var (
	managed = &GameMods{Games: make(map[config.GameID]*ModsAndFiles)}
)

type (
	ModsAndFiles struct {
		Mods map[mods.ModID]*ModFiles
	}
	GameMods struct {
		Games map[config.GameID]*ModsAndFiles
	}
	ModFiles struct {
		BackedUpFiles map[string]*mods.ModFile
		MovedFiles    map[string]*mods.ModFile
	}
)

func (mg *GameMods) Get(game config.GameDef) (mmf *ModsAndFiles, ok bool) {
	mmf, ok = mg.Games[game.ID()]
	return
}

func (mg *GameMods) Set(game config.GameDef, mmf *ModsAndFiles) {
	mg.Games[game.ID()] = mmf
}

func InitializeManagedFiles() error {
	b, err := ioutil.ReadFile(filepath.Join(config.PWD, managedXmlName))
	if err != nil {
		return nil
	}
	return json.Unmarshal(b, &managed)
}

func HasManagedFiles(game config.GameDef, modID mods.ModID) bool {
	var (
		mmf *ModsAndFiles
		ok  bool
		mf  *ModFiles
	)
	if mmf, ok = managed.Get(game); !ok {
		return false
	}
	if mf, ok = mmf.Mods[modID]; !ok {
		return false
	}
	return len(mf.MovedFiles) > 0
}

func GetModsWithManagedFiles(game config.GameDef) (mmf *ModsAndFiles) {
	var ok bool
	if mmf, ok = managed.Get(game); !ok {
		mmf = &ModsAndFiles{
			Mods: make(map[mods.ModID]*ModFiles),
		}
		managed.Set(game, mmf)
	}
	return
}

func GetManagedFiles(game config.GameDef, modID mods.ModID) (mf *ModFiles, hasManaged bool) {
	mf, hasManaged = GetModsWithManagedFiles(game).Mods[modID]
	return
}

func SaveManagedJson() error {
	b, err := json.MarshalIndent(managed, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(config.PWD, managedXmlName), b, 0777)
}
