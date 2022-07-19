package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	"github.com/kiamev/moogle-mod-manager/ui/game-select"
	"github.com/kiamev/moogle-mod-manager/ui/local"
	"github.com/kiamev/moogle-mod-manager/ui/menu"
	mod_author "github.com/kiamev/moogle-mod-manager/ui/mod-author"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

func main() {
	state.App = app.New()
	state.Window = state.App.NewWindow("Moogle Mod Manager " + browser.Version)
	state.Window.Resize(fyne.NewSize(800, 800))
	if err := managed.Initialize(); err != nil {
		dialog.ShowError(err, state.Window)
		state.Window.Close()
		return
	}

	state.RegisterMainMenu(menu.New())
	state.RegisterScreen(state.None, game_select.New())
	state.RegisterScreen(state.ModAuthor, mod_author.New())
	state.RegisterScreen(state.LocalMods, local.New())

	state.ShowScreen(state.None)
	state.Window.ShowAndRun()
}
