package entity

import (
	"time"

	"gorm.io/gorm"
)

type Top struct {
	gorm.Model
	
	Title string `gorm:"column:title" json:"title"`	
	Link  string `gorm:"column:link" json:"link"`
	Image string `gorm:"column:image" json:"image"`	
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`		
	Views int `gorm:"column:views;default:0"`
}

func NewTop(top Top) *Top {
	return &Top{
		Title: top.Title,
		Link: top.Link,
		Image: top.Image,
		CreatedAt: top.CreatedAt,
		Views: top.Views,
	}
}