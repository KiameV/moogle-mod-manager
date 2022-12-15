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
	case mods.CurseForge:
		callback = curseForge
	case mods.Nexus:
		callback = remote
	case mods.Hosted:
		callback = hosted
	default:
		return fmt.Errorf("unknown kind %v", kind)
	}

	return confirmer.ConfirmDownload(enabler, done, callback)
}

func hosted(enabler *mods.ModEnabler, done confirm.DownloadCompleteCallback, err error) {
	if err != nil {
		return
	}
	for _, ti := range enabler.ToInstall {
		if len(ti.Download.Hosted.Sources) == 0 {
			done(enabler, fmt.Errorf("%s has no download sources", ti.Download.Name))
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
	if err != nil {
		return
	}
	var (
		dir  []os.DirEntry
		path string
		name string
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
		if name, err = ti.Download.FileName(); err != nil {
			done(enabler, err)
			return
		}

		ti.Download.DownloadedArchiveLocation = nil
		for _, f := range dir {
			if name == f.Name() {
				s := filepath.Join(path, f.Name())
				ti.Download.DownloadedArchiveLocation = &s
				break
			}
		}
		if ti.Download.DownloadedArchiveLocation == nil || *ti.Download.DownloadedArchiveLocation == "" {
			done(enabler, fmt.Errorf("failed to find %s in %s", name, path))
			return
		}
	}
	done(enabler, nil)
}

func curseForge(enabler *mods.ModEnabler, done confirm.DownloadCompleteCallback, err error) {
	if err != nil {
		return
	}
	for _, ti := range enabler.ToInstall {
		if len(enabler.ToInstall) == 0 {
			done(enabler, fmt.Errorf("%s has no download sources", ti.Download.Name))
			return
		}
		for _, i := range enabler.ToInstall {
			var f string
			if f, err = ti.GetDownloadLocation(enabler.Game, enabler.TrackedMod); err != nil {
				done(enabler, err)
				return
			}
			if f, err = browser.Download(i.Download.CurseForge.Url, f); err == nil {
				// success
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
