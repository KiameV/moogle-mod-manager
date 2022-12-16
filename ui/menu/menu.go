package menu

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/configure"
	"github.com/kiamev/moogle-mod-manager/ui/local"
	a "github.com/kiamev/moogle-mod-manager/ui/mod-author"
	"github.com/kiamev/moogle-mod-manager/ui/secret"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"net/url"
)

func New() state.Screen {
	return &MainMenu{}
}

type MainMenu struct{}

func (m *MainMenu) PreDraw(fyne.Window, ...interface{}) error { return nil }

func (m *MainMenu) OnClose() {}

func (m *MainMenu) DrawAsDialog(fyne.Window) {}

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
			configure.Show(w)
		}),
		fyne.NewMenuItem("Secrets", func() {
			secret.Show(w)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Check For App Updates", func() {
			if newer, newerVersion, err := browser.CheckForUpdate(); err != nil {
				dialog.ShowError(err, w)
			} else if newer {
				dialog.ShowConfirm(
					"Update Available",
					fmt.Sprintf("Version %s is available.\nWould you like to update?", newerVersion),
					func(ok bool) {
						_ = browser.Update(newerVersion)
					}, w)
			} else {
				dialog.ShowInformation("No Updates Available", "You are running the latest version.", w)
			}
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("About", func() {
			purl, _ := url.Parse("https://www.patreon.com/kiamev")
			kurl, _ := url.Parse("https://ko-fi.com/kiamev")
			dialog.ShowCustom("About", "ok", container.NewBorder(
				widget.NewRichTextFromMarkdown(fmt.Sprintf(`
## Moogle Mod Manager %s
____________________________
Written by Kiame Vivacity

Contributors:

- Silvris`,
					browser.Version)), nil, nil, nil,
				container.NewVBox(
					widget.NewLabel("If you'd like to support the project:"),
					container.NewHBox(widget.NewLabel("- Patreon"), widget.NewHyperlink("https://www.patreon.com/kiamev", purl)),
					container.NewHBox(widget.NewLabel("- Ko-fi  "), widget.NewHyperlink("https://ko-fi.com/kiamev", kurl))),
			), w)
		}))
	menus = append(menus, file)

	author := fyne.NewMenu("Author")
	newMenu := fyne.NewMenuItem("New", nil)
	newMenu.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Hosted Mod", func() {
			state.GetScreen(state.ModAuthor).(*a.ModAuthorer).NewHostedMod()
			state.ShowScreen(state.ModAuthor)
		}),
		fyne.NewMenuItem("From Nexus", func() {
			state.GetScreen(state.ModAuthor).(*a.ModAuthorer).NewNexusMod()
			state.ShowScreen(state.ModAuthor)
		}),
		fyne.NewMenuItem("From Curseforge", func() {
			state.GetScreen(state.ModAuthor).(*a.ModAuthorer).NewCurseForgeMod()
			state.ShowScreen(state.ModAuthor)
		}),
	)
	if state.GetCurrentGUI() != state.ModAuthor {
		author.Items = append(author.Items,
			newMenu,
			fyne.NewMenuItem("Edit Mod", func() {
				if state.GetScreen(state.ModAuthor).(*a.ModAuthorer).LoadModToEdit() {
					state.ShowScreen(state.ModAuthor)
				}
			}),
			fyne.NewMenuItem("Edit Current Mod", func() {
				if state.GetCurrentGUI() == state.LocalMods {
					if tm := state.GetScreen(state.LocalMods).(local.LocalUI).GetSelected(); tm != nil {
						state.GetScreen(state.ModAuthor).(*a.ModAuthorer).EditMod(tm.Mod(), func(mod *mods.Mod) {
							tm.SetMod(mod)
							if err := tm.Save(); err != nil {
								util.ShowErrorLong(err)
							}
						})
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
