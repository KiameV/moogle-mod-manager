package local

import (
	"github.com/aarzilli/nucular"
	"github.com/kiamev/pr-modsync/mods"
	"github.com/kiamev/pr-modsync/ui/configure"
	"github.com/kiamev/pr-modsync/ui/state"
	"github.com/kiamev/pr-modsync/ui/util"
)

const (
	All      = "All"
	Enabled  = "Enabled"
	Disabled = "Disabled"
)

var (
	filters  = []string{All, Enabled, Disabled}
	filter   = 0
	selected *mods.GameMod
)

func Draw(w *nucular.Window, gm mods.GameMods) {
	w.Row(4).Static()
	w.Row(14).Static(200, 100, 100)
	w.Label(gm.GetGameName(), "LC")
	filter = w.ComboSimple(filters, filter, 12)

	w.Row(700).Static(290, 10, 500)
	if sw := w.GroupBegin("mods", nucular.WindowBorder|nucular.WindowNoScrollbar); sw != nil {
		sw.Row(20).Static(270)
		for _, m := range gm.GetMods() {
			checked := m.Enabled
			sw.CheckboxText(m.Mod.Name, &checked)
			selected = m
		}
		sw.GroupEnd()
	}
	w.Spacing(1)
	if sw := w.GroupBegin("preview", nucular.WindowBorder|nucular.WindowNoScrollbar); sw != nil {
		if selected != nil {
			util.AddFieldValueText(sw, "Name:", selected.Mod.Name)
			util.AddFieldValueText(sw, "Version:", selected.Mod.Version)
			util.AddFieldValueText(sw, "Description:", selected.Mod.Description)
			util.AddFieldValueText(sw, "Category:", selected.Mod.Category)
			util.AddFieldValueText(sw, "Release Date:", selected.Mod.ReleaseDate)

			util.DrawImg(sw, selected.Mod.Preview)

			sw.Row(20).Static(60)
			if sw.ButtonText("Enable") {
				configure.Initialize(selected, state.CurrentUI)
				state.CurrentUI = state.Configure
			}
		}
		sw.GroupEnd()
	}
}
