package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type News struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
	Link  string `json:"link"`
	Image string `json:"image"`
	Slug string `json:"slug"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Visible bool `json:"-"`
}

func NewNews(news News) *News {
	return &News{
		ID:    uuid.NewString(),
		Title: strings.TrimSpace(news.Title),
		Text:  strings.TrimSpace(news.Text),
		Link:  strings.TrimSpace(news.Link),
		Image: strings.TrimSpace(news.Image),
		Slug: SlugTitle(news.Title),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Visible: news.Visible,
	}
}

func SlugTitle(title string) string {
	return slug.Make(strings.TrimSpace(title))
}