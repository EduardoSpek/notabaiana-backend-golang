package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/middlewares"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/router"
	"github.com/gorilla/mux"
)

type ServerWeb struct {
	router *mux.Router
}

func NewServerWeb() *ServerWeb {
	return &ServerWeb{router: router.NewGorillaMux()}
}

// Start run the application
func (serverweb *ServerWeb) Start() {
	api := serverweb.router
	// Rota para servir arquivos estáticos
	api.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	api.Use(middlewares.CorsMiddleware)

	srv := &http.Server{
		Handler:      api,
		Addr:         ":" + os.Getenv("PORT"),
		ReadTimeout:  5 * time.Second,  // tempo máximo para ler o request
		WriteTimeout: 10 * time.Second, // tempo máximo para escrever a resposta
		IdleTimeout:  60 * time.Second, // conexões keep-alive
	}

	fmt.Println("O Servidor foi iniciado na porta " + os.Getenv("PORT"))
	log.Fatal(srv.ListenAndServe())

}
