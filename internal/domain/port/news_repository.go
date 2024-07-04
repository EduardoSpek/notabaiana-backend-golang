package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type NewsRepository interface {
	Create(news entity.News) (entity.News, error)
	Update(news entity.News) (entity.News, error)
	FindAll(page, limit int) ([]entity.News, error)
	FindCategory(category string, page int) ([]entity.News, error)
	FindRecent() (entity.News, error)
	NewsExists(title string) error
	GetBySlug(slug string) (entity.News, error)
	NewsTruncateTable() error
	FindAllViews() ([]entity.News, error)
	ClearViews() error
	SearchNews(page int, str_search string) []entity.News
	ClearImagePath(id string) error
	GetTotalNews() int
	GetTotalNewsBySearch(str_search string) int
	GetTotalNewsByCategory(category string) int
	NewsMake() (entity.News, error)
	CleanNews()
}
