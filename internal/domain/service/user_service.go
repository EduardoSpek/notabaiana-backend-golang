package service

import (
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

type UserService struct {
	UserRepository port.UserRepository
}

func NewUserService(user_repository port.UserRepository) *UserService {
	return &UserService{ UserRepository: user_repository }
}

func (uc *UserService) UpdateUser(id string, user entity.UserInput) (interface{}, error) {
	newuser := entity.NewUpdateUser(id, user)
	_, err := newuser.Validations()

	if err != nil {
		return &entity.User{}, err
	}

	encryptedPassword, err := utils.EncryptPassword(newuser.Password)

	if err != nil {
		return &entity.User{}, err
	}

	newuser.Password = encryptedPassword

	userUpdated, err := uc.UserRepository.Update(*newuser)

	if err != nil {
		return &entity.User{}, err
	}

	userOutput := struct{
		ID string `json:"id"`
		Email string `json:"email"`		
		Admin bool `json:"admin"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID: userUpdated.ID,
		Email: userUpdated.Email,		
		Admin: userUpdated.Admin,
		UpdatedAt: userUpdated.UpdatedAt,
	}

	return userOutput, nil
}

func (uc *UserService) CreateUser(user entity.UserInput) (interface{}, error) {
	newuser := entity.NewUser(user)
	_, err := newuser.Validations()

	if err != nil {
		return &entity.User{}, err
	}

	encryptedPassword, err := utils.EncryptPassword(newuser.Password)

	if err != nil {
		return &entity.User{}, err
	}

	newuser.Password = encryptedPassword

	userCreated, err := uc.UserRepository.Create(*newuser)

	if err != nil {
		return &entity.User{}, err
	}

	userOutput := struct{
		ID string `json:"id"`
		Email string `json:"email"`		
		Admin bool `json:"admin"`
		CreatedAt time.Time `json:"created_at"`
	}{
		ID: userCreated.ID,
		Email: userCreated.Email,		
		Admin: userCreated.Admin,
		CreatedAt: userCreated.CreatedAt,
	}

	return userOutput, nil
}