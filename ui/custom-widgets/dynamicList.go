package custom_widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Callbacks struct {
	GetItemKey    func(item interface{}) string
	GetItemFields func(item interface{}) []string
	OnEditItem    func(item interface{}, done func(result interface{}))
}

type DynamicList struct {
	list      *fyne.Container
	Items     []interface{}
	callbacks Callbacks
}

func NewDynamicList(callbacks Callbacks) *DynamicList {
	l := &DynamicList{
		list:      container.NewVBox(),
		callbacks: callbacks,
	}
	return l
}

func (l *DynamicList) AddItem(item interface{}) {
	l.Items = append(l.Items, item)
	l.createRow(item)
}

func (l *DynamicList) createRow(item interface{}) {
	r := container.NewHBox(
		widget.NewLabel(l.callbacks.GetItemKey(item)),
		widget.NewToolbar(
			// Edit
			newAction(item, theme.DocumentCreateIcon(), func(item interface{}) {
				l.callbacks.OnEditItem(item, func(result interface{}) {
					l.removeItem(item)
					l.AddItem(result)
				})
			}),
			// Remove
			newAction(item, theme.ContentRemoveIcon(), func(item interface{}) {
				l.removeItem(item)
			})),
	)
	l.list.Objects = append(l.list.Objects, r)
}

func (l *DynamicList) Draw() fyne.CanvasObject {
	return l.list
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
	for i, v := range l.Items {
		if item == v {
			l.Items = append(l.Items[:i], l.Items[i+1:]...)
			l.list.Objects = append(l.list.Objects[:i], l.list.Objects[i+1:]...)
			break
		}
	}
}
