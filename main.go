package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/Xuanwo/go-locale"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	"github.com/kiamev/moogle-mod-manager/mods/managed/authored"
	config_installer "github.com/kiamev/moogle-mod-manager/ui/config-installer"
	"github.com/kiamev/moogle-mod-manager/ui/game-select"
	"github.com/kiamev/moogle-mod-manager/ui/local"
	"github.com/kiamev/moogle-mod-manager/ui/menu"
	mod_author "github.com/kiamev/moogle-mod-manager/ui/mod-author"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

func main() {
	state.App = app.New()
	state.Window = state.App.NewWindow("Moogle Mod Manager " + browser.Version)
	state.Window.Resize(fyne.NewSize(1000, 850))
	if err := managed.Initialize(); err != nil {
		dialog.ShowError(err, state.Window)
	}
	if err := authored.Initialize(); err != nil {
		dialog.ShowError(err, state.Window)
	}

	if tag, err := locale.Detect(); err != nil {
		dialog.ShowError(err, state.Window)
	} else {
		// TODO
		println(tag.String())
		//https://github.com/nicksnyder/go-i18n/
		//https://en.wikipedia.org/wiki/IETF_language_tag
		//en-US
	}

	state.RegisterMainMenu(menu.New())
	state.RegisterScreen(state.None, game_select.New())
	state.RegisterScreen(state.ModAuthor, mod_author.New())
	state.RegisterScreen(state.LocalMods, local.New())
	state.RegisterScreen(state.ConfigInstaller, config_installer.New())

	state.ShowScreen(state.None)
	state.Window.ShowAndRun()
}
