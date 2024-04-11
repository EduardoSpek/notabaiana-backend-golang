package test

import (
	"log"
	"os"
	"testing"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	"github.com/eduardospek/bn-api/internal/service"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestCrawler(t *testing.T) {
	t.Parallel()

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Erro ao carregar arquivo .env: %v", err)
	}

	list_news := []entity.News{
		{
			ID: uuid.NewString(),
			Title: "Titulo 1",
			Text: "Texto qualquer",
			Link: "https://www.bahianoticias.com.br/holofote/noticia/73832-alinne-rosa-passa-por-transformacao-apos-lipo-lad-confira-antes-e-depois",
			Image: "https://www.bahianoticias.com.br/fotos/holofote_noticias/73832/IMAGEM_NOTICIA_original.jpg",
		},
		{
			ID: uuid.NewString(),
			Title: "Titulo 2",
			Text: "Texto qualquer 2",
			Link: "https://www.bahianoticias.com.br/holofote/noticia2/73832-alinne-rosa-passa-por-transformacao-apos-lipo-lad-confira-antes-e-depois2",
			Image: "https://www.bahianoticias.com.br/fotos/holofote_noticias/73832/IMAGEM_NOTICIA_original2.jpg",
		},
		{
			ID: uuid.NewString(),
			Title: "Titulo 3",
			Text: "Texto qualquer 3",
			Link: "https://www.bahianoticias.com.br/holofote/noticia3/73832-alinne-rosa-passa-por-transformacao-apos-lipo-lad-confira-antes-e-depois3",
			Image: "https://www.bahianoticias.com.br/fotos/holofote_noticias/73832/IMAGEM_NOTICIA_original3.jpg",
		},
	}

	crawler := service.NewCrawler(list_news)

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