package web

import (
	"net/http"

	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/middlewares"
)

func (s *ServerWeb) ContactController(contactcontroller controllers.ContactController) {
	s.router.HandleFunc("/contacts/create", contactcontroller.CreateForm).Methods("POST")
	s.router.HandleFunc("/admin/contacts/create", contactcontroller.AdminCreateForm).Methods("POST")
	s.router.HandleFunc("/admin/contacts/{id}", contactcontroller.AdminGetByID).Methods("GET")
	s.router.Handle("/admin/contacts", middlewares.JwtMiddleware(http.HandlerFunc(contactcontroller.AdminFindAll))).Methods("GET")
	s.router.Handle("/admin/contacts/deleteall", middlewares.JwtMiddleware(http.HandlerFunc(contactcontroller.AdminDeleteAll))).Methods("DELETE")
	s.router.Handle("/admin/contacts/{id}", middlewares.JwtMiddleware(http.HandlerFunc(contactcontroller.AdminDelete))).Methods("DELETE")
}
