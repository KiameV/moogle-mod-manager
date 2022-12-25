package local

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/kiamev/moogle-mod-manager/actions"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/util"
)

type (
	enableBind struct {
		binding.Bool
		parent             *localUI
		tm                 mods.TrackedMod
		start              func() bool
		modEnabledCallback func(result actions.Result)
	}
)

func newEnableBind(parent *localUI, tm mods.TrackedMod, start func() bool, modEnabledCallback func(r actions.Result)) *enableBind {
	var (
		b = &enableBind{
			parent:             parent,
			Bool:               binding.NewBool(),
			tm:                 tm,
			start:              start,
			modEnabledCallback: modEnabledCallback,
		}
	)
	_ = b.Set(tm.Enabled())
	b.AddListener(b)
	return b
}

func (b *enableBind) DataChanged() {
	var (
		isChecked, _ = b.Get()
		tmEnabled    = b.tm.Enabled()
		action       actions.Action
		err          error
	)
	if isChecked != tmEnabled {
		if !b.start() {
			return
		}
		if isChecked {
			// Enable
			if action, err = actions.New(actions.Install, state.CurrentGame, b.tm, b.ActionDone); err != nil {
				util.ShowErrorLong(err)
				_ = b.Set(false)
			} else if err = action.Run(); err != nil {
				util.ShowErrorLong(err)
				_ = b.Set(false)
			}
		} else {
			// Disable
			if action, err = actions.New(actions.Uninstall, state.CurrentGame, b.tm, b.ActionDone); err != nil {
				util.ShowErrorLong(err)
				_ = b.Set(true)
			} else if err = action.Run(); err != nil {
				util.ShowErrorLong(err)
				_ = b.Set(true)
			}
		}
	}
}

func (b *enableBind) ActionDone(r actions.Result) {
	if r.Err != nil {
		util.ShowErrorLong(r.Err)
	}
	b.modEnabledCallback(r)
	b.parent.ModList.Refresh()
}
