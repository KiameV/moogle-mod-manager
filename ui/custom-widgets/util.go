package custom_widgets

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/kiamev/moogle-mod-manager/mods"
)

func GetValueFromDataItem(di binding.DataItem) (result interface{}, ok bool) {
	switch u := di.(type) {
	case binding.Untyped:
		if i, err := u.Get(); err == nil {
			switch v := i.(type) {
			case binding.Untyped:
				result, err = v.Get()
				ok = err == nil
				return
			case *mods.Mod:
				result = v
				ok = true
				return
			case interface{}:
				result = v
				ok = true
				return
			}
		}
	}
	return nil, false
}
