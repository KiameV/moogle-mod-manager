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
	d := &modKindDef{
		entryManager: newEntryManager(),
		main:         container.NewMax(),
		hosted:       widget.NewForm(),
		nexus:        widget.NewForm(),
	}
	d.kindSelect = widget.NewRadioGroup(mods.Kinds, d.onSelectChange)
	return d
}

func (d *modKindDef) draw() fyne.CanvasObject {
	if len(d.hosted.Items) == 0 {
		d.hosted.AppendItem(d.getFormItem("Version"))
		d.hosted.AppendItem(d.getFormItem("'Mod File' Links"))
	}
	if len(d.nexus.Items) == 0 {
		d.nexus.AppendItem(d.getFormItem("Mod ID"))
	}
	d.kindSelect.SetSelected(string(mods.Hosted))

	return container.NewBorder(d.kindSelect, nil, nil, nil, d.main)
}

func (d *modKindDef) onSelectChange(kind string) {
	if kind == string(mods.Hosted) {
		d.main.RemoveAll()
		d.main.Add(d.hosted)
	} else {
		d.main.RemoveAll()
		d.main.Add(d.nexus)
	}
}

func (d *modKindDef) compile() *mods.ModKind {
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
	return &k
}

func (d *modKindDef) set(k *mods.ModKind) {
	d.Clear()
	if k == nil {
		return
	}
	if k.Kind == mods.Hosted {
		d.kindSelect.SetSelected(string(mods.Hosted))
		d.createFormItem("Version", k.Hosted.Version)
		d.createFormItem("'Mod File' Links", strings.Join(k.Hosted.ModFileLinks, ", "))
	} else {
		d.kindSelect.SetSelected(string(mods.Nexus))
		d.createFormItem("Mod ID", k.Nexus.ID)
	}
}

func (d *modKindDef) Clear() {
	d.createFormItem("Version", "")
	d.createFormItem("'Mod File' Links", "")
	d.createFormItem("Mod ID", "")
}
