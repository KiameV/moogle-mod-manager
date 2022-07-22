package local

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/ncruces/zenity"
)

func New() state.Screen {
	return &localMods{}
}

type localMods struct {
}

func (m *localMods) Draw(w fyne.Window) {
	var (
		selectable = managed.GetMods(*state.CurrentGame)
		modList    = widget.NewList(
			func() int { return len(selectable) },
			func() fyne.CanvasObject { return widget.NewLabel("") },
			func(id widget.ListItemID, object fyne.CanvasObject) {
				object.(*widget.Label).SetText(selectable[id].Mod.Name)
			})
		removeButton = widget.NewButton("Remove", func() {

		})
		modDetails = container.NewScroll(container.NewMax())
	)
	modList.OnSelected = func(id widget.ListItemID) {
		removeButton.Enable()
		modDetails.Content = m.createPreview(selectable[id].Mod)
	}

	split := container.NewHSplit(
		container.NewBorder(container.NewVBox(
			widget.NewLabelWithStyle(config.GameNameString(*state.CurrentGame), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			container.NewHBox(
				widget.NewButton("Add File", func() {
					if file, err := zenity.SelectFile(
						zenity.Title("Select a mod file"),
						zenity.FileFilter{
							Name:     "mod file",
							Patterns: []string{"*.xml", "*.json"},
						}); err == nil {
						if err = managed.AddModFromFile(*state.CurrentGame, file); err != nil {
							dialog.ShowError(err, w)
						}
					}
				}),
				widget.NewButton("Add Remote", func() {
					e := widget.NewEntry()
					dialog.ShowForm("Add Remote mod file", "Add", "Cancel",
						[]*widget.FormItem{widget.NewFormItem("URL", e)},
						func(ok bool) {
							if ok && e.Text != "" {
								if err := managed.AddModFromUrl(*state.CurrentGame, e.Text); err != nil {
									dialog.ShowError(err, w)
									return
								}
							}
						}, w)
				}),
				removeButton)), nil, nil,
			container.NewVScroll(modList)),
		modDetails)
	//split.Offset
	w.SetContent(split)
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
	if mod.Preview != "" {
		c.Add(canvas.NewImageFromURI(storage.NewFileURI(mod.Preview)))
	}
	if mod.ModCompatibility != nil && mod.ModCompatibility.HasItems() {
		c.Add(m.createCompatibility(mod.ModCompatibility))
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
