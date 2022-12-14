package util

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
	"strings"
)

type Finder interface {
	GetNewestMods(game config.GameDef, lastID int) ([]*mods.Mod, error)
	GetFromID(game config.GameDef, id int) (found bool, mod *mods.Mod, err error)
}

type ModCompiler interface {
	AppendNewMods(folder string, game config.GameDef, ms []*mods.Mod) (result []*mods.Mod, err error)
	SetFinder(finder Finder)
}

type modCompiler struct {
	finder Finder
	kind   mods.Kind
}

func NewModCompiler(kind mods.Kind) ModCompiler {
	return &modCompiler{kind: kind}
}

func (c *modCompiler) SetFinder(finder Finder) {
	c.finder = finder
}

func (c *modCompiler) AppendNewMods(folder string, game config.GameDef, ms []*mods.Mod) (result []*mods.Mod, err error) {
	var (
		lastID = c.getLastModID(ms)
		nm     []*mods.Mod
		mod    *mods.Mod
		file   string
		found  bool
	)
	if nm, err = c.finder.GetNewestMods(game, lastID); err != nil {
		return
	}

	if c.kind == mods.Nexus {
		newModsLastID := c.getLastModID(nm)
		result = ms
		for id := lastID; id < newModsLastID; id++ {
			file = filepath.Join(folder, fmt.Sprintf("%d", id), "mod.json")
			if _, err = os.Stat(file); err != nil {
				if found, mod, err = c.finder.GetFromID(game, id); found && err == nil {
					if err = util.SaveToFile(file, mod); err != nil {
						return
					}
					result = append(result, mod)
				}
			}
		}
	} else if c.kind == mods.CurseForge {
		for _, mod = range nm {
			id := strings.Split(string(mod.ModID), ".")[1]
			file = filepath.Join(folder, id, "mod.json")
			if _, err = os.Stat(file); err != nil {
				if err = util.SaveToFile(file, mod); err != nil {
					return
				}
				result = append(result, mod)
			}
		}
	} else {
		err = fmt.Errorf("invalid kind %v", c.kind)
	}
	return
}

func (c *modCompiler) getLastModID(ms []*mods.Mod) (lastID int) {
	for _, m := range ms {
		if m.ModKind.Kind == c.kind {
			id, _ := m.ModIdAsNumber()
			if int(id) > lastID {
				lastID = int(id)
			}
		}
	}
	return
}
