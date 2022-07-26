package managed

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/browser"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/decompressor"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"io/ioutil"
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
			if b, err = readFile(tm.MoogleModFile); err != nil {
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

func AddModFromFile(game config.Game, file string) (tm *model.TrackedMod, err error) {
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
		return nil, fmt.Errorf("unknown file extension: %s", ext)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load mod: %v", err)
	}
	if s := mod.Validate(); s != "" {
		return nil, fmt.Errorf("failed to load mod:\n%s", s)
	}

	tm = model.NewTrackerMod(mod, game)
	err = AddMod(game, tm)
	return
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
	err = AddMod(game, tm)
	return
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
	}

	if err = saveMoogle(game, tm); err != nil {
		return
	}

	for _, g := range tm.Mod.Games {
		i := int(config.NameToGame(g.Name))
		m := lookup[i]
		m.Mods = append(m.Mods, tm)
	}
	return saveToJson()
}

func UpdateMod(game config.Game, tm *model.TrackedMod) (err error) {
	if err = tm.GetMod().Supports(game); err != nil {
		return
	}

	if err = DisableMod(game, tm); err != nil {
		return
	}

	tm.Mod = tm.UpdatedMod
	if err = saveMoogle(game, tm); err != nil {
		return
	}

	tm.UpdatedMod = nil
	return saveToJson()
}

func GetMods(game config.Game) []*model.TrackedMod { return lookup[game].Mods }

func GetMod(game config.Game, modID string) (*model.TrackedMod, bool) {
	if mods := GetMods(game); mods != nil {
		for _, tm := range mods {
			if tm.Mod.ID == modID {
				return tm, true
			}
		}
	}
	return nil, false
}

func EnableMod(game config.Game, tm *model.TrackedMod, tis []*mods.ToInstall) (err error) {
	confirmDownloads(tis, func() {
		var (
			downloadDir string
			f           string
			d           decompressor.Decompressor
			modPath     = filepath.Join(config.Get().GetModsFullPath(game), tm.GetDirSuffix())
		)
		if downloadDir, err = createPath(filepath.Join(config.Get().GetDownloadFullPath(game), tm.GetDirSuffix())); err != nil {
			dialog.ShowError(err, state.Window)
			return
		}
		for _, ti := range tis {
			if len(ti.Download.Sources) == 0 {
				dialog.ShowError(fmt.Errorf("%s has no download sources", ti.Download.Name), state.Window)
				return
			}
			for _, source := range ti.Download.Sources {
				if f, err = browser.Download(source, downloadDir); err == nil {
					// success
					ti.Download.DownloadedLoc = f
					break
				}
			}
		}

		for _, ti := range tis {
			if ti.Download.DownloadedLoc == "" {
				dialog.ShowError(fmt.Errorf("failed to download %s", ti.Download.Sources[0]), state.Window)
				return
			}
		}

		for _, ti := range tis {
			if d, err = decompressor.NewDecompressor(ti.Download.DownloadedLoc); err != nil {
				dialog.ShowError(err, state.Window)
				return
			}
			if err = d.DecompressTo(modPath); err != nil {
				dialog.ShowError(err, state.Window)
				return
			}
		}

		for _, ti := range tis {
			if err = AddModFiles(game, tm, ti.DownloadFiles); err != nil {
				dialog.ShowError(err, state.Window)
				return
			}
		}
	})
	return
}

func createPath(path string) (string, error) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		err = fmt.Errorf("failed to create mod directory: %v", err)
		return "", err
	}
	return path, nil
}

func DisableMod(game config.Game, tm *model.TrackedMod) (err error) {
	// TODO
	return RemoveModFiles(game, tm)
}

func confirmDownloads(tis []*mods.ToInstall, callback func()) {
	sb := strings.Builder{}
	for i, ti := range tis {
		sb.WriteString(fmt.Sprintf("## Download %d\n\n", i+1))
		if len(ti.Download.Sources) == 1 {
			sb.WriteString(ti.Download.Sources[0] + "\n\n")
		} else {
			sb.WriteString("### Sources:\n\n")
			for j, s := range ti.Download.Sources {
				sb.WriteString(fmt.Sprintf(" - %d. %s\n\n", j+1, s))
			}
		}
	}
	d := dialog.NewCustomConfirm("Download Files?", "Yes", "Cancel", container.NewVScroll(widget.NewRichTextFromMarkdown(sb.String())), func(ok bool) {
		if ok {
			callback()
		}
	}, state.Window)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
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

func saveMoogle(game config.Game, tm *model.TrackedMod) (err error) {
	var (
		b []byte
		f *os.File
	)
	if b, err = json.MarshalIndent(tm.Mod, "", "\t"); err != nil {
		return
	}

	modPath := filepath.Join(config.Get().GetModsFullPath(game), tm.GetDirSuffix())
	if _, err = os.Stat(modPath); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(modPath), 0777); err != nil {
			return
		}
	}
	if f, err = os.Create(tm.MoogleModFile); err != nil {
		return
	}
	defer func() { _ = f.Close() }()
	if _, err = f.Write(b); err != nil {
		return
	}
	return
}
