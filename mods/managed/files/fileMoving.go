package files

import (
	"errors"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func addModFiles(enabler *mods.ModEnabler, mmf *managedModsAndFiles, files []*mods.DownloadFiles, cr conflictResult) (err error) {
	var (
		game     = enabler.Game
		tm       = enabler.TrackedMod
		configs  = config.Get()
		modPath  = filepath.Join(configs.GetModsFullPath(game), tm.GetDirSuffix())
		backedUp []*mods.ModFile
		moved    []*mods.ModFile
	)

	for _, df := range files {
		modDir := filepath.Join(modPath, df.DownloadName)
		if err = MoveFiles(enabler.Game, df.Files, modDir, config.Get().GetGameDir(game), configs.GetBackupFullPath(game), &backedUp, &moved, cr, false); err != nil {
			break
		}
		if err == nil {
			if err = MoveDirs(game, df.Dirs, modDir, config.Get().GetGameDir(game), configs.GetBackupFullPath(game), &backedUp, &moved, cr, false); err != nil {
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
	removeBackupFile(enabler, mf, mmf, cr.replace)
	mmf.Mods[tm.GetModID()] = mf
	return saveManagedJson()
}

func removeBackupFile(enabler *mods.ModEnabler, mf *modFiles, mmf *managedModsAndFiles, toRemove map[string]bool) {
	var (
		c    = config.Get()
		game = enabler.Game
		k    string
		ok   bool
		dir  string
		err  error
	)
	if len(toRemove) > 0 {
		// If there are skipped files, remove from backup
		for id, mod := range mmf.Mods {
			if id != enabler.TrackedMod.GetModID() {
				var toDelete []string
				for k = range mod.BackedUpFiles {
					if dir, err = c.RemoveDir(game, config.GameDirKind, k); err != nil {
						continue
					}
					if _, ok = toRemove[dir]; ok {
						toDelete = append(toDelete, k)
					}
				}
				for k = range mod.MovedFiles {
					if dir, err = c.RemoveDir(game, config.GameDirKind, k); err != nil {
						continue
					}
					if _, ok = toRemove[dir]; ok {
						toDelete = append(toDelete, k)
					}
				}
				if len(toDelete) > 0 {
					for _, f := range toDelete {
						mf.BackedUpFiles[f] = mod.BackedUpFiles[f]
						delete(mod.BackedUpFiles, f)
						delete(mod.MovedFiles, f)
					}
				}
			}
		}
	}
}

func MoveFiles(game config.Game, files []*mods.ModFile, modDir string, toDir string, backupDir string, backedUp *[]*mods.ModFile, movedFiles *[]*mods.ModFile, cr conflictResult, returnOnFail bool) (err error) {
	var (
		c   = config.Get()
		dir string
	)
	for _, f := range files {
		to := path.Join(toDir, f.To)
		if IsDir(to) {
			to = filepath.Join(to, filepath.Base(f.From))
		}
		dir = c.RemoveGameDir(game, to)
		if cr.skip[dir] {
			continue
		}
		if !cr.replace[dir] {
			if _, err = os.Stat(to); err == nil {
				if err = MoveFile(cut, to, path.Join(backupDir, f.To), backedUp); err != nil {
					if returnOnFail {
						return
					}
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

func MoveDirs(game config.Game, dirs []*mods.ModDir, modDir string, toDir string, backupDir string, replacedFiles *[]*mods.ModFile, movedFiles *[]*mods.ModFile, cr conflictResult, returnOnFail bool) (err error) {
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
				to = strings.Replace(to, fromDir, "", 1)
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
	return MoveFiles(game, mf, modDir, toDir, backupDir, replacedFiles, movedFiles, cr, returnOnFail)
}

func MoveFile(action action, from, to string, files *[]*mods.ModFile) (err error) {
	if IsDir(to) {
		to = filepath.Join(to, filepath.Base(from))
	}
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

func IsDir(path string) bool {
	return filepath.Ext(path) == ""
}

func cutFile(src, dst string) error {
	var (
		in, out *os.File
		fi      os.FileInfo
		err     error
	)
	if in, err = os.Open(src); err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	if out, err = os.Create(dst); err != nil {
		_ = in.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, in)
	_ = in.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	if err = out.Sync(); err != nil {
		return fmt.Errorf("sync error: %s", err)
	}
	if fi, err = os.Stat(src); err != nil {
		return fmt.Errorf("stat error: %s", err)
	}
	if err = os.Chmod(dst, fi.Mode()); err != nil {
		return fmt.Errorf("chmod error: %s", err)
	}
	if err = os.Remove(src); err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	var (
		in, out *os.File
		fi      os.FileInfo
		err     error
	)
	if in, err = os.Open(src); err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	if out, err = os.Create(dst); err != nil {
		_ = in.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, in)
	_ = in.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	if err = out.Sync(); err != nil {
		return fmt.Errorf("sync error: %s", err)
	}
	if fi, err = os.Stat(src); err != nil {
		return fmt.Errorf("stat error: %s", err)
	}
	if err = os.Chmod(dst, fi.Mode()); err != nil {
		return fmt.Errorf("chmod error: %s", err)
	}
	return nil
}

func removeModFiles(mf *modFiles, mmf *managedModsAndFiles, tm *mods.TrackedMod) (err error) {
	var (
		handled = make([]string, 0, len(mf.MovedFiles))
		sb      = strings.Builder{}
	)
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
		if f != nil {
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
		}
		handled = append(handled, k)
	}
	for _, h := range handled {
		delete(mf.BackedUpFiles, h)
	}

	_ = saveManagedJson()

	if err != nil {
		return
	}

	delete(mmf.Mods, tm.GetModID())
	return
}
