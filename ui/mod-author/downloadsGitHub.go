package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/discover/remote/github"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"strings"
)

type githubDownloadsDef struct {
	entry.Manager
	dlList *fyne.Container
	kinds  *mods.Kinds
}

func newGithubDownloadsDef(kinds *mods.Kinds) dlHoster {
	return &githubDownloadsDef{
		Manager: entry.NewManager(),
		kinds:   kinds,
		dlList:  container.NewVBox(),
	}
}

func (d *githubDownloadsDef) version() (string, error) {
	return github.LatestRelease(entry.Value[string](d, "owner"), entry.Value[string](d, "repo"))
}

func (d *githubDownloadsDef) compile(mod *mods.Mod) error {
	gh, err := d.compileGH()
	if err != nil {
		return err
	}
	if gh.Repo != "" && gh.Owner != "" && gh.Version != "" {
		mod.ModKind.Kinds.Add(mods.HostedGitHub)
		mod.ModKind.GitHub = gh
	}
	return nil
}

func (d *githubDownloadsDef) compileGH() (*mods.GitHub, error) {
	var (
		v, err = d.version()
		gh     = &mods.GitHub{
			Owner:   entry.Value[string](d, "owner"),
			Repo:    entry.Value[string](d, "repo"),
			Version: v,
		}
	)
	if err != nil {
		return nil, err
	}
	return gh, nil
}

func (d *githubDownloadsDef) compileDownloads() (result []*mods.Download, err error) {
	var (
		dls     []github.Download
		version string
	)
	if dls, err = github.ListDownloads(entry.Value[string](d, "owner"), entry.Value[string](d, "repo"), version); err != nil {
		return
	}
	if version, err = d.version(); err != nil {
		return
	}
	result = make([]*mods.Download, len(dls))
	for i, dl := range dls {
		name := dl.Name
		if j := strings.LastIndex(name, "."); j != -1 {
			name = name[:j]
		}
		result[i] = &mods.Download{
			Name:    name,
			Version: version,
			Hosted: &mods.HostedDownloadable{
				Sources: []string{dl.URL},
			},
		}
	}
	return
}

func (d *githubDownloadsDef) draw() *container.TabItem {
	return container.NewTabItem(
		"GitHub",
		container.NewVBox(
			widget.NewForm(
				entry.FormItem[string](d, "owner"),
				entry.FormItem[string](d, "repo")),
			container.NewHBox(widget.NewButton("Load Downloadables", func() {
				if gh, err := d.compileGH(); err != nil {
					util.ShowErrorLong(err)
					return
				} else {
					d.setGH(gh)
				}
			})),
		))
}

func (d *githubDownloadsDef) set(mod *mods.Mod) {
	d.setGH(mod.ModKind.GitHub)
}

func (d *githubDownloadsDef) setGH(gh *mods.GitHub) {
	if gh == nil {
		d.clear()
	} else {
		entry.NewEntry[string](d, entry.KindString, "owner", gh.Owner)
		entry.NewEntry[string](d, entry.KindString, "repo", gh.Repo)
		dls, err := d.compileDownloads()
		if err != nil {
			return
		}
		d.dlList.Objects = nil
		if len(dls) > 0 {
			d.kinds.Add(mods.HostedGitHub)
			for _, dl := range dls {
				d.dlList.Add(widget.NewLabel("- " + dl.Name))
			}
		}
	}
}

func (d *githubDownloadsDef) getFormItems() []*widget.FormItem {
	return []*widget.FormItem{
		entry.FormItem[string](d, "owner"),
		entry.FormItem[string](d, "repo"),
	}
}

func (d *githubDownloadsDef) clear() {
	entry.NewEntry[string](d, entry.KindString, "owner", "")
	entry.NewEntry[string](d, entry.KindString, "repo", "")
	d.kinds.Remove(mods.HostedGitHub)
}
