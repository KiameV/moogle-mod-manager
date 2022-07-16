package state

import "github.com/kiamev/pr-modsync/config"

type GUI byte

const (
	None GUI = iota
	LocalMods
	Configure
)

var (
	CurrentGame *config.Game
	CurrentUI   GUI
	Errors      []error
)
