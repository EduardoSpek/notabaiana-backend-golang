package postgres

import (
	"errors"
	"sync"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"gorm.io/gorm"
)

var (
	ErrDownloadNotFound = errors.New("download não encontrado")
)

type DownloadPostgresRepository struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

func NewDownloadPostgresRepository(db_adapter port.DBAdapter) *DownloadPostgresRepository {
	db, _ := db_adapter.Connect()
	return &DownloadPostgresRepository{db: db}
}

func (repo *DownloadPostgresRepository) Create(download *entity.Download) (*entity.Download, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	result := repo.db.Create(&download)

	if result.Error != nil {
		tx.Rollback()
		return &entity.Download{}, result.Error
	}

	tx.Commit()

	return download, nil
}

func (repo *DownloadPostgresRepository) Update(download *entity.Download) (*entity.Download, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	result := repo.db.Model(&download).Updates(map[string]interface{}{
		"category":   download.Category,
		"title":      download.Title,
		"link":       download.Link,
		"text":       download.Text,
		"image":      download.Image,
		"visible":    download.Visible,
		"updated_at": download.UpdatedAt,
	})

	if result.Error != nil {
		tx.Rollback()
		return &entity.Download{}, result.Error
	}

	tx.Commit()

	return download, nil
}

func (repo *DownloadPostgresRepository) GetByLink(link string) (*entity.Download, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var download *entity.Download
	result := repo.db.Model(&entity.Download{}).Where("link = ?", link).First(&download)

	if result.Error != nil {
		return &entity.Download{}, result.Error
	}

	tx.Commit()

	return download, nil
}

func (repo *DownloadPostgresRepository) FindAll(page, limit int) ([]*entity.Download, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	offset := (page - 1) * limit

	tx := repo.db.Begin()
	defer tx.Rollback()

	var download []*entity.Download
	repo.db.Model(&entity.Download{}).Where("visible = true").Order("created_at DESC").Limit(limit).Offset(offset).Find(&download)

	if repo.db.Error != nil {
		return []*entity.Download{}, repo.db.Error
	}

	tx.Commit()

	return download, nil
}

func (repo *DownloadPostgresRepository) GetTotalVisible() int {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var total int64
	repo.db.Model(&entity.Download{}).Where("visible = true").Count(&total)

	return int(total)

}
