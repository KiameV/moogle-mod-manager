package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v45/github"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/util"
	"golang.org/x/oauth2"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

var (
	sourceOwner = "KiameV"
	sourceRepo  = "moogle-mod-manager-mods"
	author      = "moogle-modder"
	authorName  = "Moogle Modder"
	authorEmail = "moogle-modder@hotmail.com"
)

type Committer interface {
	Submit() (url string, err error)
}

func NewCommitter(mod *mods.Mod) Committer {
	return &committer{
		client: github.NewClient(oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: pat2 + pat}))),
		mod: mod,
	}
}

type committer struct {
	client *github.Client
	mod    *mods.Mod
}

func (c *committer) Submit() (url string, err error) {
	//if *sourceOwner == "" || *sourceRepo == "" || *commitBranch == "" || *sourceFiles == "" || *authorName == "" || *authorEmail == "" {
	//	log.Fatal("You need to specify a non-empty value for the flags `-source-owner`, `-source-repo`, `-commit-branch`, `-files`, `-author-name` and `-author-email`")
	//}
	var (
		ref  *github.Reference
		tree *github.Tree
	)

	var file string
	if c.mod.Game != nil {
		if c.mod.ModKind.Kind == mods.Hosted {
			file = filepath.Join(repoGameDir(config.NameToGame(c.mod.Game.Name)), c.mod.ID)
		} else if c.mod.ModKind.Kind == mods.Nexus && c.mod.ModKind.Nexus != nil {
			file = repoNexusIDDir(config.NameToGame(c.mod.Game.Name), c.mod.ModKind.Nexus.ID)
		}
	}
	if file == "" {
		err = errors.New("unable to format remote directory")
		return
	}
	file = filepath.Join(file, "mod.json")

	if err = util.SaveToFile(file, c.mod); err != nil {
		return
	}

	var branch string
	if ref, branch, err = c.getRef(); err != nil {
		err = fmt.Errorf("unable to get/create the commit reference: %s", err)
		return
	}
	if ref == nil {
		err = errors.New("no error where returned but the reference is nil")
		return
	}

	if tree, err = c.getTree(ref, file); err != nil {
		err = fmt.Errorf("unable to create the tree based on the provided files: %s\n", err)
		return
	}

	if err = c.pushCommit(ref, tree); err != nil {
		err = fmt.Errorf("unable to create the commit: %s\n", err)
		return
	}

	if url, err = c.createPR(branch); err != nil {
		err = fmt.Errorf("unable to create the pull request: %s", err)
	}
	return
}

// getRef returns the commit branch reference object if it exists or creates it
// from the base branch before returning it.
func (c *committer) getRef() (ref *github.Reference, commitBranch string, err error) {
	commitBranch = "refs/heads/" + c.mod.BranchName()

	ctx, cnl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cnl()

	if ref, _, err = c.client.Git.GetRef(ctx, author, sourceRepo, commitBranch); err == nil {
		return
	}

	// We consider that an error means the branch has not been found and needs to
	// be created.
	//if *commitBranch == *baseBranch {
	//	return nil, errors.New("the commit branch does not exist but `-base-branch` is the same as `-commit-branch`")
	//}

	var baseRef *github.Reference
	if baseRef, _, err = c.client.Git.GetRef(ctx, author, sourceRepo, "refs/heads/main"); err != nil {
		return
	}
	newRef := &github.Reference{Ref: github.String(commitBranch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = c.client.Git.CreateRef(ctx, author, sourceRepo, newRef)
	return
}

// getTree generates the tree to commit based on the given files and the commit
// of the ref you got in getRef.
func (c *committer) getTree(ref *github.Reference, file string) (tree *github.Tree, err error) {
	// Create a tree with what to commit.
	var entries []*github.TreeEntry

	// Load each file into the tree.
	var b []byte
	if b, err = ioutil.ReadFile(file); err != nil {
		return nil, err
	}
	entries = append(entries, &github.TreeEntry{Path: github.String(file), Type: github.String("blob"), Content: github.String(string(b)), Mode: github.String("100644")})

	ctx, cnl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cnl()
	tree, _, err = c.client.Git.CreateTree(ctx, author, sourceRepo, *ref.Object.SHA, entries)
	return tree, err
}

// pushCommit creates the commit in the given reference using the given tree.
func (c *committer) pushCommit(ref *github.Reference, tree *github.Tree) (err error) {
	ctx, cnl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cnl()

	// Get the parent commit to attach the commit to.
	var parent *github.RepositoryCommit
	if parent, _, err = c.client.Repositories.GetCommit(ctx, author, sourceRepo, *ref.Object.SHA, nil); err != nil {
		return
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	a := &github.CommitAuthor{Date: &date, Name: &authorName, Email: &authorEmail}
	msg := fmt.Sprintf("%s - %s", c.mod.Name, c.mod.Version)
	commit := &github.Commit{Author: a, Message: &msg, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	var nc *github.Commit
	if nc, _, err = c.client.Git.CreateCommit(ctx, author, sourceRepo, commit); err != nil {
		return
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = nc.SHA
	_, _, err = c.client.Git.UpdateRef(ctx, author, sourceRepo, ref, false)
	return
}

// createPR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func (c *committer) createPR(commitBranch string) (url string, err error) {
	sbj := fmt.Sprintf("%s - %s", c.mod.Name, c.mod.Version)
	base := "main"
	commitBranch = author + ":" + commitBranch
	newPR := &github.NewPullRequest{
		Title:               &sbj,
		Head:                &commitBranch,
		Base:                &base,
		Body:                nil,
		MaintainerCanModify: github.Bool(true),
	}

	var pr *github.PullRequest

	ctx, cnl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cnl()
	if pr, _, err = c.client.PullRequests.Create(ctx, sourceOwner, sourceRepo, newPR); err != nil {
		if !strings.Contains(err.Error(), "pull request already exists") {
			return
		}
	}

	if err != nil {
		url = fmt.Sprintf("https://github.com/%s/%s/pull", sourceOwner, sourceRepo)
		err = nil
	} else {
		url = pr.GetHTMLURL()
	}
	return
}

const pat = "4IYjtV7j9BWmyiSJ1GRz8e"
const pat2 = "ghp_ezio5oN8qtU1fX"
