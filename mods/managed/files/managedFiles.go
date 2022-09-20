package files

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
	managed = make(map[config.Game]*managedModsAndFiles)
)

type action bool

const (
	duplicate action = false
	cut       action = true
)

func InitializeManagedFiles() error {
	b, err := ioutil.ReadFile(filepath.Join(config.PWD, managedXmlName))
	if err != nil {
		return nil
	}
	return json.Unmarshal(b, &managed)
}

type managedModsAndFiles struct {
	Mods map[mods.ModID]*modFiles
}

type modFiles struct {
	BackedUpFiles map[string]*mods.ModFile
	MovedFiles    map[string]*mods.ModFile
}

func saveManagedJson() error {
	b, err := json.MarshalIndent(managed, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(config.PWD, managedXmlName), b, 0777)
}
