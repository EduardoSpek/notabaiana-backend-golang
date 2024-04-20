package test

import (
	"testing"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	database "github.com/eduardospek/bn-api/internal/infra/database/memorydb"
	"github.com/eduardospek/bn-api/internal/service"
	"github.com/eduardospek/bn-api/internal/utils"
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
		}
	
		_, err := news_service.CreateNews(news)
	
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
	

}