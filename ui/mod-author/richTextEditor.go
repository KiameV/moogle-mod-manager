package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type richTextEditor struct {
	input   binding.String
	preview *widget.RichText
}

func newRichTextEditor() *richTextEditor {
	e := &richTextEditor{
		input:   binding.NewString(),
		preview: widget.NewRichTextWithText(""),
	}
	e.input.AddListener(e)
	return e
}

func (e *richTextEditor) Draw() fyne.CanvasObject {
	entry := widget.NewMultiLineEntry()
	entry.Bind(e.input)
	entry.Wrapping = fyne.TextWrapWord
	e.preview.Wrapping = fyne.TextWrapWord
	return container.NewVSplit(
		container.NewScroll(entry),
		container.NewVScroll(e.preview))
}

func (e *richTextEditor) DataChanged() {
	e.preview.ParseMarkdown(e.String())
}

func (e *richTextEditor) SetText(s string) {
	e.input.Set(s)
}

func (e *richTextEditor) String() string {
	s, _ := e.input.Get()
	return s
}
