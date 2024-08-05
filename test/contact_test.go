package test

import (
	"log"
	"reflect"
	"testing"

	"github.com/eduardospek/notabaiana-backend-golang/internal/adapter"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	database "github.com/eduardospek/notabaiana-backend-golang/internal/infra/database/postgres"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/joho/godotenv"
)

func isStruct(v interface{}) bool {

	// Verificando se o tipo é uma struct
	return reflect.TypeOf(v).Kind() == reflect.Struct
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

	contactDTO := entity.ContactDTO{
		Name:     "Eduardo Spek",
		Email:    "eu@vc.com",
		Title:    "Estou com uma dúvida Loren ipsun dolor sit iamet",
		Text:     "Loren ipsun dolor sit iamet",
		Answered: false,
	}

	contact := entity.NewContact(contactDTO)

	_, err := contact.Validations()

	if err != nil {
		t.Error(err)
	}

	testcases := []TestCase{
		{
			Esperado:  "Eduardo Spek",
			Recebido:  contact.Name,
			Descricao: "Validação do Name",
		},
		{
			Esperado:  "Estou com uma dúvida Loren ipsun dolor sit iamet",
			Recebido:  contact.Title,
			Descricao: "Validação do Title",
		},
		{
			Esperado:  "eu@vc.com",
			Recebido:  contact.Email,
			Descricao: "Validação do Email",
		},
		{
			Esperado:  false,
			Recebido:  contact.Answered,
			Descricao: "Respondido?",
		},
	}

	for _, teste := range testcases {
		Resultado(t, teste.Esperado, teste.Recebido, teste.Descricao)
	}

}

func TestContactService(t *testing.T) {
	t.Parallel()

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	postgres := adapter.NewPostgresAdapter()
	repo := database.NewContactPostgresRepository(postgres)
	imagedownloader := utils.NewImgDownloader()
	contact_service := service.NewContactService(repo, imagedownloader)

	var id, id2, id3 string

	t.Run("Deve Criar um novo contact", func(t *testing.T) {
		dto := entity.ContactDTO{
			Name:     "Eduardo Spek",
			Email:    "eu@vc.com.br",
			Title:    "Eu tenho uma dúvida",
			Text:     "Loren ipsun sit iamet",
			Answered: false,
		}

		newcontact, err := contact_service.AdminCreate(dto)

		if err != nil {
			t.Error(err)
		}

		id = newcontact.ID

		if newcontact.ID == "" {
			t.Error("ID vazio")
		}

		if !isStruct(newcontact) {
			t.Error()
		}

	})

	t.Run("Deve Criar um novo contact 2", func(t *testing.T) {
		dto := entity.ContactDTO{
			Name:     "Thaís Freire",
			Email:    "thais@freire.com.br",
			Title:    "Quero patrocinar vocês",
			Text:     "Loren ipsun sit iamet sit iamet",
			Answered: true,
		}

		newcontact, err := contact_service.AdminCreate(dto)

		if err != nil {
			t.Error(err)
		}

		id2 = newcontact.ID

		if !isStruct(newcontact) {
			t.Error()
		}

	})

	t.Run("Deve criar um novo contact 3", func(t *testing.T) {
		dto := entity.ContactDTO{
			Name:     "Nathan Freire",
			Email:    "nathan@freire.com.br",
			Title:    "Preciso de ajuda",
			Text:     "Loren ipsun sit iamet sit iamet sit iamet sit iamet sit iamet sit iamet",
			Answered: false,
		}

		newcontact, err := contact_service.AdminCreate(dto)

		if err != nil {
			t.Error(err)
		}

		id3 = newcontact.ID

		if !isStruct(newcontact) {
			t.Error()
		}

	})

	t.Run("Deve listar os contacts", func(t *testing.T) {
		lista, err := contact_service.AdminFindAll()

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

	t.Run("Deve obter o contact informado por ID", func(t *testing.T) {
		contact, err := contact_service.AdminGetByID(id)

		if err != nil {
			t.Error(err)
		}

		if contact.ID == "" {
			t.Error("ID vazio")
		}

		if isSliceOfStructs(contact) {
			t.Error(contact)
		}

	})

	t.Run("Deve deletar o contact por ID", func(t *testing.T) {
		err := contact_service.AdminDelete(id)

		if err != nil {
			t.Error(err)
		}

		lista, err := contact_service.AdminFindAll()

		if err != nil {
			t.Error(err)
		}

		esperado := 1

		if len(lista) < esperado {
			t.Errorf("Esperado um total de %d e retornado %d", esperado, len(lista))
		}

	})

	t.Run("Deve deletar os contacts informados", func(t *testing.T) {
		var contacts_list []entity.ContactDTO

		c1 := entity.ContactDTO{
			ID: id2,
		}

		contacts_list = append(contacts_list, c1)

		c2 := entity.ContactDTO{
			ID: id3,
		}

		contacts_list = append(contacts_list, c2)

		err := contact_service.AdminDeleteAll(contacts_list)

		if err != nil {
			t.Error(err)
		}

		lista, err := contact_service.AdminFindAll()

		if err != nil {
			t.Error(err)
		}

		esperado := 0

		if len(lista) < esperado {
			t.Errorf("Esperado um total de %d e retornado %d", esperado, len(lista))
		}

	})

}
