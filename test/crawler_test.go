package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/eduardospek/bn-api/internal/controllers"
	database "github.com/eduardospek/bn-api/internal/infra/database/memorydb"
	"github.com/eduardospek/bn-api/internal/service"
	"github.com/eduardospek/bn-api/internal/utils"
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

		testcases := []TestCase{
			{
				Esperado: "BAHIA NOTICIAS",
				Recebido: string(rss.Channel.Title),
			},
			{
				Esperado: `Anderson Leonardo recebe alta da UTI e é encaminhado para quarto: "Respirando sem ajuda de aparelhos"`,
				Recebido: string(rss.Channel.Items[0].Title),
			},
			{
				Esperado: "https://www.bahianoticias.com.br/holofote/noticia/73833-anderson-leonardo-recebe-alta-da-uti-e-e-encaminhado-para-quarto-respirando-sem-ajuda-de-aparelhos",
				Recebido: string(rss.Channel.Items[0].Link),
			},
			{
				Esperado: "https://www.bahianoticias.com.br/fotos/holofote_noticias/73833/IMAGEM_NOTICIA_original.jpg",
				Recebido: string(rss.Channel.Items[0].Media.URL),
			},
		}

		for _, teste := range testcases {
			Resultado(t, teste.Esperado, teste.Recebido)
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
		disparador := service.NewDisparador(*news, *crawler)
		controller := controllers.NewCrawlerController(*disparador)
		
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