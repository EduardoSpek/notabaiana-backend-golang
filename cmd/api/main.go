package main

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	database "github.com/eduardospek/notabaiana-backend-golang/internal/infra/database/postgres"
	"github.com/eduardospek/notabaiana-backend-golang/internal/infra/web"
	"github.com/eduardospek/notabaiana-backend-golang/internal/infra/web/controllers"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
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
	postgres := database.NewPostgres()
	db, _ := postgres.Connect()
	newsrepo := database.NewNewsPostgresRepository(db)
	imagedownloader := utils.NewImgDownloader()	
	news_service := service.NewNewsService(newsrepo, imagedownloader)	

	crawler_service := service.NewCrawler()
	copier_service := service.NewCopier(*news_service, *crawler_service)
	crawler_controller := controllers.NewCrawlerController(*copier_service)

	news_controller := controllers.NewNewsController(*news_service)

	toprepo := database.NewTopPostgresRepository(db)
	top_service := service.NewTopService(toprepo, *news_service)
	top_controller := controllers.NewTopController(*top_service)

	server := web.NewServerWeb()

	server.TopController(*top_controller)
	server.CrawlerController(*crawler_controller)
	server.NewsController(*news_controller)

	//go copier_service.Start("https://www.bahianoticias.com.br/holofote/rss.xml", 22)
	go copier_service.Start("https://www.bahianoticias.com.br/principal/rss.xml", 5)	
	// go copier_service.Start("https://www.bahianoticias.com.br/esportes/rss.xml", 33)
	// go copier_service.Start("https://www.bahianoticias.com.br/justica/rss.xml", 44)
	// go copier_service.Start("https://www.bahianoticias.com.br/hall/rss.xml", 55)
	// go copier_service.Start("https://www.bahianoticias.com.br/saude/rss.xml", 66)
	// go copier_service.Start("https://www.bahianoticias.com.br/municipios/rss.xml", 77)
		
	//Função para gerar as top notícias a cada 60 minutos
	go top_service.Start(60)

	server.Start()

}