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
	"github.com/kiamev/moogle-mod-manager/mods/managed/authored"
	"github.com/kiamev/moogle-mod-manager/mods/nexus"
	"github.com/kiamev/moogle-mod-manager/repo"
	config_installer "github.com/kiamev/moogle-mod-manager/ui/config-installer"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"github.com/ncruces/zenity"
	"golang.design/x/clipboard"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func New() state.Screen {
	dl := newDownloadsDef()
	return &ModAuthorer{
		entryManager:   newEntryManager(),
		previewDef:     newPreviewDef(),
		modKindDef:     newModKindDef(),
		modCompatsDef:  newModCompatibilityDef(),
		downloadDef:    dl,
		donationsDef:   newDonationsDef(),
		gamesDef:       newGamesDef(),
		alwaysDownload: newAlwaysDownloadDef(dl),
		configsDef:     newConfigurationsDef(dl),
		description:    newRichTextEditor(),
		releaseNotes:   newRichTextEditor(),
	}
}

type ModAuthorer struct {
	*entryManager
	modBeingEdited *mods.Mod

	previewDef     *previewDef
	modKindDef     *modKindDef
	modCompatsDef  *modCompatabilityDef
	downloadDef    *downloadsDef
	donationsDef   *donationsDef
	gamesDef       *gamesDef
	alwaysDownload *alwaysDownloadDef
	configsDef     *configurationsDef

	description  *richTextEditor
	releaseNotes *richTextEditor

	tabs *container.AppTabs
}

func (a *ModAuthorer) PreDraw() error { return nil }

func (a *ModAuthorer) OnClose() {
	if a.modBeingEdited != nil {
		a.modBeingEdited = nil
	}
}

func (a *ModAuthorer) NewMod() {
	a.updateEntries(&mods.Mod{
		ReleaseDate:         time.Now().Format("Jan 02 2006"),
		ConfigSelectionType: mods.Auto,
	})
}

func (a *ModAuthorer) NewFromNexus() {
	a.NewMod()
	e := widget.NewEntry()
	dialog.ShowForm("", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("Link", e)},
		func(ok bool) {
			if !ok {
				state.ShowPreviousScreen()
				return
			}
			m, err := nexus.GetModFromNexus(*state.CurrentGame, e.Text)
			if err != nil {
				util.ShowErrorLong(err)
				return
			}
			a.updateEntries(m)
		}, state.Window)
}

func (a *ModAuthorer) LoadModToEdit() (successfullyLoadedMod bool) {
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
		util.ShowErrorLong(err)
		return false
	}
	if path.Ext(file) == ".xml" {
		err = xml.Unmarshal(b, &mod)
	} else {
		err = json.Unmarshal(b, &mod)
	}
	a.updateEntries(&mod)
	return true
}

func (a *ModAuthorer) EditMod(mod *mods.Mod) {
	a.modBeingEdited = mod
	a.updateEntries(mod)
}

