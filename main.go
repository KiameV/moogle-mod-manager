package main

import (
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
	"github.com/aarzilli/nucular/style"
	"github.com/kiamev/pr-modsync/browser"
	"github.com/kiamev/pr-modsync/config"
	"github.com/kiamev/pr-modsync/mods"
	"github.com/kiamev/pr-modsync/ui/configure"
	"github.com/kiamev/pr-modsync/ui/local"
	state "github.com/kiamev/pr-modsync/ui/state"
	"image"
	"image/color"
	"time"
)

const version = "0.1.0"

var (
	status        string
	statusTimer   *time.Timer
	errTextEditor nucular.TextEditor
	wnd           nucular.MasterWindow
)

func main() {
	/*mo := mods.GameMod{
		Mod: mods.Mod{
			ID:               "metalliguyau.sprites.ffvipr",
			Name:             "Textbox Portraits",
			Author:           "metalliguyAU",
			Version:          "2.1",
			ReleaseDate:      "2022-07-07",
			Category:         "Textures",
			Description:      "Adds Portraits to the textboxes for the main characters! Includes Opera Omnia, OldVer and SNES Portraits, with 6 UI styles. This update now includes Kefka's portrait! To make this mod work properly, Shiryu's Classic Text Box Framework (FFVI) HAS to be installed first. Thanks to Shiryu for enabling a workaround for Kefka's portrait! It also replaces Mog, Gogo and Umaro's portraits with Cid, Leo and Gestahl's. The reason for this is because Mog, Gogo and Umaro have very little dialogue compared to Cid, Leo and Gestahl. It makes the game a lot more immersive with these new portraits! Now includes support for Chrono Trigger UI!",
			ReleaseNotes:     "",
			Link:             "https://www.nexusmods.com/finalfantasy6pixelremaster/mods/26",
			Preview:          "https://staticdelivery.nexusmods.com/mods/4335/images/26/26-1650565736-597895528.png",
			ModCompatibility: mods.ModCompatibility{},
			GameVersions:     nil,
			Downloadables: []mods.Download{
				{
					Name:        "Boarderless",
					Sources:     []string{"https://supporter-files.nexus-cdn.com/4335/26/FFVIPR_Textbox_Portraits_(Full_Frame)_metalliguy-26-2-1-1650566101.rar"},
					InstallType: "Archive",
				},
				{
					Name:        "Boardered",
					Sources:     []string{"https://supporter-files.nexus-cdn.com/4335/26/FFVIPR_Textbox_Portraits_(Full_Frame_Bordered)_metalliguy-26-1-1-1650566164.rar"},
					InstallType: "Archive",
				},
			},
			Configurations: []mods.Configuration{
				{
					Name:        "Frame Type",
					Description: "Choose a frame type for the textbox portraits.",
					Choices: []mods.Choice{
						{
							Description:           "Borderless",
							Preview:               "https://staticdelivery.nexusmods.com/mods/4335/images/26/26-1650565736-597895528.png",
							NextConfigurationName: toPtr("FullFramePortraitType"),
						},
						{
							Description:           "Bordered",
							Preview:               "https://staticdelivery.nexusmods.com/mods/4335/images/26/26-1650565780-1285802872.png",
							NextConfigurationName: toPtr("BorderedPortraitType"),
						},
					},
				},
				{
					Name:        "FullFramePortraitType",
					Description: "Choose the type of portrait to use.",
					Choices: []mods.Choice{
						{
							Description: "Stock UI SNES",
							Preview:     "https://staticdelivery.nexusmods.com/mods/4335/images/26/26-1650565736-597895528.png",
							DownloadFiles: mods.DownloadFiles{
								DownloadName: "Boarderless",
								Files: []mods.ModFile{
									{
										From: "Stock UI + SNES",
										To:   toPtr("FINAL FANTASY VI PR/FINAL FANTASY VI_Data/StreamingAssets/aa/StandaloneWindows64"),
									},
								},
							},
						},
						{
							Description: "Stock UI OO",
							Preview:     "https://staticdelivery.nexusmods.com/mods/4335/images/26/26-1650565736-597895528.png",
							DownloadFiles: mods.DownloadFiles{
								DownloadName: "Boarderless",
								Files: []mods.ModFile{
									{
										From: "Stock UI + OO",
										To:   toPtr("FINAL FANTASY VI PR/FINAL FANTASY VI_Data/StreamingAssets/aa/StandaloneWindows64"),
									},
								},
							},
						},
						{
							Description: "Stock UI Old Version",
							Preview:     "https://staticdelivery.nexusmods.com/mods/4335/images/26/26-1650565736-597895528.png",
							DownloadFiles: mods.DownloadFiles{
								DownloadName: "Boarderless",
								Files: []mods.ModFile{
									{
										From: "Stock UI + OO",
										To:   toPtr("FINAL FANTASY VI PR/FINAL FANTASY VI_Data/StreamingAssets/aa/StandaloneWindows64"),
									},
								},
							},
						},
					},
				},
			},
		},
		Enabled: false,
	}

	result, _ := xml.Marshal(mo.Mod)
	println(result)*/

	errTextEditor.Flags = nucular.EditReadOnly | nucular.EditSelectable | nucular.EditSelectable | nucular.EditMultiline
	var (
		x = config.Get().WindowX
		y = config.Get().WindowY
	)
	if x == 0 || y == 0 {
		x = config.WindowWidth
		y = config.WindowHeight
	}
	wnd = nucular.NewMasterWindowSize(0, "FF PR Mod Manager - "+version, image.Point{X: x, Y: y}, updateWindow)
	wnd.SetStyle(style.FromTable(customTheme, 1.2))
	wnd.Main()
}

