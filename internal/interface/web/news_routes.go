package web

import (
	"net/http"

	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/middlewares"
)

func (s *ServerWeb) NewsController(newscontroller *controllers.NewsController) {
	s.router.HandleFunc("/news/image", newscontroller.NewsImage).Methods("GET")

	s.router.HandleFunc("/truncate/news/{key}", newscontroller.NewsTruncateTable).Methods("GET")

	s.router.HandleFunc("/clean/news/{key}", newscontroller.CleanNews).Methods("GET")

	s.router.HandleFunc("/make/news/{key}", newscontroller.NewsMake).Methods("GET")

	s.router.Handle("/news", middlewares.AccessOriginMiddleware(http.HandlerFunc(newscontroller.News))).Methods("GET")

	s.router.Handle("/news/category/{category}/{page}", middlewares.AccessOriginMiddleware(http.HandlerFunc(newscontroller.NewsCategory))).Methods("GET")

	s.router.Handle("/news/search/{page}", middlewares.AccessOriginMiddleware(http.HandlerFunc(newscontroller.SearchNews))).Methods("GET")

	s.router.Handle("/news/{slug}", middlewares.AccessOriginMiddleware(http.HandlerFunc(newscontroller.GetNewsBySlug))).Methods("GET")

	s.router.Handle("/news/{page}/{qtd}", middlewares.AccessOriginMiddleware(http.HandlerFunc(newscontroller.News))).Methods("GET")

	s.router.Handle("/admin/news/{slug}", middlewares.JwtMiddleware(http.HandlerFunc(newscontroller.AdminGetNewsBySlug))).Methods("GET")

	s.router.Handle("/admin/news/{page}/{qtd}", middlewares.JwtMiddleware(http.HandlerFunc(newscontroller.AdminNews))).Methods("GET")

	s.router.HandleFunc("/admin/news/create", newscontroller.CreateNewsUsingTheForm).Methods("POST")

	s.router.HandleFunc("/admin/news/update/{slug}", newscontroller.UpdateNewsUsingTheForm).Methods("POST")

	s.router.Handle("/admin/news/deleteall", middlewares.JwtMiddleware(http.HandlerFunc(newscontroller.AdminDeleteAllNews))).Methods("DELETE")

	s.router.Handle("/admin/news/{id}", middlewares.JwtMiddleware(http.HandlerFunc(newscontroller.DeleteNews))).Methods("DELETE")
}
