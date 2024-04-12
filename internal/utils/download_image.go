package utils

import (
	"image"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

type ImgDownloader struct {}

func NewImgDownloader() *ImgDownloader {
	return &ImgDownloader{}
}

func (d *ImgDownloader) DownloadImage(url string) (image.Image, error) {
	// Baixa a imagem da URL
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Decodifica a imagem
	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (d *ImgDownloader) ResizeAndSaveImage(img image.Image, width, height int, outputPath string) error {
	// Redimensiona a imagem
	resizedImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	// Cria o arquivo de saída
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Salva a imagem redimensionada no arquivo
	err = jpeg.Encode(outputFile, resizedImg, nil)
	if err != nil {
		return err
	}

	return nil
}