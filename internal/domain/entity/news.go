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

	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Title     string    `gorm:"column:title" json:"title"`
	TitleAi   string    `gorm:"column:title_ai" json:"title_ai"`
	Text      string    `gorm:"column:text" json:"text"`
	Link      string    `gorm:"column:link" json:"link"`
	Image     string    `gorm:"column:image" json:"image"`
	Slug      string    `gorm:"column:slug" json:"slug"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Visible   bool      `gorm:"column:visible;" json:"visible"`
	Views     int       `gorm:"column:views;default:0" json:"views"`
	Category  string    `gorm:"column:category" json:"category"`
	Make      bool      `gorm:"column:make;default:false" json:"make"`
}

func NewNews(news News) *News {

	var slug string

	if news.TitleAi != "" {
		slug = SlugTitle(news.TitleAi)
	} else {
		slug = SlugTitle(news.Title)
	}

	return &News{
		ID:        uuid.NewString(),
		Title:     strings.TrimSpace(news.Title),
		TitleAi:   strings.TrimSpace(news.TitleAi),
		Text:      strings.TrimSpace(news.Text),
		Link:      strings.TrimSpace(news.Link),
		Image:     strings.TrimSpace(news.Image),
		Slug:      slug,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Visible:   news.Visible,
		Category:  news.Category,
	}
}

func UpdateNews(news News) *News {

	var slug string

	if news.TitleAi != "" {
		slug = SlugTitle(news.TitleAi)
	} else {
		slug = SlugTitle(news.Title)
	}

	return &News{
		ID:        strings.TrimSpace(news.ID),
		Title:     strings.TrimSpace(news.Title),
		TitleAi:   strings.TrimSpace(news.TitleAi),
		Text:      strings.TrimSpace(news.Text),
		Link:      strings.TrimSpace(news.Link),
		Image:     strings.TrimSpace(news.Image),
		Slug:      slug,
		UpdatedAt: time.Now(),
		Visible:   news.Visible,
		Category:  news.Category,
		Make:      news.Make,
	}
}

func SlugTitle(title string) string {
	return slug.Make(strings.TrimSpace(title))
}
