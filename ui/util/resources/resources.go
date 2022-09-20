package resources

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods/managed/cache"
	"path/filepath"
)

const (
	mmmRepoResources = "https://raw.githubusercontent.com/kiamev/moogle-mod-manager/master/resources/"
	resourcesDir     = "resources"
)

var (
	LogoI   fyne.CanvasObject
	LogoII  fyne.CanvasObject
	LogoIII fyne.CanvasObject
	LogoIV  fyne.CanvasObject
	LogoV   fyne.CanvasObject
	LogoVI  fyne.CanvasObject

	LogoChronoCross fyne.CanvasObject

	LogoBofIII fyne.CanvasObject
	LogoBofIV  fyne.CanvasObject
)

func Initialize() {
	LogoI = loadLogo(config.I, "1.png")
	LogoII = loadLogo(config.II, "2.png")
	LogoIII = loadLogo(config.III, "3.png")
	LogoIV = loadLogo(config.IV, "4.png")
	LogoV = loadLogo(config.V, "5.png")
	LogoVI = loadLogo(config.VI, "6.png")
	LogoChronoCross = loadLogo(config.ChronoCross, "chronocross.png")
	LogoBofIII = loadLogo(config.BofIII, "bof3.png")
	LogoBofIV = loadLogo(config.BofIV, "bof4.png")
}

func loadLogo(game config.Game, f string) fyne.CanvasObject {
	var (
		file   = filepath.Join(config.PWD, resourcesDir)
		r, err = cache.GetImage(mmmRepoResources+f, file)
		img    *canvas.Image
	)

	if err != nil {
		return createTextLogo(game)
	}

	img = canvas.NewImageFromResource(r)
	size := fyne.Size{Width: 444 * .75, Height: 176 * .75}
	img.SetMinSize(size)
	img.Resize(size)
	img.FillMode = canvas.ImageFillContain
	return img
}

func createTextLogo(game config.Game) fyne.CanvasObject {
	return widget.NewLabel(config.GameNameString(game))
}
