package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type UserRepository interface {
	Create(user entity.User) (entity.User, error)
	Update(user entity.User) (entity.User, error)
}