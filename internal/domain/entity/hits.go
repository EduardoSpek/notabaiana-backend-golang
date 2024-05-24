package entity

import "gorm.io/gorm"

type Hits struct {
	gorm.Model

	IP string `gorm:"column:ip" json:"ip"`
	Session string `gorm:"column:session" json:"session"`
	Views int `gorm:"column:views" json:"views"`

}

func NewHits(hit Hits) *Hits {
	return &Hits{
		IP: hit.IP,
		Session: hit.Session,
		Views: hit.Views,
	}
}