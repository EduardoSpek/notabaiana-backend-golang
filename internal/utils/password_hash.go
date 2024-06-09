package utils

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {

	const cost = 12

	// Gera um hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	// Retorna a senha encriptada como string
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}