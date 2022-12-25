package confirm

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
)

func ShowEnableModConfirmDialog(baseModName mods.ModName, neededMod *mods.Mod, done func(mods.Result)) {
	msg := fmt.Sprintf("[%s] requires [%s], would you like to enable it first?", baseModName, neededMod.Name)
	d := dialog.NewCustomConfirm("Enable Required Mod?", "Yes", "Cancel",
		container.NewVScroll(widget.NewRichTextFromMarkdown(msg)), func(ok bool) {
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
