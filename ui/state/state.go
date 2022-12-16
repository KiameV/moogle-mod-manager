package state

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
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
	CurrentGame  config.GameDef
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
		if ui.PopupWindow == nil {
			ui.PopupWindow = ui.App.NewWindow("Finder")
			ui.PopupWindow.Resize(config.Get().Size())
			ui.PopupWindow.SetOnClosed(func() { ui.PopupWindow = nil })
			if err := screens[gui].PreDraw(ui.PopupWindow, args); err != nil {
				dialog.ShowError(err, ui.Window)
				return
			}
			ui.PopupWindow.Show()
			ui.ShowingPopup = true
			screens[gui].DrawAsDialog(ui.PopupWindow)
		}
		return
	} else {
		if err := screens[gui].PreDraw(ui.Window, args); err != nil {
			dialog.ShowError(err, ui.Window)
			return
		}
	}
	appendGuiHistory(gui)
	mainMenu.Draw(ui.Window)
	screens[gui].Draw(ui.Window)
}

func ClosePopupWindow() {
	ui.PopupWindow.Close()
	ui.ShowingPopup = false
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
	ui.Window.MainMenu().Refresh()
	mainMenu.Draw(ui.Window)
	s.Draw(ui.Window)
}

func UpdateCurrentScreen() {
	s := screens[GetCurrentGUI()]
	s.Draw(ui.Window)
}

func RegisterScreen(gui GUI, screen Screen) {
	screens[gui] = screen
}

func RegisterMainMenu(m Screen) {
	mainMenu = m
}

func RefreshMenu() {
	ui.Window.MainMenu().Refresh()
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
