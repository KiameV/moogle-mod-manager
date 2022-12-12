package confirm

import (
	"github.com/kiamev/moogle-mod-manager/mods"
)

type cfConfirmer struct{}

func (_ *cfConfirmer) ConfirmDownload(enabler *mods.ModEnabler, competeCallback DownloadCompleteCallback, done DownloadCallback) (err error) {
	done(enabler, competeCallback, err)
	return
}
