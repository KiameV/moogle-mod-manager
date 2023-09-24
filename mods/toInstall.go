package mods

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
	"path/filepath"
)

type ToInstall struct {
	kinds         Kinds
	Download      *Download
	DownloadFiles []*DownloadFiles
	downloadDir   string
}

func NewToInstall(kinds Kinds, download *Download, downloadFiles *DownloadFiles) *ToInstall {
	return &ToInstall{
		kinds:         kinds,
		Download:      download,
		DownloadFiles: []*DownloadFiles{downloadFiles},
	}
}

func NewToInstallForMod(mod *Mod, downloadFiles []*DownloadFiles) (result []*ToInstall, err error) {
	mLookup := make(map[string]*Download)
	dfLookup := make(map[string]*DownloadFiles)
	for _, dl := range mod.Downloadables {
		mLookup[dl.Name] = dl
	}
	for _, df := range downloadFiles {
		i, ok := dfLookup[df.DownloadName]
		if !ok {
			i = &DownloadFiles{
				DownloadName: df.DownloadName,
			}
			dfLookup[df.DownloadName] = i
		}
		i.Files = append(i.Files, df.Files...)
		i.Dirs = append(i.Dirs, df.Dirs...)
	}
	for n, df := range dfLookup {
		dl := mLookup[n]
		result = append(result, NewToInstall(mod.Kinds(), dl, df))
	}
	return
}

func (ti *ToInstall) GetDownloadLocation(game config.GameDef, tm TrackedMod) (string, error) {
	if ti.downloadDir != "" {
		return ti.downloadDir, nil
	}
	if ti.kinds.IsHosted() {
		return ti.getHostedDownloadLocation(game, tm, tm.Mod().Version)
	}
	return ti.getRemoteDownloadLocation(game, tm)
}

func (ti *ToInstall) getHostedDownloadLocation(game config.GameDef, tm TrackedMod, v string) (string, error) {
	var m = tm.Mod()
	if v == "" {
		v = "nv"
	}
	if len(m.Games) > 0 && m.Category == config.Utility {
		ti.downloadDir = config.Get().GetDownloadFullPathForUtility()
	} else {
		ti.downloadDir = config.Get().GetDownloadFullPathForGame(game)
	}
	ti.downloadDir = filepath.Join(ti.downloadDir, tm.ID().AsDir(), util.CreateFileName(v))
	if err := createPath(ti.downloadDir); err != nil {
		return "", err
	}
	return ti.downloadDir, nil
}

func (ti *ToInstall) getRemoteDownloadLocation(game config.GameDef, tm TrackedMod) (string, error) {
	ti.downloadDir = filepath.Join(config.Get().GetDownloadFullPathForGame(game), tm.ID().AsDir(), util.CreateFileName(ti.Download.Version))
	if err := createPath(ti.downloadDir); err != nil {
		return "", err
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
