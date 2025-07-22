package postgres

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/config"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
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
	db := db_adapter.GetDB()
	return &NewsPostgresRepository{db: db}
}

func (repo *NewsPostgresRepository) CleanNews() ([]*entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news []*entity.News

	result := tx.Model(&entity.News{}).Where("visible = false AND created_at <= ?", time.Now().AddDate(0, 0, config.News_DisabledCleaningDays)).Order("created_at DESC").Find(&news)

	if result.Error != nil {
		return []*entity.News{}, result.Error
	}

	return news, nil

}

func (repo *NewsPostgresRepository) CleanNewsOld() ([]*entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news []*entity.News

	result := tx.Model(&entity.News{}).Where("created_at <= ?", time.Now().AddDate(0, 0, config.News_OldCleaningDays)).Order("created_at DESC").Find(&news)

	if result.Error != nil {
		return []*entity.News{}, result.Error
	}

	return news, nil

}

func (repo *NewsPostgresRepository) NewsMake() (*entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news *entity.News

	result := repo.db.Model(&entity.News{}).Where("((visible = true AND category='famosos' AND topstory = false) OR (topstory = true AND visible = true)) AND Make = false AND created_at >= ? AND created_at <= ?", time.Now().AddDate(0, 0, config.News_MakeDays), time.Now()).Order("created_at DESC").Limit(1).First(&news)

	if result.Error != nil {
		return &entity.News{}, result.Error
	}

	result = repo.db.Model(&news).Updates(map[string]interface{}{
		"Make": true,
	})

	if result.Error != nil {
		tx.Rollback()
		return &entity.News{}, result.Error
	}

	return news, nil
}

func (repo *NewsPostgresRepository) Create(news *entity.News) (*entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var n *entity.News
	nresult := repo.db.Model(&entity.News{}).Where("visible = true AND title = ?", news.Title).First(&n)

	if nresult.Error != nil {

		result := repo.db.Create(&news)

		if result.Error != nil {
			tx.Rollback()
			return &entity.News{}, result.Error
		}
	}

	tx.Commit()

	return news, nil

}

func (repo *NewsPostgresRepository) Update(news *entity.News) (*entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	result := repo.db.Model(&news).Updates(map[string]interface{}{
		"title":      news.Title,
		"title_ai":   news.TitleAi,
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
		return &entity.News{}, result.Error
	}

	tx.Commit()

	return news, nil
}

func (repo *NewsPostgresRepository) GetByID(id string) (*entity.News, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var news *entity.News
	result := repo.db.Where("id = ?", id).First(&news)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return news, nil
}

func (repo *NewsPostgresRepository) AdminGetBySlug(slug string) (*entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news *entity.News
	result := repo.db.Model(&entity.News{}).Where("slug = ?", slug).First(&news)

	if result.Error != nil {
		return &entity.News{}, result.Error
	}

	news.Views += 1

	result = repo.db.Save(&news)

	if result.Error != nil {
		return &entity.News{}, result.Error
	}

	tx.Commit()

	return news, nil
}

func (repo *NewsPostgresRepository) GetBySlug(ctx context.Context, slug string) (*entity.News, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.WithContext(ctx).Begin() // inicia transação com contexto
	defer func() {
		// rollback só se a transação ainda não foi committed
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
		}
	}()

	var news entity.News
	result := tx.Where("visible = true AND slug = ?", slug).First(&news)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	// Incrementa a view
	news.Views += 1

	result = tx.Save(&news)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &news, nil
}

func (repo *NewsPostgresRepository) SearchNews(page int, str_search string) ([]*entity.News, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	limit := config.News_PerPage
	offset := (page - 1) * limit

	var news []*entity.News
	repo.db.Model(&entity.News{}).Where("visible = true AND unaccent(title) ILIKE unaccent(?)", "%"+str_search+"%").Order("created_at DESC").Limit(limit).Offset(offset).Find(&news)

	return news, nil
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

func (repo *NewsPostgresRepository) AdminFindAll(page, limit int) ([]*entity.News, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	offset := (page - 1) * limit

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news []*entity.News
	repo.db.Model(&entity.News{}).Order("created_at DESC, visible DESC").Limit(limit).Offset(offset).Find(&news)

	if repo.db.Error != nil {
		return []*entity.News{}, repo.db.Error
	}

	tx.Commit()

	return news, nil
}

func (repo *NewsPostgresRepository) FindAll(ctx context.Context, page, limit int) ([]*entity.News, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	offset := (page - 1) * limit

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news []*entity.News
	repo.db.Model(&entity.News{}).Where("visible = true").Order("created_at DESC").Limit(limit).Offset(offset).Find(&news)

	if repo.db.Error != nil {
		return []*entity.News{}, repo.db.Error
	}

	tx.Commit()

	return news, nil
}

func (repo *NewsPostgresRepository) FindRecent() (*entity.News, error) {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news *entity.News
	tx.Model(&entity.News{}).Where("visible = true").Order("created_at DESC").First(&news)

	if tx.Error != nil {
		return &entity.News{}, tx.Error
	}

	tx.Commit()

	return news, nil
}

func (repo *NewsPostgresRepository) FindCategory(category string, page int) ([]*entity.News, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	limit := config.News_PerPage
	offset := (page - 1) * limit

	tx := repo.db.Begin()
	defer tx.Rollback()

	var news []*entity.News
	tx.Model(&entity.News{}).Where("visible = true AND category=?", category).Order("created_at DESC").Limit(limit).Offset(offset).Find(&news)

	if tx.Error != nil {
		return []*entity.News{}, tx.Error
	}

	tx.Commit()

	return news, nil
}

func (repo *NewsPostgresRepository) FindAllViews() ([]*entity.News, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var news []*entity.News

	result := repo.db.Model(&entity.News{}).Where("visible = true AND created_at >= ? AND created_at <= ?", time.Now().AddDate(0, 0, -2), time.Now()).Order("views DESC").Limit(10).Find(&news)

	if result.Error != nil {
		return []*entity.News{}, result.Error
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

	repo.db.Model(&entity.News{}).Where("id = ?", id).Update("image", "")

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

	var news *entity.News
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

	var news *entity.News
	newsSelected := repo.db.Model(&entity.News{}).Where("id = ?", id).First(&news)

	if newsSelected.Error != nil {
		return ErrNewsNotFound
	}

	repo.db.Unscoped().Delete(news)

	tx.Commit()

	return nil
}

func (repo *NewsPostgresRepository) DeleteAll(news []*entity.News) error {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	for _, b := range news {

		var newsUnic *entity.News
		newsSelected := repo.db.Model(&entity.News{}).Where("id = ?", b.ID).First(&newsUnic)

		if newsSelected.Error != nil {
			return ErrNewsNotFound
		}

		repo.db.Unscoped().Delete(newsUnic)

	}

	tx.Commit()

	return nil
}

func (repo *NewsPostgresRepository) SetVisible(visible bool, id string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	repo.db.Model(&entity.News{}).Where("id = ?", id).Update("visible", visible)

	if repo.db.Error != nil {
		return repo.db.Error
	}

	tx.Commit()

	return nil
}
