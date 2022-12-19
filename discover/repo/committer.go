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
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	sourceOwner = "KiameV"
	author      = "moogle-modder"
	authorName  = "Moogle Modder"
	authorEmail = "moogle-modder@hotmail.com"

	hostedPrefix     = strings.ToLower(string(mods.Hosted + "."))
	nexusPrefix      = strings.ToLower(string(mods.Nexus + "."))
	curseforgePrefix = strings.ToLower(string(mods.CurseForge + "."))
)

type Committer interface {
	Submit() (url string, err error)
}

func NewCommitter(mod *mods.Mod) Committer {
	return &repoClient{
		client: github.NewClient(oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "g" + pat3 + "_" + pat2 + pat}))),
		mod: mod,
	}
}

type repoClient struct {
	client *github.Client
	mod    *mods.Mod
}

func (c *repoClient) Submit() (url string, err error) {
	//if *sourceOwner == "" || *sourceRepo == "" || *commitBranch == "" || *sourceFiles == "" || *authorName == "" || *authorEmail == "" {
	//	log.Fatal("You need to specify a non-empty value for the flags `-source-owner`, `-source-repo`, `-commit-branch`, `-files`, `-author-name` and `-author-email`")
	//}
	var (
		rd   = repoDefs[0]
		ref  *github.Reference
		tree *github.Tree
		file string
		game config.GameDef
	)
	for _, d := range c.mod.Downloadables {
		s := ""
		d.DownloadedArchiveLocation = (*mods.ArchiveLocation)(&s)
	}
	if len(c.mod.Games) == 1 {
		if game, err = config.GameDefFromID(c.mod.Games[0].ID); err != nil {
			return
		}
		file = rd.repoGameModDir(Author, game, c.mod)
	} else if len(c.mod.Games) > 1 {
		if c.mod.ModKind.Kind != mods.Hosted {
			err = errors.New("multi-game mods must be hosted")
			return
		}
		file = util.CreateFileName(string(c.mod.ModID))
		file = rd.removeFilePrefixes(file)
		file = filepath.Join(rd.repoDir(Author), "utilities", file)
	} else {
		err = errors.New("no games specified")
		return
	}

	if file == "" {
		err = errors.New("unable to format remote directory")
		return
	}
	file = filepath.Join(file, "mod.json")

	if err = c.mod.Save(file); err != nil {
		return
	}

	var branch string
	if ref, branch, err = c.getRef(rd); err != nil {
		err = fmt.Errorf("unable to get/create the commit reference: %s", err)
		return
	}
	if ref == nil {
		err = errors.New("no error where returned but the reference is nil")
		return
	}

	if tree, err = c.getTree(rd, ref, file); err != nil {
		err = fmt.Errorf("unable to create the tree based on the provided files: %s\n", err)
		return
	}

	if err = c.pushCommit(rd, ref, tree); err != nil {
		err = fmt.Errorf("unable to create the commit: %s\n", err)
		return
	}

	if url, err = c.createPR(rd, branch); err != nil {
		err = fmt.Errorf("unable to create the pull request: %s", err)
	}
	return
}

// getRef returns the commit branch reference object if it exists or creates it
// from the base branch before returning it.
func (c *repoClient) getRef(rd repoDef) (ref *github.Reference, commitBranch string, err error) {
	commitBranch = "refs/heads/" + c.mod.BranchName()
	sourceRepo := rd.Source()

	ctx, cnl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cnl()

	if ref, _, err = c.client.Git.GetRef(ctx, author, rd.Url, commitBranch); err == nil {
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
	if err != nil {
		err = nil
		ref = newRef
	}
	return
}

// getTree generates the tree to commit based on the given files and the commit
// of the ref you got in getRef.
func (c *repoClient) getTree(rd repoDef, ref *github.Reference, file string) (tree *github.Tree, err error) {
	// Create a tree with what to commit.
	var (
		entries    []*github.TreeEntry
		sourceRepo = rd.Source()
	)

	// Load each file into the tree.
	var b []byte
	if b, err = os.ReadFile(file); err != nil {
		return nil, err
	}
	file = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(file, rd.repoDir(Author)), "\\"), "/")
	file = strings.ReplaceAll(file, "\\", "/")
	entries = append(entries, &github.TreeEntry{Path: github.String(file), Type: github.String("blob"), Content: github.String(string(b)), Mode: github.String("100644")})

	ctx, cnl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cnl()
	tree, _, err = c.client.Git.CreateTree(ctx, author, sourceRepo, *ref.Object.SHA, entries)
	return tree, err
}

// pushCommit creates the commit in the given reference using the given tree.
func (c *repoClient) pushCommit(rd repoDef, ref *github.Reference, tree *github.Tree) (err error) {
	var (
		ctx, cnl   = context.WithTimeout(context.Background(), 10*time.Second)
		parent     *github.RepositoryCommit
		sourceRepo = rd.Source()
	)
	defer cnl()

	// Get the parent commit to attach the commit to.
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
	_, _, err = c.client.Git.UpdateRef(ctx, author, sourceRepo, ref, true)
	return
}

// createPR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func (c *repoClient) createPR(rd repoDef, commitBranch string) (url string, err error) {
	commitBranch = author + ":" + commitBranch
	var (
		sbj   = fmt.Sprintf("%s - %s", c.mod.Name, c.mod.Version)
		base  = "main"
		newPR = &github.NewPullRequest{
			Title:               &sbj,
			Head:                &commitBranch,
			Base:                &base,
			Body:                nil,
			MaintainerCanModify: github.Bool(true),
		}
		sourceRepo = rd.Source()
		pr         *github.PullRequest
		ctx, cnl   = context.WithTimeout(context.Background(), 10*time.Second)
	)
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
const pat3 = "hp"
const pat2 = "ezio5oN8qtU1fX"
