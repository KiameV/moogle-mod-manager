package mod_author

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type (
	dlHoster interface {
		clear()
		compile(mod *mods.Mod) error
		compileDownloads() ([]*mods.Download, error)
		draw() *container.TabItem
		set(*mods.Mod)
	}
	downloads struct {
		kinds     *mods.Kinds
		dlHosters []dlHoster
	}
)

func newDownloads(games *gamesDef, kinds *mods.Kinds) *downloads {
	return &downloads{
		kinds: kinds,
		dlHosters: []dlHoster{
			newDownloadsDef(kinds),
			newGithubDownloadsDef(kinds),
			newDownloadsRemoteDef(games, mods.Nexus),
			newDownloadsRemoteDef(games, mods.CurseForge),
		},
	}
}

func (d *downloads) compileDownloads() (result []*mods.Download, err error) {
	var (
		l   = make(map[string]*mods.Download)
		dls []*mods.Download
	)
	for _, h := range d.dlHosters {
		if dls, err = h.compileDownloads(); err != nil {
			return
		}
		if len(dls) > 0 && len(l) > 0 {
			if len(dls) != len(l) {
				err = errors.New("number of downloads must be equal")
				return
			}
			for _, dl := range dls {
				if _, ok := l[dl.Name]; !ok {
					err = errors.New("download names must be equal")
					return
				}
			}
			for k, _ := range l {
				for _, dl := range dls {
					found := false
					if dl.Name == k {
						found = true
					}
					if !found {
						err = errors.New("download names must be equal")
						return
					}
				}
			}
		}
		d.addDlToMap(&l, dls)
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
	var dls []*mods.Download
	for _, h := range d.dlHosters {
		if err = h.compile(mod); err != nil {
			return
		}
	}
	if dls, err = d.compileDownloads(); err != nil {
		return
	}
	mod.Downloadables = dls
	return
}

func (d *downloads) set(mod *mods.Mod) {
	for _, h := range d.dlHosters {
		h.set(mod)
	}
}

func (d *downloads) clear() {
	for _, h := range d.dlHosters {
		h.clear()
	}
}

func (d *downloads) draw() fyne.CanvasObject {
	t := container.NewAppTabs()
	for _, h := range d.dlHosters {
		t.Append(h.draw())
	}
	return t
}
