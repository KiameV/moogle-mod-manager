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
	"github.com/kiamev/moogle-mod-manager/downloads"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/nexus"
	"github.com/kiamev/moogle-mod-manager/util"
	archiver "github.com/mholt/archiver/v4"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	modTrackerName = "tracker.json"
)

type trackedModsForGame struct {
	Game config.Game        `json:"game"`
	Mods []*mods.TrackedMod `json:"mods"`
}

var lookup = make([]*trackedModsForGame, 6)

func Initialize() (err error) {
	if err = util.LoadFromFile(filepath.Join(config.PWD, modTrackerName), &lookup); err != nil {
		// first run
		for i := range lookup {
			lookup[i] = &trackedModsForGame{Game: config.Game(i)}
		}
		return saveToJson()
	}
	for _, tms := range lookup {
		for _, tm := range tms.Mods {
			var mod *mods.Mod
			if err = util.LoadFromFile(tm.MoogleModFile, &mod); err != nil {
				return
			}
			tm.Mod = mod
		}
	}
	return initializeFiles()
}

func AddModFromFile(game config.Game, file string) (tm *mods.TrackedMod, err error) {
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

func AddModFromUrl(game config.Game, url string) (tm *mods.TrackedMod, err error) {
	var (
		mod *mods.Mod
		b   []byte
	)
	if i := strings.Index(url, "?"); i != -1 {
		url = url[:i]
	}
	if nexus.IsNexus(url) {
		if mod, err = nexus.GetModFromNexus(url); err != nil {
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

func AddMod(game config.Game, tm *mods.TrackedMod) error {
	if err := addMod(game, tm); err != nil {
		return err
	}
	return saveToJson()
}

func addMod(game config.Game, tm *mods.TrackedMod) (err error) {
	if err = tm.GetMod().Supports(game); err != nil {
		return
	}

	tm.Enabled = false
	i := int(game)
	m := lookup[i]
	for i = range m.Mods {
		if m.Mods[i].Mod.ID == tm.Mod.ID {
			return errors.New("mod already added")
		}
	}

	if err = saveMoogle(tm); err != nil {
		return
	}

	m.Mods = append(m.Mods, tm)
	return
}

func UpdateMod(game config.Game, tm *mods.TrackedMod) (err error) {
	if err = tm.GetMod().Supports(game); err != nil {
		return
	}

	if err = DisableMod(game, tm); err != nil {
		return
	}

	tm.Mod = tm.UpdatedMod
	if err = saveMoogle(tm); err != nil {
		return
	}

	tm.UpdatedMod = nil
	return saveToJson()
}

func GetMods(game config.Game) []*mods.TrackedMod {
	return lookup[game].Mods
}

func IsModEnabled(game config.Game, id mods.ModID) (mod *mods.TrackedMod, found bool, enabled bool) {
	if mod, found = TryGetMod(game, id); found {
		enabled = mod.Enabled
	} else {
		mod = nil
	}
	return
}

func TryGetMod(game config.Game, id mods.ModID) (*mods.TrackedMod, bool) {
	var m *mods.TrackedMod
	if gm := lookup[game]; gm != nil {
		for _, m = range gm.Mods {
			if m.Mod.ID == id {
				return m, true
			}
		}
	}
	return nil, false
}

func EnableMod(enabler *mods.ModEnabler) (err error) {
	return downloads.Download(enabler, enableMod)
}

func enableMod(enabler *mods.ModEnabler, err error) {
	var (
		game    = enabler.Game
		tm      = enabler.TrackedMod
		tis     = enabler.ToInstall
		modPath = filepath.Join(config.Get().GetModsFullPath(game), tm.GetDirSuffix())
	)
	if err != nil {
		tm.Enabled = false
		enabler.DoneCallback(mods.Error, err)
		return
	}
	enabler.ShowWorking()

	for _, ti := range tis {
		to := filepath.Join(modPath, ti.Download.Name)
		if err = decompress(*ti.Download.DownloadedArchiveLocation, to); err != nil {
			tm.Enabled = false
			enabler.DoneCallback(mods.Error, err)
			return
		}
		if tm.Mod.ModKind.Kind == mods.Nexus {
			var fi os.FileInfo
			sa := filepath.Join(to, "StreamingAssets")
			if fi, err = os.Stat(sa); err == nil && fi.IsDir() {
				newTo := filepath.Join(to, string(mods.GameToInstallBaseDir(game)))
				_ = os.MkdirAll(newTo, 0777)
				_ = os.Rename(sa, filepath.Join(newTo, "StreamingAssets"))
			} else if !tm.Mod.IsManuallyCreated {
				dir := filepath.Join(to, string(mods.GameToInstallBaseDir(enabler.Game)))
				if _, err = os.Stat(dir); err != nil {
					tm.Enabled = false
					enabler.DoneCallback(mods.Error, errors.New("unsupported nexus mod"))
					return
				}
			}
		}
	}

	for _, ti := range tis {
		AddModFiles(enabler, ti.DownloadFiles, func(result mods.Result, err ...error) {
			if result == mods.Error {
				tm.Enabled = false
			} else {
				tm.Enabled = true
				_ = saveToJson()
			}
			enabler.DoneCallback(result, err...)
		})
	}
}

func decompress(from string, to string) error {
	if fi, err := os.Stat(to); err == nil && fi.IsDir() {
		var fis []os.FileInfo
		if fis, err = ioutil.ReadDir(to); err == nil && len(fis) > 0 {
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

func DisableMod(game config.Game, tm *mods.TrackedMod) (err error) {
	if err = RemoveModFiles(game, tm); err != nil {
		return
	}
	tm.Enabled = false
	return saveToJson()
}

func RemoveMod(game config.Game, tm *mods.TrackedMod) error {
	gm := lookup[game].Mods
	for i, m := range gm {
		if m.Mod.ID == tm.GetModID() {
			if m.Enabled {
				if err := RemoveModFiles(game, tm); err != nil {
					return errors.New("failed to disable mod")
				}
			}
			lookup[game].Mods = append(gm[:i], gm[i+1:]...)
			return saveToJson()
		}
	}
	return fmt.Errorf("failed to find %s", tm.Mod.Name)
}

func saveToJson() error {
	return util.SaveToFile(filepath.Join(config.PWD, modTrackerName), &lookup)
}

func saveMoogle(tm *mods.TrackedMod) (err error) {
	return tm.Save()
}
