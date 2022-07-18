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
	"strings"
	"time"
)

var (
	mod     *mods.Mod
	entries = make(map[string]fyne.CanvasObject)
	gameDef = newGamesDef()
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
		ModCompatibility: nil,
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

	tabs := container.NewAppTabs(
		container.NewTabItem("Mod", form),
		container.NewTabItem("Compatibility", createCompatibilities(w)),
		container.NewTabItem("Downloadables", widget.NewForm()),
		container.NewTabItem("Donation Links", widget.NewForm()),
		container.NewTabItem("Game", gameDef.draw(w)),
		container.NewTabItem("Download Files", widget.NewForm()),
		container.NewTabItem("Configurations", widget.NewForm()))

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
	w.SetContent(container.NewVScroll(container.NewVBox(tabs, xmlButtons, jsonButtons)))
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

func getFormString(key string) string {
	e, ok := entries[key]
	if !ok {
		return ""
	}
	switch t := e.(type) {
	case *widget.Entry:
		return t.Text
	case *widget.SelectEntry:
		return t.Text
	}
	return ""
}

func getFormStrings(key string) []string {
	s := getFormString(key)
	if s != "" {
		return strings.Split(s, ", ")
	}
	return nil
}

func getFormItem(name string, customKey ...string) *widget.FormItem {
	key := name
	if len(customKey) > 0 {
		key = customKey[0]
	}
	e, _ := entries[key]
	return widget.NewFormItem(name, e)
}

func setFormItem(key string, value string) {
	e, ok := entries[key]
	if !ok {
		e = widget.NewEntry()
		//e.Resize(fyne.Size{Width: 200, Height: e.MinSize().Height})
		entries[key] = e
	}
	e.(*widget.Entry).SetText(value)
}

func setFormSelect(key string, possible []string, value string) {
	e, ok := entries[key]
	if !ok {
		e = widget.NewSelectEntry(possible)
		//e.Resize(fyne.Size{Width: 200, Height: e.MinSize().Height})
		entries[key] = e
	}
	e.(*widget.SelectEntry).SetText(value)
}

func setFormMultiLine(key string, value string) {
	e, ok := entries[key]
	if !ok {
		e = widget.NewMultiLineEntry()
		//e.Resize(fyne.Size{Width: 200, Height: e.MinSize().Height})
		entries[key] = e
	}
	e.(*widget.Entry).SetText(value)
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
	if b, err = Marshal(asJson); err != nil {
		dialog.ShowError(err, w)
		return
	}
	clipboard.Write(clipboard.FmtText, b)
}

func Marshal(asJson bool) (b []byte, err error) {
	if asJson {
		b, err = json.MarshalIndent(mod, "", "\t")
	} else {
		b, err = xml.MarshalIndent(mod, "", "\t")
	}
	return
}
