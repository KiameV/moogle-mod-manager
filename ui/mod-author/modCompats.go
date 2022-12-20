package mod_author

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	xw "fyne.io/x/fyne/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/discover"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"strings"
)

type modCompatsDef struct {
	*entryManager
	list *cw.DynamicList
	name string
	gd   *gamesDef
}

func newModCompatsDef(name string, gd *gamesDef) *modCompatsDef {
	d := &modCompatsDef{
		entryManager: newEntryManager(),
		name:         name,
		gd:           gd,
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	}, true)
	return d
}

func (d *modCompatsDef) compile() []*mods.ModCompat {
	downloads := make([]*mods.ModCompat, len(d.list.Items))
	for i, item := range d.list.Items {
		downloads[i] = item.(*mods.ModCompat)
	}
	return downloads
}

func (d *modCompatsDef) getItemKey(item interface{}) string {
	name, err := discover.GetDisplayName(state.CurrentGame, item.(*mods.ModCompat).ModID())
	if err != nil {
		name = err.Error()
	}
	return name
}

func (d *modCompatsDef) getItemFields(item interface{}) []string {
	return nil
}

func (d *modCompatsDef) onEditItem(item interface{}) {
	d.createItem(item)
}

func (d *modCompatsDef) createItem(item interface{}, done ...func(interface{})) {
	var m = item.(*mods.ModCompat)

	var (
		game config.GameDef
		err  error
	)
	if d.gd != nil && len(d.gd.list.Items) == 1 {
		game, err = config.GameDefFromID(d.gd.compile()[0].ID)
	}
	if err != nil {
		util.ShowErrorLong(errors.New("please specify a supported Games first (from the Games tab)"))
		return
	}

	modLookup, err := discover.GetModsAsLookup(game)
	if err != nil {
		util.ShowErrorLong(err)
		return
	}

	search := xw.NewCompletionEntry(nil)
	search.SetText(string(m.ModID()))
	search.OnChanged = func(s string) {
		if len(s) < 3 {
			search.HideCompletion()
		}
		s = strings.ToLower(s)
		var results []string
		for _, mod := range modLookup.All() {
			if strings.Contains(strings.ToLower(string(mod.ID())), s) || strings.Contains(strings.ToLower(string(mod.Name)), s) {
				results = append(results, string(mod.Name))
			}
		}
		search.SetOptions(results)
		search.ShowCompletion()
	}

	fd := dialog.NewForm("Edit Mod Compatibility", "Save", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Mod", search),
	}, func(ok bool) {
		if ok {
			var selected *mods.Mod
			m.Hosted = nil
			m.Nexus = nil
			if search.Text != "" {
				for _, mod := range modLookup.All() {
					if mod.Name.Contains(search.Text) {
						selected = mod
						break
					}
				}
				if selected == nil {
					// TODO
					return
				}
				switch selected.ModKind.Kind {
				case mods.Hosted:
					m.Kind = mods.Hosted
					m.Hosted = &mods.ModCompatHosted{
						ModID: selected.ModID,
					}
				case mods.Nexus:
					m.Kind = mods.Nexus
					m.Nexus = &mods.ModCompatNexus{
						ModID: selected.ModID,
					}
				case mods.CurseForge:
					m.Kind = mods.CurseForge
					m.CurseForge = &mods.ModCompatCF{
						ModID: selected.ModID,
					}
				default:
					panic(fmt.Sprint("unknown mod kind: ", selected.ModKind.Kind))
				}
			}
			if len(done) > 0 {
				done[0](m)
			}
		}
	}, ui.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *modCompatsDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewLabelWithStyle(d.name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Add", func() {
			d.createItem(&mods.ModCompat{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}

func (d *modCompatsDef) clear() {
	d.list.Clear()
}
