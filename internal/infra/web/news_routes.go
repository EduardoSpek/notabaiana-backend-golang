package web

import "github.com/eduardospek/bn-api/internal/controllers"

func (s *ServerWeb) NewsController(newscontroller controllers.NewsController) {
	s.router.HandleFunc("/news/truncate/{key}", newscontroller.NewsTruncateTable).Methods("GET")
	s.router.HandleFunc("/news", newscontroller.News).Methods("GET")
	s.router.HandleFunc("/news/{slug}", newscontroller.GetNewsBySlug).Methods("GET")
	s.router.HandleFunc("/news/{page}/{qtd}", newscontroller.News).Methods("GET")	
}