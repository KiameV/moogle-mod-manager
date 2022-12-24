package mod_author

import (
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/mod-author/entry"
)

type downloadFilesDef struct {
	entry.Manager
	downloads   *downloads
	selectEntry *entry.SelectFormEntry
	files       *filesDef
	dirs        *dirsDef
}

func newDownloadFilesDef(downloads *downloads, installType *config.InstallType) *downloadFilesDef {
	d := &downloadFilesDef{
		Manager:   entry.NewManager(),
		downloads: downloads,
		files:     newFilesDef(installType),
		dirs:      newDirsDef(installType),
	}
	var i any = entry.NewSelectEntry(d, "Download Name", "", nil)
	d.selectEntry = i.(*entry.SelectFormEntry)
	return d
}

func (d *downloadFilesDef) compile() *mods.DownloadFiles {
	return &mods.DownloadFiles{
		DownloadName: entry.Value[string](d, "Download Name"),
		Files:        d.files.compile(),
		Dirs:         d.dirs.compile(),
	}
}

/*func (d *downloadFilesDef) draw() fyne.CanvasObject {
	var possible []string
	for _, dl := range d.downloads.compileDownloads() {
		possible = append(possible, dl.Name)
	}

	entry.NewEntry[string](d, entry.KindSelect, "Download Name", possible, d.dlName)

	return container.NewVBox(
		widget.NewForm(entry.FormItem[string](d, "Download Name")),
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
	d.selectEntry.Entry.Options = possible

	return []*widget.FormItem{
		entry.FormItem[string](d, "Download Name"),
		widget.NewFormItem("Files", d.files.draw(false)),
		widget.NewFormItem("Dirs", d.dirs.draw(false)),
	}, nil
}

func (d *downloadFilesDef) clear() {
	d.selectEntry.Set("")
	d.files.clear()
	d.dirs.clear()
}

func (d *downloadFilesDef) populate(dlf *mods.DownloadFiles) {
	if dlf == nil {
		d.clear()
	} else {
		d.selectEntry.Set(dlf.DownloadName)
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
