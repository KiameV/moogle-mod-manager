package confirm

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"strings"
)

type DownloadCompleteCallback func(enabler *mods.ModEnabler, err error)

type downloadCallback func(enabler *mods.ModEnabler, completeCallback DownloadCompleteCallback, err error)

func Hosted(enabler *mods.ModEnabler, completeCallback DownloadCompleteCallback, done downloadCallback) {
	var (
		sb  = strings.Builder{}
		err error
	)
	for i, ti := range enabler.ToInstall {
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
	d := dialog.NewCustomConfirm("Download Files?", "Yes", "Cancel", container.NewVScroll(widget.NewRichTextFromMarkdown(sb.String())), func(ok bool) {
		if ok {
			done(enabler, completeCallback, err)
		}
	}, state.Window)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}
