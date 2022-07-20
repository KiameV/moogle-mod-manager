package mod_author

import (
	"encoding/json"
	"encoding/xml"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"golang.design/x/clipboard"
	"time"
)

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

	tabs       *container.AppTabs
	dlTab      *container.TabItem
	configTab  *container.TabItem
	installTab *container.TabItem
}

func (a *ModAuthorer) NewMod() {
	a.updateEntries(&mods.Mod{
		ReleaseDate: time.Now().Format("Jan 02 2006"),
	})
}

func (a *ModAuthorer) EditMod(m *mods.Mod) {
	a.updateEntries(m)
}

func (a *ModAuthorer) Draw(w fyne.Window) {
	installType := widget.NewRadioGroup([]string{"Direct Install", "Configuration"}, func(choice string) {
		title := a.tabs.Items[len(a.tabs.Items)-1].Text
		if title == "Download Files" || title == "Configurations" {
			a.tabs.RemoveIndex(len(a.tabs.Items) - 1)
		}
		a.installTab = a.configTab
		if choice == "Direct Install" {
			a.installTab = a.dlTab
		}
		a.tabs.Append(a.installTab)
	})

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
			a.getFormItem("Preview"),
			widget.NewFormItem("Install Type", installType)))

	a.tabs = container.NewAppTabs(
		container.NewTabItem("Mod", form),
		container.NewTabItem("Compatibility", a.modCompatsDef.draw()),
		container.NewTabItem("Downloadables", a.downloadDef.draw()),
		container.NewTabItem("Donation Links", a.donationsDef.draw()),
		container.NewTabItem("Game", a.gamesDef.draw()))
	a.dlTab = container.NewTabItem("Download Files", a.dlFilesDef.draw())
	a.configTab = container.NewTabItem("Configurations", a.configsDef.draw())
	a.tabs.OnSelected = func(tab *container.TabItem) {
		if tab == a.dlTab {
			tab.Content = a.dlFilesDef.draw()
		} else if tab == a.configTab {
			tab.Content = a.configsDef.draw()
		}
	}

	installType.SetSelected("Direct Install")

	xmlButtons := container.NewHBox(
		widget.NewButton("Copy XML", func() {
			a.pasteToClipboard(false)
		}),
		widget.NewButton("Save mod.xml", func() {

		}))
	jsonButtons := container.NewHBox(
		widget.NewButton("Copy Json", func() {
			a.pasteToClipboard(true)
		}),
		widget.NewButton("Save mod.json", func() {

		}))
	w.SetContent(container.NewVScroll(container.NewVBox(a.tabs, widget.NewSeparator(), xmlButtons, jsonButtons)))
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
	a.createFormItem("Preview", mod.Preview)
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
		ID:               a.getString("ID"),
		Name:             a.getString("Name"),
		Author:           a.getString("Author"),
		Version:          a.getString("Version"),
		ReleaseDate:      a.getString("ReleaseDate"),
		Category:         a.getString("Category"),
		Description:      a.getString("Description"),
		ReleaseNotes:     a.getString("ReleaseNotes"),
		Link:             a.getString("Link"),
		Preview:          a.getString("Preview"),
		ModCompatibility: a.modCompatsDef.compile(),
		Downloadables:    a.downloadDef.compile(),
		DonationLinks:    a.donationsDef.compile(),
		Game:             a.gamesDef.compile(),
	}
	if a.installTab == a.dlTab {
		m.DownloadFiles = a.dlFilesDef.compile()
		if m.DownloadFiles == nil ||
			(m.Name == "" && len(m.DownloadFiles.Files) == 0 && (len(m.DonationLinks) == 0)) {
			m.DownloadFiles = nil
		}
	} else {
		m.Configurations = a.configsDef.compile()
		if len(m.Configurations) == 0 {
			m.Configurations = nil
		}
	}
	return m
}
