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
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	wu "github.com/kiamev/moogle-mod-manager/ui/util"
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

type trackedModsForGame struct {
	Game config.Game         `json:"game"`
	Mods []*model.TrackedMod `json:"mods"`
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
	return
}

func AddModFromFile(game config.Game, file string) (tm *model.TrackedMod, err error) {
	var mod *mods.Mod
	if err = util.LoadFromFile(file, &mod); err != nil {
		return
	}
	if s := mod.Validate(); s != "" {
		return nil, fmt.Errorf("failed to load mod:\n%s", s)
	}

	tm = model.NewTrackerMod(mod, game)
	if err = AddMod(game, tm); err != nil {
		return nil, err
	}
	return tm, saveToJson()
}

func AddModFromUrl(game config.Game, url string) (tm *model.TrackedMod, err error) {
	var b []byte
	if b, err = browser.DownloadAsBytes(url); err != nil {
		return nil, err
	}
	var mod *mods.Mod
	if b[0] == '<' {
		err = xml.Unmarshal(b, &mod)
	} else {
		err = json.Unmarshal(b, &mod)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load mod: %v", err)
	}

	tm = model.NewTrackerMod(mod, game)
	if err = AddMod(game, tm); err != nil {
		return nil, err
	}
	return tm, saveToJson()
}

func AddMod(game config.Game, tm *model.TrackedMod) error {
	if err := addMod(game, tm); err != nil {
		return err
	}
	return saveToJson()
}

func addMod(game config.Game, tm *model.TrackedMod) (err error) {
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
	}

	if err = saveMoogle(tm); err != nil {
		return
	}

	for _, g := range tm.Mod.Games {
		i := int(config.NameToGame(g.Name))
		m := lookup[i]
		m.Mods = append(m.Mods, tm)
	}
	return
}

func UpdateMod(game config.Game, tm *model.TrackedMod) (err error) {
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

func GetMods(game config.Game) []*model.TrackedMod { return lookup[game].Mods }

func EnableMod(game config.Game, tm *model.TrackedMod, tis []*mods.ToInstall) (err error) {
	if err = downloads.Download(game, tm, tis, enableMod); err != nil {
		wu.ShowErrorLong(err)
	}
	return
}

func enableMod(game config.Game, tm *model.TrackedMod, tis []*mods.ToInstall, err error) {
	var (
		modPath   = filepath.Join(config.Get().GetModsFullPath(game), tm.GetDirSuffix())
		installed []*model.InstalledDownload
	)

	if err != nil {
		wu.ShowErrorLong(err)
		tm.Enabled = false
		return
	}

	for _, ti := range tis {
		if err = decompress(ti.Download.DownloadedLoc, modPath); err != nil {
			wu.ShowErrorLong(err)
			tm.Enabled = false
			return
		}
	}

	for _, ti := range tis {
		if err = AddModFiles(game, tm, ti.DownloadFiles); err != nil {
			wu.ShowErrorLong(err)
			tm.Enabled = false
			return
		}
	}

	tm.Installed = installed
	_ = saveToJson()
	return
}

func decompress(from string, to string) error {
	if filepath.Ext(from) == ".rar" {
		handler := func(ctx context.Context, f archiver.File) (err error) {
			if !f.IsDir() {
				var r io.ReadCloser
				if r, err = f.Open(); err != nil {
					return
				}
				defer func() { _ = r.Close() }()

				fp := filepath.Join(to, f.Name())
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

func DisableMod(game config.Game, tm *model.TrackedMod) (err error) {
	if err = RemoveModFiles(game, tm); err != nil {
		return
	}
	return saveToJson()
}

func RemoveMod(game config.Game, tm *model.TrackedMod) error {
	gm := lookup[game].Mods
	for i, m := range gm {
		if m.Mod.ID != tm.GetModID() {
			return fmt.Errorf("failed to find %s", tm.Mod.Name)
		}
		if m.Enabled {
			if err := RemoveModFiles(game, tm); err != nil {

			}
		}
		lookup[game].Mods = append(gm[:i], gm[i+1:]...)
		break
	}
	return saveToJson()
}

func saveToJson() error {
	return util.SaveToFile(filepath.Join(config.PWD, modTrackerName), &lookup)
}

func saveMoogle(tm *model.TrackedMod) (err error) {
	return util.SaveToFile(tm.MoogleModFile, tm.Mod)
}
