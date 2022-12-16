package managed

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	unarr "github.com/gen2brain/go-unarr"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/discover/remote"
	"github.com/kiamev/moogle-mod-manager/discover/remote/curseforge"
	"github.com/kiamev/moogle-mod-manager/discover/remote/nexus"
	"github.com/kiamev/moogle-mod-manager/downloads"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files"
	"github.com/kiamev/moogle-mod-manager/mods/managed/files/managed"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/util"
	archiver "github.com/mholt/archiver/v4"
	"io"
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
		return saveToJson()
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
	return managed.InitializeManagedFiles()
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
	return tm, saveToJson()
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
	return tm, saveToJson()
}

func AddMod(game config.GameDef, tm mods.TrackedMod) error {
	if err := addMod(game, tm); err != nil {
		return err
	}
	return saveToJson()
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

func UpdateMod(game config.GameDef, tm mods.TrackedMod) (err error) {
	if tm.UpdatedMod() == nil {
		return errors.New("no update available")
	}

	if err = tm.Mod().Supports(game); err != nil {
		return
	}

	if tm.Enabled() {
		if err = DisableMod(game, tm); err != nil {
			return
		}
	}

	tm.SetMod(tm.UpdatedMod())
	if err = saveMoogle(tm); err != nil {
		return
	}

	tm.SetUpdatedMod(nil)
	return saveToJson()
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

func EnableMod(enabler *mods.ModEnabler) (err error) {
	if err = canInstall(enabler); err != nil {
		return
	}
	return downloads.Download(enabler, enableMod)
}

func enableMod(enabler *mods.ModEnabler, err error) {
	var (
		game    = enabler.Game
		tm      = enabler.TrackedMod
		tis     = enabler.ToInstall
		modPath = filepath.Join(config.Get().GetModsFullPath(game), tm.DirSuffix())
	)
	if err != nil {
		tm.Disable()
		enabler.DoneCallback(mods.Error, err)
		return
	}
	enabler.ShowWorking()

	for _, ti := range tis {
		var (
			to   = filepath.Join(modPath, ti.Download.Name)
			kind = tm.Kind()
		)
		if err = decompress(*ti.Download.DownloadedArchiveLocation, to); err != nil {
			tm.Disable()
			enabler.DoneCallback(mods.Error, err)
			return
		}
		if kind == mods.Nexus || kind == mods.CurseForge {
			var fi os.FileInfo
			sa := filepath.Join(to, "StreamingAssets")
			if fi, err = os.Stat(sa); err == nil && fi.IsDir() {
				newTo := filepath.Join(to, string(game.BaseDir()))
				_ = os.MkdirAll(newTo, 0777)
				_ = os.Rename(sa, filepath.Join(newTo, "StreamingAssets"))
			} else if !tm.Mod().IsManuallyCreated {
				dir := filepath.Join(to, string(game.BaseDir()))
				if _, err = os.Stat(dir); err != nil {
					tm.Disable()
					enabler.DoneCallback(mods.Error, errors.New("unsupported nexus mod"))
					return
				}
			}
		}
	}

	for _, ti := range tis {
		files.AddModFiles(enabler, ti.DownloadFiles, func(result mods.Result, err ...error) {
			if result == mods.Error || result == mods.Cancel {
				tm.Disable()
			} else {
				tm.Enable()
				// Find any mods that are now disabled because all the files have been replaced by other mods
				for _, mod := range GetEnabledMods(enabler.Game) {
					if !managed.HasManagedFiles(enabler.Game, mod.ID()) {
						mod.Disable()
					}
				}
				_ = saveToJson()
			}
			enabler.DoneCallback(result, err...)
		})
	}
}

func decompress(from string, to string) error {
	if fi, err := os.Stat(to); err == nil && fi.IsDir() {
		var fis []os.DirEntry
		if fis, err = os.ReadDir(to); err == nil && len(fis) > 0 {
			return nil
		}
	}

	if filepath.Ext(from) == ".rar" {
		handler := func(ctx context.Context, f archiver.File) (err error) {
			if !f.IsDir() {
				var r io.ReadCloser
				if r, err = f.Open(); err != nil {
					return
				}
				defer func() { _ = r.Close() }()

				fp := filepath.Join(to, f.NameInArchive)
				if err = os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
					return
				}
				buf := new(strings.Builder)
				if _, err = io.Copy(buf, r); err != nil {
					return
				}
				var file *os.File
				if file, err = os.Create(fp); err != nil {
					return
				}
				defer func() { _ = file.Close() }()

				_, err = file.WriteString(buf.String())
			}
			return
		}
		f, err := os.Open(from)
		if err != nil {
			return err
		}
		return archiver.Rar{}.Extract(context.Background(), f, nil, handler)
	}
	a, err := unarr.NewArchive(from)
	if err != nil {
		return err
	}
	defer func() { _ = a.Close() }()

	if err = os.MkdirAll(to, 0777); err != nil {
		return err
	}

	_, err = a.Extract(to)
	return err
}

func DisableMod(game config.GameDef, tm mods.TrackedMod) (err error) {
	if err = canDisable(game, tm); err != nil {
		return
	}
	if err = files.RemoveModFiles(game, tm); err != nil {
		return
	}
	tm.Disable()
	return saveToJson()
}

func RemoveMod(game config.GameDef, tm mods.TrackedMod) error {
	lookup.RemoveMod(game, tm)
	return nil
}

func saveToJson() error {
	return util.SaveToFile(filepath.Join(config.PWD, modTrackerName), &lookup)
}

func saveMoogle(tm mods.TrackedMod) (err error) {
	return tm.Save()
}

func canInstall(enabler *mods.ModEnabler) error {
	var (
		tm      = enabler.TrackedMod
		c       = tm.Mod().ModCompatibility
		mc      *mods.ModCompat
		mod     mods.TrackedMod
		found   bool
		enabled bool
	)
	if c != nil {
		if len(c.Forbids) > 0 {
			for _, mc = range c.Forbids {
				if mod, found, enabled = IsModEnabled(state.CurrentGame, mc.ModID()); found && enabled {
					return fmt.Errorf("[%s] cannot be enabled because [%s] is enabled", tm.DisplayName(), mod.DisplayName())
				}
			}
		}
		if len(c.Requires) > 0 {
			for _, mc = range c.Requires {
				mod, found, enabled = IsModEnabled(state.CurrentGame, mc.ModID())
				if !found {
					return fmt.Errorf("[%s] cannot be enabled because [%s] is not enabled", tm.DisplayName(), mc.ModID())
				} else if !enabled {
					return fmt.Errorf("[%s] cannot be enabled because [%s] is not enabled", tm.DisplayName(), mod.DisplayName())
				}
			}
		}
	}
	return nil
}

func canDisable(game config.GameDef, tm mods.TrackedMod) error {
	var (
		c  = tm.Mod().ModCompatibility
		id = tm.Mod().ModID
		mc *mods.ModCompat
	)
	for _, m := range lookup.GetMods(game) {
		if m.ID() != tm.ID() && m.Enabled() {
			if c = m.Mod().ModCompatibility; c != nil {
				if len(c.Requires) > 0 {
					for _, mc = range c.Requires {
						if mc.ModID() == id {
							return fmt.Errorf("[%s] cannot be disabled because [%s] is enabled", tm.DisplayName(), m.DisplayName())
						}
					}
				}
			}
		}
	}
	return nil
}
