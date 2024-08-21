package service

import (
	"fmt"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

type TopService struct {
	HitsRepository port.HitsRepository
	TopRepository  port.TopRepository
	NewsRepository port.NewsRepository
	NewsService    *NewsService
}

func NewTopService(toprepo port.TopRepository, newsrepo port.NewsRepository, hitsrepo port.HitsRepository, newsservice *NewsService) *TopService {
	return &TopService{TopRepository: toprepo, NewsRepository: newsrepo, HitsRepository: hitsrepo, NewsService: newsservice}
}

func (t *TopService) TopCreate() {

	// var news []entity.News
	// hits, _ := t.HitsRepository.TopHits()

	// for _, hit := range hits {

	// 	new, err := t.NewsRepository.GetBySlug(hit.Session)

	// 	if err != nil {
	// 		return
	// 	}

	// 	news = append(news, new)
	// }

	news, err := t.NewsService.FindAllViews()

	if err != nil {
		fmt.Println(err)
	}

	var tops []entity.Top
	var newtop entity.Top
	var ntop entity.Top

	for _, top := range news {

		newtop = entity.Top{
			Title:     top.Title,
			TitleAi:   top.TitleAi,
			Link:      top.Link,
			Image:     top.Image,
			CreatedAt: top.CreatedAt,
			Views:     top.Views,
		}

		ntop = *entity.NewTop(newtop)

		tops = append(tops, ntop)
	}

	err = t.TopRepository.TopTruncateTable()

	if err != nil {
		fmt.Println(err)
	}

	err = t.TopRepository.Create(tops)

	if err != nil {
		fmt.Println(err)
	}

	err = t.NewsService.ClearViews()

	if err != nil {
		fmt.Println(err)
	}
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
