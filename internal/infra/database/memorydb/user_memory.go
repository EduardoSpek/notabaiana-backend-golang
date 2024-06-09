package memorydb

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
)

type UserMemoryRepository struct {
	Userdb map[string]entity.User
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{ Userdb: make(map[string]entity.User) }
}

func (r *UserMemoryRepository) Update(user entity.User) (entity.User, error) {
	updatedUser := r.Userdb[user.ID]
	updatedUser.Email = user.Email
	updatedUser.Password = user.Password
	updatedUser.UpdatedAt = user.UpdatedAt

	r.Userdb[user.ID] = updatedUser
	return updatedUser, nil
}

func (r *UserMemoryRepository) Create(news entity.User) (entity.User, error) {
	r.Userdb[news.ID] = news
	return news, nil
}