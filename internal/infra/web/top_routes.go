package web

import "github.com/eduardospek/bn-api/internal/controllers"

func (s *ServerWeb) TopController(topcontroller controllers.TopController) {
	s.router.HandleFunc("/top/{key}", topcontroller.CreateTop).Methods("GET")	
}