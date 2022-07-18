package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/kiamev/pr-modsync/mods/managed"
	"github.com/kiamev/pr-modsync/ui/game-select"
	"github.com/kiamev/pr-modsync/ui/state"
)

func main() {
	state.App = app.New()
	state.Window = state.App.NewWindow("Moogle Mod Manager")
	state.Window.Resize(fyne.NewSize(800, 600))
	if err := managed.Initialize(); err != nil {
		dialog.ShowError(err, state.Window)
		state.Window.Close()
		return
	}
	game_select.Draw(state.Window)
	state.Window.ShowAndRun()
}
