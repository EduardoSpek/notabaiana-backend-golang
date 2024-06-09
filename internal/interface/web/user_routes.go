package web

import "github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"

func (s *ServerWeb) UserController(usercontroller controllers.UserController) {	
	s.router.HandleFunc("/user", usercontroller.CreateUser).Methods("POST")	
}