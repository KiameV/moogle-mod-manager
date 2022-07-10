package util

import (
	"github.com/aarzilli/nucular"
	"image"
	"image/draw"
	"image/png"
	"os"
)

var imgs = make(map[string]*image.RGBA)

func AddFieldValueText(w *nucular.Window, label string, value string) {
	w.Row(12).Static(w.Bounds.W - 20)
	w.Label(label, "LC")
	w.Row(12).Static(w.Bounds.W - 20)
	w.Label(value, "LC")
	w.Row(8).Static()
}

func loadImg(imagePath string) *image.RGBA {
	img, ok := imgs[imagePath]
	if !ok {
		var i image.Image
		// TODO fix for internet paths
		fr, err := os.Open(imagePath)
		if err != nil {
			i = image.NewRGBA(image.Rect(0, 0, 1, 1))
		} else {
			if i, err = png.Decode(fr); err != nil {
				i = image.NewRGBA(image.Rect(0, 0, 1, 1))
			}
		}
		p := i.Bounds().Max
		x := p.X
		y := p.Y
		if x > 400 {
			x = 400
		}
		if y > 300 {
			y = 300
		}
		img = image.NewRGBA(image.Rect(0, 0, x, y))
		draw.Draw(img, img.Bounds(), i, image.Point{}, draw.Src)
		imgs[imagePath] = img
	}
	return img
}

func DrawImg(w *nucular.Window, imagePath string) {
	if imagePath == "" {
		return
	}
	var (
		img  = loadImg(imagePath)
		size = img.Bounds().Max
	)
	w.Row(size.Y).Static(size.X)
	w.Image(img)
}

/* Change to remote image loading
import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

func main() {
    url := "http://i.imgur.com/m1UIjW1.jpg"
    // don't worry about errors
    response, e := http.Get(url)
    if e != nil {
        log.Fatal(e)
    }
    defer response.Body.Close()

    //open a file for writing
    file, err := os.Create("/tmp/asdf.jpg")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // Use io.Copy to just dump the response body to the file. This supports huge files
    _, err = io.Copy(file, response.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Success!")
}
*/
