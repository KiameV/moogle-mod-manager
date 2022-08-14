package game_select

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/util/resources"
)

func New() state.Screen {
	return &GameSelect{}
}

type GameSelect struct{}

func (ui *GameSelect) PreDraw() error { return nil }

func (s *GameSelect) OnClose() {}

func (s *GameSelect) Draw(w fyne.Window) {
	w.SetContent(container.NewCenter(
		container.NewVBox(
			container.NewMax(resources.LogoI, widget.NewButton("", func() {
				state.CurrentGame = toGamePtr(config.I)
				state.ShowScreen(state.LocalMods)
			})),
			widget.NewSeparator(),
			container.NewMax(resources.LogoII, widget.NewButton("", func() {
				state.CurrentGame = toGamePtr(config.II)
				state.ShowScreen(state.LocalMods)
			})),
			widget.NewSeparator(),
			container.NewMax(resources.LogoIII, widget.NewButton("", func() {
				state.CurrentGame = toGamePtr(config.III)
				state.ShowScreen(state.LocalMods)
			})),
			widget.NewSeparator(),
			container.NewMax(resources.LogoIV, widget.NewButton("", func() {
				state.CurrentGame = toGamePtr(config.IV)
				state.ShowScreen(state.LocalMods)
			})),
			widget.NewSeparator(),
			container.NewMax(resources.LogoV, widget.NewButton("", func() {
				state.CurrentGame = toGamePtr(config.V)
				state.ShowScreen(state.LocalMods)
			})),
			widget.NewSeparator(),
			container.NewMax(resources.LogoVI, widget.NewButton("", func() {
				state.CurrentGame = toGamePtr(config.VI)
				state.ShowScreen(state.LocalMods)
			}))),
	))
}

func toGamePtr(game config.Game) *config.Game {
	return &game
}
