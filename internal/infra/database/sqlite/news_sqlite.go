package sqlite

import (
	"errors"
	"fmt"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	"github.com/eduardospek/bn-api/internal/utils"
	_ "github.com/lib/pq"
)

var (
	ErrNewsExists = errors.New("notícia já cadastrada com este título")	
    ErrNewsNotExistsWithID = errors.New("não existe notícia com este ID")
)

type NewsSQLiteRepository struct {}

func NewNewsSQLiteRepository() *NewsSQLiteRepository {
	return &NewsSQLiteRepository{}
}

// insertNews insere um novo usuário no banco de dados
func (repo *NewsSQLiteRepository) Create(news entity.News) (entity.News, error) {    
    db, err := conn.Connect()

    if err != nil {
        return entity.News{}, err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    if err := db.Create(&news).Error; err != nil {
        tx.Rollback()
        return entity.News{}, err
    }

    tx.Commit()

    return news, nil
    
}

func (repo *NewsSQLiteRepository) Update(news entity.News) (entity.News, error)  {    
    db, err := conn.Connect()

    if err != nil {
        return entity.News{}, err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    if err := db.Updates(&news).Error; err != nil {
        tx.Rollback()
        return entity.News{}, err
    }

    tx.Commit()

    updatenews, err := repo.GetById(news.ID)

    if err != nil {
        fmt.Println(err)
		return entity.News{}, err
	}

    return updatenews, err
}

func (repo *NewsSQLiteRepository) GetById(id string) (entity.News, error) {
	db, err := conn.Connect()

    if err != nil {
        return entity.News{}, err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    var news entity.News
    db.Where("id = ?", id).First(news)

    tx.Commit()

    return news, nil
}

func (repo *NewsSQLiteRepository) GetBySlug(slug string) (entity.News, error) {
	db, err := conn.Connect()

    if err != nil {
        return entity.News{}, err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    var news entity.News
    db.Where("slug = ?", slug).First(news)

    tx.Commit()

    return news, nil
}

func (repo *NewsSQLiteRepository) FindAll(page, limit int) (interface{}, error) {
	
	offset := (page - 1) * limit

    db, err := conn.Connect()

    if err != nil {
        return entity.News{}, err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    var news []entity.News
    db.Model(&entity.News{}).Where("visible = true").Limit(limit).Offset(offset).Find(&news)

    tx.Commit()

    var total int64
    db.Model(&entity.News{}).Count(&total)

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

func (repo *NewsSQLiteRepository) Delete(id string) (error) {
	
	db, err := conn.Connect()

    if err != nil {
        return err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    var news entity.News
    
    // Utilize o método `Delete` e passe o ID do registro como argumento
    err = db.Model(&news).Where("id = ?", id).Delete(&news).Error
    if err != nil {
        return err
    }
    
    tx.Commit()

    return nil

}

func (repo *NewsSQLiteRepository) NewsTruncateTable() error {
    db, err := conn.Connect()

    if err != nil {
        return err
    }

	tx := db.Begin()
    defer tx.Rollback() 

    err = db.Exec("TRUNCATE TABLE users").Error
    
    if err != nil {
        return err
    }

    tx.Commit()

    return nil
}

//VALIDATIONS
func (repo *NewsSQLiteRepository) NewsExists(title string) error {
    db, err := conn.Connect()

    if err != nil {
        return err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    var news entity.News
    err = db.Where("title = ?", title).First(news).Error

    if err != nil {
        return nil
    }

    tx.Commit()    
  
    return ErrNewsExists
}