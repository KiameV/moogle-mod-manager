package entry

import (
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type SelectFormEntry struct {
	Entry    *widget.Select
	fi       *widget.FormItem
	selected string
	bind     binding.String
}

func NewSelectFormEntry(key string, value any, possible []string) Entry[string] {
	e := &SelectFormEntry{
		selected: value.(string),
	}
	e.bind = binding.BindString(&e.selected)
	e.bind.AddListener(e)
	e.Entry = widget.NewSelect(possible, func(s string) {
		_ = e.bind.Set(s)
	})
	e.fi = widget.NewFormItem(key, e.Entry)
	return e
}

func (e *SelectFormEntry) Enable(enable bool) {
	if enable {
		e.Entry.Enable()
	} else {
		e.Entry.Disable()
	}
}

func (e *SelectFormEntry) Binding() binding.DataItem {
	return e.bind
}

func (e *SelectFormEntry) Set(value string) {
	_ = e.bind.Set(value)
}

func (e *SelectFormEntry) Value() string {
	return e.selected
}

func (e *SelectFormEntry) DataChanged() {
	if e != nil && e.bind != nil {
		if v, err := e.bind.Get(); err == nil {
			if e.Entry != nil && e.Entry.Selected != v {
				e.Entry.Selected = v
			}
		}
	}
}

func (e *SelectFormEntry) FormItem() *widget.FormItem {
	return e.fi
}

func (e *SelectFormEntry) Clear() {
	var s []string
	e.Entry.Options = s
	e.Entry.Selected = ""
}
