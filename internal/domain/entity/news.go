package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type News struct {
	gorm.Model

	ID    string `gorm:"column:id;primaryKey" json:"id"`
	Title string `gorm:"column:title" json:"title"`
	Text  string `gorm:"column:text" json:"text"`
	Link  string `gorm:"column:link" json:"link"`
	Image string `gorm:"column:image" json:"image"`
	Slug string `gorm:"column:slug" json:"slug"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	Visible bool `gorm:"column:visible;default:true" json:"-"`
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