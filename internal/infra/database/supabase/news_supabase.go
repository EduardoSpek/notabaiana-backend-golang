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

type NewsSupabaseRepository struct {}

func NewNewsSupabaseRepository() *NewsSupabaseRepository {
	return &NewsSupabaseRepository{}
}

// insertNews insere um novo usuário no banco de dados
func (repo *NewsSupabaseRepository) Create(news entity.News) (entity.News, error) {    
    db, err := conn.Connect()

    if err != nil {
        return entity.News{}, err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    db.Create(&news)
    
    if db.Error != nil {
        tx.Rollback() 
        return entity.News{}, db.Error
    }

    tx.Commit()

    return news, nil
    
}

func (repo *NewsSupabaseRepository) Update(news entity.News) (entity.News, error)  {    
    db, err := conn.Connect()

    if err != nil {
        return entity.News{}, err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    db.Updates(&news)

    if db.Error != nil {
        tx.Rollback() 
        return entity.News{}, db.Error
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
	db, err := conn.Connect()

    if err != nil {
        return entity.News{}, err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    var news entity.News
    db.Where("id = ?", id).First(news)

    if db.Error != nil {
        return entity.News{}, db.Error
    }

    tx.Commit()

    return news, nil
}

func (repo *NewsSupabaseRepository) GetBySlug(slug string) (entity.News, error) {
	db, err := conn.Connect()

    if err != nil {
        return entity.News{}, err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    var news entity.News
    db.Where("slug = ?", slug).First(news)
    
    if db.Error != nil {
        return entity.News{}, db.Error
    }

    tx.Commit()

    return news, nil
}

func (repo *NewsSupabaseRepository) FindAll(page, limit int) (interface{}, error) {
	
	offset := (page - 1) * limit

    db, err := conn.Connect()

    if err != nil {
        return entity.News{}, err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    var news []entity.News
    db.Model(&entity.News{}).Where("visible = true").Limit(limit).Offset(offset).Find(&news)

    if db.Error != nil {
        return entity.News{}, db.Error
    }

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

func (repo *NewsSupabaseRepository) Delete(id string) (error) {
	
	db, err := conn.Connect()

    if err != nil {
        return err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    var news entity.News
    
    // Utilize o método `Delete` e passe o ID do registro como argumento
    db.Model(&news).Where("id = ?", id).Delete(&news)
    
    if db.Error != nil {
        return nil
    }
    
    tx.Commit()

    return nil

}

func (repo *NewsSupabaseRepository) NewsTruncateTable() error {
    db, err := conn.Connect()

    if err != nil {
        return err
    }

	tx := db.Begin()
    defer tx.Rollback() 

    db.Exec("TRUNCATE TABLE users")
    
    if db.Error != nil {
        return nil
    }

    tx.Commit()

    return nil
}

//VALIDATIONS
func (repo *NewsSupabaseRepository) NewsExists(title string) error {
    db, err := conn.Connect()

    if err != nil {
        return err
    }
	
    tx := db.Begin()
    defer tx.Rollback()    

    var news entity.News
    db.Model(&entity.News{}).Where("title = ?", title).First(&news)

    if db.Error != nil {
        return nil
    }

    tx.Commit()    
  
    return ErrNewsExists
}