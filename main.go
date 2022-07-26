package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	"github.com/kiamev/moogle-mod-manager/mods/managed/authored"
	config_installer "github.com/kiamev/moogle-mod-manager/ui/config-installer"
	"github.com/kiamev/moogle-mod-manager/ui/configure"
	"github.com/kiamev/moogle-mod-manager/ui/game-select"
	"github.com/kiamev/moogle-mod-manager/ui/local"
	"github.com/kiamev/moogle-mod-manager/ui/menu"
	mod_author "github.com/kiamev/moogle-mod-manager/ui/mod-author"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

func main() {
	state.App = app.New()
	state.Window = state.App.NewWindow("Moogle Mod Manager " + browser.Version)
	initialize()
	/*
		var tag language.Tag
		if tag, err = locale.Detect(); err != nil {
			dialog.ShowError(err, state.Window)
		} else {
			// TODO
			println(tag.String())
			//https://github.com/nicksnyder/go-i18n/
			//https://en.wikipedia.org/wiki/IETF_language_tag
			//en-US
		}
	*/

	state.RegisterMainMenu(menu.New())
	state.RegisterScreen(state.None, game_select.New())
	state.RegisterScreen(state.ModAuthor, mod_author.New())
	state.RegisterScreen(state.LocalMods, local.New())
	state.RegisterScreen(state.ConfigInstaller, config_installer.New())

	state.ShowScreen(state.None)
	if config.Get().FirstTime {
		configure.Show(state.Window)
	}

	state.Window.ShowAndRun()
}

func initialize() {
	var err error
	if err = managed.Initialize(); err != nil {
		dialog.ShowError(err, state.Window)
	}
	if err = authored.Initialize(); err != nil {
		dialog.ShowError(err, state.Window)
	}
	configs := config.Get()
	if err = configs.Initialize(); err != nil {
		dialog.ShowError(err, state.Window)
	}
	size := fyne.NewSize(config.WindowWidth, config.WindowHeight)
	if x := configs.WindowX; x != 0 {
		size.Width = float32(x)
	}
	if y := configs.WindowY; y != 0 {
		size.Height = float32(y)
	}
	state.Window.Resize(size)

	if configs.Theme == config.LightThemeColor {
		fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	}
}
