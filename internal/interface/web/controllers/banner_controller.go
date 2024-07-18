package controllers

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/gorilla/mux"
)

type BannerController struct {
	banner_service service.BannerService
}

func NewBannerController(bannerservice service.BannerService) *BannerController {
	return &BannerController{banner_service: bannerservice}
}

func (bc *BannerController) CreateBannerUsingTheForm(w http.ResponseWriter, r *http.Request) {

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

func (bc *BannerController) GetBannerDataFromTheForm(r *http.Request) (entity.BannerDTO, []multipart.File, error) {

	var images []multipart.File

	vars := mux.Vars(r)
	id := vars["id"]

	title := r.FormValue("title")
	link := r.FormValue("link")
	html := r.FormValue("html")
	tag := r.FormValue("tag")
	visible, _ := strconv.ParseBool(r.FormValue("visible"))

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
		ID:      id,
		Title:   title,
		Link:    link,
		Html:    html,
		Tag:     tag,
		Visible: visible,
	}

	return *banner, images, nil

}
