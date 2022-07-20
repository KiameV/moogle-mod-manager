package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type downloadFilesDef struct {
	*entryManager
	downloads *downloadsDef
	dlName    string
	files     *filesDef
	dirs      *dirsDef
}

func newDownloadFilesDef(downloads *downloadsDef) *downloadFilesDef {
	return &downloadFilesDef{
		entryManager: newEntryManager(),
		downloads:    downloads,
		files:        newFilesDef(),
		dirs:         newDirsDef(),
	}
}

func (d *downloadFilesDef) compile() *mods.DownloadFiles {
	return &mods.DownloadFiles{
		DownloadName: d.getString("Download Name"),
		Files:        d.files.compile(),
		Dirs:         d.dirs.compile(),
	}
}

func (d *downloadFilesDef) draw() fyne.CanvasObject {
	var possible []string
	for _, dl := range d.downloads.compile() {
		possible = append(possible, dl.Name)
	}

	d.createFormSelect("Download Name", possible, d.dlName)

	return container.NewVBox(
		widget.NewForm(d.getFormItem("Download Name")),
		d.files.draw(true),
		d.dirs.draw(true),
	)
}

func (d *downloadFilesDef) drawAsFormItems() []*widget.FormItem {
	var possible []string
	for _, dl := range d.downloads.compile() {
		possible = append(possible, dl.Name)
	}

	d.createFormSelect("Download Name", possible, d.dlName)

	return []*widget.FormItem{
		d.getFormItem("Download Name"),
		widget.NewFormItem("Files", d.files.draw(false)),
		widget.NewFormItem("Dirs", d.dirs.draw(false)),
	}
}

func (d *downloadFilesDef) clear() {
	d.dlName = ""
	d.files.clear()
	d.dirs.clear()
}

func (d *downloadFilesDef) populate(dlf *mods.DownloadFiles) {
	if dlf == nil {
		d.clear()
	} else {
		d.dlName = dlf.DownloadName
		d.files.populate(dlf.Files)
		d.dirs.populate(dlf.Dirs)
	}
}
