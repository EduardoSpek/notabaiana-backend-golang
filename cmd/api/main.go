package main

import (
	"log"

	"github.com/eduardospek/bn-api/internal/controllers"
	database "github.com/eduardospek/bn-api/internal/infra/database/memorydb"
	"github.com/eduardospek/bn-api/internal/infra/web"
	"github.com/eduardospek/bn-api/internal/service"
	"github.com/eduardospek/bn-api/internal/utils"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
        log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
    }  
}

func main() {	

	//userRepo := database.NewUserMysqlRepository()
	newsrepo := database.NewNewsMemoryRepository()
	imagedownloader := utils.NewImgDownloader()	
	news_service := service.NewNewsService(newsrepo, imagedownloader)	

	crawler_service := service.NewCrawler()
	crawler_controller := controllers.NewCrawlerController(*news_service, *crawler_service)

	server := web.NewServerWeb()

	server.CrawlerController(*crawler_controller)

	server.Start()

}