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
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"strings"
)

type DownloadCompleteCallback func(game config.Game, tm *model.TrackedMod, tis []*mods.ToInstall, err error)

type downloadCallback func(game config.Game, downloadDir string, tm *model.TrackedMod, tis []*mods.ToInstall) error

func Hosted(game config.Game, downloadDir string, tm *model.TrackedMod, tis []*mods.ToInstall, done DownloadCompleteCallback, callback downloadCallback) {
	sb := strings.Builder{}
	for i, ti := range tis {
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
			done(game, tm, tis, callback(game, downloadDir, tm, tis))
		}
	}, state.Window)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}
