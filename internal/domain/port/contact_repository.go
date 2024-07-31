package port

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"

type ContactRepository interface {
	AdminCreate(contato entity.Contact) (entity.ContactDTO, error)
	AdminGetByID(id string) (entity.ContactDTO, error)
	AdminFindAll() ([]entity.ContactDTO, error)
	AdminDelete(id string) error
	//AdminDeleteAll(contatos []entity.ContactDTO) error
}
