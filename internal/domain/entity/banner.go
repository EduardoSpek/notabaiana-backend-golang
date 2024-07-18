package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrTitle = errors.New("title não pode ser vazio e deve ter mínimo de 2 e máximo de 80 caracteres")
)

// Input and Output DTO
type BannerDTO struct {
	ID      string
	Title   string
	Link    string
	Html    string
	Image1  string
	Image2  string
	Image3  string
	Tag     string
	Visible bool
}

type Banner struct {
	gorm.Model

	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Title     string    `gorm:"column:title" json:"title"`
	Link      string    `gorm:"column:link" json:"link"`
	Html      string    `gorm:"column:html" json:"html"`
	Image1    string    `gorm:"column:image1" json:"banner_pc"`
	Image2    string    `gorm:"column:image2" json:"banner_tablet"`
	Image3    string    `gorm:"column:image3" json:"banner_mobile"`
	Tag       string    `gorm:"column:tag" json:"tag"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Visible   bool      `gorm:"column:visible;" json:"visible"`
}

func NewBanner(banner BannerDTO) *Banner {

	return &Banner{
		ID:        uuid.NewString(),
		Title:     strings.TrimSpace(banner.Title),
		Link:      strings.TrimSpace(banner.Link),
		Html:      strings.TrimSpace(banner.Html),
		Image1:    strings.TrimSpace(banner.Image1),
		Image2:    strings.TrimSpace(banner.Image2),
		Image3:    strings.TrimSpace(banner.Image3),
		Tag:       strings.TrimSpace(banner.Tag),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Visible:   banner.Visible,
	}
}

func UpdateBanner(banner BannerDTO) *Banner {

	return &Banner{
		ID:        strings.TrimSpace(banner.ID),
		Title:     strings.TrimSpace(banner.Title),
		Link:      strings.TrimSpace(banner.Link),
		Html:      strings.TrimSpace(banner.Html),
		Image1:    strings.TrimSpace(banner.Image1),
		Image2:    strings.TrimSpace(banner.Image2),
		Image3:    strings.TrimSpace(banner.Image3),
		Tag:       strings.TrimSpace(banner.Tag),
		UpdatedAt: time.Now(),
		Visible:   banner.Visible,
	}
}

func (b *Banner) Validations() (bool, error) {

	if b.Title == "" || len(b.Title) < 2 || len(b.Title) > 80 {
		return false, ErrTitle
	}

	return true, nil

}
