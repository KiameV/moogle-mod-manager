package mod_author

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type (
	Compileable interface {
		CompileDownloads() ([]mods.Download, error)
	}
	downloads struct {
		*container.TabItem
		subKind *widget.Select
		dld     *downloadsDef
		ghd     *githubDownloadsDef
	}
)

func newDownloads(kind *mods.Kind, s *widget.Select) *downloads {
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
		mod.Version, mod.ModKind.GitHub, mod.Downloadables, err = d.ghd.compile()
	} else {
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
	return d.subKind.Selected == string(mods.HostedGitHub)
}

func (d *downloads) UpdateTab() {
	if d.TabItem == nil {
		d.TabItem = container.NewTabItem("", container.NewCenter())
	}
	switch mods.SubKind(d.subKind.Selected) {
	case mods.HostedAt:
		d.TabItem.Text = "Downloads"
		d.TabItem.Content = d.dld.draw()
	case mods.HostedGitHub:
		d.TabItem.Text = "GitHub"
		d.TabItem.Content = d.ghd.draw()
	default:
		d.TabItem.Text = "-"
		d.TabItem.Content = widget.NewLabel("Select a 'Kind'")
	}
}
