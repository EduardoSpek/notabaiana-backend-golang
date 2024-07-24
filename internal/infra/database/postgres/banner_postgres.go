package postgres

import (
	"errors"
	"fmt"
	"sync"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"gorm.io/gorm"
)

var (
	ErrBannerNotFound = errors.New("banner não encontrado")
)

type BannerPostgresRepository struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

func NewBannerPostgresRepository(db_adapter port.DBAdapter) *BannerPostgresRepository {
	db, _ := db_adapter.Connect()
	return &BannerPostgresRepository{db: db}
}

func (repo *BannerPostgresRepository) AdminFindAll() ([]entity.BannerDTO, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var banners []entity.BannerDTO
	list := repo.db.Model(&entity.Banner{}).Order("created_at DESC").Find(&banners)

	if list.Error != nil {
		return []entity.BannerDTO{}, list.Error
	}

	tx.Commit()

	return banners, nil
}

func (repo *BannerPostgresRepository) FindAll() ([]entity.BannerDTO, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var banners []entity.BannerDTO
	list := repo.db.Model(&entity.Banner{}).Where("visible = true").Order("RANDOM()").Find(&banners)

	if list.Error != nil {
		return []entity.BannerDTO{}, list.Error
	}

	tx.Commit()

	return banners, nil
}

func (repo *BannerPostgresRepository) GetByID(id string) (entity.BannerDTO, error) {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var banner entity.Banner
	bannerSelected := repo.db.Model(&entity.Banner{}).Where("id = ?", id).First(&banner)

	if bannerSelected.Error != nil {
		return entity.BannerDTO{}, ErrBannerNotFound
	}

	tx.Commit()

	dto := entity.BannerDTO{
		ID:      banner.ID,
		Title:   banner.Title,
		Link:    banner.Link,
		Html:    banner.Html,
		Image1:  banner.Image1,
		Image2:  banner.Image2,
		Image3:  banner.Image3,
		Tag:     banner.Tag,
		Visible: banner.Visible,
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

func (repo *BannerPostgresRepository) Update(banner entity.Banner) (entity.BannerDTO, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	result := repo.db.Model(&banner).Updates(map[string]interface{}{
		"title":      banner.Title,
		"link":       banner.Link,
		"html":       banner.Html,
		"image1":     banner.Image1,
		"image2":     banner.Image2,
		"image3":     banner.Image3,
		"tag":        banner.Tag,
		"visible":    banner.Visible,
		"updated_at": banner.UpdatedAt,
	})

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

func (repo *BannerPostgresRepository) Delete(id string) error {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var banner entity.Banner
	bannerSelected := repo.db.Model(&entity.Banner{}).Where("id = ?", id).First(&banner)

	if bannerSelected.Error != nil {
		return ErrBannerNotFound
	}

	del1 := utils.RemoveImage("." + banner.Image1)

	if !del1 {
		fmt.Println("Imagem 1 não deletada")
	}

	del2 := utils.RemoveImage("." + banner.Image2)

	if !del2 {
		fmt.Println("Imagem 2 não deletada")
	}

	del3 := utils.RemoveImage("." + banner.Image3)

	if !del3 {
		fmt.Println("Imagem 3 não deletada")
	}
	repo.db.Unscoped().Delete(banner)

	tx.Commit()

	return nil
}

func (repo *BannerPostgresRepository) DeleteAll(banners []entity.BannerDTO) error {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	for _, b := range banners {

		var banner entity.Banner
		bannerSelected := repo.db.Model(&entity.Banner{}).Where("id = ?", b.ID).First(&banner)

		if bannerSelected.Error != nil {
			return ErrBannerNotFound
		}

		del1 := utils.RemoveImage("." + banner.Image1)

		if !del1 {
			fmt.Println("Imagem 1 não deletada")
		}

		del2 := utils.RemoveImage("." + banner.Image2)

		if !del2 {
			fmt.Println("Imagem 2 não deletada")
		}

		del3 := utils.RemoveImage("." + banner.Image3)

		if !del3 {
			fmt.Println("Imagem 3 não deletada")
		}

		repo.db.Unscoped().Delete(banner)

	}

	tx.Commit()

	return nil
}
