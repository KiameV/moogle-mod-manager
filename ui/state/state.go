package state

import (
	"fyne.io/fyne/v2"
	"github.com/kiamev/pr-modsync/config"
)

type GUI byte

const (
	None GUI = iota
	LocalMods
	Configure
)

var (
	CurrentGame *config.Game
	CurrentUI   GUI
	App         fyne.App
	Errors      []error
)
