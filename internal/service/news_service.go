package service

import (
	"errors"
	"image"
	"net/url"
	"strings"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	"github.com/gocolly/colly"
)

var (
		ErrNoCategory = errors.New("nenhuma categoria no rss")
		AllowedDomains = "www.bahianoticias.com.br"
	)

type NewsRepository interface {
	Create(news entity.News) (entity.News, error)
	FindAll(page, limit int) (interface{}, error)
	FindCategory(category string, page int) (interface{}, error)
	NewsExists(title string) error
	GetBySlug(slug string) (entity.News, error)
	NewsTruncateTable() error
	FindAllViews() ([]entity.News, error)
	ClearViews() error
	SearchNews(page int, str_search string) interface{}
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

func (s *NewsService) SearchNews(page int, str_search string) interface{} {
	
	news := s.newsrepository.SearchNews(page, str_search)
	
	return news

}

func (s *NewsService) FindAllNews(page, limit int) interface{} {
	
	news, _ := s.newsrepository.FindAll(page, limit)
	
	return news

}

func (s *NewsService) FindNewsCategory(category string, page int) interface{} {
	
	news, _ := s.newsrepository.FindCategory(category, page)
	
	return news

}


func (s *NewsService) FindAllViews() ([]entity.News, error) {
	
	news, err := s.newsrepository.FindAllViews()

	if err != nil {
		return []entity.News{}, err
	}
	
	return news, nil

}

func (s *NewsService) ClearViews() error {
	
	err := s.newsrepository.ClearViews()

	if err != nil {
		return err
	}
	
	return nil

}

func (s *NewsService) NewsTruncateTable() error {
	
	err := s.newsrepository.NewsTruncateTable()

	if err != nil {
		return err
	}
	
	return nil

}

func (s *NewsService) SaveImage(id, url, diretorio string) error {

	url = strings.Replace(url, "_original.jpg", "_5.jpg", -1)
	
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

func (s *NewsService) GetEmded(link string) (string, string) {
	var html, conteudo, script, text, str_text string
	//var err error

	collector := colly.NewCollector(
		colly.AllowedDomains(AllowedDomains),
	)

	//Obtém o texto da notícia
	collector.OnHTML(".sc-16306eb7-3.lbjQbj", func(e *colly.HTMLElement) {
		
		str_text = e.DOM.Text()		

		text += str_text + "<br><br>"
	
	})

	//Obtém os embeds das redes sociais
	collector.OnHTML(".lazyload-placeholder", func(e *colly.HTMLElement) {
		// Obter o valor do atributo "src" da imagem
		conteudo = e.Attr("data-content")
		conteudo_decoded, err := url.QueryUnescape(conteudo)

		if err != nil {			
			return
		}	

		html +=  "<br><br>" + conteudo_decoded
	
	})

	//Obtém os scrits
	collector.OnHTML(".lazyload-scripts", func(e *colly.HTMLElement) {
		// Obter o valor do atributo "src" da imagem
		script = e.Attr("data-scripts")

		script_decoded, err := url.QueryUnescape(script)

		if err != nil {			
			return
		}

		if strings.Contains(script_decoded, "instagram") || strings.Contains(script_decoded, "twitter") || strings.Contains(script_decoded, "youtube") {
		
			if !strings.Contains(html, script_decoded) {
				html += script_decoded
			}
		}
		
	})

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
	collector.Visit(link)

	html = text + html

	return html, text
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
	text = strings.Replace(text, "@BahiaNoticias", "@", -1)
	text = strings.Replace(text, "@bhall", "@", -1)
	text = strings.Replace(text, "@bnholofote", "@", -1)
	text = strings.Replace(text, "BN Holofote", "", -1)
	text = strings.Replace(text, "Siga o @bnhall_ no Instagram e fique de olho nas principais notícias.", "", -1)
	text = strings.Replace(text, "As informações são do Metrópoles, parceiro do BN", ".", -1)
	text = strings.Replace(text, " parceiro do BN,", "", -1)
	text = strings.Replace(text, "Assine a newsletter de Esportes do BN e fique bem informado sobre o esporte na Bahia, no Brasil e no mundo!", "", -1)
	
    return text
}

func (s *NewsService) GetCategory(rss string) (string, error) {

	if  strings.Contains(rss, "holofote") {
		return "famosos", nil
	} else if  strings.Contains(rss, "esportes") {
		return "esportes", nil
	} else if  strings.Contains(rss, "justica") {
		return "justica", nil
	} else if  strings.Contains(rss, "saude") {
		return "saude", nil
	} else if  strings.Contains(rss, "municipios") {
		return "municipios", nil
	} else {
		return "", ErrNoCategory
	}

}