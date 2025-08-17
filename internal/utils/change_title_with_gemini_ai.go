package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

/* type Response struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Contents Content `json:"content"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
} */

// func ChangeTitleWithGemini(prompt, title string) (string, error) {

// 	ctx := context.Background()
// 	// Access your API key as an environment variable (see "Set up your API key" above)
// 	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("KEY_GEMINI")))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Close()

// 	// The Gemini 1.5 models are versatile and work with both text-only and multimodal prompts
// 	model := client.GenerativeModel("gemini-2.0-flash")
// 	resp, err := model.GenerateContent(ctx, genai.Text(string(prompt+" "+title)))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Convert the response to a JSON string with indentation
// 	respJSON, err := json.MarshalIndent(resp.Candidates[0].Content.Parts[0], "", "  ")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return string(respJSON[1 : len(respJSON)-1]), nil
// }

/* func ChangeTitleWithGemini(prompt, title string) (string, error) {

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash-lite-preview-06-17:generateContent?key=" + os.Getenv("KEY_GEMINI")

	title = strings.ReplaceAll(title, `"`, `\"`)
	title = strings.ReplaceAll(title, `'`, `\'`)

	jsonData := `{"contents":[{"parts":[{"text":"` + prompt + ` : ` + title + `"}]}]}`

	reqBody := bytes.NewBuffer([]byte(jsonData))

	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição POST:", err)
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return "", err
	}

	// Deserializa o JSON na struct Response
	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Println("Erro ao deserializar o JSON:", err)
		return "", err
	}

	var newtitle string
	// Acessa o valor de "text"
	for _, candidate := range response.Candidates {
		for _, part := range candidate.Contents.Parts {

			newtitle = strings.Replace(part.Text, "**", "", -1)

		}
	}

	return newtitle, nil

} */

func ChangeTitleWithGemini(prompt, title string) (string, error) {

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("KEY_GEMINI")))
	if err != nil {
		fmt.Println("Erro na key/chave do Gemini: ", err)
		return "", err
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with both text-only and multimodal prompts
	model := client.GenerativeModel("gemini-2.5-flash-lite")
	resp, err := model.GenerateContent(ctx, genai.Text(string(prompt+" "+title)))
	if err != nil {
		fmt.Println("Erro ao gerar texto no Gemini: ", err)
		return "", err
	}

	// Convert the response to a JSON string with indentation
	respJSON, err := json.Marshal(resp.Candidates[0].Content.Parts[0])
	if err != nil {
		fmt.Println("Erro ao converter JSON:", err)
		return "", err
	}

	var jsonStr string
	err = json.Unmarshal(respJSON, &jsonStr)

	if err != nil {
		fmt.Println("Erro ao deserializar string JSON:", err)
		return "", err
	}

	jsonStr = strings.Replace(jsonStr, "```json", "", -1)
	jsonStr = strings.Trim(jsonStr, "`")
	jsonStr = strings.Trim(jsonStr, " ")

	return jsonStr, nil
}
