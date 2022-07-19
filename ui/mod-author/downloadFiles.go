package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type downloadFilesDef struct {
	selected  string
	downloads *downloadsDef
	files     *filesDef
	dirs      *dirsDef
}

func newDownloadFilesDef(downloads *downloadsDef) *downloadFilesDef {
	return &downloadFilesDef{
		downloads: downloads,
	}
}

func (d *downloadFilesDef) compile() *mods.DownloadFiles {
	return &mods.DownloadFiles{
		DownloadName: getFormString("dfDlName"),
		Files:        d.files.compile(),
		Dirs:         d.dirs.compile(),
	}
}

func (d *downloadFilesDef) draw() fyne.CanvasObject {
	var possible []string
	for _, dl := range d.downloads.compile() {
		possible = append(possible, dl.Name)
	}

	setFormSelect("dfDlName", possible, d.selected)

	return container.NewVBox(
		widget.NewForm(getFormItem("Name", "dfDlName")),
		d.files.draw(),
		d.dirs.draw(),
	)
}
