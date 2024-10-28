package port

import (
	"context"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
)

type BannerRepository interface {
	Create(ctx context.Context, banner entity.Banner) (entity.BannerDTO, error)
	Update(ctx context.Context, banner entity.Banner) (entity.BannerDTO, error)
	GetByID(ctx context.Context, id string) (entity.BannerDTO, error)
	GetByTag(ctx context.Context, tag string) (entity.BannerDTO, error)
	FindAll(ctx context.Context) ([]entity.BannerDTO, error)
	AdminFindAll(ctx context.Context) ([]entity.BannerDTO, error)
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, banners []entity.BannerDTO) error
}
