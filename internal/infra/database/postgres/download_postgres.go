package postgres

import (
	"errors"
	"sync"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/config"
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
	db := db_adapter.GetDB()
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

	result := tx.Model(download).Updates(download)

	if result.Error != nil {
		return nil, result.Error
	}

	result = tx.Model(download).Updates(map[string]interface{}{
		"visible": download.Visible,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return download, nil
}

func (repo *DownloadPostgresRepository) GetByID(id string) (*entity.Download, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var download entity.Download
	result := repo.db.Where("id = ?", id).First(&download)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &download, nil
}

func (repo *DownloadPostgresRepository) GetByLink(link string) (*entity.Download, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var download entity.Download
	result := repo.db.Model(&entity.Download{}).Where("link = ?", link).First(&download)

	if result.Error != nil {
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
	if err := tx.Where("visible = true AND slug = ?", slug).Preload("Musics").First(&download).Error; err != nil {
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

func (repo *DownloadPostgresRepository) AdminGetBySlug(slug string) (*entity.Download, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer tx.Rollback()

	var download entity.Download
	if err := tx.Where("slug = ?", slug).First(&download).Error; err != nil {
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

func (repo *DownloadPostgresRepository) AdminFindAll(page, limit int) ([]*entity.Download, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	offset := (page - 1) * limit
	var downloads []*entity.Download

	err := repo.db.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&downloads).Error

	if err != nil {
		return nil, err
	}

	return downloads, nil
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
		Preload("Musics").
		Find(&downloads).Error

	if err != nil {
		return nil, err
	}

	return downloads, nil
}

func (repo *DownloadPostgresRepository) FindAllTopViews(page, limit int) ([]*entity.Download, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	offset := (page - 1) * limit
	var downloads []*entity.Download

	err := repo.db.Where("visible = true").
		Order("views DESC").
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

	limit := config.Downloads_PerPage
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

func (repo *DownloadPostgresRepository) GetTotal() (int, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var total int64
	err := repo.db.Model(&entity.Download{}).Count(&total).Error
	return int(total), err
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

	limit := config.Downloads_PerPage
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

func (repo *DownloadPostgresRepository) Delete(id string) error {
	repo.mutex.Lock() // Usamos Lock em vez de RLock porque estamos modificando dados
	defer repo.mutex.Unlock()

	return repo.db.Transaction(func(tx *gorm.DB) error {
		var download entity.Download

		if err := tx.First(&download, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrDownloadNotFound
			}
			return err
		}

		if err := tx.Unscoped().Delete(&download).Error; err != nil {
			return err
		}

		return nil
	})
}

func (repo *DownloadPostgresRepository) DeleteAll(downloads []*entity.Download) error {
	repo.mutex.Lock() // Usamos Lock em vez de RLock porque estamos modificando dados
	defer repo.mutex.Unlock()

	return repo.db.Transaction(func(tx *gorm.DB) error {
		for _, download := range downloads {
			if err := tx.First(&download, "id = ?", download.ID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrDownloadNotFound
				}
				return err
			}

			if err := tx.Unscoped().Delete(download).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *DownloadPostgresRepository) Clean() ([]*entity.Download, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var downloads []*entity.Download

	result := repo.db.Model(&entity.Download{}).Where("visible = false AND created_at <= ?", time.Now().AddDate(0, 0, -7)).Order("created_at DESC").Find(&downloads)

	if result.Error != nil {
		return []*entity.Download{}, result.Error
	}

	return downloads, nil

}
