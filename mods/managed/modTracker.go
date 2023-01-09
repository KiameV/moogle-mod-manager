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
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/util"
	"os"
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
		return save()
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
			var mod mods.Mod
			if err = mod.LoadFromFile(tm.MoogleModFile()); err != nil {
				return
			}
			tm.SetMod(&mod)
		}
	}
	return
}

func AddModFromFile(game config.GameDef, file string) (mods.TrackedMod, error) {
	mod := &mods.Mod{}
	if err := mod.LoadFromFile(file); err != nil {
		return nil, err
	}
	if s := mod.Validate(); s != "" {
		return nil, fmt.Errorf("failed to load mod:\n%s", s)
	}
	return AddMod(game, mod)
}

func AddModFromUrl(game config.GameDef, url string) (mods.TrackedMod, error) {
	var (
		mod = &mods.Mod{}
		b   []byte
		err error
	)
	if i := strings.Index(url, "?"); i != -1 {
		url = url[:i]
	}
	if nexus.IsNexus(url) {
		if _, mod, err = remote.GetFromUrl(mods.Nexus, url); err != nil {
			return nil, err
		}
	} else if curseforge.IsCurseforge(url) {
		if _, mod, err = remote.GetFromUrl(mods.CurseForge, url); err != nil {
			return nil, err
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

	return AddMod(game, mod)
}

func AddMod(game config.GameDef, mod *mods.Mod) (tm mods.TrackedMod, err error) {
	var found bool
	if tm, found = lookup.GetModByID(game, mod.ID()); found {
		return
	}
	tm = mods.NewTrackerMod(mod, state.CurrentGame)
	err = addMod(game, tm)
	return
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
	return save()
}

func GetMods(game config.GameDef) []mods.TrackedMod {
	return lookup.GetMods(game)
}

func DisableMod(tm mods.TrackedMod) error {
	tm.Disable()
	return save()
}

func EnableMod(tm mods.TrackedMod) error {
	tm.Enable()
	return save()
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
	_ = os.RemoveAll(filepath.Dir(tm.MoogleModFile()))
	for _, ti := range tm.Mod().Downloadables {
		if ti.DownloadedArchiveLocation != nil && *ti.DownloadedArchiveLocation != "" {
			dir := filepath.Dir(string(*ti.DownloadedArchiveLocation))
			dir = filepath.Dir(dir)
			if strings.Contains(dir, config.Get().DownloadDir) {
				_ = os.RemoveAll(dir)
			}
		}
	}
	return save()
}

func save() error {
	return util.SaveToFile(filepath.Join(config.PWD, modTrackerName), &lookup)
}

func saveMoogle(tm mods.TrackedMod) (err error) {
	return tm.Save()
}
