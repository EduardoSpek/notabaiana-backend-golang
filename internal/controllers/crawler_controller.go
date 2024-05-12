package controllers

import (
	"net/http"
	"os"

	"github.com/eduardospek/notabaiana-backend-golang/internal/service"
	"github.com/gorilla/mux"
)

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

	go c.Copier.Start("https://www.bahianoticias.com.br/principal/rss.xml", 10)
	go c.Copier.Start("https://www.bahianoticias.com.br/holofote/rss.xml", 20)

	msg := map[string]any{
		"ok": true,
		"message": "Notícias resgatadas!",

	}
	ResponseJson(w, msg, http.StatusOK)
	
}