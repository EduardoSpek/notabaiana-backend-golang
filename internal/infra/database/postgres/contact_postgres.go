package postgres

import (
	"errors"
	"fmt"
	"sync"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"gorm.io/gorm"
)

var (
	ErrContactEmpty    = errors.New("não existe contato com o ID informado")
	ErrContactNotFound = errors.New("contato não encontrado")
)

type ContactPostgresRepository struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

func NewContactPostgresRepository(db_adapter port.DBAdapter) *ContactPostgresRepository {
	db := db_adapter.GetDB()
	return &ContactPostgresRepository{db: db}
}

func (repo *ContactPostgresRepository) AdminCreate(contact entity.Contact) (entity.ContactDTO, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	result := repo.db.Create(&contact)

	if result.Error != nil {
		tx.Rollback()
		return entity.ContactDTO{}, result.Error
	}

	tx.Commit()

	dto := entity.ContactDTO{
		ID:       contact.ID,
		Name:     contact.Name,
		Email:    contact.Email,
		Title:    contact.Title,
		Text:     contact.Text,
		Answered: contact.Answered,
	}

	return dto, nil
}

func (repo *ContactPostgresRepository) AdminFindAll() ([]entity.ContactDTO, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var contacts []entity.ContactDTO
	list := repo.db.Model(&entity.Contact{}).Order("created_at DESC, Answered ASC").Find(&contacts)

	if list.Error != nil {
		return []entity.ContactDTO{}, list.Error
	}

	tx.Commit()

	return contacts, nil
}

func (repo *ContactPostgresRepository) AdminGetByID(id string) (entity.ContactDTO, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var contact entity.Contact
	contactSelected := repo.db.Model(&entity.Contact{}).Where("id = ?", id).First(&contact)

	if contactSelected.Error != nil {
		return entity.ContactDTO{}, ErrContactNotFound
	}

	tx.Commit()

	dto := entity.ContactDTO{
		ID:        contact.ID,
		Name:      contact.Name,
		Email:     contact.Email,
		Title:     contact.Title,
		Text:      contact.Text,
		Image:     contact.Image,
		Answered:  contact.Answered,
		CreatedAt: contact.CreatedAt.Local(),
	}

	return dto, nil
}

func (repo *ContactPostgresRepository) AdminDelete(id string) error {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	var contact entity.Contact
	contactSelected := repo.db.Model(&entity.Contact{}).Where("id = ?", id).First(&contact)

	if contactSelected.Error != nil {
		return ErrContactNotFound
	}

	err := utils.RemoveImage("." + contact.Image)

	if !err {
		fmt.Println("Contacts: Imagem não deletada")
	}

	repo.db.Unscoped().Delete(contact)

	tx.Commit()

	return nil

}

func (repo *ContactPostgresRepository) AdminDeleteAll(contacts []entity.ContactDTO) error {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	tx := repo.db.Begin()
	defer tx.Rollback()

	for _, c := range contacts {

		var contact entity.Contact
		contactSelected := repo.db.Model(&entity.Contact{}).Where("id = ?", c.ID).First(&contact)

		if contactSelected.Error != nil {
			return ErrContactNotFound
		}

		err := utils.RemoveImage("." + contact.Image)

		if !err {
			fmt.Println("Contacts: Imagem não deletada")
		}

		repo.db.Unscoped().Delete(contact)

	}

	tx.Commit()

	return nil
}
