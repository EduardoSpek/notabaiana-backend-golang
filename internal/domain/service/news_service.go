package service

import (
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/gocolly/colly"
)

var (
	ErrNotAuthorized  = errors.New("você não tem autorização para criar notícias")
	ErrDecodeImage    = errors.New("não foi possível decodificar a imagem")
	ErrCreateNews     = errors.New("não foi possível criar a notícia")
	ErrUpdateNews     = errors.New("não foi possível atualizar a notícia")
	ErrParseForm      = errors.New("erro ao obter a imagem")
	ErrWordsBlackList = errors.New("o título contém palavras bloqueadas")
	ErrNoCategory     = errors.New("nenhuma categoria no rss")
	ErrSimilarTitle   = errors.New("título similar ao recente adicionado detectado")
	AllowedDomains    = "www.bahianoticias.com.br"

	LimitPerPage = 100
)

type NewsService struct {
	hitsrepository  port.HitsRepository
	newsrepository  port.NewsRepository
	imagedownloader port.ImageDownloader
}

func NewNewsService(repository port.NewsRepository, downloader port.ImageDownloader, hits port.HitsRepository) *NewsService {
	return &NewsService{newsrepository: repository, imagedownloader: downloader, hitsrepository: hits}
}

func (s *NewsService) AdminDeleteAll(banners []entity.News) error {
	err := s.newsrepository.DeleteAll(banners)

	if err != nil {
		return err
	}

	return nil
}

func (s *NewsService) Delete(id string) error {
	err := s.newsrepository.Delete(id)

	if err != nil {
		return err
	}
	return nil
}

func (s *NewsService) StartCleanNews(minutes time.Duration) {

	go s.CleanNews()

	ticker := time.NewTicker(minutes * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		go s.CleanNews()
	}
}

func (s *NewsService) CleanNews() {

	s.newsrepository.CleanNews()

}

func (s *NewsService) NewsMake() (interface{}, error) {

	news, err := s.newsrepository.NewsMake()

	if err != nil {
		return nil, err
	}

	newsOutput := struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Link  string `json:"link"`
		Image string `json:"image"`
	}{
		ID:    news.ID,
		Title: news.Title,
		Link:  news.Link,
		Image: news.Image,
	}

	return newsOutput, nil

}

func (s *NewsService) UpdateNewsUsingTheForm(file multipart.File, newsInput entity.News) (entity.News, error) {

	oldnew, err := s.newsrepository.GetBySlug(newsInput.Slug)

	if err != nil {
		return entity.News{}, err
	}

	oldnew.Title = newsInput.Title
	oldnew.Text = newsInput.Text
	oldnew.Visible = newsInput.Visible
	oldnew.TopStory = newsInput.TopStory
	oldnew.Category = newsInput.Category

	newNews := entity.UpdateNews(oldnew)

	news := ChangeLink(*newNews)

	if file != nil {
		news = RenamePathImage(news)
	}

	new, err := s.newsrepository.Update(news)

	if err != nil {
		return entity.News{}, err
	}

	err = s.SaveImageForm(file, new)

	if err != nil {
		new.Image = ""
	}

	return new, nil

}

func (s *NewsService) CreateNewsUsingTheForm(file multipart.File, news entity.News) (entity.News, error) {

	newNews := entity.NewNews(news)

	new, err := s.CreateNews(*newNews)

	if err != nil {
		return entity.News{}, err
	}

	err = s.SaveImageForm(file, new)

	if err != nil {
		newNews.Image = ""
	}

	return new, nil

}

func (s *NewsService) SaveImageForm(file multipart.File, news entity.News) error {

	if file == nil {
		return nil
	}

	defer file.Close()

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("Erro ao obter o caminho do executável:", err)
	}

	diretorio := strings.Replace(cwd, "test", "", -1) + "/images/"
	pathImage := diretorio + news.Image

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

	err = s.imagedownloader.CropAndSaveImage(img, 400, 254, pathImage)

	if err != nil {
		fmt.Println(err)
		return ErrDecodeImage
	}

	return nil

}

