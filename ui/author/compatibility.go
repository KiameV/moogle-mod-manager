package author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/pr-modsync/mods"
	"strings"
)

type mcType bool

const (
	required  mcType = true
	forbidden mcType = false
)

func createCompatibilities(w fyne.Window) fyne.CanvasObject {
	var (
		requires *fyne.Container
		forbids  *fyne.Container
	)
	requires = container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Requires", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
				showModCompatForm(w, requires, nil, required)
			})),
	)
	forbids = container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Forbids", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
				showModCompatForm(w, forbids, nil, forbidden)
			})),
	)
	if mod.ModCompatibility != nil {
		for _, mc := range mod.ModCompatibility.Requires {
			requires.Objects = append(requires.Objects, createEditRemoveModCompatRow(w, requires, mc, required))
		}
		for _, mc := range mod.ModCompatibility.Forbids {
			forbids.Objects = append(forbids.Objects, createEditRemoveModCompatRow(w, forbids, mc, forbidden))
		}
	}
	return container.NewVScroll(container.NewVBox(requires, forbids))
}

func showModCompatForm(w fyne.Window, parent *fyne.Container, mc *mods.ModCompat, mct mcType) {
	var modID, versions, source, order string
	if mc != nil {
		modID = mc.ModID
		versions = strings.Join(mc.Versions, ", ")
		source = mc.Source
		if mc.Order != nil {
			order = string(*mc.Order)
		}
	}
	setFormItem("compatModID", modID)
	setFormItem("compatVersions", versions)
	setFormItem("compatSources", source)
	setFormSelect("compatOrder", []string{"", string(mods.Before), string(mods.After)}, order)

	var p *widget.PopUp
	p = widget.NewModalPopUp(container.NewVScroll(container.NewVBox(
		widget.NewLabelWithStyle("Mod Compatibility", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewForm(
			getFormItem("Mod ID", "compatModID"),
			getFormItem("Versions", "compatVersions"),
			getFormItem("Sources", "compatSources"),
			getFormItem("Order (optional)", "compatOrder")),
		container.NewCenter(container.NewHBox(
			widget.NewButton("Save", func() {
				if mc == nil {
					if mod.ModCompatibility == nil {
						mod.ModCompatibility = &mods.ModCompatibility{}
					}
					mc = &mods.ModCompat{
						ModID:    getFormString("compatModID"),
						Versions: strings.Split(getFormString("compatVersions"), ", "),
						Source:   getFormString("compatSources"),
					}
					order = getFormString("compatOrder")
					if order != "" {
						mc.Order = (*mods.ModCompatOrder)(&order)
					}
					mcs := getModCompats(mod, mct)
					*mcs = append(*mcs, mc)
					parent.Objects = append(parent.Objects, createEditRemoveModCompatRow(w, parent, mc, mct))
				} else {

				}
				p.Hide()
			}),
			widget.NewButton("Cancel", func() { p.Hide() }))))),
		w.Canvas())
	p.Resize(fyne.Size{Width: 600, Height: 400})
	p.Show()
}

func createEditRemoveModCompatRow(w fyne.Window, parent *fyne.Container, mc *mods.ModCompat, mct mcType) fyne.CanvasObject {
	sb := strings.Builder{}
	sb.WriteString("    " + mc.ModID)
	if len(mc.Versions) > 0 {
		sb.WriteString("\n      " + strings.Join(mc.Versions, ", "))
	}
	sb.WriteString("\n      " + mc.Source)
	if mc.Order != nil {
		sb.WriteString("\n      " + string(*mc.Order))
	}

	return container.NewHBox(
		widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
			showModCompatForm(w, parent, mc, mct)
		}),
		widget.NewButtonWithIcon("Remove", theme.ContentRemoveIcon(), func() {
			mcs := getModCompats(mod, mct)
			for i, s := range *mcs {
				if s.ModID == mc.ModID {
					*mcs = append((*mcs)[:i], (*mcs)[i+1:]...)
					break
				}
			}
			for i, c := range parent.Objects {
				if cc, ok := c.(*fyne.Container); ok {
					for _, ccc := range cc.Objects {
						if rt, k := ccc.(*widget.RichText); k && len(rt.Segments) > 0 {
							if strings.Index(rt.Segments[0].Textual(), mc.ModID) == 0 {
								parent.Objects = append(parent.Objects[:i], parent.Objects[i+1:]...)
								break
							}
						}
					}
				}
			}
		}),
		widget.NewRichTextFromMarkdown(sb.String()))
}

func getModCompats(mod *mods.Mod, mct mcType) *[]*mods.ModCompat {
	if mct == required {
		return &mod.ModCompatibility.Requires
	}
	return &mod.ModCompatibility.Forbids
}

/*func newModCompatForm() []*widget.FormItem {
	return []*widget.FormItem{}
}

func createHBoxLabelValue(label, value string) *fyne.Container {
	return container.NewHBox(
		widget.NewLabelWithStyle(label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(value))
}*/
