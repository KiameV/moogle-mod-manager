package mod_author

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"path/filepath"
	"strings"
)

type downloadsDef struct {
	entry.Manager
	list  *cw.DynamicList
	kinds *mods.Kinds
}

func newDownloadsDef(kinds *mods.Kinds) *downloadsDef {
	d := &downloadsDef{
		Manager: entry.NewManager(),
		kinds:   kinds,
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	}, true)
	return d
}

func (d *downloadsDef) compileDownloads() []*mods.Download {
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
		items []*widget.FormItem
		m     = item.(*mods.Download)
	)
	var sources []string
	if m.Hosted != nil {
		sources = m.Hosted.Sources
	}
	entry.NewEntry[string](d, entry.KindMultiLine, "Sources", strings.Join(sources, "\n"))
	entry.NewEntry[string](d, entry.KindString, "Version", m.Version)

	items = []*widget.FormItem{
		entry.FormItem[string](d, "Version"),
		entry.FormItem[string](d, "Sources"),
	}

	fd := dialog.NewForm("Edit Downloadable", "Save", "Cancel", items, func(ok bool) {
		if ok {
			m.Version = entry.Value[string](d, "Version")
			if m.Hosted == nil {
				m.Hosted = &mods.HostedDownloadable{}
			}
			m.Hosted.Sources = strings.Split(entry.Value[string](d, "Sources"), "\n")
			if len(m.Hosted.Sources) > 0 {
				m.Name = filepath.Base(m.Hosted.Sources[0])
			}
			if m.Name != "" {
				m.Name = strings.TrimSuffix(m.Name, filepath.Ext(m.Name))
			}
			//m.InstallType = mods.InstallType(entry.Value[string](d, "Install Type"))
			for _, dn := range done {
				dn(m)
			}
			d.list.Refresh()
		}
	}, ui.Window)
	fd.Resize(fyne.NewSize(600, 400))
	fd.Show()
}

func (d *downloadsDef) draw() *container.TabItem {
	return container.NewTabItem("Direct Download",
		container.NewVScroll(container.NewVBox(container.NewHBox(
			widget.NewLabelWithStyle("Downloadables", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewButton("Add", func() {
				d.createItem(&mods.Download{}, func(result interface{}) {
					d.list.AddItem(result)
				})
			})),
			d.list.Draw())))
}

func (d *downloadsDef) set(downloadables []*mods.Download) {
	d.list.Clear()
	for _, i := range downloadables {
		d.list.AddItem(i)
	}
}

func (d *downloadsDef) clear() {
	d.list.Clear()
}
