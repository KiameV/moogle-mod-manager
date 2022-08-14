package mod_preview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	"net/url"
)

type ModPreviewOptions struct {
	UpdateCallback func(mod *model.TrackedMod)
	TrackedMod     *model.TrackedMod
}

func CreatePreview(mod *mods.Mod, options ...ModPreviewOptions) fyne.CanvasObject {
	c := container.NewVBox()
	if len(options) > 0 && options[0].UpdateCallback != nil && options[0].TrackedMod != nil && options[0].TrackedMod.UpdatedMod != nil {
		c.Add(widget.NewButton("Update", func() {
			options[0].UpdateCallback(options[0].TrackedMod)
		}))
	}
	c.Add(createField("Name", mod.Name))
	c.Add(createLink("Link", mod.Link))
	c.Add(createField("Author", mod.Author))
	c.Add(createField("Version", mod.Version))
	//c.Add(createField("Category", mod.Category))
	c.Add(createField("Release Date", mod.ReleaseDate))

	tabs := container.NewAppTabs(
		container.NewTabItem("Description", widget.NewRichTextFromMarkdown(mod.Description)),
	)
	if mod.ReleaseNotes != "" {
		container.NewTabItem("Release Notes", widget.NewRichTextFromMarkdown(mod.ReleaseNotes))
	}
	if mod.ModCompatibility != nil && mod.ModCompatibility.HasItems() {
		tabs.Append(container.NewTabItem("Compatibility", createCompatibility(mod.ModCompatibility)))
	}
	if mod.DonationLinks != nil && len(mod.DonationLinks) > 0 {
		tabs.Append(container.NewTabItem("Donations", createDonationLinks(mod.DonationLinks)))
	}

	c = container.NewBorder(c, nil, nil, nil, tabs)
	if img := mod.Preview.Get(); img != nil {
		c = container.NewBorder(img, nil, nil, nil, c)
	}
	return container.NewScroll(c)
}

func createField(name, value string) *fyne.Container {
	return container.NewHBox(
		widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(value),
	)
}

func createLink(name, value string) *fyne.Container {
	url, err := url.ParseRequestURI(value)
	if err != nil {
		return createField(name, value)
	}
	return container.NewHBox(
		widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewHyperlink(value, url),
	)
}

func createCompatibility(compatibility *mods.ModCompatibility) fyne.CanvasObject {
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

func createDonationLinks(links []*mods.DonationLink) fyne.CanvasObject {
	c := container.NewVBox(
		widget.NewLabelWithStyle("Support Project", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	for _, r := range links {
		if u, err := url.Parse(r.Link); err != nil {
			c.Add(widget.NewLabel("  - " + r.Name + ": " + r.Link))
		} else {
			c.Add(container.NewHBox(widget.NewLabel("  - "+r.Name), widget.NewHyperlink(r.Link, u)))
		}
	}
	return c
}