package service

import (
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

type BannerService struct {
	BannerRepository port.BannerRepository
	imagedownloader  port.ImageDownloader
}

func NewBannerService(banner_repository port.BannerRepository, downloader port.ImageDownloader) *BannerService {
	return &BannerService{BannerRepository: banner_repository, imagedownloader: downloader}
}

func (bs *BannerService) CreateBannerUsingTheForm(images []multipart.File, banner entity.BannerDTO) (entity.BannerDTO, error) {
	newbanner := entity.NewBanner(banner)
	_, err := newbanner.Validations()

	if err != nil {
		return entity.BannerDTO{}, err
	}

	bannerWithImages := bs.SaveImages(images, *newbanner)

	bannerCreated, err := bs.BannerRepository.Create(bannerWithImages)

	if err != nil {
		return entity.BannerDTO{}, err
	}

	return bannerCreated, nil
}

func (bs *BannerService) SaveImages(images []multipart.File, banner entity.Banner) entity.Banner {
	var file string

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("Erro ao obter o caminho do executável:", err)
	}

	diretorio := strings.Replace(cwd, "test", "", -1) + "/images/banners/"

	for i, image := range images {

		file = banner.ID + "_" + string(i)
		pathFile := diretorio + file

		if i == 0 {
			err = bs.SaveImageForm(image, diretorio, file, 920, 90)

			if err != nil {
				pathFile = ""
			}

			banner.Image1 = pathFile
		} else if i == 1 {
			err = bs.SaveImageForm(image, diretorio, file, 728, 90)

			if err != nil {
				pathFile = ""
			}

			banner.Image2 = pathFile
		} else if i == 2 {
			err = bs.SaveImageForm(image, diretorio, file, 386, 386)

			if err != nil {
				pathFile = ""
			}

			banner.Image3 = pathFile
		}

	}

	return banner
}

func (bs *BannerService) SaveImageForm(file multipart.File, diretorio, id string, width, height int) error {

	if file == nil {
		return nil
	}

	defer file.Close()

	pathImage := diretorio + id + ".jpg"

	f, err := os.Create(pathImage)
	if err != nil {
		return ErrParseForm
	}
	defer f.Close()
	io.Copy(f, file)

	f, err = os.Open(pathImage)

	if err != nil {
		return ErrParseForm
	}

	// Resize the image
	img, _, err := image.Decode(f)
	if err != nil {
		return ErrDecodeImage
	}

	err = bs.imagedownloader.ResizeAndSaveImage(img, width, height, pathImage)

	if err != nil {
		fmt.Println(err)
		return ErrDecodeImage
	}

	return nil

}
