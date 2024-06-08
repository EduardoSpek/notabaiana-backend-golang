package test

import (
	"testing"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
)

func TestUserEntity(t *testing.T) {
	t.Parallel()

	userInput := entity.UserInput{
		Email: "eu@vc.com",
		Password: "q1w2e3",
	}
	
	user := entity.NewUser(userInput)

	_, err := user.Validations() 

	if err != nil {
		t.Error(err)
	}

	testcases := []TestCase{
		{
			Esperado: "eu@vc.com",
			Recebido: user.Email,
			Descricao: "Validação do Email",
		},
		{
			Esperado: false,
			Recebido: user.Admin,
			Descricao: "Validação do Administrador",
		},
		
	}

	for _, teste := range testcases {
		Resultado(t, teste.Esperado, teste.Recebido, teste.Descricao)
	}

}