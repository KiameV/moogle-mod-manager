package managed

import (
	"encoding/json"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	managedXmlName = "managed.json"
)

var (
	managed = make(map[config.Game]*managedModsAndFiles)
)

type managedModsAndFiles struct {
	Mods map[string]*modFiles
}

type modFiles struct {
	BackedUpFiles map[string]*mods.ModFile
	MovedFiles    map[string]*mods.ModFile
}

func AddModFiles(game config.Game, tm *model.TrackedMod, files []*mods.DownloadFiles) (err error) {
	var (
		mmf, ok = managed[game]
		configs = config.Get()
		modPath = filepath.Join(configs.GetModsFullPath(game), tm.GetDirSuffix())
	)
	if !ok {
		mmf = &managedModsAndFiles{
			Mods: make(map[string]*modFiles),
		}
		managed[game] = mmf
	}

	if collisions := detectCollisions(nil, files); len(collisions) > 0 {
		return fmt.Errorf("cannot enable mod as these files would collide: %s", strings.Join(collisions, ", "))
	}

	var backedUp []*mods.ModFile
	var moved []*mods.ModFile

	for _, df := range files {
		path := filepath.Join(modPath, df.DownloadName)
		if err = MoveFiles(df.Files, path, configs.DirVI, configs.GetBackupFullPath(game), &backedUp, &moved); err != nil {
			return err
		}
		if err = MoveDirs(df.Dirs, path, configs.DirVI, configs.GetBackupFullPath(game), &backedUp, &moved); err != nil {
			return err
		}
	}

	mf, found := mmf.Mods[tm.GetModID()]
	if !found {
		mf = &modFiles{
			BackedUpFiles: make(map[string]*mods.ModFile),
			MovedFiles:    make(map[string]*mods.ModFile),
		}
		mmf.Mods[tm.GetModID()] = mf
	}

	for _, f := range backedUp {
		mf.BackedUpFiles[f.From] = f
	}
	for _, f := range moved {
		mf.MovedFiles[f.To] = f
	}
	mmf.Mods[tm.GetModID()] = mf

	return saveManagedJson()
}

func RemoveModFiles(game config.Game, tm *model.TrackedMod) (err error) {
	var (
		mmf, ok = managed[game]
		mf      *modFiles
	)
	if !ok {
		return fmt.Errorf("%s is not enabled", tm.Mod.Name)
	}
	if mf, ok = mmf.Mods[tm.GetModID()]; !ok {
		return fmt.Errorf("%s is not enabled", tm.Mod.Name)
	}

	handled := make([]string, 0, len(mf.MovedFiles))
	for k, f := range mf.MovedFiles {
		if _, err = os.Stat(f.To); err == nil {
			if err = os.Remove(f.To); err != nil {
				break
			}
		}
		handled = append(handled, k)
	}
	for _, h := range handled {
		delete(mf.MovedFiles, h)
	}

	handled = make([]string, 0, len(mf.BackedUpFiles))
	for k, f := range mf.BackedUpFiles {
		if _, err = os.Stat(f.From); err == nil {
			if err = os.Remove(f.From); err != nil {
				return
			}
		}
		if err = moveFile(cut, f.To, f.From, nil); err != nil {
			break
		}
		handled = append(handled, k)
	}
	for _, h := range handled {
		delete(mf.BackedUpFiles, h)
	}

	_ = saveManagedJson()

	if err != nil {
		return err
	}

	delete(mmf.Mods, tm.GetModID())
	return
}

func detectCollisions(managedFiles map[string]bool, modFiles []*mods.DownloadFiles) (collisions []string) {
	// TODO
	return
}

func saveManagedJson() error {
	b, err := json.MarshalIndent(managed, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(managedXmlName, b, 0777)
}

type action bool

const (
	duplicate action = false
	cut       action = true
)

func MoveFiles(files []*mods.ModFile, modDir string, toDir string, backupDir string, backedUp *[]*mods.ModFile, movedFiles *[]*mods.ModFile) (err error) {
	for _, f := range files {
		to := path.Join(toDir, f.To)
		if _, err = os.Stat(to); err == nil {
			if err = moveFile(cut, to, path.Join(backupDir, f.To), backedUp); err != nil {
				return
			}
		}
		if err = moveFile(duplicate, path.Join(modDir, f.From), path.Join(toDir, f.To), movedFiles); err != nil {
			return
		}
	}
	return
}

func MoveDirs(dirs []*mods.ModDir, modDir string, toDir string, backupDir string, replacedFiles *[]*mods.ModFile, movedFiles *[]*mods.ModFile) (err error) {
	/*var (
		toDir     = config.GetModDir(game)
		backupDir = config.GetBackupDir(game)
	)
	for _, f := range files {
		if err = os.Rename(path.Join(toDir, f.To), path.Join(backupDir, f.To)); err != nil {
			break
		}
		if err = copy(path.Join(modDir, f.From), path.Join(toDir, f.To)); err != nil {
			break
		}
	}*/
	return
}

func moveFile(action action, from, to string, backedUp *[]*mods.ModFile) (err error) {
	if err = os.MkdirAll(filepath.Dir(to), 0777); err != nil {
		return
	}
	if backedUp != nil {
		*backedUp = append(*backedUp, &mods.ModFile{
			From: from,
			To:   to,
		})
	}
	if action == duplicate {
		err = copyFile(from, to)
	} else {
		err = cutFile(from, to)
	}
	if err != nil {
		err = fmt.Errorf("failed to move [%s] to [%s]: %v", from, to, err)
	}
	return
}

func cutFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		_ = in.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	_ = in.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("sync error: %s", err)
	}

	si, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat error: %s", err)
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return fmt.Errorf("chmod error: %s", err)
	}

	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		_ = in.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	_ = in.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("sync error: %s", err)
	}

	si, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat error: %s", err)
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return fmt.Errorf("chmod error: %s", err)
	}

	return nil
}
