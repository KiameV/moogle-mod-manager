package game_select

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

func New() state.Screen {
	return &GameSelect{}
}

type GameSelect struct{}

func (s *GameSelect) PreDraw(fyne.Window, ...interface{}) error { return nil }

func (s *GameSelect) OnClose() {}

func (s *GameSelect) DrawAsDialog(fyne.Window) {}

func (s *GameSelect) Draw(w fyne.Window) {
	var (
		games  = config.GameDefs()
		inputs = make([]fyne.CanvasObject, 0, len(games)*2-1)
	)
	for _, g := range games {
		inputs = append(inputs, s.createInput(g))
	}
	w.SetContent(container.New(layout.NewGridLayout(3), inputs...))
}

func (s *GameSelect) createInput(g config.GameDef) *fyne.Container {
	return container.NewMax(widget.NewButton("", func() {
		state.CurrentGame = g
		state.ShowScreen(state.LocalMods)
	}), g.Logo())
}
