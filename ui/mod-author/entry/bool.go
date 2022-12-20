package entry

import (
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type boolFormEntry struct {
	entry   *widget.Check
	fi      *widget.FormItem
	bind    binding.Bool
	checked bool
}

func newBoolFormEntry(key string, value any) Entry[bool] {
	e := &boolFormEntry{
		checked: value.(bool),
	}
	e.bind = binding.BindBool(&e.checked)
	e.entry = widget.NewCheckWithData(key, e.bind)
	e.fi = widget.NewFormItem(key, e.entry)
	return e
}

func (e *boolFormEntry) Binding() binding.DataItem {
	return e.bind
}

func (e *boolFormEntry) Enable(enable bool) {
	if enable {
		e.entry.Enable()
	} else {
		e.entry.Disable()
	}
}

func (e *boolFormEntry) Set(value bool) {
	_ = e.bind.Set(value)
}

func (e *boolFormEntry) Value() bool {
	return e.checked
}

func (e *boolFormEntry) FormItem() *widget.FormItem {
	return e.fi
}
