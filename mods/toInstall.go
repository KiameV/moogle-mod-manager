package mods

import "fmt"

type ToInstall struct {
	Download      *Download
	DownloadFiles []*DownloadFiles
}

func NewToInstall(download *Download, downloadFiles *DownloadFiles) *ToInstall {
	return &ToInstall{
		Download:      download,
		DownloadFiles: []*DownloadFiles{downloadFiles},
	}
}

func NewToInstallForMod(mod *Mod, downloadFiles []*DownloadFiles) (result []*ToInstall, err error) {
	lookup := make(map[string]*Download)
	for _, dl := range mod.Downloadables {
		lookup[dl.Name] = dl
	}

	for _, f := range downloadFiles {
		dl, ok := lookup[f.DownloadName]
		if !ok {
			return nil, fmt.Errorf("could not find download %s for mod %s", f.DownloadName, mod.Name)
		}
		result = append(result, NewToInstall(dl, f))
	}
	return
}
