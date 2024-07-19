package web

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"
)

func (s *ServerWeb) BannerController(bannercontroller controllers.BannerController) {
	s.router.HandleFunc("/banners/create", bannercontroller.CreateBannerUsingTheForm).Methods("POST")
	s.router.HandleFunc("/banners", bannercontroller.BannerList).Methods("GET")
}
