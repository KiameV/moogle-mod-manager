package mod_author

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
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
	case *cw.OpenFileDialogContainer:
		return t.OpenFileDialogHandler.Get()
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
	e, found := m.entries[key]
	if !found {
		panic(fmt.Sprintf("entry %s not found", key))
	}
	return widget.NewFormItem(name, e)
}

func (m *entryManager) getBaseDirFormItem(name string) *widget.FormItem {
	e := m.entries[baseDirKey]
	return widget.NewFormItem(name, e.(*cw.OpenFileDialogContainer).Container)
}

func (m *entryManager) getFileDialog(name string) *widget.FormItem {
	e := m.entries[name]
	return widget.NewFormItem(name, e.(*cw.OpenFileDialogContainer).Container)
}

func (m *entryManager) createBaseDir(baseDir binding.String) {
	if _, ok := m.entries[baseDirKey]; !ok {
		o := &cw.OpenDirDialog{BaseDir: baseDir, Value: baseDir}
		o.ToolbarAction = widget.NewToolbarAction(theme.FolderOpenIcon(), o.Handle)
		m.entries[baseDirKey] = &cw.OpenFileDialogContainer{
			Container:             container.NewBorder(nil, nil, nil, widget.NewToolbar(o), widget.NewEntryWithData(baseDir)),
			OpenFileDialogHandler: o,
		}
	}
}

func (m *entryManager) createFileDialog(key string, value string, baseDir binding.String, isDir bool, isRelative bool) {
	e, ok := m.entries[key]
	if !ok {
		b := binding.NewString()
		var o cw.OpenFileDialogHandler
		if isDir {
			o = &cw.OpenDirDialog{
				IsRelative: isRelative,
				Value:      b,
				BaseDir:    baseDir,
			}
		} else {
			o = &cw.OpenFileDialog{
				IsRelative: isRelative,
				Value:      b,
				BaseDir:    baseDir,
			}
		}
		o.SetAction(widget.NewToolbarAction(theme.FolderOpenIcon(), o.Handle))
		e = &cw.OpenFileDialogContainer{
			Container:             container.NewBorder(nil, nil, nil, widget.NewToolbar(o), widget.NewEntryWithData(b)),
			OpenFileDialogHandler: o,
		}
		m.entries[key] = e
	}
	switch t := e.(*cw.OpenFileDialogContainer).OpenFileDialogHandler.(type) {
	case *cw.OpenDirDialog:
		_ = t.Value.Set(value)
	case *cw.OpenFileDialog:
		_ = t.Value.Set(value)
	}
}

func (m *entryManager) newFormItem(key string, value ...string) *widget.FormItem {
	var v string
	if len(value) > 0 {
		v = value[0]
	}
	m.createFormItem(key, v)
	return m.getFormItem(key)
}

func (m *entryManager) createFormItem(key string, value string) {
	e, ok := m.entries[key]
	if !ok {
		e = widget.NewEntry()
		m.entries[key] = e
	}
	e.(*widget.Entry).SetText(value)
}

func (m *entryManager) createFormSelect(key string, possible []string, value string, onChange ...func(string)) {
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
