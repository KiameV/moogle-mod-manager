package managed

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	modXmlName  = "mod.xml"
	modJsonName = "mod.json"

	tempDir = "temp"

	modTrackerName = "tracker.json"
)

type GameMod struct {
	ModDir  string   `json:"modDir"`
	Enabled bool     `json:"enabled"`
	Mod     mods.Mod `json:"-"`
}

type gameMods struct {
	Game config.Game `json:"game"`
	Mods []*GameMod  `json:"mods"`
}

// lookup first slice is the game, second slice is the mod
var lookup = make([]*gameMods, 6)

func Initialize() error {
	var (
		f   = filepath.Join(config.PWD, modTrackerName)
		b   []byte
		err error
	)
	if _, err = os.Stat(f); err != nil {
		// ignore, probably first run
		for i := range lookup {
			lookup[i] = &gameMods{Game: config.Game(i)}
		}
		return saveToJson()
	} else {
		if b, err = ioutil.ReadFile(f); err != nil {
			return err
		}
		if err = json.Unmarshal(b, &lookup); err != nil {
			return err
		}
	}
	for _, gameMod := range lookup {
		dir := filepath.Join(config.PWD, config.GetModDir(gameMod.Game))
		for _, m := range gameMod.Mods {
			if m.Mod, err = readModDef(filepath.Join(dir, m.ModDir)); err != nil {
				return err
			}
		}
	}
	return nil
}

func AddModFromFile(game config.Game, file string) (err error) {
	var (
		b  []byte
		gm GameMod
	)
	if b, err = readFile(file); err != nil {
		return
	}
	if err = xml.Unmarshal(b, &gm.Mod); err != nil {
		if err = json.Unmarshal(b, &gm.Mod); err != nil {
			return
		}
	}
	gm.Enabled = false
	if len(gm.Mod.Game) > 1 {
		gm.ModDir = filepath.Join("mods", "shared")
	} else {
		gm.ModDir = filepath.Join("mods", config.GetGameDir(game))
	}
	for _, g := range gm.Mod.Game {
		m := lookup[config.NameToGame(g.Name)]
		m.Mods = append(m.Mods, &gm)
	}
	return saveToJson()
}

func AddModFromUrl(game config.Game, url string) error {
	file, err := browser.Download(url, filepath.Join(config.PWD, tempDir))
	if err != nil {
		return err
	}
	return AddModFromFile(game, file)
}

func GetMods(game config.Game) []*GameMod { return lookup[game].Mods }

func RemoveMod(game config.Game, modID string) error {
	gm := lookup[game].Mods
	for i, m := range gm {
		if m.Enabled {
			if err := RemoveModFiles(game, modID); err != nil {

			}
		}
		if m.Mod.ID == modID {
			gm = append(gm[:i], gm[i+1:]...)
			break
		}
	}
	return saveToJson()
}

func readModDef(dir string) (mod mods.Mod, err error) {
	if err = readXml(dir, modXmlName, &mod); err != nil {
		err = readJson(dir, modJsonName, &mod)
	}
	return
}

func readXml(dir string, name string, to interface{}) (err error) {
	var (
		f = filepath.Join(dir, name)
		b []byte
	)
	if b, err = readFile(f); err != nil {
		return
	}
	err = xml.Unmarshal(b, to)
	return
}

func readJson(dir string, name string, to interface{}) (err error) {
	var (
		f = filepath.Join(dir, name)
		b []byte
	)
	if b, err = readFile(f); err != nil {
		return
	}
	err = xml.Unmarshal(b, to)
	return
}

func readFile(f string) (b []byte, err error) {
	if _, err = os.Stat(f); err != nil {
		err = fmt.Errorf("failed to find %s: %v", f, err)
		return
	}
	if b, err = ioutil.ReadFile(f); err != nil {
		err = fmt.Errorf("failed to read %s: %v", f, err)
	}
	return
}

func saveToJson() error {
	b, err := json.Marshal(&lookup)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(config.PWD, modTrackerName), b, 0755)
}
