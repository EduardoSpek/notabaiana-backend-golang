package controllers

import (
	"net/http"
	"os"

	"github.com/eduardospek/bn-api/internal/service"
	"github.com/gorilla/mux"
)

type CrawlerController struct {
	Disparador service.DisparadorService	
}

func NewCrawlerController(disparador service.DisparadorService) *CrawlerController {
	return &CrawlerController{ Disparador: disparador }
}

func (c *CrawlerController) Crawler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key != os.Getenv("KEY") {
		return
	}

	c.Disparador.Start()

	msg := map[string]any{
		"ok": true,
		"message": "Notícias resgatadas!",

	}
	ResponseJson(w, msg, http.StatusOK)
	
}