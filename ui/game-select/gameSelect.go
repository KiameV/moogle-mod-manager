package game_select

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

func New() state.Screen {
	return &GameSelect{}
}

type GameSelect struct{}

func (s *GameSelect) OnClose() {
	
}

func (s *GameSelect) Draw(w fyne.Window) {
	w.SetContent(container.NewGridWithColumns(2,
		container.NewVBox(
			widget.NewLabelWithStyle("Select Games", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewButton("Final Fantasy I", func() {
				state.CurrentGame = toGamePtr(config.I)
				state.ShowScreen(state.LocalMods)
			}),
			widget.NewButton("Final Fantasy II", func() {
				state.CurrentGame = toGamePtr(config.II)
				state.ShowScreen(state.LocalMods)
			}),
			widget.NewButton("Final Fantasy III", func() {
				state.CurrentGame = toGamePtr(config.III)
				state.ShowScreen(state.LocalMods)
			}),
			widget.NewButton("Final Fantasy IV", func() {
				state.CurrentGame = toGamePtr(config.IV)
				state.ShowScreen(state.LocalMods)
			}),
			widget.NewButton("Final Fantasy V", func() {
				state.CurrentGame = toGamePtr(config.V)
				state.ShowScreen(state.LocalMods)
			}),
			widget.NewButton("Final Fantasy VI", func() {
				state.CurrentGame = toGamePtr(config.VI)
				state.ShowScreen(state.LocalMods)
			})),
		container.NewVBox(),
	))
}

func toGamePtr(game config.Game) *config.Game {
	return &game
}
