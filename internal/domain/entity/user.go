package entity

import (
	"errors"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrEmail = errors.New("email não pode ser vazio e deve ter mínimo de 7 e máximo de 80 caracteres")
	ErrEmailNotValid = errors.New("o email informado não é válido")
	ErrPassword = errors.New("password não pode ser vazio e deve ter mínimo de 6 e máximo de 80 caracteres")
)

type User struct {
	gorm.Model

	ID			string    	`gorm:"column:id;primaryKey" json:"id"`
	Email   	string    	`gorm:"column:email;unique;size:80;not null" json:"email"`	
	Password    string    	`gorm:"column:password;size:80;not null" json:"password"`	
	Admin  		bool      	`gorm:"column:admin;default:false" json:"admin"`
	CreatedAt	time.Time 	`gorm:"column:created_at" json:"created_at"`
	UpdatedAt	time.Time 	`gorm:"column:updated_at" json:"updated_at"`
}

type UserInput struct {
	Email		string
	Password	string
}

func NewUser(user UserInput) *User {	

	return &User{
		ID: uuid.NewString(),
		Email: user.Email,
		Password: user.Password,
		Admin: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) Validations() (bool, error) {	

	if u.Email == "" || len(u.Email) < 7 || len(u.Email) > 80 { 
		return false, ErrEmail
	}

	validEmail := utils.IsValidEmail(u.Email)	

	if !validEmail { 
		return false, ErrEmailNotValid
	 }	

	if u.Password == "" || len(u.Password) < 6 || len(u.Password) > 80 { 		
		return false, ErrPassword
	}

	return true, nil
	
}

