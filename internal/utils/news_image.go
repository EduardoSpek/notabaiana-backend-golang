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

func GetHeightPositionImage(y int, label string, face font.Face) int {

	var yy int
	// Wrap the text
	maxWidth := 645
	wrappedLines := WrapText(label, face, maxWidth)

	// Draw each line
	lineHeight := 55
	for i := range wrappedLines {
		yy = y + i*lineHeight
	}

	return yy
}

func AddLabel(img *image.RGBA, x, y int, label string, face font.Face) int {

	var yy int
	// Wrap the text
	maxWidth := 645
	wrappedLines := WrapText(label, face, maxWidth)

	// Draw each line
	lineHeight := 55
	for i, line := range wrappedLines {
		yy = y + i*lineHeight                   // Starting Y position with line spacing
		drawTextOnImage(img, line, x, yy, face) // x = 20 is the left margin
	}

	return yy
}

// WrapText wraps text within a given width
func WrapText(text string, face font.Face, maxWidth int) []string {
	words := strings.Fields(text)
	var wrappedLines []string
	var line string

	for _, word := range words {
		// Test if adding the word would exceed maxWidth
		if measureTextWidth(line+" "+word, face) > maxWidth && line != "" {
			wrappedLines = append(wrappedLines, line)
			line = word
		} else {
			if line != "" {
				line += " "
			}
			line += word
		}
	}

	if line != "" {
		wrappedLines = append(wrappedLines, line)
	}

	return wrappedLines
}

// measureTextWidth calculates the width of a string using a given font
func measureTextWidth(text string, face font.Face) int {
	var width fixed.Int26_6
	for _, x := range text {
		advance, ok := face.GlyphAdvance(x)
		if !ok {
			advance = 0
		}
		width += advance
	}
	return int(width >> 6) // Convert from fixed-point to integer
}

func drawTextOnImage(img draw.Image, text string, x, y int, face font.Face) {
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)
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
