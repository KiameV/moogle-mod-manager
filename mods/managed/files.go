package managed

import (
	"encoding/json"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	"io/ioutil"
)

const (
	managedXmlName = "managed.json"
)

var (
	managed = make(map[config.Game]*managedModsAndFiles)
)

type managedModsAndFiles struct {
	Mods          map[string]modFiles
	ReplacedFiles map[string]bool
}

type modFiles struct {
	Files []mods.ModFile
}

func AddModFiles(game config.Game, tm *model.TrackedMod, files []*mods.DownloadFiles) error {
	mmf, ok := managed[game]
	if !ok {
		mmf = &managedModsAndFiles{
			Mods:          make(map[string]modFiles),
			ReplacedFiles: make(map[string]bool),
		}
		managed[game] = mmf
	}
	/*
		for _, mf := range mmf.Mods {
			if modID == mf.ModID {
				return fmt.Errorf("%s is already enabled", modID)
			}
		}

		if collisions := detectCollisions(m.AllFiles, files); len(collisions) > 0 {
			return fmt.Errorf("cannot enable mod as these files would collide: %s", strings.Join(collisions, ", "))
		}

		m.Mods = append(m.Mods, modFiles{ModID: modID, Files: files})
		for _, f := range files {
			m.AllFiles[f] = true
		}*/
	return saveManagedJson()
}

func RemoveModFiles(game config.Game, modID string) error {
	/*m, ok := managed[game]
	if !ok {
		return nil
	}
	for i, mf := range m.Mods {
		if modID == mf.ModID {
			m.Mods[i] = m.Mods[len(m.Mods)-1]
			m.Mods = m.Mods[:len(m.Mods)-1]
			if err := io.RevertMoveFiles(mf.Files, game); err != nil {
				return err
			}
			for _, f := range mf.Files {
				delete(m.AllFiles, f)
			}
			break
		}
	}*/
	return saveManagedJson()
}

func detectCollisions(managedFiles map[string]bool, modFiles []string) (collisions []string) {
	var found bool
	for _, f := range modFiles {
		if _, found = managedFiles[f]; found {
			collisions = append(collisions, f)
		}
	}
	return
}

func saveManagedJson() error {
	b, err := json.MarshalIndent(managed, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(managedXmlName, b, 0777)
}
