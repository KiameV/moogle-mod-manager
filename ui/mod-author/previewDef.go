package mod_author

import (
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type previewDef struct {
	*entryManager
}

func newPreviewDef() *previewDef {
	d := &previewDef{
		entryManager: newEntryManager(),
	}
	d.createFormItem("Preview Url", "")
	//d.createFileDialog("Preview Local", "", state.GetBaseDirBinding(), true, true)
	//d.createFormItem("Size X", "")
	//d.createFormItem("Size Y", "")
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
	d.createFormItem("Preview Url", url)
	//d.createFileDialog("Preview Local", local, state.GetBaseDirBinding(), true, true)
	//d.createFormItem("Size X", fmt.Sprintf("%d", p.Size.X))
	//d.createFormItem("Size Y", fmt.Sprintf("%d", p.Size.Y))
}

/*func (d *previewDef) draw() fyne.CanvasObject {
	return widget.NewForm(d.getFormItems()...)
}*/

func (d *previewDef) compile() *mods.Preview {
	var (
		p = &mods.Preview{
			/*Size: &mods.Size{
				X: d.getInt("Size X"),
				Y: d.getInt("Size Y"),
			},*/
		}
		url = d.getString("Preview Url")
		//local = d.getString("Preview Local")
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
		d.getFormItem("Preview Url"),
		//d.getFileDialog("Preview Local"),
		//d.getFormItem("Size X"),
		//d.getFormItem("Size Y"),
	}
}
