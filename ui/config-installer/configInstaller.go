package config_installer

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/util"
)

type ConfigInstaller interface {
	state.Screen
	Setup(mod *mods.Mod, baseDir string, done func(mods.Result, []*mods.ToInstall) error) error
}

func New() ConfigInstaller {
	return &configInstallerUI{
		choiceContainer: container.NewVBox(),
	}
}

type configInstallerUI struct {
	mod             *mods.Mod
	toInstall       []*mods.DownloadFiles
	prevConfigs     []*mods.Configuration
	choiceContainer *fyne.Container
	baseDir         string
	done            func(mods.Result, []*mods.ToInstall) error

	currentConfig *mods.Configuration
	currentChoice *mods.Choice
}

func (ui *configInstallerUI) PreDraw(fyne.Window, ...interface{}) error { return nil }

func (ui *configInstallerUI) OnClose() {}

func (ui *configInstallerUI) DrawAsDialog(fyne.Window) {}

func (ui *configInstallerUI) Setup(mod *mods.Mod, baseDir string, done func(mods.Result, []*mods.ToInstall) error) error {
	if len(mod.Configurations) == 0 || len(mod.Configurations[0].Choices) == 0 {
		return fmt.Errorf("no configurations for %s", mod.Name)
	}
	ui.mod = mod
	for _, ui.currentConfig = range mod.Configurations {
		if ui.currentConfig.Root {
			break
		}
	}
	if ui.currentConfig == nil || !ui.currentConfig.Root {
		return errors.New("could not find root configuration")
	}
	ui.prevConfigs = make([]*mods.Configuration, 0)
	ui.baseDir = baseDir
	ui.done = done
	ui.toInstall = make([]*mods.DownloadFiles, 0)
	ui.choiceContainer.RemoveAll()
	ui.toInstall = append(ui.toInstall, mod.AlwaysDownload...)
	return nil
}

func (ui *configInstallerUI) Draw(w fyne.Window) {
	state.SetBaseDir(ui.baseDir)
	c := container.NewVBox(
		widget.NewLabelWithStyle(ui.currentConfig.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewRichTextFromMarkdown(ui.currentConfig.Description),
		ui.getChoiceSelector(func(name string) {
			for _, ui.currentChoice = range ui.currentConfig.Choices {
				if ui.currentChoice.Name == name {
					ui.choiceContainer.RemoveAll()
					ui.drawChoiceInfo()
					break
				}
			}
		}))
	buttons := container.NewHBox(
		widget.NewButton("Select", func() {
			if ui.currentChoice == nil {
				return
			}
			ui.prevConfigs = append(ui.prevConfigs, ui.currentConfig)
			ui.toInstall = append(ui.toInstall, ui.currentChoice.DownloadFiles)
			if ui.currentChoice.NextConfigurationName == nil {
				tis, err := mods.NewToInstallForMod(ui.mod.ModKind.Kind, ui.mod, ui.toInstall)
				if err != nil {
					util.ShowErrorLong(err)
					state.ShowPreviousScreen()
					return
				}
				state.ShowPreviousScreen()
				if err = ui.done(mods.Ok, tis); err != nil {
					util.ShowErrorLong(err)
					return
				}
			} else {
				for _, ui.currentConfig = range ui.mod.Configurations {
					if ui.currentConfig.Name == *ui.currentChoice.NextConfigurationName {
						break
					}
				}
				ui.currentChoice = nil
				ui.choiceContainer.RemoveAll()
				ui.Draw(w)
			}
		}))
	if len(ui.prevConfigs) > 0 {
		buttons.Add(widget.NewButton("Back", func() {
			ui.popToInstall()
			ui.currentConfig = ui.popChoice()
			ui.choiceContainer.RemoveAll()
			ui.Draw(w)
		}))
	}
	c.Add(buttons)
	if img := ui.currentConfig.Preview.Get(); img != nil {
		c = container.NewBorder(img, nil, nil, nil, c)
	}
	cnlButton := widget.NewButton("Cancel", func() {
		_ = ui.done(mods.Cancel, nil)
		state.ShowPreviousScreen()
	})
	w.SetContent(
		container.NewBorder(container.NewHBox(cnlButton), nil, nil, nil,
			container.NewBorder(c, nil, nil, nil, container.NewVScroll(ui.choiceContainer))))
}

func (ui *configInstallerUI) getChoiceSelector(onChange func(choice string)) fyne.CanvasObject {
	possible := make([]string, len(ui.currentConfig.Choices))
	for j, c := range ui.currentConfig.Choices {
		possible[j] = c.Name
	}

	st := ui.mod.ConfigSelectionType
	if st == mods.Auto {
		st = mods.Radio
		if len(ui.currentConfig.Choices) > 3 {
			st = mods.Select
		}
	}

	if st == mods.Radio {
		return widget.NewRadioGroup(possible, onChange)
	}
	return widget.NewSelect(possible, onChange)
}

func (ui *configInstallerUI) drawChoiceInfo() {
	c := container.NewVBox(
		widget.NewLabelWithStyle(ui.currentChoice.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	if ui.currentChoice.Description != "" {
		c.Add(widget.NewRichTextFromMarkdown(ui.currentChoice.Description))
	}
	if img := ui.currentChoice.Preview.Get(); img != nil {
		c = container.NewBorder(img, nil, nil, nil, c)
	}
	ui.choiceContainer.Add(c)
}

func (ui *configInstallerUI) popToInstall() {
	l := len(ui.toInstall) - 1
	if l < 0 {
		return
	}
	ui.toInstall[l] = nil
	ui.toInstall = ui.toInstall[:l]
}

func (ui *configInstallerUI) popChoice() (c *mods.Configuration) {
	l := len(ui.prevConfigs) - 1
	if l < 0 {
		return nil
	}
	c = ui.prevConfigs[l]
	ui.prevConfigs[l] = nil
	ui.prevConfigs = ui.prevConfigs[:l]
	return
}
