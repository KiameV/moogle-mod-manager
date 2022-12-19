package downloads

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"os"
	"path/filepath"
)

func Download(game config.GameDef, mod mods.TrackedMod, toInstall []*mods.ToInstall) (err error) {
	switch mod.Kind() {
	case mods.CurseForge:
		err = curseForge(game, mod, toInstall)
	case mods.Nexus:
		err = nexus(game, mod, toInstall)
	case mods.Hosted:
		err = hosted(game, mod, toInstall)
	default:
		return fmt.Errorf("unknown kind %v", mod.Kind())
	}
	return
}

func hosted(game config.GameDef, mod mods.TrackedMod, toInstall []*mods.ToInstall) error {
	var (
		f   string
		err error
	)
	for _, ti := range toInstall {
		if len(ti.Download.Hosted.Sources) == 0 {
			return fmt.Errorf("%s has no download sources", ti.Download.Name)
		}
		for _, source := range ti.Download.Hosted.Sources {
			if f, err = ti.GetDownloadLocation(game, mod); err != nil {
				return err
			}
			if f, err = browser.Download(source, f); err == nil {
				// success
				ti.Download.DownloadedArchiveLocation = (*mods.ArchiveLocation)(&f)
				break
			}
		}
		if ti.Download.DownloadedArchiveLocation == nil || *ti.Download.DownloadedArchiveLocation == "" {
			return fmt.Errorf("failed to download %s", ti.Download.Hosted.Sources[0])
		}
	}

	for _, ti := range toInstall {
		if ti.Download.DownloadedArchiveLocation == nil || *ti.Download.DownloadedArchiveLocation == "" {
			return fmt.Errorf("failed to download %s", ti.Download.Hosted.Sources[0])
		}
	}
	return nil
}

func nexus(game config.GameDef, mod mods.TrackedMod, toInstall []*mods.ToInstall) error {
	var (
		dir  []os.DirEntry
		path string
		name string
		err  error
	)
	for _, ti := range toInstall {
		if path, err = ti.GetDownloadLocation(game, mod); err != nil {
			return err
		}
		if dir, err = os.ReadDir(path); err != nil {
			return err
		}
		if name, err = ti.Download.FileName(); err != nil {
			return err
		}

		ti.Download.DownloadedArchiveLocation = nil
		for _, f := range dir {
			if name == f.Name() {
				s := filepath.Join(path, f.Name())
				ti.Download.DownloadedArchiveLocation = (*mods.ArchiveLocation)(&s)
				break
			}
		}
		if ti.Download.DownloadedArchiveLocation == nil || *ti.Download.DownloadedArchiveLocation == "" {
			return fmt.Errorf("failed to find %s in %s", name, path)
		}
	}
	return nil
}

func curseForge(game config.GameDef, mod mods.TrackedMod, toInstall []*mods.ToInstall) error {
	var (
		f   string
		err error
	)
	for _, ti := range toInstall {
		for _, i := range toInstall {
			if f, err = ti.GetDownloadLocation(game, mod); err != nil {
				return err
			}
			if f, err = browser.Download(i.Download.CurseForge.Url, f); err == nil {
				// success
				ti.Download.DownloadedArchiveLocation = (*mods.ArchiveLocation)(&f)
				break
			}
		}
		if ti.Download.DownloadedArchiveLocation == nil || *ti.Download.DownloadedArchiveLocation == "" {
			return fmt.Errorf("failed to download %s", ti.Download.Hosted.Sources[0])
		}
	}

	for _, ti := range toInstall {
		if ti.Download.DownloadedArchiveLocation == nil || *ti.Download.DownloadedArchiveLocation == "" {
			return fmt.Errorf("failed to download %s", ti.Download.Hosted.Sources[0])
		}
	}
	return nil
}
