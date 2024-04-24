package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	//link := "https://www.bahianoticias.com.br/holofote/noticia/74009-daniel-cady-relembra-desconforto-de-ivete-sangalo-com-corpo-apos-gravidez-ela-se-queixava"

	link := "https://www.bahianoticias.com.br/holofote/noticia/74005-ivete-sangalo-anuncia-leilao-de-figurino-para-ajudar-instituicao-com-atuacao-no-sertao-nordestino"
	html := getEmded(link)

	fmt.Println(html)

}

func getEmded(link string) string {
	var html, conteudo, script string
	//var err error

	collector := colly.NewCollector(
		colly.AllowedDomains("www.bahianoticias.com.br"),
	)

	collector.OnHTML(".lazyload-placeholder", func(e *colly.HTMLElement) {
		// Obter o valor do atributo "src" da imagem
		conteudo = e.Attr("data-content")
		conteudo_decoded, err := url.QueryUnescape(conteudo)

		if err != nil {			
			return
		}	

		html += conteudo_decoded
	
	})

	collector.OnHTML(".lazyload-scripts", func(e *colly.HTMLElement) {
		// Obter o valor do atributo "src" da imagem
		script = e.Attr("data-scripts")

		script_decoded, err := url.QueryUnescape(script)

		if err != nil {			
			return
		}
		
		if !strings.Contains(html, script_decoded) {
			html += script_decoded
		}
		
	})

	// Visitando a URL inicial
	collector.Visit(link)

	return html
}