package main

import (
	"github.com/eduardospek/bn-api/internal/controllers"
	database "github.com/eduardospek/bn-api/internal/infra/database/supabase"
	"github.com/eduardospek/bn-api/internal/infra/web"
	"github.com/eduardospek/bn-api/internal/service"
	"github.com/eduardospek/bn-api/internal/utils"
)

// func init() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
//         log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
//     }
// }

func main() {	

	//newsrepo := database.NewNewsSQLiteRepository()
	//newsrepo := database.NewNewsMemoryRepository()
	supabase := database.NewSupabase()
	db, _ := supabase.Connect()
	newsrepo := database.NewNewsSupabaseRepository(db)
	imagedownloader := utils.NewImgDownloader()	
	news_service := service.NewNewsService(newsrepo, imagedownloader)	

	crawler_service := service.NewCrawler()
	copier_service := service.NewCopier(*news_service, *crawler_service)
	crawler_controller := controllers.NewCrawlerController(*copier_service)

	news_controller := controllers.NewNewsController(*news_service)

	server := web.NewServerWeb()

	server.CrawlerController(*crawler_controller)
	server.NewsController(*news_controller)

	go copier_service.Start("https://www.bahianoticias.com.br/principal/rss.xml", 10)
	go copier_service.Start("https://www.bahianoticias.com.br/holofote/rss.xml", 20)
	go copier_service.Start("https://www.bahianoticias.com.br/esportes/rss.xml", 30)
	go copier_service.Start("https://www.bahianoticias.com.br/justica/rss.xml", 40)	
	server.Start()

}