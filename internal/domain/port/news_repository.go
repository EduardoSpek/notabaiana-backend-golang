package port

import (
	"context"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
)

type NewsRepository interface {
	Create(news *entity.News) (*entity.News, error)
	Update(news *entity.News) (*entity.News, error)
	FindAll(ctx context.Context, page, limit int) ([]*entity.News, error)
	AdminFindAll(page, limit int) ([]*entity.News, error)
	FindCategory(category string, page int) ([]*entity.News, error)
	FindRecent() (*entity.News, error)
	NewsExists(title string) error
	GetBySlug(ctx context.Context, slug string) (*entity.News, error)
	GetByID(id string) (*entity.News, error)
	AdminGetBySlug(slug string) (*entity.News, error)
	NewsTruncateTable() error
	FindAllViews() ([]*entity.News, error)
	ClearViews() error
	SearchNews(page int, str_search string) ([]*entity.News, error)
	ClearImagePath(id string) error
	GetTotalNews() int
	GetTotalNewsVisible() int
	GetTotalNewsBySearch(str_search string) int
	GetTotalNewsByCategory(category string) int
	NewsMake() (*entity.News, error)
	CleanNews() ([]*entity.News, error)
	CleanNewsOld() ([]*entity.News, error)
	Delete(id string) error
	DeleteAll(news []*entity.News) error
	SetVisible(visible bool, id string) error
}
