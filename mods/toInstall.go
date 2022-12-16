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
		result = append(result, NewToInstall(kind, dl, df))
	}
	return
}

func (ti *ToInstall) GetDownloadLocation(game config.GameDef, tm TrackedMod) (string, error) {
	switch ti.kind {
	case Hosted:
		return ti.getHostedDownloadLocation(game, tm)
	case Nexus, CurseForge:
		return ti.getRemoteDownloadLocation(game, tm)
	}
	panic(fmt.Sprintf("unknown kind %v", ti.kind))
}

func (ti *ToInstall) getHostedDownloadLocation(game config.GameDef, tm TrackedMod) (string, error) {
	if ti.downloadDir == "" {
		var (
			v = ti.Download.Version
			m = tm.Mod()
		)
		if v == "" {
			v = "nv"
		}
		if len(m.Games) > 0 && m.Category == Utility {
			ti.downloadDir = config.Get().GetDownloadFullPathForUtility()
		} else {
			ti.downloadDir = config.Get().GetDownloadFullPathForGame(game)
		}
		ti.downloadDir = filepath.Join(ti.downloadDir, tm.ID().AsDir(), util.CreateFileName(v))
		if err := createPath(ti.downloadDir); err != nil {
			return "", err
		}
	}
	return ti.downloadDir, nil
}

func (ti *ToInstall) getRemoteDownloadLocation(game config.GameDef, tm TrackedMod) (string, error) {
	if ti.downloadDir == "" {
		ti.downloadDir = filepath.Join(config.Get().GetDownloadFullPathForGame(game), tm.ID().AsDir(), util.CreateFileName(ti.Download.Version))
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
