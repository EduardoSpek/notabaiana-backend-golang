package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/gocolly/colly"
)

var (
		ErrWordsBlackList = errors.New("o título contém palavras bloqueadas")
		ErrNoCategory = errors.New("nenhuma categoria no rss")
		ErrSimilarTitle = errors.New("título similar ao recente adicionado detectado")
		AllowedDomains = "www.bahianoticias.com.br"		

		LimitPerPage = 100
	)

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Contents Content `json:"content"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}
		
type NewsService struct {
	newsrepository port.NewsRepository
	imagedownloader port.ImageDownloader
}

func NewNewsService(repository port.NewsRepository, downloader port.ImageDownloader) *NewsService {
	return &NewsService{ newsrepository: repository, imagedownloader: downloader }
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

	//newtitle, err := s.ChangeTitleWithGemini(new.Title)

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
	
	new, err := s.newsrepository.GetBySlug(slug)

	if err != nil {
		return entity.News{}, err
	}
	
	return new, nil

}

func (s *NewsService) SearchNews(page int, str_search string) interface{} {

	str_search = strings.Replace(str_search, " ", "%", -1)
	
	news := s.newsrepository.SearchNews(page, str_search)

	total := s.newsrepository.GetTotalNewsBySearch(str_search)

	pagination := s.Pagination(page, total)

    result := struct{
        List_news []entity.News `json:"news"`
        Pagination map[string][]int `json:"pagination"`
        Search string `json:"search"`
    }{
        List_news: news,
        Pagination: pagination,
        Search: strings.Replace(str_search, "%", " ", -1),
    }
	
	return result

}

func (s *NewsService) FindAllNews(page, limit int) interface{} {

	//Limita o total de registros que deve ser retornado
	if limit > LimitPerPage { limit = LimitPerPage }
	
	news, _ := s.newsrepository.FindAll(page, limit)

	total := s.newsrepository.GetTotalNews()

	pagination := s.Pagination(page, total)

    result := struct{
        List_news []entity.News `json:"news"`
        Pagination map[string][]int `json:"pagination"`
    }{
        List_news: news,
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

	pagination := s.Pagination(page, total)

    result := struct{
        List_news []entity.News `json:"news"`
        Pagination map[string][]int `json:"pagination"`
		Category string `json:"category"`
    }{
        List_news: news,
        Pagination: pagination,
		Category: category,
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

	err = s.imagedownloader.ResizeAndSaveImage(img, width, height, outputPath)
	
	if err != nil {
		//fmt.Println("Erro ao redimensionar e salvar a imagem:", err)
		return "",err
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
	palavras := []string {
		"Bahia Notícias",
		"BN",
		"Curtas",
		"Nota Baiana",
		"NotaBaiana",
		"notabaiana",
	}
    for _, palavra := range palavras {
        if strings.Contains(titulo, palavra) {
            return true
        }
    }
    return false
}
func changeWords(text string) string {
	text = strings.Replace(text, "Siga o @bnhall_ no Instagram e fique de olho nas principais notícias.", "", -1)

	text = strings.Replace(text, "Bahia Notícias", "BN", -1)

	text = strings.Replace(text, "@BahiaNoticias", "@", -1)

	text = strings.Replace(text, "@bnholofote", "@", -1)

	text = strings.Replace(text, "BN Holofote", "", -1)

	text = strings.Replace(text, "@bhall", "@", -1)

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
			Title: item,
			Text: texts[i],
			Image: images[i],
			Link: links[i],
			Visible: true,
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
	} else if path == "justica_colunas" || path == "hall_colunas" {
		tag = "COLUNA"
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

func (s *NewsService)  ChangeTitleWithGemini(title string) (string, error) {

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + os.Getenv("KEY_GEMINI")

	title = strings.ReplaceAll(title, `"`, `\"`)
	title = strings.ReplaceAll(title, `'`, `\'`)

	jsonData := `{"contents":[{"parts":[{"text":"Matenha o contexto e refaça o título usando palavras-chaves para melhorar o SEO. Tente gerar curiosidade para o leitor querer ler a notícia completa. Retorne apenas o título com no máximo 120 caracteres. O título para ser refeito é: ` + title + `"}]}]}`

	

	reqBody := bytes.NewBuffer([]byte(jsonData))

	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição POST:", err)
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return "", err
	}
	
	// Deserializa o JSON na struct Response
	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Println("Erro ao deserializar o JSON:", err)
		return "", err
	}	

	var newtitle string
	// Acessa o valor de "text"
	for _, candidate := range response.Candidates {
		for _, part := range candidate.Contents.Parts {
						
			newtitle = strings.Replace(part.Text, "**", "", -1)
			
		}
	}

	return newtitle, nil

}

