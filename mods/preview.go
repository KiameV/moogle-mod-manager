package mods

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/cache"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"os"
	"path/filepath"
)

type Preview struct {
	Url   *string       `json:"Url,omitempty" xml:"Url,omitempty"`
	Local *string       `json:"Local,omitempty" xml:"Local,omitempty"`
	img   *canvas.Image `json:"-" xml:"-"`
}

func (p *Preview) Get() *canvas.Image {
	if p == nil {
		return nil
	}
	if p.img == nil {
		p.img = p.GetUncachedImage()
	}
	return p.img
}

func (p *Preview) GetUncachedImage() (img *canvas.Image) {
	var (
		r   fyne.Resource
		err error
	)
	if p.Local != nil {
		f := filepath.Join(state.GetBaseDir(), *p.Local)
		if _, err = os.Stat(f); err == nil {
			r, err = fyne.LoadResourceFromPath(f)
		}
	}
	if r == nil && p.Url != nil {
		if r, err = cache.GetImage(*p.Url); err != nil {
			r, err = fyne.LoadResourceFromURLString(*p.Url)
		}
	}
	if r == nil || err != nil {
		return nil
	}
	img = canvas.NewImageFromResource(r)
	img.SetMinSize(fyne.Size{Width: float32(300), Height: float32(300)})
	img.FillMode = canvas.ImageFillContain
	return
}

func (p *Preview) GetAsButton(onClick func()) *fyne.Container {
	i := p.Get()
	if i == nil {
		return nil
	}
	return container.NewMax(i, widget.NewButton("", onClick))
}

func (p *Preview) GetAsEnlargeOnClick() *fyne.Container {
	i := p.Get()
	if i == nil {
		return nil
	}
	return container.NewBorder(nil, container.NewCenter(widget.NewButton("Enlarge", func() {
		d := dialog.NewCustom("", "Close", p.GetUncachedImage(), ui.ActiveWindow())
		d.Resize(config.Get().Size())
		d.Show()
	})), nil, nil, i)
}

func (p *Preview) GetAsImageGallery(index int, previews []*Preview, enlarge bool) *fyne.Container {
	var (
		c    = container.NewMax()
		left = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
			index = p.decrementIndex(index, len(previews))
			if img := previews[index].GetUncachedImage(); img != nil {
				c.Objects = nil
				c.Add(img)
			} else {
				index = p.incrementIndex(index, len(previews))
			}
		})
		right = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
			index = p.incrementIndex(index, len(previews))
			if img := previews[index].GetUncachedImage(); img != nil {
				c.Objects = nil
				c.Add(img)
			} else {
				index = p.decrementIndex(index, len(previews))
			}
		})
	)

	if img := previews[index].GetUncachedImage(); img != nil {
		c.Objects = nil
		c.Add(img)
	}

	if enlarge {
		bottom := container.NewCenter(widget.NewButton("Enlarge", func() {
			d := dialog.NewCustom("", "Close", previews[index].GetAsImageGallery(index, previews, false), ui.ActiveWindow())
			d.Resize(config.Get().Size())
			d.Show()
		}))
		return container.NewBorder(nil, bottom, left, right, c)
	}
	return container.NewBorder(nil, nil, left, right, c)
}

func (p *Preview) incrementIndex(i int, size int) int {
	i++
	if i == size {
		i = 0
	}
	return i
}

func (p *Preview) decrementIndex(i int, size int) int {
	i--
	if i < 0 {
		i = size - 1
	}
	return i
}
