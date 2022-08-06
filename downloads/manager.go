package downloads

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	"github.com/kiamev/moogle-mod-manager/ui/confirm"
	"os"
	"path/filepath"
)

func Download(game config.Game, tm *model.TrackedMod, tis []*model.ToInstall, done confirm.DownloadCompleteCallback) error {
	if tm.Mod.ModKind.Kind == mods.Hosted {
		confirm.Hosted(game, tm, tis, done, hosted)
	} else {
		if err := confirm.Nexus(game, tm, tis, done, nexus); err != nil {
			return err
		}
	}
	return nil
}

func hosted(game config.Game, tm *model.TrackedMod, tis []*model.ToInstall) (err error) {
	var (
		f         string
		installed []*model.InstalledDownload
	)
	for _, ti := range tis {
		if len(ti.Download.Hosted.Sources) == 0 {
			err = fmt.Errorf("%s has no download sources", ti.Download.Name)
			return
		}
		for _, source := range ti.Download.Hosted.Sources {
			if f, err = ti.GetDownloadLocation(game, tm); err != nil {
				return
			}
			if f, err = browser.Download(source, f); err == nil {
				// success
				installed = append(installed, model.NewInstalledDownload(ti.Download.Name, ti.Download.Version))
				ti.Download.DownloadedArchiveLocation = f
				break
			}
		}
	}

	for _, ti := range tis {
		if ti.Download.DownloadedArchiveLocation == "" {
			fmt.Errorf("failed to download %s", ti.Download.Hosted.Sources[0])
			return
		}
	}
	return
}

func nexus(game config.Game, tm *model.TrackedMod, tis []*model.ToInstall) (err error) {
	var (
		dir  []os.DirEntry
		path string
	)
	for _, ti := range tis {
		if path, err = ti.GetDownloadLocation(game, tm); err != nil {
			return
		}
		if dir, err = os.ReadDir(path); err != nil {
			return
		}

		ti.Download.DownloadedArchiveLocation = ""
		for _, f := range dir {
			if ti.Download.Nexus.FileName == f.Name() {
				ti.Download.DownloadedArchiveLocation = filepath.Join(path, f.Name())
				break
			}
		}
		if ti.Download.DownloadedArchiveLocation == "" {
			err = fmt.Errorf("failed to find %s in %s", ti.Download.Nexus.FileName, path)
			return
		}
	}
	return
}
