package controllers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	"github.com/gorilla/mux"
)

type BannerController struct {
	banner_service *service.BannerService
}

func NewBannerController(bannerservice *service.BannerService) *BannerController {
	return &BannerController{banner_service: bannerservice}
}

func (bc *BannerController) AdminDeleteAllBanner(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any
	var ids []string
	var banners []entity.BannerDTO

	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		ResponseJson(w, err.Error(), http.StatusNotFound)
		return
	}

	for _, id := range ids {
		banners = append(banners, entity.BannerDTO{
			ID: id,
		})
	}

	err = bc.banner_service.AdminDeleteAll(banners)

	if err != nil {
		ResponseJson(w, err.Error(), http.StatusNotFound)
		return
	}

	msg = map[string]any{
		"ok":      true,
		"message": "Todos os banners selecionados foram removidos",
		"erro":    false,
	}

	ResponseJson(w, msg, http.StatusOK)

}

func (bc *BannerController) AdminBannerList(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any

	banners, err := bc.banner_service.AdminFindAll()

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "nenhum banner encontrado",
			"erro":    err.Error(),
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, banners, http.StatusOK)

}

func (bc *BannerController) BannerList(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any

	banners, err := bc.banner_service.FindAll()

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "nenhum banner encontrado",
			"erro":    err.Error(),
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, banners, http.StatusOK)

}

func (bc *BannerController) FindBanner(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var msg map[string]any

	vars := mux.Vars(r)
	id := vars["id"]

	banner, err := bc.banner_service.FindBanner(id)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "Não existe registro com o ID informado",
			"erro":    err.Error(),
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, banner, http.StatusOK)

}

func (bc *BannerController) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var msg map[string]any

	vars := mux.Vars(r)
	id := vars["id"]

	err := TokenVerifyByHeader(w, r)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": err.Error(),
			"erro":    "não autorizado",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	err = bc.banner_service.Delete(id)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "O banner não pode ser excluído",
			"erro":    err.Error(),
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	msg = map[string]any{
		"ok":      true,
		"message": "Banner excluído",
		"erro":    false,
	}

	ResponseJson(w, msg, http.StatusOK)

}

func (bc *BannerController) UpdateBannerUsingTheForm(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any

	err := TokenVerifyByForm(w, r)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": err,
			"erro":    "não autorizado",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	bannerInput, images, err := bc.GetBannerDataFromTheForm(r)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "problema com os dados do formulário",
			"erro":    "não foi possível resgatar os dados corretamente",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	new, err := bc.banner_service.UpdateBannerUsingTheForm(images, bannerInput)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "A notícia não pode ser criada",
			"erro":    err,
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, new, http.StatusOK)

}

func (bc *BannerController) CreateBannerUsingTheForm(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any

	TokenVerifyByForm(w, r)

	bannerInput, images, err := bc.GetBannerDataFromTheForm(r)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "problema com os dados do formulário",
			"erro":    "não foi possível resgatar os dados corretamente",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	new, err := bc.banner_service.CreateBannerUsingTheForm(images, bannerInput)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "A notícia não pode ser criada",
			"erro":    err.Error(),
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, new, http.StatusOK)

}

func (bc *BannerController) GetBannerDataFromTheForm(r *http.Request) (entity.BannerDTO, []multipart.File, error) {

	var images []multipart.File

	vars := mux.Vars(r)
	id := vars["id"]

	title := r.FormValue("title")
	link := r.FormValue("link")
	html := r.FormValue("html")
	tag := r.FormValue("tag")
	visible, _ := strconv.ParseBool(r.FormValue("visible"))
	visible_image1, _ := strconv.ParseBool(r.FormValue("visible_image1"))
	visible_image2, _ := strconv.ParseBool(r.FormValue("visible_image2"))
	visible_image3, _ := strconv.ParseBool(r.FormValue("visible_image3"))

	// Parse the multipart form data
	r.ParseMultipartForm(10 << 20) // 10 MB maximum

	// Get the images from the form
	image1, _, _ := r.FormFile("image1")
	image2, _, _ := r.FormFile("image2")
	image3, _, _ := r.FormFile("image3")

	images = append(images, image1)
	images = append(images, image2)
	images = append(images, image3)

	banner := &entity.BannerDTO{
		ID:            id,
		Title:         title,
		Link:          link,
		Html:          html,
		Tag:           tag,
		Visible:       visible,
		VisibleImage1: visible_image1,
		VisibleImage2: visible_image2,
		VisibleImage3: visible_image3,
	}

	return *banner, images, nil

}
