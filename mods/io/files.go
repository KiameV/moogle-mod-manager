package io

import (
	"github.com/kiamev/pr-modsync/mods"
	"io/ioutil"
	"os"
	"path/filepath"
)

func MoveFiles(files []mods.ModFile, fromDir, toDir, backupDir string) (err error) {
	for _, f := range files {
		if err = os.Rename(filepath.Join(toDir, f.To), filepath.Join(backupDir, f.To)); err != nil {
			break
		}
		if err = copy(filepath.Join(fromDir, f.From), filepath.Join(toDir, f.To)); err != nil {
			break
		}
	}
	return
}

func RevertMoveFiles(files []string, toDir, backupDir string) (err error) {
	for _, f := range files {
		if err = os.Remove(filepath.Join(toDir, f)); err != nil {
			break
		}
		if err = os.Rename(filepath.Join(backupDir, f), filepath.Join(toDir, f)); err != nil {
			break
		}
	}
	return
}

func copy(from, to string) error {
	input, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(to, input, 0777)
}
