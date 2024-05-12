package supabase

import (
	"sync"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)


type TopPostgresRepository struct {
    db *gorm.DB
    mutex sync.RWMutex
}

func NewTopPostgresRepository(db *gorm.DB) *TopPostgresRepository {
	return &TopPostgresRepository{ db: db }
}

func (repo *TopPostgresRepository) Create(tops []entity.Top)  error {    
    repo.mutex.Lock() 
    defer repo.mutex.Unlock()
    
    tx := repo.db.Begin()
    defer tx.Rollback()    
	
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

    tx.Commit()

    return nil
    
}

func (repo *TopPostgresRepository) FindAll() ([]entity.Top, error) {
    repo.mutex.RLock() 
    defer repo.mutex.RUnlock()

	var tops []entity.Top
    
    result := repo.db.Model(&entity.Top{}).Order("views DESC").Limit(10).Find(&tops)

    if result.Error != nil {
        return []entity.Top{}, result.Error
    }

    return tops, nil

}
