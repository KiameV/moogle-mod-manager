package mod_author

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"strings"
)

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
	case *widget.SelectEntry:
		return t.Text
	case *widget.Check:
		return fmt.Sprintf("%v", t.Checked)
	}
	return ""
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

func (m *entryManager) getFormItem(name string, customKey ...string) *widget.FormItem {
	key := name
	if len(customKey) > 0 {
		key = customKey[0]
	}
	e, _ := m.entries[key]
	return widget.NewFormItem(name, e)
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
	e, ok := m.entries[key]
	if !ok {
		e = widget.NewSelectEntry(possible)
		m.entries[key] = e
	}
	e.(*widget.SelectEntry).SetText(value)
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
