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
	gameMods struct {
		GameMods map[string]*mods.ModLookupConc[*mods.TrackedModConc] `json:"Mods"`
	}
)

func newGameModLookup() gameModLookup {
	return &gameMods{GameMods: make(map[string]*mods.ModLookupConc[*mods.TrackedModConc])}
}

func (gm *gameMods) Has(game config.GameDef) bool {
	_, ok := gm.GameMods[string(game.ID())]
	return ok
}

func (gm *gameMods) HasMod(game config.GameDef, tm mods.TrackedMod) bool {
	if l, ok := gm.GameMods[string(game.ID())]; ok {
		return l.Has(tm.(*mods.TrackedModConc))
	}
	return false
}

func (gm *gameMods) Len() int {
	return len(gm.GameMods)
}

func (gm *gameMods) GetMod(game config.GameDef, tm mods.TrackedMod) (mods.TrackedMod, bool) {
	if l, found := gm.GameMods[string(game.ID())]; found {
		return l.Get(tm.(*mods.TrackedModConc))
	}
	return nil, false
}

func (gm *gameMods) GetMods(game config.GameDef) (tms []mods.TrackedMod) {
	if l, found := gm.GameMods[string(game.ID())]; found {
		all := l.All()
		tms = make([]mods.TrackedMod, len(all))
		for i, tm := range all {
			tms[i] = tm
		}
	}
	return
}

func (gm *gameMods) GetModByID(game config.GameDef, id mods.ModID) (tm mods.TrackedMod, ok bool) {
	if l, found := gm.GameMods[string(game.ID())]; found {
		return l.GetByID(id)
	}
	return
}

func (gm *gameMods) ModCount(game config.GameDef) int {
	if l, found := gm.GameMods[string(game.ID())]; found {
		return l.Len()
	}
	return 0
}

func (gm *gameMods) RemoveMod(game config.GameDef, tm mods.TrackedMod) {
	if l, found := gm.GameMods[string(game.ID())]; found {
		l.Remove(tm.(*mods.TrackedModConc))
	}
}

func (gm *gameMods) Set(game config.GameDef) {
	if _, found := gm.GameMods[string(game.ID())]; !found {
		c := mods.NewModLookup[*mods.TrackedModConc]()
		gm.GameMods[string(game.ID())] = c.(*mods.ModLookupConc[*mods.TrackedModConc])
	}
}

func (gm *gameMods) SetMod(game config.GameDef, tm mods.TrackedMod) {
	if l, found := gm.GameMods[string(game.ID())]; found {
		l.Set(tm.(*mods.TrackedModConc))
	}
}
