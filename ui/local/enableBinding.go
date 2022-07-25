package local

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/kiamev/moogle-mod-manager/mods/managed/model"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

type enableBind struct {
	binding.ExternalBool
	localUI *localUI
	mod     *model.TrackedMod
	orig    bool
}

func newEnableBind(localUI *localUI, mod *model.TrackedMod) *enableBind {
	bb := binding.BindBool(&mod.Enabled)
	b := &enableBind{
		ExternalBool: bb,
		localUI:      localUI,
		mod:          mod,
		orig:         mod.Enabled,
	}
	bb.AddListener(b)
	return b
}

func (b *enableBind) DataChanged() {
	if b.orig != b.mod.Enabled {
		if ok := b.localUI.toggleEnabled(*state.CurrentGame, b.mod); ok {
			b.orig = b.mod.Enabled
		}
	}
}
