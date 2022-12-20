package mod_author

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/discover/remote"
	"github.com/kiamev/moogle-mod-manager/discover/repo"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/authored"
	config_installer "github.com/kiamev/moogle-mod-manager/ui/config-installer"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	u "github.com/kiamev/moogle-mod-manager/util"
	"github.com/ncruces/zenity"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func New() state.Screen {
	var kind = mods.Hosted
	a := &ModAuthorer{
		kind:           &kind,
		entryManager:   newEntryManager(),
		previewDef:     newPreviewDef(),
		donationsDef:   newDonationsDef(),
		gamesDef:       newGamesDef(),
		description:    newRichTextEditor(),
		releaseNotes:   newRichTextEditor(),
		categorySelect: widget.NewSelect(mods.Categories, func(string) {}),
	}
	//a.modKindDef = newModKindDef(a.kind)
	a.modCompatsDef = newModCompatibilityDef(a.gamesDef)
	a.subKindSelect = widget.NewSelect(mods.SubKinds, func(string) {
		a.downloads.UpdateTab()
		a.tabs.Refresh()
	})
	a.downloads = newDownloads(a.kind, a.subKindSelect)
	a.alwaysDownload = newAlwaysDownloadDef(a.downloads)
	a.configsDef = newConfigurationsDef(a.downloads)
	return a
}

type ModAuthorer struct {
	*entryManager
	modID        mods.ModID
	kind         *mods.Kind
	editCallback func(*mods.Mod)

	previewDef *previewDef
	//modKindDef     *modKindDef
	modCompatsDef  *modCompatabilityDef
	donationsDef   *donationsDef
	gamesDef       *gamesDef
	alwaysDownload *alwaysDownloadDef
	configsDef     *configurationsDef

	description    *richTextEditor
	releaseNotes   *richTextEditor
	categorySelect *widget.Select
	subKindSelect  *widget.Select

	downloads *downloads

	tabs *container.AppTabs
}

func (a *ModAuthorer) PreDraw(fyne.Window, ...interface{}) error { return nil }

func (a *ModAuthorer) DrawAsDialog(fyne.Window) {}

func (a *ModAuthorer) OnClose() {
	if a.editCallback != nil {
		a.editCallback = nil
	}
}

func (a *ModAuthorer) NewHostedMod() {
	a.modID = ""
	a.updateEntries(mods.NewMod(&mods.ModDef{
		ModKind: mods.ModKind{
			Kind: mods.Hosted,
		},
		ReleaseDate:         time.Now().Format("Jan 02 2006"),
		ConfigSelectionType: mods.Auto,
	}))
}

func (a *ModAuthorer) NewNexusMod() {
	a.modID = ""
	a.updateEntries(mods.NewMod(&mods.ModDef{
		ModKind: mods.ModKind{
			Kind: mods.Nexus,
		},
		ReleaseDate:         time.Now().Format("Jan 02 2006"),
		ConfigSelectionType: mods.Auto,
	}))

	e := widget.NewEntry()
	d := dialog.NewForm("", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("Link", e)},
		func(ok bool) {
			if !ok {
				state.ShowPreviousScreen()
				return
			}
			_, m, err := remote.NewNexusClient().GetFromUrl(e.Text)
			if err != nil {
				util.ShowErrorLong(err)
				return
			}
			a.updateEntries(m)
		}, ui.Window)
	d.Resize(fyne.NewSize(400, 200))
	d.Show()
}

func (a *ModAuthorer) NewCurseForgeMod() {
	a.modID = ""
	a.updateEntries(mods.NewMod(&mods.ModDef{
		ModKind: mods.ModKind{
			Kind: mods.CurseForge,
		},
		ReleaseDate:         time.Now().Format("Jan 02 2006"),
		ConfigSelectionType: mods.Auto,
	}))

	e := widget.NewEntry()
	d := dialog.NewForm("", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("Link", e)},
		func(ok bool) {
			if !ok {
				state.ShowPreviousScreen()
				return
			}
			_, m, err := remote.NewCurseForgeClient().GetFromUrl(e.Text)
			if err != nil {
				util.ShowErrorLong(err)
				return
			}
			a.updateEntries(m)
		}, ui.Window)
	d.Resize(fyne.NewSize(400, 200))
	d.Show()
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
	if b, err = os.ReadFile(file); err != nil {
		util.ShowErrorLong(err)
		return false
	}
	if path.Ext(file) == ".xml" {
		err = xml.Unmarshal(b, &mod)
	} else {
		err = json.Unmarshal(b, &mod)
	}
	if err != nil {
		util.ShowErrorLong(err)
		return false
	}
	*a.kind = mod.ModKind.Kind
	a.modID = mod.ModID
	a.updateEntries(&mod)
	return true
}

