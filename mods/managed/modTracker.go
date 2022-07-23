package managed

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	moogleModName = "mod.moogle"

	tempDir = "temp"

	modTrackerName = "tracker.json"
)

type trackedModsForGame struct {
	Game config.Game         `json:"game"`
	Mods []*model.TrackedMod `json:"mods"`
}

// lookup first slice is the game, second slice is the mod
var lookup = make([]*trackedModsForGame, 6)

func Initialize() (err error) {
	var (
		f = filepath.Join(config.PWD, modTrackerName)
		b []byte
	)
	if _, err = os.Stat(f); err != nil {
		// ignore, probably first run
		for i := range lookup {
			lookup[i] = &trackedModsForGame{Game: config.Game(i)}
		}
		return saveToJson()
	} else {
		if b, err = ioutil.ReadFile(f); err != nil {
			return
		}
		if err = json.Unmarshal(b, &lookup); err != nil {
			return
		}
	}
	for _, tms := range lookup {
		for _, tm := range tms.Mods {
			if b, err = readFile(filepath.Join(tm.Dir, moogleModName)); err != nil {
				return
			}
			var mod *mods.Mod
			if err = json.Unmarshal(b, &mod); err != nil {
				return
			}
			tm.Mod = mod
		}
	}
	return nil
}

func AddModFromFile(game config.Game, file string) (err error) {
	var (
		b   []byte
		mod *mods.Mod
	)
	if b, err = readFile(file); err != nil {
		return
	}

	ext := filepath.Ext(file)
	if ext == ".xml" {
		err = xml.Unmarshal(b, &mod)
	} else if ext == ".json" {
		err = json.Unmarshal(b, &mod)
	} else {
		return fmt.Errorf("unknown file extension: %s", ext)
	}
	if err != nil {
		return fmt.Errorf("failed to load mod: %v", err)
	}
	if s := mod.Validate(); s != "" {
		return fmt.Errorf("failed to load mod:\n%s", s)
	}
	return AddMod(game, model.NewTrackerMod(game, mod))
}

func AddModFromUrl(game config.Game, url string) error {
	b, err := browser.DownloadAsBytes(url)
	if err != nil {
		return err
	}
	var mod *mods.Mod
	if b[0] == '<' {
		err = xml.Unmarshal(b, &mod)
	} else {
		err = json.Unmarshal(b, &mod)
	}
	if err != nil {
		return fmt.Errorf("failed to load mod: %v", err)
	}
	return AddMod(game, model.NewTrackerMod(game, mod))
}

func AddMod(game config.Game, tm *model.TrackedMod) (err error) {
	if err = tm.GetMod().Supports(game); err != nil {
		return
	}

	tm.Enabled = false
	for _, g := range tm.Mod.Games {
		i := int(config.NameToGame(g.Name))
		m := lookup[i]
		for i = range m.Mods {
			if m.Mods[i].Mod.ID == tm.Mod.ID {
				return errors.New("mod already added")
			}
		}
		m.Mods = append(m.Mods, tm)
	}

	var (
		b []byte
		f *os.File
	)
	if b, err = json.MarshalIndent(tm.Mod, "", "\t"); err != nil {
		return
	}
	if _, err = os.Stat(tm.GetDir()); os.IsNotExist(err) {
		if err = os.MkdirAll(tm.GetDir(), 0755); err != nil {
			return
		}
	}
	if f, err = os.Create(filepath.Join(tm.GetDir(), moogleModName)); err != nil {
		return
	}
	defer func() { _ = f.Close() }()
	if _, err = f.Write(b); err != nil {
		return
	}

	return saveToJson()
}

func GetMods(game config.Game) []*model.TrackedMod { return lookup[game].Mods }

func RemoveMod(game config.Game, modID string) error {
	gm := lookup[game].Mods
	for i, m := range gm {
		if m.Mod.ID != modID {
			return fmt.Errorf("failed to find %s", modID)
		}
		if m.Enabled {
			if err := RemoveModFiles(game, modID); err != nil {

			}
		}
		lookup[game].Mods = append(gm[:i], gm[i+1:]...)
		break
	}
	return saveToJson()
}

/*
func readModDef(dir string) (mod *mods.Mod, err error) {
	if err = readXml(dir, "mod.xml", &mod); err != nil {
		err = readJson(dir, "mod.json", &mod)
	}
	return
}
*/
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
	b, err := json.MarshalIndent(&lookup, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(config.PWD, modTrackerName), b, 0755)
}
