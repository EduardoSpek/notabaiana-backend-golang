package port

import (
	"gorm.io/gorm"
)

type DBAdapter interface {
	Connect() error
	GetDB() *gorm.DB
}
