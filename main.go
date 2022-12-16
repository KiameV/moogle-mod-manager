package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/discover/repo"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	"github.com/kiamev/moogle-mod-manager/mods/managed/authored"
	config_installer "github.com/kiamev/moogle-mod-manager/ui/config-installer"
	"github.com/kiamev/moogle-mod-manager/ui/configure"
	"github.com/kiamev/moogle-mod-manager/ui/discover"
	"github.com/kiamev/moogle-mod-manager/ui/game-select"
	"github.com/kiamev/moogle-mod-manager/ui/local"
	"github.com/kiamev/moogle-mod-manager/ui/menu"
	mod_author "github.com/kiamev/moogle-mod-manager/ui/mod-author"
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

	// Mod versions
	// \Steam\steamapps\appmanifest_1173820.acf
	// FF1:1173770
	// FF2:1173780
	// FF3:1173790
	// FF4:1173800
	// FF5:1173810
	// FF6:1173820
	/*
			var tag language.Tag
			if tag, err = locale.Detect(); err != nil {
				util.ShowErrorLong(err)
			} else {
				// TODO
				println(tag.String())
				//https://github.com/nicksnyder/go-i18n/
				//https://en.wikipedia.org/wiki/IETF_language_tag
				//en-US
			}



		  setup: async (discovery, context) => {
		    await checkForRequiredToolStat({
		      context,
		      discovery,
		      name: 'Memoria FF6',
		      url: 'https://github.com/Albeoris/Memoria.FFPR',
		      paths: [path.join('BepInEx', 'plugins', 'Memoria.FF6.dll')],
		    })

		    await checkForRequiredToolStat({
		      context,
		      discovery,
		      name: 'BepInEx',
		      url: 'https://github.com/Albeoris/Memoria.FFPR',
		      paths: [path.join('BepInEx', 'core', 'BepInEx.Core.dll')],
		    })

		    await checkForRequiredToolSha256({
		      context,
		      discovery,
		      name: 'UnityPlayer.dll',
		      path: 'UnityPlayer.dll',
		      url: 'https://www.dropbox.com/s/pyqpoxpl7i4i67a/UnityPlayer.7z',
		      hashes: [
		        'F1B5D1110914CEBEF9D31A935239262342DEBDE78115D90F48C640CD39673CBE',
		      ],
		    })
		  },
	*/

	state.RegisterMainMenu(menu.New())
	state.RegisterScreen(state.None, game_select.New())
	state.RegisterScreen(state.ModAuthor, mod_author.New())
	state.RegisterScreen(state.LocalMods, local.New())
	state.RegisterScreen(state.DiscoverMods, discover.New())
	state.RegisterScreen(state.ConfigInstaller, config_installer.New())

	state.ShowScreen(state.None)
	if config.Get().FirstTime {
		configure.Show(ui.Window)
	}

	if game, err := config.GameDefFromID(config.GameID(config.Get().DefaultGame)); err == nil {
		state.CurrentGame = game
		state.ShowScreen(state.LocalMods)
	}

	ui.Window.ShowAndRun()
}

func initialize() {
	var err error
	config.GetSecrets().Initialize()

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
	}

	if err = config.Initialize(repo.Dirs(repo.Read)); err != nil {
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
