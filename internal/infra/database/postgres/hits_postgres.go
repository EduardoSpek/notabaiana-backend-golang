package postgres

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"gorm.io/gorm"
)

type HitsPostgresRepository struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

func NewHitsPostgresRepository(db_adapter port.DBAdapter) *HitsPostgresRepository {
	db := db_adapter.GetDB()
	return &HitsPostgresRepository{db: db}
}

func (repo *HitsPostgresRepository) Save(hit entity.Hits) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	result := repo.db.Create(&hit)

	if result.Error != nil {
		tx.Rollback()
		fmt.Println(result.Error)
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

	result := repo.db.Model(&hit).Update("views", hit.Views)

	if result.Error != nil {
		tx.Rollback()
		fmt.Println(result.Error)
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

	if hit.IP == "" {
		fmt.Println("=== Nenhum hit encontrado ===")
		return entity.Hits{}, errors.New("=== Nenhum hit encontrado ===")
	}

	tx.Commit()

	return hit, nil
}

func (repo *HitsPostgresRepository) TopHits() ([]entity.Hits, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var hits []entity.Hits
	repo.db.Model(&entity.Hits{}).Select("session, sum(views) as total").Where("created_at >= ? AND created_at <= ?", time.Now().AddDate(0, 0, -2), time.Now()).Order("total DESC").Group("session").Limit(10).Find(&hits)

	tx.Commit()

	return hits, nil
}
