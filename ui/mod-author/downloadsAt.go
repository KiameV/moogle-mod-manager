package mod_author

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"path/filepath"
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
	}, true)
	return d
}

func (d *downloadsDef) compile() []*mods.Download {
	dls := make([]*mods.Download, len(d.list.Items))
	for i, item := range d.list.Items {
		di := item.(*mods.Download)
		if di.Hosted != nil && len(di.Hosted.Sources) > 0 {
			di.Name = filepath.Base(di.Hosted.Sources[0])
			if j := strings.LastIndex(di.Name, "."); j != -1 {
				di.Name = di.Name[:j]
			}
		}
		dls[i] = di
	}
	return dls
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
	var (
		items = []*widget.FormItem{
			d.newFormItem("Version"),
		}
		m        = item.(*mods.Download)
		k        = *d.kind
		fileName string
		fileID   string
		url      string
	)
	d.createFormItem("Version", m.Version)
	d.createFormItem("File Name", "")
	d.createFormItem("File ModID", "")
	d.createFormItem("Url", "")
	//d.createFormSelect("Install Type", mods.InstallTypes, string(m.InstallType))
	switch k {
	case mods.Nexus:
		if m.Nexus != nil {
			if m.Nexus != nil {
				fileName = m.Nexus.FileName
				fileID = fmt.Sprintf("%d", m.Nexus.FileID)
			}
		}
	case mods.CurseForge:
		if m.CurseForge != nil {
			fileName = m.CurseForge.FileName
			fileID = fmt.Sprintf("%d", m.CurseForge.FileID)
			url = m.CurseForge.Url
		}
	case mods.Hosted:
		var sources []string
		if m.Hosted != nil {
			sources = m.Hosted.Sources
		}
		d.createFormMultiLine("Sources", strings.Join(sources, "\n"))
	default:
		dialog.ShowError(fmt.Errorf("unknown mod kind: %s", *d.kind), ui.Window)
		for _, dn := range done {
			dn(nil)
		}
	}

	if k == mods.Nexus || k == mods.CurseForge {
		d.createFormItem("File Name", fileName)
		d.createFormItem("File ModID", fileID)
		items = append(items, d.getFormItem("File Name"))
		items = append(items, d.getFormItem("File ModID"))
		if k == mods.CurseForge {
			d.createFormItem("Url", url)
			items = append(items, d.getFormItem("Url"))
		}
	} else if k == mods.Hosted {
		items = append(items, d.getFormItem("Sources"))
	}

	fd := dialog.NewForm("Edit Downloadable", "Save", "Cancel", items, func(ok bool) {
		if ok {
			m.Version = d.getString("Version")
			if k == mods.Nexus {
				if m.Nexus == nil {
					m.Nexus = &mods.RemoteDownloadable{}
				}
				m.Nexus.FileName = d.getString("File Name")
				m.Nexus.FileID = d.getInt("File ModID")
				m.Name = filepath.Base(m.Nexus.FileName)
			} else if k == mods.CurseForge {
				if m.CurseForge == nil {
					m.CurseForge = &mods.CurseForgeDownloadable{}
				}
				m.CurseForge.FileName = d.getString("File Name")
				m.CurseForge.FileID = d.getInt("File ModID")
				m.CurseForge.Url = d.getString("Url")
				m.Name = filepath.Base(m.CurseForge.FileName)
			} else if k == mods.Hosted {
				if m.Hosted == nil {
					m.Hosted = &mods.HostedDownloadable{}
				}
				m.Hosted.Sources = d.getStrings("Sources", "\n")
				if len(m.Hosted.Sources) > 0 {
					m.Name = filepath.Base(m.Hosted.Sources[0])
				}
			}
			if m.Name != "" {
				m.Name = strings.TrimSuffix(m.Name, filepath.Ext(m.Name))
			}
			//m.InstallType = mods.InstallType(d.getString("Install Type"))
			for _, dn := range done {
				dn(m)
			}
			d.list.Refresh()
		}
	}, ui.Window)
	fd.Resize(fyne.NewSize(600, 400))
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
