package repo

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultRepoName = "mmmm"
	defaultRepoUrl  = "https://github.com/KiameV/moogle-mod-manager-mods"
)

var repoDefs []repoDef

type repoDef struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (d repoDef) Source() string {
	sp := strings.Split(d.Url, "/")
	return sp[len(sp)-1]
}

func (d repoDef) repoDir() string {
	return filepath.Join(config.PWD, "remote", d.Name)
}

func (d repoDef) repoGameDir(game config.Game) string {
	return filepath.Join(d.repoDir(), config.String(game))
}

func (d repoDef) repoNexusIDDir(game config.Game, id string) string {
	return filepath.Join(d.repoGameDir(game), "nexus", id)
}

func (d repoDef) repoNexusDir(game config.Game, mod *mods.Mod) string {
	return d.repoNexusIDDir(game, mod.ModKind.Nexus.ID)
}

func Initialize() (err error) {
	f := filepath.Join(config.PWD, "repo.json")
	if len(repoDefs) == 0 {
		if _, err = os.Stat(f); err != nil {
			repoDefs = []repoDef{{
				Name: defaultRepoName,
				Url:  defaultRepoUrl,
			}}
			return saveDefaultRepo(f)
		}
		if err = util.LoadFromFile(f, &repoDefs); err != nil {
			return
		}
	}
	if len(repoDefs) == 0 {
		err = fmt.Errorf("no repositories found in %s, using default repository", f)
		_ = saveDefaultRepo(f)
	}
	return
}

func saveDefaultRepo(f string) error {
	repoDefs = []repoDef{{
		Name: defaultRepoName,
		Url:  defaultRepoUrl,
	}}
	return util.SaveToFile(f, &repoDefs)
}
