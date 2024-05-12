package sqlite

import (
	"fmt"
	"os"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLite struct {}

var conn SQLite

func (repo *SQLite) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("PATCH_DB_SQLITE")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db.AutoMigrate(&entity.News{})

	return db, nil
}