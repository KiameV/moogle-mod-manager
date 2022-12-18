package local

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/kiamev/moogle-mod-manager/actions"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/util"
)

type enableBind struct {
	binding.Bool
	tm          mods.TrackedMod
	start       func() bool
	showWorking func()
	hideWorking func()
	//done        mods.DoneCallback
}

func newEnableBind(tm mods.TrackedMod, start func() bool, showWorking func(), hideWorking func() /*, done mods.DoneCallback*/) *enableBind {
	b := &enableBind{
		Bool:        binding.NewBool(),
		tm:          tm,
		start:       start,
		showWorking: showWorking,
		hideWorking: hideWorking,
		//done:        done,
	}
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
		if isChecked {
			// Enable
			if !b.start() {
				_ = b.Set(false)
				return
			}
			if action, err = actions.New(actions.Install, b.newActionParams()); err != nil {
				util.ShowErrorLong(err)
				_ = b.Set(false)
				return
			}
			if err = action.Run(); err != nil {
				util.ShowErrorLong(err)
				_ = b.Set(false)
				return
			}
		} else {
			// Disable
			if !b.start() {
				_ = b.Set(false)
				return
			}
			if action, err = actions.New(actions.Uninstall, b.newActionParams()); err != nil {
				util.ShowErrorLong(err)
				_ = b.Set(true)
				return
			}
			if err = action.Run(); err != nil {
				util.ShowErrorLong(err)
				_ = b.Set(true)
				return
			}
		}
	}
}

func (b *enableBind) newActionParams() actions.Params {
	return actions.Params{
		Game: state.CurrentGame,
		Mod:  b.tm,
		WorkingDialog: actions.WorkingDialog{
			Show: b.showWorking,
			Hide: b.hideWorking,
		},
	}
}
