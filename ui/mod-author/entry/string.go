package entry

import (
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type stringFormEntry struct {
	entry *widget.Entry
	fi    *widget.FormItem
	bind  binding.String
	text  string
}

func newStringFormEntry(key string, value any) Entry[string] {
	e := &stringFormEntry{
		text: value.(string),
	}
	e.bind = binding.BindString(&e.text)
	e.entry = widget.NewEntryWithData(e.bind)
	e.fi = widget.NewFormItem(key, e.entry)
	return e
}

func (e *stringFormEntry) Binding() binding.DataItem {
	return e.bind
}

func (e *stringFormEntry) Enable(enable bool) {
	if enable {
		e.entry.Enable()
	} else {
		e.entry.Disable()
	}
}

func (e *stringFormEntry) Set(value string) {
	_ = e.bind.Set(value)
}

func (e *stringFormEntry) Value() string {
	return e.text
}

func (e *stringFormEntry) FormItem() *widget.FormItem {
	return e.fi
}
