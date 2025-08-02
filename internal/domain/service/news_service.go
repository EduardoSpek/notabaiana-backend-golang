package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/config"
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
	ErrWordsBlackList = errors.New("o título ou texto contém palavras bloqueadas")
	ErrNoCategory     = errors.New("nenhuma categoria no rss")
	ErrSimilarTitle   = errors.New("título similar ao recente adicionado detectado")
)

type Noticia struct {
	Titulo string `json:"titulo"`
	Texto  string `json:"texto"`
}

type FindAllOutput struct {
	List_news  []*entity.NewsFindAllOutput `json:"news"`
	Pagination map[string][]int            `json:"pagination"`
}

type NewsService struct {
	hitsrepository  port.HitsRepository
	newsrepository  port.NewsRepository
	imagedownloader port.ImageDownloader
}

func NewNewsService(repository port.NewsRepository, downloader port.ImageDownloader, hits port.HitsRepository) *NewsService {
	return &NewsService{newsrepository: repository, imagedownloader: downloader, hitsrepository: hits}
}

func (s *NewsService) AdminDeleteAll(news []*entity.News) error {

	var listNews []*entity.News
	for _, n := range news {
		new, err := s.newsrepository.GetByID(n.ID)

		if err != nil {
			return err
		}

		listNews = append(listNews, new)
	}

	s.RemoveImages(listNews)

	err := s.newsrepository.DeleteAll(news)

	if err != nil {
		return err
	}

	return nil
}

