package test

import (
	"log"
	"os"
	"testing"

	"github.com/eduardospek/bn-api/internal/service"
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