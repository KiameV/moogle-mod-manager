package state

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"github.com/kiamev/moogle-mod-manager/config"
)

type GUI byte

const (
	None GUI = iota
	LocalMods
	DiscoverMods
	ModAuthor
	ConfigInstaller
)

var (
	CurrentGame *config.Game
	App         fyne.App
	Window      fyne.Window
	popupWindow fyne.Window

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
	PreDraw(w fyne.Window, args ...interface{}) error
	Draw(w fyne.Window)
	DrawAsDialog(window fyne.Window)
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

func ShowScreen(gui GUI, args ...interface{}) {
	if gui == DiscoverMods {
		if popupWindow == nil {
			popupWindow = App.NewWindow("Finder")
			popupWindow.Resize(config.Get().Size())
			popupWindow.SetOnClosed(func() { popupWindow = nil })
			popupWindow.Show()
			if err := screens[gui].PreDraw(popupWindow, args); err != nil {
				dialog.ShowError(err, Window)
				return
			}
			screens[gui].DrawAsDialog(popupWindow)
		}
		return
	} else {
		if err := screens[gui].PreDraw(Window, args); err != nil {
			dialog.ShowError(err, Window)
			return
		}
	}
	appendGuiHistory(gui)
	mainMenu.Draw(Window)
	screens[gui].Draw(Window)
}

func ClosePopupWindow() {
	popupWindow.Close()
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

func UpdateCurrentScreen() {
	s := screens[GetCurrentGUI()]
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
