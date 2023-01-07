package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
)

type previewsDef struct {
	entry.Manager
	list *cw.DynamicList
}

func newpPreviewsDef() *previewsDef {
	d := &previewsDef{
		Manager: entry.NewManager(),
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	}, true)
	return d
}

func (d *previewsDef) compile() []*mods.Preview {
	p := make([]*mods.Preview, len(d.list.Items))
	for i, item := range d.list.Items {
		p[i] = item.(*mods.Preview)
	}
	return p
}

func (d *previewsDef) getItemKey(item interface{}) string {
	u := item.(*mods.Preview).Url
	if u == nil {
		return ""
	}
	return *u
}

func (d *previewsDef) getItemFields(item interface{}) []string {
	var (
		m = item.(*mods.Preview)
		u string
	)
	if m.Url != nil {
		u = *m.Url
	}
	return []string{u}
}

func (d *previewsDef) onEditItem(item interface{}) {
	d.createItem(item)
}

func (d *previewsDef) createItem(item interface{}, done ...func(interface{})) {
	var (
		m = item.(*mods.Preview)
		u string
	)
	if m.Url != nil {
		u = *m.Url
	}
	entry.NewEntry[string](d, entry.KindString, "Image URL", u)

	fd := dialog.NewForm("Preview Images", "Save", "Cancel", []*widget.FormItem{
		entry.FormItem[string](d, "Image URL"),
	}, func(ok bool) {
		if ok {
			if m.Url == nil {
				m.Url = new(string)
			}
			*m.Url = entry.Value[string](d, "Image URL")
			if len(done) > 0 {
				done[0](m)
			}
			d.list.Refresh()
		}
	}, ui.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *previewsDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewButton("Add", func() {
			d.createItem(&mods.Preview{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}

func (d *previewsDef) set(links []*mods.Preview) {
	d.list.Clear()
	for _, i := range links {
		d.list.AddItem(i)
	}
}
