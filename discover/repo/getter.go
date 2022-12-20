package repo

import (
	"context"
	"fmt"
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
	GetUtilities() ([]*mods.Mod, error)
	Pull() error
	pull(rd repoDef) error
}

type repo struct {
	kind UseKind
}

func NewGetter(kind UseKind) Getter {
	return &repo{
		kind: kind,
	}
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
	dir := rd.repoDir(r.kind)
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
		if _, err = os.Stat(rd.repoDir(r.kind)); err != nil {
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
		if toGet.Category == mods.Utility && len(toGet.Games) > 1 {
			dir = filepath.Join(rd.repoUtilDir(r.kind), toGet.ID().AsDir())
		} else if len(toGet.Games) == 1 {
			if game, err = config.GameDefFromID(toGet.Games[0].ID); err != nil {
				return
			}
			dir = filepath.Join(rd.repoGameModDir(r.kind, game, toGet))
		} else if len(toGet.Games) > 1 {
			return nil, fmt.Errorf("%s has multiple games and is not a Utility category", toGet.Name)
		} else {
			return nil, fmt.Errorf("%s has no games", toGet.Name)
		}

		if err = util.LoadFromFile(filepath.Join(dir, "mod.json"), &mod); err == nil {
			// Success
			return
		}
		if err = util.LoadFromFile(filepath.Join(dir, "mod.xml"), &mod); err == nil {
			// Success
			return
		}
	}
	return nil, fmt.Errorf("unable to find repo file for %s", toGet.Name)
}

func (r *repo) GetMods(game config.GameDef) (result []*mods.Mod, err error) {
	var m []string
	if err = r.Pull(); err != nil {
		return
	}
	for _, rd := range repoDefs {
		if m, err = r.getMods(rd, game); err != nil {
			return nil, err
		}
	}
	return r.filesToMod(game, m)
}

func (r *repo) GetUtilities() ([]*mods.Mod, error) {
	var (
		m   []string
		err error
	)
	if err = r.Pull(); err != nil {
		return nil, err
	}
	for _, rd := range repoDefs {
		if m, err = r.getUtilities(rd); err != nil {
			return nil, err
		}
	}
	return r.filesToMod(nil, m)
}

func (r *repo) filesToMod(game config.GameDef, m []string) (result []*mods.Mod, err error) {
	for _, f := range m {
		mod := &mods.Mod{}
		if err = util.LoadFromFile(f, mod); err != nil {
			return
		}
		if game != nil {
			if ok := mod.Supports(game); ok == nil {
				result = append(result, mod)
			}
		} else {
			result = append(result, mod)
		}
	}
	return
}

func (r *repo) getMods(rd repoDef, game config.GameDef) (mods []string, err error) {
	if err = filepath.WalkDir(rd.repoGameDir(r.kind, game), func(path string, d os.DirEntry, err error) error {
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
	}); err != nil {
		return
	}
	var u []string
	if u, err = r.getUtilities(rd); err != nil {
		return
	}
	mods = append(mods, u...)
	return
}

func (r *repo) getUtilities(rd repoDef) (mods []string, err error) {
	err = filepath.WalkDir(rd.repoUtilDir(r.kind), func(path string, d os.DirEntry, err error) error {
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
	if repo, err = git.PlainOpen(rd.repoDir(r.kind)); err != nil {
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
