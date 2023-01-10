package steps

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/files"
	"github.com/kiamev/moogle-mod-manager/mods"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*type zipReader struct {
	Reader *zip.ReadCloser
	Files  map[string]*zip.File
}

func newZipReader(s string) (z *zipReader, err error) {
	z = &zipReader{Files: make(map[string]*zip.File)}
	if z.Reader, err = zip.OpenReader(s); err != nil {
		return
	}
	for _, f := range z.Reader.File {
		z.Files[f.Name] = f
	}
	return
}

func (z *zipReader) close() {
	if z.Reader != nil {
		_ = z.Reader.Close()
	}
}

func backupArchivedFiles(state *State, backupDir string) error {
	var (
		zr    map[string]*zipReader
		r     *zipReader
		found bool
		err   error
	)
	defer func() {
		for _, r = range zr {
			r.close()
		}
	}()
	for _, e := range state.ExtractedFiles {
		for _, ti := range e.FilesToInstall() {
			// Read all archives
			if ti.Skip {
				continue
			}
			if r, found = zr[*ti.archive]; !found {
				if r, err = newZipReader(*ti.archive); err != nil {
					return err
				}
				zr[*ti.archive] = r
			}
			if err = backupArchiveFile(r, backupDir, ti); err != nil {
				return err
			}
		}
	}
	return nil
}

func backupArchiveFile(r *zipReader, backupDir string, ti *FileToInstall) error {
	var (
		zf    *zip.File
		zfr   io.ReadCloser
		buf   *os.File
		found bool
		err   error
	)
	if zf, found = r.Files[ti.AbsoluteTo]; found {
		dir := filepath.Join(backupDir, *ti.archive, filepath.Dir(ti.AbsoluteTo))
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		if zfr, err = zf.Open(); err != nil {
			return err
		}
		defer func() { _ = zfr.Close() }()
		if buf, err = os.Create(dir); err != nil {
			return err
		}
		defer func() { _ = buf.Close() }()
		if _, err = io.Copy(buf, zfr); err != nil {
			return err
		}
	}
	return nil
}*/

func installDirectMoveToArchive(state *State, backupDir string) (mods.Result, error) {
	var (
		archive, rel, name, bu string
		installDir, err        = config.Get().GetDir(state.Game, config.GameDirKind)
		z7                     = config.PWD + "\\7z"
	)
	if err != nil {
		return mods.Error, err
	} else if installDir == "" {
		return mods.Error, fmt.Errorf("install directory not found")
	}

	for _, e := range state.ExtractedFiles {
		for _, ti := range e.FilesToInstall() {
			archive = filepath.Join(installDir, *ti.archive)
			if rel, err = filepath.Rel(installDir, ti.AbsoluteTo); err != nil {
				return mods.Error, err
			}
			rel = filepath.Dir(rel)
			name = filepath.Base(ti.Relative)
			// Check if file already exists in the zip file
			cmd := exec.Command(z7, "l", archive, fmt.Sprintf("%s/%s", rel, name))
			if err = cmd.Run(); err == nil {
				// Extract file and move to backup directory
				bu = filepath.Join(backupDir, *ti.archive, rel)
				if err = extractFile(z7, archive, rel, name, bu); err != nil {
					return mods.Error, err
				}
			}
			if err = archiveFile(state, z7, archive, ti.AbsoluteFrom, rel, name); err != nil {
				return mods.Error, err
			}
		}
	}
	return mods.Ok, nil
}

func extractFile(z7, archive, rel, name string, backupDir string) error {
	// Create the target directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}
	_ = os.Remove(filepath.Join(backupDir, name))
	// Extract the file to the target directory
	cmd := exec.Command(z7, "e", archive, "-o"+backupDir, fmt.Sprintf("%s/%s", rel, name))
	if err := cmd.Run(); err != nil {
		return err
	}
	// Remove the file from the zip file
	//cmd = exec.Command(z7, "d", archive, fmt.Sprintf("%s/%s", rel, name))
	//if err := cmd.Run(); err != nil {
	//	return err
	//}
	return nil
}

func archiveFile(state *State, z7 string, archive, absoluteFrom string, rel, name string) (err error) {
	// Modify File Structure
	workingDir := filepath.Dir(absoluteFrom)
	workingDir = filepath.Join(workingDir, "aexta")
	if err = os.MkdirAll(filepath.Join(workingDir, rel), 0755); err != nil {
		return
	}
	defer func() { _ = os.RemoveAll(workingDir) }()

	// Move the file to its new relative location
	to := filepath.Join(workingDir, rel, name)
	if err = os.Rename(absoluteFrom, to); err != nil {
		return
	}

	rel = strings.ReplaceAll(rel, "//", "/")
	workingDir += "\\"

	// Update the zip file
	cmd := exec.Command(z7, "u", archive, workingDir, "-r", "-y")
	if err = cmd.Run(); err != nil {
		return
	}
	files.AppendArchiveFiles(state.Game, state.Mod.ID(), archive, filepath.Join(rel, name))
	return
}
