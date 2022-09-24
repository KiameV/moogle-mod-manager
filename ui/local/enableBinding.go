package local

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
	ci "github.com/kiamev/moogle-mod-manager/ui/config-installer"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"path/filepath"
)

type enableBind struct {
	binding.Bool
	tm          *mods.TrackedMod
	start       func() bool
	showWorking func()
	done        mods.DoneCallback
}

func newEnableBind(tm *mods.TrackedMod, start func() bool, showWorking func(), done mods.DoneCallback) *enableBind {
	b := &enableBind{
		Bool:        binding.NewBool(),
		tm:          tm,
		start:       start,
		showWorking: showWorking,
		done:        done,
	}
	_ = b.Set(tm.Enabled)
	b.AddListener(b)
	return b
}

func (b *enableBind) DataChanged() {
	isChecked, _ := b.Get()
	if isChecked != b.tm.Enabled {
		if isChecked {
			if !b.start() {
				_ = b.Set(false)
				return
			}
			if err := b.EnableMod(); err != nil {
				if err != nil {
					_ = b.Set(false)
					b.done(mods.Error, err)
				}
			}
		} else {
			if !b.start() {
				_ = b.Set(false)
				return
			}
			err := b.DisableMod()
			_ = b.Set(false)
			if err != nil {
				b.done(mods.Error, err)
			} else {
				b.done(mods.Ok)
			}
		}
	}
}

func (b *enableBind) EnableMod() (err error) {
	var (
		tm  = b.tm
		tis []*mods.ToInstall
	)

	if len(b.tm.Mod.Configurations) > 0 {
		err = b.enableModWithConfig()
	} else {
		tis, err = mods.NewToInstallForMod(tm.Mod.ModKind.Kind, tm.Mod, tm.Mod.AlwaysDownload)
		if err == nil {
			// Success
			err = b.enableMod(tis)
		}
	}
	return
}

func (b *enableBind) enableModWithConfig() (err error) {
	modPath := filepath.Join(config.Get().GetModsFullPath(*state.CurrentGame), b.tm.GetDirSuffix())
	if err = state.GetScreen(state.ConfigInstaller).(ci.ConfigInstaller).Setup(b.tm.Mod, modPath, b.enableMod); err != nil {
		// Failed to set up config installer screen
		return
	}
	state.ShowScreen(state.ConfigInstaller)
	return
}

func (b *enableBind) enableMod(toInstall []*mods.ToInstall) (err error) {
	return managed.EnableMod(mods.NewModEnabler(*state.CurrentGame, b.tm, toInstall, b.OnConflict, b.showWorking, func(result mods.Result, err ...error) {
		_ = b.Set(b.tm.Enabled)
		b.done(result, err...)
	}))
}

func (b *enableBind) DisableMod() error {
	return managed.DisableMod(*state.CurrentGame, b.tm)
}

func (b *enableBind) OnConflict(conflicts []*mods.FileConflict, confirmationCallback mods.ConflictChoiceCallback) {
	f := widget.NewForm()
	for _, c := range conflicts {
		var name string
		if m, ok := managed.TryGetMod(*state.CurrentGame, c.CurrentModID); !ok {
			name = string(c.CurrentModID)
		} else {
			name = m.DisplayName
		}
		f.Items = append(f.Items, widget.NewFormItem(
			filepath.Base(c.File),
			widget.NewSelect([]string{name, b.tm.Mod.Name}, c.OnChange)))
	}
	d := dialog.NewCustomConfirm("Conflicts", "ok", "cancel", container.NewVScroll(f), func(ok bool) {
		r := mods.Ok
		if !ok {
			r = mods.Cancel
		}
		confirmationCallback(r, conflicts)
	}, state.Window)
	d.Resize(fyne.NewSize(400, 400))
	d.Show()
}
