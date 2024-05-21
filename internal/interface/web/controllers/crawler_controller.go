package controllers

import (
	"net/http"
	"os"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	"github.com/gorilla/mux"
)

var list_pages = []string{
	"https://www.bahianoticias.com.br",
	"https://www.bahianoticias.com.br/holofote",
	"https://www.bahianoticias.com.br/esportes",
	"https://www.bahianoticias.com.br/bnhall",
	"https://www.bahianoticias.com.br/justica",
	"https://www.bahianoticias.com.br/saude",
	"https://www.bahianoticias.com.br/municipios",
}

type CrawlerController struct {
	Copier service.CopierService	
}

func NewCrawlerController(copier service.CopierService) *CrawlerController {
	return &CrawlerController{ Copier: copier }
}

func (c *CrawlerController) Crawler(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	key := vars["key"]

	if key != os.Getenv("KEY") {
		return
	}

	go c.Copier.Start(list_pages, 10)	

	msg := map[string]any{
		"ok": true,
		"message": "Notícias resgatadas!",

	}
	ResponseJson(w, msg, http.StatusOK)
	
}