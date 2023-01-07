package mod_author

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/discover/remote/curseforge"
	"github.com/kiamev/moogle-mod-manager/discover/remote/nexus"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"strconv"
)

type downloadsRemoteDef struct {
	games     *gamesDef
	kind      mods.Kind
	parent    *fyne.Container
	dlList    *fyne.Container
	listItems []*mods.Download
	idEntry   entry.Entry[string]
}

func newDownloadsRemoteDef(games *gamesDef, kind mods.Kind) dlHoster {
	d := &downloadsRemoteDef{
		games:   games,
		kind:    kind,
		dlList:  container.NewVBox(),
		idEntry: entry.NewStringFormEntry(string(kind)+" Mod ID", ""),
	}
	return d
}

func (d *downloadsRemoteDef) compile(mod *mods.Mod) (err error) {
	var id int64
	if d.idEntry.Value() != "" {
		if id, err = strconv.ParseInt(d.idEntry.Value(), 10, 64); err != nil {
			return
		}
		i := int(id)
		if d.kind.Is(mods.Nexus) {
			mod.ModKind.Kinds.Add(mods.Nexus)
			mod.ModKind.NexusID = (*mods.NexusModID)(&i)
		} else if d.kind.Is(mods.CurseForge) {
			mod.ModKind.Kinds.Add(mods.CurseForge)
			mod.ModKind.CurseForgeID = (*mods.CfModID)(&i)
		}
	}
	return
}

func (d *downloadsRemoteDef) compileDownloads() (dls []*mods.Download, err error) {
	if d.idEntry.Value() != "" {
		dls = d.listItems
	}
	return
}

func (d *downloadsRemoteDef) loadDownloads() (err error) {
	var (
		g   []config.GameDef
		dls []*mods.Download
	)
	if g, err = d.games.gameDefs(); err != nil {
		return
	}
	if len(g) == 1 {
		if d.kind == mods.Nexus {
			dls, err = nexus.GetDownloads(g[0], d.idEntry.Value())
		} else if d.kind == mods.CurseForge {
			dls, err = curseforge.GetDownloads(d.idEntry.Value())
		}
	} else {
		err = errors.New("select a game this mod will work with")
		return
	}
	d.setDownloadables(dls)
	return
}

func (d *downloadsRemoteDef) draw() *container.TabItem {
	d.parent = container.NewVBox(
		widget.NewForm(d.idEntry.FormItem()),
		container.NewHBox(widget.NewButton("Load Downloadables", func() {
			if err := d.loadDownloads(); err != nil {
				util.ShowErrorLong(err)
				return
			}
		})),
		widget.NewLabel("Downloads:"),
		d.dlList,
	)
	return container.NewTabItem(string(d.kind), d.parent)
}

func (d *downloadsRemoteDef) set(mod *mods.Mod) {
	d.clear()
	if d.kind.Is(mods.Nexus) && mod.ModKind.NexusID != nil {
		d.idEntry.Set(fmt.Sprintf("%d", *mod.ModKind.NexusID))
	} else if d.kind.Is(mods.CurseForge) && mod.ModKind.CurseForgeID != nil {
		d.idEntry.Set(fmt.Sprintf("%d", *mod.ModKind.CurseForgeID))
	}
	d.setDownloadables(mod.Downloadables)
}

func (d *downloadsRemoteDef) setDownloadables(dls []*mods.Download) {
	var (
		isNexus = d.kind.Is(mods.Nexus)
		isCf    = d.kind.Is(mods.CurseForge)
	)
	d.listItems = nil
	d.dlList.Objects = nil
	for _, dl := range dls {
		if (isNexus && dl.Nexus != nil) ||
			(isCf && dl.CurseForge != nil) {
			d.listItems = append(d.listItems, dl)
			d.dlList.Add(widget.NewLabel("- " + dl.Name))
		}
	}
	d.dlList.Refresh()
	if d.parent != nil {
		d.parent.Refresh()
	}
}

func (d *downloadsRemoteDef) clear() {
	d.listItems = nil
	d.dlList.Objects = nil
	d.idEntry.Set("")
	d.dlList.Objects = nil
	if d.parent != nil {
		d.parent.Refresh()
	}
}
