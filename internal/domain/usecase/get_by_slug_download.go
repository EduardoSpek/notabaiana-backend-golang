package usecase

import (
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

type GetBySlugDownloadUsecase struct {
	DownloadRepository port.GetBySlugDownloadRepository
}

func NewGetBySlugDownloadUsecase(repository port.GetBySlugDownloadRepository) *GetBySlugDownloadUsecase {
	return &GetBySlugDownloadUsecase{
		DownloadRepository: repository,
	}
}

func (d *GetBySlugDownloadUsecase) GetBySlug(slug string) (*entity.Download, error) {

	download, err := d.DownloadRepository.GetBySlug(slug)

	if err != nil {
		return &entity.Download{}, err
	}

	return download, nil
}
