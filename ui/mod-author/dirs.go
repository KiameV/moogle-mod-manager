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
)

type dirsDef struct {
	*entryManager
	list *cw.DynamicList
}

func newDirsDef() *dirsDef {
	d := &dirsDef{
		entryManager: newEntryManager(),
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	})
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
	d.createFileDialog("From", f.From, state.GetBaseDirBinding(), true, true)
	d.createFormItem("To", f.To)
	d.createFormBool("Recursive", f.Recursive)

	fd := dialog.NewForm("Edit Directory Copy", "Save", "Cancel", []*widget.FormItem{
		d.getFileDialog("From"),
		d.getFormItem("To"),
		d.getFormItem("Recursive"),
	}, func(ok bool) {
		if ok {
			f.From = d.getString("From")
			f.To = d.getString("To")
			f.Recursive = d.getBool("Recursive")
			if len(done) > 0 {
				done[0](f)
			}
			d.list.Refresh()
		}
	}, state.Window)
	fd.Resize(fyne.NewSize(400, 400))
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
