package confirm

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/discover/remote/nexus"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"net/url"
	"os"
	"path/filepath"
)

type toDownload struct {
	uri string
	dir string
}

type nexusConfirmer struct {
	Params
}

func newNexusConfirmer(params Params) Confirmer {
	return &nexusConfirmer{Params: params}
}

func (c *nexusConfirmer) Downloads(done func(mods.Result)) (err error) {
	var (
		toDl []toDownload
	)
	for _, ti := range c.ToInstall {
		if ti.Download != nil {
			dl := toDownload{
				uri: fmt.Sprintf(nexus.NexusFileDownload, ti.Download.Nexus.FileID, c.Game.Remote().Nexus.ID),
			}
			if dl.dir, err = ti.GetDownloadLocation(c.Game, c.Mod); err != nil {
				return
			}
			if _, err = os.Stat(filepath.Join(dl.dir, ti.Download.Nexus.FileName)); err == nil {
				continue
			}
			toDl = append(toDl, dl)
		}
	}

	if len(toDl) == 0 {
		done(mods.Ok)
		return nil
	}

	return c.showDialog(toDl, done)
}

func (c *nexusConfirmer) showDialog(toDl []toDownload, done func(mods.Result)) (err error) {
	var (
		u  *url.URL
		vb = container.NewVBox(widget.NewLabelWithStyle("Download the following file from Nexus", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	)
	for _, td := range toDl {
		if u, err = url.Parse(td.uri); err != nil {
			return
		}
		vb.Add(widget.NewHyperlink(td.uri, u))

		vb.Add(widget.NewLabelWithStyle("Place download in:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))

		if u, err = url.Parse(td.dir); err != nil {
			return
		}
		vb.Add(widget.NewHyperlink(td.dir, u))
	}
	d := dialog.NewCustomConfirm("Download Files", "Done", "Cancel", container.NewVScroll(vb), func(ok bool) {
		result := mods.Ok
		if !ok {
			result = mods.Cancel
		}
		done(result)
	}, ui.Window)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
	return
}
