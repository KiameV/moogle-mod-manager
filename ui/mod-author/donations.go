package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

type donationsDef struct {
	*entryManager
	list *cw.DynamicList
}

func newDonationsDef() *donationsDef {
	d := &donationsDef{
		entryManager: newEntryManager(),
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	})
	return d
}

func (d *donationsDef) compile() []*mods.DonationLink {
	downloads := make([]*mods.DonationLink, len(d.list.Items))
	for i, item := range d.list.Items {
		downloads[i] = item.(*mods.DonationLink)
	}
	return downloads
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

func (d *donationsDef) onEditItem(item interface{}, done func(result interface{})) {
	m := item.(*mods.DonationLink)
	d.createFormItem("Name", m.Name)
	d.createFormItem("Link", m.Link)

	fd := dialog.NewForm("Edit Donation", "Save", "Cancel", []*widget.FormItem{
		d.getFormItem("Name"),
		d.getFormItem("Link"),
	}, func(ok bool) {
		if ok {
			done(&mods.DonationLink{
				Name: d.getString("Name"),
				Link: d.getString("Link"),
			})
		}
	}, state.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *donationsDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewLabelWithStyle("Donation Links", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Add", func() {
			d.onEditItem(&mods.DonationLink{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}
