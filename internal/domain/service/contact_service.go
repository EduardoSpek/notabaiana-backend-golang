package service

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

type ContactService struct {
	ContactRepository port.ContactRepository
}

func NewContactService(contact_repository port.ContactRepository) *ContactService {
	return &ContactService{ContactRepository: contact_repository}
}

func (cs *ContactService) AdminCreate(contact entity.ContactDTO) (entity.ContactDTO, error) {

	newcontact := entity.NewContact(contact)
	_, err := newcontact.Validations()

	if err != nil {
		return entity.ContactDTO{}, err
	}

	contact, err = cs.ContactRepository.AdminCreate(*newcontact)

	if err != nil {
		return entity.ContactDTO{}, err
	}
	return contact, nil
}

func (cs *ContactService) AdminFindAll() ([]entity.ContactDTO, error) {

	lista, err := cs.ContactRepository.AdminFindAll()

	if err != nil {
		return []entity.ContactDTO{}, err
	}
	return lista, nil
}

func (cs *ContactService) AdminGetByID(id string) (entity.ContactDTO, error) {

	contact, err := cs.ContactRepository.AdminGetByID(id)

	if err != nil {
		return entity.ContactDTO{}, err
	}
	return contact, nil
}

func (cs *ContactService) AdminDelete(id string) error {

	err := cs.ContactRepository.AdminDelete(id)

	if err != nil {
		return err
	}
	return nil
}

func (cs *ContactService) AdminDeleteAll(contacts []entity.ContactDTO) error {

	err := cs.ContactRepository.AdminDeleteAll(contacts)

	if err != nil {
		return err
	}
	return nil
}