func (s *NewsService) Delete(id string) error {

	news, err := s.newsrepository.GetByID(id)

	if err != nil {
		return err
	}

	removed := utils.RemoveImage("./images/" + news.Image)

	if !removed {
		fmt.Println("Delete News: não foi possível deletar a imagem")
	}

	err = s.newsrepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

// StartCleanNews remove news with false status
func (s *NewsService) StartCleanNews(minutes time.Duration) {

	go s.CleanNews()

	ticker := time.NewTicker(minutes * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		go s.CleanNews()
	}
}

func (s *NewsService) CleanNews() error {

	news, err := s.newsrepository.CleanNews()

	if err != nil {
		return err
	}

	s.RemoveImages(news)

	err = s.newsrepository.DeleteAll(news)

	if err != nil {
		return err
	}

	return nil

}

func (s *NewsService) StartCleanNewsOld(minutes time.Duration) {

	go s.CleanNewsOld()

	ticker := time.NewTicker(minutes * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		go s.CleanNewsOld()
	}
}

func (s *NewsService) CleanNewsOld() error {

	news, err := s.newsrepository.CleanNewsOld()

	if err != nil {
		return err
	}

	s.RemoveImages(news)

	err = s.newsrepository.DeleteAll(news)

	if err != nil {
		return err
	}

	return nil

}

func (s *NewsService) StartScanDuplicateNews(ctx context.Context, minutes time.Duration) {

	go s.ScanDuplicateNews(ctx)

	ticker := time.NewTicker(minutes * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		go s.ScanDuplicateNews(ctx)
	}
}

func (s *NewsService) ScanDuplicateNews(ctx context.Context) error {

	var listNewsDelete []*entity.News

	news := s.FindAllNews(ctx, 1, 100)

	newsList := news.(*FindAllOutput)

	newsCopy := newsList.List_news

	for _, n := range newsList.List_news {

		newsInList := false
		for _, nd := range listNewsDelete {
			if n.ID == nd.ID {
				newsInList = true
			}
		}

		if newsInList {
			continue
		}

		for _, n2 := range newsCopy {

			if n.ID == n2.ID {
				continue
			}

			newsInList := false
			for _, nd := range listNewsDelete {
				if n2.ID == nd.ID {
					newsInList = true
				}
			}

			if newsInList {
				continue
			}

			similarity := utils.Similarity(n.Title, n2.Title)

			if similarity > 60 {
				newsInsert := &entity.News{
					ID: n2.ID,
				}
				listNewsDelete = append(listNewsDelete, newsInsert)
			}
		}

	}

	if len(listNewsDelete) > 0 {
		for _, n := range listNewsDelete {
			s.newsrepository.SetVisible(false, n.ID)
		}
	}

	return nil

}

func (s *NewsService) RemoveImages(news []*entity.News) {

	for _, n := range news {

		if n.Image != "" {
			image := "./images/" + n.Image
			utils.RemoveImage(image)
		}
	}
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

func (s *NewsService) UpdateNewsUsingTheForm(file multipart.File, newsInput *entity.News) (*entity.News, error) {

	oldnew, err := s.newsrepository.AdminGetBySlug(newsInput.Slug)

	if err != nil {
		return &entity.News{}, err
	}

	oldnew.Title = newsInput.Title
	oldnew.TitleAi = newsInput.TitleAi
	oldnew.Text = newsInput.Text
	oldnew.Visible = newsInput.Visible
	oldnew.TopStory = newsInput.TopStory
	oldnew.Category = newsInput.Category

	newNews := entity.UpdateNews(oldnew)

	news := ChangeLink(newNews)

	if file != nil {
		news = RenamePathImage(news)
	}

	new, err := s.newsrepository.Update(news)

	if err != nil {
		return &entity.News{}, err
	}

	err = s.SaveImageForm(file, new)

	if err != nil {
		new.Image = ""
	}

	return new, nil

}

func (s *NewsService) CreateNewsUsingTheForm(file multipart.File, news *entity.News) (*entity.News, error) {

	newNews := entity.NewNews(news)

	new, err := s.CreateNews(newNews)

	if err != nil {
		return &entity.News{}, err
	}

	err = s.SaveImageForm(file, new)

	if err != nil {
		newNews.Image = ""
	}

	return new, nil

}

func (s *NewsService) SaveImageForm(file multipart.File, news *entity.News) error {

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

func (s *NewsService) CreateNews(news *entity.News) (*entity.News, error) {

	new := entity.NewNews(news)

	recent, err := s.newsrepository.FindRecent()

	if err != nil {
		return &entity.News{}, err
	}

	similarity := utils.Similarity(recent.Title, new.Title)

	if similarity > 60 {
		//fmt.Println("***titulo silimar detectado***")
		return &entity.News{}, ErrSimilarTitle
	}

	embed, text := s.GetEmded(new.Link)

	if text != "" {
		new.Text = text
	}

	new = RenamePathImage(new)
	new = ChangeLink(new)

	result := listOfBlockedWords(new.Title)

	if result {
		return &entity.News{}, ErrWordsBlackList
	}

	new.Text = changeWords(new.Text)

	err = s.newsrepository.NewsExists(new.Title)

	if err != nil {
		return &entity.News{}, err
	}

	newtext, err := utils.ChangeTitleWithGemini("Você é um jornalista, refaça este texto matendo o contexto. Mantenha os assuntos principais. Baseado no texto, crie um titulo para a notícia seguindo boas práticas de SEO. Retorne o título e o texto em formato JSON. O texto é: ", new.Text)

	if err != nil {
		fmt.Println("Erro ao obter o texto do gemini:", err)
	}

	var noticia Noticia
	err = json.Unmarshal([]byte(newtext), &noticia)
	if err == nil {
		if noticia.Titulo != "" {
			new.TitleAi = strings.TrimSpace(noticia.Titulo)
		}
		if noticia.Texto != "" { 
			new.Text = strings.TrimSpace(noticia.Texto)
		}
	}

	if embed != "" {
		new.Text += "<br><br>"
		new.Text += embed
	}

	result = listOfBlockedText(new.Text)

	if result {
		return &entity.News{}, ErrWordsBlackList
	}

	_, err = s.newsrepository.Create(new)

	if err != nil {
		return &entity.News{}, err
	}

	return new, nil
}

func (s *NewsService) AdminGetNewsBySlug(slug string) (*entity.News, error) {

	new, err := s.newsrepository.AdminGetBySlug(slug)

	if err != nil {
		return &entity.News{}, err
	}

	return new, nil

}

func (s *NewsService) GetNewsBySlug(ctx context.Context, slug string) (*entity.News, error) {

	// err := s.Hit(slug)

	// if err != nil {
	// 	return &entity.News{}, err
	// }

	new, err := s.newsrepository.GetBySlug(ctx, slug)

	if err != nil {
		return &entity.News{}, err
	}

	return new, nil

}

func (s *NewsService) Hit(r *http.Request, session string) error {

	ip := utils.GetIP(r)

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

func (s *NewsService) SearchNews(page int, str_search string) (interface{}, error) {

	var newsOutput []*entity.NewsFindAllOutput

	str_search = strings.Replace(str_search, " ", "%", -1)

	news, err := s.newsrepository.SearchNews(page, str_search)

	if err != nil {
		return nil, err
	}

	newsOutput = s.NewsConvertListOutput(news)

	total := s.newsrepository.GetTotalNewsBySearch(str_search)

	pagination := utils.Pagination(page, config.News_PerPage, total)

	result := struct {
		List_news  []*entity.NewsFindAllOutput `json:"news"`
		Pagination map[string][]int            `json:"pagination"`
		Search     string                      `json:"search"`
	}{
		List_news:  newsOutput,
		Pagination: pagination,
		Search:     strings.Replace(str_search, "%", " ", -1),
	}

	return result, nil

}

func (s *NewsService) AdminFindAllNews(page, limit int) interface{} {

	var newsOutput []*entity.NewsFindAllOutput

	//Limita o total de registros que deve ser retornado
	if limit > config.News_LimitPerPage {
		limit = config.News_LimitPerPage
	}

	news, _ := s.newsrepository.AdminFindAll(page, limit)

	newsOutput = s.NewsConvertListOutput(news)

	total := s.newsrepository.GetTotalNews()

	pagination := utils.Pagination(page, limit, total)

	result := &FindAllOutput{
		List_news:  newsOutput,
		Pagination: pagination,
	}

	return result

}

func (s *NewsService) FindAllNews(ctx context.Context, page, limit int) interface{} {

	var newsOutput []*entity.NewsFindAllOutput

	//Limita o total de registros que deve ser retornado
	if limit > config.News_LimitPerPage {
		limit = config.News_LimitPerPage
	}

	news, _ := s.newsrepository.FindAll(ctx, page, limit)

	newsOutput = s.NewsConvertListOutput(news)

	total := s.newsrepository.GetTotalNewsVisible()

	pagination := utils.Pagination(page, limit, total)

	result := &FindAllOutput{
		List_news:  newsOutput,
		Pagination: pagination,
	}

	return result

}

func (s *NewsService) FindRecent() (*entity.News, error) {

	news, err := s.newsrepository.FindRecent()

	if err != nil {
		return &entity.News{}, err
	}

	return news, nil

}

func (s *NewsService) FindNewsCategory(category string, page int) interface{} {

	var newsOutput []*entity.NewsFindAllOutput

	news, _ := s.newsrepository.FindCategory(category, page)

	newsOutput = s.NewsConvertListOutput(news)

	total := s.newsrepository.GetTotalNewsByCategory(category)

	pagination := utils.Pagination(page, config.News_PerPage, total)

	result := struct {
		List_news  []*entity.NewsFindAllOutput `json:"news"`
		Pagination map[string][]int            `json:"pagination"`
		Category   string                      `json:"category"`
	}{
		List_news:  newsOutput,
		Pagination: pagination,
		Category:   category,
	}

	return result
}

func (s *NewsService) FindAllViews() ([]*entity.News, error) {

	news, err := s.newsrepository.FindAllViews()

	if err != nil {
		return []*entity.News{}, err
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
	var htmlx, conteudo, script, text, str_text string
	//var err error

	collector := colly.NewCollector(
		colly.AllowedDomains(config.News_AllowedDomains),
	)

	//Obtém o texto da notícia
	collector.OnHTML(".sc-bcb60b9a-3.jSwwqa", func(e *colly.HTMLElement) {

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

		htmlx += "<br><br>" + conteudo_decoded

	})

	//Obtém os scrits
	collector.OnHTML(".lazyload-scripts", func(e *colly.HTMLElement) {
		// Obter o valor do atributo "src" da imagem
		script = e.Attr("data-scripts")

		script_decoded, err := url.QueryUnescape(script)

		if err != nil {
			return
		}

		if strings.Contains(script_decoded, "youtube") || strings.Contains(script_decoded, "flickr") || strings.Contains(script_decoded, "tiktok") || strings.Contains(script_decoded, "bsky") {

			if !strings.Contains(htmlx, script_decoded) {
				htmlx += script_decoded
			}
		}

	})

	collector.OnHTML("iframe", func(e *colly.HTMLElement) {
		// Obter o atributo src do iframe
		src := e.Attr("src")

		// Verificar se o src contém a URL de embed do YouTube
		if strings.Contains(src, "youtube.com/embed") {
			// Alterar os atributos width e height
			e.DOM.SetAttr("width", "100%")
			e.DOM.SetAttr("height", "320")

			// Obter o nó HTML do iframe
			nodes := e.DOM.Nodes
			if len(nodes) == 0 {
				fmt.Println("Nenhum nó encontrado para o iframe")
				return
			}
			node := nodes[0]

			// Construir o HTML do iframe manualmente
			var attrStr strings.Builder
			for _, attr := range node.Attr {
				// Garantir que width e height sejam os valores modificados
				switch attr.Key {
				case "width":
					attrStr.WriteString(` width="100%"`)
				case "height":
					attrStr.WriteString(` height="320"`)
				default:
					attrStr.WriteString(fmt.Sprintf(` %s="%s"`, attr.Key, html.EscapeString(attr.Val)))
				}
			}

			// Montar o HTML completo do iframe
			iframeHTML := fmt.Sprintf("<iframe%s></iframe>", attrStr.String())

			// Adicionar o código HTML à lista
			htmlx += `<div class="imagem_anexada">` + iframeHTML + `</div>`
		}
	})

	// Definindo o callback OnHTML
	collector.OnHTML("img", func(e *colly.HTMLElement) {
		// Obter o valor do atributo "src" da imagem
		src := e.Attr("src")

		src = strings.Replace(src, " ", "%20", -1)

		// Verificar se o valor "src" contém o endereço de destino
		if strings.Contains(src, "bahianoticias.com.br/fotos/") {

			htmlx += `<div class="imagem_anexada"><img src="` + src + `" width="100%"></div>`

		}
	})

	// Visitando a URL inicial
	collector.Visit(link)

	//html = text + htmlx

	return htmlx, text
}

func RenamePathImage(news *entity.News) *entity.News {
	news.Image = news.ID + ".jpg"
	return news
}
func ChangeLink(news *entity.News) *entity.News {
	news.Link = "/news/" + news.Slug
	return news
}
func listOfBlockedText(text string) bool {
	words := []string{
		"IMG_OFER_0",
		"JusPod",
		"BN na Bola",
		"NB na Bola",
	}

	for _, word := range words {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
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
		"cassino",
		"Cassino",
		"Previdência",
		"Aposta",
		"aposta",
		"Arena Esportiva",
		"apostar",
		"aposte",
		"palpite",
		"bônus",
		"Prisma",
		"poker",
		"pôker",
		"Pôker",
		"Poker",
		"candidat",
		"Candidat",
		"Krypto",
		"trading",
		"krypto",
		"Trading",
		"Opinião",
		"Caetano",
		"Preta Gil",
		"Gilberto Gil",
		"Nininha",
		"Travesti",
		"LGBT",
		"Vittar",
		"Anitta",
		"Luísa Sonza",
		"Luisa Sonza",
		"Kret",
		"Leitte",
		"A Dama",
		"Davi",
		"Maíra Cardi",
		"Thiago Nigro",
		"Drag",
		"Queen",
		"homof",
		"estupr",
		"violen",
		"Improta",
		"Duquesa",
		"Liniker",
		"porn",
		"Bolsonaro",
		"Paulo Vieira",
		"Fábio Porchat",
		"Erika Hilton",
		"Felipe Neto",
		"Mercury",
		"Gabriela Prioli",
		"Ludmilla",
		"Thaís Carla",
		"Thais Carla",
		"TheLotter",
		"Powerball",
		"Alô Juca",
		"Linn da Quebrada",
		"Patrícia Ramos",
		"Blackjack",
		"Bruna Louise",
		"Giovanna Ewbank",
		"Bruno Gagliasso",
		"Gkay",
		"Flora Gil",
		"Eliezer",
		"Viih Tube",
		"Camila Loures",
		"Luana Piovani",
		"MC Rebecca",
		"Brunna Gonçalves",
		"Ana Paula Renault",
		"Kart Love",
		"Juliette",
		"Oruam",
		"MC Cabelinho",
		"Poze",
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

	text = strings.Replace(text, "Bahia Notícias", "BN", -1)

	text = strings.Replace(text, "Bahia Notícia", "BN", -1)

	text = strings.Replace(text, "Bahia Noticia", "BN", -1)

	text = strings.Replace(text, "Bahia Noticias", "BN", -1)

	text = strings.Replace(text, "@BahiaNoticias", "@notabaiana", -1)

	text = strings.Replace(text, "@bnholofote", "@notabaiana", -1)

	text = strings.Replace(text, "BN Holofote", "NotaBaiana", -1)

	text = strings.Replace(text, "@bhall", "@notabaiana", -1)

	text = strings.Replace(text, "As informações são do Metrópoles, parceiro do BN", ".", -1)

	text = strings.Replace(text, " parceiro do BN,", "", -1)

	text = strings.Replace(text, "Assine a newsletter de Esportes do BN e fique bem informado sobre o esporte na Bahia, no Brasil e no mundo!", "", -1)

	text = strings.Replace(text, "Siga o BN no Google News e veja os conteúdos de maneira ainda mais rápida e ágil pelo celular ou pelo computador!", "", -1)

	text = strings.Replace(text, `<img src="https://www.bahianoticias.com.br/fotos/oferecimentos/30/IMG_OFER_0.jpg" width="100%">`, "", -1)

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

func (s *NewsService) GetNewsFromPage(link string) []*entity.News {
	var conteudo string
	//var err error

	collector := colly.NewCollector(
		colly.AllowedDomains(config.News_AllowedDomains),
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
	collector.OnHTML(".sc-24c322fd-3.giHdVV", func(e *colly.HTMLElement) {

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

	var lista []*entity.News

	if len(texts) == 0 {
		fmt.Println("Não foi possível obter os textos das notícias")
		return nil
	}

	for i, item := range titulos {

		category, _ := s.GetCategory(links[i])

		new := entity.NewNews(&entity.News{
			Title:    item,
			Text:     texts[i],
			Image:    images[i],
			Link:     links[i],
			Visible:  true,
			TopStory: false,
			Category: category,
		})
		lista = append(lista, new)
	}

	return lista
}

func (s *NewsService) GetIdFromLink(link string) (int, error) {

	partes := strings.Split(link, "/")

	partes_total := len(partes)

	var maispartes []string

	switch partes_total {
	case 6:
		maispartes = strings.Split(partes[5], "-")
	case 5:
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

	switch path {
	case "esportes_bahias":
		tag = "BAHIA"
	case "esportes_vitorias":
		tag = "VITORIA"
	case "justica_colunas":
		tag = "COLUNA"
	case "hall_colunas":
		tag = "COLUNA"
	case "holofote_colunas":
		tag = "COLUNA"
	case "principal_podcasts":
		tag = "PODCAST"
	case "hall_enjoy":
		tag = "ENJOY"
	case "hall_travellings":
		tag = "TRAVELLING"
	case "hall_business":
		tag = "BUSINESS"
	default:
		tag = "NOTICIA"
	}

	newlink := fmt.Sprintf("https://www.bahianoticias.com.br/fotos/%s/%d/IMAGEM_%s_5.jpg", path, id, tag)

	return newlink, nil
}

func (s *NewsService) CopierPage(list_pages []string) []*entity.News {
	var lista []*entity.News
	for _, page := range list_pages {
		lista_page := s.GetNewsFromPage(page)
		lista = append(lista, lista_page...)
	}

	return lista
}

func (s *NewsService) NewsConvertListOutput(news []*entity.News) []*entity.NewsFindAllOutput {

	var newsOutput []*entity.NewsFindAllOutput

	for _, n := range news {
		newsOutput = append(newsOutput, &entity.NewsFindAllOutput{
			ID:        n.ID,
			Title:     n.Title,
			TitleAi:   n.TitleAi,
			Text:      n.Text,
			Image:     n.Image,
			Link:      n.Link,
			Slug:      n.Slug,
			CreatedAt: n.CreatedAt,
		})
	}

	return newsOutput
}