func (a *ModAuthorer) Draw(w fyne.Window) {
	if dir, found := authored.GetDir(a.getString("ID")); found {
		state.SetBaseDir(dir)
	}
	items := []*widget.FormItem{
		a.getBaseDirFormItem("Working Dir"),
		a.getFormItem("ID"),
		a.getFormItem("Name"),
		a.getFormItem("Author"),
		a.getFormItem("Category"),
		a.getFormItem("Version"),
		a.getFormItem("Release Date"),
		a.getFormItem("Link"),
		a.getFormItem("Select Type"),
	}
	items = append(items, a.previewDef.getFormItems()...)

	a.tabs = container.NewAppTabs(
		container.NewTabItem("Mod", container.NewVScroll(widget.NewForm(items...))),
		container.NewTabItem("Description", a.description.Draw()),
		container.NewTabItem("Kind", container.NewVScroll(a.modKindDef.draw())),
		container.NewTabItem("Compatibility", container.NewVScroll(a.modCompatsDef.draw())),
		container.NewTabItem("Release Notes", a.releaseNotes.Draw()),
		container.NewTabItem("Downloadables", container.NewVScroll(a.downloadDef.draw())),
		container.NewTabItem("Donation Links", container.NewVScroll(a.donationsDef.draw())),
		container.NewTabItem("Games", container.NewVScroll(a.gamesDef.draw())),
		container.NewTabItem("Always Install", container.NewVScroll(a.alwaysDownload.draw())),
		container.NewTabItem("Configurations", container.NewVScroll(a.configsDef.draw())))
	a.tabs.OnSelected = func(tab *container.TabItem) {
		if len(a.configsDef.list.Items) > 0 {
			a.getFormItem("Select Type").Widget.(*widget.Select).Enable()
		} else {
			a.getFormItem("Select Type").Widget.(*widget.Select).Disable()
		}
		tab.Content.Refresh()
	}

	smi := make([]*fyne.MenuItem, 0, 4)
	smi = append(smi,
		fyne.NewMenuItem("as json", func() {
			a.saveFile(true)
		}),
		/*fyne.NewMenuItem("as xml", func() {
			a.saveFile(false)
		})*/
		fyne.NewMenuItem("submit for review", func() {
			a.submitForReview()
		}),
	)
	if a.modBeingEdited != nil {
		smi = append(smi, fyne.NewMenuItem("modify and back", func() {
			mod := a.compileMod()
			callback := func() {
				*a.modBeingEdited = *a.compileMod()
				state.ShowPreviousScreen()
			}
			if !a.validate(mod, false) {
				dialog.ShowConfirm("Continue?", "The mod is not valid, continue anyway?", func(ok bool) {
					if ok {
						callback()
					}
				}, state.Window)
			} else {
				callback()
			}
		}))
	}

	buttons := container.NewHBox(
		widget.NewButton("Back", func() {
			dialog.ShowConfirm("Go Back?", "Any unsaved data will be lost, continue?", func(ok bool) {
				if ok {
					state.ShowPreviousScreen()
				}
			}, w)
		}),
		widget.NewButton("Validate", func() {
			_ = a.validate(a.compileMod(), true)
		}),
		widget.NewButton("Test", func() {
			mod := a.compileMod()
			if len(a.configsDef.list.Items) == 0 {
				tis, err := mods.NewToInstallForMod(mod.ModKind.Kind, mod, mod.AlwaysDownload)
				if err != nil {
					util.ShowErrorLong(err)
					return
				}
				util.DisplayDownloadsAndFiles(tis)
			}
			if err := state.GetScreen(state.ConfigInstaller).(config_installer.ConfigInstaller).Setup(mod, state.GetBaseDir(), func(tis []*mods.ToInstall) error {
				util.DisplayDownloadsAndFiles(tis)
				return nil
			}); err != nil {
				util.ShowErrorLong(err)
				return
			}
			state.ShowScreen(state.ConfigInstaller)
		}),
		widget.NewSeparator(),
		cw.NewButtonWithPopups("Copy",
			fyne.NewMenuItem("as json", func() {
				a.pasteToClipboard(true)
			}), fyne.NewMenuItem("as xml", func() {
				a.pasteToClipboard(false)
			})),
		cw.NewButtonWithPopups("Save", smi...))

	w.SetContent(container.NewBorder(nil, buttons, nil, nil, a.tabs))
}

func (a *ModAuthorer) updateEntries(mod *mods.Mod) {
	a.createBaseDir(state.GetBaseDirBinding())
	a.createFormItem("ID", mod.ID)
	a.createFormItem("Name", mod.Name)
	a.createFormItem("Author", mod.Author)
	a.createFormItem("Release Date", mod.ReleaseDate)
	a.createFormItem("Category", mod.Category)
	a.createFormItem("Version", mod.Version)
	a.description.SetText(mod.Description)
	a.releaseNotes.SetText(mod.ReleaseNotes)
	a.createFormItem("Link", mod.Link)
	a.createFormSelect("Select Type", mods.SelectTypes, string(mod.ConfigSelectionType))

	if dir, ok := authored.GetDir(mod.ID); ok && dir != "" {
		a.createFormItem("Working Dir", dir)
	}

	a.previewDef.set(mod.Preview)
	a.modCompatsDef.set(mod.ModCompatibility)
	a.modKindDef.set(mod.ModKind)
	a.downloadDef.set(mod.Downloadables)
	a.donationsDef.set(mod.DonationLinks)
	a.gamesDef.set(mod.Game)
	a.alwaysDownload.set(mod.AlwaysDownload)
	a.configsDef.set(mod.Configurations)
}

