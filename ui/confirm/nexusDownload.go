package confirm

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	"github.com/kiamev/moogle-mod-manager/mods/nexus"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"net/url"
)

func Nexus(game config.Game, downloadDir string, tm *model.TrackedMod, tis []*mods.ToInstall, done DownloadCompleteCallback, callback downloadCallback) error {
	c := container.NewVBox(widget.NewLabelWithStyle("Please download the following files from Nexus", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	for _, ti := range tis {
		uri := fmt.Sprintf(nexus.NexusFileDownload, ti.Download.Nexus.FileID, nexus.IdFromGame(game))
		u, err := url.Parse(uri)
		if err != nil {
			return err
		}
		c.Add(widget.NewHyperlink(uri, u))
	}
	c.Add(widget.NewLabel("Please place all downloads in"))
	u, err := url.Parse(downloadDir)
	if err != nil {
		return err
	}
	c.Add(widget.NewHyperlink(downloadDir, u))
	d := dialog.NewCustomConfirm("Download Files", "Done", "Cancel", container.NewVScroll(c), func(ok bool) {
		if ok {
			done(game, tm, tis, callback(game, downloadDir, tm, tis))
		}
	}, state.Window)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
	return nil
}
