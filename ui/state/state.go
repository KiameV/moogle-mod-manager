package state

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/kiamev/moogle-mod-manager/config"
)

type GUI byte

const (
	None GUI = iota
	LocalMods
	ModAuthor
	ConfigInstaller
)

var (
	CurrentGame *config.Game
	App         fyne.App
	Window      fyne.Window

	guiHistories []*guiHistory
	mainMenu     Screen
	screens      = make(map[GUI]Screen)
	baseDir      = binding.NewString()
)

type guiHistory struct {
	gui     GUI
	baseDir string
}

func appendGuiHistory(gui GUI) {
	guiHistories = append(guiHistories, &guiHistory{
		gui:     gui,
		baseDir: GetBaseDir(),
	})
	SetBaseDir("")
}

type Screen interface {
	Draw(w fyne.Window)
	OnClose()
}

func GetCurrentGUI() GUI {
	if len(guiHistories) > 0 {
		return guiHistories[len(guiHistories)-1].gui
	}
	return None
}

func GetScreen(gui GUI) Screen {
	return screens[gui]
}

func ShowScreen(gui GUI) {
	appendGuiHistory(gui)
	mainMenu.Draw(Window)
	screens[gui].Draw(Window)
}

func ShowPreviousScreen() {
	var s Screen
	if len(guiHistories) > 1 {
		guiHistories = guiHistories[:len(guiHistories)-1]
		h := guiHistories[len(guiHistories)-1]
		s = screens[h.gui]
		SetBaseDir(h.baseDir)
	} else {
		s = screens[None]
		SetBaseDir("")
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

func RefreshMenu() {
	Window.MainMenu().Refresh()
}

func GetBaseDir() string {
	s, _ := baseDir.Get()
	return s
}

func GetBaseDirBinding() binding.String {
	return baseDir
}

func SetBaseDir(dir string) {
	_ = baseDir.Set(dir)
}

func SetCurrentGameFromString(s string) bool {
	var g config.Game
	switch s {
	case "I":
		g = config.I
	case "II":
		g = config.II
	case "III":
		g = config.III
	case "IV":
		g = config.IV
	case "V":
		g = config.V
	case "VI":
		g = config.VI
	default:
		return false
	}
	CurrentGame = &g
	return true
}
