package entry

import (
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type selectFormEntry struct {
	entry    *widget.Select
	fi       *widget.FormItem
	selected string
	bind     binding.String
}

func newSelectFormEntry(key string, value any, possible []string) Entry[string] {
	e := &selectFormEntry{
		selected: value.(string),
	}
	e.bind = binding.BindString(&e.selected)
	e.bind.AddListener(e)
	e.entry = widget.NewSelect(possible, func(s string) {
		_ = e.bind.Set(s)
	})
	e.fi = widget.NewFormItem(key, e.entry)
	return e
}

func (e *selectFormEntry) Enable(enable bool) {
	if enable {
		e.entry.Enable()
	} else {
		e.entry.Disable()
	}
}

func (e *selectFormEntry) Binding() binding.DataItem {
	return e.bind
}

func (e *selectFormEntry) Set(value string) {
	_ = e.bind.Set(value)
}

func (e *selectFormEntry) Value() string {
	return e.selected
}

func (e *selectFormEntry) DataChanged() {
	if v, err := e.bind.Get(); err == nil {
		if e.entry.Selected != v {
			e.entry.Selected = v
		}
	}
}

func (e *selectFormEntry) FormItem() *widget.FormItem {
	return e.fi
}
