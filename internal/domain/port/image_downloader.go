package port

import "image"

type ImageDownloader interface {
	DownloadImage(url string) (image.Image, error)
	ResizeAndSaveImage(img image.Image, width, height int, outputPath string) error
	SaveImage(img image.Image, width, height int, outputPath string) error
	CropAndSaveImage(img image.Image, width, height int, outputPath string) error
}
