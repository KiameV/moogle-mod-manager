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
	list *cw.DynamicList
	//Name        string
	//Sources     []string
	//InstallType InstallType
}

func newDonationsDef() *donationsDef {
	d := &donationsDef{}
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
	return []string{
		item.(*mods.DonationLink).Name,
		item.(*mods.DonationLink).Link,
	}
}

func (d *donationsDef) onEditItem(item interface{}, done func(result interface{})) {
	setFormItem("donationName", item.(*mods.DonationLink).Name)
	setFormMultiLine("donationLink", item.(*mods.DonationLink).Link)

	fd := dialog.NewForm("Edit Donation", "Save", "Cancel", []*widget.FormItem{
		getFormItem("Name", "donationName"),
		getFormItem("Link", "donationLink"),
	}, func(ok bool) {
		if ok {
			done(&mods.DonationLink{
				Name: getFormString("donationName"),
				Link: getFormString("donationLink"),
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
