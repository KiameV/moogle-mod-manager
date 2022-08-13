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

func Clone() (err error) {
	if _, err = os.Stat(dir()); err != nil {
		if _, err = git.PlainClone(dir(), false, &git.CloneOptions{
			URL: "https://github.com/KiameV/moogle-mod-manager-mods.git",
			//Progress: os.Stdout,
		}); err != nil {
			return
		}
	}
	return
}

func CommitMod(tm *model.TrackedMod) (url string, err error) {
	var w *git.Worktree
	if _, w, err = getWorkTree(); err != nil {
		return
	}

	if err = w.Checkout(&git.CheckoutOptions{
		Hash:   plumbing.NewHash(tm.GetBranchName()),
		Create: true,
	}); err != nil {
		return
	}

	if _, err = w.Add(tm.GetBranchName()); err != nil {
		return
	}

	msg := fmt.Sprintf("release %s for %s", tm.Mod.Version, tm.Mod.Name)
	if _, err = w.Commit(msg, &git.CommitOptions{}); err != nil {
		return
	}

	if url, err = createPR(msg); err != nil {
		return
	}
	return
}

func GetMods(game config.Game) (mods []string, err error) {

}

func dir() string {
	return filepath.Join(config.PWD, "remote")
}

func getWorkTree() (r *git.Repository, w *git.Worktree, err error) {
	if r, err = git.PlainOpen(dir()); err != nil {
		return
	}
	w, err = getWorkTreeFromRepo(r)
	return
}

func getWorkTreeFromRepo(r *git.Repository) (w *git.Worktree, err error) {
	if w, err = r.Worktree(); err != nil {
		return
	}
	if err = w.Pull(&git.PullOptions{RemoteName: "origin/main"}); err != nil {
		return
	}
	return
}

func createPR(subject string, commit string) (url string, err error) {
	const (
		author = "moogle-modder"
		repo   = "KiameV/moogle-mod-manager-mods"
	)
	var (
		client   = github.NewClient(nil)
		base     = "origin/main"
		pr       *github.PullRequest
		ctx, cnl = context.WithTimeout(context.Background(), 5*time.Second)
	)
	defer cnl()

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
