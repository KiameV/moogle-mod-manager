package author

import (
	"encoding/json"
	"encoding/xml"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/pr-modsync/mods"
	"golang.design/x/clipboard"
	"time"
)

var (
	mod     *mods.Mod
	entries = make(map[string]*widget.Entry)
)

func DrawNewMod(w fyne.Window) {
	mod = &mods.Mod{
		ID:               "",
		Name:             "",
		Author:           "",
		Version:          "",
		ReleaseDate:      time.Now().Format("Jan 02 2006"),
		Category:         "",
		Description:      "",
		ReleaseNotes:     "",
		Link:             "",
		Preview:          "",
		ModCompatibility: mods.ModCompatibility{},
		Downloadables:    nil,
		DonationLinks:    nil,
		Game:             nil,
		DownloadFiles:    nil,
		Configurations:   nil,
	}
	updateEntries()
	draw(w)
}

func DrawEditMod(w fyne.Window, m *mods.Mod) {
	mod = m
	updateEntries()
	draw(w)
}

func draw(w fyne.Window) {
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
	xmlButtons := container.NewHBox(
		widget.NewButton("Copy XML", func() {
			pasteToClipboard(w, false)
		}),
		widget.NewButton("Save mod.xml", func() {

		}))
	jsonButtons := container.NewHBox(
		widget.NewButton("Copy Json", func() {
			pasteToClipboard(w, true)
		}),
		widget.NewButton("Save mod.json", func() {

		}))
	w.SetContent(container.NewVScroll(container.NewVBox(form, xmlButtons, jsonButtons)))
}

func updateEntries() {
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

func getFormItem(name string) *widget.FormItem {
	e, _ := entries[name]
	return widget.NewFormItem(name, e)
}

func setFormItem(name string, value string) {
	e, ok := entries[name]
	if !ok {
		e = widget.NewEntry()
		entries[name] = e
	}
	e.SetText(value)
}

func pasteToClipboard(w fyne.Window, asJson bool) {
	var (
		b   []byte
		err error
	)
	if err = clipboard.Init(); err != nil {
		dialog.ShowError(err, w)
		return
	}
	if b, err = Marshal(w, asJson); err != nil {
		dialog.ShowError(err, w)
		return
	}
	clipboard.Write(clipboard.FmtText, b)
}

func Marshal(w fyne.Window, asJson bool) (b []byte, err error) {
	if asJson {
		b, err = json.MarshalIndent(mod, "", "\t")
	} else {
		b, err = xml.MarshalIndent(mod, "", "\t")
	}
	return
}
