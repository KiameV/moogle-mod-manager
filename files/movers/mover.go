package movers

import (
	"errors"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type (
	Mover interface {
		CompileFilesToMove(dir string) (files []string, err error)
		MoveFiles(files []string, from, destination, backupDir string) (backups []string, err error)
	}
)

func NewMover(game config.GameDef, tm mods.TrackedMod) (mover Mover, err error) {
	switch tm.Mod().InstallType(game) {
	case config.Move:
		mover = &directMover{}
	case config.MoveToArchive:
		err = errors.New("archive install type not implemented")
	default:
		err = fmt.Errorf("unknown install type [%s]", tm.Mod().InstallType(game))
	}
	return
}
