package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	"github.com/eduardospek/bn-api/internal/utils"
)

type CopierService struct {
	news_service    NewsService
	crawler_service CrawlerService
}

func NewCopier(newsservice NewsService, crawlerservice CrawlerService) *CopierService {
	return &CopierService{news_service: newsservice, crawler_service: crawlerservice}
}

func (c *CopierService) Start(rss string, minutes time.Duration) {
	
	go c.Run(rss)

	ticker := time.NewTicker(minutes * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
		go c.Run(rss)
	}
}

func (c *CopierService) Run(rss_url string) {

	err := os.MkdirAll("images", os.ModePerm)
	if err != nil {
		fmt.Println("Erro ao criar pasta:", err)
		return
	}

	cwd, err := os.Getwd()
	diretorio := strings.Replace(cwd, "test", "", -1) + "/images/"

	if err != nil {
		fmt.Println("Erro ao obter o caminho do executável:", err)
	}

	rss := c.crawler_service.GetRSS(rss_url)
	category, _ := c.news_service.GetCategory(rss_url)

	var lista []entity.News
	var page string

	for _, item := range rss.Channel.Items {
		n := entity.News{
			Title: item.Title,
			Text:  item.Description,
			Link:  item.Link,
			Image: item.Media.URL,
			Visible: true,
			Category: category,
		}
		
		lista = append(lista, n)
	}

	page = "https://www.bahianoticias.com.br/holofote"
	lista_page := c.news_service.GetNewsFromPage(page)
	lista = append(lista, lista_page...)

	for _, n := range lista {

		embed, text := c.news_service.GetEmded(n.Link)

		if text != "" {
			n.Text = text			
		}	
		
		if embed != "" {
			n.Text += "<br><br>"
			n.Text += embed
		}		

		new, err := c.news_service.CreateNews(n)

		if err == nil {			

			outputPath, err := c.news_service.SaveImage(new.ID, n.Image, diretorio)

			if err != nil {
				fmt.Println("Erro ao Salvar Image: ", err)
			}

			// if strings.Contains(new.Title, "hackeada") {
			// 	err = os.Remove(outputPath)
			// 	if err != nil {
			// 		fmt.Println("não foi possível remover o rquivo")
			// 	}
			// }	
			
			fileExists := utils.FileExsists(outputPath)

			if !fileExists {
				err := c.news_service.ClearImagePath(new.ID)

				if err != nil {
					fmt.Println("não foi possível atualizar o caminho da imagem")
				}
			}
			

		}

	}



}