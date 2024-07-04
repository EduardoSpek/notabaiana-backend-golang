package utils

import (
	"image"
	"image/color"
	"image/draw"
	"net/http"
	"os"
	"strings"

	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func DownloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ResizeImage(img image.Image, width, height uint) image.Image {
	return resize.Resize(width, height, img, resize.Lanczos3)
}

func OverlayImage(base, overlay image.Image, offsetX, offsetY int) *image.RGBA {
	b := base.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, base, image.Point{}, draw.Src)

	o := overlay.Bounds()
	point := image.Pt(offsetX, offsetY)
	draw.Draw(m, o.Add(point), overlay, image.Point{}, draw.Over)

	return m
}

func AddLabel(img *image.RGBA, x, y int, label string, face font.Face) {

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  fixed.P(x, y),
	}

	words := strings.Split(label, " ")
	line := ""
	for i, word := range words {
		if (i+1)%4 == 0 {
			d.DrawString(line + " " + word)
			y += 50 // Adjust line height as needed
			d.Dot = fixed.P(x, y)
			line = ""
		} else {
			line += " " + word
		}
	}

	if line != "" {
		d.DrawString(line)
	}
}

func LoadFont(path string, size float64) (font.Face, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fnt, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
