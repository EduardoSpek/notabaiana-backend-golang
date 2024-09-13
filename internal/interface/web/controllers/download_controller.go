package controllers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/usecase"
	"github.com/eduardospek/notabaiana-backend-golang/internal/infra/database/postgres"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/gorilla/mux"
)

var (
	ErrDecodeImage = errors.New("não foi possível decodificar a imagem")
)

type DownloadController struct {
	DownloadRepository *postgres.DownloadPostgresRepository
	ImageDownloader    *utils.ImgDownloader
}

func NewDownloadController(repository *postgres.DownloadPostgresRepository, imagedownloader *utils.ImgDownloader) *DownloadController {
	return &DownloadController{DownloadRepository: repository, ImageDownloader: imagedownloader}
}

func (bc *DownloadController) CreateDownloadUsingTheForm(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any
	var downloadCreated *entity.Download

	//TokenVerifyByForm(w, r)

	downloadInput, image, err := bc.GetDownloadDataFromTheForm(r)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "problema com os dados do formulário",
			"erro":    "não foi possível resgatar os dados corretamente",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	createDownloadUsecase := usecase.NewCreateDownloadUsecase(bc.DownloadRepository)
	downloadCreated, err = createDownloadUsecase.Create(downloadInput)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "problema ao cadastrar os dados no banco",
			"erro":    "não foi possível cadatrar os dados corretamente",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	imgSaved, err := SaveImageForm(image, downloadCreated.Image, "downloads", 300, 300)

	if err != nil {
		fmt.Println(err)
	}

	if !imgSaved {
		downloadCreated.Image = ""
		updateDownloadUsecase := usecase.NewUpdateDownloadUsecase(bc.DownloadRepository)
		downloadUpdated, err := updateDownloadUsecase.Update(downloadCreated)

		if err != nil {
			fmt.Println(err)
		}

		downloadCreated = downloadUpdated
	}

	ResponseJson(w, downloadCreated, http.StatusOK)

}

func (bc *DownloadController) GetDownloadDataFromTheForm(r *http.Request) (*entity.Download, multipart.File, error) {

	vars := mux.Vars(r)
	id := vars["id"]

	category := r.FormValue("category")
	title := r.FormValue("title")
	text := r.FormValue("text")
	link := r.FormValue("link")
	visible, _ := strconv.ParseBool(r.FormValue("visible"))

	// Parse the multipart form data
	r.ParseMultipartForm(10 << 20) // 10 MB maximum

	// Get the images from the form
	image, _, _ := r.FormFile("image")

	download := &entity.Download{
		ID:       id,
		Category: category,
		Title:    title,
		Text:     text,
		Link:     link,
		Visible:  visible,
	}

	return download, image, nil

}

func (bc *DownloadController) FindAll(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pageStr := vars["page"]
	qtdStr := vars["qtd"]

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(qtdStr)
	if err != nil {
		limit = 24
	}

	downloadFindAll := usecase.NewFindAllDownloadUsecase(bc.DownloadRepository)
	downloads, err := downloadFindAll.FindAll(page, limit)

	if err != nil {
		msg := map[string]any{
			"ok":      false,
			"message": "não foi possível obter a lista de downloads",
			"erro":    "erro ao buscar lista",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, downloads, http.StatusOK)

}

func (bc *DownloadController) GetBySlug(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	slug := vars["slug"]

	downloadUsecase := usecase.NewGetBySlugDownloadUsecase(bc.DownloadRepository)
	download, err := downloadUsecase.GetBySlug(slug)

	if err != nil {
		msg := map[string]any{
			"ok":      false,
			"message": "não foi possível obter os dados do download",
			"erro":    "erro ao buscar registro",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, download, http.StatusOK)

}

func (bc *DownloadController) Search(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pageStr := vars["page"]
	str_search := r.URL.Query().Get("search")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	downloadUsecase := usecase.NewSearchDownloadUsecase(bc.DownloadRepository)
	downloads, err := downloadUsecase.Search(page, str_search)

	if err != nil {
		msg := map[string]any{
			"ok":      false,
			"message": "não foi possível obter a lista de downloads",
			"erro":    "erro ao buscar lista",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, downloads, http.StatusOK)

}

func (bc *DownloadController) FindCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	category := vars["category"]
	pageStr := vars["page"]

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	downloadUsecase := usecase.NewFindCategoryDownloadUsecase(bc.DownloadRepository)
	downloads, err := downloadUsecase.FindCategory(category, page)

	if err != nil {
		msg := map[string]any{
			"ok":      false,
			"message": "não foi possível obter a lista de downloads desta categoria",
			"erro":    "erro ao buscar lista por categoria",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, downloads, http.StatusOK)

}
