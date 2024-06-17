package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/service"
	database "github.com/eduardospek/notabaiana-backend-golang/internal/infra/database/memorydb"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/middlewares"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func TestUserEntity(t *testing.T) {
	t.Parallel()

	userInput := entity.UserInput{
		Email:    "eu@vc.com",
		Password: "q1w2e3",
	}

	user := entity.NewUser(userInput)

	_, err := user.Validations()

	if err != nil {
		t.Error(err)
	}

	testcases := []TestCase{
		{
			Esperado:  "eu@vc.com",
			Recebido:  user.Email,
			Descricao: "Validação do Email",
		},
		{
			Esperado:  false,
			Recebido:  user.Admin,
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

	repo := database.NewUserMemoryRepository()
	user_service := service.NewUserService(repo)
	controller := controllers.NewUserController(*user_service)

	var responseRoute entity.User

	var token struct {
		Token string
	}

	t.Run("Deve Criar um novo usuário", func(t *testing.T) {

		user := &entity.UserInput{
			Email:    "eu@vc.com",
			Password: "q1w2e3",
			Admin:    true,
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

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", controller.CreateUser).Methods("POST")

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Esperado: %v - Recebido: %v",
				http.StatusOK, status)
		}

		err = json.NewDecoder(rr.Body).Decode(&responseRoute)

		if err != nil {
			t.Fatalf("Erro ao decodificar resposta JSON: %v", err)
		}

	})

	t.Run("Deve fazer o login", func(t *testing.T) {

		user := &entity.UserInput{
			Email:    "eu@vc.com",
			Password: "q1w2e3",
		}

		userJson, err := json.Marshal(user)

		if err != nil {
			t.Fatalf("Erro ao converter usuário para JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(userJson))

		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/login", controller.Login).Methods("POST")

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Esperado: %v - Recebido: %v",
				http.StatusOK, status)
		}

		err = json.NewDecoder(rr.Body).Decode(&token)

		if err != nil {
			t.Fatalf("Erro ao decodificar resposta JSON: %v", err)
		}
	})

	t.Run("Deve atualizar um usuário", func(t *testing.T) {

		user := &entity.UserInput{
			Email:    "vc@laele.com",
			Password: "p0o9i8u7",
		}

		userJson, err := json.Marshal(user)

		if err != nil {
			t.Fatalf("Erro ao converter usuário para JSON: %v", err)
		}

		req, err := http.NewRequest("PUT", "/user/"+responseRoute.ID, bytes.NewBuffer(userJson))

		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token.Token)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.Handle("/user/{id}", middlewares.JwtMiddleware(http.HandlerFunc(controller.UpdateUser))).Methods("PUT")

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Esperado: %v - Recebido: %v",
				http.StatusOK, status)
		}

		err = json.NewDecoder(rr.Body).Decode(&responseRoute)

		if err != nil {
			t.Fatalf("Erro ao decodificar resposta JSON: %v", err)
		}
	})

	t.Run("Deve checar se o usuário autenticado é admin", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/accesscheck", nil)

		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token.Token)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.Handle("/accesscheck", middlewares.JwtMiddleware(http.HandlerFunc(controller.AccessCheck))).Methods("GET")

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Esperado: %v - Recebido: %v",
				http.StatusOK, status)
		}

		err = json.NewDecoder(rr.Body).Decode(&responseRoute)

		if err != nil {
			t.Fatalf("Erro ao decodificar resposta JSON: %v", err)
		}
	})
}
