package repo

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/kiamev/moogle-mod-manager/config"
	"os"
	"path/filepath"
	"time"
)

type repo struct{}

func (r repo) Clone() (err error) {
	dir := repoDir()
	if _, err = os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return
		}
		if _, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL: "https://github.com/KiameV/moogle-mod-manager-mods.git",
			//Progress: os.Stdout,
		}); err != nil {
			return
		}
	}
	return
}

func (r repo) Pull() (err error) {
	_, _, err = r.getWorkTree()
	return

}

func (r repo) GetMods(game config.Game) (mods []string, overrides []string, err error) {
	if _, err = os.Stat(repoDir()); err != nil {
		if err = r.Clone(); err != nil {
			return
		}
	} else if err = r.Pull(); err != nil {
		return
	}
	err = filepath.WalkDir(repoGameDir(game), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "mod.json" || d.Name() == "mod.xml" {
			mods = append(mods, path)
		}
		if d.Name() == "override.json" || d.Name() == "override.xml" {
			overrides = append(overrides, path)
		}
		return nil
	})
	return
}

func (r repo) getWorkTree() (repo *git.Repository, w *git.Worktree, err error) {
	if repo, err = git.PlainOpen(repoDir()); err != nil {
		return
	}
	w, err = r.getWorkTreeFromRepo(repo)
	return
}

func (repo) getWorkTreeFromRepo(r *git.Repository) (w *git.Worktree, err error) {
	if w, err = r.Worktree(); err != nil {
		return
	}
	ctx, cnl := context.WithTimeout(context.Background(), time.Second*5)
	defer cnl()
	if err = w.PullContext(ctx, &git.PullOptions{RemoteName: "origin"}); err != nil && err != git.NoErrAlreadyUpToDate {
		return
	}
	err = nil
	return
}
