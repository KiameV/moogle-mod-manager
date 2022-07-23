package local

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/ncruces/zenity"
)

type LocalUI interface {
	state.Screen
	GetSelected() *model.TrackedMod
}

func New() LocalUI {
	return &localMods{}
}

type localMods struct {
	selectedMod *model.TrackedMod
}

func (m *localMods) OnClose() {

}

func (m *localMods) GetSelected() *model.TrackedMod {
	return m.selectedMod
}

func (m *localMods) Draw(w fyne.Window) {
	var (
		selectable = managed.GetMods(*state.CurrentGame)
		modList    = widget.NewList(
			func() int { return len(selectable) },
			func() fyne.CanvasObject {
				return container.NewHBox(widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}))
			},
			func(id widget.ListItemID, object fyne.CanvasObject) {
				object.(*fyne.Container).Objects[0].(*widget.Label).SetText(selectable[id].Mod.Name)
			})
		addButton = cw.NewButtonWithPopups("Add",
			fyne.NewMenuItem("From File", func() { m.addFromFile() }),
			fyne.NewMenuItem("From URL", func() { m.addFromUrl() }))
		removeButton = widget.NewButton("Remove", func() {})
		modDetails   = container.NewScroll(container.NewMax())
	)
	removeButton.Disable()
	modList.OnSelected = func(id widget.ListItemID) {
		m.selectedMod = selectable[id]
		removeButton.Enable()
		modDetails.Content = m.createPreview(m.selectedMod.Mod)
		modDetails.Refresh()
	}
	modList.OnUnselected = func(id widget.ListItemID) {
		m.selectedMod = nil
		removeButton.Disable()
		modDetails.Hide()
	}

	buttons := container.NewHBox(addButton, widget.NewSeparator(), removeButton)

	split := container.NewHSplit(
		container.NewVScroll(modList),
		modDetails)
	split.SetOffset(0.3)

	w.SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle(config.GameNameString(*state.CurrentGame), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			buttons,
		), nil, nil, nil,
		split))
}

func (m *localMods) createPreview(mod *mods.Mod) fyne.CanvasObject {
	c := container.NewVBox(
		m.createField("Name", mod.Name),
		m.createMultiLineField("Description", mod.Description),
		m.createField("Version", mod.Version),
		m.createField("Link", mod.Link),
		m.createField("Author", mod.Author),
		m.createField("Category", mod.Category),
		m.createField("Release Date", mod.ReleaseDate),
	)
	if mod.ReleaseNotes != "" {
		c.Add(m.createMultiLineField("Release Notes", mod.ReleaseDate))
	}
	if mod.ModCompatibility != nil && mod.ModCompatibility.HasItems() {
		c.Add(m.createCompatibility(mod.ModCompatibility))
	}

	if img := mod.Preview.Get(); img != nil {
		c = container.NewBorder(img, nil, nil, nil, c)
	}
	return c
}

func (m *localMods) createField(name, value string) *fyne.Container {
	return container.NewHBox(
		widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(value),
	)
}

func (m *localMods) createMultiLineField(name, value string) *fyne.Container {
	return container.NewHBox(
		widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewRichTextFromMarkdown(value),
	)
}

func (m *localMods) createCompatibility(compatibility *mods.ModCompatibility) fyne.CanvasObject {
	c := container.NewVBox(
		widget.NewLabelWithStyle("Compatibility", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	if len(compatibility.Requires) > 0 {
		c.Add(widget.NewLabelWithStyle("  Requires", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		for _, r := range compatibility.Requires {
			c.Add(widget.NewLabel("  - " + r.Name + ": " + r.Source))
		}
	}
	if len(compatibility.Requires) > 0 {
		c.Add(widget.NewLabelWithStyle("  Forbids", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		for _, r := range compatibility.Requires {
			c.Add(widget.NewLabel("  - " + r.Name + ": " + r.Source))
		}
	}
	return c
}

func (m *localMods) addFromFile() {
	if file, err := zenity.SelectFile(
		zenity.Title("Select a mod file"),
		zenity.FileFilter{
			Name:     "mod file",
			Patterns: []string{"*.xml", "*.json"},
		}); err == nil {
		if err = managed.AddModFromFile(*state.CurrentGame, file); err != nil {
			dialog.ShowError(err, state.Window)
			return
		}
		state.Window.Content().Refresh()
	}
}

func (m *localMods) addFromUrl() {
	e := widget.NewEntry()
	dialog.ShowForm("Add Remote mod file", "Add", "Cancel",
		[]*widget.FormItem{widget.NewFormItem("URL", e)},
		func(ok bool) {
			if ok && e.Text != "" {
				if err := managed.AddModFromUrl(*state.CurrentGame, e.Text); err != nil {
					dialog.ShowError(err, state.Window)
					return
				}
				state.Window.Content().Refresh()
			}
		}, state.Window)
}
