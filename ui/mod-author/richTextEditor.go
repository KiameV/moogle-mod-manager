package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type richTextEditor struct {
	s       *string
	input   binding.ExternalString
	preview *widget.RichText
}

func newRichTextEditor() *richTextEditor {
	s := ""
	e := &richTextEditor{
		s:       &s,
		preview: widget.NewRichTextWithText(""),
	}
	e.input = binding.BindString(e.s)
	e.input.AddListener(e)
	return e
}

func (e *richTextEditor) Draw() fyne.CanvasObject {
	entry := widget.NewMultiLineEntry()
	entry.Bind(e.input)
	return container.NewVSplit(
		container.NewScroll(entry),
		container.NewScroll(e.preview))
}

func (e *richTextEditor) DataChanged() {
	e.preview.ParseMarkdown(*e.s)
}

func (e *richTextEditor) SetText(s string) {
	*e.s = s
	e.preview.ParseMarkdown(s)
}

func (e *richTextEditor) String() string {
	return *e.s
}
