package mod_author

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"strings"
)

type gamesDef struct {
	entry.Manager
	list      *cw.DynamicList
	gameAdded func(config.GameID)
}

func newGamesDef(gameAdded func(config.GameID)) *gamesDef {
	d := &gamesDef{
		Manager:   entry.NewManager(),
		gameAdded: gameAdded,
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.editItem,
	}, false)
	return d
}

func (d *gamesDef) compile() (games []*mods.Game) {
	games = make([]*mods.Game, len(d.list.Items))
	for i, item := range d.list.Items {
		games[i] = item.(*mods.Game)
	}
	return games
}

func (d *gamesDef) gameDefs() (games []config.GameDef, err error) {
	games = make([]config.GameDef, len(d.list.Items))
	for i, item := range d.list.Items {
		if games[i], err = config.GameDefFromID(item.(*mods.Game).ID); err != nil {
			return
		}
	}
	return
}

func (d *gamesDef) getItemKey(item interface{}) string {
	return string(item.(*mods.Game).ID)
}

func (d *gamesDef) getItemFields(item interface{}) []string {
	versions := item.(*mods.Game).Versions
	if len(versions) == 0 {
		return nil
	}
	result := make([]string, len(versions))
	for i, v := range versions {
		result[i] = string(v.Version)
	}
	return []string{strings.Join(result, ", ")}
}

func (d *gamesDef) editItem(item interface{}) {
	d.createItem(item)
}

func (d *gamesDef) createItem(item interface{}, done ...func(interface{})) {
	g := item.(*mods.Game)
	entry.NewSelectEntry(d, "Games", string(g.ID), config.GameIDs())
	versions := g.Versions
	var v string
	if versions != nil {
		s := make([]string, len(versions))
		for i, ver := range versions {
			s[i] = string(ver.Version)
		}
		v = strings.Join(s, ", ")
	}
	entry.NewEntry[string](d, entry.KindString, "Versions", v)

	fd := dialog.NewForm("Edit Games", "Save", "Cancel", []*widget.FormItem{
		entry.FormItem[string](d, "Games"),
		entry.FormItem[string](d, "Versions"),
	}, func(ok bool) {
		if ok {
			g.ID = config.GameID(entry.Value[string](d, "Games"))
			selected := strings.Split(entry.Value[string](d, "Versions"), ",")
			g.Versions = make([]config.Version, len(selected))
			for i, s := range selected {
				g.Versions[i] = config.Version{Version: config.VersionID(s)}
			}
			if len(done) > 0 {
				done[0](g)
			}
			d.list.Refresh()
			if d.gameAdded != nil {
				d.gameAdded(g.ID)
			}
		}
	}, ui.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *gamesDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewLabelWithStyle("Games", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Add", func() {
			d.createItem(&mods.Game{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}

func (d *gamesDef) set(games []*mods.Game) {
	d.list.Clear()
	for _, g := range games {
		d.list.AddItem(g)
	}
}

func (d *gamesDef) AuthorHintDir() string {
	for _, g := range d.compile() {
		if gd, err := config.GameDefFromID(g.ID); err == nil {
			return fmt.Sprintf("To %s/", gd.AuthorHintDir())
		}
	}
	return "To"
}
