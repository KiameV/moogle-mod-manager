package mod_author

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
)

type (
	Compileable interface {
		CompileDownloads() ([]mods.Download, error)
	}
	downloads struct {
		*container.TabItem
		subKind entry.Entry[string]
		dld     *downloadsDef
		ghd     *githubDownloadsDef
	}
)

func newDownloads(kind *mods.Kind, s entry.Entry[string]) *downloads {
	d := &downloads{
		subKind: s,
		dld:     newDownloadsDef(kind),
		ghd:     newGithubDownloadsDef(),
	}
	d.UpdateTab()
	return d
}

func (d *downloads) compileDownloads() (result []*mods.Download, err error) {
	if d.isGithub() {
		_, _, result, err = d.ghd.compile()
	} else {
		result = d.dld.compile()
	}
	return
}

func (d *downloads) compile(mod *mods.Mod) (err error) {
	if d.isGithub() {
		mod.ModKind.Kind = mods.Hosted
		sk := mods.HostedGitHub
		mod.ModKind.SubKind = &sk
		mod.Version, mod.ModKind.GitHub, mod.Downloadables, err = d.ghd.compile()
	} else {
		if mod.Kind().Is(mods.Hosted) {
			sk := mods.HostedAt
			mod.ModKind.SubKind = &sk
		}
		mod.Downloadables = d.dld.compile()
	}
	return
}

func (d *downloads) set(mod *mods.Mod) {
	d.ghd.set(mod.ModKind.GitHub)
	d.dld.set(mod.Downloadables)
}

func (d *downloads) clear() {
	if d.isGithub() {
		d.ghd.clear()
	}
}

func (d *downloads) isGithub() bool {
	return d.subKind.Value() == string(mods.HostedGitHub)
}

func (d *downloads) UpdateTab() {
	if d.TabItem == nil {
		d.TabItem = container.NewTabItem("", container.NewCenter())
	}
	switch d.subKind.Value() {
	case string(mods.HostedAt):
		d.TabItem.Text = "Downloads"
		d.TabItem.Content = d.dld.draw()
	case string(mods.HostedGitHub):
		d.TabItem.Text = "GitHub"
		d.TabItem.Content = d.ghd.draw()
	default:
		d.TabItem.Text = "-"
		d.TabItem.Content = widget.NewLabel("Select a 'Kind'")
	}
}
