package supabase

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	"github.com/eduardospek/bn-api/internal/utils"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	ErrNewsExists = errors.New("notícia já cadastrada com este título")	
    ErrNewsNotExistsWithID = errors.New("não existe notícia com este ID")
)

type NewsPostgresRepository struct {
    db *gorm.DB
}

func NewNewsPostgresRepository(db *gorm.DB) *NewsPostgresRepository {
	return &NewsPostgresRepository{ db: db }
}

// insertNews insere um novo usuário no banco de dados
func (repo *NewsPostgresRepository) Create(news entity.News) (entity.News, error) {    
    
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

func (repo *NewsPostgresRepository) Update(news entity.News) (entity.News, error)  {    
	
    tx := repo.db.Begin()
    defer tx.Rollback()    

    result := repo.db.Updates(&news)
    
    if result.Error != nil {
        tx.Rollback() 
        return entity.News{}, result.Error
    }

    tx.Commit()

    updatenews, err := repo.GetById(news.ID)

    if err != nil {
        fmt.Println(err)
		return entity.News{}, err
	}

    return updatenews, err
}

func (repo *NewsPostgresRepository) GetById(id string) (entity.News, error) {	

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

func (repo *NewsPostgresRepository) SearchNews(page int, str_search string) interface{} {

    limit := 10
    offset := (page - 1) * limit

    var count int64
    repo.db.Model(&entity.News{}).Where("visible = true AND LOWER(title) LIKE ?", "%"+strings.ToLower(str_search)+"%").Count(&count)

	var news []entity.News
    repo.db.Model(&entity.News{}).Where("visible = true AND LOWER(title) LIKE ?", "%"+strings.ToLower(str_search)+"%").Order("created_at DESC").Limit(limit).Offset(offset).Find(&news)

	total := int(count)  

    pagination := utils.Pagination(page, total)

    result := struct{
        List_news []entity.News `json:"news"`
        Pagination map[string][]int `json:"pagination"`
        Search string `json:"search"`
    }{
        List_news: news,
        Pagination: pagination,
        Search: str_search,
    }

	return result
}

func (repo *NewsPostgresRepository) FindAll(page, limit int) (interface{}, error) {
	
	offset := (page - 1) * limit

    tx := repo.db.Begin()
    defer tx.Rollback()    

    var news []entity.News
    repo.db.Model(&entity.News{}).Where("visible = true").Order("created_at DESC").Limit(limit).Offset(offset).Find(&news)

    if repo.db.Error != nil {
        return entity.News{}, repo.db.Error
    }

    tx.Commit()

    var total int64
    repo.db.Model(&entity.News{}).Count(&total)

    pagination := utils.Pagination(page, int(total))

    result := struct{
        List_news []entity.News `json:"news"`
        Pagination map[string][]int `json:"pagination"`
    }{
        List_news: news,
        Pagination: pagination,
    }
    
    return result, nil
}

func (repo *NewsPostgresRepository) FindCategory(category string, page int) (interface{}, error) {
	
    limit := 10
    offset := (page - 1) * limit

    tx := repo.db.Begin()
    defer tx.Rollback()    

    var news []entity.News
    repo.db.Model(&entity.News{}).Where("visible = true AND category=?", category).Order("created_at DESC").Limit(limit).Offset(offset).Find(&news)

    if repo.db.Error != nil {
        return entity.News{}, repo.db.Error
    }

    tx.Commit()

    var total int64
    repo.db.Model(&entity.News{}).Count(&total)

    pagination := utils.Pagination(page, int(total))

    result := struct{
        List_news []entity.News `json:"news"`
        Pagination map[string][]int `json:"pagination"`
		Category string `json:"category"`
    }{
        List_news: news,
        Pagination: pagination,
		Category: category,
    }

	return result, nil
}

func (repo *NewsPostgresRepository) FindAllViews() ([]entity.News, error) {	

    var news []entity.News
    
    result := repo.db.Model(&entity.News{}).Where("visible = true AND created_at >= ? AND created_at <= ?", time.Now().AddDate(0, 0, -1), time.Now()).Order("created_at DESC, views DESC").Limit(10).Find(&news)

    if result.Error != nil {
        return []entity.News{}, result.Error
    }

    return news, nil
}

func (repo *NewsPostgresRepository) ClearViews() error {	
    
    result := repo.db.Exec("UPDATE news SET views = 0")

    if result.Error != nil {
        return result.Error
    }

    return nil
}

func (repo *NewsPostgresRepository) Delete(id string) (error) {

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

func (repo *NewsPostgresRepository) NewsTruncateTable() error {
    
	tx := repo.db.Begin()
    defer tx.Rollback() 

    repo.db.Exec("TRUNCATE TABLE users")
    
    if repo.db.Error != nil {
        return nil
    }

    tx.Commit()

    return nil
}

//VALIDATIONS
func (repo *NewsPostgresRepository) NewsExists(title string) error {
	
    var news entity.News
    result := repo.db.Model(&entity.News{}).Where("title = ?", title).First(&news)

    if result.Error != nil {
        return nil
    }
  
    return ErrNewsExists
}