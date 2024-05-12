package test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eduardospek/notabaiana-backend-golang/internal/controllers"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	database "github.com/eduardospek/notabaiana-backend-golang/internal/infra/database/memorydb"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type TestCase struct {
	Esperado any
	Recebido any
	Descricao string
}

func Resultado(t *testing.T, esperado any, recebido any, descricao string) {
    t.Helper()
    if esperado != recebido {
        t.Errorf("Descricao: %s | Esperado: %s | Recebido: %s", descricao, esperado, recebido)
    }
}

func TestNewsEntity(t *testing.T) {
	t.Parallel()
	
	news := entity.News{		
		Title: "Titulo",
		Text: "Texto",
		Link: "http://www.eduardospek.com.br",
		Image: "https://www.bahianoticias.com.br/fotos/holofote_noticias/73825/IMAGEM_NOTICIA_original.jpg",
		Visible: true,
	}

	n := entity.NewNews(news)

	testcases := []TestCase{
		{
			Esperado: news.Title,
			Recebido: n.Title,
			Descricao: "Title",
		},
		{
			Esperado: news.Text,
			Recebido: n.Text,
			Descricao: "Text",
		},
		{
			Esperado: news.Link,
			Recebido: n.Link,
			Descricao: "Link",
		},
		{
			Esperado: news.Image,
			Recebido: n.Image,
			Descricao: "Image",
		},
		{
			Esperado: true,
			Recebido: n.Visible,
			Descricao: "Visible",
		},
	}

	for _, teste := range testcases {
		Resultado(t, teste.Esperado, teste.Recebido, teste.Descricao)
	}

}

func TestNewsService(t *testing.T) {
	t.Parallel()

	news_repo := database.NewNewsMemoryRepository()
	imagedownloader := utils.NewImgDownloader()
	news_service := service.NewNewsService(news_repo, imagedownloader)

	t.Run("Deve criar uma nova noticia no banco", func (t *testing.T)  {
		news := entity.News{		
			Title: "Titulo",
			Text: "Texto",
			Link: "http://www.eduardospek.com.br",
			Image: "https://www.bahianoticias.com.br/fotos/holofote_noticias/73825/IMAGEM_NOTICIA_original.jpg",
			Visible: true,
			Category: "holofote",
		}
	
		_, err := news_service.CreateNews(news)
	
		if err != nil {
			t.Error(err)
		}

		news = entity.News{		
			Title: "Eduardo Spek",
			Text: "Texto",
			Link: "http://www.eduardospek.com.br",
			Image: "https://www.bahianoticias.com.br/fotos/holofote_noticias/73825/IMAGEM_NOTICIA_original.jpg",
			Visible: true,
			Category: "principal",
		}
	
		_, err = news_service.CreateNews(news)
	
		if err != nil {
			t.Error(err)
		}

		news = entity.News{		
			Title: "Eduardo Spek na tela da globo",
			Text: "Texto",
			Link: "http://www.eduardospek.com.br",
			Image: "https://www.bahianoticias.com.br/fotos/holofote_noticias/73825/IMAGEM_NOTICIA_original.jpg",
			Visible: true,
			Category: "holofote",
		}
	
		_, err = news_service.CreateNews(news)
	
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Deve listar as noticias do banco", func (t *testing.T)  {

		lista := news_service.FindAllNews(1, 4)
		
		newsList := lista.(struct{
			List_news []entity.News `json:"news"`
			Pagination map[string][]int `json:"pagination"`
		})		

		if string(newsList.List_news[0].Title) != "Titulo" {
			t.Error("Erro: Não foi possível retornar as notícias")			
		}
	
	})

	t.Run("Deve buscar noticias do banco", func (t *testing.T)  {

		str_search := "Eduardo"

		lista := news_service.SearchNews(1, str_search)
		
		newsList := lista.(struct{
			List_news []entity.News `json:"news"`
			Pagination map[string][]int `json:"pagination"`
			Search string `json:"search"`
		})

		var passou bool = false

		for _, news := range newsList.List_news {
			if strings.Contains(news.Title, str_search) {
				passou = true				
			}
			
		}

		if len(newsList.List_news) < 2 {
			t.Error("Erro: Não encontrou os dois registros")			
		}

		if !passou {
			t.Error("Erro: Não foi possível retornar as notícias da busca")			
		}
	
	})
}

func TestNewsController(t *testing.T) {
	t.Parallel()

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Erro ao carregar arquivo .env: %v", err)
	}
	

	t.Run("Deve buscar noticias no banco e retornar status 200", func(t *testing.T) {

		str_search := "Eduardo"

		req, err := http.NewRequest("GET", "/news/busca/1?search=" + str_search, nil)
		
		if err != nil {
			t.Fatal(err)
		}

		repo := database.NewNewsMemoryRepository()
		imagedownloader := utils.NewImgDownloader()
		news_service := service.NewNewsService(repo, imagedownloader)		
		controller := controllers.NewNewsController(*news_service)

		news := entity.News{		
			Title: "Eduardo Spek",
			Text: "Texto",
			Link: "http://www.eduardospek.com.br",
			Image: "https://www.bahianoticias.com.br/fotos/holofote_noticias/73825/IMAGEM_NOTICIA_original.jpg",
			Visible: true,
		}
	
		_, err = news_service.CreateNews(news)
	
		if err != nil {
			t.Error(err)
		}
		
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/news/busca/{page}", controller.SearchNews)

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		fmt.Println(rr.Body.String())

		expected := str_search

		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}		

	})
}