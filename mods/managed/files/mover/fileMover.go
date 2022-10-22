package mover

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/action"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/conflict"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/managed"
)

type FileMover interface {
	AddModFiles(enabler *mods.ModEnabler, mmf *managed.ManagedModsAndFiles, files []*mods.DownloadFiles, cr conflict.Result) (err error)
	MoveFiles(game config.Game, files []*mods.ModFile, modDir string, toDir string, backupDir string, backedUp *[]*mods.ModFile, movedFiles *[]*mods.ModFile, cr conflict.Result, returnOnFail bool) (err error)
	MoveDirs(game config.Game, dirs []*mods.ModDir, modDir string, toDir string, backupDir string, replacedFiles *[]*mods.ModFile, movedFiles *[]*mods.ModFile, cr conflict.Result, returnOnFail bool) (err error)
	MoveFile(a action.FileAction, from, to string, files *[]*mods.ModFile) (err error)
	IsDir(path string) bool
	RemoveModFiles(mf *managed.ModFiles, mmf *managed.ManagedModsAndFiles, tm *mods.TrackedMod) (err error)
}

func NewFileMover(game config.Game) (mover FileMover) {
	switch game {
	case config.ChronoCross:
		// TODO Chrono Cross
		//mover = &archiveFileMover{}
	case config.BofIII, config.BofIV:
		// TODO BoF
		panic("BoF not implemented")
	default:
		mover = &basicFileMover{}
	}
	return
}
