package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/discover/remote/nexus"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
	"github.com/kiamev/moogle-mod-manager/ui/util"
)

type downloadsRemoteDef struct {
	games   *gamesDef
	kind    mods.Kind
	kinds   *mods.Kinds
	parent  *fyne.Container
	dlList  *fyne.Container
	idEntry entry.Entry[string]
}

func newDownloadsRemoteDef(games *gamesDef, kind mods.Kind, kinds *mods.Kinds) *downloadsRemoteDef {
	d := &downloadsRemoteDef{
		games:   games,
		kinds:   kinds,
		dlList:  container.NewVBox(),
		idEntry: entry.NewStringFormEntry(string(kind)+" Mod ID", ""),
	}
	return d
}

func (d *downloadsRemoteDef) compileDownloads() (dls []*mods.Download, err error) {
	var g []config.GameDef
	if g, err = d.games.gameDefs(); err != nil {
		return
	}
	if len(g) == 1 {
		if d.kind == mods.Nexus {
			dls, err = nexus.GetDownloads(g[0], d.idEntry.Value())
		} else if d.kind == mods.CurseForge {
			dls, err = nexus.GetDownloads(g[0], d.idEntry.Value())
		}
	}
	return
}

func (d *downloadsRemoteDef) draw() *container.TabItem {
	d.dlList = container.NewVBox()
	d.parent = container.NewVBox(
		widget.NewForm(d.idEntry.FormItem()),
		widget.NewButton("Load Downloadables", func() {
			dls, err := d.compileDownloads()
			if err != nil {
				util.ShowErrorLong(err)
				return
			}
			d.set(dls)
		}),
	)
	return container.NewTabItem(string(d.kind), d.parent)
}

func (d *downloadsRemoteDef) set(dls []*mods.Download) {
	d.dlList.Objects = nil
	for _, dl := range dls {
		d.dlList.Add(widget.NewLabel("- " + dl.Name))
	}
	d.parent.Refresh()
}

func (d *downloadsRemoteDef) clear() {
	d.idEntry.Set("")
	d.dlList.Objects = nil
	d.parent.Refresh()
}
