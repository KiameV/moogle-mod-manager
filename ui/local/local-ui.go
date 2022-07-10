package local

import (
	"github.com/aarzilli/nucular"
	"github.com/kiamev/pr-modsync/mods"
	"image"
	"image/draw"
	"image/png"
	"os"
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
			addFieldValueText(sw, "Name:", selected.Mod.Name)
			addFieldValueText(sw, "Version:", selected.Mod.Version)
			addFieldValueText(sw, "Description:", selected.Mod.Description)
			addFieldValueText(sw, "Category:", selected.Mod.Category)
			addFieldValueText(sw, "Release Date:", selected.Mod.ReleaseDate)

			fr, _ := os.Open(selected.Mod.Preview)
			img, _ := png.Decode(fr)
			i := image.NewRGBA(image.Rect(0, 0, 400, 225))
			draw.Draw(i, i.Bounds(), img, image.Point{}, draw.Src)

			sw.Row(225).Static(400)
			sw.Image(i)

			sw.Row(20).Static(60)
			sw.ButtonText("Enable")

		}
		sw.GroupEnd()
	}
}

func addFieldValueText(w *nucular.Window, label string, value string) {
	w.Row(12).Static(w.Bounds.W - 20)
	w.Label(label, "LC")
	w.Row(12).Static(w.Bounds.W - 20)
	w.Label(value, "LC")
	w.Row(8).Static()
}
