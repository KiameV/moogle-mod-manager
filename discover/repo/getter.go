package repo

import (
	"context"
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
	"time"
)

type Getter interface {
	GetMod(*mods.Mod) (*mods.Mod, error)
	GetMods(game config.GameDef) ([]*mods.Mod, error)
	Pull() error
	pull(rd repoDef) error
}

type repo struct{}

func NewGetter() Getter {
	return &repo{}
}

func (r *repo) clone() (err error) {
	for _, rd := range repoDefs {
		if err = r.cloneRepo(rd); err == nil {
			return
		}
	}
	return
}

func (r *repo) cloneRepo(rd repoDef) (err error) {
	dir := rd.repoDir()
	if _, err = os.Stat(filepath.Join(dir, ".git")); err != nil {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return
		}
		if _, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL: rd.Url,
			//Progress: os.Stdout,
		}); err != nil {
			return
		}
	}
	return
}

func (r *repo) Pull() (err error) {
	for _, rd := range repoDefs {
		if _, err = os.Stat(rd.repoDir()); err != nil {
			if err = r.clone(); err != nil {
				return
			}
		} else if err = r.pull(rd); err != nil {
			break
		}
	}
	return
}

func (r *repo) pull(rd repoDef) error {
	_, w, err := r.getWorkTree(rd)
	if err == nil {
		err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Force:      true,
		})
		if err == git.NoErrAlreadyUpToDate {
			err = nil
		}
	}
	return err
}

func (r *repo) GetMod(toGet *mods.Mod) (mod *mods.Mod, err error) {
	var (
		dir  string
		game config.GameDef
	)
	for _, rd := range repoDefs {
		if toGet.Category == mods.Utility {
			dir = filepath.Join(rd.repoUtilDir(), toGet.DirectoryName())
		} else if len(toGet.Games) == 1 {
			if game, err = config.GameDefFromID(toGet.Games[0].ID); err != nil {
				return
			}
			dir = filepath.Join(rd.repoGameDir(game), toGet.DirectoryName())
		} else if len(toGet.Games) > 1 {
			return nil, errors.New(toGet.Name + " has multiple games and is not a Utility category")
		} else {
			return nil, errors.New(toGet.Name + " has no games")
		}

		if err = util.LoadFromFile(filepath.Join(dir, "mod.json"), &mod); err == nil {
			return
		}
		if err = util.LoadFromFile(filepath.Join(dir, "mod.xml"), &mod); err == nil {
			return
		}
	}
	return nil, errors.New("unable to find repo file for " + toGet.Name)
}

func (r *repo) GetMods(game config.GameDef) (result []*mods.Mod, err error) {
	var (
		m  []string
		ok error
	)
	for _, rd := range repoDefs {
		if m, err = r.getMods(rd, game); err != nil {
			return nil, err
		}
		for _, f := range m {
			mod := &mods.Mod{}
			if err = util.LoadFromFile(f, mod); err != nil {
				return
			}
			if ok = mod.Supports(game); ok == nil {
				result = append(result, mod)
			}
		}
	}
	return
}

func (r *repo) getMods(rd repoDef, game config.GameDef) (mods []string, err error) {
	if err = r.Pull(); err != nil {
		return
	}
	err = filepath.WalkDir(rd.repoGameDir(game), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "mod.json" || d.Name() == "mod.xml" {
			mods = append(mods, path)
		}
		return nil
	})
	err = filepath.WalkDir(rd.repoUtilDir(), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "mod.json" || d.Name() == "mod.xml" {
			mods = append(mods, path)
		}
		return nil
	})
	return
}

func (r *repo) getWorkTree(rd repoDef) (repo *git.Repository, w *git.Worktree, err error) {
	if repo, err = git.PlainOpen(rd.repoDir()); err != nil {
		return
	}
	w, err = r.getWorkTreeFromRepo(repo)
	return
}

func (*repo) getWorkTreeFromRepo(r *git.Repository) (w *git.Worktree, err error) {
	if w, err = r.Worktree(); err != nil {
		return
	}
	ctx, cnl := context.WithTimeout(context.Background(), time.Second*5)
	defer cnl()
	if err = w.PullContext(ctx, &git.PullOptions{
		RemoteName: "origin",
		Force:      true,
	}); err != nil && err != git.NoErrAlreadyUpToDate {
		_, _ = r.ResolveRevision("origin/main")
	}
	err = nil
	return
}
