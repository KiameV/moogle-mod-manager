package mods

import (
	"encoding/xml"
	"fmt"
	"github.com/kiamev/pr-modsync/config"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	modXmlName     = "mod.xml"
	modEnabledName = "mod.enabled"
)

type GameMods interface {
	GetGameName() string
	GetMods() []*GameMod
}

type gameMods struct {
	config.Game
	mods []*GameMod
}

type GameMod struct {
	Mod     Mod
	Enabled bool
}

var lookup = make(map[config.Game]GameMods)

func GetGameMods(game config.Game) (gm GameMods, err error) {
	var ok bool
	if gm, ok = lookup[game]; !ok {
		if gm, err = newGameMods(game); err == nil {
			lookup[game] = gm
		}
	}
	return
}

func newGameMods(game config.Game) (GameMods, error) {
	gm := &gameMods{Game: game}
	return gm, gm.loadMods()
}

func (gm *gameMods) GetMods() []*GameMod { return gm.mods }

func (gm *gameMods) loadMods() (err error) {
	var (
		dir   = config.Get().GetModDir(gm.Game)
		files []os.FileInfo
	)
	if _, err = os.Stat(dir); err != nil {
		// No mods loaded yet
		err = nil
		return
	}

	if files, err = ioutil.ReadDir(dir); err != nil {
		return
	}
	for _, f := range files {
		if f.IsDir() {
			if err = gm.readModDir(filepath.Join(dir, f.Name())); err != nil {
				return
			}
		}
	}
	return
}

func (gm *gameMods) readModDir(dir string) (err error) {
	var (
		b       []byte
		gameMod GameMod
	)
	if b, err = os.ReadFile(filepath.Join(dir, modXmlName)); err != nil {
		err = fmt.Errorf("failed to find %s for mod in %s: %v", modXmlName, dir, err)
		return
	}
	if err = xml.Unmarshal(b, &gameMod.Mod); err != nil {
		err = fmt.Errorf("failed to read %s for mod in %s: %v", modXmlName, dir, err)
		return
	}
	gameMod.Mod.Preview = filepath.Join(dir, gameMod.Mod.Preview)

	_, err = os.Stat(filepath.Join(dir, modEnabledName))
	gameMod.Enabled = err == nil
	err = nil

	gm.mods = append(gm.mods, &gameMod)
	return
}

func (gm *gameMods) GetGameName() (name string) {
	switch gm.Game {
	case config.I:
		name = "I"
	case config.II:
		name = "II"
	case config.III:
		name = "III"
	case config.IV:
		name = "IV"
	case config.V:
		name = "V"
	case config.VI:
		name = "VI"
	}
	name = fmt.Sprintf("Final Fantasy %s PR Mods", name)
	return
}
