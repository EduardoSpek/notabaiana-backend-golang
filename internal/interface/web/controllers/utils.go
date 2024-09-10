package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

var (
	ErrToken = errors.New("acesso não autorizado")
)

func ResponseJson(w http.ResponseWriter, data any, statusCode int) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.New("responseJson: não foi possível converter para json")
	}

	// Escrevendo a resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)

	return nil

}

func TokenVerifyByForm(w http.ResponseWriter, r *http.Request) error {

	token := r.FormValue("token")

	if token == "" {
		return ErrToken
	}

	claims, err := utils.ValidateJWT(token)

	if err != nil {
		return ErrToken
	}

	if !claims.Admin {
		return ErrToken
	}

	return nil
}
func TokenVerifyByHeader(w http.ResponseWriter, r *http.Request) error {

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		return ErrToken
	}

	tokenStr = tokenStr[len("Bearer "):]

	claims, err := utils.ValidateJWT(tokenStr)

	if err != nil {
		return ErrToken
	}

	if !claims.Admin {
		return ErrToken
	}

	return nil

}

func SaveImageForm(file multipart.File, filename string, pasta string) error {

	if file == nil {
		return nil
	}

	defer file.Close()

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("Erro ao obter o caminho do executável:", err)
	}

	diretorio := strings.Replace(cwd, "test", "", -1) + "/images/" + pasta
	pathImage := diretorio + filename

	f, err := os.Create(pathImage)
	if err != nil {
		return ErrParseForm
	}
	defer f.Close()
	io.Copy(f, file)

	f, err = os.Open(pathImage)

	if err != nil {
		return ErrParseForm
	}

	// Resize the image
	img, _, err := image.Decode(f)
	if err != nil {
		return ErrDecodeImage
	}

	ImgDownloader := utils.NewImgDownloader()

	err = ImgDownloader.CropAndSaveImage(img, 400, 254, pathImage)

	if err != nil {
		fmt.Println(err)
		return ErrDecodeImage
	}

	return nil

}
