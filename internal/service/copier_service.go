package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/eduardospek/bn-api/internal/domain/entity"
)

type CopierService struct {
	news_service    NewsService
	crawler_service CrawlerService
}

func NewCopier(newsservice NewsService, crawlerservice CrawlerService) *CopierService {
	return &CopierService{news_service: newsservice, crawler_service: crawlerservice}
}

func (c *CopierService) Start() {

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

	rss := c.crawler_service.GetRSS(os.Getenv("URL_RSS"))

	for _, item := range rss.Channel.Items {
		n := entity.News{
			Title: item.Title,
			Text:  item.Description,
			Link:  item.Link,
			Image: item.Media.URL,
		}		

		new, err := c.news_service.CreateNews(n)

		if err != nil {
			fmt.Println("Erro ao Salvar News: ", err)
		} else {

			err = c.news_service.SaveImage(new.ID, n.Image, diretorio)

			if err != nil {
				fmt.Println("Erro ao Salvar Image: ", err)
			}

		}

	}

}