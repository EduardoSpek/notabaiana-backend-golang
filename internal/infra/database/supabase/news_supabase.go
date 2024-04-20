package sqlite

import (
	"errors"
	"fmt"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	"github.com/eduardospek/bn-api/internal/utils"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	ErrNewsExists = errors.New("notícia já cadastrada com este título")	
    ErrNewsNotExistsWithID = errors.New("não existe notícia com este ID")
)

type NewsSupabaseRepository struct {
    db *gorm.DB
}

func NewNewsSupabaseRepository(db *gorm.DB) *NewsSupabaseRepository {
	return &NewsSupabaseRepository{ db: db }
}

// insertNews insere um novo usuário no banco de dados
func (repo *NewsSupabaseRepository) Create(news entity.News) (entity.News, error) {    
    
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

func (repo *NewsSupabaseRepository) Update(news entity.News) (entity.News, error)  {    
	
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

func (repo *NewsSupabaseRepository) GetById(id string) (entity.News, error) {	

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

func (repo *NewsSupabaseRepository) GetBySlug(slug string) (entity.News, error) {
	
	
    tx := repo.db.Begin()
    defer tx.Rollback()    

    var news entity.News
    repo.db.Model(&entity.News{}).Where("slug = ?", slug).First(&news)
    
    if repo.db.Error != nil {
        return entity.News{}, repo.db.Error
    }

    tx.Commit()

    return news, nil
}

func (repo *NewsSupabaseRepository) FindAll(page, limit int) (interface{}, error) {
	
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

func (repo *NewsSupabaseRepository) Delete(id string) (error) {

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

func (repo *NewsSupabaseRepository) NewsTruncateTable() error {
    
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
func (repo *NewsSupabaseRepository) NewsExists(title string) error {
	
    tx := repo.db.Begin()
    defer tx.Rollback()    

    var news entity.News
    repo.db.Model(&entity.News{}).Where("title = ?", title).First(&news)

    if repo.db.Error != nil {
        return nil
    }

    tx.Commit()    
  
    return ErrNewsExists
}