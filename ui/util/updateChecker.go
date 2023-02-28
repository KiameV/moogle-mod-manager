package util

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/dialog"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
)

func PromptForUpdateAsNeeded(ignoreNoUpdate bool) {
	time.Sleep(time.Second)
	if newer, newerVersion, err := browser.CheckForUpdate(); err != nil {
		if !ignoreNoUpdate {
			dialog.ShowError(err, ui.ActiveWindow())
		}
	} else if newer {
		dialog.ShowConfirm(
			"Update Available",
			fmt.Sprintf("Version %s is available.\nWould you like to update?", newerVersion),
			func(ok bool) {
				if ok {
					_ = browser.Update(newerVersion)
				}
			}, ui.ActiveWindow())
	} else if !ignoreNoUpdate {
		dialog.ShowInformation("No Updates Available", "You are running the latest version.", ui.ActiveWindow())
	}
}
