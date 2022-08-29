package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"strings"
)

type modCompatsDef struct {
	*entryManager
	list *cw.DynamicList
	name string
}

func newModCompatsDef(name string) *modCompatsDef {
	d := &modCompatsDef{
		entryManager: newEntryManager(),
		name:         name,
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	})
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
	return item.(*mods.ModCompat).DisplayName()
}

func (d *modCompatsDef) getItemFields(item interface{}) []string {
	return nil
}

func (d *modCompatsDef) onEditItem(item interface{}) {
	d.createItem(item)
}

func (d *modCompatsDef) createItem(item interface{}, done ...func(interface{})) {
	m := item.(*mods.ModCompat)
	d.createFormItem("Mod ID", m.ModID())
	d.createFormItem("Versions", strings.Join(m.Versions, ", "))
	//order := mods.None
	//if m.Order != nil {
	//	order = *m.Order
	//}
	//d.createFormSelect("Order", mods.ModCompatOrders, string(order))

	fd := dialog.NewForm("Edit Mod Compatibility", "Save", "Cancel", []*widget.FormItem{
		d.getFormItem("Kind"),
		d.getFormItem("Mod ID"),
		d.getFormItem("Versions"),
	}, func(ok bool) {
		if ok {
			m.Hosted = nil
			m.Nexus = nil
			m.Versions = d.getStrings("Versions", ",")
			if d.getString("kind") == string(mods.Hosted) {
				m.Kind = mods.Hosted
				m.Hosted = &mods.ModCompatHosted{
					ModID: d.getString("Mod ID"),
				}
			} else {
				m.Kind = mods.Nexus
				m.Nexus = &mods.ModCompatNexus{
					ModID: d.getString("Mod ID"),
				}
			}
			//o := d.getString("Order")
			//if o != "" && o != string(mods.None) {
			//	m.Order = (*mods.ModCompatOrder)(&o)
			//} else {
			//	m.Order = nil
			//}
			if len(done) > 0 {
				done[0](m)
			}
		}
	}, state.Window)
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