//Pagination recebe a página atual e o total de noticias para retornar a páginação de resultado
func (s *NewsService) Pagination(currentPage, totalNews int) map[string][]int {
	
	// Calcula o total de páginas
	totalPages := int(math.Ceil(float64(totalNews) / 10)) 

	// Garante que a página atual esteja dentro dos limites
	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}

	previousPages := []int{}
	nextPages := []int{}

	if currentPage == totalPages {
		if currentPage > 2 {
			previousPages = []int{currentPage - 2, currentPage - 1}
			nextPages = []int{}
		} else {
			previousPages = []int{currentPage - 1}
			nextPages = []int{}
		}
	} else if currentPage-2 > 2 && currentPage+2 <= totalPages {
		previousPages = []int{currentPage - 2, currentPage - 1}
		nextPages = []int{currentPage + 1, currentPage + 2}
	} else if currentPage == 1 && currentPage == totalPages {
		previousPages = []int{}
		nextPages = []int{}
	} else if currentPage == 1 && totalPages < 3 {
		previousPages = []int{}
		nextPages = []int{currentPage + 1}
	} else if currentPage == 1 && totalPages > 2 {
		previousPages = []int{}
		nextPages = []int{currentPage + 1, currentPage + 2}
	} else if currentPage == 2 && currentPage == totalPages {
		previousPages = []int{1}
		nextPages = []int{}
	} else if currentPage == 2 && totalPages < 4 {
		previousPages = []int{currentPage - 1}
		nextPages = []int{currentPage + 1}
	} else if currentPage == 2 && totalPages > 3 {
		previousPages = []int{currentPage - 1}
		nextPages = []int{currentPage + 1, currentPage + 2}
	} else if currentPage == 3 && currentPage == totalPages {
		previousPages = []int{1, 2}
		nextPages = []int{}
	} else if currentPage == 3 && totalPages < 5 {
		previousPages = []int{1, 2}
		nextPages = []int{currentPage + 1}
	} else if currentPage == 3 && totalPages > 4 {
		previousPages = []int{1, 2}
		nextPages = []int{currentPage + 1, currentPage + 2}
	} else if currentPage == 4 && currentPage == totalPages {
		previousPages = []int{currentPage - 2, currentPage - 1}
		nextPages = []int{}
	} else if currentPage == 4 && totalPages < 6 {
		previousPages = []int{currentPage - 2, currentPage - 1}
		nextPages = []int{currentPage + 1}
	} else if currentPage == 4 && totalPages > 5 {
		previousPages = []int{currentPage - 2, currentPage - 1}
		nextPages = []int{currentPage + 1, currentPage + 2}
	} else if currentPage == 5 && totalPages < 7 {
		previousPages = []int{currentPage - 2, currentPage - 1}
		nextPages = []int{currentPage + 1}
	} else if currentPage > 3 && totalPages > currentPage && currentPage+1 <= totalPages {
		previousPages = []int{currentPage - 2, currentPage - 1}
		nextPages = []int{currentPage + 1}
	}

	result := map[string][]int{
		"previousPages": previousPages,
		"currentPage":   {currentPage},
		"nextPages":     nextPages,
		"totalPages":    {totalPages},
	}

	return result
}