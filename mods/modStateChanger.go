package mods

import (
	"github.com/kiamev/moogle-mod-manager/config"
)

type Result byte

const (
	_ Result = iota
	Ok
	Cancel
	Error
	Working
)

type (
	DoneCallback           func(result Result, err ...error)
	ConflictChoiceCallback func(result Result, choices []*FileConflict, err ...error)
	OnConflict             func(conflicts []*FileConflict, choiceCallback ConflictChoiceCallback)
)

type FileConflict struct {
	File         string
	CurrentModID ModID
	NewModID     ModID
	ChoiceName   string
}

func (c *FileConflict) OnChange(selected string) {
	c.ChoiceName = selected
}

type ModEnabler struct {
	Game         config.GameDef
	TrackedMod   TrackedMod
	ToInstall    []*ToInstall
	OnConflict   OnConflict
	ShowWorking  func()
	DoneCallback DoneCallback
}

func (e *ModEnabler) Kind() Kind {
	return e.TrackedMod.Kind()
}

func NewModEnabler(game config.GameDef, tm TrackedMod, toInstall []*ToInstall, onConflict OnConflict, showWorking func(), doneCallback DoneCallback) *ModEnabler {
	return &ModEnabler{
		Game:         game,
		TrackedMod:   tm,
		ToInstall:    toInstall,
		OnConflict:   onConflict,
		ShowWorking:  showWorking,
		DoneCallback: doneCallback,
	}
}
