package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	"github.com/gorilla/mux"
)

type NewsController struct {
	news_service service.NewsService	
}

func NewNewsController(newsservice service.NewsService) *NewsController {
	return &NewsController{ news_service: newsservice }
}

func (c *NewsController) NewsCreateByForm(w http.ResponseWriter, r *http.Request) {

	//gcaptcha := r.FormValue("g-recaptcha-response")

	// Captura o token enviado pelo cliente
	token := r.FormValue("g-recaptcha-response")
        
	// Faça uma solicitação POST para a API de verificação do reCAPTCHA v3 do Google
	response, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", 
		map[string][]string{
			"secret":   {"6LdrROwpAAAAAPNbLdsY6XI6kI5R_xhV_2831cKJ"},
			"response": {token},
		})
	
	if err != nil {
		fmt.Println("Erro ao fazer a solicitação:", err)
		return
	}
	defer response.Body.Close()

	// Lê a resposta da API
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	// Decodifica a resposta JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Erro ao decodificar a resposta:", err)
		return
	}

	// Verifica se a resposta foi bem-sucedida
	success := result["success"].(bool)
	if success {
		fmt.Println("Token validado com sucesso!")
		// Faça o que você quiser aqui
	} else {
		fmt.Println("Falha na validação do token.")
	}
 	

	// new, err := c.news_service.NewsCreateByForm(r)

	// if err != nil {
	// 	msg := map[string]any{
	// 		"ok": false,
	// 		"message": "A notícia não pode ser criada",
	// 		"erro": err,
	// 	}
	// 	ResponseJson(w, msg, http.StatusNotFound)
	// 	return
	// }
	
	// ResponseJson(w, new, http.StatusOK)
	
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

func (c *NewsController) NewsCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	category := vars["category"]
	pageStr := vars["page"]

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1		
	}

	listnews := c.news_service.FindNewsCategory(category, page)
	
	ResponseJson(w, listnews, http.StatusOK)
	
}

func (c *NewsController) NewsTruncateTable(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["key"]

	if key != os.Getenv("KEY") {
		return
	}

	err := c.news_service.NewsTruncateTable()

	var msg map[string]any

	if err != nil {
		msg = map[string]any{
			"ok": false,
			"message": "Não foi possível limpar a tabela news",
		}
		ResponseJson(w, msg, http.StatusOK)
		return
	}

	msg = map[string]any{
		"ok": true,
		"message": "Tabela news Limpada com sucesso!",
	}

	ResponseJson(w, msg, http.StatusOK)

}

func (c *NewsController) SearchNews(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pageStr := vars["page"]
	str_search := r.URL.Query().Get("search")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1		
	}

	listnews := c.news_service.SearchNews(page, str_search)
	
	ResponseJson(w, listnews, http.StatusOK)
	
}