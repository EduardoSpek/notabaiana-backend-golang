package service

import (
	"fmt"
	"image"

	"github.com/eduardospek/bn-api/internal/domain/entity"
)

type NewsRepository interface {
	Create(news entity.News) error
	FindAll() []entity.News
}

type ImageDownloader interface {
	DownloadImage(url string) (image.Image, error)
	ResizeAndSaveImage(img image.Image, width, height int, outputPath string) error
}

type NewsService struct {
	newsrepository NewsRepository
	imagedownloader ImageDownloader
}

func NewNewsService(repository NewsRepository, downloader ImageDownloader) *NewsService {
	return &NewsService{ newsrepository: repository, imagedownloader: downloader }
}

func (s *NewsService) CreateNews(news entity.News) (entity.News, error) {
	
	new := *entity.NewNews(news)
	
	err := s.newsrepository.Create(new)
	
	if err != nil {
		return entity.News{}, err
	}
	
	return new, nil
}

func (s *NewsService) FindAllNews() []entity.News {
	
	news := s.newsrepository.FindAll()
	
	return news

}

func (s *NewsService) SaveImage(id, url, diretorio string) error {
	
	img, err := s.imagedownloader.DownloadImage(url)
	
	if err != nil {
		fmt.Println("Erro ao baixar a imagem:", err)
		return err
	}
	
	outputPath := diretorio + id + ".jpg"

	err = s.imagedownloader.ResizeAndSaveImage(img, 400, 267, outputPath)
	
	if err != nil {
		fmt.Println("Erro ao redimensionar e salvar a imagem:", err)
		return err
	}

	fmt.Println("Imagem redimensionada e salva com sucesso em", outputPath)
	
	return nil

}