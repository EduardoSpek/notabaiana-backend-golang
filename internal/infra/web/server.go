package web

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eduardospek/bn-api/internal/infra/web/router"
	"github.com/eduardospek/bn-api/internal/middlewares"
	"github.com/gorilla/mux"
)

type ServerWeb struct {
	router *mux.Router
}

func NewServerWeb () *ServerWeb {
	return &ServerWeb{ router: router.NewGorillaMux() }
}

// Start run the application
func (serverweb *ServerWeb) Start() {
	api := serverweb.router

	api.Use(middlewares.CorsMiddleware)

	fmt.Println("O Servidor foi iniciado na porta "+ os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), api))

}