func (a *ModAuthorer) EditMod(mod *mods.Mod, editCallback func(*mods.Mod)) {
	a.modID = mod.ModID
	*a.kind = mod.ModKind.Kind
	a.editCallback = editCallback
	a.updateEntries(mod)
}

func (a *ModAuthorer) Draw(w fyne.Window) {
	if dir, found := authored.GetDir(a.modID); found {
		state.SetBaseDir(dir)
	}

	switch *a.kind {
	case mods.Hosted:
		a.tabs = a.createHostedInputs()
	case mods.CurseForge, mods.Nexus:
		a.tabs = a.createRemoteInputs()
	default:
		panic("invalid mod kind")
	}

	/*a.tabs.OnSelected = func(tab *container.TabItem) {
		if len(a.configsDef.list.Items) > 0 {
			a.getFormItem("Select Type").Widget.(*widget.Select).Enable()
		} else {
			a.getFormItem("Select Type").Widget.(*widget.Select).Disable()
		}
		tab.Content.Refresh()
	}*/

	smi := make([]*fyne.MenuItem, 0, 4)
	smi = append(smi,
		fyne.NewMenuItem("as json (local save)", func() {
			a.saveFile(asJson)
		}),
		/*fyne.NewMenuItem("as xml", func() {
			a.saveFile(asXml)
		})*/
		fyne.NewMenuItem("submit for review (public release)", func() {
			a.submitForReview()
		}),
	)
	if a.editCallback != nil {
		smi = append(smi, fyne.NewMenuItem("modify and back (local save)", func() {
			mod, err := a.compileMod()
			if err != nil {
				util.ShowErrorLong(err)
				return
			}
			callback := func() {
				m, e := a.compileMod()
				if e != nil {
					util.ShowErrorLong(e)
					return
				}
				a.editCallback(m)
				state.ShowPreviousScreen()
			}
			if !a.validate(mod, false) {
				dialog.ShowConfirm("Continue?", "The mod is not valid, continue anyway?", func(ok bool) {
					if ok {
						callback()
					}
				}, ui.Window)
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
			m, err := a.compileMod()
			if err != nil {
				util.ShowErrorLong(err)
				return
			}
			_ = a.validate(m, true)
		}),
		widget.NewButton("Test", func() {
			var (
				tis      []*mods.ToInstall
				mod, err = a.compileMod()
			)
			if err != nil {
				util.ShowErrorLong(err)
				return
			}

			if tis, err = mods.NewToInstallForMod(mod.ModKind.Kind, mod, mod.AlwaysDownload); err != nil {
				util.ShowErrorLong(err)
				return
			}

			if len(a.configsDef.list.Items) == 0 {
				util.DisplayDownloadsAndFiles(tis)
			} else {
				if err = state.GetScreen(state.ConfigInstaller).(config_installer.ConfigInstaller).Setup(mod, state.GetBaseDir(), func(_ mods.Result, tis []*mods.ToInstall) error {
					util.DisplayDownloadsAndFiles(tis)
					return nil
				}); err != nil {
					util.ShowErrorLong(err)
					return
				}
				state.ShowScreen(state.ConfigInstaller)
			}
		}),
		widget.NewSeparator(),
		cw.NewButtonWithPopups("Manual Edit",
			fyne.NewMenuItem("copy as json", func() {
				a.writeToClipboard(asJson)
			}), fyne.NewMenuItem("paste as json", func() {
				a.readFromClipboard(asJson)
			})),
		widget.NewSeparator(),
		cw.NewButtonWithPopups("Save", smi...))

	w.SetContent(container.NewBorder(nil, buttons, nil, nil, a.tabs))
}

func (a *ModAuthorer) updateEntries(mod *mods.Mod) {
	*a.kind = mod.ModKind.Kind
	a.createBaseDir(state.GetBaseDirBinding())
	a.modID = mod.ModID
	a.createFormItem("Name", string(mod.Name))
	a.createFormItem("Author", mod.Author)
	a.createFormItem("Release Date", mod.ReleaseDate)
	a.categorySelect.Selected = string(mod.Category)
	a.categorySelect.Refresh()
	a.subKindSelect.Selected = string(mod.SubKind())
	a.subKindSelect.Refresh()
	a.createFormItem("Version", mod.Version)
	a.description.SetText(mod.Description)
	a.releaseNotes.SetText(mod.ReleaseNotes)
	a.createFormItem("Link", mod.Link)
	//a.createFormSelect("Select Type", mods.SelectTypes, string(mod.ConfigSelectionType))

	a.createFormItem("Working Dir", config.PWD)
	if dir, ok := authored.GetDir(mod.ModID); ok && dir != "" {
		a.createFormItem("Working Dir", dir)
	}

	a.previewDef.set(mod.Preview)

	a.modCompatsDef.set(mod.ModCompatibility)
	//a.modKindDef.set(&mod.ModKind)
	a.downloads.set(mod)
	a.donationsDef.set(mod.DonationLinks)
	a.gamesDef.set(mod.Games)
	a.alwaysDownload.set(mod.AlwaysDownload)
	a.configsDef.set(mod.Configurations)

	a.downloads.set(mod)
}

type As byte

const (
	asJson As = iota
	asXml
)

func (a *ModAuthorer) writeToClipboard(as As) {
	var (
		b        []byte
		mod, err = a.compileMod()
	)
	if err != nil {
		util.ShowErrorLong(err)
		return
	}
	callback := func() {
		if b, err = a.Marshal(mod, asJson); err != nil {
			util.ShowErrorLong(err)
			return
		}
		_ = clipboard.WriteAll(string(b))
	}
	if !a.validate(mod, false) {
		dialog.ShowConfirm("Continue?", "The mod is not valid, continue anyway?", func(ok bool) {
			if ok {
				callback()
			}
		}, ui.Window)
	} else {
		callback()
	}
}

func (a *ModAuthorer) readFromClipboard(as As) {
	var (
		mod    mods.Mod
		s, err = clipboard.ReadAll()
	)
	if err != nil {
		util.ShowErrorLong(err)
		return
	}
	if as == asJson {
		err = json.Unmarshal([]byte(s), &mod)
	} else {
		err = xml.Unmarshal([]byte(s), &mod)
	}
	if err != nil {
		util.ShowErrorLong(err)
		return
	}
	a.updateEntries(&mod)
}

func (a *ModAuthorer) Marshal(mod *mods.Mod, as As) (b []byte, err error) {
	if as == asJson {
		b, err = json.MarshalIndent(mod, "", "\t")
	} else {
		b, err = xml.MarshalIndent(mod, "", "\t")
	}
	return
}

func (a *ModAuthorer) compileMod() (m *mods.Mod, err error) {
	m = mods.NewMod(&mods.ModDef{
		Name:         mods.ModName(a.getString("Name")),
		Author:       a.getString("Author"),
		ReleaseDate:  a.getString("Release Date"),
		Category:     mods.Category(a.categorySelect.Selected),
		Version:      a.getString("Version"),
		Description:  a.description.String(),
		ReleaseNotes: a.releaseNotes.String(),
		Link:         a.getString("Link"),
		ModKind: mods.ModKind{
			Kind: *a.kind,
		},
		Preview: a.previewDef.compile(),
		//ModKind:      *a.modKindDef.compile(),
		//ConfigSelectionType: mods.SelectType(a.getString("Select Type")),
		ConfigSelectionType: mods.Auto,
		ModCompatibility:    a.modCompatsDef.compile(),
		DonationLinks:       a.donationsDef.compile(),
		Games:               a.gamesDef.compile(),
		IsManuallyCreated:   true,
	})

	if err = a.downloads.compile(m); err != nil {
		return
	}

	switch *a.kind {
	case mods.Hosted:
		name := u.CreateFileName(string(m.Name))
		author := u.CreateFileName(m.Author)
		if name != "" && author != "" {
			m.ModID = mods.ModID(strings.ToLower(fmt.Sprintf("%s.%s", name, author)))
		}
	case mods.Nexus:
		m.ModID = mods.NewModID(mods.Nexus, string(a.modID))
	case mods.CurseForge:
		m.ModID = mods.NewModID(mods.CurseForge, string(a.modID))
	default:
		panic("invalid mod kind")
	}

	m.AlwaysDownload = a.alwaysDownload.compile()
	if len(m.AlwaysDownload) == 0 {
		m.AlwaysDownload = nil
	}
	for _, ad := range m.AlwaysDownload {
		for _, f := range ad.Files {
			f.From = trimNewLine(f.From)
			f.To = trimNewLine(f.To)
		}
		for _, d := range ad.Dirs {
			d.From = trimNewLine(d.From)
			d.To = trimNewLine(d.To)
		}
	}

	m.Configurations = a.configsDef.compile()
	if len(m.Configurations) == 0 {
		m.Configurations = nil
	}
	for _, conf := range m.Configurations {
		for _, c := range conf.Choices {
			if c.DownloadFiles != nil {
				for _, f := range c.DownloadFiles.Files {
					f.From = trimNewLine(f.From)
					f.To = trimNewLine(f.To)
				}
				for _, d := range c.DownloadFiles.Dirs {
					d.From = trimNewLine(d.From)
					d.To = trimNewLine(d.To)
				}
			}
		}
	}

	_ = authored.SetDir(m.ModID, state.GetBaseDir())
	return
}

func trimNewLine(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return strings.TrimSpace(s)
}

func (a *ModAuthorer) submitForReview() {
	var (
		mod, err = a.compileMod()
		pr       string
	)
	if err != nil {
		util.ShowErrorLong(err)
		return
	}

	if !a.validate(mod, false) {
		dialog.ShowInformation("Invalid Mod Def", "The mod is not valid, please fix it first.", ui.Window)
	}

	/*if err = repo.NewGetter(repo.Author).Pull(); err != nil {
		util.ShowErrorLong(err)
		return
	}*/

	if pr, err = repo.NewCommitter(mod).Submit(); err != nil {
		util.ShowErrorLong(err)
		return
	}

	prUrl, _ := url.Parse(pr)
	dialog.ShowCustom(
		"Successfully submitted mod",
		"ok",
		container.NewMax(widget.NewHyperlink(pr, prUrl)), ui.Window)
}

func (a *ModAuthorer) saveFile(asJson As) error {
	mod, err := a.compileMod()
	if err != nil {
		return err
	}
	if !a.validate(mod, false) {
		dialog.ShowConfirm("Continue?", "The mod is not valid, continue anyway?", func(ok bool) {
			if ok {
				a.save(mod, asJson)
			}
		}, ui.Window)
	} else {
		a.save(mod, asJson)
	}
	return nil
}

func (a *ModAuthorer) save(mod *mods.Mod, json As) {
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

	if json == asJson {
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
	if !strings.Contains(file, ext) {
		file = file + ext
	}
	if _, err = os.Stat(file); err == nil {
		dialog.ShowConfirm("Replace File?", "Replace "+file+"?", func(b bool) {
			save = b
		}, ui.Window)
	}
	if save {
		if err = os.WriteFile(file, b, 0755); err != nil {
			util.ShowErrorLong(err)
		}
	}
}

func (a *ModAuthorer) validate(mod *mods.Mod, showMessage bool) bool {
	s := mod.Validate()
	if showMessage {
		if s != "" {
			dialog.ShowError(errors.New(s), ui.Window)
		} else {
			dialog.ShowInformation("", "Mod is valid", ui.Window)
		}
	}
	return s == ""
}

func (a *ModAuthorer) createHostedInputs() *container.AppTabs {
	var entries = []*widget.FormItem{
		a.getBaseDirFormItem("Working Dir"),
		a.getFormItem("Name"),
		a.getFormItem("Author"),
	}
	if a.kind.Is(mods.Hosted) {
		entries = append(entries, widget.NewFormItem("Kind", a.subKindSelect))
	}
	entries = append(entries,
		widget.NewFormItem("Category", a.categorySelect),
		a.getFormItem("Version"),
		a.getFormItem("Release Date"),
		a.getFormItem("Link"))
	entries = append(entries, a.previewDef.getFormItems()...)

	return container.NewAppTabs(
		container.NewTabItem("Mod", container.NewVScroll(widget.NewForm(entries...))),
		container.NewTabItem("Description", a.description.Draw()),
		//container.NewTabItem("Kind", container.NewVScroll(a.modKindDef.draw())),
		container.NewTabItem("Compatibility", container.NewVScroll(a.modCompatsDef.draw())),
		container.NewTabItem("Release Notes", a.releaseNotes.Draw()),
		a.downloads.TabItem,
		container.NewTabItem("Donation Links", container.NewVScroll(a.donationsDef.draw())),
		container.NewTabItem("Games", container.NewVScroll(a.gamesDef.draw())),
		container.NewTabItem("Always Install", container.NewVScroll(a.alwaysDownload.draw())),
		container.NewTabItem("Configurations", container.NewVScroll(a.configsDef.draw())))
}

func (a *ModAuthorer) createRemoteInputs() *container.AppTabs {
	var entries = []*widget.FormItem{
		a.getBaseDirFormItem("Working Dir"),
		a.getFormItem("Name"),
		widget.NewFormItem("Category", a.categorySelect),
		//a.getFormItem("Select Type"),
	}
	entries = append(entries, a.previewDef.getFormItems()...)

	return container.NewAppTabs(
		container.NewTabItem("Mod", container.NewVScroll(widget.NewForm(entries...))),
		container.NewTabItem("Description", a.description.Draw()),
		container.NewTabItem("Compatibility", container.NewVScroll(a.modCompatsDef.draw())),
		//container.NewTabItem("Release Notes", a.releaseNotes.Draw()),
		//container.NewTabItem("Downloadables", container.NewVScroll(a.downloadDef.draw())),
		container.NewTabItem("Donation Links", container.NewVScroll(a.donationsDef.draw())),
		//container.NewTabItem("Games", container.NewVScroll(a.gamesDef.draw())),
		container.NewTabItem("Always Install", container.NewVScroll(a.alwaysDownload.draw())),
		container.NewTabItem("Configurations", container.NewVScroll(a.configsDef.draw())))
}