func (s *NewsService) CreateNews(news entity.News) (entity.News, error) {

	new := *entity.NewNews(news)

	recent, err := s.newsrepository.FindRecent()

	if err != nil {
		return entity.News{}, err
	}

	similarity := utils.Similarity(recent.Title, new.Title)

	if similarity > 70 {
		fmt.Println("***titulo silimar detectado***")
		return entity.News{}, ErrSimilarTitle
	}

	new = RenamePathImage(new)
	new = ChangeLink(new)

	result := listOfBlockedWords(new.Title)

	if result {
		return entity.News{}, ErrWordsBlackList
	}

	new.Text = changeWords(new.Text)

	err = s.newsrepository.NewsExists(new.Title)

	if err != nil {
		return entity.News{}, err
	}

	//newtitle, err := utils.ChangeTitleWithGemini("Refaça este título e mantanha o contexto. Utilize palavras-chave para melhorar o SEO. O texto não deve ultrapassar os 127 caracteres. O título é", new.Title)

	//if err == nil && newtitle != "" {
	//	new.TitleAi = strings.TrimSpace(newtitle)
	//}

	_, err = s.newsrepository.Create(new)

	if err != nil {
		return entity.News{}, err
	}

	return new, nil
}

func (s *NewsService) GetNewsBySlug(slug string) (entity.News, error) {

	// err := s.Hit(slug)

	// if err != nil {
	// 	return entity.News{}, err
	// }

	new, err := s.newsrepository.GetBySlug(slug)

	if err != nil {
		return entity.News{}, err
	}

	return new, nil

}

