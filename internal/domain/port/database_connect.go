package port

import "gorm.io/gorm"

type DBAdapter interface {
	Connect() (*gorm.DB, error)
}