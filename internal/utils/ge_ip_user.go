package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetIP() string {	

	// URL da API
	url := "https://api.ipify.org?format=text"

	// Fazer a requisição HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Erro ao fazer a requisição: %v", err)
	}
	defer resp.Body.Close()

	// Ler o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro ao ler o corpo da resposta: %v", err)
	}

	ip := string(body)

	fmt.Println(ip)

	return ip
}