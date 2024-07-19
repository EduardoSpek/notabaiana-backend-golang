package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

var (
	ErrToken = errors.New("acesso não autorizado")
)

func ResponseJson(w http.ResponseWriter, data any, statusCode int) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.New("responseJson: não foi possível converter para json")
	}

	// Escrevendo a resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)

	return nil

}

func TokenVerifyByForm(w http.ResponseWriter, r *http.Request) error {

	token := r.FormValue("token")

	if token == "" {
		return ErrToken
	}

	claims, err := utils.ValidateJWT(token)

	if err != nil {
		return ErrToken
	}

	if !claims.Admin {
		return ErrToken
	}

	return nil
}
