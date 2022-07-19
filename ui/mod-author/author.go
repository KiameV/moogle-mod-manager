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
		modCompatsDef: newModCompatibilityDef(),
		downloadDef:   dl,
		donationsDef:  newDonationsDef(),
		gamesDef:      newGamesDef(),
		dlFilesDef:    newDownloadFilesDef(dl),
	}
}

type ModAuthorer struct {
	modCompatsDef *modCompatsDef
	downloadDef   *downloadsDef
	donationsDef  *donationsDef
	gamesDef      *gamesDef
	dlFilesDef    *downloadFilesDef

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
	installType := widget.NewSelect([]string{"Direct Install", "Configuration"}, func(choice string) {
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
			getFormItem("ID"),
			getFormItem("Name"),
			getFormItem("Author"),
			getFormItem("Version"),
			getFormItem("Release Date"),
			getFormItem("Category"),
			getFormItem("Description"),
			getFormItem("Release Notes"),
			getFormItem("Link"),
			getFormItem("Preview"),
			widget.NewFormItem("Install Type", installType)))

	a.tabs = container.NewAppTabs(
		container.NewTabItem("Mod", form),
		container.NewTabItem("Compatibility", a.modCompatsDef.draw()),
		container.NewTabItem("Downloadables", a.downloadDef.draw()),
		container.NewTabItem("Donation Links", a.donationsDef.draw()),
		container.NewTabItem("Game", a.gamesDef.draw()))
	a.dlTab = container.NewTabItem("Download Files", a.dlFilesDef.draw())
	a.configTab = container.NewTabItem("Configurations", widget.NewForm())
	a.tabs.OnSelected = func(tab *container.TabItem) {
		if tab == a.dlTab {
			tab.Content = a.dlFilesDef.draw()
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
	setFormItem("ID", mod.ID)
	setFormItem("Name", mod.Name)
	setFormItem("Author", mod.Author)
	setFormItem("Version", mod.Version)
	setFormItem("Release Date", mod.ReleaseDate)
	setFormItem("Category", mod.Category)
	setFormItem("Description", mod.Description)
	setFormMultiLine("Release Notes", mod.ReleaseNotes)
	setFormItem("Link", mod.Link)
	setFormItem("Preview", mod.Preview)
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
		ID:               getFormString("ID"),
		Name:             getFormString("Name"),
		Author:           getFormString("Author"),
		Version:          getFormString("Version"),
		ReleaseDate:      getFormString("ReleaseDate"),
		Category:         getFormString("Category"),
		Description:      getFormString("Description"),
		ReleaseNotes:     getFormString("ReleaseNotes"),
		Link:             getFormString("Link"),
		Preview:          getFormString("Preview"),
		ModCompatibility: a.modCompatsDef.compile(),
		Downloadables:    a.downloadDef.compile(),
		DonationLinks:    a.donationsDef.compile(),
		Game:             a.gamesDef.compile(),
	}

	return m
}
