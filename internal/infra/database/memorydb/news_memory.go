package memorydb

import (
	"errors"
	"strings"

	"github.com/eduardospek/bn-api/internal/domain/entity"
)

var (
	ErrNewExists = errors.New("notícia já cadastrada com este título")
	ErrNotNewSlug = errors.New("não há notícia com este slug")
)

type NewsMemoryRepository struct {
	Newsdb map[string]entity.News
}

func NewNewsMemoryRepository() *NewsMemoryRepository {
	return &NewsMemoryRepository{ Newsdb: make(map[string]entity.News) }
}

func (r *NewsMemoryRepository) Create(news entity.News) (entity.News, error) {
	r.Newsdb[news.ID] = news
	return news, nil
}
func (r *NewsMemoryRepository) GetBySlug(slug string) (entity.News, error) {	
	for _, n := range r.Newsdb {
		if slug == n.Slug {
			return n, nil
		}
	}
	return entity.News{}, ErrNotNewSlug
}

func (r *NewsMemoryRepository) FindAll(page, limit int) ([]entity.News, error) {
	var news []entity.News
	for _, n := range r.Newsdb {
		news = append(news, n)
	}
	return news, nil
}

//VALIDATIONS
func (r *NewsMemoryRepository) NewsExists(title string) error {
	title = strings.TrimSpace(title)
    for _, new := range r.Newsdb {
        if new.Title == title {			
            return ErrNewExists
        }
    }	
    return nil
}