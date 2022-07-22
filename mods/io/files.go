package io

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"io/ioutil"
	"os"
	"path/filepath"
)

func MoveFiles(files []mods.ModFile, modDir string, game config.Game) (err error) {
	var (
		toDir     = config.GetModDir(game)
		backupDir = config.GetBackupDir(game)
	)
	for _, f := range files {
		if err = os.Rename(filepath.Join(toDir, f.To), filepath.Join(backupDir, f.To)); err != nil {
			break
		}
		if err = copy(filepath.Join(modDir, f.From), filepath.Join(toDir, f.To)); err != nil {
			break
		}
	}
	return
}

func RevertMoveFiles(files []string, game config.Game) (err error) {
	var (
		toDir     = config.GetModDir(game)
		backupDir = config.GetBackupDir(game)
	)
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
