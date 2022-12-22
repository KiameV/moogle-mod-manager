package entry

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
)

const baseDirKey = "__base_dir"

type (
	Kind    byte
	em      map[string]any
	Manager interface {
		get(key string) (any, bool)
		set(key string, value any)
	}
	manager struct {
		entries em
	}
	Entry[T any] interface {
		Enable(bool)
		Binding() binding.DataItem
		Set(t T)
		Value() T
		FormItem() *widget.FormItem
	}
	valuer[T any] interface {
		Value() T
	}
)

const (
	_ Kind = iota
	KindBool
	KindString
	KindMultiLine
)

func NewManager() Manager {
	return &manager{entries: make(map[string]any)}
}

func (m *manager) get(key string) (any, bool) {
	e, ok := m.entries[key]
	return e, ok
}

func (m *manager) set(key string, value any) {
	m.entries[key] = value
}

func Value[T any](m Manager, key string) (t T) {
	e, ok := m.get(key)
	if ok {
		switch en := e.(type) {
		case Entry[T]:
			t = en.Value()
		default:
			panic("unknown type")
		}
	}
	return
}

func DialogValue(m Manager, key string) string {
	e, ok := m.get(key)
	if ok {
		switch en := e.(type) {
		case *cw.OpenFileDialogContainer:
			return en.OpenFileDialogHandler.Get()
		default:
			panic("unknown type")
		}
	}
	return ""
}

func NewEntry[T any](m Manager, kind Kind, key string, value T) Entry[T] {
	e, found := m.get(key)
	if !found {
		switch kind {
		case KindBool:
			e = newBoolFormEntry(key, value)
		case KindString:
			e = NewStringFormEntry(key, value)
		case KindMultiLine:
			e = newMultiLineFormEntry(key, value)
		default:
			panic(fmt.Sprintf("unknown entry kind %d", kind))
		}
		m.set(key, e)
	} else {
		e.(Entry[T]).Set(value)
	}
	return e.(Entry[T])
}

func NewSelectEntry(m Manager, key string, value string, possible []string) Entry[string] {
	e, found := m.get(key)
	if !found {
		e = NewSelectFormEntry(key, value, possible)
		m.set(key, e)
	} else {
		e.(*SelectFormEntry).Entry.Options = possible
		e.(*SelectFormEntry).selected = value
	}
	return e.(Entry[string])
}

func GetEntry[T any](m Manager, key string) Entry[T] {
	e, ok := m.get(key)
	if !ok {
		return nil
	}
	return e.(Entry[T])
}

func FormItem[T any](m Manager, key string, values ...T) *widget.FormItem {
	e := GetEntry[T](m, key)
	if len(values) > 0 {
		e.Set(values[0])
	}
	return e.FormItem()
}

func GetBaseDirFormItem(m Manager, name string) *widget.FormItem {
	e, _ := m.get(baseDirKey)
	return widget.NewFormItem(name, e.(*cw.OpenFileDialogContainer).Container)
}

func GetFileDialog(m Manager, name string) *widget.FormItem {
	e, _ := m.get(name)
	return widget.NewFormItem(name, e.(*cw.OpenFileDialogContainer).Container)
}

func CreateBaseDir(m Manager, baseDir binding.String) {
	if _, ok := m.get(baseDirKey); !ok {
		o := &cw.OpenDirDialog{BaseDir: baseDir, Value: baseDir}
		o.ToolbarAction = widget.NewToolbarAction(theme.FolderOpenIcon(), o.Handle)
		m.set(baseDirKey, &cw.OpenFileDialogContainer{
			Container:             container.NewBorder(nil, nil, nil, widget.NewToolbar(o), widget.NewEntryWithData(baseDir)),
			OpenFileDialogHandler: o,
		})
	}
}

func CreateFileDialog(m Manager, key string, value string, baseDir binding.String, isDir bool, isRelative bool) {
	e, ok := m.get(key)
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
		m.set(key, e)
	}
	switch t := e.(*cw.OpenFileDialogContainer).OpenFileDialogHandler.(type) {
	case *cw.OpenDirDialog:
		_ = t.Value.Set(value)
	case *cw.OpenFileDialog:
		_ = t.Value.Set(value)
	}
}

func FormItemFileDialog(m Manager, key string) *widget.FormItem {
	e, _ := m.get(key)
	return widget.NewFormItem(key, e.(fyne.CanvasObject))
}
