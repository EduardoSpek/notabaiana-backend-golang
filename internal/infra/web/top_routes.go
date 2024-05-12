package web

import "github.com/eduardospek/notabaiana-backend-golang/internal/controllers"

func (s *ServerWeb) TopController(topcontroller controllers.TopController) {
	s.router.HandleFunc("/top/{key}", topcontroller.CreateTop).Methods("GET")
	s.router.HandleFunc("/top", topcontroller.FindAll).Methods("GET")
}