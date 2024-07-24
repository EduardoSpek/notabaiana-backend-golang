package web

import (
	"net/http"

	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/middlewares"
)

func (s *ServerWeb) BannerController(bannercontroller controllers.BannerController) {
	s.router.HandleFunc("/admin/banners/update/{id}", bannercontroller.UpdateBannerUsingTheForm).Methods("POST")
	s.router.HandleFunc("/admin/banners/create", bannercontroller.CreateBannerUsingTheForm).Methods("POST")
	s.router.HandleFunc("/admin/banners/{id}", bannercontroller.FindBanner).Methods("GET")
	s.router.Handle("/admin/banners", middlewares.JwtMiddleware(http.HandlerFunc(bannercontroller.AdminBannerList))).Methods("GET")
	s.router.HandleFunc("/banners", bannercontroller.BannerList).Methods("GET")
	s.router.Handle("/admin/banners/deleteall", middlewares.JwtMiddleware(http.HandlerFunc(bannercontroller.AdminDeleteAllBanner))).Methods("DELETE")
	s.router.Handle("/admin/banners/{id}", middlewares.JwtMiddleware(http.HandlerFunc(bannercontroller.DeleteBanner))).Methods("DELETE")
}
