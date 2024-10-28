package postgres

import (
	"errors"
	"sync"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("usuário não encontrado")
)

type UserPostgresRepository struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

func NewUserPostgresRepository(db_adapter port.DBAdapter) *UserPostgresRepository {
	db := db_adapter.GetDB()
	return &UserPostgresRepository{db: db}
}

func (repo *UserPostgresRepository) GetByID(id string) (entity.User, error) {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var user entity.User
	repo.db.Model(&entity.User{}).Where("id = ?", id).First(&user)

	if repo.db.Error != nil {
		return entity.User{}, ErrUserNotFound
	}

	tx.Commit()

	return user, nil

}

func (repo *UserPostgresRepository) GetByEmail(email string) (entity.User, error) {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var user entity.User
	repo.db.Model(&entity.User{}).Where("email = ?", email).First(&user)

	if repo.db.Error != nil {
		return entity.User{}, ErrUserNotFound
	}

	tx.Commit()

	return user, nil
}

func (repo *UserPostgresRepository) Update(user entity.User) (entity.User, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	result := repo.db.Model(&user).Updates(map[string]interface{}{
		"email":    user.Email,
		"password": user.Password,
	})

	if result.Error != nil {
		tx.Rollback()
		return entity.User{}, ErrUserNotFound
	}

	tx.Commit()

	return user, nil
}

func (repo *UserPostgresRepository) Create(user entity.User) (entity.User, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	result := repo.db.Create(&user)

	if result.Error != nil {
		tx.Rollback()
		return entity.User{}, result.Error
	}

	tx.Commit()

	return user, nil
}
