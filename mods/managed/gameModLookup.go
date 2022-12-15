package managed

import (
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type (
	gameModLookup interface {
		Has(game config.GameDef) bool
		HasMod(game config.GameDef, tm mods.TrackedMod) bool
		Len() int
		GetMod(game config.GameDef, tm mods.TrackedMod) (found mods.TrackedMod, ok bool)
		GetMods(game config.GameDef) (found []mods.TrackedMod)
		GetModByID(game config.GameDef, id mods.ModID) (mods.TrackedMod, bool)
		ModCount(game config.GameDef) int
		RemoveMod(game config.GameDef, tm mods.TrackedMod)
		Set(game config.GameDef)
		SetMod(game config.GameDef, tm mods.TrackedMod)
	}
	gameModLookup_ struct {
		Games []*gameMods `json:"Games"`
	}
	gameMods struct {
		Game config.GameID          `json:"GameID"`
		Mods []*mods.TrackedModConc `json:"Mods"`
	}
)

func newGameModLookup() gameModLookup {
	return &gameModLookup_{}
}

func (l *gameModLookup_) Has(game config.GameDef) bool {
	_, ok := l.getGameMods(game)
	return ok
}

func (l *gameModLookup_) HasMod(game config.GameDef, tm mods.TrackedMod) bool {
	if gm, ok := l.getGameMods(game); ok {
		if _, ok = gm.getMod(tm.ID()); ok {
			return true
		}
	}
	return false
}

func (l *gameModLookup_) Len() int {
	return len(l.Games)
}

func (l *gameModLookup_) GetMod(game config.GameDef, tm mods.TrackedMod) (found mods.TrackedMod, ok bool) {
	var gm *gameMods
	if gm, ok = l.getGameMods(game); ok {
		found, ok = gm.getMod(tm.ID())
	}
	return
}

func (l *gameModLookup_) GetMods(game config.GameDef) (found []mods.TrackedMod) {
	if gm, ok := l.getGameMods(game); ok {
		found = make([]mods.TrackedMod, len(gm.Mods))
		for i, tm := range gm.Mods {
			found[i] = tm
		}
	}
	return
}

func (l *gameModLookup_) GetModByID(game config.GameDef, id mods.ModID) (tm mods.TrackedMod, ok bool) {
	var gm *gameMods
	if gm, ok = l.getGameMods(game); ok {
		for _, tm = range gm.Mods {
			if tm.ID() == id {
				ok = true
				break
			}
		}
	}
	return
}

func (l *gameModLookup_) ModCount(game config.GameDef) int {
	if gm, ok := l.getGameMods(game); ok {
		return len(gm.Mods)
	}
	return 0
}

func (l *gameModLookup_) RemoveMod(game config.GameDef, tm mods.TrackedMod) {
	if gm, ok := l.getGameMods(game); ok {
		for i, mod := range gm.Mods {
			if mod.ID() == tm.ID() {
				gm.Mods = append(gm.Mods[:i], gm.Mods[i+1:]...)
				break
			}
		}
	}
}

func (l *gameModLookup_) Set(game config.GameDef) {
	for _, gm := range l.Games {
		if gm.Game == game.ID() {
			return
		}
	}
	l.Games = append(l.Games, &gameMods{Game: game.ID()})
}

func (l *gameModLookup_) SetMod(game config.GameDef, tm mods.TrackedMod) {
	if gm, ok := l.getGameMods(game); ok {
		if _, ok = gm.getMod(tm.ID()); ok {
			return
		}
		gm.Mods = append(gm.Mods, tm.(*mods.TrackedModConc))
	}
}

func (l *gameModLookup_) getGameMods(game config.GameDef) (gm *gameMods, ok bool) {
	for _, gm = range l.Games {
		if gm.Game == game.ID() {
			ok = true
			break
		}
	}
	return
}

func (gm *gameMods) getMod(id mods.ModID) (tm mods.TrackedMod, ok bool) {
	for _, tm = range gm.Mods {
		if tm.ID() == id {
			ok = true
			break
		}
	}
	return
}
