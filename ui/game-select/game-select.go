package game_select

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/pr-modsync/config"
	"github.com/kiamev/pr-modsync/ui/local"
	"github.com/kiamev/pr-modsync/ui/menu"
	"github.com/kiamev/pr-modsync/ui/state"
)

func Draw(w fyne.Window) {
	menu.Add(w)
	w.SetContent(container.NewGridWithColumns(2,
		container.NewVBox(
			widget.NewLabelWithStyle("Select Game", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewButton("Final Fantasy I", func() {
				state.CurrentGame = toGamePtr(config.I)
				state.CurrentUI = state.LocalMods
				local.Draw(w)
			}),
			widget.NewButton("Final Fantasy II", func() {
				state.CurrentGame = toGamePtr(config.I)
				state.CurrentUI = state.LocalMods
				local.Draw(w)
			}),
			widget.NewButton("Final Fantasy III", func() {
				state.CurrentGame = toGamePtr(config.I)
				state.CurrentUI = state.LocalMods
				local.Draw(w)
			}),
			widget.NewButton("Final Fantasy IV", func() {
				state.CurrentGame = toGamePtr(config.I)
				state.CurrentUI = state.LocalMods
				local.Draw(w)
			}),
			widget.NewButton("Final Fantasy V", func() {
				state.CurrentGame = toGamePtr(config.I)
				state.CurrentUI = state.LocalMods
				local.Draw(w)
			}),
			widget.NewButton("Final Fantasy VI", func() {
				state.CurrentGame = toGamePtr(config.I)
				state.CurrentUI = state.LocalMods
				local.Draw(w)
			})),
		container.NewVBox(),
	))
}

func toGamePtr(game config.Game) *config.Game {
	return &game
}
