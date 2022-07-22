package mod_author

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	config_installer "github.com/kiamev/moogle-mod-manager/ui/config-installer"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"github.com/ncruces/zenity"
	"golang.design/x/clipboard"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type installType string

const (
	directInstall  installType = "Direct Install"
	configurations installType = "Configuration"
)

var possibleInstallTypes = []string{string(directInstall), string(configurations)}

func New() state.Screen {
	dl := newDownloadsDef()
	return &ModAuthorer{
		entryManager:  newEntryManager(),
		modCompatsDef: newModCompatibilityDef(),
		downloadDef:   dl,
		donationsDef:  newDonationsDef(),
		gamesDef:      newGamesDef(),
		dlFilesDef:    newDownloadFilesDef(dl),
		configsDef:    newConfigurationsDef(dl),
	}
}

type ModAuthorer struct {
	*entryManager

	modCompatsDef *modCompatabilityDef
	downloadDef   *downloadsDef
	donationsDef  *donationsDef
	gamesDef      *gamesDef
	dlFilesDef    *downloadFilesDef
	configsDef    *configurationsDef

	tabs *container.AppTabs
}

func (a *ModAuthorer) NewMod() {
	a.updateEntries(&mods.Mod{
		ReleaseDate:         time.Now().Format("Jan 02 2006"),
		ConfigSelectionType: mods.Auto,
	})
}

func (a *ModAuthorer) EditMod() (successfullyLoadedMod bool) {
	var (
		file, err = zenity.SelectFile(
			zenity.Title("Load mod"),
			zenity.FileFilter{
				Name:     "mod",
				Patterns: []string{"*.xml", "*.json"},
			})
		b   []byte
		mod mods.Mod
	)
	if err != nil {
		return false
	}
	if b, err = ioutil.ReadFile(file); err != nil {
		dialog.ShowError(err, state.Window)
		return false
	}
	if filepath.Ext(file) == ".xml" {
		err = xml.Unmarshal(b, &mod)
	} else {
		err = json.Unmarshal(b, &mod)
	}
	a.updateEntries(&mod)
	return true
}

func (a *ModAuthorer) Draw(w fyne.Window) {
	form := container.NewVBox(
		widget.NewForm(
			a.getFormItem("ID"),
			a.getFormItem("Name"),
			a.getFormItem("Author"),
			a.getFormItem("Version"),
			a.getFormItem("Release Date"),
			a.getFormItem("Category"),
			a.getFormItem("Description"),
			a.getFormItem("Release Notes"),
			a.getFormItem("Link"),
			a.getFormItem("Mod File Links"),
			a.getFormItem("Preview"),
			a.getFormItem("Select Type")))

	a.tabs = container.NewAppTabs(
		container.NewTabItem("Mod", form),
		container.NewTabItem("Compatibility", a.modCompatsDef.draw()),
		container.NewTabItem("Downloadables", a.downloadDef.draw()),
		container.NewTabItem("Donation Links", a.donationsDef.draw()),
		container.NewTabItem("Games", a.gamesDef.draw()),
		container.NewTabItem("Always Install", a.dlFilesDef.draw()),
		container.NewTabItem("Configurations", a.configsDef.draw()))
	a.tabs.OnSelected = func(tab *container.TabItem) {
		if len(a.configsDef.list.Items) > 0 {
			a.getFormItem("Select Type").Widget.(*widget.Select).Enable()
		} else {
			a.getFormItem("Select Type").Widget.(*widget.Select).Disable()
		}
		tab.Content.Refresh()
	}

	validateButton := container.NewHBox(
		widget.NewButton("Validate", func() {
			a.validate()
		}),
		widget.NewButton("Test", func() {
			mod := a.compileMod()
			if len(a.configsDef.list.Items) == 0 {
				util.DisplayDownloadsAndFiles(mod, nil)
			}
			if err := state.GetScreen(state.ConfigInstaller).(config_installer.ConfigInstaller).Setup(mod, true); err != nil {
				dialog.ShowError(err, state.Window)
				return
			}
			state.ShowScreen(state.ConfigInstaller)
		}))

	xmlButtons := container.NewHBox(
		widget.NewButton("Copy XML", func() {
			a.pasteToClipboard(false)
		}),
		widget.NewButton("Save mod.xml", func() {
			a.saveFile(false)
		}))
	jsonButtons := container.NewHBox(
		widget.NewButton("Copy Json", func() {
			a.pasteToClipboard(true)
		}),
		widget.NewButton("Save mod.json", func() {
			a.saveFile(true)
		}))
	w.SetContent(container.NewVScroll(container.NewVBox(a.tabs, widget.NewSeparator(), validateButton, xmlButtons, jsonButtons)))
}

