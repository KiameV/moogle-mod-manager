package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/discover/remote/github"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
	"strings"
)

type githubDownloadsDef struct {
	entry.Manager
}

func newGithubDownloadsDef() *githubDownloadsDef {
	return &githubDownloadsDef{
		Manager: entry.NewManager(),
	}
}

func (d *githubDownloadsDef) version() (string, error) {
	return github.LatestRelease(entry.Value[string](d, "owner"), entry.Value[string](d, "repo"))
}

func (d *githubDownloadsDef) compile() (version string, gh *mods.GitHub, result []*mods.Download, err error) {
	var dls []github.Download
	if version, err = d.version(); err != nil {
		return
	}
	if dls, err = github.ListDownloads(entry.Value[string](d, "owner"), entry.Value[string](d, "repo"), version); err != nil {
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
	gh = &mods.GitHub{
		Owner: entry.Value[string](d, "owner"),
		Repo:  entry.Value[string](d, "repo"),
	}
	return
}

func (d *githubDownloadsDef) draw() fyne.CanvasObject {
	return container.NewVBox(widget.NewForm(
		entry.FormItem[string](d, "owner"),
		entry.FormItem[string](d, "repo")))
}

func (d *githubDownloadsDef) set(gh *mods.GitHub) {
	if gh == nil {
		d.clear()
	} else {
		entry.NewEntry[string](d, entry.KindString, "owner", gh.Owner)
		entry.NewEntry[string](d, entry.KindString, "repo", gh.Repo)
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

}
