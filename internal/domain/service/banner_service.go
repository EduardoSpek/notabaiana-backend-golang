package service

import (
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

var (
	banner_pc_dimensions     = [2]int{1300, 190}
	banner_tablet_dimensions = [2]int{726, 106}
	banner_mobile_dimensions = [2]int{386, 386}
)

type BannerService struct {
	BannerRepository port.BannerRepository
	imagedownloader  port.ImageDownloader
}

func NewBannerService(banner_repository port.BannerRepository, downloader port.ImageDownloader) *BannerService {
	return &BannerService{BannerRepository: banner_repository, imagedownloader: downloader}
}

func (bs *BannerService) FindBanner(id string) (entity.BannerDTO, error) {
	banner, err := bs.BannerRepository.GetByID(id)

	if err != nil {
		return entity.BannerDTO{}, err
	}
	return banner, nil
}

func (bs *BannerService) Delete(id string) error {
	err := bs.BannerRepository.Delete(id)

	if err != nil {
		return err
	}
	return nil
}

func (bs *BannerService) AdminFindAll() (interface{}, error) {

	banners, err := bs.BannerRepository.AdminFindAll()

	if err != nil {
		return nil, err
	}

	result := struct {
		Banners []entity.BannerDTO `json:"banners"`
	}{
		Banners: banners,
	}

	return result, nil

}

func (bs *BannerService) FindAll() (interface{}, error) {

	banners, err := bs.BannerRepository.FindAll()

	if err != nil {
		return nil, err
	}

	result := struct {
		Banners []entity.BannerDTO `json:"banners"`
	}{
		Banners: banners,
	}

	return result, nil

}

func (bs *BannerService) UpdateBannerUsingTheForm(images []multipart.File, banner entity.BannerDTO) (entity.BannerDTO, error) {

	currentBanner, err := bs.BannerRepository.GetByID(banner.ID)

	if err != nil {
		return entity.BannerDTO{}, err
	}

	currentBanner.Title = banner.Title
	currentBanner.Link = banner.Link
	currentBanner.Html = banner.Html
	currentBanner.Tag = banner.Tag
	currentBanner.Visible = banner.Visible

	newbanner := entity.UpdateBanner(currentBanner)
	_, err = newbanner.Validations()

	if err != nil {
		return entity.BannerDTO{}, err
	}

	bannerWithImages := bs.SaveImages(images, *newbanner)

	bannerCreated, err := bs.BannerRepository.Update(bannerWithImages)

	if err != nil {
		return entity.BannerDTO{}, err
	}

	return bannerCreated, nil
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
	var err error

	diretorio := "/images/banners/"

	for i, image := range images {

		file = banner.ID + "_" + strconv.Itoa(i) + ".jpg"
		pathFile := diretorio + file

		if i == 0 {
			err = bs.SaveImageForm(image, diretorio, file, banner_pc_dimensions[0], banner_pc_dimensions[1])

			if err != nil {
				pathFile = ""
			}

			if image != nil {
				banner.Image1 = pathFile
			}
		} else if i == 1 {
			err = bs.SaveImageForm(image, diretorio, file, banner_tablet_dimensions[0], banner_tablet_dimensions[1])

			if err != nil {
				pathFile = ""
			}

			if image != nil {
				banner.Image2 = pathFile
			}
		} else if i == 2 {
			err = bs.SaveImageForm(image, diretorio, file, banner_mobile_dimensions[0], banner_mobile_dimensions[1])

			if err != nil {
				pathFile = ""
			}

			if image != nil {
				banner.Image3 = pathFile
			}
		}

	}

	return banner
}

func (bs *BannerService) SaveImageForm(file multipart.File, diretorio, filename string, width, height int) error {

	if file == nil {
		return nil
	}

	defer file.Close()

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("Erro ao obter o caminho do executável:", err)
	}

	pathImage := strings.Replace(cwd, "test", "", -1) + diretorio + filename

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
