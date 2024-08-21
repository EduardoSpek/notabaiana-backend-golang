package web

import (
	"net/http"

	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/middlewares"
)

func (s *ServerWeb) UserController(usercontroller *controllers.UserController) {
	//s.router.HandleFunc("/user", usercontroller.CreateUser).Methods("POST")
	s.router.Handle("/user/{id}", middlewares.JwtMiddleware(http.HandlerFunc(usercontroller.UpdateUser))).Methods("PUT")
	s.router.HandleFunc("/login", usercontroller.Login).Methods("POST")
	s.router.HandleFunc("/accesscheck", usercontroller.AccessCheck).Methods("GET")
}
