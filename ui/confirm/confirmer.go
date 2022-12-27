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
	k := params.Mod.Kinds()
	if k.Is(mods.Nexus) {
		return newNexusConfirmer(params)
	}
	if k.IsHosted() {
		return newHostedConfirmer(params)
	}
	return newBypassConfirmer(params)
}
