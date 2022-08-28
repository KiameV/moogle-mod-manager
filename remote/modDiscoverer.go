package remote

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
)

func GetMods(game config.Game) ([]*mods.Mod, error) {
	return getNexusMods(game)
}
