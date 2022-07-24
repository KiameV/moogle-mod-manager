package io

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"io/ioutil"
	"os"
	"path"
)

func MoveFiles(files []mods.ModFile, modDir string, game config.Game) (err error) {
	var (
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
	}
	return
}

func RevertMoveFiles(files []string, game config.Game) (err error) {
	var (
		toDir     = config.GetModDir(game)
		backupDir = config.GetBackupDir(game)
	)
	for _, f := range files {
		if err = os.Remove(path.Join(toDir, f)); err != nil {
			break
		}
		if err = os.Rename(path.Join(backupDir, f), path.Join(toDir, f)); err != nil {
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
