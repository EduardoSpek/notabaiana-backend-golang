package controllers

import (
	"net/http"
	"strconv"

	"github.com/eduardospek/bn-api/internal/service"
	"github.com/gorilla/mux"
)

type NewsController struct {
	news_service service.NewsService	
}

func NewNewsController(newsservice service.NewsService) *NewsController {
	return &NewsController{ news_service: newsservice }
}

func (c *NewsController) GetNewsBySlug(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	slug := vars["slug"]

	new, err := c.news_service.GetNewsBySlug(slug)

	if err != nil {
		msg := map[string]any{
			"ok": false,
			"message": "não há notícia com este slug",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}
	
	ResponseJson(w, new, http.StatusOK)
	
}

func (c *NewsController) News(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pageStr := vars["page"]
	qtdStr := vars["qtd"]

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1		
	}

	limit, err := strconv.Atoi(qtdStr)
	if err != nil {
		limit = 10		
	}

	listnews := c.news_service.FindAllNews(page, limit)
	
	ResponseJson(w, listnews, http.StatusOK)
	
}