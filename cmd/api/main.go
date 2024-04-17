package main

import (
	"log"

	"github.com/eduardospek/bn-api/internal/controllers"
	database "github.com/eduardospek/bn-api/internal/infra/database/sqlite"
	"github.com/eduardospek/bn-api/internal/infra/web"
	"github.com/eduardospek/bn-api/internal/service"
	"github.com/eduardospek/bn-api/internal/utils"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
        log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
    }  
}

func main() {	

	newsrepo := database.NewNewsSQLiteRepository()
	//newsrepo := database.NewNewsMemoryRepository()
	imagedownloader := utils.NewImgDownloader()	
	news_service := service.NewNewsService(newsrepo, imagedownloader)	

	crawler_service := service.NewCrawler()
	copier_service := service.NewCopier(*news_service, *crawler_service)
	crawler_controller := controllers.NewCrawlerController(*copier_service)

	news_controller := controllers.NewNewsController(*news_service)

	server := web.NewServerWeb()

	server.CrawlerController(*crawler_controller)
	server.NewsController(*news_controller)

	server.Start()

}