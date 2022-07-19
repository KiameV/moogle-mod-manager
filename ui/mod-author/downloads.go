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

type downloadsDef struct {
	list *cw.DynamicList
	//Name        string
	//Sources     []string
	//InstallType InstallType
}

func newDownloadsDef() *downloadsDef {
	d := &downloadsDef{}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	})
	return d
}

func (d *downloadsDef) compile() []*mods.Download {
	downloads := make([]*mods.Download, len(d.list.Items))
	for i, item := range d.list.Items {
		downloads[i] = item.(*mods.Download)
	}
	return downloads
}

func (d *downloadsDef) getItemKey(item interface{}) string {
	return item.(*mods.Download).Name
}

func (d *downloadsDef) getItemFields(item interface{}) []string {
	return []string{
		item.(*mods.Download).Name,
		strings.Join(item.(*mods.Download).Sources, ", "),
		string(item.(*mods.Download).InstallType),
	}
}

func (d *downloadsDef) onEditItem(item interface{}, done func(result interface{})) {
	setFormItem("dlName", item.(*mods.Download).Name)
	setFormMultiLine("dlSources", strings.Join(item.(*mods.Download).Sources, "\n"))
	setFormSelect("dlInstallType", mods.InstallTypes, string(item.(*mods.Download).InstallType))

	fd := dialog.NewForm("Edit Downloadable", "Save", "Cancel", []*widget.FormItem{
		getFormItem("Name", "dlName"),
		getFormItem("Sources", "dlSources"),
		getFormItem("Install Type", "dlInstallType"),
	}, func(ok bool) {
		if ok {
			done(&mods.Download{
				Name:        getFormString("dlName"),
				Sources:     getFormStrings("gameDefVersion", "\n"),
				InstallType: mods.InstallType(getFormString("dlInstallType")),
			})
		}
	}, state.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *downloadsDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewLabelWithStyle("Downloadables", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Add", func() {
			d.onEditItem(&mods.Download{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}
