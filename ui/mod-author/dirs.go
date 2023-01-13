package mod_author

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"strings"
)

type dirsDef struct {
	entry.Manager
	list        *cw.DynamicList
	installType *config.InstallType
	gamesDef    *gamesDef
}

func newDirsDef(installType *config.InstallType, gamesDef *gamesDef) *dirsDef {
	d := &dirsDef{
		Manager:     entry.NewManager(),
		installType: installType,
		gamesDef:    gamesDef,
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	}, false)
	return d
}

func (d *dirsDef) compile() []*mods.ModDir {
	dl := make([]*mods.ModDir, len(d.list.Items))
	for i, item := range d.list.Items {
		dl[i] = item.(*mods.ModDir)
	}
	return dl
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
	entry.NewEntry[string](d, entry.KindString, d.gamesDef.AuthorHintDir(), f.To)
	entry.NewEntry[bool](d, entry.KindBool, "Recursive", f.Recursive)
	s := ""
	if f.ToArchive != nil {
		s = *f.ToArchive
	}
	entry.NewEntry[string](d, entry.KindString, "To Archive", s)

	items := []*widget.FormItem{
		entry.GetFileDialog(d, "From"),
		entry.FormItem[string](d, d.gamesDef.AuthorHintDir()),
		entry.FormItem[bool](d, "Recursive"),
	}
	if d.installType.Is(config.MoveToArchive) {
		items = append(items, entry.FormItem[string](d, "To Archive"))
	}

	fd := dialog.NewForm("Edit Directory Copy", "Save", "Cancel", items,
		func(ok bool) {
			if ok {
				f.From = cleanPath(entry.DialogValue(d, "From"))
				f.To = cleanPath(entry.Value[string](d, d.gamesDef.AuthorHintDir()))
				f.Recursive = entry.Value[bool](d, "Recursive")
				if s = entry.Value[string](d, "To Archive"); s == "" {
					f.ToArchive = nil
				} else {
					f.ToArchive = &s
				}
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
