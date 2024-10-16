package web

import (
	"net/http"

	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/middlewares"
)

func (s *ServerWeb) DownloadController(downloadcontroller *controllers.DownloadController) {
	s.router.Handle("/downloads/category/{category}/{page}", middlewares.AccessOriginMiddleware(http.HandlerFunc(downloadcontroller.FindCategory))).Methods("GET")

	s.router.Handle("/downloads/topviews/{page}/{qtd}", middlewares.AccessOriginMiddleware(http.HandlerFunc(downloadcontroller.FindAllTopViews))).Methods("GET")

	s.router.Handle("/downloads/search/{page}", middlewares.AccessOriginMiddleware(http.HandlerFunc(downloadcontroller.Search))).Methods("GET")

	s.router.Handle("/download/{slug}", middlewares.AccessOriginMiddleware(http.HandlerFunc(downloadcontroller.GetBySlug))).Methods("GET")

	s.router.Handle("/downloads/{page}/{qtd}", middlewares.AccessOriginMiddleware(http.HandlerFunc(downloadcontroller.FindAll))).Methods("GET")

	s.router.Handle("/admin/downloads/{slug}", middlewares.JwtMiddleware(http.HandlerFunc(downloadcontroller.AdminGetBySlug))).Methods("GET")

	s.router.HandleFunc("/admin/downloads/create", downloadcontroller.CreateDownloadUsingTheForm).Methods("POST")

	s.router.HandleFunc("/admin/downloads/update/{slug}", downloadcontroller.UpdateDownloadUsingTheForm).Methods("POST")

	s.router.Handle("/admin/downloads/deleteall", middlewares.JwtMiddleware(http.HandlerFunc(downloadcontroller.DeleteAll))).Methods("DELETE")

	s.router.Handle("/admin/downloads/{id}", middlewares.JwtMiddleware(http.HandlerFunc(downloadcontroller.Delete))).Methods("DELETE")

	s.router.Handle("/admin/downloads/{page}/{qtd}", middlewares.JwtMiddleware(http.HandlerFunc(downloadcontroller.AdminFindAll))).Methods("GET")
}
