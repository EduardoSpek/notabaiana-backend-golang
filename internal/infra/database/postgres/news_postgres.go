package postgres

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	ErrNewsExists          = errors.New("notícia já cadastrada com este título")
	ErrNewsNotExistsWithID = errors.New("não existe notícia com este ID")
	ErrNewsNotFound        = errors.New("notícia não encontrada")
)

type NewsPostgresRepository struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

func NewNewsPostgresRepository(db_adapter port.DBAdapter) *NewsPostgresRepository {
	db, _ := db_adapter.Connect()
	return &NewsPostgresRepository{db: db}
}

func (repo *NewsPostgresRepository) CleanNews() {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news []entity.News

	repo.db.Model(&entity.News{}).Where("visible = false AND created_at <= ?", time.Now().AddDate(0, 0, -7)).Order("created_at DESC").Find(&news)

	for _, n := range news {

		if n.Image != "" {
			image := "./images/" + n.Image
			utils.RemoveImage(image)
		}

		repo.db.Unscoped().Delete(n)
	}
}

func (repo *NewsPostgresRepository) NewsMake() (entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news entity.News

	result := repo.db.Model(&entity.News{}).Where("((visible = true AND category='famosos' AND topstory = false) OR (topstory = true AND visible = true)) AND Make = false AND created_at >= ? AND created_at <= ?", time.Now().AddDate(0, 0, -2), time.Now()).Order("created_at DESC").Limit(1).First(&news)

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

func (repo *NewsPostgresRepository) Create(news entity.News) (entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var n entity.News
	nresult := repo.db.Model(&entity.News{}).Where("visible = true AND title = ?", news.Title).First(&n)

	if nresult.Error != nil {

		result := repo.db.Create(&news)

		if result.Error != nil {
			tx.Rollback()
			return entity.News{}, result.Error
		}
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
		"topstory":   news.TopStory,
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

func (repo *NewsPostgresRepository) AdminGetBySlug(slug string) (entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news entity.News
	result := repo.db.Model(&entity.News{}).Where("slug = ?", slug).First(&news)

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
	repo.db.Model(&entity.News{}).Count(&total)

	return int(total)

}

func (repo *NewsPostgresRepository) GetTotalNewsVisible() int {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var total int64
	repo.db.Model(&entity.News{}).Where("visible = true").Count(&total)

	return int(total)

}

func (repo *NewsPostgresRepository) AdminFindAll(page, limit int) ([]entity.News, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	offset := (page - 1) * limit

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news []entity.News
	repo.db.Model(&entity.News{}).Order("created_at DESC").Limit(limit).Offset(offset).Find(&news)

	if repo.db.Error != nil {
		return []entity.News{}, repo.db.Error
	}

	tx.Commit()

	return news, nil
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

func (repo *NewsPostgresRepository) Delete(id string) error {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news entity.News
	newsSelected := repo.db.Model(&entity.News{}).Where("id = ?", id).First(&news)

	if newsSelected.Error != nil {
		return ErrNewsNotFound
	}

	del1 := utils.RemoveImage("." + news.Image)

	if !del1 {
		fmt.Println("Imagem não deletada")
	}
	repo.db.Unscoped().Delete(news)

	tx.Commit()

	return nil
}

func (repo *NewsPostgresRepository) DeleteAll(news []entity.News) error {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	for _, b := range news {

		var news entity.News
		newsSelected := repo.db.Model(&entity.News{}).Where("id = ?", b.ID).First(&news)

		if newsSelected.Error != nil {
			return ErrNewsNotFound
		}

		del1 := utils.RemoveImage("." + news.Image)

		if !del1 {
			fmt.Println("Imagem da notícia não deletada")
		}

		repo.db.Unscoped().Delete(news)

	}

	tx.Commit()

	return nil
}
