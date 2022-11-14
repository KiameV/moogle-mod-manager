package downloads

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/confirm"
	"os"
	"path/filepath"
)

func Download(enabler *mods.ModEnabler, done confirm.DownloadCompleteCallback) error {
	var (
		kind      = enabler.Kind()
		confirmer = confirm.NewConfirmer(kind)
		callback  confirm.DownloadCallback
	)
	switch kind {
	case mods.CurseForge, mods.Nexus:
		callback = remote
	case mods.Hosted:
		callback = hosted
	default:
		return fmt.Errorf("unknown kind %v", kind)
	}

	return confirmer.ConfirmDownload(enabler, done, callback)
}

func hosted(enabler *mods.ModEnabler, done confirm.DownloadCompleteCallback, err error) {
	var installed []*mods.InstalledDownload

	for _, ti := range enabler.ToInstall {
		if len(ti.Download.Hosted.Sources) == 0 {
			err = fmt.Errorf("%s has no download sources", ti.Download.Name)
			done(enabler, err)
			return
		}
		for _, source := range ti.Download.Hosted.Sources {
			var f string
			if f, err = ti.GetDownloadLocation(enabler.Game, enabler.TrackedMod); err != nil {
				done(enabler, err)
				return
			}
			if f, err = browser.Download(source, f); err == nil {
				// success
				installed = append(installed, mods.NewInstalledDownload(ti.Download.Name, ti.Download.Version))
				ti.Download.DownloadedArchiveLocation = &f
				break
			}
		}
	}

	for _, ti := range enabler.ToInstall {
		if ti.Download.DownloadedArchiveLocation == nil || *ti.Download.DownloadedArchiveLocation == "" {
			done(enabler, fmt.Errorf("failed to download %s", ti.Download.Hosted.Sources[0]))
			return
		}
	}
	done(enabler, nil)
}

func remote(enabler *mods.ModEnabler, done confirm.DownloadCompleteCallback, err error) {
	var (
		dir  []os.DirEntry
		path string
	)
	for _, ti := range enabler.ToInstall {
		if path, err = ti.GetDownloadLocation(enabler.Game, enabler.TrackedMod); err != nil {
			done(enabler, err)
			return
		}
		if dir, err = os.ReadDir(path); err != nil {
			done(enabler, err)
			return
		}

		ti.Download.DownloadedArchiveLocation = nil
		for _, f := range dir {
			if ti.Download.Nexus.FileName == f.Name() {
				s := filepath.Join(path, f.Name())
				ti.Download.DownloadedArchiveLocation = &s
				break
			}
		}
		if ti.Download.DownloadedArchiveLocation == nil || *ti.Download.DownloadedArchiveLocation == "" {
			done(enabler, fmt.Errorf("failed to find %s in %s", ti.Download.Nexus.FileName, path))
			return
		}
	}
	done(enabler, nil)
}
