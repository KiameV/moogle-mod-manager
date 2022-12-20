package mod_author

import (
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
)

type previewDef struct {
	entry.Manager
}

func newPreviewDef() *previewDef {
	d := &previewDef{
		Manager: entry.NewManager(),
	}
	entry.NewEntry[string](d, entry.KindString, "Preview Url", "")
	//d.createFileDialog("Preview Local", "", state.GetBaseDirBinding(), true, true)
	//entry.NewEntry[string](d, entry.KindString, "Size X", "")
	//entry.NewEntry[string](d, entry.KindString, "Size Y", "")
	return d
}

func (d *previewDef) set(p *mods.Preview) {
	if p == nil {
		p = &mods.Preview{}
	}
	var url string
	if p.Url != nil {
		url = *p.Url
	}
	/*var local string
	if p.Local != nil {
		local = *p.Local
	}*/
	entry.NewEntry[string](d, entry.KindString, "Preview Url", url)
	//d.createFileDialog("Preview Local", local, state.GetBaseDirBinding(), true, true)
	//entry.NewEntry[string](d, entry.KindString, "Size X", fmt.Sprintf("%d", p.Size.X))
	//entry.NewEntry[string](d, entry.KindString, "Size Y", fmt.Sprintf("%d", p.Size.Y))
}

/*func (d *previewDef) draw() fyne.CanvasObject {
	return widget.NewForm(d.getFormItems()...)
}*/

func (d *previewDef) compile() *mods.Preview {
	var (
		p = &mods.Preview{
			/*Size: &mods.Size{
				X: d.GetInt("Size X"),
				Y: d.GetInt("Size Y"),
			},*/
		}
		url = entry.Value[string](d, "Preview Url")
		//local = entry.Value[string](d, "Preview Local")
	)
	if url != "" {
		p.Url = &url
	}
	//if local != "" {
	//	p.Local = &local
	//}
	if p.Url == nil && p.Local == nil {
		p = nil
	}
	return p
}

func (d *previewDef) getFormItems() []*widget.FormItem {
	return []*widget.FormItem{
		entry.FormItem[string](d, "Preview Url"),
		//d.GetFileDialog("Preview Local"),
		//entry.FormItem[string](d, "Size X"),
		//entry.FormItem[string](d, "Size Y"),
	}
}
