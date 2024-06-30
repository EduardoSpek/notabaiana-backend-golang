package postgres

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	ErrNewsExists          = errors.New("notícia já cadastrada com este título")
	ErrNewsNotExistsWithID = errors.New("não existe notícia com este ID")
)

type NewsPostgresRepository struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

func NewNewsPostgresRepository(db_adapter port.DBAdapter) *NewsPostgresRepository {
	db, _ := db_adapter.Connect()
	return &NewsPostgresRepository{db: db}
}

func (repo *NewsPostgresRepository) NewsMake() (entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news entity.News

	result := repo.db.Model(&entity.News{}).Where("visible = true AND Make = false AND category = 'famosos' AND created_at >= ? AND created_at <= ?", time.Now().AddDate(0, 0, -2), time.Now()).Order("created_at DESC").Limit(1).First(&news)

	if result.Error != nil {
		return entity.News{}, result.Error
	}

	result = repo.db.Model(&news).Updates(map[string]interface{}{
		"Make": true,
	})

	if result.Error != nil {
		tx.Rollback()
		return entity.News{}, result.Error
	}

	return news, nil
}

// insertNews insere um novo usuário no banco de dados
func (repo *NewsPostgresRepository) Create(news entity.News) (entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	result := repo.db.Create(&news)

	if result.Error != nil {
		tx.Rollback()
		return entity.News{}, result.Error
	}

	tx.Commit()

	return news, nil

}

func (repo *NewsPostgresRepository) Update(news entity.News) (entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	result := repo.db.Model(&news).Updates(map[string]interface{}{
		"title":      news.Title,
		"text":       news.Text,
		"visible":    news.Visible,
		"category":   news.Category,
		"slug":       news.Slug,
		"link":       news.Link,
		"image":      news.Image,
		"updated_at": news.UpdatedAt,
	})

	if result.Error != nil {
		tx.Rollback()
		return entity.News{}, result.Error
	}

	tx.Commit()

	return news, nil
}

func (repo *NewsPostgresRepository) GetById(id string) (entity.News, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news entity.News
	repo.db.Model(&entity.News{}).Where("id = ?", id).First(&news)

	if repo.db.Error != nil {
		return entity.News{}, repo.db.Error
	}

	tx.Commit()

	return news, nil
}

func (repo *NewsPostgresRepository) GetBySlug(slug string) (entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news entity.News
	result := repo.db.Model(&entity.News{}).Where("visible = true AND slug = ?", slug).First(&news)

	if result.Error != nil {
		return entity.News{}, result.Error
	}

	news.Views += 1

	result = repo.db.Save(&news)

	if result.Error != nil {
		return entity.News{}, result.Error
	}

	tx.Commit()

	return news, nil
}

func (repo *NewsPostgresRepository) SearchNews(page int, str_search string) []entity.News {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	limit := 16
	offset := (page - 1) * limit

	var news []entity.News
	repo.db.Model(&entity.News{}).Where("visible = true AND unaccent(title) ILIKE unaccent(?)", "%"+str_search+"%").Order("created_at DESC").Limit(limit).Offset(offset).Find(&news)

	return news
}

func (repo *NewsPostgresRepository) GetTotalNewsBySearch(str_search string) int {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var count int64
	repo.db.Model(&entity.News{}).Where("visible = true AND unaccent(title) ILIKE unaccent(?)", "%"+str_search+"%").Count(&count)

	total := int(count)

	return total

}

func (repo *NewsPostgresRepository) GetTotalNewsByCategory(category string) int {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var total int64
	repo.db.Model(&entity.News{}).Where("visible = true AND category=?", category).Count(&total)

	return int(total)

}

func (repo *NewsPostgresRepository) GetTotalNews() int {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var total int64
	repo.db.Model(&entity.News{}).Where("visible = true").Count(&total)

	return int(total)

}

func (repo *NewsPostgresRepository) FindAll(page, limit int) ([]entity.News, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	offset := (page - 1) * limit

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news []entity.News
	repo.db.Model(&entity.News{}).Where("visible = true").Order("created_at DESC").Limit(limit).Offset(offset).Find(&news)

	if repo.db.Error != nil {
		return []entity.News{}, repo.db.Error
	}

	tx.Commit()

	return news, nil
}

func (repo *NewsPostgresRepository) FindRecent() (entity.News, error) {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news entity.News
	repo.db.Model(&entity.News{}).Where("visible = true").Order("created_at DESC").First(&news)

	if repo.db.Error != nil {
		return entity.News{}, repo.db.Error
	}

	tx.Commit()

	return news, nil
}

func (repo *NewsPostgresRepository) FindCategory(category string, page int) ([]entity.News, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	limit := 16
	offset := (page - 1) * limit

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news []entity.News
	repo.db.Model(&entity.News{}).Where("visible = true AND category=?", category).Order("created_at DESC").Limit(limit).Offset(offset).Find(&news)

	if repo.db.Error != nil {
		return []entity.News{}, repo.db.Error
	}

	tx.Commit()

	return news, nil
}

func (repo *NewsPostgresRepository) FindAllViews() ([]entity.News, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var news []entity.News

	result := repo.db.Model(&entity.News{}).Where("visible = true AND created_at >= ? AND created_at <= ?", time.Now().AddDate(0, 0, -2), time.Now()).Order("views DESC").Limit(10).Find(&news)

	if result.Error != nil {
		return []entity.News{}, result.Error
	}

	return news, nil
}

func (repo *NewsPostgresRepository) ClearViews() error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	result := repo.db.Exec("UPDATE news SET views = 0")

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *NewsPostgresRepository) Delete(id string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news entity.News

	// Utilize o método `Delete` e passe o ID do registro como argumento
	repo.db.Model(&news).Where("id = ?", id).Delete(&news)

	if repo.db.Error != nil {
		return nil
	}

	tx.Commit()

	return nil

}

func (repo *NewsPostgresRepository) ClearImagePath(id string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news entity.News
	news.ID = id

	repo.db.Model(&news).Update("image", "")

	if repo.db.Error != nil {
		return repo.db.Error
	}

	tx.Commit()

	return nil
}

func (repo *NewsPostgresRepository) NewsTruncateTable() error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	repo.db.Exec("TRUNCATE TABLE news")

	if repo.db.Error != nil {
		return nil
	}

	tx.Commit()

	return nil
}

// VALIDATIONS
func (repo *NewsPostgresRepository) NewsExists(title string) error {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	title = strings.ToLower(title)

	var news entity.News
	result := repo.db.Model(&entity.News{}).Where("LOWER(title) = ?", title).First(&news)

	if result.Error != nil {
		return nil
	}

	return ErrNewsExists
}
