package memorydb

import (
	"errors"
	"strings"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
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

func (r *NewsMemoryRepository) FindAll(page, limit int) (interface{}, error) {
	var news []entity.News
	for _, n := range r.Newsdb {
		news = append(news, n)
	}

	pagination := utils.Pagination(page, len(r.Newsdb))

    result := struct{
        List_news []entity.News `json:"news"`
        Pagination map[string][]int `json:"pagination"`
    }{
        List_news: news,
        Pagination: pagination,
    }

	return result, nil
}

func (r *NewsMemoryRepository) FindCategory(category string, page int) (interface{}, error) {
	var news []entity.News
	for _, n := range r.Newsdb {
		if category == n.Category {
			news = append(news, n)
		}
	}

	pagination := utils.Pagination(page, len(r.Newsdb))

    result := struct{
        List_news []entity.News `json:"news"`
        Pagination map[string][]int `json:"pagination"`
		Category string `json:"category"`
    }{
        List_news: news,
        Pagination: pagination,
		Category: category,
    }

	return result, nil
}

func (r *NewsMemoryRepository) SearchNews(page int, str_search string) interface{} {
	var news []entity.News
	for _, n := range r.Newsdb {
		if strings.Contains(n.Title, str_search) {
			news = append(news, n)
		}
	}

	total := len(r.Newsdb)

    pagination := utils.Pagination(page, int(total))

    result := struct{
        List_news []entity.News `json:"news"`
        Pagination map[string][]int `json:"pagination"`
        Search string `json:"search"`
    }{
        List_news: news,
        Pagination: pagination,
        Search: str_search,
    }

	return result
}

func (r *NewsMemoryRepository) FindAllViews() ([]entity.News, error) {
	var news []entity.News
	for _, n := range r.Newsdb {
		news = append(news, n)
	}
	return news, nil
}
func (r *NewsMemoryRepository) NewsTruncateTable() error {
    
	r.Newsdb = make(map[string]entity.News)

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