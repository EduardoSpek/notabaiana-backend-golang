package controllers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/gorilla/mux"
)

type ContactController struct {
	contact_service service.ContactService
}

func NewContactController(contactservice service.ContactService) *ContactController {
	return &ContactController{contact_service: contactservice}
}

func (bc *ContactController) AdminDeleteAll(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any
	var ids []string
	var contacts []entity.ContactDTO

	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		ResponseJson(w, err.Error(), http.StatusNotFound)
		return
	}

	for _, id := range ids {
		contacts = append(contacts, entity.ContactDTO{
			ID: id,
		})
	}

	err = bc.contact_service.AdminDeleteAll(contacts)

	if err != nil {
		ResponseJson(w, err.Error(), http.StatusNotFound)
		return
	}

	msg = map[string]any{
		"ok":      true,
		"message": "Todos os contacts selecionados foram removidos",
		"erro":    false,
	}

	ResponseJson(w, msg, http.StatusOK)

}

func (bc *ContactController) AdminFindAll(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any

	contacts, err := bc.contact_service.AdminFindAll()

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "nenhum contact encontrado",
			"erro":    err.Error(),
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, contacts, http.StatusOK)

}

func (bc *ContactController) AdminGetByID(w http.ResponseWriter, r *http.Request) {
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

	contact, err := bc.contact_service.AdminGetByID(id)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "Não existe registro com o ID informado",
			"erro":    err.Error(),
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, contact, http.StatusOK)

}

func (bc *ContactController) AdminDelete(w http.ResponseWriter, r *http.Request) {
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

	err = bc.contact_service.AdminDelete(id)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "O contact não pode ser excluído",
			"erro":    err.Error(),
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	msg = map[string]any{
		"ok":      true,
		"message": "Contact excluído",
		"erro":    false,
	}

	ResponseJson(w, msg, http.StatusOK)

}

func (bc *ContactController) CreateForm(w http.ResponseWriter, r *http.Request) {

	success := utils.GoogleRecaptchaVerify(r)

	if success {

		var msg map[string]any

		contactInput, images, err := bc.GetContactDataForm(r)

		if err != nil {
			msg = map[string]any{
				"ok":      false,
				"message": "problema com os dados do formulário",
				"erro":    "não foi possível resgatar os dados corretamente",
			}
			ResponseJson(w, msg, http.StatusNotFound)
			return
		}

		new, err := bc.contact_service.AdminCreateForm(images, contactInput)

		if err != nil {
			msg = map[string]any{
				"ok":      false,
				"message": "A menssagem não pode ser criada",
				"erro":    err.Error(),
			}
			ResponseJson(w, msg, http.StatusNotFound)
			return
		}

		ResponseJson(w, new, http.StatusOK)
	} else {
		msg := map[string]any{
			"ok":      false,
			"message": "Token do captcha inválido",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}
}

func (bc *ContactController) AdminCreateForm(w http.ResponseWriter, r *http.Request) {

	var msg map[string]any

	TokenVerifyByForm(w, r)

	contactInput, images, err := bc.GetContactDataForm(r)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "problema com os dados do formulário",
			"erro":    "não foi possível resgatar os dados corretamente",
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	new, err := bc.contact_service.AdminCreateForm(images, contactInput)

	if err != nil {
		msg = map[string]any{
			"ok":      false,
			"message": "A menssagem não pode ser criada",
			"erro":    err.Error(),
		}
		ResponseJson(w, msg, http.StatusNotFound)
		return
	}

	ResponseJson(w, new, http.StatusOK)

}

func (bc *ContactController) GetContactDataForm(r *http.Request) (entity.ContactDTO, multipart.File, error) {

	vars := mux.Vars(r)
	id := vars["id"]

	name := r.FormValue("name")
	email := r.FormValue("email")
	title := r.FormValue("title")
	text := r.FormValue("text")
	answered, _ := strconv.ParseBool(r.FormValue("answered"))

	// Parse the multipart form data
	r.ParseMultipartForm(10 << 20) // 10 MB maximum

	// Get the images from the form
	image, _, _ := r.FormFile("image")

	contact := &entity.ContactDTO{
		ID:       id,
		Name:     name,
		Email:    email,
		Title:    title,
		Text:     text,
		Answered: answered,
	}

	return *contact, image, nil

}
