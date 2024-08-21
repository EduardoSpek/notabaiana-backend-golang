package web

import "github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"

func (s *ServerWeb) CrawlerController(crawlercontroller *controllers.CrawlerController) {
	s.router.HandleFunc("/crawler/{key}", crawlercontroller.Crawler).Methods("GET")
}
