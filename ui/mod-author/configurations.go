package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
)

type configurationsDef struct {
	*entryManager
	list       *cw.DynamicList
	choicesDef *choicesDef
	previewDef *previewDef
}

func newConfigurationsDef(dlDef *downloads) *configurationsDef {
	d := &configurationsDef{
		entryManager: newEntryManager(),
		previewDef:   newPreviewDef(),
	}
	d.choicesDef = newChoicesDef(dlDef, d)
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
	d.createFormItem("Name", c.Name)
	d.createFormMultiLine("Description", c.Description)
	d.createFormBool("Root", c.Root)
	d.previewDef.set(c.Preview)
	d.choicesDef.populate(c.Choices)

	items := []*widget.FormItem{
		d.getFormItem("Name"),
		d.getFormItem("Description"),
		d.getFormItem("Root"),
	}
	items = append(items, d.previewDef.getFormItems()...)
	items = append(items, widget.NewFormItem("Choices", d.choicesDef.draw(false)))

	fd := dialog.NewForm("Edit Configuration", "Save", "Cancel", items, func(ok bool) {
		if ok {
			c.Name = d.getString("Name")
			c.Description = d.getString("Description")
			c.Root = d.getBool("Root")
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
