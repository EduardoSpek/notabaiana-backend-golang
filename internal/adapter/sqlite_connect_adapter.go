package adapter

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLiteAdapter struct {
	db *gorm.DB
}

func NewSQLiteAdapter() (*SQLiteAdapter, error) {
	adapter := &SQLiteAdapter{}
	err := adapter.Connect()
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

func (repo *SQLiteAdapter) CloseDB() {
	db, err := repo.db.DB()

	if err != nil {
		fmt.Println(err)
	}

	db.Close()
}

func (repo *SQLiteAdapter) GetDB() *gorm.DB {
	return repo.db
}

func (repo *SQLiteAdapter) Connect() error {

	//connStr := "user=" + os.Getenv("POSTGRES_USERNAME") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " host=" + os.Getenv("POSTGRES_HOST") + " port=" + os.Getenv("POSTGRES_PORT") + " dbname=" + os.Getenv("POSTGRES_DB") + ""

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Limite para considerar uma consulta lenta
			LogLevel:                  logger.Silent, // Nível de log (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: true,          // Ignorar erros de "registro não encontrado"
			Colorful:                  false,         // Desativar saída colorida
		},
	)

	db, err := gorm.Open(sqlite.Open(os.Getenv("PATCH_DB_SQLITE")), &gorm.Config{
		Logger:      newLogger,
		PrepareStmt: false,
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	repo.db = db

	db.AutoMigrate(&entity.News{}, &entity.Top{}, &entity.Hits{}, &entity.User{}, &entity.Banner{}, &entity.Contact{}, &entity.Download{}, &entity.Music{})

	sqlDB, err := db.DB()

	if err != nil {
		fmt.Println(err)
		return err
	}

	sqlDB.SetMaxOpenConns(10) // número máximo de conexões abertas
	sqlDB.SetMaxIdleConns(5)  // número máximo de conexões ociosas
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}
