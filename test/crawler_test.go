package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/eduardospek/notabaiana-backend-golang/internal/controllers"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	database "github.com/eduardospek/notabaiana-backend-golang/internal/infra/database/memorydb"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func TestCrawler(t *testing.T) {
	t.Parallel()

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Erro ao carregar arquivo .env: %v", err)
	}

	crawler := service.NewCrawler()

	t.Run("Deve obter os dados do RSS", func(t *testing.T) {
		
		rss := crawler.GetRSS(os.Getenv("URL_RSS"))	
		
		title := rss.Channel.Items[0].Title

		if title == "" {
			t.Error("Erro: Não foi possível obter o RSS")
		}
		
	})
}

func TestCrawlerController(t *testing.T) {
	t.Parallel()

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Erro ao carregar arquivo .env: %v", err)
	}

	t.Run("Deve cadastrar as noticias no banco e retornar status 200", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/crawler/" + os.Getenv("KEY"), nil)
		
		if err != nil {
			t.Fatal(err)
		}

		repo := database.NewNewsMemoryRepository()
		imagedownloader := utils.NewImgDownloader()
		news := service.NewNewsService(repo, imagedownloader)
		crawler := service.NewCrawler()
		copier := service.NewCopier(*news, *crawler)
		controller := controllers.NewCrawlerController(*copier)
		
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/crawler/{key}", controller.Crawler)

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := `{"message":"Notícias resgatadas!","ok":true}`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}		

	})
}