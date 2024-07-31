package memorydb

import (
	"errors"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
)

var (
	ErrContactEmpty = errors.New("não existe contato com o ID informado")
)

type ContactMemoryRepository struct {
	db map[string]entity.Contact
}

func NewContactMemoryRepository() *ContactMemoryRepository {
	return &ContactMemoryRepository{db: make(map[string]entity.Contact)}
}

func (c *ContactMemoryRepository) AdminCreate(contato entity.Contact) (entity.ContactDTO, error) {
	c.db[contato.ID] = contato
	dto := entity.ContactDTO{
		ID:       contato.ID,
		Name:     contato.Name,
		Email:    contato.Email,
		Title:    contato.Title,
		Text:     contato.Text,
		Answered: contato.Answered,
	}
	return dto, nil
}

func (c *ContactMemoryRepository) AdminFindAll() ([]entity.ContactDTO, error) {
	var lista []entity.ContactDTO
	for _, contato := range c.db {
		dto := entity.ContactDTO{
			ID:       contato.ID,
			Name:     contato.Name,
			Email:    contato.Email,
			Title:    contato.Title,
			Text:     contato.Text,
			Answered: contato.Answered,
		}
		lista = append(lista, dto)
	}

	return lista, nil
}

func (c *ContactMemoryRepository) AdminGetByID(id string) (entity.ContactDTO, error) {
	contato := c.db[id]

	if contato.ID == "" {
		return entity.ContactDTO{}, ErrContactEmpty
	}

	dto := entity.ContactDTO{
		ID:       contato.ID,
		Name:     contato.Name,
		Email:    contato.Email,
		Title:    contato.Title,
		Text:     contato.Text,
		Answered: contato.Answered,
	}

	return dto, nil
}

func (c *ContactMemoryRepository) AdminDelete(id string) error {

	if _, exists := c.db[id]; exists {
		// Remove o item do mapa
		delete(c.db, id)
		return nil
	} else {
		return ErrContactEmpty
	}
}

func (c *ContactMemoryRepository) AdminDeleteAll(contacts []entity.ContactDTO) error {

	for _, cc := range contacts {

		delete(c.db, cc.ID)

	}

	return nil
}