func (a *ModAuthorer) updateEntries(mod *mods.Mod) {
	a.createFormItem("ID", mod.ID)
	a.createFormItem("Name", mod.Name)
	a.createFormItem("Author", mod.Author)
	a.createFormItem("Version", mod.Version)
	a.createFormItem("Release Date", mod.ReleaseDate)
	a.createFormItem("Category", mod.Category)
	a.createFormItem("Description", mod.Description)
	a.createFormMultiLine("Release Notes", mod.ReleaseNotes)
	a.createFormItem("Link", mod.Link)
	a.createFormMultiLine("Mod File Links", strings.Join(mod.ModFileLinks, "\n"))
	a.createFormItem("Preview", mod.Preview)
	a.createFormSelect("Select Type", mods.SelectTypes, string(mod.ConfigSelectionType))

	a.modCompatsDef.set(mod.ModCompatibility)
	a.downloadDef.set(mod.Downloadables)
	a.donationsDef.set(mod.DonationLinks)
	a.gamesDef.set(mod.Games)
	a.dlFilesDef.set(mod.DownloadFiles)
	a.configsDef.set(mod.Configurations)
}

func (a *ModAuthorer) pasteToClipboard(asJson bool) {
	var (
		b   []byte
		err error
	)
	if err = clipboard.Init(); err != nil {
		dialog.ShowError(err, state.Window)
		return
	}
	if b, err = a.Marshal(asJson); err != nil {
		dialog.ShowError(err, state.Window)
		return
	}
	clipboard.Write(clipboard.FmtText, b)
}

func (a *ModAuthorer) Marshal(asJson bool) (b []byte, err error) {
	if asJson {
		b, err = json.MarshalIndent(a.compileMod(), "", "\t")
	} else {
		b, err = xml.MarshalIndent(a.compileMod(), "", "\t")
	}
	return
}

func (a *ModAuthorer) compileMod() (mod *mods.Mod) {
	m := &mods.Mod{
		ID:                  a.getString("ID"),
		Name:                a.getString("Name"),
		Author:              a.getString("Author"),
		Version:             a.getString("Version"),
		ReleaseDate:         a.getString("Release Date"),
		Category:            a.getString("Category"),
		Description:         a.getString("Description"),
		ReleaseNotes:        a.getString("Release Notes"),
		Link:                a.getString("Link"),
		ModFileLinks:        strings.Split(a.getString("Mod File Links"), "\n"),
		Preview:             a.getString("Preview"),
		ConfigSelectionType: mods.SelectType(a.getString("Select Type")),
		ModCompatibility:    a.modCompatsDef.compile(),
		Downloadables:       a.downloadDef.compile(),
		DonationLinks:       a.donationsDef.compile(),
		Games:               a.gamesDef.compile(),
	}
	m.DownloadFiles = a.dlFilesDef.compile()
	if m.DownloadFiles == nil ||
		(m.Name == "" && len(m.DownloadFiles.Files) == 0 && (len(m.DonationLinks) == 0)) {
		m.DownloadFiles = nil
	}
	m.Configurations = a.configsDef.compile()
	if len(m.Configurations) == 0 {
		m.Configurations = nil
	}
	return m
}

func (a *ModAuthorer) saveFile(asJson bool) {
	var (
		b, err = a.Marshal(asJson)
		ext    string
		file   string
		save   = true
	)
	if err != nil {
		dialog.ShowError(err, state.Window)
		return
	}

	if asJson {
		ext = ".json"
	} else {
		ext = ".xml"
	}

	if file, err = zenity.SelectFileSave(
		zenity.Title("Save mod"+ext),
		zenity.Filename("mod"+ext),
		zenity.FileFilter{
			Name:     "*" + ext,
			Patterns: []string{"*" + ext},
		}); err != nil {
		return
	}
	if strings.Index(file, ext) == -1 {
		file = file + ext
	}
	if _, err = os.Stat(file); err == nil {
		dialog.ShowConfirm("Replace File?", "Replace "+file+"?", func(b bool) {
			save = b
		}, state.Window)
	}
	if save {
		if err = ioutil.WriteFile(file, b, 0755); err != nil {
			dialog.ShowError(err, state.Window)
		}
	}
}

func (a *ModAuthorer) validate() {
	s := a.compileMod().Validate()
	if s != "" {
		dialog.ShowError(errors.New(s), state.Window)
	} else {
		dialog.ShowInformation("", "Mod is valid", state.Window)
	}
}
