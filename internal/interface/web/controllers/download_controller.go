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

	TokenVerifyByForm(w, r)

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
	downloadCreated, err := createDownloadUsecase.Create(downloadInput)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "problema ao cadastrar os dados no banco",
			"erro":    "não foi possível cadatrar os dados corretamente",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	err = SaveImageForm(image, downloadCreated.Image, "downloads")

	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "A notícia não pode ser criada",
			"erro":    err.Error(),
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
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
