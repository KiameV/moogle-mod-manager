package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

type configurationsDef struct {
	*entryManager
	list       *cw.DynamicList
	choicesDef *choicesDef
}

func newConfigurationsDef(dlDef *downloadsDef) *configurationsDef {
	d := &configurationsDef{
		entryManager: newEntryManager(),
	}
	d.choicesDef = newChoicesDef(dlDef, d)
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	})
	return d
}

func (d *configurationsDef) compile() []*mods.Configuration {
	downloads := make([]*mods.Configuration, len(d.list.Items))
	for i, item := range d.list.Items {
		downloads[i] = item.(*mods.Configuration)
	}
	return downloads
}

func (d *configurationsDef) getItemKey(item interface{}) string {
	return item.(*mods.Configuration).Name
}

func (d *configurationsDef) getItemFields(item interface{}) []string {
	c := item.(*mods.Configuration)
	return []string{
		c.Name,
		c.Description,
		c.Preview,
	}
}

func (d *configurationsDef) onEditItem(item interface{}, done func(result interface{})) {
	c := item.(*mods.Configuration)
	d.createFormItem("Name", c.Name)
	d.createFormMultiLine("Description", c.Description)
	d.createFormItem("Preview", c.Preview)
	d.choicesDef.populate(c.Choices)

	fd := dialog.NewForm("Edit Configuration", "Save", "Cancel", []*widget.FormItem{
		d.getFormItem("Name"),
		d.getFormItem("Description"),
		d.getFormItem("Preview"),
		widget.NewFormItem("Choices", d.choicesDef.draw(false)),
	}, func(ok bool) {
		if ok {
			done(&mods.Configuration{
				Name:        d.getString("Name"),
				Description: d.getString("Description"),
				Preview:     d.getString("Preview"),
				Choices:     d.choicesDef.compile(),
			})
		}
	}, state.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *configurationsDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewLabelWithStyle("Configurations", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Add", func() {
			d.onEditItem(&mods.Configuration{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}

func (d *configurationsDef) set(configurations []*mods.Configuration) {
	d.list.Clear()
	if configurations != nil {
		for _, c := range configurations {
			d.list.AddItem(c)
		}
	}
}
