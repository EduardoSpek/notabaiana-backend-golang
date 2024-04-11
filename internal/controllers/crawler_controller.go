package controllers

import (
	"net/http"
	"os"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	"github.com/eduardospek/bn-api/internal/service"
)

type CrawlerController struct {
	news_service service.NewsService
	crawler_service service.CrawlerService
}

func NewCrawlerController(newsservice service.NewsService, crawlerservice service.CrawlerService) *CrawlerController {
	return &CrawlerController{ news_service: newsservice, crawler_service: crawlerservice }
}

func (c *CrawlerController) Crawler(w http.ResponseWriter, r http.Request) {

	rss := c.crawler_service.GetRSS(os.Getenv("URL_RSS"))

	for _, item := range rss.Channel.Items {
		n := entity.News{
			Title: item.Title,
			Text: item.Description,
			Link: item.Link,
			Image: item.Media.URL,
		}
		
		err := c.news_service.CreateNews(n)

		if err != nil {
			continue
		}
	}

	msg := map[string]any{
		"ok": true,
		"message": "Notícias resgatadas!",

	}
	ResponseJson(w, msg, http.StatusNotFound)
	
}