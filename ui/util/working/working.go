package working

import (
	"fyne.io/fyne/v2/dialog"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
)

var workingDialog dialog.Dialog

func ShowDialog() {
	if workingDialog == nil {
		if w := ui.ActiveWindow(); w != nil {
			workingDialog = dialog.NewInformation("Working", "Working...", w)
			workingDialog.Show()
		}
	}
}

func HideDialog() {
	if workingDialog != nil {
		workingDialog.Hide()
		workingDialog = nil
	}
}
