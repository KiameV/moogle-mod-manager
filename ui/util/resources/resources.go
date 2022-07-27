package resources

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/decompressor"
	"os"
	"path/filepath"
)

const resourcesDir = "resources"

var (
	LogoI   fyne.CanvasObject
	LogoII  fyne.CanvasObject
	LogoIII fyne.CanvasObject
	LogoIV  fyne.CanvasObject
	LogoV   fyne.CanvasObject
	LogoVI  fyne.CanvasObject
)

func Initialize() {
	resources := filepath.Join(config.PWD, resourcesDir)
	downloadResources(resources)
	LogoI = loadLogo(config.I, filepath.Join(resources, "1.png"))
	LogoII = loadLogo(config.II, filepath.Join(resources, "2.png"))
	LogoIII = loadLogo(config.III, filepath.Join(resources, "3.png"))
	LogoIV = loadLogo(config.IV, filepath.Join(resources, "4.png"))
	LogoV = loadLogo(config.V, filepath.Join(resources, "5.png"))
	LogoVI = loadLogo(config.VI, filepath.Join(resources, "6.png"))
}

func downloadResources(resources string) {
	var (
		f   string
		d   decompressor.Decompressor
		err error
	)
	if _, err = os.Stat(resources); err != nil {
		if err = os.Mkdir(resources, 0777); err != nil {
			return
		}
	}
	if _, err = os.Stat(filepath.Join(resources, "1.png")); err != nil {
		if f, err = browser.Download("https://github.com/KiameV/moogle-mod-manager/blob/main/resources/logos.zip?raw=true", "./resources"); err != nil {
			return
		}
		defer func() {
			_ = os.Remove(f)
		}()
		if d, err = decompressor.NewDecompressor(f); err != nil {
			return
		}
		if err = d.DecompressTo(resources); err != nil {
			return
		}
	}
}

func loadLogo(game config.Game, f string) fyne.CanvasObject {
	var (
		r   fyne.Resource
		err error
	)
	if _, err = os.Stat(f); err != nil {
		return createTextLogo(game)
	}
	if r, err = fyne.LoadResourceFromPath(f); err != nil {
		return createTextLogo(game)
	}
	img := canvas.NewImageFromResource(r)
	//size := fyne.Size{Width: float32(444), Height: float32(176)}
	size := fyne.Size{Width: 444 * .75, Height: 176 * .75}
	img.SetMinSize(size)
	img.Resize(size)
	img.FillMode = canvas.ImageFillContain
	return img
}

func createTextLogo(game config.Game) fyne.CanvasObject {
	return widget.NewLabel(config.GameNameString(game))
}
