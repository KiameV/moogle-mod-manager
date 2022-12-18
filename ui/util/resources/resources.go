package resources

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/cache"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/util"
	"path/filepath"
)

const (
	mmmRepoResources = "https://raw.githubusercontent.com/kiamev/moogle-mod-manager/master/resources/"
	resourcesDir     = "resources"
)

var (
	Icon fyne.Resource
)

func Initialize(games []config.GameDef) {
	for _, g := range games {
		if g.LogoPath() != "" {
			g.SetLogo(loadLogo(g))
		}
	}
	Icon, _ = loadImage("icon16.png")
}

func loadLogo(game config.GameDef) fyne.CanvasObject {
	var (
		r, err = loadImage(game.LogoPath())
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

func loadImage(f string) (fyne.Resource, error) {
	if util.FileExists(f) {
		return fyne.LoadResourceFromPath(f)
	}
	return cache.GetImage(mmmRepoResources+f, filepath.Join(config.PWD, resourcesDir))
}

func createTextLogo(game config.GameDef) fyne.CanvasObject {
	return widget.NewLabel(string(game.Name()))
}
