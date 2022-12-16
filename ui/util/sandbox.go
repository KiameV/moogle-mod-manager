package util

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"strings"
)

func DisplayDownloadsAndFiles(toInstall []*mods.ToInstall) {
	defer func() {
		if err := recover(); err != nil {
			ShowErrorLong(fmt.Errorf("%v", err))
		}
	}()
	sb := strings.Builder{}
	for _, ti := range toInstall {
		sb.WriteString(fmt.Sprintf("## %s\n\n", ti.Download.Name))
		sb.WriteString("### Sources:\n\n")
		if ti.Download.Hosted != nil {
			for _, s := range ti.Download.Hosted.Sources {
				sb.WriteString(fmt.Sprintf("  - %s\n\n", s))
			}
		} else if ti.Download.Nexus != nil {
			sb.WriteString(fmt.Sprintf("  - %s\n\n", ti.Download.Nexus.FileName))
		} else if ti.Download.CurseForge != nil {
			sb.WriteString(" - ")
		}
		sb.WriteString("### Files:\n\n")
		for _, dl := range ti.DownloadFiles {
			for _, f := range dl.Files {
				sb.WriteString(fmt.Sprintf("  - %s -> %s\n\n", f.From, f.To))
			}
			sb.WriteString("### Dirs:\n\n")
			for _, dir := range dl.Dirs {
				sb.WriteString(fmt.Sprintf("  - %s -> %s | Recursive %v\n\n", dir.From, dir.To, dir.Recursive))
			}
		}
		sb.WriteString("_____________________\n\n")
	}
	d := dialog.NewCustom("Downloads and File/Dir Copies", "ok", container.NewVScroll(widget.NewRichTextFromMarkdown(sb.String())), ui.Window)
	d.Resize(fyne.NewSize(600, 600))
	d.Show()
}
