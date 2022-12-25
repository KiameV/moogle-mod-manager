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
	prevConfigs     []*mods.Configuration
	choices         []*mods.Choice
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
	ui.choices = make([]*mods.Choice, 0)
	ui.choiceContainer.RemoveAll()
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
			ui.choices = append(ui.choices, ui.currentChoice)
			if ui.currentChoice.NextConfigurationName == nil {
				tis, err := mods.NewToInstallForMod(ui.mod.ModKind.Kind, ui.mod, ui.uniqueToInstall())
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
			if len(ui.choices) > 0 {
				ui.choices = ui.choices[:len(ui.choices)-1]
			}
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

func (ui *configInstallerUI) uniqueToInstall() []*mods.DownloadFiles {
	var (
		l  = make(map[string]*mods.DownloadFiles)
		df *mods.DownloadFiles
	)
	for _, df = range ui.mod.AlwaysDownload {
		l[df.DownloadName] = df
	}
	for _, c := range ui.choices {
		df = c.DownloadFiles
		if df != nil && df.DownloadName != "" {
			if len(df.Dirs) > 0 || len(df.Files) > 0 {
				if to, found := l[df.DownloadName]; !found {
					l[df.DownloadName] = df
				} else {
					l[df.DownloadName] = merge(to, df)
				}
			}
		}
	}
	result := make([]*mods.DownloadFiles, 0, len(l))
	for _, df = range l {
		result = append(result, df)
	}
	return result
}

func merge(df1 *mods.DownloadFiles, df2 *mods.DownloadFiles) *mods.DownloadFiles {
	var (
		m     = make(map[string]bool)
		dirs  = make([]*mods.ModDir, 0, len(df1.Dirs)+len(df2.Dirs))
		files = make([]*mods.ModFile, 0, len(df1.Files)+len(df2.Files))
	)
	for _, d := range df2.Dirs {
		m[d.To] = true
		dirs = append(dirs, d)
	}
	for _, d := range df1.Dirs {
		if !m[d.To] {
			dirs = append(dirs, d)
		}
	}

	m = make(map[string]bool)
	for _, f := range df2.Files {
		m[f.To] = true
		files = append(files, f)
	}
	for _, f := range df1.Files {
		if !m[f.To] {
			files = append(files, f)
		}
	}
	return &mods.DownloadFiles{
		DownloadName: df1.DownloadName,
		Dirs:         dirs,
		Files:        files,
	}
}
