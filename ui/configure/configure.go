package configure

import (
	"fmt"
	"github.com/aarzilli/nucular"
	"github.com/kiamev/pr-modsync/mods"
	"github.com/kiamev/pr-modsync/mods/managed"
	"github.com/kiamev/pr-modsync/ui/state"
	"github.com/kiamev/pr-modsync/ui/util"
)

var (
	gm            *managed.GameMod
	currentConfig mods.Configuration
	previousUI    state.GUI
)

func Initialize(gameMod *managed.GameMod, callingUI state.GUI) {
	if len(gameMod.Mod.Configurations) == 0 {
		state.CurrentUI = callingUI
	}
	gm = gameMod
	currentConfig = gm.Mod.Configurations[0]
	previousUI = callingUI
}

func Draw(w *nucular.Window) {
	w.Row(6).Static()
	w.Row(14).Static(w.Bounds.W - 20)
	w.Label(gm.Mod.Name, "LC")

	w.Row(6).Static()
	w.Row(12).Static(w.Bounds.W - 20)
	w.Label(currentConfig.Name, "LC")

	w.Row(6).Static()
	w.Row(12).Static(w.Bounds.W - 20)
	w.Label(currentConfig.Description, "LC")

	util.DrawImg(w, currentConfig.Preview)

	for i, c := range currentConfig.Choices {
		w.Row(300).Static(20, 380)
		w.Spacing(1)
		if sw := w.GroupBegin(fmt.Sprintf("%s-%d", currentConfig.Name, i), nucular.WindowBorder|nucular.WindowNoScrollbar); sw != nil {
			sw.Row(12).Static(sw.Bounds.W - 20)
			sw.Label(c.Description, "LC")
			util.DrawImg(sw, c.Preview)
			sw.Row(12).Static(50)
			if sw.ButtonText("Select") {
				if c.NextConfigurationName != nil && *c.NextConfigurationName != "" {
					for _, cfg := range gm.Mod.Configurations {
						if *c.NextConfigurationName == cfg.Name {
							currentConfig = cfg
							break
						}
					}
				} else {
					// TODO Copy files
					state.CurrentUI = previousUI
				}
			}
			sw.GroupEnd()
		}
	}
}
