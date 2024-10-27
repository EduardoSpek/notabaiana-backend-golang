package adapter

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresAdapter struct{}

func NewPostgresAdapter() *PostgresAdapter {
	return &PostgresAdapter{}
}

func (repo *PostgresAdapter) Connect() (*gorm.DB, error) {

	connStr := "user=" + os.Getenv("POSTGRES_USERNAME") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " host=" + os.Getenv("POSTGRES_HOST") + " port=" + os.Getenv("POSTGRES_PORT") + " dbname=" + os.Getenv("POSTGRES_DB") + ""

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Limite para considerar uma consulta lenta
			LogLevel:                  logger.Silent, // Nível de log (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: true,          // Ignorar erros de "registro não encontrado"
			Colorful:                  false,         // Desativar saída colorida
		},
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger:      newLogger,
		PrepareStmt: false,
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db.AutoMigrate(&entity.News{}, &entity.Top{}, &entity.Hits{}, &entity.User{}, &entity.Banner{}, &entity.Contact{}, &entity.Download{}, &entity.Music{})

	return db, nil
}
