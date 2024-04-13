package web

import "github.com/eduardospek/bn-api/internal/controllers"

func (s *ServerWeb) CrawlerController(crawlercontroller controllers.CrawlerController) {
	s.router.HandleFunc("/crawler/{key}", crawlercontroller.Crawler).Methods("GET")
}