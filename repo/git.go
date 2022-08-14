package repo

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v45/github"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
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

func (r repo) CommitMod(tm *model.TrackedMod) (url string, err error) {
	var (
		w      *git.Worktree
		branch = tm.GetBranchName()
	)
	if _, w, err = r.getWorkTree(); err != nil {
		return
	}

	if err = w.Checkout(&git.CheckoutOptions{
		Hash:   plumbing.NewHash(branch),
		Create: true,
	}); err != nil {
		return
	}

	if _, err = w.Add(branch); err != nil {
		return
	}

	msg := fmt.Sprintf("release %s for %s", tm.Mod.Version, tm.Mod.Name)
	if _, err = w.Commit(msg, &git.CommitOptions{}); err != nil {
		return
	}

	if url, err = r.createPR(msg, branch); err != nil {
		return
	}
	return
}

func (r repo) Pull() (err error) {
	_, _, err = r.getWorkTree()
	return

}

func (r repo) GetMods(game config.Game) (mods []string, overrides []string, err error) {
	dir := repoDir()
	if _, err = os.Stat(dir); err != nil {
		if err = r.Clone(); err != nil {
			return
		}
	} else if err = r.Pull(); err != nil {
		return
	}
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
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

func (repo) createPR(subject string, commit string) (url string, err error) {
	const (
		author = "moogle-modder"
		repo   = "KiameV/moogle-mod-manager-mods"
	)
	var (
		client   = github.NewClient(nil)
		base     = "origin/main"
		prs      []*github.PullRequest
		pr       *github.PullRequest
		ctx, cnl = context.WithTimeout(context.Background(), 5*time.Second)
	)
	defer cnl()

	if prs, _, err = client.PullRequests.List(ctx, author, repo, &github.PullRequestListOptions{}); err != nil {
		return
	}

	for _, pr = range prs {
		if pr.GetTitle() == subject {
			url = pr.GetHTMLURL()
			return
		}
	}

	if pr, _, err = client.PullRequests.Create(ctx, author, repo, &github.NewPullRequest{
		Title:               &subject,
		Head:                &commit,
		Base:                &base,
		Body:                &commit,
		MaintainerCanModify: github.Bool(true),
	}); err != nil {
		return
	}
	url = pr.GetHTMLURL()
	return
}
