package mover

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/action"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/conflict"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/managed"
)

type FileMover interface {
	AddModFiles(enabler *mods.ModEnabler, mmf *managed.ManagedModsAndFiles, files []*mods.DownloadFiles, cr conflict.Result) (err error)
	MoveFiles(game config.GameDef, files []*mods.ModFile, modDir string, toDir string, backupDir string, backedUp *[]*mods.ModFile, movedFiles *[]*mods.ModFile, cr conflict.Result, returnOnFail bool) (err error)
	MoveDirs(game config.GameDef, dirs []*mods.ModDir, modDir string, toDir string, backupDir string, replacedFiles *[]*mods.ModFile, movedFiles *[]*mods.ModFile, cr conflict.Result, returnOnFail bool) (err error)
	MoveFile(a action.FileAction, from, to string, files *[]*mods.ModFile) (err error)
	IsDir(path string) bool
	RemoveModFiles(mf *managed.ModFiles, mmf *managed.ManagedModsAndFiles, tm *mods.TrackedMod) (err error)
}

func NewFileMover(mod *mods.Mod, game config.GameDef) (mover FileMover) {
	var it config.InstallType
	if mod.InstallType != nil {
		it = *mod.InstallType
	} else {
		it = game.DefaultInstallType
	}
	switch it {
	case config.Move:
		mover = &basicFileMover{}
	default:
		panic(fmt.Sprintf("unknown install type: %s", it))
	}
	return
}
