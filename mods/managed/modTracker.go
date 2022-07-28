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
	"github.com/kiamev/moogle-mod-manager/util"
	"io"
	"net/http"
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
	confirmDownloads(tm, tis, func() {
		var (
			downloadDir string
			f           string
			d           decompressor.Decompressor
			modPath     = filepath.Join(config.Get().GetModsFullPath(game), tm.GetDirSuffix())
			installed   []*model.InstalledDownload
		)
		if downloadDir, err = createPath(filepath.Join(config.Get().GetDownloadFullPath(game), tm.GetDirSuffix())); err != nil {
			dialog.ShowError(err, state.Window)
			return
		}

		if tm.Mod.ModKind.Kind == mods.Hosted {
			for _, ti := range tis {
				if len(ti.Download.Hosted.Sources) == 0 {
					dialog.ShowError(fmt.Errorf("%s has no download sources", ti.Download.Name), state.Window)
					return
				}
				for _, source := range ti.Download.Hosted.Sources {
					if f, err = browser.Download(source, filepath.Join(downloadDir, util.CreateFileName(ti.Download.Version))); err == nil {
						// success
						installed = append(installed, model.NewInstalledDownload(ti.Download.Name, ti.Download.Version))
						ti.Download.DownloadedLoc = f
						break
					}
				}
			}

			for _, ti := range tis {
				if ti.Download.DownloadedLoc == "" {
					dialog.ShowError(fmt.Errorf("failed to download %s", ti.Download.Hosted.Sources[0]), state.Window)
					return
				}
			}
		} else {
			for _, ti := range tis {
				if ti.Download.DownloadedLoc == "" {
					dialog.ShowError(fmt.Errorf("failed to download %s", ti.Download.Nexus.FileName), state.Window)
					return
				}
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

		tm.Installed = installed
		_ = saveToJson()
	})
	// Don't save here, save in the callback
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
	if err = RemoveModFiles(game, tm); err != nil {
		return
	}
	return saveToJson()
}

func confirmDownloads(tm *model.TrackedMod, tis []*mods.ToInstall, callback func()) {
	if tm.Mod.ModKind.Kind == mods.Nexus {
		for _, ti := range tis {
			resp, err := http.Get(ti.Download.DownloadedLoc)
			if err != nil {
				return
			}
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}
			b = b
		}
		callback()
		return
	}
	sb := strings.Builder{}
	for i, ti := range tis {
		sb.WriteString(fmt.Sprintf("## Download %d\n\n", i+1))
		if len(ti.Download.Hosted.Sources) == 1 {
			sb.WriteString(ti.Download.Hosted.Sources[0] + "\n\n")
		} else {
			sb.WriteString("### Sources:\n\n")
			for j, s := range ti.Download.Hosted.Sources {
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

func saveToJson() error {
	return util.SaveToFile(filepath.Join(config.PWD, modTrackerName), &lookup)
}

func saveMoogle(tm *model.TrackedMod) (err error) {
	return util.SaveToFile(tm.MoogleModFile, tm.Mod)
}
