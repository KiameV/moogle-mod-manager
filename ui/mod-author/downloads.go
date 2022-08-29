package mod_author

import (
	"fmt"
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
	*entryManager
	list *cw.DynamicList
	kind *mods.Kind
}

func newDownloadsDef(kind *mods.Kind) *downloadsDef {
	d := &downloadsDef{
		entryManager: newEntryManager(),
		kind:         kind,
	}
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
	dl := item.(*mods.Download)
	if dl.Version == "" {
		return dl.Name
	}
	return fmt.Sprintf("%s - %s", dl.Name, dl.Version)
}

func (d *downloadsDef) getItemFields(item interface{}) []string {
	return []string{
		item.(*mods.Download).Name,
		//strings.Join(item.(*mods.Download).Sources, ", "),
		//string(item.(*mods.Download).InstallType),
	}
}

func (d *downloadsDef) onEditItem(item interface{}) {
	d.createItem(item)
}

func (d *downloadsDef) createItem(item interface{}, done ...func(interface{})) {
	m := item.(*mods.Download)
	d.createFormItem("Name", m.Name)
	d.createFormItem("Version", m.Version)
	d.createFormItem("File Name", "")
	d.createFormItem("File ID", "")
	//d.createFormSelect("Install Type", mods.InstallTypes, string(m.InstallType))
	if *d.kind == mods.Nexus {
		if m.Nexus != nil {
			var fileName, fileID string
			if m.Nexus != nil {
				fileName = m.Nexus.FileName
				fileID = fmt.Sprintf("%d", m.Nexus.FileID)
			}
			d.createFormItem("File Name", fileName)
			d.createFormItem("File ID", fileID)
		}
	}
	if *d.kind == mods.Hosted {
		var sources []string
		if m.Hosted != nil {
			sources = m.Hosted.Sources
		}
		d.createFormMultiLine("Sources", strings.Join(sources, "\n"))
	}

	items := []*widget.FormItem{
		d.getFormItem("Name"),
		d.getFormItem("Version"),
	}
	if *d.kind == mods.Nexus {
		items = append(items, d.getFormItem("File Name"))
		items = append(items, d.getFormItem("File ID"))
	}
	if *d.kind == mods.Hosted {
		items = append(items, d.getFormItem("Sources"))
	}

	fd := dialog.NewForm("Edit Downloadable", "Save", "Cancel", items, func(ok bool) {
		if ok {
			m.Name = d.getString("Name")
			m.Version = d.getString("Version")
			if *d.kind == mods.Nexus {
				if m.Nexus == nil {
					m.Nexus = &mods.NexusDownloadable{}
				}
				m.Nexus.FileName = d.getString("File Name")
				m.Nexus.FileID = d.getInt("File ID")
			} else if *d.kind == mods.Hosted {
				if m.Hosted == nil {
					m.Hosted = &mods.HostedDownloadable{}
				}
				m.Hosted.Sources = d.getStrings("Sources", "\n")
			}
			//m.InstallType = mods.InstallType(d.getString("Install Type"))
			if len(done) > 0 {
				done[0](m)
			}
			d.list.Refresh()
		}
	}, state.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *downloadsDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewLabelWithStyle("Downloadables", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Add", func() {
			d.createItem(&mods.Download{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}

func (d *downloadsDef) set(downloadables []*mods.Download) {
	d.list.Clear()
	for _, i := range downloadables {
		d.list.AddItem(i)
	}
}
