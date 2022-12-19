package confirm

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"os"
	"path/filepath"
	"strings"
)

type hostedConfirmer struct {
	Params
}

func newHostedConfirmer(params Params) Confirmer {
	return &hostedConfirmer{Params: params}
}

func (c *hostedConfirmer) Downloads(done func(mods.Result)) (err error) {
	var sb = strings.Builder{}
	for i, ti := range c.ToInstall {
		if c.alreadyDownloaded(ti) {
			continue
		}
		sb.WriteString(fmt.Sprintf("## Download %d\n\n", i+1))
		if len(ti.Download.Hosted.Sources) == 1 {
			sb.WriteString(ti.Download.Hosted.Sources[0] + "\n\n")
		} else {
			sb.WriteString("### Sources:\n\n")
			for j, s := range ti.Download.Hosted.Sources {
				sb.WriteString(fmt.Sprintf(" - %d. %s\n\n", j+1, s))
			}
		}
	}
	if sb.Len() == 0 {
		done(mods.Ok)
		return
	}

	d := dialog.NewCustomConfirm("Download Files?", "Yes", "Cancel", container.NewVScroll(widget.NewRichTextFromMarkdown(sb.String())), func(ok bool) {
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

func (c *hostedConfirmer) alreadyDownloaded(ti *mods.ToInstall) bool {
	file := strings.Split(ti.Download.Hosted.Sources[0], "/")
	file = strings.Split(file[len(file)-1], "?")
	dir, _ := ti.GetDownloadLocation(c.Game, c.Mod)
	_, err := os.Stat(filepath.Join(dir, file[0]))
	return err == nil
}
