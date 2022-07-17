package menu

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/kiamev/pr-modsync/browser"
	a "github.com/kiamev/pr-modsync/ui/author"
	"github.com/kiamev/pr-modsync/ui/state"
)

func Add(w fyne.Window) {
	file := fyne.NewMenu("File")
	if state.CurrentUI == state.LocalMods {
		file.Items = append(file.Items,
			fyne.NewMenuItem("Add Mod From URL", func() {

			}),
			fyne.NewMenuItem("Add Mod From File", func() {

			}),
			fyne.NewMenuItemSeparator())
	}
	file.Items = append(file.Items,
		fyne.NewMenuItem("Configure", func() {
			/*widget.NewButton("Dark", func() {
				a.Settings().SetTheme(theme.DarkTheme())
			}),
			widget.NewButton("Light", func() {
				a.Settings().SetTheme(theme.LightTheme())
			}),*/
		}),
		fyne.NewMenuItem("Check For Updates", func() {
			if newer, newerVersion, err := browser.CheckForUpdate(); err != nil {
				dialog.ShowError(err, w)
			} else if newer {
				dialog.ShowConfirm(
					"Update Available",
					fmt.Sprintf("Version %s is available.\nWould you like to update?", newerVersion),
					func(ok bool) {
						browser.Update(newerVersion)
					}, w)
			} else {
				dialog.ShowInformation("No Updates Available", "You are running the latest version.", w)
			}
		}))
	author := fyne.NewMenu("Author")
	author.Items = append(author.Items,
		fyne.NewMenuItem("New Mod", func() {
			a.DrawNewMod(w)
		}))
	w.SetMainMenu(fyne.NewMainMenu(file, author))
}
