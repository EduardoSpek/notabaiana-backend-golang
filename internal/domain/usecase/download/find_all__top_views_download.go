package usecase

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

type FindAllTopViewsDownloadUsecase struct {
	DownloadRepository port.DownloadRepository
}

func NewFindAllTopViewsDownloadUsecase(repository port.DownloadRepository) *FindAllTopViewsDownloadUsecase {
	return &FindAllTopViewsDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *FindAllTopViewsDownloadUsecase) FindAllTopViews(page, limit int) (interface{}, error) {

	if limit > perPage {
		limit = perPage
	}

	downloads, err := d.DownloadRepository.FindAllTopViews(page, limit)

	if err != nil {
		return nil, err
	}

	// total, err := d.DownloadRepository.GetTotalVisible()

	// if err != nil {
	// 	return nil, err
	// }

	// pagination := utils.Pagination(page, limit, total)

	// result := struct {
	// 	Downloads  []*entity.Download `json:"downloads"`
	// 	Pagination map[string][]int   `json:"pagination"`
	// }{
	// 	Downloads:  downloads,
	// 	Pagination: pagination,
	// }

	return downloads, nil
}
