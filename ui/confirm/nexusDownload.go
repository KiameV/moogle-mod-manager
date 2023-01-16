package confirm

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/kiamev/moogle-mod-manager/mods"
	"github.com/kiamev/moogle-mod-manager/ui/state/ui"
	"github.com/kiamev/moogle-mod-manager/ui/util"
	"os"
	"path/filepath"
	"time"
)

type toDownload struct {
	uri      string
	dir      string
	fileName string
}

type manualDownloadConfirmer struct {
	Params
}

func newManualDownloadConfirmer(params Params) Confirmer {
	return &manualDownloadConfirmer{Params: params}
}

func (c *manualDownloadConfirmer) Downloads(done func(mods.Result)) (err error) {
	var (
		toDl     []toDownload
		fileName string
	)
	for _, ti := range c.ToInstall {
		fileName, _ = ti.Download.FileName()
		if ti.Download != nil {
			dl := toDownload{
				fileName: fileName,
			}
			if dl.dir, err = ti.GetDownloadLocation(c.Game, c.Mod); err != nil {
				return
			}
			if _, err = os.Stat(filepath.Join(dl.dir, fileName)); err == nil {
				continue
			}
			toDl = append(toDl, dl)
		}
	}

	if len(toDl) == 0 {
		done(mods.Ok)
		return nil
	}

	return c.showDialog(toDl, done)
}

func (c *manualDownloadConfirmer) showDialog(toDl []toDownload, done func(mods.Result)) (err error) {
	var (
		fi   []*widget.FormItem
		rows []*downloadRow
	)

	for i, td := range toDl {
		r := newDownloadRow(&td)
		rows = append(rows, r)
		text := "Place download in:"
		if len(toDl) == 1 && clipboard.WriteAll(td.dir) == nil {
			text += " (copied to clipboard)"
		}

		fi = append(fi, widget.NewFormItem(fmt.Sprintf("%d:", i+1), r))
		fi = append(fi, widget.NewFormItem("",
			widget.NewLabelWithStyle("Download the following file/s:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})))
		fi = append(fi, widget.NewFormItem("",
			util.CreateUrlRow(td.uri)))
		fi = append(fi, widget.NewFormItem("",
			widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})))
		fi = append(fi, widget.NewFormItem("",
			util.CreateUrlRow(td.dir)))
	}

	fi = append(fi, widget.NewFormItem("", container.NewCenter(widget.NewButton("Check", func() {
		for _, r := range rows {
			if e := r.Validate(); e != nil {
				util.ShowErrorLong(e)
				return
			}
		}
	}))))
	d := dialog.NewForm("Download Files", "Done", "Cancel", fi, func(ok bool) {
		result := mods.Ok
		if !ok {
			result = mods.Cancel
		}
		done(result)
	}, ui.Window)
	d.SetOnClosed(func() {
		for _, r := range rows {
			r.stop = true
		}
		for _, r := range rows {
			r.SetOnValidationChanged(nil)
		}
	})
	d.Resize(fyne.NewSize(500, 450))
	d.Show()
	return
}

type downloadRow struct {
	fyne.Validatable
	fyne.Widget
	fileNeeded        string
	stop              bool
	validatedCallback func(error)
}

func newDownloadRow(td *toDownload) *downloadRow {
	r := &downloadRow{
		fileNeeded: filepath.Join(td.dir, td.fileName),
		Widget:     widget.NewLabel("Found"),
	}
	if r.Found() {
		r.Show()
	} else {
		r.Hide()
	}
	r.start()
	return r
}

func (r *downloadRow) start() {
	go func() {
		var err error
		for !r.stop {
			err = r.Validate()
			if err == nil {
				r.Show()
			} else {
				r.Hide()
			}
			if r.validatedCallback != nil {
				r.validatedCallback(err)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

func (r *downloadRow) Stop() {
	r.stop = true
}

func (r *downloadRow) Found() bool {
	_, err := os.Stat(r.fileNeeded)
	return err == nil
}

func (r *downloadRow) Validate() error {
	if !r.Found() {
		return fmt.Errorf("The following file was not found:\n%s", r.fileNeeded)
	}
	return nil
}

func (r *downloadRow) SetOnValidationChanged(validatedCallback func(error)) {
	r.validatedCallback = validatedCallback
}
