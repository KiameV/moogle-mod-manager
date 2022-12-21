package mod_author

import (
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
)

type previewDef struct {
	e entry.Entry[string]
}

func newPreviewDef() *previewDef {
	return &previewDef{
		e: entry.NewStringFormEntry("Preview Url", ""),
	}
}

func (d *previewDef) set(p *mods.Preview) {
	if p == nil || p.Url == nil {
		d.e.Set("")
	} else {
		d.e.Set(*p.Url)
	}
}

func (d *previewDef) compile() *mods.Preview {
	var p mods.Preview
	if url := d.e.Value(); url != "" {
		p.Url = &url
	}
	return &p
}

func (d *previewDef) getFormItems() []*widget.FormItem {
	return []*widget.FormItem{
		d.e.FormItem(),
	}
}
