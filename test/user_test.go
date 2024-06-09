package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	database "github.com/eduardospek/notabaiana-backend-golang/internal/infra/database/memorydb"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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


func TestUserService(t *testing.T) {
	t.Parallel()

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Erro ao carregar arquivo .env: %v", err)
	}

	t.Run("Deve Criar um novo usuário", func(t *testing.T) {

		user := &entity.UserInput{
			Email: "eduardospekoficial@gmail.com",
			Password: "q1w2e3",
		}

		userJson, err := json.Marshal(user)

		if err != nil {
			t.Fatalf("Erro ao converter usuário para JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(userJson))
		
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")


		repo := database.NewUserMemoryRepository()
		user_service := service.NewUserService(repo)		
		controller := controllers.NewUserController(*user_service)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", controller.CreateUser).Methods("POST")

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		fmt.Println(rr.Body.String())
		t.Fail()
	})
}