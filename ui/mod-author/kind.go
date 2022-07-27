package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"strings"
)

type modKindDef struct {
	*entryManager
	hosted     *widget.Form
	nexus      *widget.Form
	main       *fyne.Container
	kindSelect *widget.RadioGroup
}

func newModKindDef() *modKindDef {
	return &modKindDef{
		entryManager: newEntryManager(),
		main:         container.NewMax(),
		hosted:       widget.NewForm(),
		nexus:        widget.NewForm(),
	}
}

func (d *modKindDef) draw() fyne.CanvasObject {
	d.Clear()
	if len(d.hosted.Items) == 0 {
		d.hosted.AppendItem(d.getFormItem("Version"))
		d.hosted.AppendItem(d.getFormItem("'Mod File' Links"))
	}
	if len(d.nexus.Items) == 0 {
		d.nexus.AppendItem(d.getFormItem("Mod ID"))
	}
	d.kindSelect = widget.NewRadioGroup(mods.Kinds, func(kind string) {
		if kind == string(mods.Hosted) {
			d.main.RemoveAll()
			d.main.Add(d.hosted)
		} else {
			d.main.RemoveAll()
			d.main.Add(d.nexus)
		}
	})
	d.kindSelect.SetSelected(string(mods.Hosted))

	return container.NewBorder(d.kindSelect, nil, nil, nil, d.main)
}

func (d *modKindDef) compile() mods.ModKind {
	k := mods.ModKind{}
	switch d.kindSelect.Selected {
	case string(mods.Hosted):
		k.Kind = mods.Hosted
		k.Hosted = &mods.HostedModKind{
			Version:      d.getString("Version"),
			ModFileLinks: d.getStrings("'Mod File' Links", ","),
		}
	default: // string(mods.Nexus):
		k.Kind = mods.Nexus
		k.Nexus = &mods.NexusModKind{
			ID: d.getString("Mod ID"),
		}
	}
	return k
}

func (d *modKindDef) set(k *mods.ModKind) {
	d.Clear()
	if k.Kind == mods.Hosted {
		d.createFormItem("Version", k.Hosted.Version)
		d.createFormItem("'Mod File' Links", strings.Join(k.Hosted.ModFileLinks, ", "))
		d.kindSelect.SetSelected(string(mods.Hosted))
	} else {
		d.createFormItem("Mod ID", k.Nexus.ID)
		d.kindSelect.SetSelected(string(mods.Nexus))
	}
}

func (d *modKindDef) Clear() {
	d.createFormItem("Version", "")
	d.createFormItem("'Mod File' Links", "")
	d.createFormItem("Mod ID", "")
}