func (a *ModAuthorer) pasteToClipboard(asJson bool) {
	var (
		b   []byte
		err error
	)
	if err = clipboard.Init(); err != nil {
		util.ShowErrorLong(err)
		return
	}
	mod := a.compileMod()
	callback := func() {
		if b, err = a.Marshal(mod, asJson); err != nil {
			util.ShowErrorLong(err)
			return
		}
		clipboard.Write(clipboard.FmtText, b)
	}
	if !a.validate(mod, false) {
		dialog.ShowConfirm("Continue?", "The mod is not valid, continue anyway?", func(ok bool) {
			if ok {
				callback()
			}
		}, state.Window)
	} else {
		callback()
	}
}

func (a *ModAuthorer) Marshal(mod *mods.Mod, asJson bool) (b []byte, err error) {
	if asJson {
		b, err = json.MarshalIndent(mod, "", "\t")
	} else {
		b, err = xml.MarshalIndent(mod, "", "\t")
	}
	return
}

func (a *ModAuthorer) compileMod() (mod *mods.Mod) {
	m := &mods.Mod{
		ID:                  a.getString("ID"),
		Name:                a.getString("Name"),
		Author:              a.getString("Author"),
		ReleaseDate:         a.getString("Release Date"),
		Category:            a.getString("Category"),
		Version:             a.getString("Version"),
		Description:         a.description.String(),
		ReleaseNotes:        a.releaseNotes.String(),
		Link:                a.getString("Link"),
		Preview:             a.previewDef.compile(),
		ModKind:             a.modKindDef.compile(),
		ConfigSelectionType: mods.SelectType(a.getString("Select Type")),
		ModCompatibility:    a.modCompatsDef.compile(),
		Downloadables:       a.downloadDef.compile(),
		DonationLinks:       a.donationsDef.compile(),
		Game:                a.gamesDef.compile(),
		IsManuallyCreated:   true,
	}

	m.AlwaysDownload = a.alwaysDownload.compile()
	if len(m.AlwaysDownload) == 0 {
		m.AlwaysDownload = nil
	}

	m.Configurations = a.configsDef.compile()
	if len(m.Configurations) == 0 {
		m.Configurations = nil
	}
	authored.SetDir(m.ID, state.GetBaseDir())
	return m
}

func (a *ModAuthorer) submitForReview() {
	mod := a.compileMod()
	if !a.validate(mod, false) {
		dialog.ShowInformation("Invalid Mod Def", "The mod is not valid, please fix it first.", state.Window)
	}
	ur, err := repo.NewCommitter(mod).Submit()
	if err != nil {
		util.ShowErrorLong(err)
	} else {
		u, _ := url.Parse(ur)
		dialog.ShowCustom(
			"Successfully submitted mod",
			"ok",
			container.NewMax(widget.NewHyperlink(ur, u)), state.Window)
	}
}

func (a *ModAuthorer) saveFile(asJson bool) {
	mod := a.compileMod()
	if !a.validate(mod, false) {
		dialog.ShowConfirm("Continue?", "The mod is not valid, continue anyway?", func(ok bool) {
			if ok {
				a.save(mod, asJson)
			}
		}, state.Window)
	} else {
		a.save(mod, asJson)
	}
}

func (a *ModAuthorer) save(mod *mods.Mod, asJson bool) {
	var (
		b    []byte
		ext  string
		file string
		save = true
		err  error
	)
	if b, err = a.Marshal(mod, asJson); err != nil {
		util.ShowErrorLong(err)
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
			util.ShowErrorLong(err)
		}
	}
}

func (a *ModAuthorer) validate(mod *mods.Mod, showMessage bool) bool {
	s := mod.Validate()
	if showMessage {
		if s != "" {
			dialog.ShowError(errors.New(s), state.Window)
		} else {
			dialog.ShowInformation("", "Mod is valid", state.Window)
		}
	}
	return s == ""
}
