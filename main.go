package main

import (
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
	"github.com/aarzilli/nucular/style"
	"github.com/kiamev/pr-modsync/browser"
	"github.com/kiamev/pr-modsync/config"
	"github.com/kiamev/pr-modsync/mods"
	"github.com/kiamev/pr-modsync/ui/local"
	state "github.com/kiamev/pr-modsync/ui/state"
	"image"
	"image/color"
	"time"
)

const version = "0.1.0"

var (
	status        string
	statusTimer   *time.Timer
	errTextEditor nucular.TextEditor
	wnd           nucular.MasterWindow
)

func main() {
	errTextEditor.Flags = nucular.EditReadOnly | nucular.EditSelectable | nucular.EditSelectable | nucular.EditMultiline
	var (
		x = config.Get().WindowX
		y = config.Get().WindowY
	)
	if x == 0 || y == 0 {
		x = config.WindowWidth
		y = config.WindowHeight
	}
	wnd = nucular.NewMasterWindowSize(0, "FF PR Mod Manager - "+version, image.Point{X: x, Y: y}, updateWindow)
	wnd.SetStyle(style.FromTable(customTheme, 1.2))
	wnd.Main()
}

func updateWindow(w *nucular.Window) {
	var err error
	w.MenubarBegin()
	w.Row(12).Static(100, 100, 300, 200)
	if w := w.Menu(label.TA("Game", "LC"), 100, nil); w != nil {
		w.Row(12).Dynamic(1)
		if w.MenuItem(label.TA("I", "LC")) {
			state.Game = config.I
			state.CurrentUI = state.LocalMods
			w.Close()
		} else if w.MenuItem(label.TA("II", "LC")) {
			state.Game = config.II
			state.CurrentUI = state.LocalMods
			w.Close()
		} else if w.MenuItem(label.TA("III", "LC")) {
			state.Game = config.III
			state.CurrentUI = state.LocalMods
			w.Close()
		} else if w.MenuItem(label.TA("IV", "LC")) {
			state.Game = config.IV
			state.CurrentUI = state.LocalMods
			w.Close()
		} else if w.MenuItem(label.TA("V", "LC")) {
			state.Game = config.V
			state.CurrentUI = state.LocalMods
			w.Close()
		} else if w.MenuItem(label.TA("VI", "LC")) {
			state.Game = config.VI
			state.CurrentUI = state.LocalMods
			w.Close()
		}
	}

	if state.Game != config.None && state.CurrentUI == state.LocalMods {
		if w.MenuItem(label.TA("Add Mod", "LC")) {
			w.Close()
		}
	} else {
		w.Spacing(1)
	}

	if w := w.Menu(label.TA("Check For Update", "LC"), 300, nil); w != nil {
		var hasNewer bool
		var latest string
		if hasNewer, latest, err = browser.CheckForUpdate(version); err != nil {
			popupErr(w, err)
		}
		if hasNewer {
			browser.Update(latest)
		} else {
			status = "version is current"
		}
		w.Close()
	}

	popupErr(w, state.Errors...)

	if status != "" {
		w.Label("Status: "+status, "RC")
		if statusTimer != nil {
			statusTimer.Stop()
		}
		statusTimer = time.AfterFunc(2*time.Second, func() { status = "" })

	} else {
		w.Spacing(1)
	}
	w.MenubarEnd()

	if state.Game != config.None {
		var gm mods.GameMods
		if gm, err = mods.GetGameMods(state.Game); err != nil {
			popupErr(w, err)
			return
		}
		switch state.CurrentUI {
		case state.LocalMods:
			local.Draw(w, gm)
		case state.RemoteMods:
		case state.Config:
		case state.None:
			break
		}
	}
}

func popupErr(w *nucular.Window, errs ...error) {
	if len(errs) == 0 {
		return
	}
	for _, err := range errs {
		errTextEditor.Buffer = []rune(err.Error())
		w.Master().PopupOpen("Error", nucular.WindowMovable|nucular.WindowTitle|nucular.WindowDynamic, rect.Rect{X: 20, Y: 100, W: 700, H: 600}, true,
			func(w *nucular.Window) {
				w.Row(300).Dynamic(1)
				errTextEditor.Edit(w)
				w.Row(25).Dynamic(1)
				if w.Button(label.T("OK"), false) {
					w.Close()
				}
			})
	}
}

var customTheme = style.ColorTable{
	ColorText:                  color.RGBA{0, 0, 0, 255},
	ColorWindow:                color.RGBA{255, 255, 255, 255},
	ColorHeader:                color.RGBA{242, 242, 242, 255},
	ColorHeaderFocused:         color.RGBA{0xc3, 0x9a, 0x9a, 255},
	ColorBorder:                color.RGBA{0, 0, 0, 255},
	ColorButton:                color.RGBA{185, 185, 185, 255},
	ColorButtonHover:           color.RGBA{215, 215, 215, 255},
	ColorButtonActive:          color.RGBA{200, 200, 200, 255},
	ColorToggle:                color.RGBA{225, 225, 225, 255},
	ColorToggleHover:           color.RGBA{200, 200, 200, 255},
	ColorToggleCursor:          color.RGBA{30, 30, 30, 255},
	ColorSelect:                color.RGBA{175, 175, 175, 255},
	ColorSelectActive:          color.RGBA{190, 190, 190, 255},
	ColorSlider:                color.RGBA{190, 190, 190, 255},
	ColorSliderCursor:          color.RGBA{215, 215, 215, 255},
	ColorSliderCursorHover:     color.RGBA{235, 235, 235, 255},
	ColorSliderCursorActive:    color.RGBA{225, 225, 225, 255},
	ColorProperty:              color.RGBA{225, 225, 225, 255},
	ColorEdit:                  color.RGBA{245, 245, 245, 255},
	ColorEditCursor:            color.RGBA{0, 0, 0, 255},
	ColorCombo:                 color.RGBA{225, 225, 225, 255},
	ColorChart:                 color.RGBA{160, 160, 160, 255},
	ColorChartColor:            color.RGBA{45, 45, 45, 255},
	ColorChartColorHighlight:   color.RGBA{255, 0, 0, 255},
	ColorScrollbar:             color.RGBA{180, 180, 180, 255},
	ColorScrollbarCursor:       color.RGBA{140, 140, 140, 255},
	ColorScrollbarCursorHover:  color.RGBA{150, 150, 150, 255},
	ColorScrollbarCursorActive: color.RGBA{160, 160, 160, 255},
	ColorTabHeader:             color.RGBA{210, 210, 210, 255},
}
