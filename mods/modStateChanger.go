package mods

import (
	"github.com/kiamev/moogle-mod-manager/config"
)

type DoneCallback func(err error)

type ModEnabler struct {
	Game         config.Game
	TrackedMod   *TrackedMod
	ToInstall    []*ToInstall
	DoneCallback DoneCallback
}

func NewModEnabler(game config.Game, tm *TrackedMod, toInstall []*ToInstall, doneCallback DoneCallback) *ModEnabler {
	return &ModEnabler{
		Game:         game,
		TrackedMod:   tm,
		ToInstall:    toInstall,
		DoneCallback: doneCallback,
	}
}
