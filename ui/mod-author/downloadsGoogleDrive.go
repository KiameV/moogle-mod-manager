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

type googleDriveDownloadsDef struct {
	entry.Manager
	list  *cw.DynamicList
	kinds *mods.Kinds
}

func newGoogleDriveDownloadsDef(kinds *mods.Kinds) dlHoster {
	d := &googleDriveDownloadsDef{
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

func (d *googleDriveDownloadsDef) compile(mod *mods.Mod) error {
	if len(d.list.Items) > 0 {
		mod.ModKind.Kinds.Add(mods.GoogleDrive)
	}
	return nil
}

func (d *googleDriveDownloadsDef) compileDownloads() ([]*mods.Download, error) {
	dls := make([]*mods.Download, len(d.list.Items))
	for i, item := range d.list.Items {
		di := item.(*mods.Download)
		if di.GoogleDrive != nil {
			di.Name = di.GoogleDrive.Name
			if j := strings.LastIndex(di.Name, "."); j != -1 {
				di.Name = di.Name[:j]
			}
		}
		dls[i] = di
	}
	return dls, nil
}

func (d *googleDriveDownloadsDef) getItemKey(item interface{}) string {
	dl := item.(*mods.Download)
	return fmt.Sprintf(dl.Name)
}

func (d *googleDriveDownloadsDef) getItemFields(item interface{}) []string {
	return []string{
		item.(*mods.Download).Name,
	}
}

func (d *googleDriveDownloadsDef) onEditItem(item interface{}) {
	d.createItem(item)
}

func (d *googleDriveDownloadsDef) createItem(item interface{}, done ...func(interface{})) {
	var (
		items    []*widget.FormItem
		m        = item.(*mods.Download)
		fileName string
		url      string
	)
	if gd := m.GoogleDrive; gd != nil {
		fileName = gd.Name
		url = gd.Url
	}
	entry.NewEntry[string](d, entry.KindString, "File Name", fileName)
	entry.NewEntry[string](d, entry.KindString, "Version", m.Version)
	entry.NewEntry[string](d, entry.KindString, "URL", url)

	items = []*widget.FormItem{
		entry.FormItem[string](d, "File Name"),
		entry.FormItem[string](d, "Version"),
		entry.FormItem[string](d, "URL"),
	}

	fd := dialog.NewForm("Edit Downloadable", "Save", "Cancel", items, func(ok bool) {
		if ok {
			m.Version = entry.Value[string](d, "Version")
			if m.GoogleDrive == nil {
				m.GoogleDrive = &mods.GoogleDriveDownloadable{}
			}
			fileName = strings.TrimSpace(entry.Value[string](d, "File Name"))
			m.GoogleDrive.Name = fileName
			m.GoogleDrive.Url = entry.Value[string](d, "URL")
			m.Name = strings.TrimSuffix(fileName, filepath.Ext(fileName))
			for _, dn := range done {
				dn(m)
			}
			d.list.Refresh()
		}
	}, ui.Window)
	fd.Resize(fyne.NewSize(600, 400))
	fd.Show()
}

func (d *googleDriveDownloadsDef) draw() *container.TabItem {
	return container.NewTabItem("Google Drive",
		container.NewVScroll(container.NewVBox(container.NewHBox(
			widget.NewLabelWithStyle("Downloadables", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewButton("Add", func() {
				d.createItem(&mods.Download{}, func(result interface{}) {
					d.list.AddItem(result)
					d.kinds.Add(mods.GoogleDrive)
				})
			})),
			d.list.Draw())))
}

func (d *googleDriveDownloadsDef) set(mod *mods.Mod) {
	d.clear()
	if mod.ModKind.Kinds.Is(mods.GoogleDrive) {
		for _, i := range mod.Downloadables {
			if i.GoogleDrive != nil {
				d.list.AddItem(i)
				d.kinds.Add(mods.GoogleDrive)
			}
		}
	}
}

func (d *googleDriveDownloadsDef) clear() {
	d.list.Clear()
	d.kinds.Remove(mods.GoogleDrive)
}
