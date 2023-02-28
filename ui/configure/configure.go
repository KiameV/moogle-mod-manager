package configure

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
)

func Show(w fyne.Window, done func()) {
	configs := *config.Get()
	items := []*widget.FormItem{
		createSelectRow("Default GameDef", &configs.DefaultGame, config.GameIDs()...),
		createCheckboxRow("Check For M3 Updates on Start", configs.CheckForM3UpdateOnStart),
		createCheckboxRow("Delete Downloads After Install", &configs.DeleteDownloadAfterInstall),
	}
	for _, g := range config.GameDefs() {
		var (
			gd *config.GameDir
			ok bool
		)
		if gd, ok = configs.GameDirs[string(g.ID())]; !ok {
			gd = &config.GameDir{}
			configs.GameDirs[string(g.ID())] = gd
		}
		items = append(items, createDirRow(string(g.ID()+" Dir"), &gd.Dir))
	}
	items = append(items, createDirRow("Download Dir", &configs.DownloadDir))
	items = append(items, createDirRow("Backup Dir", &configs.BackupDir))
	items = append(items, createDirRow("Image Cache Dir", &configs.ImgCacheDir))

	d := dialog.NewForm("Configure", "Save", "Cancel", items, func(ok bool) {
		if ok {
			configs.FirstTime = false
			_ = os.MkdirAll(configs.ModsDir, 0777)
			_ = os.MkdirAll(configs.BackupDir, 0777)
			_ = os.MkdirAll(configs.DownloadDir, 0777)
			_ = os.MkdirAll(configs.ImgCacheDir, 0777)
			if err := configs.Save(); err != nil {
				dialog.ShowError(err, w)
				return
			}
			config.Set(configs)
		}
		if done != nil {
			done()
		}
	}, w)
	d.Resize(fyne.NewSize(800, 400))
	d.Show()
}

func createSelectRow(label string, value *string, options ...string) *widget.FormItem {
	if *value == "" {
		*value = options[0]
	}
	sel := widget.NewSelect(options, func(s string) {
		*value = s
	})
	sel.SetSelected(*value)
	return widget.NewFormItem(label, sel)
}

func createCheckboxRow(label string, value *bool, options ...string) *widget.FormItem {
	return widget.NewFormItem(label, widget.NewCheckWithData("", binding.BindBool(value)))
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
