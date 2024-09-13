package usecase

import (
	"strings"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

type SearchDownloadUsecase struct {
	DownloadRepository port.DownloadRepository
}

func NewSearchDownloadUsecase(repository port.DownloadRepository) *SearchDownloadUsecase {
	return &SearchDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *SearchDownloadUsecase) Search(page int, str_search string) (interface{}, error) {

	limit := 24

	str_search = strings.Replace(str_search, " ", "%", -1)

	downloads, err := d.DownloadRepository.Search(page, str_search)

	if err != nil {
		return nil, err
	}

	total, err := d.DownloadRepository.GetTotalSearch(str_search)

	if err != nil {
		return nil, err
	}

	pagination := utils.Pagination(page, limit, total)

	result := struct {
		Downloads  []*entity.Download `json:"downloads"`
		Pagination map[string][]int   `json:"pagination"`
		Search     string             `json:"search"`
	}{
		Downloads:  downloads,
		Pagination: pagination,
		Search:     strings.Replace(str_search, "%", " ", -1),
	}

	return result, nil
}
