package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type ContactRepository interface {
	AdminCreate(contact entity.Contact) (entity.ContactDTO, error)
	AdminGetByID(id string) (entity.ContactDTO, error)
	AdminFindAll() ([]entity.ContactDTO, error)
	AdminDelete(id string) error
	AdminDeleteAll(contacts []entity.ContactDTO) error
}
