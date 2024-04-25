package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

var AllowedDomains = "www.bahianoticias.com.br"

func main() {
	//link := "https://www.bahianoticias.com.br/holofote/noticia/74009-daniel-cady-relembra-desconforto-de-ivete-sangalo-com-corpo-apos-gravidez-ela-se-queixava"

	link := "https://www.bahianoticias.com.br/noticia/291607-confusao-entre-equipes-de-vereadoras-termina-com-homem-baleado-em-pernambuco"
	
	text := getText(link)

	fmt.Println(text)

}

func getText(link string) string {
	var html, conteudo string
	//var err error

	collector := colly.NewCollector(
		colly.AllowedDomains(AllowedDomains),
	)

	//Obtém o texto da notícia
	collector.OnHTML(".sc-16306eb7-3.lbjQbj", func(e *colly.HTMLElement) {
		
		conteudo = e.DOM.Text()		

		html += conteudo
	
	})

	collector.Visit(link)

	return html
}

// func getEmded(link string) string {
// 	var html, conteudo string
// 	//var err error

// 	collector := colly.NewCollector(
// 		colly.AllowedDomains(AllowedDomains),
// 	)

// 	//Obtém o texto da notícia
// 	collector.OnHTML(".sc-16306eb7-3.lbjQbj", func(e *colly.HTMLElement) {
		
// 		conteudo = e.DOM.Text()		

// 		html += conteudo
	
// 	})

// 	// collector.OnHTML(".lazyload-placeholder", func(e *colly.HTMLElement) {
// 	// 	// Obter o valor do atributo "src" da imagem
// 	// 	conteudo = e.Attr("data-content")
// 	// 	conteudo_decoded, err := url.QueryUnescape(conteudo)

// 	// 	if err != nil {			
// 	// 		return
// 	// 	}	

// 	// 	html += conteudo_decoded
	
// 	// })

// 	// collector.OnHTML(".lazyload-scripts", func(e *colly.HTMLElement) {
// 	// 	// Obter o valor do atributo "src" da imagem
// 	// 	script = e.Attr("data-scripts")

// 	// 	script_decoded, err := url.QueryUnescape(script)

// 	// 	if err != nil {			
// 	// 		return
// 	// 	}
		
// 	// 	if !strings.Contains(html, script_decoded) {
// 	// 		html += script_decoded
// 	// 	}
		
// 	// })

// 	// Visitando a URL inicial
// 	collector.Visit(link)

// 	return html
// }