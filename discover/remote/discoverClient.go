package remote

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/remote"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
)

type discoverClient interface {
	GetMods(game *config.Game) (result []*mods.Mod, err error)
}

type implClient struct {
	client     remote.Client
	folderName string
}

func newNexusClient() discoverClient {
	return &implClient{
		client:     remote.NewNexusClient(),
		folderName: "nexus",
	}
}

func newCurseForgeClient() discoverClient {
	return &implClient{
		client:     remote.NewCurseForgeClient(),
		folderName: "cf",
	}
}

func (c *implClient) getDir(game config.Game) string {
	return filepath.Join(config.PWD, "remote", config.String(game), c.folderName)
}

func (c *implClient) GetMods(game *config.Game) (result []*mods.Mod, err error) {
	if game == nil {
		return
	}
	dir := c.getDir(*game)
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
	return c.appendNewNexusMods(*game, result)
}

func (c *implClient) appendNewNexusMods(game config.Game, ms []*mods.Mod) (result []*mods.Mod, err error) {
	var (
		lastID = c.getLastModID(ms)
		nm     []*mods.Mod
		mod    *mods.Mod
		file   string
		found  bool
		nc     = remote.NewNexusClient()
	)
	if nm, err = nc.GetNewestMods(game, lastID); err != nil {
		return
	}

	newModsLastID := c.getLastModID(nm)
	result = ms
	for id := lastID; id < newModsLastID; id++ {
		// First time getting mods, get them all
		file = filepath.Join(c.getDir(game), fmt.Sprintf("%d", id), "mod.json")
		if _, err = os.Stat(file); err != nil {
			if found, mod, err = nc.GetFromID(game, id); found && err == nil {
				if err = util.SaveToFile(file, mod); err != nil {
					return
				}
				result = append(result, mod)
			}
		}
	}
	return
}

func (c *implClient) getLastModID(ms []*mods.Mod) (lastID int) {
	for _, m := range ms {
		if m.ModKind.Kind == mods.Nexus || m.ModKind.Kind == mods.CurseForge {
			id, _ := m.ModIdAsNumber()
			if int(id) > lastID {
				lastID = int(id)
			}
		}
	}
	return
}
