package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/gorilla/mux"
)

type UserController struct {
	user_service service.UserService
}

func NewUserController(userservice service.UserService) *UserController {
	return &UserController{user_service: userservice}
}

func (u *UserController) CheckUser(w http.ResponseWriter, r *http.Request) {    
    var msg map[string]any
        
    tokenStr := r.Header.Get("Authorization")
    if tokenStr == "" {
        msg = map[string]any{
            "ok": false,
            "message": "acesso não autorizado",
            "erro": "token é necessário",
        }
        ResponseJson(w, msg, http.StatusForbidden)
        return
    }

    tokenStr = tokenStr[len("Bearer "):]

    claims, err := utils.ValidateJWT(tokenStr)

    if err != nil {
        msg = map[string]any{
            "ok": false,
            "message": "acesso não autorizado",
            "erro": "token inválido",
        }
        ResponseJson(w, msg, http.StatusForbidden)
        return
    }

    if !claims.Admin {
        msg = map[string]any{
            "ok": false,
            "message": "acesso não autorizado",
            "erro": "sem permissão de admin",
        }
        ResponseJson(w, msg, http.StatusNotFound)
        return
    }  
    
    msg = map[string]any{
        "ok": true,
        "message": "acesso autorizado",
        "erro": nil,
    }

	ResponseJson(w, msg, http.StatusOK)
    
}

func (u *UserController) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
    var userInput entity.UserInput	
    
    err := json.NewDecoder(r.Body).Decode(&userInput)
    if err!= nil {
        ResponseJson(w, err.Error(), http.StatusNotFound)
        return
    }

    userToken, err := u.user_service.Login(userInput)

	if err != nil {
        msg := map[string]any{
            "ok": false,
            "message": "Não foi possível efetuar o login",
            "erro": err.Error(),
        }
		ResponseJson(w, msg, http.StatusNotFound)
        return
	}

	ResponseJson(w, userToken, http.StatusOK)
}

func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    var userInput entity.UserInput

	vars := mux.Vars(r)
	id := vars["id"]
    
    err := json.NewDecoder(r.Body).Decode(&userInput)
    if err!= nil {
        ResponseJson(w, err.Error(), http.StatusNotFound)
        return
    }
    
    tokenStr := r.Header.Get("Authorization")
    if tokenStr == "" {
        ResponseJson(w, "acesso não autorizado", http.StatusForbidden)
        return
    }

    tokenStr = tokenStr[len("Bearer "):]

    claims, err := utils.ValidateJWT(tokenStr)

    if err != nil {
        ResponseJson(w, "acesso não autorizado: token inválido", http.StatusForbidden)
        return
    }

    if strings.TrimSpace(id) != strings.TrimSpace(claims.ID) && !claims.Admin {
        ResponseJson(w, "acesso não autorizado: usuário não identificado", http.StatusNotFound)
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
