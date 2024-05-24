package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type TopRepository interface {
	Create(tops []entity.Top) error
	TopTruncateTable() error
	FindAll() ([]entity.Top, error)
}
