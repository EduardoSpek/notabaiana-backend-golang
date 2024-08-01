package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrName = errors.New("campo email não pode estar vazio e deve ter no máximo 80 caracteres")
	ErrText = errors.New("campo mensagem não pode estar vazio")
)

// Input and Output DTO
type ContactDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
	Answered  bool      `json:"answered"`
	CreatedAt time.Time `json:"created_at"`
}

type Contact struct {
	gorm.Model

	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Answered  bool      `json:"answered"`
}

func NewContact(contact ContactDTO) *Contact {

	return &Contact{
		ID:        uuid.NewString(),
		Name:      strings.TrimSpace(contact.Name),
		Email:     strings.TrimSpace(contact.Email),
		Title:     strings.TrimSpace(contact.Title),
		Text:      strings.TrimSpace(contact.Text),
		Image:     strings.TrimSpace(contact.Image),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Answered:  contact.Answered,
	}
}

func UpdateContact(contact ContactDTO) *Contact {

	return &Contact{
		ID:        strings.TrimSpace(contact.ID),
		Name:      strings.TrimSpace(contact.Name),
		Email:     strings.TrimSpace(contact.Email),
		Title:     strings.TrimSpace(contact.Title),
		Text:      strings.TrimSpace(contact.Text),
		Image:     strings.TrimSpace(contact.Image),
		UpdatedAt: time.Now(),
		Answered:  contact.Answered,
	}
}

func (c *Contact) Validations() (bool, error) {

	if c.Name == "" || len(c.Name) < 2 || len(c.Name) > 80 {
		return false, ErrName
	}

	if c.Title == "" || len(c.Title) < 2 || len(c.Title) > 80 {
		return false, ErrTitle
	}

	if c.Email == "" || len(c.Email) < 7 || len(c.Email) > 80 {
		return false, ErrEmail
	}

	validEmail := utils.IsValidEmail(c.Email)

	if !validEmail {
		return false, ErrEmailNotValid
	}

	if c.Text == "" || len(c.Text) < 2 {
		return false, ErrText
	}

	return true, nil

}
