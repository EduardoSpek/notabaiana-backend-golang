package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eduardospek/notabaiana-backend-golang/internal/adapter"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/usecase"
	database "github.com/eduardospek/notabaiana-backend-golang/internal/infra/database/postgres"
	"github.com/eduardospek/notabaiana-backend-golang/internal/interface/web/controllers"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func TestDownloadEntity(t *testing.T) {
	t.Parallel()

	downloadDTO := entity.Download{
		Category: "pagode",
		Title:    "Harmnonia do Samba",
		Link:     "https://www.suamusica.com.br/harmonia-do-samba",
		Text:     "Loren ipsun dolor sit iamet",
		Visible:  false,
	}

	download := entity.NewDownload(downloadDTO)

	_, err := download.Validations()

	if err != nil {
		t.Error(err)
	}

	testcases := []TestCase{
		{
			Esperado:  "Harmnonia do Samba",
			Recebido:  download.Title,
			Descricao: "Validação do título",
		},
		{
			Esperado:  "Loren ipsun dolor sit iamet",
			Recebido:  download.Text,
			Descricao: "Validação do texto",
		},
		{
			Esperado:  "https://www.suamusica.com.br/harmonia-do-samba",
			Recebido:  download.Link,
			Descricao: "Validação do Email",
		},
		{
			Esperado:  false,
			Recebido:  download.Visible,
			Descricao: "Visível",
		},
	}

	for _, teste := range testcases {
		Resultado(t, teste.Esperado, teste.Recebido, teste.Descricao)
	}

}

func TestDownloadUsecase(t *testing.T) {
	t.Parallel()

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	postgres := adapter.NewPostgresAdapter()
	repo := database.NewDownloadPostgresRepository(postgres)
	//imagedownloader := utils.NewImgDownloader()

	t.Run("Deve Criar um novo download", func(t *testing.T) {

		dto := &entity.Download{
			Category: "pagode",
			Title:    "Harmnonia do Samba 1",
			Link:     "https://www.suamusica.com.br/harmonia-do-samba",
			Text:     "Loren ipsun dolor sit iamet",
		}

		createDownloadUsecase := usecase.NewCreateDownloadUsecase(repo)
		downloadCreated, err := createDownloadUsecase.Create(dto)

		if err != nil {
			t.Error(err)
		}

		if downloadCreated.ID == "" {
			t.Error("ID vazio")
		}

		if !isStruct(*downloadCreated) {
			t.Error()
		}

	})
}

func TestDownloadController(t *testing.T) {
	t.Parallel()

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	postgres := adapter.NewPostgresAdapter()
	repo := database.NewDownloadPostgresRepository(postgres)
	imagedownloader := utils.NewImgDownloader()

	var responseRoute *entity.Download

	t.Run("Deve Criar um novo download", func(t *testing.T) {

		dto := &entity.Download{
			Category: "pagode",
			Title:    "Harmnonia do Samba 1",
			Link:     "https://www.suamusica.com.br/harmonia-do-samba",
			Text:     "Loren ipsun dolor sit iamet",
		}

		controller := controllers.NewDownloadController(repo, imagedownloader)

		formData := url.Values{}
		formData.Set("category", dto.Category)
		formData.Set("title", dto.Title)
		formData.Set("link", dto.Link)
		formData.Set("text", dto.Text)

		if err != nil {
			t.Fatalf("Erro ao converter usuário para JSON: %v", err)
		}

		req, err := http.NewRequest("POST", "/admin/downloads/create", strings.NewReader(formData.Encode()))

		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/admin/downloads/create", controller.CreateDownloadUsingTheForm).Methods("POST")

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

	t.Run("Deve Criar um novo download com imagem", func(t *testing.T) {

		dto := &entity.Download{
			Category: "pagode",
			Title:    "Harmnonia do Samba 11",
			Link:     "https://www.suamusica.com.br/harmonia-do-samba11",
			Text:     "Loren ipsun dolor sit iamet",
		}

		controller := controllers.NewDownloadController(repo, imagedownloader)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("nome", "João")
		_ = writer.WriteField("email", "joao@example.com")

		writer.WriteField("category", dto.Category)
		writer.WriteField("title", dto.Title)
		writer.WriteField("link", dto.Link)
		writer.WriteField("text", dto.Text)

		// Adicione o arquivo de imagem
		file, err := os.Open("../files/base_image.jpg") // Certifique-se de que este arquivo existe
		if err != nil {
			//t.Fatal(err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("image", filepath.Base(file.Name()))
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			t.Fatal(err)
		}

		// Feche o writer multipart
		err = writer.Close()
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/admin/downloads/create", body)

		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/admin/downloads/create", controller.CreateDownloadUsingTheForm).Methods("POST")

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
