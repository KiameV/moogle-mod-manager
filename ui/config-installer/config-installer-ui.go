package config_installer

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state"
)

type ConfigInstaller interface {
	state.Screen
	Setup(mod *mods.Mod, baseDir string, callback func([]*mods.ToInstall) error) error
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
	callback        func([]*mods.ToInstall) error

	currentConfig *mods.Configuration
	currentChoice *mods.Choice
}

func (i *configInstallerUI) OnClose() {

}

func (i *configInstallerUI) Setup(mod *mods.Mod, baseDir string, callback func([]*mods.ToInstall) error) error {
	if len(mod.Configurations) == 0 || len(mod.Configurations[0].Choices) == 0 {
		return fmt.Errorf("no configurations for %s", mod.Name)
	}
	i.mod = mod
	for _, i.currentConfig = range mod.Configurations {
		if i.currentConfig.Root {
			break
		}
	}
	if i.currentConfig == nil || !i.currentConfig.Root {
		return errors.New("could not find root configuration")
	}
	i.prevConfigs = make([]*mods.Configuration, 0)
	i.baseDir = baseDir
	i.callback = callback
	i.toInstall = make([]*mods.DownloadFiles, 0)
	i.choiceContainer.RemoveAll()
	for _, dl := range mod.AlwaysDownload {
		i.toInstall = append(i.toInstall, dl)
	}
	return nil
}

func (i *configInstallerUI) Draw(w fyne.Window) {
	state.SetBaseDir(i.baseDir)
	c := container.NewVBox(
		widget.NewLabelWithStyle(i.currentConfig.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewRichTextFromMarkdown(i.currentConfig.Description),
		i.getChoiceSelector(func(name string) {
			for _, i.currentChoice = range i.currentConfig.Choices {
				if i.currentChoice.Name == name {
					i.choiceContainer.RemoveAll()
					i.drawChoiceInfo()
					break
				}
			}
		}))
	buttons := container.NewHBox(
		widget.NewButton("Select", func() {
			if i.currentChoice == nil {
				return
			}
			i.prevConfigs = append(i.prevConfigs, i.currentConfig)
			i.toInstall = append(i.toInstall, i.currentChoice.DownloadFiles)
			if i.currentChoice.NextConfigurationName == nil {
				tis, err := mods.NewToInstallForMod(i.mod, i.toInstall)
				if err != nil {
					dialog.ShowError(err, w)
					state.ShowPreviousScreen()
					return
				}
				if err = i.callback(tis); err != nil {
					dialog.ShowError(err, w)
					return
				}
				state.ShowPreviousScreen()
			} else {
				for _, i.currentConfig = range i.mod.Configurations {
					if i.currentConfig.Name == *i.currentChoice.NextConfigurationName {
						break
					}
				}
				i.currentChoice = nil
				i.choiceContainer.RemoveAll()
				i.Draw(w)
			}
		}))
	if len(i.prevConfigs) > 0 {
		buttons.Add(widget.NewButton("Back", func() {
			i.popToInstall()
			i.currentConfig = i.popChoice()
			i.choiceContainer.RemoveAll()
			i.Draw(w)
		}))
	}
	c.Add(buttons)
	if img := i.currentConfig.Preview.Get(); img != nil {
		c = container.NewBorder(img, nil, nil, nil, c)
	}
	w.SetContent(container.NewBorder(c, nil, nil, nil, container.NewVScroll(i.choiceContainer)))
}

func (i *configInstallerUI) getChoiceSelector(onChange func(choice string)) fyne.CanvasObject {
	possible := make([]string, len(i.currentConfig.Choices))
	for j, c := range i.currentConfig.Choices {
		possible[j] = c.Name
	}

	st := i.mod.ConfigSelectionType
	if st == mods.Auto {
		st = mods.Radio
		if len(i.currentConfig.Choices) > 3 {
			st = mods.Select
		}
	}

	if st == mods.Radio {
		return widget.NewRadioGroup(possible, onChange)
	}
	return widget.NewSelect(possible, onChange)
}

func (i *configInstallerUI) drawChoiceInfo() {
	c := container.NewVBox(
		widget.NewLabelWithStyle(i.currentChoice.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	if i.currentChoice.Description != "" {
		c.Add(widget.NewRichTextFromMarkdown(i.currentChoice.Description))
	}
	if img := i.currentChoice.Preview.Get(); img != nil {
		c = container.NewBorder(img, nil, nil, nil, c)
	}
	i.choiceContainer.Add(c)
}

func (i *configInstallerUI) popToInstall() {
	l := len(i.toInstall) - 1
	if l < 0 {
		return
	}
	i.toInstall[l] = nil
	i.toInstall = i.toInstall[:l]
}

func (i *configInstallerUI) popChoice() (c *mods.Configuration) {
	l := len(i.prevConfigs) - 1
	if l < 0 {
		return nil
	}
	c = i.prevConfigs[l]
	i.prevConfigs[l] = nil
	i.prevConfigs = i.prevConfigs[:l]
	return
}
