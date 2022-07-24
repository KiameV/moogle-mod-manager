package custom_widgets

import "fyne.io/fyne/v2/data/binding"

func GetValueFromDataItem(di binding.DataItem) (result interface{}, ok bool) {
	switch u := di.(type) {
	case binding.Untyped:
		if i, err := u.Get(); err == nil {
			var ut binding.Untyped
			if ut, ok = i.(binding.Untyped); ok {
				result, err = ut.Get()
				return result, err == nil
			}
		}
	}
	return nil, false
}
