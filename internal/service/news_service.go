package service

import (
	"errors"
	"image"
	"strings"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	"github.com/gocolly/colly"
)

type NewsRepository interface {
	Create(news entity.News) (entity.News, error)
	FindAll(page, limit int) (interface{}, error)
	NewsExists(title string) error
	GetBySlug(slug string) (entity.News, error)
	NewsTruncateTable() error
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

	result := containsWordsInTitle(new.Title)

	if result {
		return entity.News{}, errors.New("título com palavra bloqueada")
	}
	
	new.Text = changeWords(new.Text)

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

func (s *NewsService) FindAllNews(page, limit int) interface{} {
	
	news, _ := s.newsrepository.FindAll(page, limit)
	
	return news

}

func (s *NewsService) NewsTruncateTable() error {
	
	err := s.newsrepository.NewsTruncateTable()

	if err != nil {
		return err
	}
	
	return nil

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

func (s *NewsService) GetImagesPage(url string) string {

	var html string

	collector := colly.NewCollector(
        colly.AllowedDomains("www.bahianoticias.com.br"),
    )	

    // Definindo o callback OnHTML
	collector.OnHTML("img", func(e *colly.HTMLElement) {
		// Obter o valor do atributo "src" da imagem
		src := e.Attr("src")

		src = strings.Replace(src, " ", "%20", -1)
	
		// Verificar se o valor "src" contém o endereço de destino
		if strings.Contains(src, "bahianoticias.com.br/fotos/") {
			
			html += `<div class="imagem_anexada"><img src="` + src + `" width="100%"></div>`
						
		}
	})

    // Visitando a URL inicial
    collector.Visit(url)

	return html

}

func RenamePathImage(news entity.News) entity.News {
	news.Image = news.ID + ".jpg"
	return news
}
func ChangeLink(news entity.News) entity.News {
	news.Link = "/news/" + news.Slug
	return news
}
func containsWordsInTitle(titulo string) bool {
	palavras := []string {
		"Bahia Notícias",
		"BN",
		"Curtas e Venenosas",
	}
    for _, palavra := range palavras {
        if strings.Contains(titulo, palavra) {
            return true
        }
    }
    return false
}
func changeWords(text string) string {
	text = strings.Replace(text, "Bahia Notícias", "BN", -1)
    return text
}