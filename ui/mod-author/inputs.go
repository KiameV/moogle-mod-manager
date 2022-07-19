package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"strings"
)

var entries = make(map[string]fyne.CanvasObject)

func getFormString(key string) string {
	e, ok := entries[key]
	if !ok {
		return ""
	}
	switch t := e.(type) {
	case *widget.Entry:
		return t.Text
	case *widget.SelectEntry:
		return t.Text
	}
	return ""
}

func getFormStrings(key string, split string) []string {
	s := getFormString(key)
	if s != "" {
		sl := strings.Split(s, split)
		for i, j := range sl {
			sl[i] = strings.TrimSpace(j)
		}
	}
	return nil
}

func getFormItem(name string, customKey ...string) *widget.FormItem {
	key := name
	if len(customKey) > 0 {
		key = customKey[0]
	}
	e, _ := entries[key]
	return widget.NewFormItem(name, e)
}

func setFormItem(key string, value string) {
	e, ok := entries[key]
	if !ok {
		e = widget.NewEntry()
		//e.Resize(fyne.Size{Width: 200, Height: e.MinSize().Height})
		entries[key] = e
	}
	e.(*widget.Entry).SetText(value)
}

func setFormSelect(key string, possible []string, value string) {
	e, ok := entries[key]
	if !ok {
		e = widget.NewSelectEntry(possible)
		//e.Resize(fyne.Size{Width: 200, Height: e.MinSize().Height})
		entries[key] = e
	}
	e.(*widget.SelectEntry).SetText(value)
}

func setFormMultiLine(key string, value string) {
	e, ok := entries[key]
	if !ok {
		e = widget.NewMultiLineEntry()
		//e.Resize(fyne.Size{Width: 200, Height: e.MinSize().Height})
		entries[key] = e
	}
	e.(*widget.Entry).SetText(value)
}
