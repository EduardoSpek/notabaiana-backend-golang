package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GoogleRecaptchaVerify(r *http.Request) bool {
	// Captura o token enviado pelo cliente
	token := r.FormValue("g-recaptcha-response")

	// Faça uma solicitação POST para a API de verificação do reCAPTCHA v3 do Google
	response, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		map[string][]string{
			"secret":   {os.Getenv("KEY_GOOGLE_RECAPTCHA")},
			"response": {token},
		})

	if err != nil {
		fmt.Println("Erro ao fazer a solicitação:", err)
		return false
	}
	defer response.Body.Close()

	// Lê a resposta da API
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return false
	}

	// Decodifica a resposta JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Erro ao decodificar a resposta:", err)
		return false
	}

	// Verifica se a resposta foi bem-sucedida
	success := result["success"].(bool)

	return success
}
