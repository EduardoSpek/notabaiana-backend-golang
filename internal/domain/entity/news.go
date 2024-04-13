package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type News struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
	Link  string `json:"link"`
	Image string `json:"image"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewNews(news News) *News {
	return &News{
		ID:    uuid.NewString(),
		Title: strings.TrimSpace(news.Title),
		Text:  strings.TrimSpace(news.Text),
		Link:  strings.TrimSpace(news.Link),
		Image: strings.TrimSpace(news.Image),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),	
	}
}