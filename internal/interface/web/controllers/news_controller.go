package controllers

import (
	"errors"
	"fmt"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/gorilla/mux"
)

var (
	ErrParseForm = errors.New("erro ao obter a imagem")
)

type NewsController struct {
	news_service service.NewsService
}

func NewNewsController(newsservice service.NewsService) *NewsController {
	return &NewsController{news_service: newsservice}
}

func (c *NewsController) CleanNews(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any

	vars := mux.Vars(r)
	key := vars["key"]

	if key != os.Getenv("KEY") {
		return
	}

	c.news_service.CleanNews()

	msg = map[string]any{
		"ok":      true,
		"message": "notícias inativas removidas",
	}

	ResponseJson(w, msg, http.StatusOK)

}

func (c *NewsController) NewsMake(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any

	vars := mux.Vars(r)
	key := vars["key"]

	if key != os.Getenv("KEY") {
		return
	}

	news, err := c.news_service.NewsMake()

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "não há notícia para resgatar",
			"erro":    "sem novas notícias",
		}
		ResponseJson(w, msg, http.StatusForbidden)
		return
	}

	ResponseJson(w, news, http.StatusOK)

}

func (c *NewsController) NewsImage(w http.ResponseWriter, r *http.Request) {
	imageURL := r.URL.Query().Get("image")
	title := r.URL.Query().Get("title")
	var numberLines int

	err := os.MkdirAll("files", os.ModePerm)
	if err != nil {
		fmt.Println("Erro ao criar pasta:", err)
		return
	}

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("Erro ao obter o caminho do executável:", err)
	}

	diretorio := strings.Replace(cwd, "test", "", -1) + "/files/"

	totalWords := strings.Split(title, " ")
	if len(totalWords) > 1 {
		numberLines = len(totalWords) / 5
	}

	baseImgFile, err := os.Open(diretorio + "base_image.jpg")
	if err != nil {
		http.Error(w, "Could not open base image", http.StatusInternalServerError)
		return
	}
	defer baseImgFile.Close()

	baseImg, err := jpeg.Decode(baseImgFile)
	if err != nil {
		http.Error(w, "Could not decode base image", http.StatusInternalServerError)
		return
	}

	overlayImg, err := utils.DownloadImage(imageURL)
	if err != nil {
		http.Error(w, "Could not download overlay image", http.StatusInternalServerError)
		return
	}

	distaceY := 220 + (int(numberLines) * 60)
	resizedOverlay := utils.ResizeImage(overlayImg, 645, 405)
	finalImg := utils.OverlayImage(baseImg, resizedOverlay, 36, distaceY)

	fontFace, err := utils.LoadFont(diretorio+"roboto-latin-700-normal.ttf", 42)
	if err != nil {
		http.Error(w, "Could not load font", http.StatusInternalServerError)
		return
	}

	utils.AddLabel(finalImg, 26, 170, title, fontFace)

	w.Header().Set("Content-Disposition", "attachment; filename=final_image.jpg")
	w.Header().Set("Content-Type", "image/jpeg")
	jpeg.Encode(w, finalImg, nil)
}

func (s *NewsController) GetNewsDataFromTheForm(r *http.Request) (entity.News, multipart.File, error) {

	vars := mux.Vars(r)
	slug := vars["slug"]

	title := r.FormValue("title")
	text := r.FormValue("text")
	category := r.FormValue("category")
	id := r.FormValue("id")
	visible, _ := strconv.ParseBool(r.FormValue("visible"))

	// Parse the multipart form data
	r.ParseMultipartForm(10 << 20) // 10 MB maximum

	// Get the file from the form
	file, _, _ := r.FormFile("image")

	new := &entity.News{
		ID:       id,
		Title:    title,
		Text:     text,
		Visible:  visible,
		Category: category,
		Slug:     slug,
	}

	return *new, file, nil

}

func (c *NewsController) UpdateNewsUsingTheForm(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any
	token := r.FormValue("token")

	if token == "" {
		msg = map[string]any{
			"ok":      false,
			"message": "acesso não autorizado",
			"erro":    "token é necessário",
		}
		ResponseJson(w, msg, http.StatusForbidden)
		return
	}

	claims, err := utils.ValidateJWT(token)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "acesso não autorizado",
			"erro":    "token inválido",
		}
		ResponseJson(w, msg, http.StatusForbidden)
		return
	}

	if !claims.Admin {
		msg = map[string]any{
			"ok":      false,
			"message": "acesso não autorizado",
			"erro":    "sem permissão de admin",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	newsInput, file, err := c.GetNewsDataFromTheForm(r)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "problema com os dados do formulário",
			"erro":    "não foi possível resgatar os dados corretamente",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	new, err := c.news_service.UpdateNewsUsingTheForm(file, newsInput)

	if err != nil {
		msg := map[string]any{
			"ok":      false,
			"message": "A notícia não pode ser atualizada!",
			"erro":    err,
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, new, http.StatusOK)

}

func (c *NewsController) CreateNewsUsingTheForm(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any
	token := r.FormValue("token")

	if token == "" {
		msg = map[string]any{
			"ok":      false,
			"message": "acesso não autorizado",
			"erro":    "token é necessário",
		}
		ResponseJson(w, msg, http.StatusForbidden)
		return
	}

	claims, err := utils.ValidateJWT(token)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "acesso não autorizado",
			"erro":    "token inválido",
		}
		ResponseJson(w, msg, http.StatusForbidden)
		return
	}

	if !claims.Admin {
		msg = map[string]any{
			"ok":      false,
			"message": "acesso não autorizado",
			"erro":    "sem permissão de admin",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	newsInput, file, err := c.GetNewsDataFromTheForm(r)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "problema com os dados do formulário",
			"erro":    "não foi possível resgatar os dados corretamente",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	new, err := c.news_service.CreateNewsUsingTheForm(file, newsInput)

	if err != nil {
		msg := map[string]any{
			"ok":      false,
			"message": "A notícia não pode ser criada",
			"erro":    err,
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, new, http.StatusOK)

}

func (c *NewsController) GetNewsBySlug(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	slug := vars["slug"]

	new, err := c.news_service.GetNewsBySlug(slug)

	if err != nil {
		msg := map[string]any{
			"ok":      false,
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
			"ok":      false,
			"message": "Não foi possível limpar a tabela news",
		}
		ResponseJson(w, msg, http.StatusOK)
		return
	}

	msg = map[string]any{
		"ok":      true,
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

func (c *NewsController) AdminNews(w http.ResponseWriter, r *http.Request) {

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

	listnews := c.news_service.AdminFindAllNews(page, limit)

	ResponseJson(w, listnews, http.StatusOK)

}
