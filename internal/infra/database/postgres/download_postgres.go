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
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer tx.Rollback()

	if err := tx.Create(download).Error; err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return download, nil
}

func (repo *DownloadPostgresRepository) Update(download *entity.Download) (*entity.Download, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer tx.Rollback()

	result := tx.Model(download).Updates(map[string]interface{}{
		"category":   download.Category,
		"title":      download.Title,
		"link":       download.Link,
		"text":       download.Text,
		"image":      download.Image,
		"visible":    download.Visible,
		"updated_at": download.UpdatedAt,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return download, nil
}

func (repo *DownloadPostgresRepository) GetByLink(link string) (*entity.Download, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var download entity.Download
	result := repo.db.Where("link = ?", link).First(&download)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &download, nil
}

func (repo *DownloadPostgresRepository) GetBySlug(slug string) (*entity.Download, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer tx.Rollback()

	var download entity.Download
	if err := tx.Where("visible = true AND slug = ?", slug).First(&download).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	download.Views++
	if err := tx.Save(&download).Error; err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &download, nil
}

func (repo *DownloadPostgresRepository) FindAll(page, limit int) ([]*entity.Download, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	offset := (page - 1) * limit
	var downloads []*entity.Download

	err := repo.db.Where("visible = true").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&downloads).Error

	if err != nil {
		return nil, err
	}

	return downloads, nil
}

func (repo *DownloadPostgresRepository) FindCategory(category string, page int) ([]*entity.Download, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	limit := PerPage
	offset := (page - 1) * limit
	var downloads []*entity.Download

	err := repo.db.Where("visible = true AND category = ?", category).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&downloads).Error

	if err != nil {
		return nil, err
	}

	return downloads, nil
}

func (repo *DownloadPostgresRepository) GetTotalVisible() (int, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var total int64
	err := repo.db.Model(&entity.Download{}).Where("visible = true").Count(&total).Error
	return int(total), err
}

func (repo *DownloadPostgresRepository) GetTotalFindCategory(category string) (int, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var total int64
	err := repo.db.Model(&entity.Download{}).Where("visible = true AND category = ?", category).Count(&total).Error
	return int(total), err
}

func (repo *DownloadPostgresRepository) GetTotalSearch(strSearch string) (int, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var count int64
	err := repo.db.Model(&entity.Download{}).
		Where("visible = true AND unaccent(title) ILIKE unaccent(?)", "%"+strSearch+"%").
		Count(&count).Error
	return int(count), err
}

func (repo *DownloadPostgresRepository) Search(page int, strSearch string) ([]*entity.Download, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	limit := PerPage
	offset := (page - 1) * limit
	var downloads []*entity.Download

	err := repo.db.Where("visible = true AND unaccent(title) ILIKE unaccent(?)", "%"+strSearch+"%").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&downloads).Error

	if err != nil {
		return nil, err
	}

	return downloads, nil
}
