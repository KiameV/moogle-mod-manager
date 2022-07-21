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
	return item.(*mods.ModCompat).ModID
}

func (d *modCompatsDef) getItemFields(item interface{}) []string {
	m := item.(*mods.ModCompat)
	s := []string{
		m.ModID,
		m.Source,
	}
	if len(m.Versions) > 0 {
		s = append(s, strings.Join(m.Versions, ", "))
	}
	if m.Order != nil && *m.Order != mods.None {
		s = append(s, string(*m.Order))
	}
	return s
}

func (d *modCompatsDef) onEditItem(item interface{}, done func(result interface{})) {
	m := item.(*mods.ModCompat)
	d.createFormItem("Mod ID", m.ModID)
	d.createFormItem("Source", m.Source)
	d.createFormItem("Versions", strings.Join(m.Versions, ", "))
	order := mods.None
	if m.Order != nil {
		order = *m.Order
	}
	d.createFormSelect("Order", mods.ModCompatOrders, string(order))

	fd := dialog.NewForm("Edit Mod Compatibility", "Save", "Cancel", []*widget.FormItem{
		d.getFormItem("Mod ID"),
		d.getFormItem("Source"),
		d.getFormItem("Versions"),
		d.getFormItem("Order"),
	}, func(ok bool) {
		if ok {
			mc := &mods.ModCompat{
				ModID:    d.getString("Mod ID"),
				Source:   d.getString("Source"),
				Versions: d.getStrings("Versions", ","),
			}
			o := d.getString("Order")
			if o != "" && o != string(mods.None) {
				mc.Order = (*mods.ModCompatOrder)(&o)
			}
			done(mc)
		}
	}, state.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *modCompatsDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewLabelWithStyle(d.name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Add", func() {
			d.onEditItem(&mods.ModCompat{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}

func (d *modCompatsDef) clear() {
	d.list.Clear()
}
