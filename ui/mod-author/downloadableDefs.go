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

type downloadablesDef struct {
	list *cw.DynamicList
	//Name        string
	//Sources     []string
	//InstallType InstallType
}

func newDownloadablesDef() *downloadablesDef {
	d := &downloadablesDef{}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	})
	return d
}

func (d *downloadablesDef) compile() []*mods.Download {
	downloads := make([]*mods.Download, len(d.list.Items))
	for i, item := range d.list.Items {
		downloads[i] = item.(*mods.Download)
	}
	return downloads
}

func (d *downloadablesDef) getItemKey(item interface{}) string {
	return item.(*mods.Download).Name
}

func (d *downloadablesDef) getItemFields(item interface{}) []string {
	return []string{
		item.(*mods.Download).Name,
		strings.Join(item.(*mods.Download).Sources, ", "),
		string(item.(*mods.Download).InstallType),
	}
}

func (d *downloadablesDef) onEditItem(item interface{}, done func(result interface{})) {
	setFormItem("dlableName", item.(*mods.Download).Name)
	setFormMultiLine("dlableSources", strings.Join(item.(*mods.Download).Sources, "\n"))
	setFormSelect("dlableInstallType", mods.InstallTypes, string(item.(*mods.Download).InstallType))

	fd := dialog.NewForm("Edit Downloadable", "Save", "Cancel", []*widget.FormItem{
		getFormItem("Name", "dlableName"),
		getFormItem("Sources", "dlableSources"),
		getFormItem("Install Type", "dlableInstallType"),
	}, func(ok bool) {
		if ok {
			done(&mods.Download{
				Name:        getFormString("dlableName"),
				Sources:     getFormStrings("gameDefVersion", "\n"),
				InstallType: mods.InstallType(getFormString("dlableInstallType")),
			})
		}
	}, state.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *downloadablesDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewLabelWithStyle("Downloadables", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Add", func() {
			d.onEditItem(&mods.Download{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}
