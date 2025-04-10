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
	ID            string `json:"id"`
	Title         string `json:"title"`
	Link          string `json:"link"`
	Html          string `json:"html"`
	Image1        string `json:"banner_pc"`
	Image2        string `json:"banner_tablet"`
	Image3        string `json:"banner_mobile"`
	Tag           string `json:"tag"`
	Visible       bool   `json:"visible"`
	VisibleImage1 bool   `json:"visible_image1"`
	VisibleImage2 bool   `json:"visible_image2"`
	VisibleImage3 bool   `json:"visible_image3"`
}

type Banner struct {
	gorm.Model

	ID            string    `gorm:"column:id;primaryKey" json:"id"`
	Title         string    `gorm:"column:title" json:"title"`
	Link          string    `gorm:"column:link" json:"link"`
	Html          string    `gorm:"column:html" json:"html"`
	Image1        string    `gorm:"column:image1" json:"banner_pc"`
	Image2        string    `gorm:"column:image2" json:"banner_tablet"`
	Image3        string    `gorm:"column:image3" json:"banner_mobile"`
	Tag           string    `gorm:"column:tag" json:"tag"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
	Visible       bool      `gorm:"column:visible;" json:"visible"`
	VisibleImage1 bool      `gorm:"column:visible_image1;default:true" json:"visible_image1"`
	VisibleImage2 bool      `gorm:"column:visible_image2;default:true" json:"visible_image2"`
	VisibleImage3 bool      `gorm:"column:visible_image3;default:true" json:"visible_image3"`
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

func (b *Banner) ToDTO() BannerDTO {
	return BannerDTO{
		ID:            b.ID,
		Title:         b.Title,
		Link:          b.Link,
		Html:          b.Html,
		Image1:        b.Image1,
		Image2:        b.Image2,
		Image3:        b.Image3,
		Tag:           b.Tag,
		Visible:       b.Visible,
		VisibleImage1: b.VisibleImage1,
		VisibleImage2: b.VisibleImage2,
		VisibleImage3: b.VisibleImage3,
	}
}

func UpdateBanner(banner BannerDTO) *Banner {

	return &Banner{
		ID:            strings.TrimSpace(banner.ID),
		Title:         strings.TrimSpace(banner.Title),
		Link:          strings.TrimSpace(banner.Link),
		Html:          strings.TrimSpace(banner.Html),
		Image1:        strings.TrimSpace(banner.Image1),
		Image2:        strings.TrimSpace(banner.Image2),
		Image3:        strings.TrimSpace(banner.Image3),
		Tag:           strings.TrimSpace(banner.Tag),
		UpdatedAt:     time.Now(),
		Visible:       banner.Visible,
		VisibleImage1: banner.VisibleImage1,
		VisibleImage2: banner.VisibleImage2,
		VisibleImage3: banner.VisibleImage3,
	}
}

func (b *Banner) Validations() (bool, error) {

	if b.Title == "" || len(b.Title) < 2 || len(b.Title) > 80 {
		return false, ErrTitle
	}

	return true, nil

}
