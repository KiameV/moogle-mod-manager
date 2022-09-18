package remote

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/nexus"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
)

func getDir(game config.Game) string {
	return filepath.Join(config.PWD, "remote", config.String(game), "nexus")
}

func GetNexusMods(game *config.Game) (result []*mods.Mod, err error) {
	if game == nil {
		return
	}
	dir := getDir(*game)
	_ = os.MkdirAll(dir, 0777)
	if err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "mod.json" || d.Name() == "mod.xml" {
			m := &mods.Mod{}
			if err = util.LoadFromFile(path, m); err != nil {
				return err
			}
			result = append(result, m)
		}
		return nil
	}); err != nil {
		return
	}
	return appendNewNexusMods(*game, result)
}

func appendNewNexusMods(game config.Game, ms []*mods.Mod) (result []*mods.Mod, err error) {
	var (
		lastID = getLastNexusModID(ms)
		nm     []*mods.Mod
		mod    *mods.Mod
		file   string
	)
	if nm, err = nexus.GetNewestMods(game, lastID); err != nil {
		return
	}

	newModsLastID := getLastNexusModID(nm)
	result = ms
	for id := lastID; id < newModsLastID; id++ {
		// First time getting mods, get them all
		file = filepath.Join(getDir(game), fmt.Sprintf("%d", id), "mod.json")
		if _, err = os.Stat(file); err != nil {
			if mod, err = nexus.GetModFromNexusByID(game, id); err == nil {
				if err = util.SaveToFile(file, mod); err != nil {
					return
				}
				result = append(result, mod)
			}
		}
	}
	return
}

func getLastNexusModID(ms []*mods.Mod) (lastID int) {
	for _, m := range ms {
		if m.ModKind.Kind == mods.Nexus {
			id, _ := m.ModIdAsNumber()
			if int(id) > lastID {
				lastID = int(id)
			}
		}
	}
	return
}
