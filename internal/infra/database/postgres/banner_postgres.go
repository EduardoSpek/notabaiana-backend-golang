package postgres

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrBannerNotFound    = errors.New("banner não encontrado")
	ErrTransactionFailed = errors.New("falha na transação")
	ErrDatabaseOperation = errors.New("falha na operação do banco de dados")
)

// BannerPostgresRepository implementa as operações de banco de dados para Banner
type BannerPostgresRepository struct {
	db     *gorm.DB
	mutex  sync.RWMutex
	logger *zap.Logger
}

// NewBannerPostgresRepository cria uma nova instância do repositório
func NewBannerPostgresRepository(db_adapter port.DBAdapter, logger *zap.Logger) (*BannerPostgresRepository, error) {
	db := db_adapter.GetDB()
	return &BannerPostgresRepository{
		db:     db,
		logger: logger,
	}, nil
}

// beginTx inicia uma nova transação com tratamento de erro apropriado
func (repo *BannerPostgresRepository) beginTx() (*gorm.DB, error) {
	tx := repo.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrTransactionFailed, tx.Error)
	}
	return tx, nil
}

// deleteImages remove as imagens associadas a um banner
func (repo *BannerPostgresRepository) deleteImages(banner entity.Banner) {
	images := map[string]string{
		"Image1": banner.Image1,
		"Image2": banner.Image2,
		"Image3": banner.Image3,
	}

	for name, path := range images {
		if path != "" {
			if !utils.RemoveImage("." + path) {
				repo.logger.Warn("falha ao deletar imagem",
					zap.String("image", name),
					zap.String("path", path))
			}
		}
	}
}

// AdminFindAll retorna todos os banners para administração
func (repo *BannerPostgresRepository) AdminFindAll(ctx context.Context) ([]entity.BannerDTO, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx, err := repo.beginTx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var banners []entity.Banner
	if err := tx.WithContext(ctx).
		Model(&entity.Banner{}).
		Order("created_at DESC").
		Find(&banners).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	dtos := make([]entity.BannerDTO, len(banners))
	for i, banner := range banners {
		dtos[i] = banner.ToDTO()
	}

	return dtos, nil
}

// FindAll retorna todos os banners visíveis em ordem aleatória
func (repo *BannerPostgresRepository) FindAll(ctx context.Context) ([]entity.BannerDTO, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx, err := repo.beginTx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var banners []entity.Banner
	if err := tx.WithContext(ctx).
		Model(&entity.Banner{}).
		Where("visible = true").
		Order("RANDOM()").
		Find(&banners).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	dtos := make([]entity.BannerDTO, len(banners))
	for i, banner := range banners {
		dtos[i] = banner.ToDTO()
	}

	return dtos, nil
}

// GetByID retorna um banner específico por ID
func (repo *BannerPostgresRepository) GetByID(ctx context.Context, id string) (entity.BannerDTO, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx, err := repo.beginTx()
	if err != nil {
		return entity.BannerDTO{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var banner entity.Banner
	if err := tx.WithContext(ctx).
		Model(&entity.Banner{}).
		Where("id = ?", id).
		First(&banner).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.BannerDTO{}, ErrBannerNotFound
		}
		return entity.BannerDTO{}, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	if err := tx.Commit().Error; err != nil {
		return entity.BannerDTO{}, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	return banner.ToDTO(), nil
}

// GetByTag retorna um banner específico por tag
func (repo *BannerPostgresRepository) GetByTag(ctx context.Context, tag string) (entity.BannerDTO, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx, err := repo.beginTx()
	if err != nil {
		return entity.BannerDTO{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var banner entity.Banner
	if err := tx.WithContext(ctx).
		Model(&entity.Banner{}).
		Where("tag = ?", tag).
		First(&banner).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.BannerDTO{}, ErrBannerNotFound
		}
		return entity.BannerDTO{}, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	if err := tx.Commit().Error; err != nil {
		return entity.BannerDTO{}, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	return banner.ToDTO(), nil
}

// Create cria um novo banner
func (repo *BannerPostgresRepository) Create(ctx context.Context, banner entity.Banner) (entity.BannerDTO, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx, err := repo.beginTx()
	if err != nil {
		return entity.BannerDTO{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.WithContext(ctx).Create(&banner).Error; err != nil {
		tx.Rollback()
		return entity.BannerDTO{}, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	if err := tx.Commit().Error; err != nil {
		return entity.BannerDTO{}, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	return banner.ToDTO(), nil
}

// Update atualiza um banner existente
func (repo *BannerPostgresRepository) Update(ctx context.Context, banner entity.Banner) (entity.BannerDTO, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx, err := repo.beginTx()
	if err != nil {
		return entity.BannerDTO{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updates := map[string]interface{}{
		"title":          banner.Title,
		"link":           banner.Link,
		"html":           banner.Html,
		"image1":         banner.Image1,
		"image2":         banner.Image2,
		"image3":         banner.Image3,
		"tag":            banner.Tag,
		"visible":        banner.Visible,
		"visible_image1": banner.VisibleImage1,
		"visible_image2": banner.VisibleImage2,
		"visible_image3": banner.VisibleImage3,
		"updated_at":     banner.UpdatedAt,
	}

	if err := tx.WithContext(ctx).
		Model(&banner).
		Updates(updates).Error; err != nil {
		tx.Rollback()
		return entity.BannerDTO{}, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	if err := tx.Commit().Error; err != nil {
		return entity.BannerDTO{}, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	return banner.ToDTO(), nil
}

// Delete remove um banner específico
func (repo *BannerPostgresRepository) Delete(ctx context.Context, id string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx, err := repo.beginTx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var banner entity.Banner
	if err := tx.WithContext(ctx).
		Model(&entity.Banner{}).
		Where("id = ?", id).
		First(&banner).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBannerNotFound
		}
		return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	repo.deleteImages(banner)

	if err := tx.WithContext(ctx).Unscoped().Delete(&banner).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	return tx.Commit().Error
}

// DeleteAll remove múltiplos banners
func (repo *BannerPostgresRepository) DeleteAll(ctx context.Context, banners []entity.BannerDTO) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx, err := repo.beginTx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, bannerDTO := range banners {
		var banner entity.Banner
		if err := tx.WithContext(ctx).
			Model(&entity.Banner{}).
			Where("id = ?", bannerDTO.ID).
			First(&banner).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBannerNotFound
			}
			return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
		}

		repo.deleteImages(banner)

		if err := tx.WithContext(ctx).Unscoped().Delete(&banner).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
		}
	}

	return tx.Commit().Error
}
