package state

import (
	"fyne.io/fyne/v2"
	"github.com/kiamev/moogle-mod-manager/config"
)

type GUI byte

const (
	None GUI = iota
	LocalMods
	Configure
	ModAuthor
	ConfigInstaller
)

var (
	CurrentGame *config.Game
	App         fyne.App
	Window      fyne.Window

	guiHistory []GUI
	mainMenu   Screen
	screens    = make(map[GUI]Screen)
)

type Screen interface {
	Draw(w fyne.Window)
}

func GetCurrentGUI() GUI {
	if len(guiHistory) > 0 {
		return guiHistory[len(guiHistory)-1]
	}
	return None
}

func GetScreen(gui GUI) Screen {
	return screens[gui]
}

func ShowScreen(gui GUI) {
	guiHistory = append(guiHistory, gui)
	mainMenu.Draw(Window)
	screens[gui].Draw(Window)
}

func ShowPreviousScreen() {
	var s Screen
	if len(guiHistory) > 1 {
		guiHistory = guiHistory[:len(guiHistory)-1]
		s = screens[guiHistory[len(guiHistory)-1]]
	} else {
		s = screens[None]
	}
	Window.MainMenu().Refresh()
	mainMenu.Draw(Window)
	s.Draw(Window)
}

func RegisterScreen(gui GUI, screen Screen) {
	screens[gui] = screen
}

func RegisterMainMenu(m Screen) {
	mainMenu = m
}
