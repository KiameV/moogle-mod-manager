package managed

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
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

func AddModFiles(game config.Game, tm *mods.TrackedMod, files []*mods.DownloadFiles) (err error) {
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

	if err = detectCollisions(mmf.Mods, files); err != nil {
		return
	}

	var backedUp []*mods.ModFile
	var moved []*mods.ModFile

	for _, df := range files {
		modDir := filepath.Join(modPath, df.DownloadName)
		if err = MoveFiles(df.Files, modDir, config.Get().GetGameDir(game), configs.GetBackupFullPath(game), &backedUp, &moved, false); err != nil {
			break
		}
		if err == nil {
			if err = MoveDirs(game, df.Dirs, modDir, config.Get().GetGameDir(game), configs.GetBackupFullPath(game), &backedUp, &moved, false); err != nil {
				break
			}
		}
	}

	if err != nil {
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%v\n", err))
		for _, f := range moved {
			if e := os.Remove(f.To); e != nil {
				sb.WriteString(fmt.Sprintf("failed to remove [%s]\n", f.To))
			}
		}
		for _, f := range backedUp {
			if e := MoveFile(cut, f.To, f.From, nil); e != nil {
				sb.WriteString(fmt.Sprintf("failed to restore [%s] from [%s]\n", f.To, f.From))
			}
		}
		return errors.New(fmt.Sprintf("%s: %v", sb.String(), err))
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

func RemoveModFiles(game config.Game, tm *mods.TrackedMod) (err error) {
	var (
		mmf, ok = managed[game]
		mf      *modFiles
		sb      = strings.Builder{}
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
				sb.WriteString(fmt.Sprintf("failed to remove [%s]: %v\n", f.To, err))
				err = nil
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
				sb.WriteString(fmt.Sprintf("failed to remove [%s]: %v\n", f.To, err))
				err = nil
			}
		}
		if err = MoveFile(cut, f.To, f.From, nil); err != nil {
			sb.WriteString(fmt.Sprintf("failed to move [%s] to [%s]: %v\n", f.To, f.From, err))
			err = nil
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

func detectCollisions(managedFiles map[string]*modFiles, modFiles []*mods.DownloadFiles) error {
	/*fileToMod := make(map[string]string)
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
	return nil
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

func initializeFiles() error {
	b, err := ioutil.ReadFile(filepath.Join(config.PWD, managedXmlName))
	if err != nil {
		return nil
	}
	return json.Unmarshal(b, &managed)
}

func saveManagedJson() error {
	b, err := json.MarshalIndent(managed, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(config.PWD, managedXmlName), b, 0777)
}

type action bool

const (
	duplicate action = false
	cut       action = true
)

func MoveFiles(files []*mods.ModFile, modDir string, toDir string, backupDir string, backedUp *[]*mods.ModFile, movedFiles *[]*mods.ModFile, returnOnFail bool) (err error) {
	for _, f := range files {
		to := path.Join(toDir, f.To)
		if _, err = os.Stat(to); err == nil {
			if err = MoveFile(cut, to, path.Join(backupDir, f.To), backedUp); err != nil {
				if returnOnFail {
					return
				}
			}
		}
		if err = MoveFile(duplicate, path.Join(modDir, f.From), path.Join(toDir, f.To), movedFiles); err != nil {
			if returnOnFail {
				return
			}
		}
	}
	return
}

func MoveDirs(game config.Game, dirs []*mods.ModDir, modDir string, toDir string, backupDir string, replacedFiles *[]*mods.ModFile, movedFiles *[]*mods.ModFile, returnOnFail bool) (err error) {
	var (
		mf   []*mods.ModFile
		from string
		to   string
	)
	toBaseDir := mods.GameToInstallBaseDir(game)
	modDir = strings.ReplaceAll(modDir, "\\", "/")
	for _, d := range dirs {
		fromDir := strings.ReplaceAll(d.From, "\\", "/")
		for len(fromDir) > 0 && (fromDir[0] == '.' || fromDir[0] == '/') {
			fromDir = fromDir[1:]
		}
		if err = filepath.Walk(filepath.Join(modDir, d.From),
			func(path string, info os.FileInfo, err error) error {
				if returnOnFail {
					err = nil
				}
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}

				from = strings.ReplaceAll(path, "\\", "/")
				from = strings.ReplaceAll(from, modDir, "")

				to = strings.ReplaceAll(from, modDir, "")
				to = filepath.Join(d.To, to)
				to = strings.ReplaceAll(to, "\\", "/")
				to = strings.TrimLeft(to, "/")
				c := strings.Count(to, string(toBaseDir)+"/")
				if c == 0 && strings.HasPrefix(to, mods.StreamingAssetsDir) {
					to = filepath.Join(string(toBaseDir), to)
				} else if c > 1 {
					to = strings.Replace(to, string(toBaseDir)+"/", "", 1)
				}

				mf = append(mf, &mods.ModFile{
					From: from,
					To:   to,
				})
				return nil
			}); err != nil {
			return
		}
	}
	return MoveFiles(mf, modDir, toDir, backupDir, replacedFiles, movedFiles, returnOnFail)
}

func MoveFile(action action, from, to string, files *[]*mods.ModFile) (err error) {
	if err = os.MkdirAll(filepath.Dir(to), 0777); err != nil {
		return
	}
	if action == duplicate {
		err = copyFile(from, to)
	} else {
		err = cutFile(from, to)
	}
	if err != nil {
		err = fmt.Errorf("failed to move [%s] to [%s]: %v", from, to, err)
		return
	}
	if files != nil {
		*files = append(*files, &mods.ModFile{
			From: from,
			To:   to,
		})
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
