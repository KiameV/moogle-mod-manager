package custom_widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
)

type Callbacks struct {
	GetItemKey    func(item interface{}) string
	GetItemFields func(item interface{}) []string
	OnEditItem    func(item interface{})
}

type DynamicList struct {
	list          *fyne.Container
	Items         []interface{}
	callbacks     Callbacks
	confirmDelete bool
}

func NewDynamicList(callbacks Callbacks, confirmDelete bool) *DynamicList {
	return &DynamicList{
		list:          container.NewVBox(),
		callbacks:     callbacks,
		confirmDelete: confirmDelete,
	}
}

func (l *DynamicList) AddItem(item interface{}) {
	l.Items = append(l.Items, item)
	l.list.Objects = append(l.list.Objects, l.createRow(item))
}

func (l *DynamicList) createRow(item interface{}) fyne.CanvasObject {
	return container.NewHBox(
		widget.NewLabel(l.callbacks.GetItemKey(item)),
		widget.NewToolbar(
			// Edit
			newAction(item, theme.DocumentCreateIcon(), func(item interface{}) {
				l.callbacks.OnEditItem(item)
			}),
			// Remove
			newAction(item, theme.ContentRemoveIcon(), func(item interface{}) {
				l.removeItem(item)
			})),
	)
}

func (l *DynamicList) Draw() fyne.CanvasObject {
	return l.list
}

func (l *DynamicList) Refresh() {
	for i, item := range l.Items {
		l.list.Objects[i] = l.createRow(item)
	}
}

func (l *DynamicList) Reset() {
	l.Items = make([]interface{}, 0)
	l.list.Objects = make([]fyne.CanvasObject, 0)
}

type Action struct {
	*widget.ToolbarAction
}

func newAction(item interface{}, icon fyne.Resource, onActivated func(item interface{})) *Action {
	return &Action{
		ToolbarAction: widget.NewToolbarAction(icon, func() { onActivated(item) }),
	}
}

func (l *DynamicList) removeItem(item interface{}) {
	if l.confirmDelete {
		dialog.NewConfirm("Delete Item?", "Are you sure you want to delete this item?", func(ok bool) {
			if ok {
				l.removeFromList(item)
			}
		}, ui.Window).Show()
	} else {
		l.removeFromList(item)
	}
}

func (l *DynamicList) removeFromList(item interface{}) {
	for i, v := range l.Items {
		if item == v {
			l.Items = append(l.Items[:i], l.Items[i+1:]...)
			l.list.Objects = append(l.list.Objects[:i], l.list.Objects[i+1:]...)
			break
		}
	}
}

func (l *DynamicList) Clear() {
	l.Items = nil
	l.list.Objects = nil
}
