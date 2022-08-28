package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"strings"
)

type modKindDef struct {
	*entryManager
	kind *mods.Kind
}

func newModKindDef(kind *mods.Kind) *modKindDef {
	d := &modKindDef{
		entryManager: newEntryManager(),
		kind:         kind,
	}
	return d
}

func (d *modKindDef) draw() fyne.CanvasObject {
	switch *d.kind {
	case mods.Hosted:
		return widget.NewForm(d.getFormItem("'Mod File' Links"))
	case mods.Nexus:
		return widget.NewForm(d.getFormItem("Mod ID"))
	}
	panic("unknown mod kind")
}

func (d *modKindDef) compile() *mods.ModKind {
	k := mods.ModKind{}
	switch *d.kind {
	case mods.Hosted:
		k.Kind = mods.Hosted
		k.Hosted = &mods.HostedModKind{
			ModFileLinks: d.getStrings("'Mod File' Links", ","),
		}
	case mods.Nexus:
		k.Kind = mods.Nexus
		k.Nexus = &mods.NexusModKind{
			ID: d.getString("Mod ID"),
		}
	default:
		panic("unknown mod kind")
	}
	return &k
}

func (d *modKindDef) set(k *mods.ModKind) {
	d.Clear()
	if k == nil {
		return
	}
	if k.Kind == mods.Hosted {
		d.createFormItem("'Mod File' Links", strings.Join(k.Hosted.ModFileLinks, ", "))
	} else {
		d.createFormItem("Mod ID", k.Nexus.ID)
	}
}

func (d *modKindDef) Clear() {
	d.createFormItem("'Mod File' Links", "")
	d.createFormItem("Mod ID", "")
}
