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
	authorDir       = "author"
	repoDir         = "repo"
	defaultRepoName = "mmmm"
	defaultRepoUrl  = "https://github.com/KiameV/moogle-mod-manager-mods"
)

var repoDefs []repoDef

type (
	UseKind byte
	repoDef struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
)

const (
	_ UseKind = iota
	Author
	Read
)

func (d repoDef) Source() string {
	sp := strings.Split(d.Url, "/")
	return sp[len(sp)-1]
}

func (d repoDef) repoDir(k UseKind) string {
	dir := repoDir
	if k == Author {
		dir = authorDir
	}
	return filepath.Join(config.PWD, dir, d.Name)
}

func (d repoDef) repoUtilDir(k UseKind) string {
	return filepath.Join(d.repoDir(k), "utilities")
}

func (d repoDef) repoGameDir(k UseKind, game config.GameDef) string {
	if game == nil {
		return ""
	}
	return filepath.Join(d.repoDir(k), string(game.ID()))
}

func (d repoDef) repoGameModDir(k UseKind, game config.GameDef, mod *mods.Mod) string {
	return filepath.Join(d.repoGameDir(k, game), strings.ToLower(string(mod.Kind())), d.removeFilePrefixes(strings.ToLower(mod.ID().AsDir())))
}

func (d repoDef) removeFilePrefixes(s string) string {
	s = strings.TrimPrefix(s, hostedPrefix)
	s = strings.TrimPrefix(s, nexusPrefix)
	s = strings.TrimPrefix(s, curseforgePrefix)
	return s
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

func Dirs(k UseKind) (dirs []string) {
	dirs = make([]string, len(repoDefs))
	for i, rd := range repoDefs {
		dirs[i] = rd.repoDir(k)
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
