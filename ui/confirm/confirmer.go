package confirm

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type (
	Params struct {
		Game      config.GameDef
		Mod       mods.TrackedMod
		ToInstall []*mods.ToInstall
	}
	Confirmer interface {
		Downloads(done func(mods.Result)) error
	}
)

func NewParams(game config.GameDef, mod mods.TrackedMod, toInstall []*mods.ToInstall) Params {
	return Params{
		Game:      game,
		Mod:       mod,
		ToInstall: toInstall,
	}
}

func NewConfirmer(params Params) Confirmer {
	switch params.Mod.Kind() {
	case mods.Nexus:
		return newNexusConfirmer(params)
	case mods.Hosted:
		return newHostedConfirmer(params)
	}
	return newBypassConfirmer(params)
}
