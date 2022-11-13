package confirm

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type (
	DownloadCompleteCallback func(enabler *mods.ModEnabler, err error)
	DownloadCallback         func(enabler *mods.ModEnabler, completeCallback DownloadCompleteCallback, err error)
	Confirmer                interface {
		ConfirmDownload(enabler *mods.ModEnabler, competeCallback DownloadCompleteCallback, done DownloadCallback) (err error)
	}
)

func NewConfirmer(kind mods.Kind) Confirmer {
	switch kind {
	case mods.Nexus:
		return &nexusConfirmer{}
	case mods.CurseForge:
		return &cfConfirmer{}
	case mods.Hosted:
		return &hostedConfirmer{}
	}
	panic(fmt.Sprintf("unknown kind %v", kind))
}
