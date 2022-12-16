package mover

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/action"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/conflict"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/managed"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type basicFileMover struct{}

func (m *basicFileMover) AddModFiles(enabler *mods.ModEnabler, mmf *managed.ModsAndFiles, files []*mods.DownloadFiles, cr conflict.Result) (err error) {
	var (
		game     = enabler.Game
		tm       = enabler.TrackedMod
		configs  = config.Get()
		modPath  = filepath.Join(configs.GetModsFullPath(game), tm.ID().AsDir())
		backedUp []*mods.ModFile
		moved    []*mods.ModFile
	)

	for _, df := range files {
		var (
			modDir     = filepath.Join(modPath, df.DownloadName)
			installDir string
		)
		if installDir, err = configs.GetDir(game, config.GameDirKind); err != nil {
			break
		}
		if err = m.MoveFiles(enabler.Game, df.Files, modDir, installDir, configs.GetBackupFullPath(game), &backedUp, &moved, cr, false); err != nil {
			break
		}
		if err == nil {
			if err = m.MoveDirs(game, df.Dirs, modDir, installDir, configs.GetBackupFullPath(game), &backedUp, &moved, cr, false); err != nil {
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
			if e := m.MoveFile(action.Cut, f.To, f.From, nil); e != nil {
				sb.WriteString(fmt.Sprintf("failed to restore [%s] from [%s]\n", f.To, f.From))
			}
		}
		return fmt.Errorf("%s: %v", sb.String(), err)
	}

	mf, found := mmf.Mods[tm.ID()]
	if !found {
		mf = &managed.ModFiles{
			BackedUpFiles: make(map[string]*mods.ModFile),
			MovedFiles:    make(map[string]*mods.ModFile),
		}
		mmf.Mods[tm.ID()] = mf
	}

	for _, f := range backedUp {
		mf.BackedUpFiles[f.From] = f
	}
	for _, f := range moved {
		mf.MovedFiles[f.To] = f
	}
	m.removeBackupFile(enabler, mf, mmf, cr.Replace)
	mmf.Mods[tm.ID()] = mf
	return managed.SaveManagedJson()
}

func (m *basicFileMover) removeBackupFile(enabler *mods.ModEnabler, mf *managed.ModFiles, mmf *managed.ModsAndFiles, toRemove map[string]bool) {
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
			if id != enabler.TrackedMod.ID() {
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

func (m *basicFileMover) MoveFiles(game config.GameDef, files []*mods.ModFile, modDir string, toDir string, backupDir string, backedUp *[]*mods.ModFile, movedFiles *[]*mods.ModFile, cr conflict.Result, returnOnFail bool) (err error) {
	var (
		c   = config.Get()
		dir string
	)
	for _, f := range files {
		to := path.Join(toDir, f.To)
		if m.IsDir(to) {
			to = filepath.Join(to, filepath.Base(f.From))
		}
		if dir, err = c.RemoveGameDir(game, to); err != nil {
			return
		}
		if cr.Skip[dir] {
			continue
		}
		if !cr.Replace[dir] {
			if _, err = os.Stat(to); err == nil {
				if err = m.MoveFile(action.Cut, to, path.Join(backupDir, f.To), backedUp); err != nil {
					if returnOnFail {
						return
					}
				}
			}
		}
		if err = m.MoveFile(action.Duplicate, path.Join(modDir, f.From), path.Join(toDir, f.To), movedFiles); err != nil {
			if returnOnFail {
				return
			}
		}
	}
	return
}

func (m *basicFileMover) MoveDirs(game config.GameDef, dirs []*mods.ModDir, modDir string, toDir string, backupDir string, replacedFiles *[]*mods.ModFile, movedFiles *[]*mods.ModFile, cr conflict.Result, returnOnFail bool) (err error) {
	var (
		mf   []*mods.ModFile
		from string
		to   string
	)
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
				c := strings.Count(to, string(game.BaseDir())+"/")
				if c == 0 && strings.HasPrefix(to, config.StreamingAssetsDir) {
					to = filepath.Join(string(game.BaseDir()), to)
				} else if c > 1 {
					to = strings.Replace(to, string(game.BaseDir())+"/", "", 1)
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
	return m.MoveFiles(game, mf, modDir, toDir, backupDir, replacedFiles, movedFiles, cr, returnOnFail)
}

func (m *basicFileMover) MoveFile(a action.FileAction, from, to string, files *[]*mods.ModFile) (err error) {
	if m.IsDir(to) {
		to = filepath.Join(to, filepath.Base(from))
	}
	if err = os.MkdirAll(filepath.Dir(to), 0777); err != nil {
		return
	}
	if a == action.Duplicate {
		err = m.copyFile(from, to)
	} else {
		err = m.cutFile(from, to)
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

func (m *basicFileMover) IsDir(path string) bool {
	return filepath.Ext(path) == ""
}

func (m *basicFileMover) cutFile(src, dst string) error {
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

func (m *basicFileMover) copyFile(src, dst string) error {
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

func (m *basicFileMover) RemoveModFiles(mf *managed.ModFiles, mmf *managed.ModsAndFiles, tm mods.TrackedMod) (err error) {
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
			if err = m.MoveFile(action.Cut, f.To, f.From, nil); err != nil {
				sb.WriteString(fmt.Sprintf("failed to move [%s] to [%s]: %v\n", f.To, f.From, err))
				err = nil
			}
		}
		handled = append(handled, k)
	}
	for _, h := range handled {
		delete(mf.BackedUpFiles, h)
	}

	_ = managed.SaveManagedJson()

	if err != nil {
		return
	}

	delete(mmf.Mods, tm.ID())
	return
}
