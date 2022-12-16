package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"strings"
)

type gamesDef struct {
	*entryManager
	list *cw.DynamicList
}

func newGamesDef() *gamesDef {
	d := &gamesDef{
		entryManager: newEntryManager(),
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
	d.createFormSelect("Games", config.GameIDs(), string(g.ID))
	versions := g.Versions
	var v string
	if versions != nil {
		s := make([]string, len(versions))
		for i, ver := range versions {
			s[i] = string(ver.Version)
		}
		v = strings.Join(s, ", ")
	}
	d.createFormItem("Versions", v)

	fd := dialog.NewForm("Edit Games", "Save", "Cancel", []*widget.FormItem{
		d.getFormItem("Games"),
		d.getFormItem("Versions"),
	}, func(ok bool) {
		if ok {
			g.ID = config.GameID(d.getString("Games"))
			selected := d.getStrings("Versions", ",")
			g.Versions = make([]config.Version, len(selected))
			for i, s := range selected {
				g.Versions[i] = config.Version{Version: config.VersionID(s)}
			}
			if len(done) > 0 {
				done[0](g)
			}
			d.list.Refresh()
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
