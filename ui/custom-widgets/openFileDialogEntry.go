package custom_widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/ncruces/zenity"
	"strings"
)

type OpenFileDialogHandler interface {
	widget.ToolbarItem
	Handle()
	Get() string
	SetAction(a *widget.ToolbarAction)
}

type OpenFileDialogContainer struct {
	*fyne.Container
	OpenFileDialogHandler OpenFileDialogHandler
}

type OpenDirDialog struct {
	*widget.ToolbarAction
	BaseDir    binding.String
	Value      binding.String
	IsRelative bool
}

func (o *OpenDirDialog) Get() string {
	s, _ := o.Value.Get()
	return s
}

func (o *OpenDirDialog) SetAction(a *widget.ToolbarAction) { o.ToolbarAction = a }

func (o *OpenDirDialog) Handle() {
	dir := config.PWD
	if o.BaseDir != nil {
		dir, _ = o.BaseDir.Get()
	}
	s, err := zenity.SelectFile(
		zenity.Title("Select file"),
		zenity.Filename(dir),
		zenity.Directory())
	if err == nil {
		if o.IsRelative {
			s = strings.ReplaceAll(s, dir, "")
			s = strings.ReplaceAll(s, "\\", "/")
			if len(s) == 0 || (len(s) > 0 && s[0] == '/') {
				s = "." + s
			}
		}
		_ = o.Value.Set(s)
	}
}

type OpenFileDialog struct {
	*widget.ToolbarAction
	BaseDir    binding.String
	Value      binding.String
	IsRelative bool
}

func (o *OpenFileDialog) Get() string {
	s, _ := o.Value.Get()
	return s
}

func (o *OpenFileDialog) SetAction(a *widget.ToolbarAction) { o.ToolbarAction = a }

func (o *OpenFileDialog) Handle() {
	dir := config.PWD
	if o.BaseDir != nil {
		dir, _ = o.BaseDir.Get()
	}
	s, err := zenity.SelectFile(
		zenity.Title("Select file"),
		zenity.Filename(dir),
		zenity.FileFilter{
			Name:     "All files",
			Patterns: []string{"*"},
		})
	if err == nil {
		if o.IsRelative {
			s = strings.ReplaceAll(s, dir, "")
			s = strings.ReplaceAll(s, "\\", "/")
			if len(s) > 0 && s[0] == '/' {
				s = "." + s
			}
		}
		_ = o.Value.Set(s)
	}
}
