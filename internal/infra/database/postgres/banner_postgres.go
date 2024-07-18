package postgres

import (
	"errors"
	"sync"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"gorm.io/gorm"
)

var (
	ErrBannerNotFound = errors.New("usuário não encontrado")
)

type BannerPostgresRepository struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

func NewBannerPostgresRepository(db_adapter port.DBAdapter) *BannerPostgresRepository {
	db, _ := db_adapter.Connect()
	return &BannerPostgresRepository{db: db}
}

func (repo *BannerPostgresRepository) GetByID(id string) (entity.BannerDTO, error) {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var banner entity.Banner
	repo.db.Model(&entity.Banner{}).Where("id = ?", id).First(&banner)

	if repo.db.Error != nil {
		return entity.BannerDTO{}, ErrBannerNotFound
	}

	tx.Commit()

	dto := entity.BannerDTO{
		ID:     banner.ID,
		Title:  banner.Title,
		Link:   banner.Link,
		Html:   banner.Html,
		Image1: banner.Image1,
		Image2: banner.Image2,
		Image3: banner.Image3,
		Tag:    banner.Tag,
	}

	return dto, nil
}

func (repo *BannerPostgresRepository) GetByTag(tag string) (entity.BannerDTO, error) {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var banner entity.Banner
	repo.db.Model(&entity.Banner{}).Where("tag = ?", tag).First(&banner)

	if repo.db.Error != nil {
		return entity.BannerDTO{}, ErrBannerNotFound
	}

	tx.Commit()

	dto := entity.BannerDTO{
		ID:     banner.ID,
		Title:  banner.Title,
		Link:   banner.Link,
		Html:   banner.Html,
		Image1: banner.Image1,
		Image2: banner.Image2,
		Image3: banner.Image3,
		Tag:    banner.Tag,
	}

	return dto, nil
}

func (repo *BannerPostgresRepository) Create(banner entity.Banner) (entity.BannerDTO, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	result := repo.db.Create(&banner)

	if result.Error != nil {
		tx.Rollback()
		return entity.BannerDTO{}, result.Error
	}

	tx.Commit()

	dto := entity.BannerDTO{
		ID:     banner.ID,
		Title:  banner.Title,
		Link:   banner.Link,
		Html:   banner.Html,
		Image1: banner.Image1,
		Image2: banner.Image2,
		Image3: banner.Image3,
		Tag:    banner.Tag,
	}

	return dto, nil
}
