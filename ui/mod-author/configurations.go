package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
)

type configurationsDef struct {
	entry.Manager
	list       *cw.DynamicList
	choicesDef *choicesDef
	previewDef *previewDef
}

func newConfigurationsDef(dlDef *downloads, installType *config.InstallType) *configurationsDef {
	d := &configurationsDef{
		Manager:    entry.NewManager(),
		previewDef: newPreviewDef(),
	}
	d.choicesDef = newChoicesDef(dlDef, d, installType)
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	}, true)
	return d
}

func (d *configurationsDef) compile() []*mods.Configuration {
	cfgs := make([]*mods.Configuration, len(d.list.Items))
	for i, item := range d.list.Items {
		cfgs[i] = item.(*mods.Configuration)
	}
	return cfgs
}

func (d *configurationsDef) getItemKey(item interface{}) string {
	c := item.(*mods.Configuration)
	if c.Root {
		return c.Name + " (root)"
	}
	return c.Name
}

func (d *configurationsDef) getItemFields(item interface{}) []string {
	c := item.(*mods.Configuration)
	return []string{
		c.Name,
		c.Description,
	}
}

func (d *configurationsDef) onEditItem(item interface{}) {
	d.createItem(item)
}

func (d *configurationsDef) createItem(item interface{}, done ...func(interface{})) {
	c := item.(*mods.Configuration)
	entry.NewEntry[string](d, entry.KindString, "Name", c.Name)
	entry.NewEntry[string](d, entry.KindMultiLine, "Description", c.Description)
	entry.NewEntry[bool](d, entry.KindBool, "Root", c.Root)
	d.previewDef.set(c.Preview)
	d.choicesDef.populate(c.Choices)

	items := []*widget.FormItem{
		entry.FormItem[string](d, "Name"),
		entry.FormItem[string](d, "Description"),
		entry.FormItem[bool](d, "Root"),
	}
	items = append(items, d.previewDef.getFormItems()...)
	items = append(items, widget.NewFormItem("Choices", d.choicesDef.draw(false)))

	fd := dialog.NewForm("Edit Configuration", "Save", "Cancel", items, func(ok bool) {
		if ok {
			c.Name = entry.Value[string](d, "Name")
			c.Description = entry.Value[string](d, "Description")
			c.Root = entry.Value[bool](d, "Root")
			c.Preview = d.previewDef.compile()
			c.Choices = d.choicesDef.compile()
			if len(done) > 0 {
				done[0](c)
			}
			d.list.Refresh()
		}
	}, ui.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *configurationsDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewLabelWithStyle("Configurations", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Add", func() {
			d.createItem(&mods.Configuration{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}

func (d *configurationsDef) set(configurations []*mods.Configuration) {
	d.list.Clear()
	for _, c := range configurations {
		d.list.AddItem(c)
	}
}
