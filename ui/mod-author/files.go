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
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
)

type filesDef struct {
	entry.Manager
	list *cw.DynamicList
}

func newFilesDef() *filesDef {
	d := &filesDef{
		Manager: entry.NewManager(),
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	}, false)
	return d
}

func (d *filesDef) compile() []*mods.ModFile {
	downloads := make([]*mods.ModFile, len(d.list.Items))
	for i, item := range d.list.Items {
		downloads[i] = item.(*mods.ModFile)
	}
	return downloads
}

func (d *filesDef) getItemKey(item interface{}) string {
	f := item.(*mods.ModFile)
	return fmt.Sprintf("%s -> %s", f.From, f.To)
}

func (d *filesDef) getItemFields(item interface{}) []string {
	f := item.(*mods.ModFile)
	return []string{
		f.From,
		f.To,
	}
}

func (d *filesDef) onEditItem(item interface{}) {
	d.createItem(item)
}

func (d *filesDef) createItem(item interface{}, done ...func(interface{})) {
	f := item.(*mods.ModFile)
	entry.CreateFileDialog(d, "From", f.From, state.GetBaseDirBinding(), false, true)
	entry.NewEntry[string](d, entry.KindString, "To FF PR/", f.To)

	fd := dialog.NewForm("Edit File Copy", "Save", "Cancel", []*widget.FormItem{
		entry.GetFileDialog(d, "From"),
		entry.FormItem[string](d, "To FF PR/"),
	}, func(ok bool) {
		if ok {
			f.From = cleanPath(entry.Value[string](d, "From"))
			f.To = cleanPath(entry.Value[string](d, "To FF PR/"))
			if len(done) > 0 {
				done[0](f)
			}
			d.list.Refresh()
		}
	}, ui.Window)
	fd.Resize(fyne.NewSize(600, 400))
	fd.Show()
}

func (d *filesDef) draw(label bool) fyne.CanvasObject {
	c := container.NewHBox()
	if label {
		c.Add(widget.NewLabelWithStyle("Files", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	c.Add(widget.NewButton("Add", func() {
		d.createItem(&mods.ModFile{}, func(result interface{}) {
			d.list.AddItem(result)
		})
	}))
	return container.NewVBox(
		c,
		d.list.Draw())
}

func (d *filesDef) clear() {
	d.list.Clear()
}

func (d *filesDef) populate(files []*mods.ModFile) {
	d.list.Clear()
	for _, f := range files {
		d.list.AddItem(f)
	}
}
