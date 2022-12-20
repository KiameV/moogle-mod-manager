package mod_author

import (
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type downloadFilesDef struct {
	*entryManager
	downloads *downloads
	dlName    string
	files     *filesDef
	dirs      *dirsDef
}

func newDownloadFilesDef(downloads *downloads) *downloadFilesDef {
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

/*func (d *downloadFilesDef) draw() fyne.CanvasObject {
	var possible []string
	for _, dl := range d.downloads.compileDownloads() {
		possible = append(possible, dl.Name)
	}

	d.createFormSelect("Download Name", possible, d.dlName)

	return container.NewVBox(
		widget.NewForm(d.getFormItem("Download Name")),
		d.files.draw(true),
		d.dirs.draw(true),
	)
}*/

func (d *downloadFilesDef) getFormItems() ([]*widget.FormItem, error) {
	var (
		possible []string
		dls, err = d.downloads.compileDownloads()
	)
	if err != nil {
		return nil, err
	}
	for _, dl := range dls {
		possible = append(possible, dl.Name)
	}

	d.createFormSelect("Download Name", possible, d.dlName)

	return []*widget.FormItem{
		d.getFormItem("Download Name"),
		widget.NewFormItem("Files", d.files.draw(false)),
		widget.NewFormItem("Dirs", d.dirs.draw(false)),
	}, nil
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

/*func (d *downloadFilesDef) set(df *mods.DownloadFiles) {
	d.dlName = ""
	d.files.clear()
	d.dirs.clear()
	if df != nil {
		d.dlName = df.DownloadName
		d.files.populate(df.Files)
		d.dirs.populate(df.Dirs)
	}
}*/
