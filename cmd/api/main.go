package main

import (
	"fmt"
	"os"

	"github.com/eduardospek/notabaiana-backend-golang/internal/adapter"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	database "github.com/eduardospek/notabaiana-backend-golang/internal/infra/database/postgres"
	"github.com/eduardospek/notabaiana-backend-golang/internal/infra/notifications"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

func init() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	// }

	err := os.MkdirAll("files", os.ModePerm)
	if err != nil {
		fmt.Println("Erro ao criar pasta:", err)
		return
	}

	err = os.MkdirAll("images", os.ModePerm)
	if err != nil {
		fmt.Println("Erro ao criar pasta:", err)
	}

	err = os.MkdirAll("images/banners", os.ModePerm)
	if err != nil {
		fmt.Println("Erro ao criar pasta:", err)
	}

	err = os.MkdirAll("images/contacts", os.ModePerm)
	if err != nil {
		fmt.Println("Erro ao criar pasta:", err)
	}

	err = os.MkdirAll("images/downloads", os.ModePerm)
	if err != nil {
		fmt.Println("Erro ao criar pasta:", err)
	}
}

var (
	list_pages = []string{
		"https://www.bahianoticias.com.br",
		"https://www.bahianoticias.com.br/holofote",
		"https://www.bahianoticias.com.br/esportes",
		"https://www.bahianoticias.com.br/justica",
		"https://www.bahianoticias.com.br/saude",
		"https://www.bahianoticias.com.br/municipios",
	}
	list_downloads = []string{
		"https://suamusica.com.br/_next/data/webid-1017/pt-BR/categorias/pagode.json?category=pagode",
		"https://suamusica.com.br/_next/data/webid-1017/pt-BR/categorias/swingueira.json?category=swingueira",
		"https://suamusica.com.br/_next/data/webid-1017/pt-BR/categorias/samba.json?category=samba",
		"https://suamusica.com.br/_next/data/webid-1017/pt-BR/categorias/forro.json?category=forro",
		"https://suamusica.com.br/_next/data/webid-1017/pt-BR/categorias/axe.json?category=axe",
		"https://suamusica.com.br/_next/data/webid-1017/pt-BR/categorias/arrocha.json?category=arrocha",
		"https://suamusica.com.br/_next/data/webid-1017/pt-BR/categorias/piseiro.json?category=piseiro",
		"https://suamusica.com.br/_next/data/webid-1017/pt-BR/categorias/arrochadeira.json?category=arrochadeira",
		"https://suamusica.com.br/_next/data/webid-1017/pt-BR/categorias/variados.json?category=variados",
	}
)

func main() {

	//newsrepo := database.NewNewsSQLiteRepository()
	//newsrepo := database.NewNewsMemoryRepository()
	postgres := adapter.NewPostgresAdapter()
	newsrepo := database.NewNewsPostgresRepository(postgres)
	imagedownloader := utils.NewImgDownloader()
	hitrepo := database.NewHitsPostgresRepository(postgres)
	news_service := service.NewNewsService(newsrepo, imagedownloader, hitrepo)

	crawler_service := service.NewCrawler()
	copier_service := service.NewCopier(news_service, crawler_service)
	crawler_controller := controllers.NewCrawlerController(copier_service)

	download_repository := database.NewDownloadPostgresRepository(postgres)
	copier_downloads := service.NewCopierDownload(download_repository, imagedownloader)

	news_controller := controllers.NewNewsController(news_service)

	toprepo := database.NewTopPostgresRepository(postgres)
	top_service := service.NewTopService(toprepo, newsrepo, hitrepo, news_service)
	top_controller := controllers.NewTopController(top_service)

	user_repo := database.NewUserPostgresRepository(postgres)
	user_service := service.NewUserService(user_repo)
	user_controller := controllers.NewUserController(user_service)

	banner_repo := database.NewBannerPostgresRepository(postgres)
	banner_service := service.NewBannerService(banner_repo, imagedownloader)
	banner_controller := controllers.NewBannerController(banner_service)

	var list_notifications []port.EmailPort
	email_notifications := notifications.NewGmailSMTP()
	ntfy_notifications := notifications.NewNtfyMobilePushNotifications()
	list_notifications = append(list_notifications, email_notifications, ntfy_notifications)
	notifications := notifications.NewNotifications(list_notifications)

	contact_repo := database.NewContactPostgresRepository(postgres)
	contact_service := service.NewContactService(contact_repo, imagedownloader, notifications)
	contact_controller := controllers.NewContactController(*contact_service)

	server := web.NewServerWeb()

	server.UserController(user_controller)
	server.TopController(top_controller)
	server.CrawlerController(crawler_controller)
	server.NewsController(news_controller)
	server.BannerController(banner_controller)
	server.ContactController(contact_controller)

	go copier_service.Start(list_pages, 3)
	go copier_downloads.Start(list_downloads, 30)

	//Função para gerar as top notícias a cada 60 minutos
	go top_service.Start(60)

	//Função para limpar as notícias inativas
	go news_service.StartCleanNews(60 * 24)

	server.Start()

}
