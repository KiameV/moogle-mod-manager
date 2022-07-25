package cache

import (
	"fyne.io/fyne/v2"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
)

func GetImage(url string) (r fyne.Resource, err error) {
	var (
		key    = util.CreateFileName(url)
		imgDir = config.Get().ImgCacheDir
		fp     = filepath.Join(imgDir, key)
		_      = os.MkdirAll(fp, 0777)
		file   string
	)
	if file, err = browser.Download(url, fp); err != nil {
		return
	}
	return fyne.LoadResourceFromPath(file)
}
