package usecase

import (
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

	downloads := d.DownloadRepository.Search(page, str_search)

	total := d.DownloadRepository.GetTotalSearch(str_search)

	pagination := utils.Pagination(page, limit, total)

	result := struct {
		Downloads  []*entity.Download `json:"downloads"`
		Pagination map[string][]int   `json:"pagination"`
	}{
		Downloads:  downloads,
		Pagination: pagination,
	}

	return result, nil
}
