package util

import (
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"strings"
)

func DisplayDownloadsAndFiles(mod *mods.Mod, toInstall []*mods.DownloadFiles) {
	sb := strings.Builder{}
	for dl, dlf := range compileDownloadFiles(mod, toInstall) {
		sb.WriteString(fmt.Sprintf("Download: %s\n\n", dl.Name))
		sb.WriteString("  Sources:\n\n")
		for _, s := range dl.Sources {
			sb.WriteString(fmt.Sprintf("  - %s\n\n", s))
		}
		sb.WriteString("  Files:\n\n")
		for _, f := range dlf.Files {
			sb.WriteString(fmt.Sprintf("  - %s -> %s\n\n", f.From, f.To))
		}
		sb.WriteString("  Dirs:\n\n")
		for _, dir := range dlf.Dirs {
			sb.WriteString(fmt.Sprintf("  - %s -> %s | Recursive %v\n\n", dir.From, dir.To, dir.Recursive))
		}
		break
	}
	dialog.ShowCustom("Downloads and File/Dir Copies", "ok", widget.NewRichTextFromMarkdown(sb.String()), state.Window)
	state.ShowPreviousScreen()
}

func compileDownloadFiles(mod *mods.Mod, toInstall []*mods.DownloadFiles) map[*mods.Download]*mods.DownloadFiles {
	dlf := make(map[*mods.Download]*mods.DownloadFiles)
	dl := make(map[string]*mods.Download)
	for _, d := range mod.Downloadables {
		dl[d.Name] = d
	}
	add(mod.DownloadFiles, dl, dlf)
	for _, ti := range toInstall {
		add(ti, dl, dlf)
	}
	return dlf
}

func add(ti *mods.DownloadFiles, dl map[string]*mods.Download, dlf map[*mods.Download]*mods.DownloadFiles) {
	var (
		d  *mods.Download
		f  *mods.DownloadFiles
		ok bool
	)
	if ti == nil {
		return
	}
	if d, ok = dl[ti.DownloadName]; !ok {
		return
	}
	if f, ok = dlf[d]; !ok {
		f = &mods.DownloadFiles{DownloadName: ti.DownloadName}
		dlf[d] = f
	}
	for _, df := range ti.Files {
		f.Files = append(f.Files, df)
	}
	for _, dd := range ti.Dirs {
		f.Dirs = append(f.Dirs, dd)
	}
}
