package configure

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"os"
)

func Show(w fyne.Window) {
	configs := *config.Get()
	items := []*widget.FormItem{
		createDirRow("FF I Dir", &configs.DirI),
		createDirRow("FF II Dir", &configs.DirII),
		createDirRow("FF III Dir", &configs.DirIII),
		createDirRow("FF IV Dir", &configs.DirIV),
		createDirRow("FF V Dir", &configs.DirV),
		createDirRow("FF VI Dir", &configs.DirVI),
		createDirRow("Mods Dir", &configs.ModsDir),
		createDirRow("Download Dir", &configs.DownloadDir),
		createDirRow("Backup Dir", &configs.BackupDir),
		createDirRow("Image Cache Dir", &configs.ImgCacheDir),
	}
	d := dialog.NewForm("Configure", "Save", "Cancel", items, func(ok bool) {
		if ok {
			c := config.Get()
			*c = configs
			c.FirstTime = false
			_ = os.MkdirAll(c.ModsDir, 0777)
			_ = os.MkdirAll(c.BackupDir, 0777)
			_ = os.MkdirAll(c.DownloadDir, 0777)
			_ = os.MkdirAll(c.ImgCacheDir, 0777)
			if err := c.Save(); err != nil {
				dialog.ShowError(err, w)
				return
			}
		}
	}, w)
	d.Resize(fyne.NewSize(800, 400))
	d.Show()
}

func createDirRow(label string, value *string) *widget.FormItem {
	b := binding.BindString(value)
	o := &cw.OpenDirDialog{
		IsRelative: false,
		Value:      b,
	}
	o.SetAction(widget.NewToolbarAction(theme.FolderOpenIcon(), o.Handle))
	c := &cw.OpenFileDialogContainer{
		Container:             container.NewBorder(nil, nil, nil, widget.NewToolbar(o), widget.NewEntryWithData(b)),
		OpenFileDialogHandler: o,
	}
	return widget.NewFormItem(label, c.Container)
}
