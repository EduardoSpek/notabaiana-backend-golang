package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
)

type UserController struct {
	user_service service.UserService
}

func NewUserController(userservice service.UserService) *UserController {
	return &UserController{user_service: userservice}
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
