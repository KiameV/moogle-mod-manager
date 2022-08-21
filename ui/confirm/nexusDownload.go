package confirm

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/nexus"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"net/url"
)

func Nexus(enabler *mods.ModEnabler, competeCallback DownloadCompleteCallback, done downloadCallback) (err error) {
	var (
		uri string
		u   *url.URL
		c   = container.NewVBox(widget.NewLabelWithStyle("Download the following file from Nexus", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		dll string
	)
	for _, ti := range enabler.ToInstall {
		uri = fmt.Sprintf(nexus.NexusFileDownload, ti.Download.Nexus.FileID, nexus.GameToID(enabler.Game))
		if u, err = url.Parse(uri); err != nil {
			return
		}
		c.Add(widget.NewHyperlink(uri, u))

		c.Add(widget.NewLabelWithStyle("Place download in:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))

		if dll, err = ti.GetDownloadLocation(enabler.Game, enabler.TrackedMod); err != nil {
			return
		}
		if u, err = url.Parse(dll); err != nil {
			return
		}
		c.Add(widget.NewHyperlink(dll, u))
	}
	d := dialog.NewCustomConfirm("Download Files", "Done", "Cancel", container.NewVScroll(c), func(ok bool) {
		if ok {
			done(enabler, competeCallback, err)
		}
	}, state.Window)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
	return
}
