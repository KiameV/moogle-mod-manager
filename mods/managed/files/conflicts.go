package files

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type DoneCallback func(result mods.Result, skip map[string]bool, err ...error)

func ResolveConflicts(enabler *mods.ModEnabler, managedFiles map[mods.ModID]*modFiles, modFiles []*mods.DownloadFiles, done DoneCallback) {
	c := config.Get()
	fileToMod := make(map[string]mods.ModID)
	for modID, mf := range managedFiles {
		for _, f := range mf.MovedFiles {
			fileToMod[c.RemoveGameDir(enabler.Game, f.To)] = modID
		}
	}
	toInstall, err := compileFilesToMove(modFiles)
	if err != nil {
		done(mods.Error, nil, err)
	}

	detectCollisions(enabler, toInstall, fileToMod, done)
	return
}

func compileFilesToMove(modFiles []*mods.DownloadFiles) (toInstall []string, err error) {
	for _, mf := range modFiles {
		for _, f := range mf.Files {
			to := f.To
			if filepath.Ext(to) == "" {
				to = filepath.Join(to, filepath.Base(f.From))
			}
			toInstall = append(toInstall, strings.ReplaceAll(to, "\\", "/"))
		}
		for _, d := range mf.Dirs {
			if d.Recursive {
				_ = filepath.WalkDir(d.From, func(path string, d fs.DirEntry, err error) error {
					if d.IsDir() {
						return nil
					}
					toInstall = append(toInstall, strings.ReplaceAll(path, "\\", "/"))
					return nil
				})
			} else {
				var de []fs.DirEntry
				if de, err = os.ReadDir(d.From); err != nil {
					return
				}
				for _, e := range de {
					if e.IsDir() {
						continue
					}
					toInstall = append(toInstall, strings.ReplaceAll(e.Name(), "\\", "/"))
				}
			}
		}
	}
	return
}

func detectCollisions(enabler *mods.ModEnabler, toInstall []string, installedFiles map[string]mods.ModID, done DoneCallback) {
	var (
		newModID   = enabler.TrackedMod.GetModID()
		collisions []*mods.FileConflict
		id         mods.ModID
		found      bool
		skip       = make(map[string]bool)
	)
	for _, ti := range toInstall {
		if id, found = installedFiles[ti]; found {
			collisions = append(collisions, &mods.FileConflict{
				File:         ti,
				CurrentModID: id,
				NewModID:     newModID,
			})
		}
	}
	if len(collisions) > 0 {
		enabler.OnConflict(collisions, func(result mods.Result, choices []*mods.FileConflict, err ...error) {
			if result == mods.Error {
				done(result, nil, err...)
				return
			}
			if result == mods.Cancel {
				done(result, nil)
				return
			}
			for _, c := range choices {
				if c.ChoiceName != enabler.TrackedMod.DisplayName {
					skip[c.File] = true
				}
			}
			done(mods.Ok, skip)
		})
	} else {
		done(mods.Ok, nil)
	}
}
