package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrLink = errors.New("o link não é válido")
)

type Download struct {
	gorm.Model

	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Category  string    `gorm:"column:category;index:idx_visible_category" json:"category"`
	Title     string    `gorm:"column:title;index:idx_visible_title" json:"title"`
	Text      string    `gorm:"column:text" json:"text"`
	Link      string    `gorm:"column:link;unique;index:idx_link_visible" json:"link"`
	Image     string    `gorm:"column:image" json:"image"`
	Slug      string    `gorm:"column:slug;unique;index:idx_slug_visible" json:"slug"`
	Views     int       `gorm:"column:views;default:0;index:,sort:desc" json:"views"`
	Downloads int       `gorm:"column:downloads;default:0" json:"downloads"`
	Visible   bool      `gorm:"column:visible;default:true;index:idx_visible_title" json:"visible"`
	Make      bool      `gorm:"column:make;default:false" json:"make"`
	CreatedAt time.Time `gorm:"column:created_at;index:,sort:desc" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Musics    []*Music  `gorm:"foreignKey:DownloadID" json:"musics"`
}

type Music struct {
	gorm.Model

	File       string `json:"file"`
	Path       string `gorm:"column:path;unique" json:"path"`
	Position   int    `json:"position"`
	DownloadID string `gorm:"type:string"`
}

func NewDownload(download Download) *Download {

	slug := SlugTitle(download.Title)

	return &Download{
		ID:        uuid.NewString(),
		Category:  download.Category,
		Title:     strings.TrimSpace(download.Title),
		Text:      strings.TrimSpace(download.Text),
		Link:      strings.TrimSpace(download.Link),
		Image:     strings.TrimSpace(download.Image),
		Slug:      slug,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Musics:    download.Musics,
	}
}

func UpdateDownload(download Download) *Download {

	slug := SlugTitle(download.Title)

	return &Download{
		ID:        strings.TrimSpace(download.ID),
		Category:  download.Category,
		Title:     strings.TrimSpace(download.Title),
		Text:      strings.TrimSpace(download.Text),
		Link:      strings.TrimSpace(download.Link),
		Image:     strings.TrimSpace(download.Image),
		Slug:      slug,
		Visible:   download.Visible,
		Make:      download.Make,
		UpdatedAt: time.Now(),
	}
}

func (v *Download) Validations() (bool, error) {

	if v.Title == "" || len(v.Title) < 2 || len(v.Title) > 144 {
		return false, ErrTitle
	}

	if v.Link == "" || len(v.Link) < 2 || len(v.Link) > 144 {
		return false, ErrLink
	}

	isValidUrl := utils.IsValidURL(v.Link)

	if !isValidUrl {
		return false, ErrLink
	}

	return true, nil

}
