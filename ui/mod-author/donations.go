package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
)

type donationsDef struct {
	entry.Manager
	list *cw.DynamicList
}

func newDonationsDef() *donationsDef {
	d := &donationsDef{
		Manager: entry.NewManager(),
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	}, true)
	return d
}

func (d *donationsDef) compile() []*mods.DonationLink {
	dls := make([]*mods.DonationLink, len(d.list.Items))
	for i, item := range d.list.Items {
		dls[i] = item.(*mods.DonationLink)
	}
	return dls
}

func (d *donationsDef) getItemKey(item interface{}) string {
	return item.(*mods.DonationLink).Name
}

func (d *donationsDef) getItemFields(item interface{}) []string {
	m := item.(*mods.DonationLink)
	return []string{
		m.Name,
		m.Link,
	}
}

func (d *donationsDef) onEditItem(item interface{}) {
	d.createItem(item)
}

func (d *donationsDef) createItem(item interface{}, done ...func(interface{})) {
	m := item.(*mods.DonationLink)
	entry.NewEntry[string](d, entry.KindString, "Name", m.Name)
	entry.NewEntry[string](d, entry.KindString, "Link", m.Link)

	fd := dialog.NewForm("Edit Donation", "Save", "Cancel", []*widget.FormItem{
		entry.FormItem[string](d, "Name"),
		entry.FormItem[string](d, "Link"),
	}, func(ok bool) {
		if ok {
			m.Name = entry.Value[string](d, "Name")
			m.Link = entry.Value[string](d, "Link")
			if len(done) > 0 {
				done[0](m)
			}
			d.list.Refresh()
		}
	}, ui.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *donationsDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewLabelWithStyle("Donation Links", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Add", func() {
			d.createItem(&mods.DonationLink{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}

func (d *donationsDef) set(links []*mods.DonationLink) {
	d.list.Clear()
	for _, i := range links {
		d.list.AddItem(i)
	}
}
