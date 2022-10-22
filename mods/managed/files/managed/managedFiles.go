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
	managed = make(map[config.Game]*ManagedModsAndFiles)
)

type (
	ManagedModsAndFiles struct {
		Mods map[mods.ModID]*ModFiles
	}
	ModFiles struct {
		BackedUpFiles map[string]*mods.ModFile
		MovedFiles    map[string]*mods.ModFile
	}
)

func InitializeManagedFiles() error {
	b, err := ioutil.ReadFile(filepath.Join(config.PWD, managedXmlName))
	if err != nil {
		return nil
	}
	return json.Unmarshal(b, &managed)
}

func HasManagedFiles(game config.Game, modID mods.ModID) bool {
	var (
		mmf *ManagedModsAndFiles
		ok  bool
		mf  *ModFiles
	)
	if mmf, ok = managed[game]; !ok {
		return false
	}
	if mf, ok = mmf.Mods[modID]; !ok {
		return false
	}
	return len(mf.MovedFiles) > 0
}

func GetModsWithManagedFiles(game config.Game) (mmf *ManagedModsAndFiles) {
	var ok bool
	if mmf, ok = managed[game]; !ok {
		mmf = &ManagedModsAndFiles{
			Mods: make(map[mods.ModID]*ModFiles),
		}
	}
	return
}

func GetManagedFiles(game config.Game, modID mods.ModID) (mf *ModFiles, hasManaged bool) {
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
