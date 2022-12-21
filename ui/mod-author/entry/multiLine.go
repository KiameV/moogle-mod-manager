package entry

import (
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type multiLineFormEntry struct {
	entry *widget.Entry
	fi    *widget.FormItem
	bind  binding.String
	text  string
}

func newMultiLineFormEntry(key string, value any) Entry[string] {
	e := &multiLineFormEntry{
		text: value.(string),
	}
	e.bind = binding.BindString(&e.text)
	e.entry = widget.NewMultiLineEntry()
	e.fi = widget.NewFormItem(key, e.entry)
	return e
}

func (e *multiLineFormEntry) Enable(enable bool) {
	if enable {
		e.entry.Enable()
	} else {
		e.entry.Disable()
	}
}

func (e *multiLineFormEntry) Binding() binding.DataItem {
	return e.bind
}

func (e *multiLineFormEntry) Set(value string) {
	_ = e.bind.Set(value)
}

func (e *multiLineFormEntry) Value() string {
	return e.text
}

func (e *multiLineFormEntry) FormItem() *widget.FormItem {
	return e.fi
}
