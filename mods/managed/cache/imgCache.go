package cache

import (
	"fyne.io/fyne/v2"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
)

func GetImage(url string, imgDirOverride ...string) (r fyne.Resource, err error) {
	var (
		key    = util.CreateFileName(url)
		imgDir = getImgDir(imgDirOverride...)
		fp     = filepath.Join(imgDir, key)
		_      = os.MkdirAll(fp, 0777)
		file   string
	)
	if file, err = browser.Download(url, fp); err != nil {
		return
	}
	return fyne.LoadResourceFromPath(file)
}

func getImgDir(imgDirOverride ...string) string {
	var imgDir = config.Get().ImgCacheDir
	if len(imgDirOverride) > 0 && imgDirOverride[0] != "" {
		imgDir = imgDirOverride[0]
	}
	return imgDir
}
