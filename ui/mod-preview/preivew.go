package mod_preview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/discover"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"net/url"
	"strings"
)

type ModPreviewOptions struct {
	UpdateCallback func(mod mods.TrackedMod)
	TrackedMod     mods.TrackedMod
}

func CreatePreview(mod *mods.Mod, options ...ModPreviewOptions) fyne.CanvasObject {
	c := container.NewVBox()
	if len(options) > 0 && options[0].UpdateCallback != nil && options[0].TrackedMod != nil && options[0].TrackedMod.UpdatedMod() != nil {
		c.Add(widget.NewButton("Update", func() {
			options[0].UpdateCallback(options[0].TrackedMod)
		}))
	}
	c.Add(createField("Name", string(mod.Name)))
	c.Add(createLink("Link", mod.Link))
	c.Add(createField("Author", mod.Author))
	c.Add(createField("Version", mod.Version))
	if mod.Category != "" {
		c.Add(createField("Category", string(mod.Category)))
	}
	c.Add(createField("Release Date", mod.ReleaseDate))

	text := widget.NewRichTextFromMarkdown(strings.ReplaceAll(mod.Description, "\r", ""))
	text.Wrapping = fyne.TextWrapWord
	tabs := container.NewAppTabs(
		container.NewTabItem("Description", text),
	)
	if mod.ReleaseNotes != "" {
		text = widget.NewRichTextFromMarkdown(strings.ReplaceAll(mod.ReleaseNotes, "\r", ""))
		text.Wrapping = fyne.TextWrapWord
		container.NewTabItem("Release Notes", text)
	}
	if mod.ModCompatibility != nil && mod.ModCompatibility.HasItems() {
		if len(mod.Games) > 0 {
			if game, err := config.GameDefFromID(mod.Games[0].ID); err == nil {
				tabs.Append(container.NewTabItem("Compatibility", createCompatibility(game, mod.ModCompatibility)))
			}
		}
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

func createCompatibility(game config.GameDef, compatibility *mods.ModCompatibility) fyne.CanvasObject {
	var (
		c = container.NewVBox(
			widget.NewLabelWithStyle("Compatibility", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		)
		name string
		err  error
	)

	// Requires
	if len(compatibility.Requires) > 0 {
		c.Add(widget.NewLabelWithStyle("  Requires", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		for _, r := range compatibility.Requires {
			name, err = discover.GetDisplayName(game, r.ModID())
			if err != nil {
				util.ShowErrorLong(err)
			}
			c.Add(widget.NewLabel("  - " + name))
		}
	}

	// Forbids
	if len(compatibility.Forbids) > 0 {
		c.Add(widget.NewLabelWithStyle("  Forbids", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		for _, r := range compatibility.Forbids {
			name, err = discover.GetDisplayName(game, r.ModID())
			if err != nil {
				util.ShowErrorLong(err)
			}
			c.Add(widget.NewLabel("  - " + name))
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
