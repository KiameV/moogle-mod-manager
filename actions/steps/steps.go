package steps

import (
	"fmt"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/mods/managed"
)

type (
	State struct {
		Game       config.GameDef
		Mod        mods.TrackedMod
		Downloaded []string
		Extracted  []string
		ToInstall  []string
	}
	Step func(state *State) error
)

func NewState(game config.GameDef, mod mods.TrackedMod) *State {
	return &State{
		Game: game,
		Mod:  mod,
	}
}

func VerifyEnable(state *State) error {
	var (
		tm      = state.Mod
		c       = tm.Mod().ModCompatibility
		mc      *mods.ModCompat
		mod     mods.TrackedMod
		found   bool
		enabled bool
	)
	if c != nil {
		if len(c.Forbids) > 0 {
			for _, mc = range c.Forbids {
				if mod, found, enabled = managed.IsModEnabled(state.Game, mc.ModID()); found && enabled {
					return fmt.Errorf("[%s] cannot be enabled because [%s] is enabled", tm.DisplayName(), mod.DisplayName())
				}
			}
		}
		if len(c.Requires) > 0 {
			for _, mc = range c.Requires {
				mod, found, enabled = managed.IsModEnabled(state.Game, mc.ModID())
				if !found {
					return fmt.Errorf("[%s] cannot be enabled because [%s] is not enabled", tm.DisplayName(), mc.ModID())
				} else if !enabled {
					return fmt.Errorf("[%s] cannot be enabled because [%s] is not enabled", tm.DisplayName(), mod.DisplayName())
				}
			}
		}
	}
	return nil
}

func VerifyDisable(state *State) error {
	var (
		tm  = state.Mod
		mod = tm.Mod()
		c   = mod.ModCompatibility
		id  = mod.ModID
		mc  *mods.ModCompat
	)
	for _, m := range managed.GetMods(state.Game) {
		if m.ID() != tm.ID() && m.Enabled() {
			if c = m.Mod().ModCompatibility; c != nil {
				if len(c.Requires) > 0 {
					for _, mc = range c.Requires {
						if mc.ModID() == id {
							return fmt.Errorf("[%s] cannot be disabled because [%s] is enabled", tm.DisplayName(), m.DisplayName())
						}
					}
				}
			}
		}
	}
	return nil
}

func Download(state *State) error {
	return nil
}

func Extract(state *State) error {
	return nil
}

func PreInstall(state *State) error {
	return nil
}

func Conflicts(state *State) error {
	return nil
}

func Install(state *State) error {
	return nil
}

func InstallMod(state *State) error {
	return nil
}

func UninstallMove(state *State) error {
	return nil
}

func EnableMod(state *State) error {
	state.Mod.Enable()
	return managed.Save()
}

func DisableMod(state *State) error {
	state.Mod.Disable()
	return managed.Save()
}
