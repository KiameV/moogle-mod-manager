package downloads

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	"github.com/kiamev/moogle-mod-manager/ui/confirm"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
)

func Download(game config.Game, tm *model.TrackedMod, tis []*mods.ToInstall, done confirm.DownloadCompleteCallback) error {
	downloadDir, err := createPath(filepath.Join(config.Get().GetDownloadFullPath(game), tm.GetDirSuffix()))
	if err != nil {
		return err
	}

	if tm.Mod.ModKind.Kind == mods.Hosted {
		confirm.Hosted(game, downloadDir, tm, tis, done, hosted)
	} else {
		for _, ti := range tis {
			ti.Download.DownloadedLoc = filepath.Join(downloadDir, util.CreateFileName(ti.Download.Version))
		}
		if err = confirm.Nexus(game, downloadDir, tm, tis, done, nexus); err != nil {
			return err
		}
	}
	return nil
}

func hosted(game config.Game, downloadDir string, tm *model.TrackedMod, tis []*mods.ToInstall) (err error) {
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
			if f, err = browser.Download(source, filepath.Join(downloadDir, util.CreateFileName(ti.Download.Version))); err == nil {
				// success
				installed = append(installed, model.NewInstalledDownload(ti.Download.Name, ti.Download.Version))
				ti.Download.DownloadedLoc = f
				break
			}
		}
	}

	for _, ti := range tis {
		if ti.Download.DownloadedLoc == "" {
			fmt.Errorf("failed to download %s", ti.Download.Hosted.Sources[0])
			return
		}
	}
	return
}

func nexus(game config.Game, downloadDir string, tm *model.TrackedMod, tis []*mods.ToInstall) (err error) {
	var dir []os.DirEntry
	if dir, err = os.ReadDir(downloadDir); err != nil {
		return err
	}
	for _, ti := range tis {
		ti.Download.DownloadedLoc = ""
		for _, f := range dir {
			if ti.Download.Nexus.FileName == f.Name() {
				ti.Download.DownloadedLoc = filepath.Join(downloadDir, f.Name())
				break
			}
		}
		if ti.Download.DownloadedLoc == "" {
			err = fmt.Errorf("failed to find %s in %s", ti.Download.Nexus.FileName, downloadDir)
			return err
		}
	}
	return err
}

func createPath(path string) (string, error) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		err = fmt.Errorf("failed to create mod directory: %v", err)
		return "", err
	}
	return path, nil
}
