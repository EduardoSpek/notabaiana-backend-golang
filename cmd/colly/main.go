package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	"github.com/gocolly/colly"
)

var AllowedDomains = "www.bahianoticias.com.br"

func main() {
	//link := "https://www.bahianoticias.com.br/holofote/noticia/74009-daniel-cady-relembra-desconforto-de-ivete-sangalo-com-corpo-apos-gravidez-ela-se-queixava"

	link := "https://www.bahianoticias.com.br/holofote"
	
	text := GetNewsFromPage(link)

	fmt.Println(text)

	// link := "https://www.bahianoticias.com.br/folha/noticia/275994-globo-fecha-contrato-e-vai-transmitir-festival-de-parintins-2024-para-todo-o-brasil"

	// image, err := GetImageLink(link)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(image)

}

func GetNewsFromPage(link string) interface{} {
	var conteudo string
	//var err error

	collector := colly.NewCollector(
		colly.AllowedDomains(AllowedDomains),
	)

	var titulos []string
	var texts []string	
	var links []string
	var images []string

	//Obtém titulos
	collector.OnHTML("h3.sc-b4c8ccf3-1.ireAxk", func(e *colly.HTMLElement) {
		
		conteudo = e.DOM.Text()		

		titulos = append(titulos, conteudo)
	
	})

	//Obtém links
	collector.OnHTML(".sc-b4c8ccf3-0.fsXNOt a", func(e *colly.HTMLElement) {
		
		conteudo = e.Attr("href")

		conteudo = "https://www.bahianoticias.com.br" + conteudo

		image, err := GetImageLink(conteudo)

		if err != nil {
			return
		}
		
		images = append(images, image)

		links = append(links, conteudo)
	
	})

	//Obtém textos
	collector.OnHTML(".sc-81cf810-3.gCNTHg", func(e *colly.HTMLElement) {
		
		conteudo = e.DOM.Text()		

		texts = append(texts, conteudo)
	
	})

	//Obtém images 
	collector.OnHTML(".sc-81cf810-2.hiSMeg div span img", func(e *colly.HTMLElement) {
		
		conteudo = e.Attr("src")

		conteudo_decoded, err := url.QueryUnescape(conteudo)		

		if err != nil {			
			return
		}	

		images = append(images, conteudo_decoded)
	
	})
	
	// Visitando a URL inicial
	collector.Visit(link)

	var lista []entity.News

	for i, item := range titulos {
		new := entity.NewNews(entity.News{
			Title: item,
			Text: texts[i],
			Image: images[i],
			Link: links[i],
			Visible: true,
		})
		lista = append(lista, *new)
	}

	return lista
}

func GetIdFromLink(link string) (int, error) {
	
	partes := strings.Split(link, "/")

	partes_total := len(partes)

	var maispartes []string

	if partes_total == 6 {
		maispartes = strings.Split(partes[5], "-")
	} else if partes_total == 5 {
		maispartes = strings.Split(partes[4], "-")
	}
		
	id, err := strconv.Atoi(maispartes[0])

	if err != nil {
		return 0, err
	}

	return id, nil
}
func ReturnPathFromLink(link string) (string, error) {
	if strings.Contains(link, "folha/noticia") {
		return "folha_noticias", nil
	} else if strings.Contains(link, "holofote/noticia") {
		return "holofote_noticias", nil
	} else if strings.Contains(link, "municipios/noticia") {
		return "municipios_noticias", nil
	} else if strings.Contains(link, "saude/noticia") {
		return "saude_noticias", nil
	} else if strings.Contains(link, "justica/noticia") {
		return "justica_noticias", nil
	} else if strings.Contains(link, "bnhall/noticia") {
		return "hall_noticias", nil
	} else if strings.Contains(link, "esportes/vitoria") {
		return "esportes_vitorias", nil
	} else if strings.Contains(link, "esportes/bahia") {
		return "esportes_bahias", nil
	} else if strings.Contains(link, "esportes/noticia") {
		return "esportes_noticias", nil
	} else {
		return "principal_noticias", nil
	}	
}

func GetImageLink(link string) (string, error) {
	id, err := GetIdFromLink(link)

	if err != nil {
		return "", err
	}

	path, err := ReturnPathFromLink(link)

	if err != nil {
		return "", err
	}

	newlink := fmt.Sprintf("https://www.bahianoticias.com.br/fotos/%s/%d/IMAGEM_NOTICIA_5.jpg", path, id)

	return newlink, nil
}

// func getText(link string) string {
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

// 	collector.Visit(link)

// 	return html
// }

// func getEmded(link string) string {
// 	var html, conteudo, script string
// 	//var err error

// 	collector := colly.NewCollector(
// 		colly.AllowedDomains(AllowedDomains),
// 	)

// 	//Obtém o texto da notícia
// 	// collector.OnHTML(".sc-16306eb7-3.lbjQbj", func(e *colly.HTMLElement) {
		
// 	// 	conteudo = e.DOM.Text()		

// 	// 	html += conteudo
	
// 	// })

// 	collector.OnHTML(".lazyload-placeholder", func(e *colly.HTMLElement) {
// 		// Obter o valor do atributo "src" da imagem
// 		conteudo = e.Attr("data-content")
// 		conteudo_decoded, err := url.QueryUnescape(conteudo)

// 		if err != nil {			
// 			return
// 		}	

// 		html += conteudo_decoded
	
// 	})

// 	collector.OnHTML(".lazyload-scripts", func(e *colly.HTMLElement) {
// 		// Obter o valor do atributo "src" da imagem
// 		script = e.Attr("data-scripts")

// 		script_decoded, err := url.QueryUnescape(script)

// 		if err != nil {			
// 			return
// 		}
		
// 		if !strings.Contains(html, script_decoded) {
// 			html += script_decoded
// 		}
		
// 	})

// 	// Visitando a URL inicial
// 	collector.Visit(link)

// 	return html
// }