func (s *NewsService) Hit(session string) error {

	ip := utils.GetIP()

	_, err := s.hitsrepository.Get(ip, session)

	if err != nil {

		newhit := entity.Hits{
			IP:      ip,
			Session: session,
			Views:   1,
		}

		err = s.hitsrepository.Save(newhit)

		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (s *NewsService) SearchNews(page int, str_search string) interface{} {

	str_search = strings.Replace(str_search, " ", "%", -1)

	news := s.newsrepository.SearchNews(page, str_search)

	total := s.newsrepository.GetTotalNewsBySearch(str_search)

	pagination := utils.Pagination(page, total)

	result := struct {
		List_news  []entity.News    `json:"news"`
		Pagination map[string][]int `json:"pagination"`
		Search     string           `json:"search"`
	}{
		List_news:  news,
		Pagination: pagination,
		Search:     strings.Replace(str_search, "%", " ", -1),
	}

	return result

}

func (s *NewsService) AdminFindAllNews(page, limit int) interface{} {

	//Limita o total de registros que deve ser retornado
	if limit > LimitPerPage {
		limit = LimitPerPage
	}

	news, _ := s.newsrepository.AdminFindAll(page, limit)

	total := s.newsrepository.GetTotalNews()

	pagination := utils.Pagination(page, total)

	result := struct {
		List_news  []entity.News    `json:"news"`
		Pagination map[string][]int `json:"pagination"`
	}{
		List_news:  news,
		Pagination: pagination,
	}

	return result

}

func (s *NewsService) FindAllNews(page, limit int) interface{} {

	//Limita o total de registros que deve ser retornado
	if limit > LimitPerPage {
		limit = LimitPerPage
	}

	news, _ := s.newsrepository.FindAll(page, limit)

	total := s.newsrepository.GetTotalNewsVisible()

	pagination := utils.Pagination(page, total)

	result := struct {
		List_news  []entity.News    `json:"news"`
		Pagination map[string][]int `json:"pagination"`
	}{
		List_news:  news,
		Pagination: pagination,
	}

	return result

}

func (s *NewsService) FindRecent() (entity.News, error) {

	news, err := s.newsrepository.FindRecent()

	if err != nil {
		return entity.News{}, err
	}

	return news, nil

}

func (s *NewsService) FindNewsCategory(category string, page int) interface{} {

	news, _ := s.newsrepository.FindCategory(category, page)

	total := s.newsrepository.GetTotalNewsByCategory(category)

	pagination := utils.Pagination(page, total)

	result := struct {
		List_news  []entity.News    `json:"news"`
		Pagination map[string][]int `json:"pagination"`
		Category   string           `json:"category"`
	}{
		List_news:  news,
		Pagination: pagination,
		Category:   category,
	}

	return result
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

func (s *NewsService) ClearImagePath(id string) error {

	err := s.newsrepository.ClearImagePath(id)

	if err != nil {
		return err
	}

	return nil

}

func (s *NewsService) SaveImage(id, url, diretorio string) (string, error) {

	url = strings.Replace(url, "_original.jpg", "_5.jpg", -1)

	img, err := s.imagedownloader.DownloadImage(url)

	if err != nil {
		//fmt.Println("Erro ao baixar a imagem:", err)
		return "", err
	}

	outputPath := diretorio + id + ".jpg"

	width := 400
	height := int(float64(img.Bounds().Dy()) * (float64(width) / float64(img.Bounds().Dx())))

	err = s.imagedownloader.CropAndSaveImage(img, width, height, outputPath)

	if err != nil {
		//fmt.Println("Erro ao redimensionar e salvar a imagem:", err)
		return "", err
	}

	//fmt.Println("Imagem redimensionada e salva com sucesso em", outputPath)

	return outputPath, nil

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

		html += "<br><br>" + conteudo_decoded

	})

	//Obtém os scrits
	collector.OnHTML(".lazyload-scripts", func(e *colly.HTMLElement) {
		// Obter o valor do atributo "src" da imagem
		script = e.Attr("data-scripts")

		script_decoded, err := url.QueryUnescape(script)

		if err != nil {
			return
		}

		if strings.Contains(script_decoded, "instagram") || strings.Contains(script_decoded, "twitter") || strings.Contains(script_decoded, "youtube") || strings.Contains(script_decoded, "flickr") {

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

	//html = text + html

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
func listOfBlockedWords(titulo string) bool {
	palavras := []string{
		"Bahia Notícias",
		"Bahia Notícia",
		"Bahia Noticias",
		"Bahia Noticia",
		"BN",
		"Curtas",
		"Nota Baiana",
		"NotaBaiana",
		"notabaiana",
		"apple-touch-icon.png",
		"Davidson pelo mundo",
		"Davidson Pelo Mundo",
		"Davidson pelo Mundo",
		"Lula",
		"PT",
		"PDT",
		"Avante Brasil",
		"PSDB",
		"Jaques Wagner",
		"Jerônimo",
		"LGBT",
		"Jaques Wagner",
		"Bolsonaro",
		"PL",
		"Govern",
		"govern",
		"convenção",
		"convenções",
		"candidat",
		"Nininha",
		"Kret",
		"Travesti",
		"travesti",
		"drag queen",
		"Drag Queen",
		"boyceta",
		"mulher trans",
	}
	for _, palavra := range palavras {
		if strings.Contains(titulo, palavra) {
			return true
		}
	}
	return false
}
func changeWords(text string) string {
	text = strings.Replace(text, " ", " ", -1)
	text = strings.Replace(text, "Siga o @bnhall_ no Instagram e fique de olho nas principais notícias.", "", -1)

	text = strings.Replace(text, "BN", "NB", -1)

	text = strings.Replace(text, "Bahia Notícias", "NB", -1)

	text = strings.Replace(text, "Bahia Notícia", "NB", -1)

	text = strings.Replace(text, "Bahia Noticia", "NB", -1)

	text = strings.Replace(text, "Bahia Noticias", "NB", -1)

	text = strings.Replace(text, "@BahiaNoticias", "@", -1)

	text = strings.Replace(text, "@bnholofote", "@", -1)

	text = strings.Replace(text, "BN Holofote", "", -1)

	text = strings.Replace(text, "@bhall", "@", -1)

	text = strings.Replace(text, "As informações são do Metrópoles, parceiro do NB", ".", -1)

	text = strings.Replace(text, " parceiro do NB,", "", -1)

	text = strings.Replace(text, "Assine a newsletter de Esportes do NB e fique bem informado sobre o esporte na Bahia, no Brasil e no mundo!", "", -1)

	text = strings.Replace(text, "Siga o NB no Google News e veja os conteúdos de maneira ainda mais rápida e ágil pelo celular ou pelo computador!", "", -1)

	return text
}

func (s *NewsService) GetCategory(rss string) (string, error) {

	if strings.Contains(rss, "holofote") {
		return "famosos", nil
	} else if strings.Contains(rss, "esportes") {
		return "esportes", nil
	} else if strings.Contains(rss, "justica") {
		return "justica", nil
	} else if strings.Contains(rss, "saude") {
		return "saude", nil
	} else if strings.Contains(rss, "municipios") {
		return "municipios", nil
	} else {
		return "", ErrNoCategory
	}

}

func (s *NewsService) GetNewsFromPage(link string) []entity.News {
	var conteudo string
	//var err error

	collector := colly.NewCollector(
		colly.AllowedDomains(AllowedDomains),
	)

	var titulos []string
	var texts []string
	var links []string
	var images []string

	//Obtém titulos
	collector.OnHTML("h3.sc-b4c8ccf3-1.ireAxk", func(e *colly.HTMLElement) {

		conteudo = e.DOM.Text()

		titulos = append(titulos, conteudo)

	})

	//Obtém links
	collector.OnHTML(".sc-b4c8ccf3-0.fsXNOt a", func(e *colly.HTMLElement) {

		conteudo = e.Attr("href")

		conteudo = "https://www.bahianoticias.com.br" + conteudo

		image, err := s.GetImageLink(conteudo)

		if err != nil {
			return
		}

		images = append(images, image)

		links = append(links, conteudo)

	})

	//Obtém textos
	collector.OnHTML(".sc-81cf810-3.gCNTHg", func(e *colly.HTMLElement) {

		conteudo = e.DOM.Text()

		texts = append(texts, conteudo)

	})

	//Obtém images
	collector.OnHTML(".sc-81cf810-2.hiSMeg div span img", func(e *colly.HTMLElement) {

		conteudo = e.Attr("src")

		conteudo_decoded, err := url.QueryUnescape(conteudo)

		if err != nil {
			return
		}

		images = append(images, conteudo_decoded)

	})

	// Visitando a URL inicial
	collector.Visit(link)

	var lista []entity.News

	for i, item := range titulos {

		category, _ := s.GetCategory(links[i])

		new := entity.NewNews(entity.News{
			Title:    item,
			Text:     texts[i],
			Image:    images[i],
			Link:     links[i],
			Visible:  true,
			TopStory: false,
			Category: category,
		})
		lista = append(lista, *new)
	}

	return lista
}

func (s *NewsService) GetIdFromLink(link string) (int, error) {

	partes := strings.Split(link, "/")

	partes_total := len(partes)

	var maispartes []string

	if partes_total == 6 {
		maispartes = strings.Split(partes[5], "-")
	} else if partes_total == 5 {
		maispartes = strings.Split(partes[4], "-")
	}

	id, err := strconv.Atoi(maispartes[0])

	if err != nil {
		return 0, err
	}

	return id, nil
}
func (s *NewsService) ReturnPathFromLink(link string) (string, error) {
	if strings.Contains(link, "folha/noticia") {
		return "folha_noticias", nil
	} else if strings.Contains(link, "holofote/noticia") {
		return "holofote_noticias", nil
	} else if strings.Contains(link, "municipios/noticia") {
		return "municipios_noticias", nil
	} else if strings.Contains(link, "saude/noticia") {
		return "saude_noticias", nil
	} else if strings.Contains(link, "justica/noticia") {
		return "justica_noticias", nil
	} else if strings.Contains(link, "bnhall/noticia") {
		return "hall_noticias", nil
	} else if strings.Contains(link, "bnhall/enjoy") {
		return "hall_enjoy", nil
	} else if strings.Contains(link, "esportes/vitoria") {
		return "esportes_vitorias", nil
	} else if strings.Contains(link, "esportes/bahia") {
		return "esportes_bahias", nil
	} else if strings.Contains(link, "esportes/noticia") {
		return "esportes_noticias", nil
	} else {
		return "principal_noticias", nil
	}
}

func (s *NewsService) GetImageLink(link string) (string, error) {
	id, err := s.GetIdFromLink(link)

	if err != nil {
		return "", err
	}

	path, err := s.ReturnPathFromLink(link)

	if err != nil {
		return "", err
	}

	var tag string
	tag = "NOTICIA"
	if path == "esportes_bahias" {
		tag = "BAHIA"
	} else if path == "esportes_vitorias" {
		tag = "VITORIA"
	} else if path == "justica_colunas" || path == "hall_colunas" || path == "holofote_colunas" {
		tag = "COLUNA"
	} else if path == "principal_podcasts" {
		tag = "PODCAST"
	} else if path == "hall_enjoy" {
		tag = "ENJOY"
	} else if path == "hall_travellings" {
		tag = "TRAVELLING"
	} else if path == "hall_business" {
		tag = "BUSINESS"
	}

	newlink := fmt.Sprintf("https://www.bahianoticias.com.br/fotos/%s/%d/IMAGEM_%s_5.jpg", path, id, tag)

	return newlink, nil
}

func (s *NewsService) CopierPage(list_pages []string) []entity.News {
	var lista []entity.News
	for _, page := range list_pages {
		lista_page := s.GetNewsFromPage(page)
		lista = append(lista, lista_page...)
	}

	return lista
}
