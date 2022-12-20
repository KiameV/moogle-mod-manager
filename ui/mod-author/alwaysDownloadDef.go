package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
)

type alwaysDownloadDef struct {
	*entryManager
	list             *cw.DynamicList
	downloadFilesDef *downloadFilesDef
}

func newAlwaysDownloadDef(downloads *downloads) *alwaysDownloadDef {
	d := &alwaysDownloadDef{
		entryManager:     newEntryManager(),
		downloadFilesDef: newDownloadFilesDef(downloads),
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	}, true)
	return d
}

func (d *alwaysDownloadDef) compile() []*mods.DownloadFiles {
	downloads := make([]*mods.DownloadFiles, len(d.list.Items))
	for i, item := range d.list.Items {
		downloads[i] = item.(*mods.DownloadFiles)
	}
	return downloads
}

func (d *alwaysDownloadDef) getItemKey(item interface{}) string {
	dlf := item.(*mods.DownloadFiles)
	return dlf.DownloadName
}

func (d *alwaysDownloadDef) getItemFields(item interface{}) []string {
	return []string{}
}

func (d *alwaysDownloadDef) onEditItem(item interface{}) {
	d.createItem(item)
}

func (d *alwaysDownloadDef) createItem(item interface{}, done ...func(interface{})) {
	dlf := item.(*mods.DownloadFiles)
	d.downloadFilesDef.populate(dlf)

	fi, err := d.downloadFilesDef.getFormItems()
	if err != nil {
		dialog.ShowError(err, ui.Window)
		return
	}

	fd := dialog.NewForm("Edit Download Files", "Save", "Cancel", fi,
		func(ok bool) {
			if ok {
				result := d.downloadFilesDef.compile()
				*dlf = *result
				if len(done) > 0 {
					done[0](dlf)
				}
				d.list.Refresh()
			}
		}, ui.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *alwaysDownloadDef) draw() fyne.CanvasObject {
	return container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Always Download", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewButton("Add", func() {
				d.createItem(&mods.DownloadFiles{}, func(result interface{}) {
					d.list.AddItem(result)
				})
			})),
		d.list.Draw())
}

func (d *alwaysDownloadDef) set(alwaysDownload []*mods.DownloadFiles) {
	d.list.Clear()
	for _, f := range alwaysDownload {
		d.list.AddItem(f)
	}
}
