package web

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"
)

func (s *ServerWeb) DownloadController(downloadcontroller *controllers.DownloadController) {
	s.router.HandleFunc("/downloads/category/{category}/{page}", downloadcontroller.FindCategory).Methods("GET")
	s.router.HandleFunc("/downloads/search/{page}", downloadcontroller.Search).Methods("GET")
	s.router.HandleFunc("/download/{slug}", downloadcontroller.GetBySlug).Methods("GET")
	s.router.HandleFunc("/downloads/{page}/{qtd}", downloadcontroller.FindAll).Methods("GET")
	s.router.HandleFunc("/admin/downloads/create", downloadcontroller.CreateDownloadUsingTheForm).Methods("POST")
}
