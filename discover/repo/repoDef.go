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
	return filepath.Join(config.PWD, "repo", d.Name)
}

func (d repoDef) repoUtilDir() string {
	return filepath.Join(d.repoDir(), "utilities")
}

func (d repoDef) repoGameDir(game config.GameDef) string {
	return filepath.Join(d.repoDir(), string(game.ID()))
}

func (d repoDef) repoGameModDir(game config.GameDef, kind mods.Kind, id mods.ModID) string {
	return filepath.Join(d.repoGameDir(game), string(kind), string(id))
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

func Dirs() (dirs []string) {
	dirs = make([]string, len(repoDefs))
	for i, rd := range repoDefs {
		dirs[i] = rd.repoDir()
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
