package utils

import "math"

// Pagination recebe a página atual e o total de noticias para retornar a páginação de resultado
func Pagination(currentPage, total int) map[string][]int {

	// Calcula o total de páginas
	totalPages := int(math.Ceil(float64(total) / 16))

	// Garante que a página atual esteja dentro dos limites
	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}

	previousPages := []int{}
	nextPages := []int{}

	if currentPage == totalPages {
		if currentPage > 2 {
			previousPages = []int{currentPage - 2, currentPage - 1}
			nextPages = []int{}
		} else {
			previousPages = []int{currentPage - 1}
			nextPages = []int{}
		}
	} else if currentPage-2 > 2 && currentPage+2 <= totalPages {
		previousPages = []int{currentPage - 2, currentPage - 1}
		nextPages = []int{currentPage + 1, currentPage + 2}
	} else if currentPage == 1 && currentPage == totalPages {
		previousPages = []int{}
		nextPages = []int{}
	} else if currentPage == 1 && totalPages < 3 {
		previousPages = []int{}
		nextPages = []int{currentPage + 1}
	} else if currentPage == 1 && totalPages > 2 {
		previousPages = []int{}
		nextPages = []int{currentPage + 1, currentPage + 2}
	} else if currentPage == 2 && currentPage == totalPages {
		previousPages = []int{1}
		nextPages = []int{}
	} else if currentPage == 2 && totalPages < 4 {
		previousPages = []int{currentPage - 1}
		nextPages = []int{currentPage + 1}
	} else if currentPage == 2 && totalPages > 3 {
		previousPages = []int{currentPage - 1}
		nextPages = []int{currentPage + 1, currentPage + 2}
	} else if currentPage == 3 && currentPage == totalPages {
		previousPages = []int{1, 2}
		nextPages = []int{}
	} else if currentPage == 3 && totalPages < 5 {
		previousPages = []int{1, 2}
		nextPages = []int{currentPage + 1}
	} else if currentPage == 3 && totalPages > 4 {
		previousPages = []int{1, 2}
		nextPages = []int{currentPage + 1, currentPage + 2}
	} else if currentPage == 4 && currentPage == totalPages {
		previousPages = []int{currentPage - 2, currentPage - 1}
		nextPages = []int{}
	} else if currentPage == 4 && totalPages < 6 {
		previousPages = []int{currentPage - 2, currentPage - 1}
		nextPages = []int{currentPage + 1}
	} else if currentPage == 4 && totalPages > 5 {
		previousPages = []int{currentPage - 2, currentPage - 1}
		nextPages = []int{currentPage + 1, currentPage + 2}
	} else if currentPage == 5 && totalPages < 7 {
		previousPages = []int{currentPage - 2, currentPage - 1}
		nextPages = []int{currentPage + 1}
	} else if currentPage > 3 && totalPages > currentPage && currentPage+1 <= totalPages {
		previousPages = []int{currentPage - 2, currentPage - 1}
		nextPages = []int{currentPage + 1}
	}

	result := map[string][]int{
		"previousPages": previousPages,
		"currentPage":   {currentPage},
		"nextPages":     nextPages,
		"totalPages":    {totalPages},
	}

	return result
}
