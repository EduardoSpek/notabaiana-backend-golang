package entity

import "github.com/google/uuid"

type News struct {
	ID    string
	Title string
	Text  string
	Link  string
	Image string
}

func NewNews(news News) *News {
	return &News{
		ID:    uuid.NewString(),
		Title: news.Title,
		Text:  news.Text,
		Link:  news.Link,
		Image: news.Image,
	}
}