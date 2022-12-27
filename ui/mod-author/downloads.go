package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type (
	Compileable interface {
		CompileDownloads() ([]mods.Download, error)
	}
	downloads struct {
		kinds *mods.Kinds
		dld   *downloadsDef
		ghd   *githubDownloadsDef
		nrd   *downloadsRemoteDef
		cfd   *downloadsRemoteDef
	}
)

func newDownloads(games *gamesDef, kinds *mods.Kinds) *downloads {
	return &downloads{
		kinds: kinds,
		dld:   newDownloadsDef(kinds),
		ghd:   newGithubDownloadsDef(kinds),
		nrd:   newDownloadsRemoteDef(games, mods.Nexus, kinds),
		cfd:   newDownloadsRemoteDef(games, mods.CurseForge, kinds),
	}
}

func (d *downloads) compileDownloads() (result []*mods.Download, err error) {
	var (
		l = make(map[string]*mods.Download)
	)
	if d.kinds.Is(mods.HostedGitHub) {
		if result, err = d.ghd.compileDownloads(); err != nil {
			return
		}
		d.addDlToMap(&l, result)
	}
	if d.kinds.Is(mods.HostedAt) {
		d.addDlToMap(&l, d.dld.compileDownloads())
	}
	if d.kinds.Is(mods.Nexus) {
		if result, err = d.nrd.compileDownloads(); err != nil {
			return
		}
		d.addDlToMap(&l, result)
	}
	if d.kinds.Is(mods.CurseForge) {
		if result, err = d.cfd.compileDownloads(); err != nil {
			return
		}
		d.addDlToMap(&l, result)
	}

	result = make([]*mods.Download, 0, len(l))
	for _, dl := range l {
		result = append(result, dl)
	}
	return
}

func (d *downloads) addDlToMap(l *map[string]*mods.Download, dls []*mods.Download) {
	for _, dl := range dls {
		(*l)[dl.Name] = dl
	}
}

func (d *downloads) compile(mod *mods.Mod) (err error) {
	if d.kinds.Is(mods.HostedGitHub) {
		if mod.Version, mod.ModKind.GitHub, err = d.ghd.compile(); err != nil {
			return
		}
	} else {
		mod.ModKind.GitHub = nil
	}
	mod.Downloadables, err = d.compileDownloads()
	return
}

func (d *downloads) set(mod *mods.Mod) {
	d.ghd.set(mod.ModKind.GitHub)
	d.dld.set(mod.Downloadables)
}

func (d *downloads) clear() {
	d.dld.clear()
	d.ghd.clear()
	d.nrd.clear()
	d.cfd.clear()
}

func (d *downloads) draw() fyne.CanvasObject {
	return container.NewAppTabs(
		d.nrd.draw(),
		d.cfd.draw(),
		d.ghd.draw(),
		d.dld.draw())
}
