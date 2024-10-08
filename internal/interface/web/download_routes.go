package web

import (
	"net/http"

	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/middlewares"
)

func (s *ServerWeb) DownloadController(downloadcontroller *controllers.DownloadController) {
	s.router.HandleFunc("/downloads/category/{category}/{page}", downloadcontroller.FindCategory).Methods("GET")
	s.router.HandleFunc("/downloads/topviews/{page}/{qtd}", downloadcontroller.FindAllTopViews).Methods("GET")
	s.router.HandleFunc("/downloads/search/{page}", downloadcontroller.Search).Methods("GET")
	s.router.HandleFunc("/download/{slug}", downloadcontroller.GetBySlug).Methods("GET")
	s.router.HandleFunc("/downloads/{page}/{qtd}", downloadcontroller.FindAll).Methods("GET")
	s.router.Handle("/admin/downloads/{slug}", middlewares.JwtMiddleware(http.HandlerFunc(downloadcontroller.AdminGetBySlug))).Methods("GET")
	s.router.HandleFunc("/admin/downloads/create", downloadcontroller.CreateDownloadUsingTheForm).Methods("POST")
	s.router.HandleFunc("/admin/downloads/update/{slug}", downloadcontroller.UpdateDownloadUsingTheForm).Methods("POST")
	s.router.Handle("/admin/downloads/deleteall", middlewares.JwtMiddleware(http.HandlerFunc(downloadcontroller.DeleteAll))).Methods("DELETE")
	s.router.Handle("/admin/downloads/{id}", middlewares.JwtMiddleware(http.HandlerFunc(downloadcontroller.Delete))).Methods("DELETE")
	s.router.Handle("/admin/downloads/{page}/{qtd}", middlewares.JwtMiddleware(http.HandlerFunc(downloadcontroller.AdminFindAll))).Methods("GET")
}
