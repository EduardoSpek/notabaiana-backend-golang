package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

type CopierService struct {
	news_service    *NewsService
	crawler_service *CrawlerService
}

func NewCopier(newsservice *NewsService, crawlerservice *CrawlerService) *CopierService {
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

	lista := c.news_service.CopierPage(list_pages)

	for _, n := range lista {

		go func() {

			new, err := c.news_service.CreateNews(n)

			if err == nil {

				if new.ID != "" {

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
			}
		}()

	}
}
