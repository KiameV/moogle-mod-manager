package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/kiamev/moogle-mod-manager/mods"
)

type modCompatabilityDef struct {
	requires *modCompatsDef
	forbids  *modCompatsDef
}

func newModCompatibilityDef() *modCompatabilityDef {
	return &modCompatabilityDef{
		requires: newModCompatsDef("Requires"),
		forbids:  newModCompatsDef("Forbids"),
	}
}

func (d *modCompatabilityDef) draw() fyne.CanvasObject {
	return container.NewVScroll(container.NewVBox(
		d.requires.draw(),
		d.forbids.draw(),
	))
}

func (d *modCompatabilityDef) compile() *mods.ModCompatibility {
	return &mods.ModCompatibility{
		Requires: d.requires.compile(),
		Forbids:  d.forbids.compile(),
	}
}
