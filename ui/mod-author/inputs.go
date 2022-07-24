package mod_author

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
	"strconv"
	"strings"
)

const baseDirKey = "__base_dir"

func newEntryManager() *entryManager {
	return &entryManager{
		entries: make(map[string]fyne.CanvasObject),
	}
}

type entryManager struct {
	entries map[string]fyne.CanvasObject
}

func (m *entryManager) getBool(key string) bool {
	e, ok := m.entries[key]
	if !ok {
		return false
	}
	var c *widget.Check
	if c, ok = e.(*widget.Check); !ok {
		return false
	}
	return c.Checked
}

func (m *entryManager) getString(key string) string {
	e, ok := m.entries[key]
	if !ok {
		return ""
	}
	switch t := e.(type) {
	case *widget.Entry:
		return t.Text
	case *widget.Select:
		return t.Selected
	case *widget.Check:
		return fmt.Sprintf("%v", t.Checked)
	case *ofd:
		return t.ofh.Get()
	}
	return ""
}

func (m *entryManager) getInt(key string) (i int) {
	if s := m.getString(key); s != "" {
		if j, err := strconv.ParseUint(s, 10, 32); err == nil {
			i = int(j)
		}
	}
	return
}

func (m *entryManager) getStrings(key string, split string) []string {
	s := m.getString(key)
	if s != "" {
		sl := strings.Split(s, split)
		for i, j := range sl {
			sl[i] = strings.TrimSpace(j)
		}
		return sl
	}
	return nil
}

func (m *entryManager) getFormItem(name string) *widget.FormItem {
	key := name
	e, _ := m.entries[key]
	return widget.NewFormItem(name, e)
}

func (m *entryManager) getBaseDirFormItem(name string) *widget.FormItem {
	e, _ := m.entries[baseDirKey]
	return widget.NewFormItem(name, e.(*ofd).Container)
}

func (m *entryManager) getFileDialog(name string) *widget.FormItem {
	e, _ := m.entries[name]
	return widget.NewFormItem(name, e.(*ofd).Container)
}

func (m *entryManager) createBaseDir(baseDir binding.String) {
	if _, ok := m.entries[baseDirKey]; !ok {
		o := &openDir{BaseDir: baseDir, Value: baseDir}
		o.ToolbarAction = widget.NewToolbarAction(theme.FolderOpenIcon(), o.Handle)
		m.entries[baseDirKey] = &ofd{
			Container: container.NewBorder(nil, nil, nil, widget.NewToolbar(o), widget.NewEntryWithData(baseDir)),
			ofh:       o,
		}
	}
}

func (m *entryManager) createFileDialog(key string, value string, baseDir binding.String, isDir bool, isRelative bool) {
	e, ok := m.entries[key]
	if !ok {
		b := binding.NewString()
		var o openFileHandler
		if isDir {
			o = &openDir{
				IsRelative: isRelative,
				Value:      b,
				BaseDir:    baseDir,
			}
		} else {
			o = &openFile{
				IsRelative: isRelative,
				Value:      b,
				BaseDir:    baseDir,
			}
		}
		o.SetAction(widget.NewToolbarAction(theme.FolderOpenIcon(), o.Handle))
		e = &ofd{
			Container: container.NewBorder(nil, nil, nil, widget.NewToolbar(o), widget.NewEntryWithData(b)),
			ofh:       o,
		}
		m.entries[key] = e
	}
	switch t := e.(*ofd).ofh.(type) {
	case *openDir:
		_ = t.Value.Set(value)
	case *openFile:
		_ = t.Value.Set(value)
	}
}

func (m *entryManager) createFormItem(key string, value string) {
	e, ok := m.entries[key]
	if !ok {
		e = widget.NewEntry()
		m.entries[key] = e
	}
	e.(*widget.Entry).SetText(value)
}

func (m *entryManager) createFormSelect(key string, possible []string, value string) {
	e := widget.NewSelect(possible, func(string) {})
	m.entries[key] = e
	e.SetSelected(value)
}

func (m *entryManager) createFormMultiLine(key string, value string) {
	e, ok := m.entries[key]
	if !ok {
		e = widget.NewMultiLineEntry()
		m.entries[key] = e
	}
	e.(*widget.Entry).SetText(value)
}

func (m *entryManager) createFormBool(key string, value bool) {
	e, ok := m.entries[key]
	if !ok {
		e = widget.NewCheck(key, func(bool) {})
		m.entries[key] = e
	}
	e.(*widget.Check).SetChecked(value)
}

type openFileHandler interface {
	widget.ToolbarItem
	Handle()
	Get() string
	SetAction(a *widget.ToolbarAction)
}

type ofd struct {
	*fyne.Container
	ofh openFileHandler
}

type openDir struct {
	*widget.ToolbarAction
	BaseDir    binding.String
	Value      binding.String
	IsRelative bool
}

func (o *openDir) Get() string {
	s, _ := o.Value.Get()
	return s
}

func (o *openDir) SetAction(a *widget.ToolbarAction) { o.ToolbarAction = a }

func (o *openDir) Handle() {
	dir, _ := o.BaseDir.Get()
	s, err := zenity.SelectFile(
		zenity.Title("Select file"),
		zenity.Filename(dir),
		zenity.Directory())
	if err == nil {
		if o.IsRelative {
			dir, _ = o.BaseDir.Get()
			s = strings.ReplaceAll(s, dir, "")
			s = strings.ReplaceAll(s, "\\", "/")
			if len(s) == 0 || (len(s) > 0 && s[0] == '/') {
				s = "." + s
			}
		}
		_ = o.Value.Set(s)
	}
}

type openFile struct {
	*widget.ToolbarAction
	BaseDir    binding.String
	Value      binding.String
	IsRelative bool
}

func (o *openFile) Get() string {
	s, _ := o.Value.Get()
	return s
}

func (o *openFile) SetAction(a *widget.ToolbarAction) { o.ToolbarAction = a }

func (o *openFile) Handle() {
	dir, _ := o.BaseDir.Get()
	s, err := zenity.SelectFile(
		zenity.Title("Select file"),
		zenity.Filename(dir),
		zenity.FileFilter{
			Name:     "All files",
			Patterns: []string{"*"},
		})
	if err == nil {
		if o.IsRelative {
			dir, _ = o.BaseDir.Get()
			s = strings.ReplaceAll(s, dir, "")
			s = strings.ReplaceAll(s, "\\", "/")
			if len(s) > 0 && s[0] == '/' {
				s = "." + s
			}
		}
		_ = o.Value.Set(s)
	}
}
