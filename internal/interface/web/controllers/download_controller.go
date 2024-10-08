package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	usecase "github.com/eduardospek/notabaiana-backend-golang/internal/domain/usecase/download"
	"github.com/gorilla/mux"
)

var (
	ErrDecodeImage = errors.New("não foi possível decodificar a imagem")
)

type DownloadController struct {
	DownloadRepository port.DownloadRepository
	ImageDownloader    port.ImageDownloader
	Cache              port.CachePort
}

func NewDownloadController(repository port.DownloadRepository, imagedownloader port.ImageDownloader, cache port.CachePort) *DownloadController {
	return &DownloadController{DownloadRepository: repository, ImageDownloader: imagedownloader, Cache: cache}
}

func (bc *DownloadController) CreateDownloadUsingTheForm(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any
	var downloadCreated *entity.Download

	err := TokenVerifyByForm(w, r)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "você não tem autorização",
			"erro":    "token inválido",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

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

	imgSaved, err := SaveImageForm(image, downloadCreated.ID+".jpg", "downloads", 300, 300)

	if err != nil {
		fmt.Println("IMG SAVED:", err)
	}

	if !imgSaved {
		downloadCreated.Image = ""
		updateDownloadUsecase := usecase.NewUpdateDownloadUsecase(bc.DownloadRepository)
		downloadUpdated, err := updateDownloadUsecase.Update(downloadCreated)

		if err != nil {
			fmt.Println(err)
		}

		downloadCreated = downloadUpdated
	} else {
		downloadCreated.Image = downloadCreated.ID + ".jpg"
		updateDownloadUsecase := usecase.NewUpdateDownloadUsecase(bc.DownloadRepository)
		updateDownloadUsecase.Update(downloadCreated)
	}

	bc.Cache.Cleanup()

	ResponseJson(w, downloadCreated, http.StatusOK)

}

func (bc *DownloadController) UpdateDownloadUsingTheForm(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any
	var downloadUpdated *entity.Download
	var imgSaved = false

	err := TokenVerifyByForm(w, r)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "você não tem autorização",
			"erro":    "token inválido",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

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

	updateDownloadUsecase := usecase.NewUpdateDownloadUsecase(bc.DownloadRepository)
	downloadUpdated, err = updateDownloadUsecase.Update(downloadInput)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "problema ao atualizar os dados no banco",
			"erro":    "não foi possível atualizar os dados corretamente",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	imgSaved, err = SaveImageForm(image, downloadUpdated.ID+".jpg", "downloads", 300, 300)

	if err != nil {
		fmt.Println("Erro ao Salvar imagem: ", err)
		downloadUpdated.Image = ""
		updateDownloadUsecase := usecase.NewUpdateDownloadUsecase(bc.DownloadRepository)
		downloadUp, err := updateDownloadUsecase.Update(downloadUpdated)

		if err != nil {
			fmt.Println(err)
		}

		downloadUpdated = downloadUp
	}

	if imgSaved {
		downloadUpdated.Image = downloadUpdated.ID + ".jpg"
		updateDownloadUsecase := usecase.NewUpdateDownloadUsecase(bc.DownloadRepository)
		updateDownloadUsecase.Update(downloadUpdated)
	}

	bc.Cache.Cleanup()

	ResponseJson(w, downloadUpdated, http.StatusOK)

}

func (bc *DownloadController) GetDownloadDataFromTheForm(r *http.Request) (*entity.Download, multipart.File, error) {

	vars := mux.Vars(r)
	slug := vars["slug"]

	id := r.FormValue("id")
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
		Slug:     slug,
		Visible:  visible,
	}

	return download, image, nil

}

func (bc *DownloadController) AdminFindAll(w http.ResponseWriter, r *http.Request) {

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

	downloadFindAll := usecase.NewAdminFindAllDownloadUsecase(bc.DownloadRepository)
	downloads, err := downloadFindAll.AdminFindAll(page, limit)

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

	cacheString := fmt.Sprintf("downloads:%d:%d", page, limit)

	if valor, existe := bc.Cache.Get(cacheString); existe {
		ResponseJson(w, valor, http.StatusOK)
		return
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

	bc.Cache.Set(cacheString, downloads)

	ResponseJson(w, downloads, http.StatusOK)

}

func (bc *DownloadController) FindAllTopViews(w http.ResponseWriter, r *http.Request) {

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

	cacheString := fmt.Sprintf("downloadsTopViews:%d:%d", page, limit)

	if valor, existe := bc.Cache.Get(cacheString); existe {
		ResponseJson(w, valor, http.StatusOK)
		return
	}

	downloadFindAllTopViews := usecase.NewFindAllTopViewsDownloadUsecase(bc.DownloadRepository)
	downloads, err := downloadFindAllTopViews.FindAllTopViews(page, limit)

	if err != nil {
		msg := map[string]any{
			"ok":      false,
			"message": "não foi possível obter a lista de downloads",
			"erro":    "erro ao buscar lista",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	bc.Cache.Set(cacheString, downloads)

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

func (bc *DownloadController) AdminGetBySlug(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	slug := vars["slug"]

	downloadUsecase := usecase.NewAdminGetBySlugDownloadUsecase(bc.DownloadRepository)
	download, err := downloadUsecase.AdminGetBySlug(slug)

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

	cacheString := fmt.Sprintf("downloadFindCategory:%s:%d", category, page)

	if valor, existe := bc.Cache.Get(cacheString); existe {
		ResponseJson(w, valor, http.StatusOK)
		return
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

	bc.Cache.Set(cacheString, downloads)

	ResponseJson(w, downloads, http.StatusOK)

}

func (bc *DownloadController) Delete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	bc.Cache.Cleanup()

	downloadUsecase := usecase.NewDeleteDownloadUsecase(bc.DownloadRepository)
	err := downloadUsecase.Delete(id)

	if err != nil {
		msg := map[string]any{
			"ok":      false,
			"message": "não foi possível deletar o registro",
			"erro":    "erro ao deletar registro",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	msg := map[string]any{
		"ok":      true,
		"message": "registro deletado com sucesso",
		"erro":    nil,
	}

	ResponseJson(w, msg, http.StatusOK)

}

func (bc *DownloadController) DeleteAll(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any
	var ids []string
	var downloads []*entity.Download

	bc.Cache.Cleanup()

	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		ResponseJson(w, err.Error(), http.StatusNotFound)
		return
	}

	for _, id := range ids {
		downloads = append(downloads, &entity.Download{
			ID: id,
		})
	}

	downloadUsecase := usecase.NewDeleteAllDownloadUsecase(bc.DownloadRepository)
	err = downloadUsecase.DeleteAll(downloads)

	if err != nil {
		ResponseJson(w, err.Error(), http.StatusNotFound)
		return
	}

	msg = map[string]any{
		"ok":      true,
		"message": "Os Downloads selecionados foram removidos",
		"erro":    false,
	}

	ResponseJson(w, msg, http.StatusOK)

}
