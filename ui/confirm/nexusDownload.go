package confirm

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	"github.com/kiamev/moogle-mod-manager/mods/nexus"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"net/url"
)

func Nexus(game config.Game, tm *model.TrackedMod, tis []*model.ToInstall, done DownloadCompleteCallback, callback downloadCallback) (err error) {
	var (
		uri string
		u   *url.URL
		c   = container.NewVBox(widget.NewLabelWithStyle("Download the following file from Nexus", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		dll string
	)
	for _, ti := range tis {
		uri = fmt.Sprintf(nexus.NexusFileDownload, ti.Download.Nexus.FileID, nexus.GameToID(game))
		if u, err = url.Parse(uri); err != nil {
			return
		}
		c.Add(widget.NewHyperlink(uri, u))

		c.Add(widget.NewLabelWithStyle("Place download in:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))

		if dll, err = ti.GetDownloadLocation(game, tm); err != nil {
			return
		}
		if u, err = url.Parse(dll); err != nil {
			return
		}
		c.Add(widget.NewHyperlink(dll, u))
	}
	d := dialog.NewCustomConfirm("Download Files", "Done", "Cancel", container.NewVScroll(c), func(ok bool) {
		if ok {
			done(game, tm, tis, callback(game, tm, tis))
		}
	}, state.Window)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
	return
}
