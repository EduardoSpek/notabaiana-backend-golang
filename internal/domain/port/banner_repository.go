package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type BannerRepository interface {
	Create(banner entity.Banner) (entity.BannerDTO, error)
	GetByID(id string) (entity.BannerDTO, error)
	GetByTag(tag string) (entity.BannerDTO, error)
	FindAll() ([]entity.BannerDTO, error)
}
