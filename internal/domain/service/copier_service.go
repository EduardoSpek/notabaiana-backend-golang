package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

type CopierService struct {
	news_service    NewsService
	crawler_service CrawlerService
}

func NewCopier(newsservice NewsService, crawlerservice CrawlerService) *CopierService {
	return &CopierService{news_service: newsservice, crawler_service: crawlerservice}
}

func (c *CopierService) Start(rss []string, minutes time.Duration) {

	go c.Run(rss)

	ticker := time.NewTicker(minutes * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		go c.Run(rss)
	}
}

func (c *CopierService) Run(list_pages []string) {

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("Erro ao obter o caminho do executável:", err)
	}

	diretorio := strings.Replace(cwd, "test", "", -1) + "/images/"

	// rss := c.crawler_service.GetRSS(rss_url)
	// category, _ := c.news_service.GetCategory(rss_url)

	//var page string

	// for _, item := range rss.Channel.Items {
	// 	n := entity.News{
	// 		Title: item.Title,
	// 		Text:  item.Description,
	// 		Link:  item.Link,
	// 		Image: item.Media.URL,
	// 		Visible: true,
	// 		Category: category,
	// 	}

	// 	lista = append(lista, n)
	// }

	lista := c.news_service.CopierPage(list_pages)

	for _, n := range lista {

		go func() {
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

				fileExists := utils.FileExsists(outputPath)

				if !fileExists {
					err := c.news_service.ClearImagePath(new.ID)

					if err != nil {
						fmt.Println("não foi possível atualizar o caminho da imagem")
					}
				}
			}
		}()
	}
}
