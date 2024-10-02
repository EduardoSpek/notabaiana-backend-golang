package memorydb

import (
	"errors"
	"strings"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
)

var (
	ErrNewExists  = errors.New("notícia já cadastrada com este título")
	ErrNotNewSlug = errors.New("não há notícia com este slug")
)

type NewsMemoryRepository struct {
	Newsdb map[string]*entity.News
}

func NewNewsMemoryRepository() *NewsMemoryRepository {
	return &NewsMemoryRepository{Newsdb: make(map[string]*entity.News)}
}

func (r *NewsMemoryRepository) CleanNews() ([]*entity.News, error) {

	return nil, nil

}

func (r *NewsMemoryRepository) NewsMake() (*entity.News, error) {
	return &entity.News{}, nil
}

func (r *NewsMemoryRepository) Update(news *entity.News) (*entity.News, error) {
	r.Newsdb[news.ID] = news
	return news, nil
}

func (r *NewsMemoryRepository) Create(news *entity.News) (*entity.News, error) {
	r.Newsdb[news.ID] = news
	return news, nil
}
func (r *NewsMemoryRepository) AdminGetBySlug(slug string) (*entity.News, error) {
	for _, n := range r.Newsdb {
		if slug == n.Slug {
			return n, nil
		}
	}
	return &entity.News{}, ErrNotNewSlug
}
func (r *NewsMemoryRepository) GetByID(id string) (*entity.News, error) {
	for _, n := range r.Newsdb {
		if id == n.ID {
			return n, nil
		}
	}
	return &entity.News{}, ErrNotNewSlug
}
func (r *NewsMemoryRepository) GetBySlug(slug string) (*entity.News, error) {
	for _, n := range r.Newsdb {
		if slug == n.Slug {
			return n, nil
		}
	}
	return &entity.News{}, ErrNotNewSlug
}

func (r *NewsMemoryRepository) AdminFindAll(page, limit int) ([]*entity.News, error) {
	var news []*entity.News
	for _, n := range r.Newsdb {
		news = append(news, n)
	}

	return news, nil
}

func (r *NewsMemoryRepository) FindAll(page, limit int) ([]*entity.News, error) {
	var news []*entity.News
	for _, n := range r.Newsdb {
		news = append(news, n)
	}

	return news, nil
}

func (r *NewsMemoryRepository) FindRecent() (*entity.News, error) {
	var news *entity.News

	var latestTime time.Time

	// Iterando sobre o mapa
	for _, v := range r.Newsdb {
		if v.CreatedAt.After(latestTime) {
			latestTime = v.CreatedAt
			news = v
		}
	}

	return news, nil
}

func (r *NewsMemoryRepository) FindCategory(category string, page int) ([]*entity.News, error) {
	var news []*entity.News
	for _, n := range r.Newsdb {
		if category == n.Category {
			news = append(news, n)
		}
	}

	return news, nil
}

func (r *NewsMemoryRepository) SearchNews(page int, str_search string) ([]*entity.News, error) {
	var news []*entity.News
	for _, n := range r.Newsdb {
		if strings.Contains(n.Title, str_search) {
			news = append(news, n)
		}
	}

	return news, nil
}

func (r *NewsMemoryRepository) GetTotalNewsBySearch(str_search string) int {
	var total int = 0
	for _, n := range r.Newsdb {
		if strings.Contains(n.Title, str_search) {
			total++
		}
	}

	return total
}

func (r *NewsMemoryRepository) GetTotalNewsByCategory(category string) int {
	var total int = 0
	for _, n := range r.Newsdb {
		if n.Category == category {
			total++
		}
	}

	return total
}

func (r *NewsMemoryRepository) Delete(id string) error {
	return nil
}
func (r *NewsMemoryRepository) DeleteAll(news []*entity.News) error {
	return nil
}

func (r *NewsMemoryRepository) GetTotalNews() int {
	total := len(r.Newsdb)

	return total

}

func (r *NewsMemoryRepository) GetTotalNewsVisible() int {
	total := len(r.Newsdb)

	return total

}

func (r *NewsMemoryRepository) FindAllViews() ([]*entity.News, error) {
	var news []*entity.News
	for _, n := range r.Newsdb {
		news = append(news, n)
	}
	return news, nil
}
func (r *NewsMemoryRepository) NewsTruncateTable() error {

	r.Newsdb = make(map[string]*entity.News)

	return nil
}

func (repo *NewsMemoryRepository) ClearViews() error {
	return nil
}

func (r *NewsMemoryRepository) ClearImagePath(id string) error {

	for _, n := range r.Newsdb {
		if id == n.ID {
			n.Image = ""
			r.Newsdb[n.ID] = n
			return nil
		}
	}

	return errors.New("não foi possível atualizar a Image")
}

// VALIDATIONS
func (r *NewsMemoryRepository) NewsExists(title string) error {
	title = strings.TrimSpace(title)
	for _, new := range r.Newsdb {
		if new.Title == title {
			return ErrNewExists
		}
	}
	return nil
}
