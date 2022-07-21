package config_installer

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"strings"
)

type ConfigInstaller interface {
	state.Screen
	Setup(mod *mods.Mod, isSandbox bool) error
}

func New() ConfigInstaller {
	return &configInstallerUI{
		choiceDesc: container.NewVBox(),
	}
}

type configInstallerUI struct {
	mod         *mods.Mod
	isSandbox   bool
	toInstall   []*mods.DownloadFiles
	prevConfigs []*mods.Configuration
	choiceDesc  *fyne.Container

	currentConfig *mods.Configuration
	currentChoice *mods.Choice
}

func (i *configInstallerUI) Setup(mod *mods.Mod, isSandbox bool) error {
	if len(mod.Configurations) == 0 || len(mod.Configurations[0].Choices) == 0 {
		return fmt.Errorf("no configurations for %s", mod.Name)
	}
	i.mod = mod
	i.currentConfig = mod.Configurations[0]
	i.isSandbox = isSandbox
	return nil
}

func (i *configInstallerUI) Draw(w fyne.Window) {
	c := container.NewVBox(
		widget.NewLabelWithStyle(i.currentConfig.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)
	if i.currentConfig.Preview != "" {
		c.Add(canvas.NewImageFromURI(storage.NewFileURI(i.currentConfig.Preview)))
	}
	c.Add(widget.NewRichTextFromMarkdown(i.currentConfig.Description))
	c.Add(i.getChoiceSelector(func(name string) {
		for _, i.currentChoice = range i.currentConfig.Choices {
			if i.currentChoice.Name == name {
				i.choiceDesc.RemoveAll()
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
				dlfs := i.compileDownloadFiles()
				if i.isSandbox {
					sb := strings.Builder{}
					for dl, dlf := range dlfs {
						sb.WriteString(fmt.Sprintf("Download: %s\n\n", dl.Name))
						sb.WriteString("  Sources:\n\n")
						for _, s := range dl.Sources {
							sb.WriteString(fmt.Sprintf("  - %s\n\n", s))
						}
						sb.WriteString("  Files:\n\n")
						for _, f := range dlf.Files {
							sb.WriteString(fmt.Sprintf("  - %s -> %s\n\n", f.From, f.To))
						}
						sb.WriteString("  Dirs:\n\n")
						for _, dir := range dlf.Dirs {
							sb.WriteString(fmt.Sprintf("  - %s -> %s | Recursive %v\n\n", dir.From, dir.To, dir.Recursive))
						}
						break
					}
					dialog.ShowCustom("Downloads and File/Dir Copies", "ok", widget.NewRichTextFromMarkdown(sb.String()), state.Window)
					state.ShowPreviousScreen()
				} else {
					// TODO
				}
			} else {
				for _, i.currentConfig = range i.mod.Configurations {
					if i.currentConfig.Name == *i.currentChoice.NextConfigurationName {
						break
					}
				}
				i.currentChoice = nil
				i.choiceDesc.RemoveAll()
				i.Draw(w)
			}
		}))
	if len(i.prevConfigs) > 0 {
		buttons.Add(widget.NewButton("Back", func() {
			i.popToInstall()
			i.currentConfig = i.popChoice()
			i.choiceDesc.RemoveAll()
			i.Draw(w)
		}))
	}
	c.Add(buttons)
	w.SetContent(container.NewBorder(c, nil, nil, nil, container.NewVScroll(i.choiceDesc)))
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
	c := container.NewVBox()
	c.Add(widget.NewSeparator())
	c.Add(widget.NewLabelWithStyle(i.currentChoice.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	if i.currentChoice.Description != "" {
		c.Add(widget.NewRichTextFromMarkdown(i.currentChoice.Description))
	}
	if i.currentChoice.Preview != "" {
		c.Add(canvas.NewImageFromURI(storage.NewFileURI(i.currentChoice.Preview)))
	}
	i.choiceDesc.Add(container.NewMax(c))
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

func (i *configInstallerUI) compileDownloadFiles() map[*mods.Download]*mods.DownloadFiles {
	dlf := make(map[*mods.Download]*mods.DownloadFiles)
	dl := make(map[string]*mods.Download)
	for _, d := range i.mod.Downloadables {
		dl[d.Name] = d
	}
	for _, ti := range i.toInstall {
		d := dl[ti.DownloadName]
		f, ok := dlf[d]
		if !ok {
			dlf[d] = &mods.DownloadFiles{DownloadName: ti.DownloadName}
		}
		for _, df := range ti.Files {
			f.Files = append(f.Files, df)
		}
		for _, dd := range ti.Dirs {
			f.Dirs = append(f.Dirs, dd)
		}
	}
	return dlf
}
