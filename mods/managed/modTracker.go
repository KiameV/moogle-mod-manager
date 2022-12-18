package managed

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/discover/remote"
	"github.com/kiamev/moogle-mod-manager/discover/remote/curseforge"
	"github.com/kiamev/moogle-mod-manager/discover/remote/nexus"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/util"
	"path/filepath"
	"strings"
)

const (
	modTrackerName = "tracker.json"
)

var (
	lookup = newGameModLookup()
)

func Initialize(games []config.GameDef) (err error) {
	if err = util.LoadFromFile(filepath.Join(config.PWD, modTrackerName), &lookup); err != nil {
		// first run
		for _, game := range games {
			lookup.Set(game)
		}
		return Save()
	}

	if len(games) != lookup.Len() {
		for _, game := range games {
			if !lookup.Has(game) {
				lookup.Set(game)
			}
		}
	}

	for _, game := range games {
		for _, tm := range lookup.GetMods(game) {
			var mod *mods.Mod
			if err = util.LoadFromFile(tm.MoogleModFile(), &mod); err != nil {
				return
			}
			tm.SetMod(mod)
		}
	}
	return
}

func AddModFromFile(game config.GameDef, file string) (tm mods.TrackedMod, err error) {
	var mod *mods.Mod
	if err = util.LoadFromFile(file, &mod); err != nil {
		return
	}
	if s := mod.Validate(); s != "" {
		return nil, fmt.Errorf("failed to load mod:\n%s", s)
	}

	tm = mods.NewTrackerMod(mod, game)
	if err = AddMod(game, tm); err != nil {
		return nil, err
	}
	return tm, Save()
}

func AddModFromUrl(game config.GameDef, url string) (tm mods.TrackedMod, err error) {
	var (
		mod *mods.Mod
		b   []byte
	)
	if i := strings.Index(url, "?"); i != -1 {
		url = url[:i]
	}
	if nexus.IsNexus(url) {
		if _, mod, err = remote.GetFromUrl(mods.Nexus, url); err != nil {
			return
		}
	} else if curseforge.IsCurseforge(url) {
		if _, mod, err = remote.GetFromUrl(mods.CurseForge, url); err != nil {
			return
		}
	} else {
		if b, err = browser.DownloadAsBytes(url); err != nil {
			return nil, err
		}
		if b[0] == '<' {
			err = xml.Unmarshal(b, &mod)
		} else {
			err = json.Unmarshal(b, &mod)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to load mod: %v", err)
		}
	}

	tm = mods.NewTrackerMod(mod, game)
	if err = AddMod(game, tm); err != nil {
		return nil, err
	}
	return tm, Save()
}

func AddMod(game config.GameDef, tm mods.TrackedMod) error {
	if err := addMod(game, tm); err != nil {
		return err
	}
	return Save()
}

func addMod(game config.GameDef, tm mods.TrackedMod) (err error) {
	if err = tm.Mod().Supports(game); err != nil {
		return
	}

	//tm.Disable()
	if lookup.HasMod(game, tm) {
		return errors.New("mod already added")
	}

	if err = saveMoogle(tm); err != nil {
		return
	}

	lookup.SetMod(game, tm)
	return
}

func GetMods(game config.GameDef) []mods.TrackedMod {
	return lookup.GetMods(game)
}

func GetEnabledMods(game config.GameDef) (result []mods.TrackedMod) {
	for _, tm := range lookup.GetMods(game) {
		if tm.Enabled() {
			result = append(result, tm)
		}
	}
	return
}

func IsModEnabled(game config.GameDef, id mods.ModID) (mod mods.TrackedMod, found bool, enabled bool) {
	if mod, found = TryGetMod(game, id); found {
		enabled = mod.Enabled()
	} else {
		mod = nil
	}
	return
}

func TryGetMod(game config.GameDef, id mods.ModID) (m mods.TrackedMod, found bool) {
	m, found = lookup.GetModByID(game, id)
	return
}

func RemoveMod(game config.GameDef, tm mods.TrackedMod) error {
	lookup.RemoveMod(game, tm)
	return nil
}

func Save() error {
	return util.SaveToFile(filepath.Join(config.PWD, modTrackerName), &lookup)
}

func saveMoogle(tm mods.TrackedMod) (err error) {
	return tm.Save()
}
