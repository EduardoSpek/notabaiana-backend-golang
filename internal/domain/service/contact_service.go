package service

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

type ContactService struct {
	ContactRepository port.ContactRepository
}

func NewContactService(contato_repository port.ContactRepository) *ContactService {
	return &ContactService{ContactRepository: contato_repository}
}

func (cs *ContactService) AdminCreate(contato entity.ContactDTO) (entity.ContactDTO, error) {

	newcontato := entity.NewContact(contato)
	_, err := newcontato.Validations()

	if err != nil {
		return entity.ContactDTO{}, err
	}

	contato, err = cs.ContactRepository.AdminCreate(*newcontato)

	if err != nil {
		return entity.ContactDTO{}, err
	}
	return contato, nil
}

func (cs *ContactService) AdminFindAll() ([]entity.ContactDTO, error) {

	lista, err := cs.ContactRepository.AdminFindAll()

	if err != nil {
		return []entity.ContactDTO{}, err
	}
	return lista, nil
}

func (cs *ContactService) AdminGetByID(id string) (entity.ContactDTO, error) {

	contato, err := cs.ContactRepository.AdminGetByID(id)

	if err != nil {
		return entity.ContactDTO{}, err
	}
	return contato, nil
}

func (cs *ContactService) AdminDelete(id string) error {

	err := cs.ContactRepository.AdminDelete(id)

	if err != nil {
		return err
	}
	return nil
}
