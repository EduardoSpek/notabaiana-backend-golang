package usecase

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

type FindCategoryDownloadUsecase struct {
	DownloadRepository port.DownloadRepository
}

func NewFindCategoryDownloadUsecase(repository port.DownloadRepository) *FindCategoryDownloadUsecase {
	return &FindCategoryDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *FindCategoryDownloadUsecase) FindCategory(str_search string, page int) (interface{}, error) {

	limit := 24

	downloads, err := d.DownloadRepository.FindCategory(str_search, page)

	if err != nil {
		return nil, err
	}

	total, err := d.DownloadRepository.GetTotalFindCategory(str_search)

	if err != nil {
		return nil, err
	}

	pagination := utils.Pagination(page, limit, total)

	result := struct {
		Downloads  []*entity.Download `json:"downloads"`
		Pagination map[string][]int   `json:"pagination"`
		Category   string             `json:"category"`
	}{
		Downloads:  downloads,
		Pagination: pagination,
		Category:   str_search,
	}

	return result, nil
}
