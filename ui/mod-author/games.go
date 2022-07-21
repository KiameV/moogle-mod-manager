package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state"
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
	})
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
	return string(item.(*mods.Game).Name)
}

func (d *gamesDef) getItemFields(item interface{}) []string {
	v := item.(*mods.Game).Versions
	if len(v) == 0 {
		return nil
	}
	return []string{strings.Join(v, ", ")}
}

func (d *gamesDef) editItem(item interface{}, done func(result interface{})) {
	d.createFormSelect("Game", []string{
		config.String(config.I),
		config.String(config.II),
		config.String(config.III),
		config.String(config.IV),
		config.String(config.V),
		config.String(config.VI),
	}, config.String(config.NameToGame(item.(*mods.Game).Name)))
	var v string
	versions := item.(*mods.Game).Versions
	if versions != nil {
		v = strings.Join(versions, ", ")
	}
	d.createFormItem("Versions", v)

	fd := dialog.NewForm("Edit Game", "Save", "Cancel", []*widget.FormItem{
		d.getFormItem("Game"),
		d.getFormItem("Versions"),
	}, func(ok bool) {
		if ok {
			done(&mods.Game{
				Name:     config.GameToName(config.FromString(d.getString("Game"))),
				Versions: d.getStrings("Versions", ","),
			})
		}
	}, state.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *gamesDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewLabelWithStyle("Games", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Add", func() {
			d.editItem(&mods.Game{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}

func (d *gamesDef) set(game []*mods.Game) {
	d.list.Clear()
	for _, i := range game {
		d.list.AddItem(i)
	}
}
