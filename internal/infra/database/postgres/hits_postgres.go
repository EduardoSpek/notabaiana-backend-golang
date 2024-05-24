package postgres

import (
	"sync"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"gorm.io/gorm"
)

type HitsPostgresRepository struct {
	db *gorm.DB
    mutex sync.RWMutex
}

func NewHitsPostgresRepository(db_adapter port.DBAdapter) *HitsPostgresRepository {
	db, _ := db_adapter.Connect()
	return &HitsPostgresRepository{ db: db }
}

func (repo *HitsPostgresRepository) Save(hit entity.Hits) error {
	repo.mutex.Lock() 
    defer repo.mutex.Unlock()
    
    tx := repo.db.Begin()
    defer tx.Rollback()    

    result := repo.db.Create(&hit)
    
    if result.Error != nil {
        tx.Rollback() 
        return result.Error
    }

    tx.Commit()

    return nil
}

func (repo *HitsPostgresRepository) Update(hit entity.Hits) error {
	repo.mutex.Lock() 
    defer repo.mutex.Unlock()
    
    tx := repo.db.Begin()
    defer tx.Rollback()    

    result := repo.db.Updates(&hit)
    
    if result.Error != nil {
        tx.Rollback() 
        return result.Error
    }

    tx.Commit()

    return nil
}

func (repo *HitsPostgresRepository) Get(ip string, session string) (entity.Hits, error) {
	repo.mutex.RLock() 
    defer repo.mutex.RUnlock()

    tx := repo.db.Begin()
    defer tx.Rollback()    

    var hit entity.Hits
    repo.db.Model(&entity.Hits{}).Where("ip = ? AND session = ?", ip, session).First(&hit)

    if repo.db.Error != nil {
        return entity.Hits{}, repo.db.Error
    }

    tx.Commit()

    return hit, nil
}