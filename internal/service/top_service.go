package service

import (
	"time"

	"github.com/eduardospek/bn-api/internal/domain/entity"
)

type TopRepository interface {
	Create(tops []entity.Top) error
	FindAll() ([]entity.Top, error)
}

type TopService struct {
	TopRepository TopRepository
	NewsService NewsService
}

func NewTopService(newsservice NewsService) *TopService {
	return &TopService{  NewsService: newsservice }
}

func (t *TopService) TopCreate() error {

	news, err := t.NewsService.FindAllViews()

	if err != nil {
		return err
	}

	var tops []entity.Top
	
	for _, top := range news {
		
		newtop := &entity.Top{
			Title: top.Title,
			Link: top.Link,
			Image: top.Image,
			CreatedAt: top.CreatedAt,
		}

		tops = append(tops, *newtop)
	}

	err = t.TopRepository.Create(tops)

	if err != nil {
		return err
	}

	return nil
}

func (t *TopService) FindAll() ([]entity.Top, error) {
	tops, err := t.TopRepository.FindAll()

	if err != nil {
		return []entity.Top{}, err
	}

	return tops, nil
}

func (t *TopService) Start(minutes time.Duration) {
	
	go t.TopCreate()

	ticker := time.NewTicker(minutes * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
		go t.TopCreate()
	}
}