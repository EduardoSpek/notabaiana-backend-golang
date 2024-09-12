package usecase

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

var perPage = 24

type FindAllDownloadUsecase struct {
	DownloadRepository port.DownloadRepository
}

func NewFindAllDownloadUsecase(repository port.DownloadRepository) *FindAllDownloadUsecase {
	return &FindAllDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *FindAllDownloadUsecase) FindAll(page, limit int) (interface{}, error) {

	if limit > perPage {
		limit = perPage
	}

	downloads, err := d.DownloadRepository.FindAll(page, limit)

	if err != nil {
		return nil, err
	}

	total := d.DownloadRepository.GetTotalVisible()

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
