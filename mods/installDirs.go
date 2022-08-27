package mods

import "github.com/kiamev/moogle-mod-manager/config"

type InstallBaseDir string

const (
	InstallDir_I   InstallBaseDir = "FINAL FANTASY_Data"
	InstallDir_II  InstallBaseDir = "FINAL FANTASY II_Data"
	InstallDir_III InstallBaseDir = "FINAL FANTASY III_Data"
	InstallDir_IV  InstallBaseDir = "FINAL FANTASY IV_Data"
	InstallDir_V   InstallBaseDir = "FINAL FANTASY V_Data"
	InstallDir_VI  InstallBaseDir = "FINAL FANTASY VI_Data"

	StreamingAssetsDir = "StreamingAssets"
)

func GameToInstallBaseDir(game config.Game) InstallBaseDir {
	switch game {
	case config.I:
		return InstallDir_I
	case config.II:
		return InstallDir_II
	case config.III:
		return InstallDir_III
	case config.IV:
		return InstallDir_IV
	case config.V:
		return InstallDir_V
	default:
		return InstallDir_VI
	}
}
