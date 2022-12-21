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
	"strings"
)

type dirsDef struct {
	entry.Manager
	list *cw.DynamicList
}

func newDirsDef() *dirsDef {
	d := &dirsDef{
		Manager: entry.NewManager(),
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	}, false)
	return d
}

func (d *dirsDef) compile() []*mods.ModDir {
	downloads := make([]*mods.ModDir, len(d.list.Items))
	for i, item := range d.list.Items {
		downloads[i] = item.(*mods.ModDir)
	}
	return downloads
}

func (d *dirsDef) getItemKey(item interface{}) string {
	f := item.(*mods.ModDir)
	return fmt.Sprintf("%s -> %s", f.From, f.To)
}

func (d *dirsDef) getItemFields(item interface{}) []string {
	f := item.(*mods.ModDir)
	return []string{
		f.From,
		f.To,
	}
}

func (d *dirsDef) onEditItem(item interface{}) {
	d.createItem(item)
}

func (d *dirsDef) createItem(item interface{}, done ...func(interface{})) {
	f := item.(*mods.ModDir)
	entry.CreateFileDialog(d, "From", f.From, state.GetBaseDirBinding(), true, true)
	entry.NewEntry[string](d, entry.KindString, "To FF PR/", f.To)
	entry.NewEntry[bool](d, entry.KindBool, "Recursive", f.Recursive)

	fd := dialog.NewForm("Edit Directory Copy", "Save", "Cancel", []*widget.FormItem{
		entry.GetFileDialog(d, "From"),
		entry.FormItem[string](d, "To FF PR/"),
		entry.FormItem[bool](d, "Recursive"),
	}, func(ok bool) {
		if ok {
			f.From = cleanPath(entry.DialogValue(d, "From"))
			f.To = cleanPath(entry.Value[string](d, "To FF PR/"))
			f.Recursive = entry.Value[bool](d, "Recursive")
			if len(done) > 0 {
				done[0](f)
			}
			d.list.Refresh()
		}
	}, ui.Window)
	fd.Resize(fyne.NewSize(600, 400))
	fd.Show()
}

func (d *dirsDef) draw(label bool) fyne.CanvasObject {
	c := container.NewHBox()
	if label {
		c.Add(widget.NewLabelWithStyle("Dirs", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	c.Add(widget.NewButton("Add", func() {
		d.createItem(&mods.ModDir{}, func(result interface{}) {
			d.list.AddItem(result)
		})
	}))
	return container.NewVBox(
		c,
		d.list.Draw())
}

func (d *dirsDef) clear() {
	d.list.Clear()
}

func (d *dirsDef) populate(dirs []*mods.ModDir) {
	d.list.Clear()
	for _, dir := range dirs {
		d.list.AddItem(dir)
	}
}

func cleanPath(s string) string {
	s = strings.ReplaceAll(s, "\\", "/")
	return strings.ReplaceAll(s, "//", "/")
}
