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

func (r *repo) Clone() (err error) {
	for _, rd := range repoDefs {
		if err = r.clone(rd); err == nil {
			return
		}
	}
	return
}

func (r *repo) clone(rd repoDef) (err error) {
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
		if err = r.pull(rd); err != nil {
			break
		}
	}
	return
}

func (r *repo) pull(rd repoDef) (err error) {
	_, _, err = r.getWorkTree(rd)
	return

}

func (r *repo) GetMods(game config.Game) (mods []string, err error) {
	var m []string
	for _, rd := range repoDefs {
		if m, err = r.getMods(rd, game); err != nil {
			return nil, err
		}
		mods = append(mods, m...)
	}
	return
}

func (r *repo) getMods(rd repoDef, game config.Game) (mods []string, err error) {
	if _, err = os.Stat(rd.repoDir()); err != nil {
		if err = r.Clone(); err != nil {
			return
		}
	} else if err = r.Pull(); err != nil {
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
