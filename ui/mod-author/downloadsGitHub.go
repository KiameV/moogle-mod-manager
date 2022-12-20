package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/discover/remote/github"
	"github.com/kiamev/moogle-mod-manager/mods"
	"strings"
)

type githubDownloadsDef struct {
	*entryManager
}

func newGithubDownloadsDef() *githubDownloadsDef {
	return &githubDownloadsDef{
		entryManager: newEntryManager(),
	}
}

func (d *githubDownloadsDef) version() (string, error) {
	return github.LatestRelease(d.getString("owner"), d.getString("repo"))
}

func (d *githubDownloadsDef) compile() (version string, gh *mods.GitHub, result []*mods.Download, err error) {
	var dls []github.Download
	if version, err = d.version(); err != nil {
		return
	}
	if dls, err = github.ListDownloads(d.getString("owner"), d.getString("repo"), version); err != nil {
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
		Owner: d.getString("owner"),
		Repo:  d.getString("repo"),
	}
	return
}

func (d *githubDownloadsDef) draw() fyne.CanvasObject {
	return container.NewVBox(widget.NewForm(
		d.getFormItem("owner"),
		d.getFormItem("repo")))
}

func (d *githubDownloadsDef) set(gh *mods.GitHub) {
	if gh == nil {
		d.clear()
	} else {
		d.createFormItem("owner", gh.Owner)
		d.createFormItem("repo", gh.Repo)
	}
}

func (d *githubDownloadsDef) getFormItems() []*widget.FormItem {
	return []*widget.FormItem{
		d.getFormItem("owner"),
		d.getFormItem("repo"),
	}
}

func (d *githubDownloadsDef) clear() {
	d.createFormItem("owner", "")
	d.createFormItem("repo", "")

}
