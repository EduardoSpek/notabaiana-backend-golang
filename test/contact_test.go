package test

import (
	"reflect"
	"testing"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	database "github.com/eduardospek/notabaiana-backend-golang/internal/infra/database/memorydb"
)

func isStruct(v interface{}) bool {
	// Obtendo o tipo da variável
	t := reflect.TypeOf(v)

	// Verificando se o tipo é uma struct
	return t.Kind() == reflect.Struct
}

func isSliceOfStructs(v interface{}) bool {
	// Obtendo o tipo da variável
	t := reflect.TypeOf(v)

	// Verificando se o tipo é um slice ou array
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return false
	}

	// Obtendo o tipo do elemento do slice/array
	elemType := t.Elem()

	// Verificando se o tipo do elemento é uma struct
	return elemType.Kind() == reflect.Struct
}

func TestContactEntity(t *testing.T) {
	t.Parallel()

	contatoDTO := entity.ContactDTO{
		Name:     "Eduardo Spek",
		Email:    "eu@vc.com",
		Title:    "Estou com uma dúvida Loren ipsun dolor sit iamet",
		Text:     "Loren ipsun dolor sit iamet",
		Answered: false,
	}

	contato := entity.NewContact(contatoDTO)

	_, err := contato.Validations()

	if err != nil {
		t.Error(err)
	}

	testcases := []TestCase{
		{
			Esperado:  "Eduardo Spek",
			Recebido:  contato.Name,
			Descricao: "Validação do Name",
		},
		{
			Esperado:  "Estou com uma dúvida Loren ipsun dolor sit iamet",
			Recebido:  contato.Title,
			Descricao: "Validação do Title",
		},
		{
			Esperado:  "eu@vc.com",
			Recebido:  contato.Email,
			Descricao: "Validação do Email",
		},
		{
			Esperado:  false,
			Recebido:  contato.Answered,
			Descricao: "Respondido?",
		},
	}

	for _, teste := range testcases {
		Resultado(t, teste.Esperado, teste.Recebido, teste.Descricao)
	}

}

func TestContactService(t *testing.T) {
	t.Parallel()

	repo := database.NewContactMemoryRepository()
	contato_service := service.NewContactService(repo)

	var id, id2 string

	t.Run("Deve Criar um novo contato", func(t *testing.T) {
		dto := entity.ContactDTO{
			Name:     "Eduardo Spek",
			Email:    "eu@vc.com.br",
			Title:    "Eu tenho uma dúvida",
			Text:     "Loren ipsun sit iamet",
			Answered: false,
		}

		newcontato, err := contato_service.AdminCreate(dto)

		if err != nil {
			t.Error(err)
		}

		id = newcontato.ID

		if !isStruct(newcontato) {
			t.Error()
		}

	})

	t.Run("Deve Criar um novo contato 2", func(t *testing.T) {
		dto := entity.ContactDTO{
			Name:     "Thaís Freire",
			Email:    "thais@freire.com.br",
			Title:    "Quero patrocinar vocês",
			Text:     "Loren ipsun sit iamet sit iamet",
			Answered: true,
		}

		newcontato, err := contato_service.AdminCreate(dto)

		if err != nil {
			t.Error(err)
		}

		id2 = newcontato.ID

		if !isStruct(newcontato) {
			t.Error()
		}

	})

	t.Run("Deve listar os contatos", func(t *testing.T) {
		lista, err := contato_service.AdminFindAll()

		if err != nil {
			t.Error(err)
		}

		esperado := 2

		if len(lista) < esperado {
			t.Errorf("Esperado um total de %d e retornado %d", esperado, len(lista))
		}

		if !isSliceOfStructs(lista) {
			t.Error(lista)
		}

	})

	t.Run("Deve obter o contato informado por ID", func(t *testing.T) {
		contato, err := contato_service.AdminGetByID(id)

		if err != nil {
			t.Error(err)
		}

		if isSliceOfStructs(contato) {
			t.Error(contato)
		}

	})

	t.Run("Deve deletar o contato por ID", func(t *testing.T) {
		err := contato_service.AdminDelete(id)

		if err != nil {
			t.Error(err)
		}

		lista, err := contato_service.AdminFindAll()

		if err != nil {
			t.Error(err)
		}

		esperado := 1

		if len(lista) < esperado {
			t.Errorf("Esperado um total de %d e retornado %d", esperado, len(lista))
		}

	})

	t.Run("Deve deletar os contatos informados", func(t *testing.T) {
		var contacts_list []entity.ContactDTO

		c1 := entity.ContactDTO{
			ID: id,
		}

		contacts_list = append(contacts_list, c1)

		c2 := entity.ContactDTO{
			ID: id2,
		}

		contacts_list = append(contacts_list, c2)

		err := contato_service.AdminDeleteAll(contacts_list)

		if err != nil {
			t.Error(err)
		}

		lista, err := contato_service.AdminFindAll()

		if err != nil {
			t.Error(err)
		}

		esperado := 0

		if len(lista) < esperado {
			t.Errorf("Esperado um total de %d e retornado %d", esperado, len(lista))
		}

	})

}
