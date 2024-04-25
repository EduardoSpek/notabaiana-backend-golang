package supabase

import (
	"fmt"
	"os"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (repo *Postgres) Connect() (*gorm.DB, error) {
	
	connStr := "user="+ os.Getenv("POSTGRES_USERNAME") +" password="+ os.Getenv("POSTGRES_PASSWORD") +" host="+ os.Getenv("POSTGRES_HOST") +" port="+ os.Getenv("POSTGRES_PORT") +" dbname="+ os.Getenv("POSTGRES_DB") +""
	
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db.AutoMigrate(&entity.News{}, &entity.Top{})

	return db, nil
}