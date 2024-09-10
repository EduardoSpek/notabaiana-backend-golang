package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type NewsFindAllOutput struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	TitleAi   string    `json:"title_ai"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
	Link      string    `json:"link"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type News struct {
	gorm.Model

	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Title     string    `gorm:"column:title;index:idx_visible_title" json:"title"`
	TitleAi   string    `gorm:"column:title_ai" json:"title_ai"`
	Text      string    `gorm:"column:text" json:"text"`
	Link      string    `gorm:"column:link" json:"link"`
	Image     string    `gorm:"column:image" json:"image"`
	Slug      string    `gorm:"column:slug;index:idx_slug_visible" json:"slug"`
	CreatedAt time.Time `gorm:"column:created_at;index:,sort:desc" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Visible   bool      `gorm:"column:visible;index:idx_visible_title" json:"visible"`
	TopStory  bool      `gorm:"column:topstory;" json:"topstory"`
	Views     int       `gorm:"column:views;default:0;index:,sort:desc" json:"views"`
	Category  string    `gorm:"column:category;index:idx_visible_category" json:"category"`
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
		TopStory:  false,
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
		TopStory:  news.TopStory,
		Category:  news.Category,
		Make:      news.Make,
	}
}

func SlugTitle(title string) string {
	return slug.Make(strings.TrimSpace(title))
}
