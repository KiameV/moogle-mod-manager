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

type choicesDef struct {
	*entryManager
	list       *cw.DynamicList
	dlfDef     *downloadFilesDef
	configDef  *configurationsDef
	previewDef *previewDef
}

func newChoicesDef(dlDef *downloads, configDef *configurationsDef) *choicesDef {
	d := &choicesDef{
		entryManager: newEntryManager(),
		dlfDef:       newDownloadFilesDef(dlDef),
		configDef:    configDef,
		previewDef:   newPreviewDef(),
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	}, true)
	return d
}

func (d *choicesDef) compile() []*mods.Choice {
	choices := make([]*mods.Choice, len(d.list.Items))
	for i, item := range d.list.Items {
		choices[i] = item.(*mods.Choice)
	}
	return choices
}

func (d *choicesDef) getItemKey(item interface{}) string {
	return item.(*mods.Choice).Name
}

func (d *choicesDef) getItemFields(item interface{}) []string {
	c := item.(*mods.Choice)
	sl := []string{
		c.Name,
		c.Description,
	}
	if c.NextConfigurationName != nil {
		sl = append(sl, *c.NextConfigurationName)
	}
	if c.DownloadFiles != nil {
		sl = append(sl, c.DownloadFiles.DownloadName)
	}
	return sl
}

func (d *choicesDef) onEditItem(item interface{}) {
	d.createItem(item)
}

func (d *choicesDef) createItem(item interface{}, done ...func(interface{})) {
	var (
		c          = item.(*mods.Choice)
		configs    = d.configDef.compile()
		possible   = d.getPossibleConfigs(configs)
		nextConfig = ""
	)
	d.dlfDef.populate(c.DownloadFiles)

	if c.NextConfigurationName != nil {
		nextConfig = *c.NextConfigurationName
	}

	d.createFormItem("Name", c.Name)
	d.createFormItem("Description", c.Description)
	d.createFormSelect("Next Configuration", possible, nextConfig)
	d.previewDef.set(c.Preview)
	if c.DownloadFiles != nil {
		d.dlfDef.populate(c.DownloadFiles)
	}

	form := []*widget.FormItem{
		d.getFormItem("Name"),
		d.getFormItem("Description"),
		d.getFormItem("Next Configuration"),
	}
	form = append(form, d.previewDef.getFormItems()...)

	dls, err := d.dlfDef.getFormItems()
	if err != nil {
		dialog.ShowError(err, ui.Window)
	} else {
		form = append(form, dls...)
	}

	fd := dialog.NewForm("Edit Choice", "Save", "Cancel", form, func(ok bool) {
		if ok {
			c.Name = d.getString("Name")
			c.Description = d.getString("Description")
			c.Preview = d.previewDef.compile()
			c.DownloadFiles = d.dlfDef.compile()
			if d.getString("Next Configuration") != "" {
				s := d.getString("Next Configuration")
				c.NextConfigurationName = &s
			} else {
				c.NextConfigurationName = nil
			}
			if len(done) > 0 {
				done[0](c)
			}
			d.list.Refresh()
		}
	}, ui.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *choicesDef) draw(includeLabel bool) fyne.CanvasObject {
	c := container.NewVBox()
	if includeLabel {
		c.Add(widget.NewLabelWithStyle("Choices", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	c.Add(widget.NewButton("Add", func() {
		d.createItem(&mods.Choice{}, func(result interface{}) {
			d.list.AddItem(result)
		})
	}))
	c.Add(d.list.Draw())
	return c
}

func (d *choicesDef) getPossibleConfigs(configs []*mods.Configuration) (possible []string) {
	if len(configs) > 0 {
		possible = make([]string, len(configs)+1)
		possible[0] = ""
		for i, cfg := range d.configDef.compile() {
			possible[i+1] = cfg.Name
		}
	}
	return
}

func (d *choicesDef) populate(choices []*mods.Choice) {
	d.list.Clear()
	for _, c := range choices {
		d.list.AddItem(c)
	}
}
