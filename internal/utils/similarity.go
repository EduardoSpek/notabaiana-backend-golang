package utils

import (
	"math"
	"strings"
)

// Função para calcular a distância de Levenshtein entre duas strings
func levenshtein(a, b string) int {
	la := len(a)
	lb := len(b)

	// Criação de uma matriz 2D para armazenar as distâncias
	d := make([][]int, la+1)
	for i := range d {
		d[i] = make([]int, lb+1)
	}

	// Inicialização das bordas da matriz
	for i := 0; i <= la; i++ {
		d[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		d[0][j] = j
	}

	// Preenchimento da matriz
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			d[i][j] = min(
				d[i-1][j]+1,    // deleção
				d[i][j-1]+1,    // inserção
				d[i-1][j-1]+cost, // substituição
			)
		}
	}

	return d[la][lb]
}

// Função auxiliar para encontrar o valor mínimo entre três números
func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

// Função para calcular a porcentagem de similaridade entre dois títulos
func Similarity(a, b string) float64 {
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	distance := levenshtein(a, b)
	maxLen := math.Max(float64(len(a)), float64(len(b)))
	return (1 - float64(distance)/maxLen) * 100
}

// func main() {
// 	title1 := "Exemplo de título de notícia"
// 	title2 := "Exemplo de título da notícia"
// 	similarityPercentage := similarity(title1, title2)
// 	fmt.Printf("Similaridade: %.2f%%\n", similarityPercentage)
// }
