package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/config/secrets"
	"github.com/kiamev/moogle-mod-manager/discover/repo"
	"github.com/kiamev/moogle-mod-manager/files"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	"github.com/kiamev/moogle-mod-manager/mods/managed/authored"
	config_installer "github.com/kiamev/moogle-mod-manager/ui/config-installer"
	"github.com/kiamev/moogle-mod-manager/ui/configure"
	"github.com/kiamev/moogle-mod-manager/ui/discover"
	"github.com/kiamev/moogle-mod-manager/ui/game-select"
	"github.com/kiamev/moogle-mod-manager/ui/local"
	"github.com/kiamev/moogle-mod-manager/ui/menu"
	mod_author "github.com/kiamev/moogle-mod-manager/ui/mod-author"
	"github.com/kiamev/moogle-mod-manager/ui/secret"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"github.com/kiamev/moogle-mod-manager/ui/util/resources"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			_ = os.WriteFile("log.txt", []byte(err.(string)), 0644)
		}
	}()

	if os.Getenv("profile") == "true" {
		f, err := os.Create(filepath.Join(config.PWD, "cpuprofile"))
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	ui.App = app.New()
	ui.Window = ui.App.NewWindow("Moogle Mod Manager " + browser.Version)
	initialize()

	state.RegisterMainMenu(menu.New())
	state.RegisterScreen(state.None, game_select.New())
	state.RegisterScreen(state.ModAuthor, mod_author.New())
	state.RegisterScreen(state.LocalMods, local.New())
	state.RegisterScreen(state.DiscoverMods, discover.New())
	state.RegisterScreen(state.ConfigInstaller, config_installer.New())

	state.ShowScreen(state.None)
	if config.Get().FirstTime {
		configure.Show(ui.Window, func() {
			secret.Show(ui.Window)
		})
	}

	if game, err := config.GameDefFromID(config.GameID(config.Get().DefaultGame)); err == nil {
		state.CurrentGame = game
		state.ShowScreen(state.LocalMods)
	}

	ui.Window.ShowAndRun()
}

func initialize() {
	var err error
	secrets.Initialize()

	if err = repo.Initialize(); err != nil {
		util.ShowErrorLong(err)
	}

	configs := config.Get()
	if err = configs.Initialize(); err != nil {
		util.ShowErrorLong(err)
	}

	ui.Window.Resize(config.Get().Size())
	ui.Window.SetMaster()

	if configs.Theme == config.LightThemeColor {
		fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	}

	if err = repo.NewGetter(repo.Read).Pull(); err != nil {
		util.ShowErrorLong(err)
		return
	}

	if err = config.Initialize(repo.Dirs(repo.Read)); err != nil {
		util.ShowErrorLong(err)
	}

	if err = files.Initialize(); err != nil {
		util.ShowErrorLong(err)
	}

	if err = managed.Initialize(config.GameDefs()); err != nil {
		util.ShowErrorLong(err)
	}
	if err = authored.Initialize(); err != nil {
		util.ShowErrorLong(err)
	}

	configs.InitializeGames(config.GameDefs())
	resources.Initialize(config.GameDefs())
	if resources.Icon != nil {
		ui.Window.SetIcon(resources.Icon)
	}
}
