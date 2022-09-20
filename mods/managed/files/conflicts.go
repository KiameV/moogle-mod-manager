package files

import (
	"github.com/kiamev/moogle-mod-manager/mods"
)

func ResolveConflicts(enabler *mods.ModEnabler, managedFiles map[mods.ModID]*modFiles, modFiles []*mods.DownloadFiles, done mods.DoneCallback) {
	/*fileToMod := make(map[string]mods.ModID)
	for modID, mf := range managedFiles {
		for _, f := range mf.MovedFiles {
			fileToMod[f.To] = modID
		}
	}
	collisions := detectCollisionsFromMap(fileToMod, modFiles)
	if len(collisions) > 0 {
		sb := strings.Builder{}
		for modID, fs := range collisions {
			sb.WriteString(modID)
			sb.WriteByte('\n')
			for _, f := range fs {
				sb.WriteString(f)
				sb.WriteByte('\n')
			}
		}
		return fmt.Errorf("cannot enable mod as these files would collide:\n%s", sb.String())
	}*/
	done(mods.Ok)
	return
}

func detectCollisionsFromMap(rootDir string, fileToMod map[string]string, modFiles []*mods.DownloadFiles) (collisions map[string][]string) {
	/*for _, mf := range modFiles {
		for _, f := range mf.Files {
			if modID, found := fileToMod[f.To]; found {
				collisions[modID] = append(collisions[modID], f.To)
			}
		}
		for _, d := range mf.Dirs {
			_ = filepath.WalkDir(rootDir)
		}
	}*/
	return
}
