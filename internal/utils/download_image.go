package utils

import (
	"image"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

type ImgDownloader struct{}

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

func (d *ImgDownloader) CropAndSaveImage(img image.Image, width, height int, outputPath string) error {

	// Definir a largura e a altura desejadas do recorte
	cropWidth, cropHeight := width, height

	// Obter as dimensões da imagem original
	srcWidth := img.Bounds().Max.X
	srcHeight := img.Bounds().Max.Y

	// Calcular as proporções da imagem original e do recorte desejado
	srcAspect := float64(srcWidth) / float64(srcHeight)
	cropAspect := float64(cropWidth) / float64(cropHeight)

	// Definir as dimensões da área de recorte para incluir o máximo possível da imagem original
	if srcAspect > cropAspect {
		// A imagem original é mais larga do que a proporção do recorte
		height = srcHeight
		width = int(float64(height) * cropAspect)
	} else {
		// A imagem original é mais alta do que a proporção do recorte
		width = srcWidth
		height = int(float64(width) / cropAspect)
	}

	// Calcular as coordenadas do ponto superior esquerdo para centralizar o recorte
	x0 := (srcWidth - width) / 2
	y0 := (srcHeight - height) / 2

	// Realizar o recorte
	cropped := imaging.Crop(img, image.Rect(x0, y0, x0+width, y0+height))

	// Redimensionar o recorte para as dimensões desejadas
	resized := imaging.Resize(cropped, cropWidth, cropHeight, imaging.Lanczos)

	// Salvar a nova imagem recortada e redimensionada
	err := imaging.Save(resized, outputPath)
	if err != nil {
		return err
	}

	return nil

}
