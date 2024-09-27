package adapter

import (
	"fmt"
	"os"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresAdapter struct{}

func NewPostgresAdapter() *PostgresAdapter {
	return &PostgresAdapter{}
}

func (repo *PostgresAdapter) Connect() (*gorm.DB, error) {

	connStr := "user=" + os.Getenv("POSTGRES_USERNAME") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " host=" + os.Getenv("POSTGRES_HOST") + " port=" + os.Getenv("POSTGRES_PORT") + " dbname=" + os.Getenv("POSTGRES_DB") + ""

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db.AutoMigrate(&entity.News{}, &entity.Top{}, &entity.Hits{}, &entity.User{}, &entity.Banner{}, &entity.Contact{}, &entity.Music{}, &entity.Download{})

	return db, nil
}