func updateWindow(w *nucular.Window) {
	var err error
	w.MenubarBegin()
	w.Row(12).Static(100, 100, 300, 200)
	if w := w.Menu(label.TA("Game", "LC"), 100, nil); w != nil {
		w.Row(12).Dynamic(1)
		if w.MenuItem(label.TA("I", "LC")) {
			state.Game = config.I
			state.CurrentUI = state.LocalMods
			w.Close()
		} else if w.MenuItem(label.TA("II", "LC")) {
			state.Game = config.II
			state.CurrentUI = state.LocalMods
			w.Close()
		} else if w.MenuItem(label.TA("III", "LC")) {
			state.Game = config.III
			state.CurrentUI = state.LocalMods
			w.Close()
		} else if w.MenuItem(label.TA("IV", "LC")) {
			state.Game = config.IV
			state.CurrentUI = state.LocalMods
			w.Close()
		} else if w.MenuItem(label.TA("V", "LC")) {
			state.Game = config.V
			state.CurrentUI = state.LocalMods
			w.Close()
		} else if w.MenuItem(label.TA("VI", "LC")) {
			state.Game = config.VI
			state.CurrentUI = state.LocalMods
			w.Close()
		}
	}

	if state.Game != config.None && state.CurrentUI == state.LocalMods {
		if w.MenuItem(label.TA("Add Mod", "LC")) {
			w.Close()
		}
	} else {
		w.Spacing(1)
	}

	if w := w.Menu(label.TA("Check For Update", "LC"), 300, nil); w != nil {
		var hasNewer bool
		var latest string
		if hasNewer, latest, err = browser.CheckForUpdate(version); err != nil {
			popupErr(w, err)
		}
		if hasNewer {
			browser.Update(latest)
		} else {
			status = "version is current"
		}
		w.Close()
	}

	popupErr(w, state.Errors...)

	if status != "" {
		w.Label("Status: "+status, "RC")
		if statusTimer != nil {
			statusTimer.Stop()
		}
		statusTimer = time.AfterFunc(2*time.Second, func() { status = "" })

	} else {
		w.Spacing(1)
	}
	w.MenubarEnd()

	if state.Game != config.None {
		switch state.CurrentUI {
		case state.LocalMods:
			var gm mods.GameMods
			if gm, err = mods.GetGameMods(state.Game); err != nil {
				popupErr(w, err)
				return
			}
			local.Draw(w, gm)
		case state.Configure:
			configure.Draw(w)
		case state.None:
			break
		}
	}
}

func popupErr(w *nucular.Window, errs ...error) {
	if len(errs) == 0 {
		return
	}
	for _, err := range errs {
		errTextEditor.Buffer = []rune(err.Error())
		w.Master().PopupOpen("Error", nucular.WindowMovable|nucular.WindowTitle|nucular.WindowDynamic, rect.Rect{X: 20, Y: 100, W: 700, H: 600}, true,
			func(w *nucular.Window) {
				w.Row(300).Dynamic(1)
				errTextEditor.Edit(w)
				w.Row(25).Dynamic(1)
				if w.Button(label.T("OK"), false) {
					w.Close()
				}
			})
	}
}

var customTheme = style.ColorTable{
	ColorText:                  color.RGBA{0, 0, 0, 255},
	ColorWindow:                color.RGBA{255, 255, 255, 255},
	ColorHeader:                color.RGBA{242, 242, 242, 255},
	ColorHeaderFocused:         color.RGBA{0xc3, 0x9a, 0x9a, 255},
	ColorBorder:                color.RGBA{0, 0, 0, 255},
	ColorButton:                color.RGBA{185, 185, 185, 255},
	ColorButtonHover:           color.RGBA{215, 215, 215, 255},
	ColorButtonActive:          color.RGBA{200, 200, 200, 255},
	ColorToggle:                color.RGBA{225, 225, 225, 255},
	ColorToggleHover:           color.RGBA{200, 200, 200, 255},
	ColorToggleCursor:          color.RGBA{30, 30, 30, 255},
	ColorSelect:                color.RGBA{175, 175, 175, 255},
	ColorSelectActive:          color.RGBA{190, 190, 190, 255},
	ColorSlider:                color.RGBA{190, 190, 190, 255},
	ColorSliderCursor:          color.RGBA{215, 215, 215, 255},
	ColorSliderCursorHover:     color.RGBA{235, 235, 235, 255},
	ColorSliderCursorActive:    color.RGBA{225, 225, 225, 255},
	ColorProperty:              color.RGBA{225, 225, 225, 255},
	ColorEdit:                  color.RGBA{245, 245, 245, 255},
	ColorEditCursor:            color.RGBA{0, 0, 0, 255},
	ColorCombo:                 color.RGBA{225, 225, 225, 255},
	ColorChart:                 color.RGBA{160, 160, 160, 255},
	ColorChartColor:            color.RGBA{45, 45, 45, 255},
	ColorChartColorHighlight:   color.RGBA{255, 0, 0, 255},
	ColorScrollbar:             color.RGBA{180, 180, 180, 255},
	ColorScrollbarCursor:       color.RGBA{140, 140, 140, 255},
	ColorScrollbarCursorHover:  color.RGBA{150, 150, 150, 255},
	ColorScrollbarCursorActive: color.RGBA{160, 160, 160, 255},
	ColorTabHeader:             color.RGBA{210, 210, 210, 255},
}

func toPtr(s string) *string {
	return &s
}
