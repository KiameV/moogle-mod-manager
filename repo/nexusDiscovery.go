package repo

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/nexus"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
	"strconv"
)

func getNexusMods(game config.Game, ms []*mods.Mod) ([]*mods.Mod, error) {
	newMods, err := addNewNexusMods(game, ms)
	return append(ms, newMods...), err
}

func addNewNexusMods(game config.Game, ms []*mods.Mod) (newMods []*mods.Mod, err error) {
	var (
		lastID = getLastNexusModID(ms)
		mod    *mods.Mod
		file   string
	)
	if newMods, err = nexus.GetNewestMods(game, lastID); err != nil {
		return
	}

	newModsLastID := getLastNexusModID(newMods)
	newMods = make([]*mods.Mod, 0, newModsLastID-lastID+1)
	for id := lastID; id < newModsLastID; id++ {
		// First time getting mods, get them all
		file = filepath.Join(repoDefs[0].repoNexusIDDir(game, fmt.Sprintf("%d", id)), "mod.json")
		if _, err = os.Stat(file); err != nil {
			if mod, err = nexus.GetModFromNexusByID(game, id); err == nil {
				newMods = append(newMods, mod)
				if err = util.SaveToFile(file, mod); err != nil {
					return
				}
			}
		}
	}
	return
}

func getLastNexusModID(ms []*mods.Mod) (lastID int) {
	for _, m := range ms {
		if m.ModKind.Kind == mods.Nexus {
			id, _ := strconv.ParseInt(m.ModKind.Nexus.ID, 10, 64)
			if int(id) > lastID {
				lastID = int(id)
			}
		}
	}
	return
}

/*
const lastIDsFile = "lastNexusIDs.json"
type nexusLastIDs struct {
	iID   []int
	iiID  []int
	iiiID []int
	ivID  []int
	vID   []int
	viID  []int
}

func (d *nexusDiscovery) initialize() {
	if d.lastIDs == nil {
		d.lastIDs = &nexusLastIDs{}
		_ = util.LoadFromFile(filepath.Join(config.PWD, lastIDsFile), d.lastIDs)
	}
}

func (d *nexusDiscovery) update() {
	_ = util.SaveToFile(filepath.Join(config.PWD, lastIDsFile), d.lastIDs)
}
*/
