package usecase

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

type AdminFindAllDownloadUsecase struct {
	DownloadRepository port.AdminFindAllDownloadRepository
}

func NewAdminFindAllDownloadUsecase(repository port.AdminFindAllDownloadRepository) *AdminFindAllDownloadUsecase {
	return &AdminFindAllDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *AdminFindAllDownloadUsecase) AdminFindAll(page, limit int) (interface{}, error) {

	if limit > perPage {
		limit = perPage
	}

	downloads, err := d.DownloadRepository.AdminFindAll(page, limit)

	if err != nil {
		return nil, err
	}

	total, err := d.DownloadRepository.GetTotal()

	if err != nil {
		return nil, err
	}

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
