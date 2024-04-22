package controllers

import (
	"net/http"
	"os"

	"github.com/eduardospek/bn-api/internal/service"
	"github.com/gorilla/mux"
)

type TopController struct {
	TopService service.TopService	
}

func NewTopController(topservice service.TopService) *TopController {
	return &TopController{ TopService: topservice }
}

func (t *TopController) CreateTop(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["key"]

	if key != os.Getenv("KEY") {
		return
	}

	go t.TopService.TopCreate()	

	msg := map[string]any{
		"ok": true,
		"message": "Criado o Top Notícias",

	}
	ResponseJson(w, msg, http.StatusOK)
	
	
}