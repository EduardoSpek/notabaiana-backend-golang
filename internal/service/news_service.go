package service

import "github.com/eduardospek/bn-api/internal/domain/entity"

type NewsRepository interface {
	Create(news entity.News) error
	FindAll() []entity.News
}

type NewsService struct {
	newsrepository NewsRepository
}

func NewNewsService(repository NewsRepository) *NewsService {
	return &NewsService{ newsrepository: repository }
}

func (s *NewsService) CreateNews(news entity.News) error {
	err := s.newsrepository.Create(news)
	if err != nil {
		return err
	}
	return nil
}

func (s *NewsService) FindAllNews() []entity.News {
	news := s.newsrepository.FindAll()
	return news
}