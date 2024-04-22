package supabase

import (
	"fmt"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)


type TopPostgresRepository struct {
    db *gorm.DB
}

func NewTopPostgresRepository(db *gorm.DB) *TopPostgresRepository {
	return &TopPostgresRepository{ db: db }
}

// insertNews insere um novo usuário no banco de dados
func (repo *TopPostgresRepository) Create(tops []entity.Top)  error {    
    
    tx := repo.db.Begin()
    defer tx.Rollback()

    fmt.Println("Chegou no Truncate")
	
	result := repo.db.Exec("TRUNCATE TABLE tops")

	if result.Error != nil {
        tx.Rollback()
        return result.Error
    }

    result = repo.db.Create(&tops)
    
    if result.Error != nil {
        tx.Rollback() 
        return result.Error
    }

    fmt.Println("Criou")

    tx.Commit()

    return nil
    
}

func (repo *TopPostgresRepository) FindAll() ([]entity.Top, error) {  

	var tops []entity.Top
    
    result := repo.db.Model(&entity.Top{}).Order("views DESC").Limit(10).Find(&tops)

    if result.Error != nil {
        return []entity.Top{}, result.Error
    }

    return tops, nil

}
