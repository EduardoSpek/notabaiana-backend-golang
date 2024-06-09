package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	"github.com/gorilla/mux"
)

type UserController struct {
	user_service service.UserService
}

func NewUserController(userservice service.UserService) *UserController {
	return &UserController{user_service: userservice}
}

func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    var userInput entity.UserInput

	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Println("ID Controller", id)
    
    err := json.NewDecoder(r.Body).Decode(&userInput)
    if err!= nil {
        ResponseJson(w, err.Error(), http.StatusNotFound)
        return
    }

    userUpdated, err := u.user_service.UpdateUser(id, userInput)

	if err != nil {
		ResponseJson(w, err.Error(), http.StatusNotFound)
        return
	}

	ResponseJson(w, userUpdated, http.StatusOK)
    
}

func (u *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    var userInput entity.UserInput
    
    err := json.NewDecoder(r.Body).Decode(&userInput)
    if err!= nil {
        ResponseJson(w, err.Error(), http.StatusNotFound)
        return
    }

    userCreated, err := u.user_service.CreateUser(userInput)

	if err != nil {
		ResponseJson(w, err.Error(), http.StatusNotFound)
        return
	}

	ResponseJson(w, userCreated, http.StatusOK)
    
}
