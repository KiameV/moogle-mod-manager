package menu

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/ui/local"
	a "github.com/kiamev/moogle-mod-manager/ui/mod-author"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

func New() state.Screen {
	return &MainMenu{}
}

type MainMenu struct{}

func (m *MainMenu) OnClose() {

}

func (m *MainMenu) Draw(w fyne.Window) {
	file := fyne.NewMenu("File")
	var menus []*fyne.Menu
	if state.GetCurrentGUI() == state.LocalMods {
		file.Items = append(file.Items,
			fyne.NewMenuItem("Select Games", func() {
				state.ShowPreviousScreen()
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
	menus = append(menus, file)

	author := fyne.NewMenu("Author")
	if state.GetCurrentGUI() != state.ModAuthor {
		author.Items = append(author.Items,
			fyne.NewMenuItem("New Mod", func() {
				state.GetScreen(state.ModAuthor).(*a.ModAuthorer).NewMod()
				state.ShowScreen(state.ModAuthor)
			}),
			fyne.NewMenuItem("Edit Mod", func() {
				if state.GetScreen(state.ModAuthor).(*a.ModAuthorer).LoadModToEdit() {
					state.ShowScreen(state.ModAuthor)
				}
			}),
			fyne.NewMenuItem("Edit Current Mod", func() {
				if state.GetCurrentGUI() == state.LocalMods {
					if tm := state.GetScreen(state.LocalMods).(local.LocalUI).GetSelected(); tm != nil {
						state.GetScreen(state.ModAuthor).(*a.ModAuthorer).EditMod(tm.Mod)
						state.ShowScreen(state.ModAuthor)
					}
				}
			}))
	} else {
		author.Items = append(author.Items,
			fyne.NewMenuItem("Close", func() {
				state.GetScreen(state.ModAuthor).OnClose()
				state.ShowPreviousScreen()
			}))
	}
	menus = append(menus, author)
	w.SetMainMenu(fyne.NewMainMenu(menus...))
}
