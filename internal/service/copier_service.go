package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eduardospek/bn-api/internal/domain/entity"
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

	for _, item := range rss.Channel.Items {
		n := entity.News{
			Title: item.Title,
			Text:  item.Description,
			Link:  item.Link,
			Image: item.Media.URL,
			Visible: true,
		}
		
		html_imagens_anexadas := c.news_service.GetImagesPage(n.Link)

		if html_imagens_anexadas != "" {
			n.Text += "<br><br>"
			n.Text += html_imagens_anexadas
		}

		new, err := c.news_service.CreateNews(n)

		if err == nil {			

			err = c.news_service.SaveImage(new.ID, n.Image, diretorio)

			if err != nil {
				fmt.Println("Erro ao Salvar Image: ", err)
			}			

		}

	}

}