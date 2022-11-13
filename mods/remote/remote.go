package remote

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/remote/curseforge"
	"github.com/kiamev/moogle-mod-manager/mods/remote/nexus"
)

type Client interface {
	GetFromMod(in *mods.Mod) (found bool, mod *mods.Mod, err error)
	GetFromID(game config.Game, id int) (found bool, mod *mods.Mod, err error)
	GetFromUrl(url string) (found bool, mod *mods.Mod, err error)
	GetNewestMods(game config.Game, lastID int) (result []*mods.Mod, err error)
}

func NewCurseForgeClient() Client {
	return &curseforge.CurseForgeClient{}
}

func NewNexusClient() Client {
	return &nexus.NexusClient{}
}
