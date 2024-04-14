package service

import (
	"image"

	"github.com/eduardospek/bn-api/internal/domain/entity"
)

type NewsRepository interface {
	Create(news entity.News) (entity.News, error)
	FindAll(page, limit int) ([]entity.News, error)
	NewsExists(title string) error
	GetBySlug(slug string) (entity.News, error)
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
	new = RenamePathImage(new)
	new = ChangeLink(new)

	err := s.newsrepository.NewsExists(new.Title)

	if err != nil {
		return entity.News{}, err
	}
	
	_, err = s.newsrepository.Create(new)
	
	if err != nil {
		return entity.News{}, err
	}
	
	return new, nil
}

func (s *NewsService) GetNewsBySlug(slug string) (entity.News, error) {
	
	new, err := s.newsrepository.GetBySlug(slug)

	if err != nil {
		return entity.News{}, err
	}
	
	return new, nil

}

func (s *NewsService) FindAllNews(page, limit int) []entity.News {
	
	news, _ := s.newsrepository.FindAll(page, limit)
	
	return news

}

func (s *NewsService) SaveImage(id, url, diretorio string) error {
	
	img, err := s.imagedownloader.DownloadImage(url)
	
	if err != nil {
		//fmt.Println("Erro ao baixar a imagem:", err)
		return err
	}
	
	outputPath := diretorio + id + ".jpg"

	width := 400
	height := int(float64(img.Bounds().Dy()) * (float64(width) / float64(img.Bounds().Dx()))) 

	err = s.imagedownloader.ResizeAndSaveImage(img, width, height, outputPath)
	
	if err != nil {
		//fmt.Println("Erro ao redimensionar e salvar a imagem:", err)
		return err
	}

	//fmt.Println("Imagem redimensionada e salva com sucesso em", outputPath)
	
	return nil

}

func RenamePathImage(news entity.News) entity.News {
	news.Image = news.ID + ".jpg"
	return news
}
func ChangeLink(news entity.News) entity.News {
	news.Link = "/news/" + news.Slug
	return news
}