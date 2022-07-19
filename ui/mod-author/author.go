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
	return &ModAuthorer{
		modCompatsDef:    newModCompatibilityDef(),
		downloadablesDef: newDownloadablesDef(),
		gamesDef:         newGamesDef(),
	}
}

type ModAuthorer struct {
	modCompatsDef    *modCompatsDef
	downloadablesDef *downloadablesDef
	gamesDef         *gamesDef
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
	form := widget.NewForm(
		getFormItem("ID"),
		getFormItem("Name"),
		getFormItem("Author"),
		getFormItem("Version"),
		getFormItem("ReleaseDate"),
		getFormItem("Category"),
		getFormItem("Description"),
		getFormItem("ReleaseNotes"),
		getFormItem("Link"),
		getFormItem("Preview"),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Mod", form),
		container.NewTabItem("Compatibility", a.modCompatsDef.draw()),
		container.NewTabItem("Downloadables", a.downloadablesDef.draw()),
		container.NewTabItem("Donation Links", widget.NewForm()),
		container.NewTabItem("Game", a.gamesDef.draw()),
		container.NewTabItem("Download Files", widget.NewForm()),
		container.NewTabItem("Configurations", widget.NewForm()))

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
	w.SetContent(container.NewVScroll(container.NewVBox(tabs, xmlButtons, jsonButtons)))
}

func (a *ModAuthorer) updateEntries(mod *mods.Mod) {
	setFormItem("ID", mod.ID)
	setFormItem("Name", mod.Name)
	setFormItem("Author", mod.Author)
	setFormItem("Version", mod.Version)
	setFormItem("ReleaseDate", mod.ReleaseDate)
	setFormItem("Category", mod.Category)
	setFormItem("Description", mod.Description)
	setFormItem("ReleaseNotes", mod.ReleaseNotes)
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
	return &mods.Mod{
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
		Downloadables:    a.downloadablesDef.compile(),
		DonationLinks:    nil,
		Game:             a.gamesDef.compile(),
		DownloadFiles:    nil,
		Configurations:   nil,
	}
}
