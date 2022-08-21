package mods

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
)

type ToInstall struct {
	kind          Kind
	Download      *Download
	DownloadFiles []*DownloadFiles
	downloadDir   string
}

func NewToInstall(kind Kind, download *Download, downloadFiles *DownloadFiles) *ToInstall {
	return &ToInstall{
		kind:          kind,
		Download:      download,
		DownloadFiles: []*DownloadFiles{downloadFiles},
	}
}

func NewToInstallForMod(kind Kind, mod *Mod, downloadFiles []*DownloadFiles) (result []*ToInstall, err error) {
	lookup := make(map[string]*Download)
	for _, dl := range mod.Downloadables {
		lookup[dl.Name] = dl
	}
	for _, f := range downloadFiles {
		dl, _ := lookup[f.DownloadName]
		result = append(result, NewToInstall(kind, dl, f))
	}
	return
}

func (ti *ToInstall) GetDownloadLocation(game config.Game, tm *TrackedMod) (string, error) {
	if ti.kind == Hosted {
		return ti.getHostedDownloadLocation(game, tm)
	}
	return ti.getNexusDownloadLocation(game, tm)
}

func (ti *ToInstall) getHostedDownloadLocation(game config.Game, tm *TrackedMod) (string, error) {
	if ti.downloadDir == "" {
		v := ti.Download.Version
		if v == "" {
			v = "nv"
		}
		ti.downloadDir = filepath.Join(config.Get().GetDownloadFullPath(game), tm.GetDirSuffix(), util.CreateFileName(v))
		if err := createPath(ti.downloadDir); err != nil {
			return "", err
		}
	}
	return ti.downloadDir, nil
}

func (ti *ToInstall) getNexusDownloadLocation(game config.Game, tm *TrackedMod) (string, error) {
	if ti.downloadDir == "" {
		ti.downloadDir = filepath.Join(config.Get().GetDownloadFullPath(game), tm.GetDirSuffix(), util.CreateFileName(ti.Download.Version))
		if err := createPath(ti.downloadDir); err != nil {
			return "", err
		}
	}
	return ti.downloadDir, nil
}

func createPath(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		err = fmt.Errorf("failed to create mod directory: %v", err)
		return err
	}
	return nil
}
